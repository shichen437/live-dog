package params

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
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
