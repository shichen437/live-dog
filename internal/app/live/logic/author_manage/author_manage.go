package author_manage

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "github.com/shichen437/live-dog/api/v1/live"
	"github.com/shichen437/live-dog/internal/app/common/consts"
	"github.com/shichen437/live-dog/internal/app/live/dao"
	"github.com/shichen437/live-dog/internal/app/live/model/do"
	"github.com/shichen437/live-dog/internal/app/live/service"
	"github.com/shichen437/live-dog/internal/pkg/media_parser"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

func init() {
	service.RegisterAuthorManage(New())
}

func New() *sAuthorManage {
	return &sAuthorManage{}
}

type sAuthorManage struct {
}

func (s *sAuthorManage) List(ctx context.Context, req *v1.GetAuthorInfoListReq) (res *v1.GetAuthorInfoListRes, err error) {
	res = &v1.GetAuthorInfoListRes{}
	m := dao.AuthorInfo.Ctx(ctx)
	m = m.OrderDesc(dao.AuthorInfo.Columns().Id)
	res.Total, err = m.Count()
	err = m.Page(req.PageNum, req.PageSize).Scan(&res.Rows)
	utils.WriteErrLogT(ctx, err, consts.ListF)
	return res, err
}

func (s *sAuthorManage) Get(ctx context.Context, req *v1.GetAuthorInfoReq) (*v1.GetAuthorInfoRes, error) {
	return nil, nil
}

func (s *sAuthorManage) New(ctx context.Context, req *v1.PostAuthorInfoReq) (res *v1.PostAuthorInfoRes, err error) {
	if req.Url == "" {
		return nil, gerror.New("url is empty")
	}
	parser, err := media_parser.NewParser(req.Url)
	if err != nil {
		return nil, err
	}
	info, err := parser.ParseUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	// 检查是否已经存在
	authorInfo, err := dao.AuthorInfo.Ctx(ctx).
		Where(dao.AuthorInfo.Columns().UniqueId, info.UniqueId).
		Where(dao.AuthorInfo.Columns().Platform, info.Platform).
		One()
	if err != nil || authorInfo != nil {
		return nil, gerror.New("已添加过该平台博主")
	}
	result, err := dao.AuthorInfo.Ctx(ctx).Insert(do.AuthorInfo{
		Nickname:       info.Nickname,
		UniqueId:       info.UniqueId,
		Signature:      info.Signature,
		AvatarUrl:      info.Avatar,
		Platform:       info.Platform,
		FollowerCount:  info.FollowerCount,
		FollowingCount: info.FollowingCount,
		Ip:             info.IpLocation,
		Refer:          info.Refer,
		CreateTime:     gtime.Now(),
	})
	utils.WriteErrLogT(ctx, err, consts.AddF)
	infoId, err := result.LastInsertId()
	dao.AuthorInfoHistory.Ctx(ctx).Insert(do.AuthorInfoHistory{
		AuthorId:           infoId,
		LastFollowerCount:  info.FollowerCount,
		LastFollowingCount: info.FollowingCount,
		Num:                0,
		Day:                info.CurrentDay,
		CreateTime:         gtime.Now(),
	})
	return
}

func (s *sAuthorManage) Delete(ctx context.Context, req *v1.DeleteAuthorInfoReq) (*v1.DeleteAuthorInfoRes, error) {
	return nil, nil
}

func (s *sAuthorManage) Trend(ctx context.Context, req *v1.GetAuthorTrendReq) (*v1.GetAuthorTrendRes, error) {
	return nil, nil
}
