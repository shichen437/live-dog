package bilibili

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"github.com/shichen437/live-dog/internal/pkg/download"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

type BilibiliDownloader struct {
	downloadParams download.DownloadParams
}

type builder struct{}

func (b *builder) Build(params *download.DownloadParams) (download.Downloader, error) {
	return &BilibiliDownloader{
		downloadParams: *params,
	}, nil
}

func init() {
	download.Register(platform, &builder{})
}

func (d *BilibiliDownloader) DownMediaFile(ctx context.Context) (*download.DownloadResult, error) {
	taskID := uuid.New().String()
	// 视频下载
	if d.downloadParams.Type == "video" {
		return d.downloadVideo(ctx, taskID)
	}
	return nil, gerror.New("不支持当前类型下载")
}

func (d *BilibiliDownloader) downloadVideo(ctx context.Context, taskID string) (*download.DownloadResult, error) {
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
		outputTempPath, err := download.GetOutputPath(true, d.downloadParams)
		if err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		if err := os.MkdirAll(outputTempPath, os.ModePerm); err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		pm.UpdateProgress(taskID, download.DownloadStatusRunning)
		tempVideo, err := download.GetOutputFilename("m4s", utils.GenRandomString(6, randomStr), outputTempPath, d.downloadParams)
		if err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		_, err = utils.GetDownloadFile(downloadCtx, &utils.DownloadFileRequest{
			Url:       d.downloadParams.Url,
			CookieMap: utils.GetCookieMap(platform, d.downloadParams.Referer),
			Refer:     d.downloadParams.Referer,
			Filename:  tempVideo,
		})
		if err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		tempAudio, err := download.GetOutputFilename("m4s", utils.GenRandomString(6, randomStr), outputTempPath, d.downloadParams)
		if err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		defer func() {
			os.Remove(tempVideo)
			os.Remove(tempAudio)
		}()
		_, err = utils.GetDownloadFile(downloadCtx, &utils.DownloadFileRequest{
			Url:       d.downloadParams.AudioUrl,
			CookieMap: utils.GetCookieMap(platform, d.downloadParams.Referer),
			Refer:     d.downloadParams.Referer,
			Filename:  tempAudio,
		})
		if err != nil {
			g.Log().Error(downloadCtx, err)
			pm.SetError(taskID, err.Error())
			return
		}
		// 合并
		filename, err := download.GetOutputFilename("mp4", utils.GenRandomString(6, randomStr), outputPath, d.downloadParams)
		_, err = utils.NewFFmpegBuilder().
			Input(tempVideo).
			Input(tempAudio).
			CopyCodec().
			FastStart().
			AddDefaultThreads().
			Output(filename).
			Execute(downloadCtx)
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
