package params

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/live-dog/internal/pkg/utils"
	"github.com/tidwall/gjson"
)

var (
	mixinKeyEncTab = []int{
		46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
		33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
		61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
		36, 20, 34, 44, 52,
	}
)

func GetWtsParams(params map[string]string) string {
	wts := params["wts"]
	encodeQuery := getEncodeQuery(params)
	wRid := getWrid(encodeQuery)
	params["wts"] = wts
	params["w_rid"] = wRid

	var parts []string
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(parts, "&")
}

func WbiSignURL(urlStr string, params map[string]string) (string, error) {
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	imgKey, subKey := getWbiKeys()
	query := urlObj.Query()
	newParams := encWbi(params, imgKey, subKey)
	for k, v := range newParams {
		query.Set(k, v)
	}
	urlObj.RawQuery = query.Encode()
	newUrlStr := urlObj.String()
	return newUrlStr, nil
}

func getEncodeQuery(params map[string]string) string {
	params["wts"] = params["wts"] + "ea1db124af3c7062474693fa704f4ff8"

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	filteredParams := make(map[string]string)
	for _, k := range keys {
		v := params[k]
		filteredValue := filterSpecialChars(v)
		filteredParams[k] = filteredValue
	}

	values := url.Values{}
	for k, v := range filteredParams {
		values.Add(k, v)
	}

	return values.Encode()
}

func filterSpecialChars(s string) string {
	specialChars := "!'()*"
	result := ""
	for _, c := range s {
		if !strings.ContainsRune(specialChars, c) {
			result += string(c)
		}
	}
	return result
}

func getWrid(e string) string {
	hash := md5.Sum([]byte(e))
	return hex.EncodeToString(hash[:])
}

func encWbi(params map[string]string, imgKey, subKey string) map[string]string {
	mixinKey := getMixinKey(imgKey + subKey)
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	params["wts"] = currTime

	// Sort keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Remove unwanted characters
	for k, v := range params {
		v = sanitizeString(v)
		params[k] = v
	}

	// Build URL parameters
	query := url.Values{}
	for _, k := range keys {
		query.Set(k, params[k])
	}
	queryStr := query.Encode()

	// Calculate w_rid
	hash := md5.Sum([]byte(queryStr + mixinKey))
	params["w_rid"] = hex.EncodeToString(hash[:])
	return params
}

func getMixinKey(orig string) string {
	var str strings.Builder
	for _, v := range mixinKeyEncTab {
		if v < len(orig) {
			str.WriteByte(orig[v])
		}
	}
	return str.String()[:32]
}

func sanitizeString(s string) string {
	unwantedChars := []string{"!", "'", "(", ")", "*"}
	for _, char := range unwantedChars {
		s = strings.ReplaceAll(s, char, "")
	}
	return s
}

func getWbiKeys() (string, string) {
	c := g.Client()
	c.SetCookieMap(utils.GetCookieMap("bilibili", "https://www.bilibili.com"))
	re, err := c.Get(gctx.GetInitCtx(), "https://api.bilibili.com/x/web-interface/nav")
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), err)
		return "", ""
	}
	req := re.Request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	resp, err := c.Do(req)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), err)
		return "", ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), err)
		return "", ""
	}
	json := string(body)
	imgURL := gjson.Get(json, "data.wbi_img.img_url").String()
	subURL := gjson.Get(json, "data.wbi_img.sub_url").String()
	imgKey := strings.Split(strings.Split(imgURL, "/")[len(strings.Split(imgURL, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(subURL, "/")[len(strings.Split(subURL, "/"))-1], ".")[0]
	return imgKey, subKey
}
