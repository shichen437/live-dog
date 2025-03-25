package douyin

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/shichen437/live-dog/internal/pkg/lives"
	"github.com/shichen437/live-dog/internal/pkg/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/tidwall/gjson"
)

const (
	domain                     = "live.douyin.com"
	platform                   = "douyin"
	randomCookieChars          = "1234567890abcdef"
	mainInfoLineCatcherRegex   = `self.__pace_f.push\(\[1,\s*"[^:]*:([^<]*,null,\{\\"state\\"[^<]*\])\\n"\]\)`
	commonInfoLineCatcherRegex = `self.__pace_f.push\(\[1,\s*\"(\{.*\})\"\]\)`
)

func init() {
	lives.Register(domain, &builder{})
}

type builder struct{}

func (b *builder) Build(url *url.URL, liveId int) (lives.Live, error) {
	return &Douyin{
		Url:         url,
		LiveId:      liveId,
		Platform:    platform,
		RespCookies: make(map[string]string),
	}, nil
}

type Douyin struct {
	Url         *url.URL
	LiveId      int
	Platform    string
	RespCookies map[string]string
}

func (l *Douyin) GetLiveId() int {
	return l.LiveId
}

func (l *Douyin) GetPlatform() string {
	return l.Platform
}

func (l *Douyin) GetRefer() string {
	return l.Url.String()
}

func (l *Douyin) GetInfo() (info *lives.RoomInfo, err error) {
	info = &lives.RoomInfo{}
	body, err := l.getRoomWebPageResp()
	if err != nil {
		err = gerror.New("get room web page failed")
		return
	}
	json, err := getMainInfoLine(body)
	if err != nil {
		return
	}
	info.RoomName = json.Get("state.roomStore.roomInfo.room.title").String()
	info.Anchor = json.Get("state.roomStore.roomInfo.anchor.nickname").String()
	info.Platform = l.Platform
	isStreaming := json.Get("state.roomStore.roomInfo.room.status_str").String() == "2"
	info.LiveStatus = isStreaming
	if !isStreaming {
		return
	}
	streamIdPath := "state.streamStore.streamData.H264_streamData.common.stream"
	streamId := json.Get(streamIdPath).String()
	streamUrlInfos, _ := getStreamInfo(body, streamId)
	info.StreamInfos = streamUrlInfos
	return
}

func getStreamInfo(body string, streamId string) (infos []*lives.StreamUrlInfo, err error) {
	streamUrlInfos := make([]*lives.StreamUrlInfo, 0, 10)
	reg2, _ := regexp.Compile(commonInfoLineCatcherRegex)
	match2 := reg2.FindAllStringSubmatch(body, -1)
	for _, item := range match2 {
		if len(item) < 2 {
			return
		}
		commonJson := gjson.Parse(gjson.Parse(fmt.Sprintf(`"%s"`, item[1])).String())
		if !commonJson.Exists() {
			return
		}
		if !commonJson.Get("common").Exists() {
			continue
		}
		commonStreamId := commonJson.Get("common.stream").String()
		if commonStreamId == "" {
			return
		}
		if commonStreamId != streamId {
			continue
		}
		commonJson.Get("data").ForEach(func(key, value gjson.Result) bool {
			flv := value.Get("main.flv").String()
			var Url *url.URL
			Url, err = url.Parse(flv)
			if err != nil {
				err = gerror.New("invalid url")
				return false
			}
			paramsString := value.Get("main.sdk_params").String()
			paramsJson := gjson.Parse(paramsString)
			var description strings.Builder
			paramsJson.ForEach(func(key, value gjson.Result) bool {
				description.WriteString(key.String())
				description.WriteString(": ")
				description.WriteString(value.String())
				description.WriteString("\n")
				return true
			})
			Resolution := 0
			resolution := strings.Split(paramsJson.Get("resolution").String(), "x")
			if len(resolution) == 2 {
				x, err := strconv.Atoi(resolution[0])
				if err != nil {
					return true
				}
				y, err := strconv.Atoi(resolution[1])
				if err != nil {
					return true
				}
				Resolution = x * y
			}
			Vbitrate := int(paramsJson.Get("vbitrate").Int())
			streamUrlInfos = append(streamUrlInfos, &lives.StreamUrlInfo{
				Name:        key.String(),
				Description: description.String(),
				Url:         Url,
				Resolution:  Resolution,
				Vbitrate:    Vbitrate,
			})
			return true
		})
		sort.Slice(streamUrlInfos, func(i, j int) bool {
			if streamUrlInfos[i].Resolution != streamUrlInfos[j].Resolution {
				return streamUrlInfos[i].Resolution > streamUrlInfos[j].Resolution
			} else {
				return streamUrlInfos[i].Vbitrate > streamUrlInfos[j].Vbitrate
			}
		})
	}
	infos = streamUrlInfos
	return
}

func (l *Douyin) getRoomWebPageResp() (body string, err error) {
	c := g.Client()
	cookieMap := l.assembleCookieMap()
	c.SetCookieMap(cookieMap)
	req, err := c.Get(gctx.GetInitCtx(), l.Url.String())
	g.Log().Info(gctx.GetInitCtx(), "Get Room Web Page: "+l.Url.String())
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), err.Error())
		return
	}
	cookieWithOdinTt := fmt.Sprintf("odin_tt=%s; %s", utils.GenRandomString(160, randomCookieChars), req.Header.Get("Cookie"))
	req.Header.Set("Cookie", cookieWithOdinTt)
	c2 := g.Client()
	resp, err := c2.Do(req.Request)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), err.Error())
		return
	}
	switch code := resp.StatusCode; code {
	case http.StatusOK:
		body, err = utils.Text(resp)
		for _, cookie := range resp.Cookies() {
			l.RespCookies[cookie.Name] = cookie.Value
		}
		return
	default:
		err = gerror.Newf(`http response error`)
		return
	}
}

func (l *Douyin) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	jar.SetCookies(l.Url, utils.GetCookieList(platform))
	cookies := jar.Cookies(l.Url)
	cookieMap := make(map[string]string)
	cookieMap["__ac_nonce"] = utils.GenRandomString(21, randomCookieChars)
	for k, v := range l.RespCookies {
		cookieMap[k] = v
	}
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func getMainInfoLine(body string) (json *gjson.Result, err error) {
	reg, err := regexp.Compile(mainInfoLineCatcherRegex)
	if err != nil {
		return
	}
	match := reg.FindAllStringSubmatch(body, -1)
	if match == nil {
		err = fmt.Errorf("0 match for mainInfoLineCatcherRegex: %s", mainInfoLineCatcherRegex)
		return
	}
	for _, item := range match {
		if len(item) < 2 {
			continue
		}
		mainInfoLine := item[1]

		// 获取房间信息
		mainJson := gjson.Parse(fmt.Sprintf(`"%s"`, mainInfoLine))
		if !mainJson.Exists() {
			continue
		}

		mainJson = gjson.Parse(mainJson.String()).Get("3")
		if !mainJson.Exists() {
			continue
		}

		if mainJson.Get("state.roomStore.roomInfo.room.status_str").Exists() {
			json = &mainJson
			return
		}
	}
	return nil, fmt.Errorf("MainInfoLine not found")
}
