package cmd

import (
	"context"
	"os"

	"github.com/shichen437/live-dog/internal/app/live/service"
	"github.com/shichen437/live-dog/internal/pkg/crons"
	"github.com/shichen437/live-dog/internal/pkg/events"
	"github.com/shichen437/live-dog/internal/pkg/listeners"
	"github.com/shichen437/live-dog/internal/pkg/lives"
	"github.com/shichen437/live-dog/internal/pkg/recorders"
	"github.com/shichen437/live-dog/internal/pkg/utils"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
)

func LiveMonitor() error {
	ctx := gctx.GetInitCtx()
	g.Log().Info(ctx, "live monitor start!")
	models, _ := service.LiveManage().GetLiveModels4Init(ctx)
	gLiveModel := new(lives.GLiveModel)
	gLiveModel.Lives = make(map[int]lives.Live)
	gLiveModel.StartTimeMap = make(map[int]*gtime.Time)
	gLiveModel.ModelsMap = make(map[int]*lives.LiveModel)
	for _, m := range models {
		u, err := m.ParseUrl()
		if err != nil {
			g.Log().Error(ctx, "failed to parse url: "+m.LiveManage.RoomUrl)
			continue
		}
		l, _ := lives.New(u, m.GetLiveID())
		gLiveModel.Lives[m.GetLiveID()] = l
		gLiveModel.ModelsMap[m.GetLiveID()] = m
	}
	gLiveModel.CookieMap = make(map[string]string)
	cookies, _ := service.LiveCookie().GetAllCookie4Init(ctx)
	for _, c := range cookies {
		gLiveModel.CookieMap[c.Platform] = c.Cookie
	}
	ctx = context.WithValue(ctx, utils.Key, gLiveModel)
	events.NewDispatcher(ctx)
	lm := listeners.NewManager(ctx)
	rm := recorders.NewManager(ctx)
	gctx.SetInitCtx(ctx)
	if err := lm.Start(ctx); err != nil {
		g.Log().Error(ctx, "failed to start listener manager")
	}
	if err := rm.Start(ctx); err != nil {
		g.Log().Error(ctx, "failed to start recorder manager")
	}
	for _, l := range gLiveModel.Lives {
		m := gLiveModel.ModelsMap[l.GetLiveId()]
		if m.LiveManage.MonitorType == 2 {
			if err := lm.AddListener(ctx, l); err != nil {
				g.Log().Error(ctx, "failed to add listener")
			}
		}
		if m.LiveManage.MonitorType == 1 {
			crons.AddCron(m.GetLiveID())
		}
	}
	go func() {
		gproc.AddSigHandlerShutdown(shutdown)
		gproc.Listen()
	}()
	return nil
}

func shutdown(sig os.Signal) {
	g.Log().Info(gctx.GetInitCtx(), "live monitor shutdown!")
	utils.CleanBrowserPool()
	global := utils.GetGlobal(gctx.GetInitCtx())
	global.ListenerManager.Close(gctx.GetInitCtx())
	global.RecorderManager.Close(gctx.GetInitCtx())
	global.WaitGroup.Wait()
}
