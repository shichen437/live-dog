package bilibili

type BilibiliVideoDetail struct {
	Aid        string
	Author     string
	AuthorUid  string
	Cid        string
	Desc       string
	VideoCover string
}

type BilibiliMediaData struct {
	Audios []*BilibiliMediaDetail `json:"audios"`
	Videos []*BilibiliMediaDetail `json:"videos"`
}

type BilibiliMediaDetail struct {
	Codec       string   `json:"codec"`
	Height      int      `json:"height"`
	Mirrors     []string `json:"mirrors"`
	Quality     int      `json:"quality"`
	QualityDesc string   `json:"quality_desc"`
	Url         string   `json:"url"`
	Width       int      `json:"width"`
}
