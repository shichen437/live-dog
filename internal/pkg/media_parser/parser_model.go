package media_parser

var (
	platformSet = []string{"douyin"}
	BaseReg     = `http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`
)

type MediaInfo struct {
	Platform      string `json:"platform"`
	VideoID       string `json:"video_id"`
	Author        string `json:"author"`
	AuthorUid     string `json:"author_uid"`
	Desc          string `json:"desc"`
	Type          string `json:"type"`
	VideoUrl      string `json:"video_url"`
	VideoCoverUrl string `json:"video_cover_url"`
	MusicUrl      string `json:"music_url"`
	MusicCoverUrl string `json:"music_cover_url"`
}

type DownloadResult struct {
	VideoID string
	Video   []byte
}
