package author_manage

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "github.com/shichen437/live-dog/api/v1/live"
	"github.com/shichen437/live-dog/internal/app/common/consts"
	"github.com/shichen437/live-dog/internal/app/live/dao"
	"github.com/shichen437/live-dog/internal/app/live/model/do"
	"github.com/shichen437/live-dog/internal/app/live/model/entity"
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
	if req.Nickname != "" {
		m = m.WhereLike(dao.AuthorInfo.Columns().Nickname, "%"+req.Nickname+"%")
	}
	m = m.OrderDesc(dao.AuthorInfo.Columns().Id)
	res.Total, err = m.Count()
	err = m.Page(req.PageNum, req.PageSize).Scan(&res.Rows)
	utils.WriteErrLogT(ctx, err, consts.ListF)
	return res, err
}

func (s *sAuthorManage) Get(ctx context.Context, req *v1.GetAuthorInfoReq) (*v1.GetAuthorInfoRes, error) {
	if req.Id == 0 {
		return nil, gerror.New("id is empty")
	}
	authorInfo, err := dao.AuthorInfo.Ctx(ctx).WherePri(req.Id).One()
	if err != nil || authorInfo == nil {
		return nil, gerror.New("未获取到有效信息")
	}
	res := &v1.GetAuthorInfoRes{}
	authorInfo.Struct(res)
	return res, nil
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
	if err != nil || info == nil || info.UniqueId == "" {
		g.Log().Error(ctx, err)
		return nil, gerror.New("解析博主信息失败")
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
	ids := utils.ParamStrToSlice(req.Id, ",")
	_, e := dao.AuthorInfo.Ctx(ctx).WhereIn(dao.AuthorInfo.Columns().Id, ids).Delete()
	utils.WriteErrLogT(ctx, e, consts.DeleteF)
	_, e = dao.AuthorInfoHistory.Ctx(ctx).WhereIn(dao.AuthorInfoHistory.Columns().AuthorId, ids).Delete()
	utils.WriteErrLogT(ctx, e, consts.DeleteF)
	return &v1.DeleteAuthorInfoRes{}, nil
}

func (s *sAuthorManage) Trend(ctx context.Context, req *v1.GetAuthorTrendReq) (*v1.GetAuthorTrendRes, error) {
	if req.Id == 0 {
		return nil, gerror.New("id is empty")
	}
	var list []*entity.AuthorInfoHistory
	m := dao.AuthorInfoHistory.Ctx(ctx)
	m = m.Where(dao.AuthorInfoHistory.Columns().AuthorId, req.Id)
	m = m.OrderDesc(dao.AuthorInfoHistory.Columns().Id)
	if req.Range != nil {
		m = m.Limit(0, *req.Range)
	}
	err := m.Scan(&list)
	if err != nil {
		return nil, gerror.New("获取数据失败")
	}
	res := &v1.GetAuthorTrendRes{
		Days:   make([]string, 0, len(list)),
		Counts: make([]int64, 0, len(list)),
		Nums:   make([]int, 0, len(list)),
	}
	for i := len(list) - 1; i >= 0; i-- {
		item := list[i]
		res.Days = append(res.Days, item.Day)
		res.Counts = append(res.Counts, item.LastFollowerCount)
		res.Nums = append(res.Nums, item.Num)
	}
	return res, nil
}
