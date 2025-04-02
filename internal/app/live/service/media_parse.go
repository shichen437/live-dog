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
	IMediaParse interface {
		Parse(ctx context.Context, req *v1.PostMediaParseReq) (res *v1.PostMediaParseRes, err error)
		List(ctx context.Context, req *v1.GetMediaParseListReq) (res *v1.GetMediaParseListRes, err error)
		Get(ctx context.Context, req *v1.GetMediaParseReq) (res *v1.GetMediaParseRes, err error)
		Delete(ctx context.Context, req *v1.DeleteMediaParseReq) (res *v1.DeleteMediaParseRes, err error)
		Download(ctx context.Context, req *v1.GetDownloadMediaReq) (res *v1.GetDownloadMediaRes, err error)
	}
)

var (
	localMediaParse IMediaParse
)

func MediaParse() IMediaParse {
	if localMediaParse == nil {
		panic("implement not found for interface IMediaParse, forgot register?")
	}
	return localMediaParse
}

func RegisterMediaParse(i IMediaParse) {
	localMediaParse = i
}
