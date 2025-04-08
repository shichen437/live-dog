package media_parser

var (
	platformSet = []string{"douyin", "bilibili"}
	BaseReg     = `http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`
)

type MediaInfo struct {
	Platform       string `json:"platform"`
	VideoID        string `json:"video_id"`
	Refer          string `json:"refer"`
	Author         string `json:"author"`
	AuthorUid      string `json:"author_uid"`
	Desc           string `json:"desc"`
	Type           string `json:"type"`
	VideoUrl       string `json:"video_url"`
	VideoCoverUrl  string `json:"video_cover_url"`
	MusicUrl       string `json:"music_url"`
	MusicCoverUrl  string `json:"music_cover_url"`
	ImagesUrl      string `json:"images_url"`
	ImagesCoverUrl string `json:"images_cover_url"`
	VideoData      string `json:"video_data"`
}

type UserInfo struct {
	UniqueId       string `json:"unique_id"`
	Platform       string `json:"platform"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
	Signature      string `json:"signature"`
	IpLocation     string `json:"ip_location"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
	Refer          string `json:"refer"`
	CurrentDay     string `json:"current_day"`
}
