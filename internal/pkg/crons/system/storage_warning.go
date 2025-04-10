package system

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shichen437/live-dog/internal/app/monitor/model/entity"
	"github.com/shichen437/live-dog/internal/app/monitor/service"
	mp "github.com/shichen437/live-dog/internal/pkg/message_push"
	"github.com/shirou/gopsutil/v3/disk"
)

var (
	invokeTarget = "storageWarning"
	title        = "空间预警"
)

type StorageWarningParams struct {
	Threshold float64 `json:"threshold"`
}

func StorageWarning(jobId int64, params, jobName string) {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		service.SysJob().AddLog(gctx.New(), &entity.SysJobLog{
			JobId:         jobId,
			JobName:       jobName,
			InvokeTarget:  invokeTarget,
			JobMessage:    title,
			Status:        "1",
			ExceptionInfo: err.Error(),
		})
		return
	}
	threshold := 90.00
	if gjson.Valid(params) {
		param := &StorageWarningParams{}
		err := gconv.Struct(params, param)
		if err == nil && param.Threshold > 0 {
			threshold = param.Threshold
		}
	}
	if diskInfo.UsedPercent > threshold {
		mp.CustomPush(gctx.New(), &mp.MessageModel{
			Title:   title,
			Content: "存储空间已达到" + gconv.String(threshold) + "%",
		})
	}
	// 任务日志
	service.SysJob().AddLog(gctx.New(), &entity.SysJobLog{
		JobId:         jobId,
		JobName:       jobName,
		InvokeTarget:  invokeTarget,
		JobMessage:    "执行成功",
		Status:        "0",
		ExceptionInfo: "",
	})
}
