// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/shichen437/live-dog/api/v1/live"
)

type (
	IMediaDownload interface {
		List(ctx context.Context, req *v1.GetDownloadRecordReq) (res *v1.GetDownloadRecordRes, err error)
		ListFromCache(ctx context.Context, req *v1.GetDownloadRecordCacheReq) (res *v1.GetDownloadRecordCacheRes, err error)
		Delete(ctx context.Context, req *v1.DeleteDownloadRecordReq) (res *v1.DeleteDownloadRecordRes, err error)
	}
)

var (
	localMediaDownload IMediaDownload
)

func MediaDownload() IMediaDownload {
	if localMediaDownload == nil {
		panic("implement not found for interface IMediaDownload, forgot register?")
	}
	return localMediaDownload
}

func RegisterMediaDownload(i IMediaDownload) {
	localMediaDownload = i
}
