package download

import (
	"bytes"
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

var (
	downloaderBuilders sync.Map
)

type Downloader interface {
	DownMediaFile(ctx context.Context) (*DownloadResult, error)
}

type DownloaderBuilder interface {
	Build(*DownloadParams) (Downloader, error)
}

func Register(domain string, b DownloaderBuilder) {
	downloaderBuilders.Store(domain, b)
}

func NewDownloader(params *DownloadParams) (Downloader, error) {
	builder, ok := downloaderBuilders.Load(params.Platform)
	if !ok {
		return nil, gerror.New("不支持当前平台下载")
	}
	return builder.(DownloaderBuilder).Build(params)
}

func GetOutputInfo(format, randomString string, isTemp bool, data interface{}) (string, string, error) {
	outputPath, err := GetOutputPath(isTemp, data)
	if err != nil {
		return "", "", gerror.New("failed to get outputPath")
	}
	filename, err := GetOutputFilename(format, randomString, outputPath, data)
	if err != nil {
		return "", "", gerror.New("failed to get outputFilename")
	}
	return outputPath, filename, nil
}

func GetOutputPath(isTemp bool, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	outTmpl := utils.GetDownloadPathTemplate(isTemp)
	err := outTmpl.Execute(buf, data)
	if err != nil {
		return "", gerror.New("failed to get outputPath template")
	}
	outputPath := buf.String()
	return outputPath, nil
}

func GetOutputFilename(format, randomString, outputPath string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	filenameTmpl := utils.GetDownloadFilenameTemplate(outputPath, format, randomString)
	err := filenameTmpl.Execute(buf, data)
	if err != nil {
		return "", gerror.New("failed to get downloadFilename template")
	}
	filename := buf.String()
	return filename, nil
}
