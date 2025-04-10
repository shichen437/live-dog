package bilibili

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/live-dog/internal/pkg/media_parser"
	"github.com/shichen437/live-dog/internal/pkg/params"
	"github.com/shichen437/live-dog/internal/pkg/utils"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
)

type BilibiliParser struct {
	Url string
}

type builder struct{}

func (b *builder) Build(url string) (media_parser.MeidaParser, error) {
	return &BilibiliParser{
		Url: url,
	}, nil
}

func init() {
	media_parser.Register(platform, &builder{})
}

func (b *BilibiliParser) ParseURL(ctx context.Context) (*media_parser.MediaInfo, error) {
	bvid := strings.TrimSuffix(strings.TrimPrefix(utils.FindFirstMatch(b.Url, `video/BV[^?]+`), "video/"), "/")
	if bvid == "" {
		return nil, gerror.New("不支持的B站链接")
	}
	detail := b.getAidAndCid(ctx, bvid)
	if detail == nil {
		return nil, gerror.New("解析视频详情失败")
	}
	mediaData, err := b.getPlayUrlInfo(ctx, detail, bvid)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("解析视频链接失败")
	}
	mJsonData, err := json.Marshal(mediaData)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("格式化视频数据失败")
	}
	return &media_parser.MediaInfo{
		Platform:      platform,
		Refer:         b.Url,
		VideoID:       bvid,
		Author:        detail.Author,
		AuthorUid:     detail.AuthorUid,
		Desc:          detail.Desc,
		Type:          "video",
		VideoCoverUrl: detail.VideoCover,
		VideoData:     string(mJsonData),
	}, nil
}

func (b *BilibiliParser) ParseUserInfo(ctx context.Context) (*media_parser.UserInfo, error) {
	var uid string
	if strings.Contains(b.Url, "space.bilibili.com") {
		uid = strings.TrimPrefix(utils.FindFirstMatch(b.Url, `space\.bilibili\.com/(\d+)`), "space.bilibili.com/")
	}
	if uid == "" {
		return nil, gerror.New("不支持的B站链接")
	}
	userInfo, err := b.getUserStatInfo(ctx, uid)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("获取用户信息失败")
	}
	accessId, err := b.getAccessId(ctx)
	if err != nil || accessId == "" {
		g.Log().Error(ctx, err)
		return nil, gerror.New("获取access_id失败")
	}
	reqUrl, err := params.WbiSignURL(userProfileInfoUrl, g.MapStrStr{
		"mid":     uid,
		"w_webid": accessId,
	})
	if err != nil || reqUrl == "" {
		g.Log().Error(ctx, err)
		return nil, gerror.New("获取用户信息失败")
	}
	c := g.Client()
	c.SetHeaderMap(bilibiliHeaders)
	c.SetCookieMap(utils.GetCookieMap(platform, b.Url))
	resp, err := c.Get(ctx, reqUrl)
	defer resp.Close()
	if err != nil || resp.StatusCode != 200 {
		return nil, gerror.New("请求用户信息失败")
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return nil, gerror.New("请求用户信息失败")
	}
	jsonData := gjson.Parse(body)
	if jsonData.Get("code").Int() != 0 {
		return nil, gerror.New("请求用户信息失败")
	}
	userInfo.Nickname = jsonData.Get("data.name").String()
	userInfo.Avatar = jsonData.Get("data.face").String()
	userInfo.Signature = jsonData.Get("data.sign").String()
	return userInfo, nil
}

func (b *BilibiliParser) getPlayUrlInfo(ctx context.Context, detail *BilibiliVideoDetail, bvid string) (*BilibiliMediaData, error) {
	c := g.Client()
	c.SetHeaderMap(bilibiliHeaders)
	c.SetCookieMap(utils.GetCookieMap(platform, b.Url))
	payloadMap := g.MapStrStr{
		"avid":  detail.Aid,
		"bvid":  bvid,
		"cid":   detail.Cid,
		"qn":    "127",
		"fnval": "4048",
		"fourk": "1",
		"fnver": "0",
		"otype": "json",
	}
	payload := params.GetWtsParams(payloadMap)
	resp, err := c.Get(ctx, videoPlayUrl+payload)
	if err != nil || resp.StatusCode != 200 {
		return nil, gerror.New("请求视频链接失败")
	}
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	if jsonData.Get("code").Int() != 0 {
		return nil, gerror.New("请求视频链接失败")
	}
	if !jsonData.Get("data.dash.video").Exists() {
		return nil, gerror.New("请求视频链接失败")
	}
	var videos, audios []*BilibiliMediaDetail
	videoData := jsonData.Get("data.dash.video").Array()
	for _, v := range videoData {
		var mirrors []string
		if v.Get("backup_url").Exists() {
			for _, mirror := range v.Get("backup_url").Array() {
				mirrors = append(mirrors, mirror.String())
			}
		}
		videos = append(videos, &BilibiliMediaDetail{
			Url:         v.Get("base_url").String(),
			Mirrors:     mirrors,
			Width:       int(v.Get("width").Int()),
			Height:      int(v.Get("height").Int()),
			Quality:     int(v.Get("id").Int()),
			Codec:       bilibiliCodecMap[int(v.Get("codecid").Int())],
			QualityDesc: getQualityDesc(v),
		})
	}
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Quality > videos[j].Quality
	})
	audioData := jsonData.Get("data.dash.audio").Array()
	for _, v := range audioData {
		var mirrors []string
		if v.Get("backup_url").Exists() {
			for _, mirror := range v.Get("backup_url").Array() {
				mirrors = append(mirrors, mirror.String())
			}
		}
		audios = append(audios, &BilibiliMediaDetail{
			Url:     v.Get("base_url").String(),
			Mirrors: mirrors,
			Width:   int(v.Get("width").Int()),
			Height:  int(v.Get("height").Int()),
			Quality: int(v.Get("id").Int()),
			Codec:   "mp4a",
		})
	}
	sort.Slice(audios, func(i, j int) bool {
		return audios[i].Quality > audios[j].Quality
	})
	mediaData := &BilibiliMediaData{
		Videos: videos,
		Audios: audios,
	}
	return mediaData, nil
}

func (b *BilibiliParser) getAidAndCid(ctx context.Context, bvid string) *BilibiliVideoDetail {
	c := g.Client()
	c.SetHeaderMap(bilibiliHeaders)
	c.SetCookieMap(utils.GetCookieMap(platform, b.Url))
	payload := "bvid=" + bvid
	resp, err := c.Get(ctx, videoDetailPath+payload)
	if err != nil || resp.StatusCode != 200 {
		return nil
	}
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	if jsonData.Get("code").Int() != 0 {
		return nil
	}
	avid := jsonData.Get("data.aid").String()
	pageData := jsonData.Get("data.pages").Array()
	var cid string
	if len(pageData) > 0 {
		cid = pageData[0].Get("cid").String()
	}
	if avid == "" || cid == "" {
		return nil
	}
	return &BilibiliVideoDetail{
		Aid:        avid,
		Cid:        cid,
		Author:     jsonData.Get("data.owner.name").String(),
		AuthorUid:  jsonData.Get("data.owner.mid").String(),
		Desc:       jsonData.Get("data.title").String(),
		VideoCover: jsonData.Get("data.pic").String(),
	}
}

func (b *BilibiliParser) getUserStatInfo(ctx context.Context, uid string) (*media_parser.UserInfo, error) {
	c := g.Client()
	c.SetHeaderMap(bilibiliHeaders)
	c.SetCookieMap(utils.GetCookieMap(platform, b.Url))
	payload := "vmid=" + uid
	resp, err := c.Get(ctx, userProfileStatUrl+payload)
	defer resp.Close()
	if err != nil {
		return nil, gerror.New("请求用户信息失败")
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return nil, gerror.New("请求用户信息失败")
	}
	jsonData := gjson.Parse(body)
	if jsonData.Get("code").Int() != 0 {
		return nil, gerror.New("请求用户信息失败")
	}
	info := jsonData.Get("data")
	userInfo := &media_parser.UserInfo{
		UniqueId:       uid,
		FollowerCount:  int(info.Get("follower").Int()),
		FollowingCount: int(info.Get("following").Int()),
		Platform:       platform,
		Refer:          b.Url,
	}
	return userInfo, nil
}

func getQualityDesc(v gjson.Result) string {
	return "[" + bilibiliCodecMap[int(v.Get("codecid").Int())] + "][" + v.Get("width").String() + "x" + v.Get("height").String() + "]"
}

func (b *BilibiliParser) getAccessId(ctx context.Context) (string, error) {
	c := g.Client()
	c.SetHeaderMap(bilibiliHeaders)
	c.SetCookieMap(utils.GetCookieMap(platform, b.Url))
	resp, err := c.Get(ctx, b.Url)
	if err != nil {
		return "", gerror.New("请求用户信息失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", gerror.New("请求用户信息失败")
	}

	// 解析HTML
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", gerror.New("failed to parse HTML")
	}

	// 查找__RENDER_DATA__脚本内容
	var renderDataScript string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "__RENDER_DATA__" {
					if n.FirstChild != nil {
						renderDataScript = n.FirstChild.Data
					}
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if renderDataScript == "" {
		return "", gerror.New("failed to find render data script")
	}

	// URL解码并解析JSON
	decoded, err := url.QueryUnescape(renderDataScript)
	if err != nil {
		return "", gerror.New("failed to decode render data")
	}
	var renderData map[string]string
	if err := json.Unmarshal([]byte(decoded), &renderData); err != nil {
		return "", gerror.New("failed to unmarshal render data")
	}
	if _, ok := renderData["access_id"]; !ok {
		return "", gerror.New("failed to find access_id")
	}
	return renderData["access_id"], nil
}
