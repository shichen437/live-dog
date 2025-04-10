package system

import (
	"context"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/live-dog/internal/app/live/dao"
	"github.com/shichen437/live-dog/internal/app/live/model/entity"
	mEntity "github.com/shichen437/live-dog/internal/app/monitor/model/entity"
	"github.com/shichen437/live-dog/internal/app/monitor/service"
	"github.com/shichen437/live-dog/internal/pkg/media_parser"
)

type updateHistoryModel struct {
	AuthorInfo  *entity.AuthorInfo
	CurrentDay  string
	PreviousDay string
	Info        *media_parser.UserInfo
}

func FollowerTrend(jobId int64, jobName string) {
	ctx := gctx.GetInitCtx()
	authorList := getAuthorList(ctx)
	if authorList == nil {
		g.Log().Error(ctx, "[定时任务]获取博主列表失败")
		service.SysJob().AddLog(gctx.New(), &mEntity.SysJobLog{
			JobId:         jobId,
			JobName:       jobName,
			InvokeTarget:  "followerTrend",
			JobMessage:    "获取博主列表为空",
			Status:        "0",
			ExceptionInfo: "",
		})
		return
	}
	currentDay := gtime.Now().Format("Y-m-d")
	previousDay := gtime.Now().AddDate(0, 0, -1).Format("Y-m-d")
	for _, author := range authorList {
		if author.Refer == "" {
			continue
		}
		parser, err := media_parser.NewParser(author.Refer)
		if err != nil {
			continue
		}
		time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)
		info, err := parser.ParseUserInfo(ctx)
		if err != nil || info == nil || info.UniqueId != author.UniqueId {
			g.Log().Error(ctx, "[定时任务]获取博主信息失败", err)
			continue
		}
		addOrUpdateAuthorInfoHistory(ctx, &updateHistoryModel{
			AuthorInfo:  author,
			CurrentDay:  currentDay,
			PreviousDay: previousDay,
			Info:        info,
		})
		updateAuthorInfo(ctx, info, author.Id)
	}
	service.SysJob().AddLog(gctx.New(), &mEntity.SysJobLog{
		JobId:         jobId,
		JobName:       jobName,
		InvokeTarget:  "followerTrend",
		JobMessage:    "执行成功",
		Status:        "0",
		ExceptionInfo: "",
	})
}

func getAuthorList(ctx context.Context) []*entity.AuthorInfo {
	dao.AuthorInfo.Ctx(ctx)
	result, err := dao.AuthorInfo.Ctx(ctx).All()
	if err != nil {
		return nil
	}
	var authorList []*entity.AuthorInfo
	result.Structs(&authorList)
	return authorList
}

func getAuthorInfoHistory(ctx context.Context, authorId int64, day string) (history *entity.AuthorInfoHistory, err error) {
	result, err := dao.AuthorInfoHistory.Ctx(ctx).
		Where(dao.AuthorInfoHistory.Columns().AuthorId, authorId).
		Where(dao.AuthorInfoHistory.Columns().Day, day).
		One()
	if err != nil {
		return nil, err
	}
	if result != nil {
		result.Struct(history)
		return history, nil
	}
	return nil, nil
}

func addOrUpdateAuthorInfoHistory(ctx context.Context, model *updateHistoryModel) {
	currentHistory, err := getAuthorInfoHistory(ctx, int64(model.AuthorInfo.Id), model.CurrentDay)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "[定时任务]获取今日历史数据失败", err)
		return
	}
	previousHistory, err := getAuthorInfoHistory(ctx, int64(model.AuthorInfo.Id), model.PreviousDay)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "[定时任务]获取历史数据失败", err)
		return
	}
	num := 0
	if previousHistory != nil {
		num = model.Info.FollowerCount - int(previousHistory.LastFollowerCount)
	}
	if currentHistory == nil {
		// 新增
		dao.AuthorInfoHistory.Ctx(ctx).Insert(entity.AuthorInfoHistory{
			AuthorId:           model.AuthorInfo.Id,
			Day:                model.CurrentDay,
			Num:                num,
			LastFollowerCount:  model.AuthorInfo.FollowerCount,
			LastFollowingCount: model.AuthorInfo.FollowingCount,
			CreateTime:         gtime.Now(),
		})
	} else {
		// 更新
		dao.AuthorInfoHistory.Ctx(ctx).WherePri(currentHistory.Id).Update(entity.AuthorInfoHistory{
			Num:                num,
			LastFollowerCount:  model.AuthorInfo.FollowerCount,
			LastFollowingCount: model.AuthorInfo.FollowingCount,
			UpdateTime:         gtime.Now(),
		})
	}
}

func updateAuthorInfo(ctx context.Context, info *media_parser.UserInfo, authorId int) {
	dao.AuthorInfo.Ctx(ctx).WherePri(authorId).Update(entity.AuthorInfo{
		Nickname:       info.Nickname,
		AvatarUrl:      info.Avatar,
		Signature:      info.Signature,
		FollowerCount:  int64(info.FollowerCount),
		FollowingCount: info.FollowingCount,
		UpdateTime:     gtime.Now(),
	})
	return
}
