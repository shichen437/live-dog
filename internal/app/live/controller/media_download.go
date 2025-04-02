package live

import (
	"context"

	v1 "github.com/shichen437/live-dog/api/v1/live"
	"github.com/shichen437/live-dog/internal/app/live/service"
)

type mediaDownloadController struct {
}

var MediaDownload = mediaDownloadController{}

func (c *mediaDownloadController) List(ctx context.Context, req *v1.GetDownloadRecordReq) (res *v1.GetDownloadRecordRes, err error) {
	res, err = service.MediaDownload().List(ctx, req)
	return
}

func (c *mediaDownloadController) ListFromCache(ctx context.Context, req *v1.GetDownloadRecordCacheReq) (res *v1.GetDownloadRecordCacheRes, err error) {
	res, err = service.MediaDownload().ListFromCache(ctx, req)
	return
}

func (c *mediaDownloadController) Delete(ctx context.Context, req *v1.DeleteDownloadRecordReq) (res *v1.DeleteDownloadRecordRes, err error) {
	res, err = service.MediaDownload().Delete(ctx, req)
	return
}
