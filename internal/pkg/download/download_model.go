package download

import "github.com/gogf/gf/v2/os/gtime"

type DownloadParams struct {
	Platform     string   `json:"platform"`     // 平台
	Title        string   `json:"title"`        // 视频标题
	Type         string   `json:"type"`         // 视频类型
	Referer      string   `json:"referer"`      // 视频页面地址
	CoverUrl     string   `json:"coverUrl"`     // 视频封面地址
	QualityDesc  string   `json:"qualityDesc"`  // 视频质量描述
	Url          string   `json:"url"`          // 视频（流）地址
	Mirrors      []string `json:"mirrors"`      // 视频镜像地址
	Codec        string   `json:"codec"`        // 视频编码格式
	AudioUrl     string   `json:"audioUrl"`     // 音频地址
	AudioMirrors []string `json:"audioMirrors"` // 音频镜像地址
	AudioCodec   string   `json:"audioCodec"`   // 音频编码格式
	ImageUrls    []string `json:"imageUrls"`    // 图集
}

type DownloadResult struct {
	TaskID     string `json:"taskID"`
	OutputPath string `json:"outputPath"`
}

type DownloadStatus string

const (
	DownloadStatusPending     DownloadStatus = "pending"     // 待下载
	DownloadStatusRunning     DownloadStatus = "running"     // 下载中
	DownloadStatusConverting  DownloadStatus = "converting"  // 转换中
	DownloadStatusCompleted   DownloadStatus = "completed"   // 已完成
	DownloadStatusPartSucceed DownloadStatus = "partSucceed" // 部分成功
	DownloadStatusError       DownloadStatus = "error"       // 下载失败
)

type DownloadProgress struct {
	ErrorMsg   string         `json:"error"`
	OutputPath string         `json:"outputPath"`
	Progress   float64        `json:"progress"`
	StartTime  *gtime.Time    `json:"startTime"`
	Status     DownloadStatus `json:"status"`
	TaskID     string         `json:"taskID"`
	Title      string         `json:"title"`
	TotalSize  int64          `json:"totalSize"`
	UpdateTime *gtime.Time    `json:"updateTime"`
}
