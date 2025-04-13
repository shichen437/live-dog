package sys_job

import (
	"context"

	v1 "github.com/shichen437/live-dog/api/v1/monitor"
	"github.com/shichen437/live-dog/internal/app/common/consts"
	"github.com/shichen437/live-dog/internal/app/monitor/dao"
	"github.com/shichen437/live-dog/internal/app/monitor/model/do"
	"github.com/shichen437/live-dog/internal/app/monitor/model/entity"
	"github.com/shichen437/live-dog/internal/app/monitor/service"
	"github.com/shichen437/live-dog/internal/pkg/crons"
	"github.com/shichen437/live-dog/internal/pkg/crons/system"
	"github.com/shichen437/live-dog/internal/pkg/utils"

	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	service.RegisterSysJob(New())
}

func New() *sSysJob {
	return &sSysJob{}
}

type sSysJob struct{}

func (s *sSysJob) List(ctx context.Context, req *v1.GetJobListReq) (result *v1.GetJobListRes, err error) {
	result = &v1.GetJobListRes{}

	if req.PageNum == 0 {
		req.PageNum = 1
	}
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}

	err = g.Try(ctx, func(ctx context.Context) {
		var list []*entity.SysJob
		m := dao.SysJob.Ctx(ctx)
		if req.JobName != "" {
			m = m.WhereLike(dao.SysJob.Columns().JobName, "%"+req.JobName+"%")
		}
		if req.Status != nil {
			m = m.Where(dao.SysJob.Columns().Status, *req.Status)
		}
		m = m.OrderDesc(dao.SysJob.Columns().CreateTime)
		result.Total, err = m.Count()
		err = m.Page(req.PageNum, req.PageSize).Scan(&list)
		result.Rows = list
	})
	return
}

func (s *sSysJob) LogList(ctx context.Context, req *v1.GetJobLogListReq) (result *v1.GetJobLogListRes, err error) {
	result = &v1.GetJobLogListRes{}

	if req.PageNum == 0 {
		req.PageNum = 1
	}
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}

	err = g.Try(ctx, func(ctx context.Context) {
		var list []*entity.SysJobLog
		m := dao.SysJobLog.Ctx(ctx)
		if req.JobId != 0 {
			m = m.Where(dao.SysJobLog.Columns().JobId, req.JobId)
		}
		if req.JobName != "" {
			m = m.WhereLike(dao.SysJobLog.Columns().JobName, "%"+req.JobName+"%")
		}
		if req.Params != nil {
			if req.Params["beginTime"] != "" && req.Params["endTime"] != "" {
				m = m.WhereBetween(dao.SysJobLog.Columns().CreateTime, req.Params["beginTime"], req.Params["endTime"])
			}
		}
		if req.Status != nil {
			m = m.Where(dao.SysJobLog.Columns().Status, req.Status)
		}
		m = m.OrderDesc(dao.SysJobLog.Columns().CreateTime)
		result.Total, err = m.Count()
		err = m.Page(req.PageNum, req.PageSize).Scan(&list)
		result.Rows = list
	})
	return
}

func (s *sSysJob) Get(ctx context.Context, req *v1.GetJobDetailReq) (result *v1.GetJobDetailRes, err error) {
	if req.JobId == 0 {
		err = utils.TError(ctx, consts.IDEmpty)
		return
	}
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.SysJob.Ctx(ctx).Where(dao.SysJob.Columns().JobId, req.JobId).Scan(&result)
	})
	return
}

func (s *sSysJob) Add(ctx context.Context, req *v1.PostJobReq) (result *v1.PostJobRes, err error) {
	adminName := gconv.String(ctx.Value(consts.CtxAdminName))
	err = g.Try(ctx, func(ctx context.Context) {
		//添加角色信息
		r, err := dao.SysJob.Ctx(ctx).Insert(do.SysJob{
			JobName:        req.JobName,
			InvokeTarget:   req.InvokeTarget,
			CronExpression: req.CronExpression,
			MisfirePolicy:  req.MisfirePolicy,
			Concurrent:     req.Concurrent,
			Status:         req.Status,
			Type:           0,
			JobParams:      req.JobParams,
			Remark:         req.Remark,
			CreateBy:       adminName,
			CreateTime:     gtime.Now(),
		})
		utils.WriteErrLogT(ctx, err, consts.AddF)
		id, _ := r.LastInsertId()
		if id > 0 && req.Status == "0" {
			crons.AddSystemJob(id)
		}
	})
	return
}

func (s *sSysJob) Update(ctx context.Context, req *v1.PutJobReq) (result *v1.PutJobRes, err error) {
	if req.JobId == 0 {
		err = utils.TError(ctx, consts.IDEmpty)
		return
	}
	adminName := gconv.String(ctx.Value(consts.CtxAdminName))
	err = g.Try(ctx, func(ctx context.Context) {
		//更新用户信息
		_, e := dao.SysJob.Ctx(ctx).WherePri(req.JobId).Update(do.SysJob{
			JobName:        req.JobName,
			CronExpression: req.CronExpression,
			MisfirePolicy:  req.MisfirePolicy,
			Concurrent:     req.Concurrent,
			Status:         req.Status,
			JobParams:      req.JobParams,
			Remark:         req.Remark,
			UpdateTime:     gtime.Now(),
			UpdateBy:       adminName,
		})
		utils.WriteErrLogT(ctx, e, consts.UpdateF)
		crons.RemoveSystemJob(req.JobId)
		if req.Status == "0" {
			crons.AddSystemJob(req.JobId)
		}
	})
	return
}

func (s *sSysJob) Delete(ctx context.Context, req *v1.DeleteJobReq) (result *v1.DeleteJobRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		ids := utils.ParamStrToSlice(req.JobId, ",")
		//删除用户信息
		_, e := dao.SysJob.Ctx(ctx).WhereIn(dao.SysJob.Columns().JobId, ids).Delete()
		utils.WriteErrLogT(ctx, e, consts.DeleteF)
		//删除定时任务
		for _, id := range ids {
			crons.RemoveSystemJob(id)
		}
	})
	return
}

func (s *sSysJob) DeleteLog(ctx context.Context, req *v1.DeleteJobLogReq) (result *v1.DeleteJobLogRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		ids := utils.ParamStrToSlice(req.JobLogId, ",")
		_, e := dao.SysJobLog.Ctx(ctx).WhereIn(dao.SysJobLog.Columns().JobLogId, ids).Delete()
		utils.WriteErrLogT(ctx, e, consts.DeleteF)
	})
	return
}

func (s *sSysJob) RunJob(ctx context.Context, req *v1.PutJobRunReq) (result *v1.PutJobRunRes, err error) {
	if req.JobId == 0 {
		err = utils.TError(ctx, consts.IDEmpty)
		return
	}
	r, err := s.GetJobDetail(ctx, req.JobId)
	if err != nil || r == nil {
		return
	}
	switch r.InvokeTarget {
	case "storageWarning":
		go system.StorageWarning(r.JobId, r.JobParams, r.JobName)
	case "followerTrend":
		go system.FollowerTrend(r.JobId, r.JobName)
	default:
		return
	}
	return
}

func (s *sSysJob) ChangeStatus(ctx context.Context, req *v1.PutJobStatusReq) (result *v1.PutJobStatusRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, e := dao.SysJob.Ctx(ctx).WherePri(req.JobId).Update(do.SysJob{
			Status: req.Status,
		})
		utils.WriteErrLogT(ctx, e, consts.UpdateF)
		crons.RemoveSystemJob(req.JobId)
		if req.Status == "0" {
			crons.AddSystemJob(req.JobId)
		}
	})
	return
}

func (s *sSysJob) ListAll4Init(ctx context.Context) (result []*entity.SysJob, err error) {
	m := dao.SysJob.Ctx(ctx)
	err = g.Try(ctx, func(ctx context.Context) {
		m = m.Where(dao.SysJob.Columns().Status, "0")
		err = m.Scan(&result)
	})
	if len(result) == 0 {
		return nil, nil
	}
	return
}

func (s *sSysJob) GetJobDetail(ctx context.Context, jobId int64) (result *entity.SysJob, err error) {
	m := dao.SysJob.Ctx(ctx)
	err = g.Try(ctx, func(ctx context.Context) {
		m = m.Where(dao.SysJob.Columns().JobId, jobId)
		err = m.Scan(&result)
	})
	return
}

func (s *sSysJob) AddLog(ctx context.Context, jobLog *entity.SysJobLog) (err error) {
	jobLog.CreateTime = gtime.Now()
	_, err = dao.SysJobLog.Ctx(ctx).Insert(jobLog)
	return
}
