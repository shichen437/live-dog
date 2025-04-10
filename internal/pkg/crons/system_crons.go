package crons

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/live-dog/internal/app/monitor/service"
	"github.com/shichen437/live-dog/internal/pkg/crons/system"
)

func AddSystemJob(jobId int64) {
	result, err := service.SysJob().GetJobDetail(gctx.New(), jobId)
	if err != nil || result == nil || result.Status != "0" || result.Type != 0 {
		return
	}
	switch result.InvokeTarget {
	case "storageWarning":
		gcron.Add(strings.Trim(result.CronExpression, " "), func() {
			g.Log().Info(gctx.New(), "添加系统定时任务-", result.JobName)
			system.StorageWarning(result.JobId, result.JobParams, result.JobName)
		}, "Cron-System-Job-"+strconv.Itoa(int(jobId)))
	case "followerTrend":
		gcron.Add(strings.Trim(result.CronExpression, " "), func() {
			g.Log().Info(gctx.New(), "添加系统定时任务-", result.JobName)
			system.FollowerTrend(result.JobId, result.JobName)
		}, "Cron-System-Job-"+strconv.Itoa(int(jobId)))
	default:
	}

}

func RemoveSystemJob(jobId int64) {
	g.Log().Info(gctx.New(), "移除系统定时任务-", jobId)
	job := gcron.Search("Cron-System-Job-" + strconv.Itoa(int(jobId)))
	if job != nil {
		gcron.Remove(job.Name)
	}
}
