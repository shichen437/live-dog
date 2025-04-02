package download

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
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

func (pm *ProgressManager) CreateTask(taskID, filename string) *DownloadProgress {
	progress := &DownloadProgress{
		TaskID:     taskID,
		Status:     DownloadStatusPending,
		Progress:   0,
		Speed:      "0 B/s",
		Filename:   filename,
		StartTime:  time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	pm.progressMap.Store(taskID, progress)
	return progress
}

func (pm *ProgressManager) GetProgress(taskID string) (*DownloadProgress, error) {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return nil, gerror.Newf("任务 %s 不存在", taskID)
	}
	return value.(*DownloadProgress), nil
}

func (pm *ProgressManager) UpdateProgress(taskID string, status DownloadStatus, progress float64, speed string, downloaded, totalSize int64, eta string) error {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return gerror.Newf("任务 %s 不存在", taskID)
	}

	dp := value.(*DownloadProgress)
	dp.Status = status
	dp.Progress = progress
	dp.Speed = speed
	dp.Downloaded = downloaded
	dp.TotalSize = totalSize
	dp.EstimatedETA = eta
	dp.UpdateTime = time.Now().Unix()

	pm.progressMap.Store(taskID, dp)
	return nil
}

func (pm *ProgressManager) SetError(taskID, errMsg string) error {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return gerror.Newf("任务 %s 不存在", taskID)
	}

	dp := value.(*DownloadProgress)
	dp.Status = DownloadStatusError
	dp.Error = errMsg
	dp.UpdateTime = time.Now().Unix()

	pm.progressMap.Store(taskID, dp)
	return nil
}

func (pm *ProgressManager) SetCompleted(taskID string) error {
	value, ok := pm.progressMap.Load(taskID)
	if !ok {
		return gerror.Newf("任务 %s 不存在", taskID)
	}

	dp := value.(*DownloadProgress)
	dp.Status = DownloadStatusCompleted
	dp.Progress = 100
	dp.UpdateTime = time.Now().Unix()

	pm.progressMap.Store(taskID, dp)
	return nil
}

func (pm *ProgressManager) ListAllTasks() []*DownloadProgress {
	var tasks []*DownloadProgress

	pm.progressMap.Range(func(key, value interface{}) bool {
		tasks = append(tasks, value.(*DownloadProgress))
		return true
	})

	return tasks
}

func (pm *ProgressManager) CleanupTasks(olderThan time.Duration) {
	now := time.Now().Unix()
	threshold := now - int64(olderThan.Seconds())

	pm.progressMap.Range(func(key, value interface{}) bool {
		dp := value.(*DownloadProgress)
		if (dp.Status == DownloadStatusCompleted || dp.Status == DownloadStatusError) && dp.UpdateTime < threshold {
			pm.progressMap.Delete(key)
		}
		return true
	})
}
