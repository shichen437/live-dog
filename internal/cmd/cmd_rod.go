package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

func InitBrowserPool() {
	g.Log().Info(gctx.New(), "初始化浏览器池")
	utils.InitBrowserPool(1)
	g.Log().Info(gctx.New(), "初始化浏览器池完成")
}
