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
	IFileManage interface {
		List(ctx context.Context, req *v1.GetFileInfoListReq) (res *v1.GetFileInfoListRes, err error)
		Delete(ctx context.Context, req *v1.DeleteFileInfoReq) (res *v1.DeleteFileInfoRes, err error)
		Play(ctx context.Context, req *v1.GetFilePlayReq) (res *v1.GetFilePlayRes, err error)
	}
)

var (
	localFileManage IFileManage
)

func FileManage() IFileManage {
	if localFileManage == nil {
		panic("implement not found for interface IFileManage, forgot register?")
	}
	return localFileManage
}

func RegisterFileManage(i IFileManage) {
	localFileManage = i
}
