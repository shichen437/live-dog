package douyin

var (
	platform          = "douyin"
	domain            = "https://www.douyin.com"
	randomCookieChars = "1234567890abcdef"
	videoDetailPath   = "https://www.iesdouyin.com/share/video/"
	userProfilePath   = "https://www.douyin.com/aweme/v1/web/user/profile/other/?"
	douyinHeaders     = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Referer":         "https://www.douyin.com",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
	}
)
