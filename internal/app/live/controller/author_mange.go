package live

import (
	"context"

	v1 "github.com/shichen437/live-dog/api/v1/live"
	"github.com/shichen437/live-dog/internal/app/live/service"
)

type authorManageController struct {
}

var AuthorManage = authorManageController{}

func (f *authorManageController) List(ctx context.Context, req *v1.GetAuthorInfoListReq) (res *v1.GetAuthorInfoListRes, err error) {
	res, err = service.AuthorManage().List(ctx, req)
	return
}

func (f *authorManageController) New(ctx context.Context, req *v1.PostAuthorInfoReq) (res *v1.PostAuthorInfoRes, err error) {
	res, err = service.AuthorManage().New(ctx, req)
	return
}

func (f *authorManageController) Delete(ctx context.Context, req *v1.DeleteAuthorInfoReq) (res *v1.DeleteAuthorInfoRes, err error) {
	res, err = service.AuthorManage().Delete(ctx, req)
	return
}

func (f *authorManageController) Get(ctx context.Context, req *v1.GetAuthorInfoReq) (res *v1.GetAuthorInfoRes, err error) {
	res, err = service.AuthorManage().Get(ctx, req)
	return
}

func (f *authorManageController) Trend(ctx context.Context, req *v1.GetAuthorTrendReq) (res *v1.GetAuthorTrendRes, err error) {
	res, err = service.AuthorManage().Trend(ctx, req)
	return
}
