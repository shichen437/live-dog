package system

import (
	"context"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/live-dog/internal/app/live/dao"
	"github.com/shichen437/live-dog/internal/app/live/model/do"
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
		var info *media_parser.UserInfo
		var parseErr error
		maxRetries := 3

		for i := 0; i < maxRetries; i++ {
			info, parseErr = parser.ParseUserInfo(ctx)
			if parseErr == nil && info != nil && info.UniqueId == author.UniqueId {
				break
			}
			if i < maxRetries-1 {
				time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)
			}
		}

		if parseErr != nil || info == nil || info.UniqueId != author.UniqueId {
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
		time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)
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
	history = &entity.AuthorInfoHistory{}
	if err != nil {
		return nil, err
	}
	if result != nil {
		if err = result.Struct(&history); err != nil {
			return nil, err
		}
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
			LastFollowerCount:  int64(model.Info.FollowerCount),
			LastFollowingCount: model.Info.FollowingCount,
			CreateTime:         gtime.Now(),
		})
	} else {
		// 更新
		dao.AuthorInfoHistory.Ctx(ctx).
			Data(do.AuthorInfoHistory{
				Num:                num,
				LastFollowerCount:  int64(model.Info.FollowerCount),
				LastFollowingCount: model.Info.FollowingCount,
				UpdateTime:         gtime.Now(),
			}).
			Where(dao.AuthorInfoHistory.Columns().Id, currentHistory.Id).
			Update()
	}
}

func updateAuthorInfo(ctx context.Context, info *media_parser.UserInfo, authorId int) {
	dao.AuthorInfo.Ctx(ctx).Where(dao.AuthorInfo.Columns().Id, authorId).Update(do.AuthorInfo{
		Nickname:       info.Nickname,
		AvatarUrl:      info.Avatar,
		Signature:      info.Signature,
		FollowerCount:  int64(info.FollowerCount),
		FollowingCount: info.FollowingCount,
		UpdateTime:     gtime.Now(),
	})
	return
}
