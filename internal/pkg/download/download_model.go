package download

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
	DownloadStatusPending   DownloadStatus = "pending"
	DownloadStatusRunning   DownloadStatus = "running"
	DownloadStatusCompleted DownloadStatus = "completed"
	DownloadStatusError     DownloadStatus = "error"
)

type DownloadProgress struct {
	TaskID       string         `json:"taskID"`
	Status       DownloadStatus `json:"status"`
	Progress     float64        `json:"progress"`
	Speed        string         `json:"speed"`
	TotalSize    int64          `json:"totalSize"`
	Downloaded   int64          `json:"downloaded"`
	EstimatedETA string         `json:"estimatedETA"`
	Filename     string         `json:"filename"`
	Error        string         `json:"error,omitempty"`
	StartTime    int64          `json:"startTime"`
	UpdateTime   int64          `json:"updateTime"`
}
