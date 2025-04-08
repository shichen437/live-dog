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
	IAuthorManage interface {
		List(ctx context.Context, req *v1.GetAuthorInfoListReq) (*v1.GetAuthorInfoListRes, error)
		Get(ctx context.Context, req *v1.GetAuthorInfoReq) (*v1.GetAuthorInfoRes, error)
		New(ctx context.Context, req *v1.PostAuthorInfoReq) (*v1.PostAuthorInfoRes, error)
		Delete(ctx context.Context, req *v1.DeleteAuthorInfoReq) (*v1.DeleteAuthorInfoRes, error)
		Trend(ctx context.Context, req *v1.GetAuthorTrendReq) (*v1.GetAuthorTrendRes, error)
	}
)

var (
	localAuthorManage IAuthorManage
)

func AuthorManage() IAuthorManage {
	if localAuthorManage == nil {
		panic("implement not found for interface IAuthorManage, forgot register?")
	}
	return localAuthorManage
}

func RegisterAuthorManage(i IAuthorManage) {
	localAuthorManage = i
}
