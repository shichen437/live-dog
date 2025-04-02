package cmd

import (
	"context"
	"fmt"
	"os"

	admin "github.com/shichen437/live-dog/internal/app/admin/controller"
	"github.com/shichen437/live-dog/internal/app/admin/dao"
	"github.com/shichen437/live-dog/internal/app/admin/model/do"
	"github.com/shichen437/live-dog/internal/app/admin/model/entity"
	"github.com/shichen437/live-dog/internal/app/common/consts"
	common "github.com/shichen437/live-dog/internal/app/common/controller"
	"github.com/shichen437/live-dog/internal/app/common/service"
	live "github.com/shichen437/live-dog/internal/app/live/controller"
	monitor "github.com/shichen437/live-dog/internal/app/monitor/controller"
	"github.com/shichen437/live-dog/internal/pkg/sse"
	"github.com/shichen437/live-dog/internal/pkg/utils"

	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmode"
)

var (
	// Main is the main command.
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server of simple goframe demos",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.Log().Info(ctx, "Starting HTTP Server...")
			gfToken, err := GetGtoken(ctx)
			if err != nil {
				return err
			}
			s := g.Server()
			s.SetGraceful(true)
			s.SetIndexFolder(true)
			s.SetServerRoot(utils.GetFileUploadPath())
			s.AddSearchPath(utils.GetFileUploadPath())
			os.MkdirAll(utils.GetFileUploadPath()+consts.AvatarPrefix, os.ModePerm)
			s.AddStaticPath(consts.AvatarPrefix, utils.GetFileUploadPath()+consts.AvatarPrefix)
			s.Use(service.Middleware().HandlerResponse) //返回数据处理
			s.SetSwaggerUITemplate(consts.MySwaggerUITemplate)
			s.BindHandler("/sse", sse.HandleSSE)
			s.Group("/", func(group *ghttp.RouterGroup) {

				// Group middlewares.
				group.Middleware(
					service.Middleware().Ctx,
					ghttp.MiddlewareCORS,
				)
				// Register route handlers. 不需要登录
				group.Bind(
					common.Captcha,
				)
				// Special handler that needs authentication.
				group.Group("/", func(group *ghttp.RouterGroup) {
					err = gfToken.Middleware(ctx, group)
					if err != nil {
						panic(err)
					}
					group.Middleware(service.Middleware().Auth)
					bindRoute(group)
				})
			})
			LiveMonitor()
			JobInit()
			// Just run the server.
			s.Run()
			return nil
		},
	}
)

func GetGtoken(ctx context.Context) (gfToken *gtoken.GfToken, err error) {
	// 启动gtoken
	gfToken = &gtoken.GfToken{
		ServerName:       consts.ServerName,
		CacheMode:        3,
		LoginPath:        "/login",
		LoginBeforeFunc:  LoginFunc,
		LogoutPath:       "post:/logout",
		AuthPaths:        g.SliceStr{"/user", "/getInfo"}, // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{},
		AuthAfterFunc:    AuthAfterFunc,
		MultiLogin:       consts.MultiLogin,
	}
	err = gfToken.Start()
	return
}
func LoginFunc(r *ghttp.Request) (string, interface{}) {
	username := r.Get("userName").String()
	password := r.Get("password").String()
	code := r.Get("code").String()
	verifyKey := r.Get("verifyKey").String()

	//判断验证码是否正确
	//开发环境取消验证码
	debug := gmode.IsDevelop()
	if !debug {
		if code == "" || verifyKey == "" {
			r.Response.WriteJson(gtoken.Fail(utils.T(r.Context(), consts.ErrLoginCodeFailMsg)))
			r.ExitAll()
		}
		if !service.Captcha().VerifyString(verifyKey, code) {
			r.Response.WriteJson(gtoken.Fail(utils.T(r.Context(), consts.ErrLoginCodeFailMsg)))
			r.ExitAll()
		}
	}
	if username == "" || password == "" {
		r.Response.WriteJson(gtoken.Fail(utils.T(r.Context(), consts.ErrLoginFailMsg)))
		r.ExitAll()
	}
	ctx := context.TODO()
	var users *entity.SysUser
	enc, _ := utils.Encrypt(ctx, password)
	err := dao.SysUser.Ctx(ctx).Where(do.SysUser{
		UserName: username,
		Password: enc,
	}).Scan(&users)
	if err != nil || users == nil {
		r.Response.WriteJson(gtoken.Fail(utils.T(r.Context(), consts.ErrLoginFailMsg)))
		r.ExitAll()
	}
	// 唯一标识，扩展参数user data
	return fmt.Sprintf("%s%d", consts.GTokenAdminPrefix, users.UserId), users
}

func AuthAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	var users entity.SysUser
	err := gconv.Struct(respData.GetString("data"), &users)
	if err != nil {
		r.Response.WriteJson(gtoken.Unauthorized(utils.T(r.Context(), consts.ErrAuthFailMsg), nil))
		return
	}
	r.SetCtxVar(consts.CtxAdminId, users.UserId)
	r.SetCtxVar(consts.CtxAdminName, users.UserName)
	r.Middleware.Next()
}

func bindRoute(group *ghttp.RouterGroup) {
	group.Bind(
		admin.DictData,
		admin.DictType,
		admin.Menu,
		admin.Role,
		admin.SysConfig,
		admin.SysNotice,
		admin.SysUser,
		admin.PushChannel,
		monitor.OperLog,
		monitor.ServerInfo,
		monitor.SysJob,
		live.StatDaily,
		live.LiveManage,
		live.LiveCookie,
		live.FileManage,
		live.LiveHistory,
		live.MediaParse,
		live.MediaDownload,
	)
}
