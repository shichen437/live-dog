package douyin

import (
	"bytes"
	"context"
	"fmt"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/live-dog/internal/pkg/media_parser"
	"github.com/shichen437/live-dog/internal/pkg/utils"
	"github.com/tidwall/gjson"
)

type DouyinParser struct {
	Url string
}

type builder struct{}

func (b *builder) Build(url string) (media_parser.MeidaParser, error) {
	return &DouyinParser{
		Url: url,
	}, nil
}

func init() {
	media_parser.Register(platform, &builder{})
}

func (d *DouyinParser) ParseURL(ctx context.Context) (*media_parser.MediaInfo, error) {
	// 提取网址
	mediaUrl := d.getMediaUrl(ctx)
	if mediaUrl == "" {
		return nil, gerror.New("请求链接失败, 请检查链接或Cookie")
	}
	// 视频
	if strings.Contains(mediaUrl, "/video/") {
		videoId := strings.TrimPrefix(utils.FindFirstMatch(mediaUrl, `video/(\d+)?`), "video/")
		videoInfo, err := getVideoInfo(ctx, videoId)
		if err != nil {
			return nil, err
		}
		return videoInfo, nil
	}
	// 图集
	if strings.Contains(mediaUrl, "/note/") {
		videoId := strings.TrimPrefix(utils.FindFirstMatch(mediaUrl, `note/(\d+)?`), "note/")
		videoInfo, err := getVideoInfo(ctx, videoId)
		if err != nil {
			return nil, err
		}
		return videoInfo, nil
	}
	// 主页
	if strings.Contains(mediaUrl, "/user/") {
		return nil, gerror.New("请跳转博主管理页面")
	}
	return nil, gerror.New("不支持的抖音链接")
}

func (d *DouyinParser) ParseUserInfo(ctx context.Context) (*media_parser.UserInfo, error) {
	if !strings.Contains(d.Url, "/user/") ||
		!strings.Contains(d.getMediaUrl(ctx), "/user/") {
		return nil, gerror.New("不支持的抖音链接")
	}
	userInfo, err := getUserInfo(d.Url)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (d *DouyinParser) getMediaUrl(ctx context.Context) string {
	c := g.Client()
	c.SetHeaderMap(douyinHeaders)
	c.SetCookieMap(d.assembleCookieMap())
	resp, err := c.Get(ctx, d.Url)
	if err != nil {
		g.Log().Error(ctx, err)
		return ""
	}
	mediaUrl := resp.Response.Request.URL.String()
	d.Url = mediaUrl
	return mediaUrl
}

func getUserInfo(url string) (*media_parser.UserInfo, error) {
	var userInfo *media_parser.UserInfo
	var err error

	u := launcher.New().MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	router := browser.HijackRequests()

	router.MustAdd("*/user/profile/other/*", func(hj *rod.Hijack) {
		hj.MustLoadResponse()
		body := hj.Response.Payload().Body
		jsonData := gjson.ParseBytes(body)
		if !jsonData.Exists() {
			err = gerror.New("未能获取到用户信息")
			return
		}
		var avatar string
		if jsonData.Get("user.avatar_medium.url_list").Exists() {
			arr := jsonData.Get("user.avatar_medium.url_list").Array()
			avatar = arr[0].String()
		}
		userInfo = &media_parser.UserInfo{
			UniqueId:       jsonData.Get("user.sec_uid").String(),
			Platform:       platform,
			Nickname:       jsonData.Get("user.nickname").String(),
			Avatar:         avatar,
			Signature:      jsonData.Get("user.signature").String(),
			IpLocation:     jsonData.Get("user.ip_location").String(),
			FollowerCount:  int(jsonData.Get("user.follower_count").Int()),
			FollowingCount: int(jsonData.Get("user.following_count").Int()),
			Refer:          url,
		}
	})

	go router.Run()

	page := browser.MustPage(url)
	page.MustWaitLoad()
	time.Sleep(time.Second * 5)
	return userInfo, err
}

func getVideoInfo(ctx context.Context, videoId string) (*media_parser.MediaInfo, error) {
	c := g.Client()
	c.SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1 Edg/122.0.0.0")
	url := fmt.Sprintf("https://www.iesdouyin.com/share/video/%s", videoId)
	c.SetCookie("odin_tt", utils.GenRandomString(160, randomCookieChars))
	resp, err := c.Get(ctx, url)
	if err != nil {
		return nil, gerror.New("请求视频链接失败")
	}
	defer resp.Close()
	// 修正拼写错误
	re := regexp.MustCompile(`window._ROUTER_DATA\s*=\s*(.*?)</script>`)
	match := re.FindSubmatch(resp.ReadAll())
	if len(match) < 2 {
		return nil, gerror.New("未获取到有效信息")
	}
	jsonBytes := bytes.TrimSpace(match[1])
	data := gjson.GetBytes(jsonBytes, "loaderData.video_(id)/page.videoInfoRes.item_list.0")
	if !data.Exists() {
		return nil, gerror.New("未获取到有效信息")
	}
	// 获取图集图片地址
	imagesObj := data.Get("images").Array()
	images := make([]string, 0, len(imagesObj))
	for _, imageItem := range imagesObj {
		imageUrl := imageItem.Get("url_list.0").String()
		if len(imageUrl) > 0 {
			images = append(images, imageUrl)
		}
	}
	if len(images) > 0 {
		imagesInfo := &media_parser.MediaInfo{
			Platform:       platform,
			VideoID:        videoId,
			Author:         data.Get("author.nickname").String(),
			AuthorUid:      data.Get("author.sec_uid").String(),
			Desc:           data.Get("desc").String(),
			Type:           "note",
			ImagesUrl:      strings.Join(images, ","),
			ImagesCoverUrl: data.Get("video.cover.url_list.0").String(),
		}
		return imagesInfo, nil
	}
	// 获取视频播放地址
	videoInfo := &media_parser.MediaInfo{
		Platform:      platform,
		VideoID:       videoId,
		Author:        data.Get("author.nickname").String(),
		AuthorUid:     data.Get("author.sec_uid").String(),
		Desc:          data.Get("desc").String(),
		Type:          "video",
		VideoUrl:      strings.ReplaceAll(data.Get("video.play_addr.url_list.0").String(), "playwm", "play"),
		VideoCoverUrl: data.Get("video.cover.url_list.0").String(),
	}
	return videoInfo, nil
}

func (d *DouyinParser) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	url, _ := url.Parse(d.Url)
	jar.SetCookies(url, utils.GetCookieList(platform))
	cookies := jar.Cookies(url)
	cookieMap := make(map[string]string)
	cookieMap["__ac_nonce"] = utils.GenRandomString(21, randomCookieChars)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}
