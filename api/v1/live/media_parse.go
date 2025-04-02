package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/live-dog/api/v1/common"
	"github.com/shichen437/live-dog/internal/app/live/model/entity"
)

type GetMediaParseListReq struct {
	g.Meta `path:"/media/parse/list" method:"get" tags:"媒体解析" summary:"解析列表"`
	common.PageReq
	Author string `p:"author"`
}
type GetMediaParseListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.MediaParse `json:"rows"`
	Total  int                  `json:"total"`
}

type GetMediaParseReq struct {
	g.Meta `path:"/media/parse/{id}" method:"get" tags:"媒体解析" summary:"解析详情"`
	Id     int `p:"id"  v:"required"`
}

type GetMediaParseRes struct {
	g.Meta `mime:"application/json"`
	*entity.MediaParse
}

type PostMediaParseReq struct {
	g.Meta `path:"/media/parse" method:"post" tags:"媒体解析" summary:"解析媒体"`
	Url    string `p:"url" v:"required"`
}

type PostMediaParseRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteMediaParseReq struct {
	g.Meta `path:"/media/parse/{id}" method:"delete" tags:"媒体解析" summary:"删除解析"`
	Id     string `p:"id"  v:"required"`
}

type DeleteMediaParseRes struct {
	g.Meta `mime:"application/json"`
}

type GetDownloadMediaReq struct {
	g.Meta      `path:"/media/parse/download" method:"get" tags:"媒体解析" summary:"下载解析"`
	Id          int    `p:"id"  v:"required"`
	QualityDesc string `p:"quality_desc"`
}

type GetDownloadMediaRes struct {
	g.Meta `mime:"application/json"`
}
