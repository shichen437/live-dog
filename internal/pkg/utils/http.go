package utils

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

// GetClientIp 获取客户端IP
func GetClientIp(ctx context.Context) string {
	return g.RequestFromCtx(ctx).GetClientIp()
}

// GetUserAgent 获取user-agent
func GetUserAgent(ctx context.Context) string {
	return ghttp.RequestFromCtx(ctx).Header.Get("User-Agent")
}

func Text(r *http.Response) (string, error) {
	if r.Body == nil {
		return "", io.EOF
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			g.Log().Error(gctx.GetInitCtx(), "Error closing response body:", err)
		}
	}()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GetCookieList(platform string) []*http.Cookie {
	cookiesList := make([]*http.Cookie, 0)
	cookie := GetGlobal(gctx.GetInitCtx()).CookieMap[platform]
	if cookie == "" {
		return cookiesList
	}
	for _, cStr := range strings.Split(cookie, ";") {
		cArr := strings.SplitN(cStr, "=", 2)
		if len(cArr) != 2 {
			continue
		}
		cookiesList = append(cookiesList, &http.Cookie{
			Name:  strings.TrimSpace(cArr[0]),
			Value: strings.TrimSpace(cArr[1]),
		})
	}
	return cookiesList
}

func GetCookieMap(platform, refer string) map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	url, _ := url.Parse(refer)
	jar.SetCookies(url, GetCookieList(platform))
	cookies := jar.Cookies(url)
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

type DownloadFileRequest struct {
	CookieMap map[string]string
	Filename  string
	Refer     string
	Url       string
	UserAgent string
}

func GetDownloadFile(ctx context.Context, req *DownloadFileRequest) (int64, error) {
	client := g.Client()
	if req.UserAgent != "" {
		client.SetAgent(req.UserAgent)
	} else {
		client.SetAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	}
	if req.CookieMap != nil {
		client.SetCookieMap(req.CookieMap)
	}
	if req.Refer != "" {
		client.SetHeader("Referer", req.Refer)
	}
	resp, err := client.Get(ctx, req.Url)
	defer resp.Close()
	if err != nil {
		g.Log().Error(ctx, err)
		return 0, gerror.New("HTTP请求失败")
	}
	if resp.StatusCode != 200 {
		g.Log().Error(ctx, "HTTP请求失败，状态码:", resp.StatusCode)
		return 0, gerror.New("HTTP请求失败")
	}
	if resp.ContentLength == 0 {
		return 0, gerror.New("文件大小为0")
	}
	if req.Filename != "" {
		gfile.PutBytes(req.Filename, resp.ReadAll())
	}
	return resp.ContentLength, nil
}
