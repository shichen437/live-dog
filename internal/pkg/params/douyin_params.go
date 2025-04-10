package params

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/live-dog/internal/pkg/utils"
	"golang.org/x/exp/rand"
)

var (
	randomCookieChars = "1234567890abcdef"
	ttwid_url         = "https://ttwid.bytedance.com/ttwid/union/register/"
	ttwid_data        = `{"region":"cn","aid":1768,"needFid":false,"service":"www.ixigua.com","migrate_info":{"ticket":"","source":"node"},"cbUrlProtocol":"https","union":true}`
)

func GetABogus(params, ua string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		g.Log().Info(gctx.New(), wd)
	}
	jsFilePath := filepath.Join("internal/pkg/params/js/douyin", "a_bogus.js")
	jsContent, err := os.ReadFile(jsFilePath)
	if err != nil {
		return "", gerror.New("read js file failed!")
	}

	vm := vmPool.Get().(*goja.Runtime)
	defer vmPool.Put(vm)

	// 执行 JavaScript 文件内容
	_, err = vm.RunString(string(jsContent))
	if err != nil {
		return "", gerror.New("exec js file failed!")
	}

	// 获取 generate_a_bogus 方法
	generateABogus, ok := goja.AssertFunction(vm.Get("sign_detail"))
	if !ok {
		return "", gerror.New("read js method failed!")
	}
	result, err := generateABogus(goja.Undefined(), vm.ToValue(params), vm.ToValue(ua))
	if err != nil {
		return "", gerror.New("get js param failed!")
	}
	return url.QueryEscape(result.String()), nil
}

func GetOdintt() string {
	return utils.GenRandomString(160, randomCookieChars)
}

func GetMsToken() string {
	return utils.GenRandomString(107, randomCookieChars)
}

func GetTtwid() string {
	c := g.Client()
	resp, err := c.Post(gctx.New(), ttwid_url, ttwid_data)
	if err != nil {
		return ""
	}
	return resp.GetCookie("ttwid")
}

func GetVerifyFp() string {
	e := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	t := len(e)
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	var base36 string
	for milliseconds > 0 {
		remainder := int(milliseconds % 36)
		if remainder < 10 {
			base36 = strconv.Itoa(remainder) + base36
		} else {
			base36 = string([]byte{byte('a' + remainder - 10)}) + base36
		}
		milliseconds /= 36
	}
	r := base36

	o := make([]string, 36)
	o[8], o[13], o[18], o[23] = "_", "_", "_", "_"
	o[14] = "4"

	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < 36; i++ {
		if o[i] == "" {
			n := int(rand.Float64() * float64(t))
			if i == 19 {
				n = 3&n | 8
			}
			o[i] = string(e[n])
		}
	}
	return "verify_" + r + "_" + strings.Join(o, "")
}

func GetWebID(userAgent string) string {
	api := "https://mcs.zijieapi.com/webid"
	headers := map[string]string{
		"User-Agent":   userAgent,
		"Content-Type": "application/json",
	}

	data := map[string]interface{}{
		"app_id":         6383,
		"url":            "https://www.douyin.com/",
		"user_agent":     userAgent,
		"referer":        "https://www.douyin.com/",
		"user_unique_id": "",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	c := g.Client()
	c.SetHeaderMap(headers)
	resp, err := c.Post(gctx.New(), api, string(jsonData))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	if webID, ok := result["web_id"].(string); ok {
		return webID
	}

	return ""
}

func GetProfileParamsMap(uid string) map[string]string {
	paramsMap := getBaseParams()
	paramsMap["sec_user_id"] = uid
	return paramsMap
}

func ConvertParamsToQueryString(params map[string]interface{}) string {
	queryString := ""
	for key, value := range params {
		queryString += key + "=" + fmt.Sprint(value) + "&"
	}
	return queryString[:len(queryString)-1]
}

func getBaseParams() map[string]string {
	fp := GetVerifyFp()
	return g.MapStrStr{
		"device_platform":             "webapp",
		"aid":                         "6383",
		"channel":                     "channel_pc_web",
		"pc_client_type":              "1",
		"publish_video_strategy_type": "2",
		"pc_libra_divert":             "Mac",
		"version_code":                "170400",
		"version_name":                "17.4.0",
		"cookie_enabled":              "true",
		"screen_width":                "1512",
		"screen_height":               "982",
		"browser_language":            "zh-CN",
		"browser_platform":            "MacIntel",
		"browser_name":                "Chrome",
		"browser_version":             "123.0.0.0",
		"browser_online":              "true",
		"engine_name":                 "Blink",
		"engine_version":              "123.0.0.0",
		"os_name":                     "Mac OS",
		"os_version":                  "10.15.7",
		"cpu_core_num":                "10",
		"device_memory":               "8",
		"platform":                    "PC",
		"downlink":                    "10",
		"effective_type":              "4g",
		"round_trip_time":             "100",
		"fp":                          fp,
		"verifyFp":                    fp,
	}
}
