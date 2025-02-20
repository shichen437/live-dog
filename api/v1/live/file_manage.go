package v1

import (
	"github.com/shichen437/live-dog/internal/app/live/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetFileInfoListReq struct {
	g.Meta   `path:"/file/manage/list" method:"get" tags:"文件管理" summary:"文件列表"`
	Filename string `p:"filename"`
	Path     string `p:"path"`
}
type GetFileInfoListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*model.FileInfo `json:"rows"`
}

type DeleteFileInfoReq struct {
	g.Meta    `path:"/file/manage" method:"delete" tags:"文件管理" summary:"删除文件"`
	Filenames []string `p:"filenames"`
	Path      string   `p:"path"`
}
type DeleteFileInfoRes struct {
	g.Meta `mime:"application/json"`
}

type GetFilePlayReq struct {
	g.Meta `path:"/file/manage/play" method:"get" tags:"文件管理" summary:"视频播放"`
	Path   string `p:"path" v:"required#请输入文件路径"`
}

type GetFilePlayRes struct {
	g.Meta `mime:"video/*"`
}
