package douyin

import (
	"context"
	"os"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"github.com/shichen437/live-dog/internal/pkg/download"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

type DouyinDownloader struct {
	downloadParams download.DownloadParams
}

type builder struct{}

func (b *builder) Build(params *download.DownloadParams) (download.Downloader, error) {
	return &DouyinDownloader{
		downloadParams: *params,
	}, nil
}

func init() {
	download.Register(platform, &builder{})
}

func (d *DouyinDownloader) DownMediaFile(ctx context.Context) (*download.DownloadResult, error) {
	taskID := uuid.New().String()
	// 视频下载
	if d.downloadParams.Type == "video" {
		return d.downloadVideo(ctx, taskID)
	}
	// 图集下载
	if d.downloadParams.Type == "note" {
		return d.downloadNote(ctx, taskID)
	}
	return nil, gerror.New("不支持当前类型下载")
}

func (d *DouyinDownloader) downloadVideo(ctx context.Context, taskID string) (*download.DownloadResult, error) {
	outputPath, filename, err := download.GetOutputInfo("mp4", utils.GenRandomString(6, randomStr), false, d.downloadParams)
	pm := download.GetProgressManager()
	pm.CreateTask(taskID, d.downloadParams.Title, filename)
	if err != nil {
		g.Log().Error(ctx, err)
		pm.SetError(taskID, err.Error())
		return nil, gerror.New("获取输出路径失败")
	}
	go func() {
		downloadCtx := context.Background()
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		_, err := utils.GetDownloadFile(downloadCtx, &utils.DownloadFileRequest{
			Url:      d.downloadParams.Url,
			Filename: filename,
		})
		if err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		pm.SetCompleted(taskID)
	}()
	return &download.DownloadResult{
		TaskID:     taskID,
		OutputPath: outputPath,
	}, nil
}

func (d *DouyinDownloader) downloadNote(ctx context.Context, taskID string) (*download.DownloadResult, error) {
	rs := utils.GenRandomString(6, randomStr)
	outputPath, err := download.GetOutputPath(false, d.downloadParams)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("获取输出路径失败")
	}
	pm := download.GetProgressManager()
	pm.CreateTask(taskID, d.downloadParams.Title, outputPath)
	go func() {
		downloadCtx := context.Background()
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		var flag = 0
		pm.UpdateProgress(taskID, download.DownloadStatusRunning)
		for idx, imageUrl := range d.downloadParams.ImageUrls {
			var rse string
			if idx == 0 {
				rse = rs
			} else {
				rse = rs + "-" + strconv.Itoa(idx)
			}
			filename, err := download.GetOutputFilename("jpg", rse, outputPath, d.downloadParams)
			g.Log().Info(downloadCtx, filename)
			if err != nil {
				g.Log().Error(downloadCtx, err)
				continue
			}

			_, err = utils.GetDownloadFile(downloadCtx, &utils.DownloadFileRequest{
				Url:      imageUrl,
				Filename: filename,
			})
			if err != nil {
				g.Log().Error(downloadCtx, err)
				continue
			}
			flag++
		}
		if flag == len(d.downloadParams.ImageUrls) {
			pm.SetCompleted(taskID)
		} else if flag == 0 {
			pm.SetError(taskID, "下载失败")
		} else {
			pm.SetPartCompleted(taskID, "部分请求未成功，请检查日志")
		}

	}()
	return &download.DownloadResult{
		TaskID:     taskID,
		OutputPath: outputPath,
	}, nil
}
