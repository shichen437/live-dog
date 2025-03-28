package media_parse

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
	service.RegisterMediaParse(New())
}

func New() *sMediaParse {
	return &sMediaParse{}
}

type sMediaParse struct {
}

func (s *sMediaParse) Parse(ctx context.Context, req *v1.PostMediaParseReq) (res *v1.PostMediaParseRes, err error) {
	res = &v1.PostMediaParseRes{}
	if req.Url == "" {
		return nil, gerror.New("url不能为空")
	}
	parser, err := media_parser.NewParser(req.Url)
	if err != nil {
		return
	}
	info, err := parser.ParseURL(ctx)
	if err != nil {
		return
	}
	mediaInfo := do.MediaParse{
		Platform:   info.Platform,
		Author:     info.Author,
		AuthorUid:  info.AuthorUid,
		MediaId:    info.VideoID,
		Desc:       info.Desc,
		Type:       info.Type,
		CreateTime: gtime.Now(),
	}
	if info.Type == "video" {
		mediaInfo.VideoUrl = info.VideoUrl
		mediaInfo.VideoCoverUrl = info.VideoCoverUrl
		mediaInfo.VideoData = info.VideoData
		dao.MediaParse.Ctx(ctx).Insert(mediaInfo)
	}
	if info.Type == "note" {
		mediaInfo.ImagesUrl = info.ImagesUrl
		mediaInfo.ImagesCoverUrl = info.ImagesCoverUrl
		dao.MediaParse.Ctx(ctx).Insert(mediaInfo)
	}
	return nil, nil
}

func (s *sMediaParse) List(ctx context.Context, req *v1.GetMediaParseListReq) (res *v1.GetMediaParseListRes, err error) {
	res = &v1.GetMediaParseListRes{}
	var list []*entity.MediaParse
	m := dao.MediaParse.Ctx(ctx)
	if req.Author != "" {
		m = m.WhereLike(dao.MediaParse.Columns().Author, "%"+req.Author+"%")
	}
	m = m.OrderDesc(dao.MediaParse.Columns().Id)
	res.Total, err = m.Count()
	utils.WriteErrLogT(ctx, err, consts.ListF)
	if res.Total > 0 {
		err = m.Page(req.PageNum, req.PageSize).Scan(&list)
		utils.WriteErrLogT(ctx, err, consts.ListF)
		res.Rows = list
	}
	return
}

func (s *sMediaParse) Get(ctx context.Context, req *v1.GetMediaParseReq) (res *v1.GetMediaParseRes, err error) {
	res = &v1.GetMediaParseRes{}
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.MediaParse.Ctx(ctx).Where(dao.MediaParse.Columns().Id, req.Id).Scan(&res)
		utils.WriteErrLogT(ctx, err, consts.GetF)
	})
	return
}

func (s *sMediaParse) Delete(ctx context.Context, req *v1.DeleteMediaParseReq) (res *v1.DeleteMediaParseRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		ids := utils.ParamStrToSlice(req.Id, ",")
		_, e := dao.MediaParse.Ctx(ctx).WhereIn(dao.MediaParse.Columns().Id, ids).Delete()
		utils.WriteErrLogT(ctx, e, consts.DeleteF)
	})
	return
}
