package media_parser_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/live-dog/internal/pkg/media_parser"
	_ "github.com/shichen437/live-dog/internal/pkg/media_parser/douyin"
)

func TestDownload(t *testing.T) {
	loader, err := media_parser.NewParser("0.28 01/02 baa:/ B@t.eO 新年新气象。2025一切顺利！# 诸事顺利 # 一切都在慢慢变好  https://v.douyin.com/i5Wk7anF/ 复制此链接，打开Dou音搜索，直接观看视频！")
	if err != nil {
		t.Error(err)
	}
	loader.ParseURL(gctx.New())
}
