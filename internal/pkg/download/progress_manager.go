package download

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/live-dog/internal/app/live/dao"
	"github.com/shichen437/live-dog/internal/app/live/model/do"
	"github.com/shichen437/live-dog/internal/pkg/sse"
)

var (
	progressManager     *ProgressManager
	progressManagerOnce sync.Once
)

type ProgressManager struct {
	progressMap sync.Map
}

func GetProgressManager() *ProgressManager {
	progressManagerOnce.Do(func() {
		progressManager = &ProgressManager{}
	})
	return progressManager
}

func (pm *ProgressManager) CreateTask(taskID, title, outputPath string) *DownloadProgress {
	progress := &DownloadProgress{
		TaskID:     taskID,
		Title:      title,
		Status:     DownloadStatusPending,
		OutputPath: outputPath,
		StartTime:  gtime.Now(),
		UpdateTime: gtime.Now(),
	}
	pm.progressMap.Store(taskID, progress)
	newRecord(progress)
	return progress
}

func (pm *ProgressManager) GetProgress(taskID string) (*DownloadProgress, error) {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return nil, gerror.Newf("任务 %s 不存在", taskID)
	}
	return value.(*DownloadProgress), nil
}

func (pm *ProgressManager) UpdateProgress(taskID string, status DownloadStatus) error {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return gerror.Newf("任务 %s 不存在", taskID)
	}

	dp := value.(*DownloadProgress)
	dp.Status = status
	dp.UpdateTime = gtime.Now()

	pm.progressMap.Store(taskID, dp)

	updateRecord(dp)
	return nil
}

func (pm *ProgressManager) SetError(taskID, errMsg string) error {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return gerror.Newf("任务 %s 不存在", taskID)
	}

	dp := value.(*DownloadProgress)
	dp.Status = DownloadStatusError
	dp.ErrorMsg = errMsg
	dp.UpdateTime = gtime.Now()

	pm.progressMap.Store(taskID, dp)
	updateRecord(dp)
	return nil
}

func (pm *ProgressManager) SetCompleted(taskID string) error {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return gerror.Newf("任务 %s 不存在", taskID)
	}

	dp := value.(*DownloadProgress)
	dp.Status = DownloadStatusCompleted
	dp.UpdateTime = gtime.Now()

	pm.progressMap.Store(taskID, dp)

	updateRecord(dp)
	return nil
}

func (pm *ProgressManager) SetPartCompleted(taskID, errMsg string) error {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return gerror.Newf("任务 %s 不存在", taskID)
	}

	dp := value.(*DownloadProgress)
	dp.Status = DownloadStatusPartSucceed
	dp.ErrorMsg = errMsg
	dp.UpdateTime = gtime.Now()

	pm.progressMap.Store(taskID, dp)
	updateRecord(dp)
	return nil
}

func (pm *ProgressManager) ListAllTasks(limit int) []*DownloadProgress {
	var tasks []*DownloadProgress

	pm.progressMap.Range(func(key, value interface{}) bool {
		tasks = append(tasks, value.(*DownloadProgress))
		return true
	})

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartTime.Timestamp() > tasks[j].StartTime.Timestamp()
	})

	if limit > 0 && limit < len(tasks) {
		tasks = tasks[:limit]
	}

	return tasks
}

func (pm *ProgressManager) CleanupTasks(olderThan time.Duration) {
	now := gtime.Timestamp()
	threshold := now - int64(olderThan.Seconds())

	pm.progressMap.Range(func(key, value interface{}) bool {
		dp := value.(*DownloadProgress)
		if dp.UpdateTime.Timestamp() < threshold {
			pm.progressMap.Delete(key)
		}
		return true
	})
}

func updateRecord(dp *DownloadProgress) {
	dao.DownloadRecord.Ctx(context.Background()).Where(dao.DownloadRecord.Columns().TaskId, dp.TaskID).Update(do.DownloadRecord{
		Status:     dp.Status,
		UpdateTime: dp.UpdateTime,
		ErrorMsg:   dp.ErrorMsg,
	})
	broadcastProgress()
}

func newRecord(dp *DownloadProgress) {
	dao.DownloadRecord.Ctx(context.Background()).Insert(&do.DownloadRecord{
		TaskId:     dp.TaskID,
		Title:      dp.Title,
		Status:     string(DownloadStatusPending),
		Output:     dp.OutputPath,
		StartTime:  dp.StartTime,
		UpdateTime: dp.UpdateTime,
	})
	broadcastProgress()
}

func broadcastProgress() {
	msg := gjson.MustEncodeString(g.Map{
		"event": "download",
	})
	sse.BroadcastMessage(msg)
}
