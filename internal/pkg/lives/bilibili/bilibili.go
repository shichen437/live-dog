package bilibili

import (
	"fmt"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/live-dog/internal/pkg/lives"
	"github.com/shichen437/live-dog/internal/pkg/utils"
	"github.com/tidwall/gjson"
)

const (
	domain       = "live.bilibili.com"
	platform     = "bilibili"
	userAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
	roomInitUrl  = "https://api.live.bilibili.com/room/v1/Room/room_init"
	roomApiUrl   = "https://api.live.bilibili.com/room/v1/Room/get_info"
	userApiUrl   = "https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room"
	liveApiUrlv2 = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo"
	parseAddr    = "data.playurl_info.playurl.stream.0.format.0.codec.0"
)

func init() {
	lives.Register(domain, &builder{})
}

type builder struct{}

func (b *builder) Build(url *url.URL, liveId int) (lives.Live, error) {
	return &Bilibili{
		Url:         url,
		LiveId:      liveId,
		Platform:    platform,
		RespCookies: make(map[string]string),
	}, nil
}

type Bilibili struct {
	Url         *url.URL
	LiveId      int
	Platform    string
	RoomID      string
	RespCookies map[string]string
}

func (l *Bilibili) GetLiveId() int {
	return l.LiveId
}

func (l *Bilibili) GetPlatform() string {
	return l.Platform
}

func (l *Bilibili) GetRefer() string {
	return l.Url.String()
}

func (l *Bilibili) GetInfo() (info *lives.RoomInfo, err error) {
	if l.RoomID == "" && l.parseRoomID() != nil {
		return nil, gerror.New("parse room id failed")
	}
	info, err = l.getRoomInfo()
	if err != nil {
		return nil, err
	}
	err = l.getUserInfo(info)
	if err != nil {
		return nil, err
	}
	if info.LiveStatus {
		streamInfos, err := l.getStreamInfo()
		if err != nil || len(streamInfos) == 0 {
			return nil, gerror.New("get stream info failed")
		}
		info.StreamInfos = streamInfos
	}
	return info, nil
}

func (l *Bilibili) getStreamInfo() (infos []*lives.StreamUrlInfo, err error) {
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	payload := fmt.Sprintf(`?room_id=%s&protocol=0,1&format=0,1,2&codec=0,1&qn=10000&platform=web&ptype=8&dolby=5&panorama=1`, l.RoomID)
	resp, err := c.Get(gctx.GetInitCtx(), liveApiUrlv2+payload)
	defer resp.Close()
	if err != nil || resp.StatusCode != 200 {
		return nil, gerror.New("get stream info failed")
	}
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	baseUrl := jsonData.Get(parseAddr + ".base_url").String()
	jsonArr := jsonData.Get(parseAddr + ".url_info").Array()
	streamUrlInfos := make([]*lives.StreamUrlInfo, 0, 10)
	for _, v := range jsonArr {
		hosts := gjson.Get(v.String(), "host").String()
		queries := gjson.Get(v.String(), "extra").String()
		streamUrl, err := url.Parse(hosts + baseUrl + queries)
		if err != nil {
			continue
		}
		streamUrlInfos = append(streamUrlInfos, &lives.StreamUrlInfo{
			Url:                  streamUrl,
			HeadersForDownloader: l.getHeadersForDownloader(),
		})
	}
	return streamUrlInfos, nil
}

func (l *Bilibili) getUserInfo(info *lives.RoomInfo) error {
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	resp, err := c.Get(gctx.GetInitCtx(), userApiUrl, g.Map{
		"roomid": l.RoomID,
	})
	defer resp.Close()
	if err != nil || resp.StatusCode != 200 {
		return gerror.New("get room info failed")
	}
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	if jsonData.String() == "" || jsonData.Get("code").Int() != 0 {
		return gerror.New("get room info failed")
	}
	info.Anchor = jsonData.Get("data.info.uname").String()
	return nil
}

func (l *Bilibili) getRoomInfo() (*lives.RoomInfo, error) {
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	resp, err := c.Get(gctx.GetInitCtx(), roomApiUrl, g.Map{
		"room_id": l.RoomID,
		"from":    "room",
	})
	defer resp.Close()
	if err != nil || resp.StatusCode != 200 {
		return nil, gerror.New("get room info failed")
	}
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	if jsonData.String() == "" || jsonData.Get("code").Int() != 0 {
		return nil, gerror.New("get room info failed")
	}
	info := &lives.RoomInfo{
		Platform:   l.Platform,
		RoomName:   jsonData.Get("data.title").String(),
		LiveStatus: jsonData.Get("data.live_status").Int() == 1,
	}
	return info, nil
}

func (l *Bilibili) parseRoomID() error {
	paths := strings.Split(l.Url.Path, "/")
	if len(paths) < 2 {
		return gerror.New("wrong url")
	}
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	resp, err := c.Get(gctx.GetInitCtx(), roomInitUrl, g.Map{
		"id": paths[1],
	})
	if err != nil || resp.StatusCode != 200 {
		return gerror.New("get room info failed")
	}
	body, err := utils.Text(resp.Response)
	fmt.Println("Response Body:", body)
	jsonData := gjson.Parse(body)
	if jsonData.String() == "" || jsonData.Get("code").Int() != 0 {
		return gerror.New("get room info failed")
	}
	l.RoomID = jsonData.Get("data.room_id").String()
	return nil
}

func (l *Bilibili) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	jar.SetCookies(l.Url, utils.GetCookieList(platform))
	cookies := jar.Cookies(l.Url)
	cookieMap := make(map[string]string)
	for k, v := range l.RespCookies {
		cookieMap[k] = v
	}
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func (l *Bilibili) getHeadersForDownloader() map[string]string {
	return map[string]string{
		"User-Agent": userAgent,
		"Referer":    l.Url.String(),
	}
}
