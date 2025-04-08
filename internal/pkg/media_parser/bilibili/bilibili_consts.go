package bilibili

import "github.com/gogf/gf/v2/frame/g"

var (
	platform = "bilibili"

	videoDetailPath = "https://api.bilibili.com/x/web-interface/view?"
	videoPlayUrl    = "https://api.bilibili.com/x/player/playurl?"

	userProfileInfoUrl = "https://api.bilibili.com/x/space/wbi/acc/info?"
	userProfileStatUrl = "https://api.bilibili.com/x/relation/stat?"

	bilibiliHeaders = g.MapStrStr{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Referer":         "https://space.bilibili.com/",
		"Origin":          "https://www.bilibili.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
	}

	bilibiliCodecMap = g.MapIntStr{
		7:  "avc",
		12: "hevc",
		13: "av1",
	}
)
