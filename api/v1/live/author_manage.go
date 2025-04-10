package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/live-dog/api/v1/common"
	"github.com/shichen437/live-dog/internal/app/live/model/entity"
)

type GetAuthorInfoListReq struct {
	g.Meta `path:"/author/manage/list" method:"get" tags:"博主信息" summary:"博主列表"`
	common.PageReq
	Nickname string `p:"nickname"`
}

type GetAuthorInfoListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.AuthorInfo `json:"rows"`
	Total  int                  `json:"total"`
}

type GetAuthorInfoReq struct {
	g.Meta `path:"/author/manage/{id}" method:"get" tags:"博主信息" summary:"博主详情"`
	Id     int `p:"id"`
}

type GetAuthorInfoRes struct {
	g.Meta `mime:"application/json"`
	*entity.AuthorInfo
}

type PostAuthorInfoReq struct {
	g.Meta `path:"/author/manage/" method:"post" tags:"博主信息" summary:"添加博主"`
	Url    string `p:"url"`
}

type PostAuthorInfoRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteAuthorInfoReq struct {
	g.Meta `path:"/author/manage/{id}" method:"delete" tags:"博主信息" summary:"删除博主"`
	Id     string `p:"id"`
}
type DeleteAuthorInfoRes struct {
	g.Meta `mime:"application/json"`
}

type GetAuthorTrendReq struct {
	g.Meta `path:"/author/manage/trend" method:"get" tags:"博主信息" summary:"博主粉丝数据"`
	Id     int  `p:"id"`
	Range  *int `p:"range"`
}

type GetAuthorTrendRes struct {
	g.Meta `mime:"application/json"`
	Days   []string `json:"days"`
	Counts []int64  `json:"counts"`
	Nums   []int    `json:"nums"`
}
