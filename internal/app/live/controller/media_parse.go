package live

import (
	"context"

	v1 "github.com/shichen437/live-dog/api/v1/live"
	"github.com/shichen437/live-dog/internal/app/live/service"
)

type mediaParseController struct {
}

var MediaParse = mediaParseController{}

func (c *mediaParseController) Parse(ctx context.Context, req *v1.PostMediaParseReq) (res *v1.PostMediaParseRes, err error) {
	res, err = service.MediaParse().Parse(ctx, req)
	return
}

func (c *mediaParseController) List(ctx context.Context, req *v1.GetMediaParseListReq) (res *v1.GetMediaParseListRes, err error) {
	res, err = service.MediaParse().List(ctx, req)
	return
}

func (c *mediaParseController) Get(ctx context.Context, req *v1.GetMediaParseReq) (res *v1.GetMediaParseRes, err error) {
	res, err = service.MediaParse().Get(ctx, req)
	return
}

func (c *mediaParseController) Delete(ctx context.Context, req *v1.DeleteMediaParseReq) (res *v1.DeleteMediaParseRes, err error) {
	res, err = service.MediaParse().Delete(ctx, req)
	return
}

func (c *mediaParseController) Download(ctx context.Context, req *v1.GetDownloadMediaReq) (res *v1.GetDownloadMediaRes, err error) {
	res, err = service.MediaParse().Download(ctx, req)
	return
}
