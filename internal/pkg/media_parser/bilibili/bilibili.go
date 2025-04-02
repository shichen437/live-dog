package bilibili

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/live-dog/internal/pkg/media_parser"
	"github.com/shichen437/live-dog/internal/pkg/params"
	"github.com/shichen437/live-dog/internal/pkg/utils"
	"github.com/tidwall/gjson"
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

func getQualityDesc(v gjson.Result) string {
	return "[" + bilibiliCodecMap[int(v.Get("codecid").Int())] + "][" + v.Get("width").String() + "x" + v.Get("height").String() + "]"
}
