package media_download

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/shichen437/live-dog/api/v1/live"
	"github.com/shichen437/live-dog/internal/app/common/consts"
	"github.com/shichen437/live-dog/internal/app/live/dao"
	"github.com/shichen437/live-dog/internal/app/live/model/entity"
	"github.com/shichen437/live-dog/internal/app/live/service"
	"github.com/shichen437/live-dog/internal/pkg/download"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

func init() {
	service.RegisterMediaDownload(New())
}

func New() *sMediaDownload {
	return &sMediaDownload{}
}

type sMediaDownload struct {
}

func (s *sMediaDownload) List(ctx context.Context, req *v1.GetDownloadRecordReq) (res *v1.GetDownloadRecordRes, err error) {
	res = &v1.GetDownloadRecordRes{}
	var list []*entity.DownloadRecord
	m := dao.DownloadRecord.Ctx(ctx)
	if req.Status != "" {
		m = m.Where(dao.DownloadRecord.Columns().Status, req.Status)
	}
	m = m.OrderDesc(dao.DownloadRecord.Columns().Id)
	res.Total, err = m.Count()
	utils.WriteErrLogT(ctx, err, consts.ListF)
	if res.Total > 0 {
		err = m.Page(req.PageNum, req.PageSize).Scan(&list)
		utils.WriteErrLogT(ctx, err, consts.ListF)
		res.Rows = list
	}
	return
}

func (s *sMediaDownload) ListFromCache(ctx context.Context, req *v1.GetDownloadRecordCacheReq) (res *v1.GetDownloadRecordCacheRes, err error) {
	res = &v1.GetDownloadRecordCacheRes{}
	var list []*entity.DownloadRecord
	pm := download.GetProgressManager()
	taskList := pm.ListAllTasks(10)
	for _, task := range taskList {
		list = append(list, &entity.DownloadRecord{
			TaskId:   task.TaskID,
			Title:    task.Title,
			Status:   string(task.Status),
			ErrorMsg: task.ErrorMsg,
		})
	}
	res.Rows = list
	return
}

func (s *sMediaDownload) Delete(ctx context.Context, req *v1.DeleteDownloadRecordReq) (res *v1.DeleteDownloadRecordRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		ids := utils.ParamStrToSlice(req.Id, ",")
		_, e := dao.DownloadRecord.Ctx(ctx).WhereIn(dao.DownloadRecord.Columns().Id, ids).Delete()
		utils.WriteErrLogT(ctx, e, consts.DeleteF)
	})
	return
}
