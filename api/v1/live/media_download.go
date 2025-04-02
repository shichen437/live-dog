package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/live-dog/api/v1/common"
	"github.com/shichen437/live-dog/internal/app/live/model/entity"
)

type GetDownloadRecordReq struct {
	g.Meta `path:"/media/download/list" method:"get" tags:"下载中心" summary:"下载记录列表"`
	common.PageReq
	Status string `p:"status"`
}

type GetDownloadRecordRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.DownloadRecord `json:"rows"`
	Total  int                      `json:"total"`
}

type GetDownloadRecordCacheReq struct {
	g.Meta `path:"/media/download/listCache" method:"get" tags:"下载中心" summary:"下载记录列表(缓存)"`
}

type GetDownloadRecordCacheRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.DownloadRecord `json:"rows"`
}

type DeleteDownloadRecordReq struct {
	g.Meta `path:"/media/download/{id}" method:"delete" tags:"下载中心" summary:"删除下载记录"`
	Id     string `p:"id"  v:"required"`
}

type DeleteDownloadRecordRes struct {
	g.Meta `mime:"application/json"`
}
