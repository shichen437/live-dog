package utils

import (
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	pool      rod.Pool[rod.Browser]
	poolMutex sync.Mutex // 添加互斥锁
)

func GenHeadlessBrowser() (*rod.Browser, error) {
	l := launcher.New().
		NoSandbox(true).
		Headless(true).
		Leakless(true)

	path, ok := launcher.LookPath()
	if ok && path != "" {
		g.Log().Info(gctx.New(), "使用本地浏览器", path)
		l = l.Bin(path)
	}

	u, err := l.Launch()
	if err != nil {
		return nil, gerror.New("启动浏览器失败")
	}
	g.Log().Info(gctx.New(), "浏览器启动成功", u)
	return rod.New().ControlURL(u).MustConnect(), nil
}

func InitBrowserPool(poolNum int) {
	if poolNum <= 0 {
		return
	}
	pool = rod.NewBrowserPool(poolNum)
	create := func() *rod.Browser {
		b, err := GenHeadlessBrowser()
		if err != nil {
			g.Log().Error(gctx.New(), "创建浏览器失败", err)
			return nil
		}
		return b
	}
	var wg sync.WaitGroup
	for i := 0; i < poolNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			browser := pool.MustGet(create)
			defer pool.Put(browser)
		}()
	}
	wg.Wait()
}

func GetBrowser() (*rod.Browser, error) {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	create := func() *rod.Browser {
		b, err := GenHeadlessBrowser()
		if err != nil {
			g.Log().Error(gctx.New(), "创建浏览器失败", err)
			return nil
		}
		return b
	}

	browser := pool.MustGet(create)
	if browser == nil {
		return nil, gerror.New("从浏览器池获取浏览器失败")
	}
	return browser, nil
}

func ReturnBrowser(browser *rod.Browser) {
	if pool != nil && browser != nil {
		poolMutex.Lock()
		pool.Put(browser)
		poolMutex.Unlock()
	}
}

func CleanBrowserPool() {
	if pool == nil {
		return
	}
	pool.Cleanup(func(p *rod.Browser) {
		p.MustClose()
	})
}
