// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaParse is the golang structure for table media_parse.
type MediaParse struct {
	Id             int64       `json:"id"             orm:"id"               description:"媒体解析主键 ID"`
	Platform       string      `json:"platform"       orm:"platform"         description:"平台"`
	Author         string      `json:"author"         orm:"author"           description:"作者名称"`
	AuthorUid      string      `json:"authorUid"      orm:"author_uid"       description:"作者 UID"`
	Desc           string      `json:"desc"           orm:"desc"             description:"媒体描述"`
	MediaId        string      `json:"mediaId"        orm:"media_id"         description:"媒体 ID"`
	Type           string      `json:"type"           orm:"type"             description:"媒体类型"`
	VideoUrl       string      `json:"videoUrl"       orm:"video_url"        description:"视频 url"`
	VideoCoverUrl  string      `json:"videoCoverUrl"  orm:"video_cover_url"  description:"视频封面 url"`
	VideoData      string      `json:"videoData"      orm:"video_data"       description:"视频数据"`
	MusicUrl       string      `json:"musicUrl"       orm:"music_url"        description:"音乐 url"`
	MusicCoverUrl  string      `json:"musicCoverUrl"  orm:"music_cover_url"  description:"音乐封面 url"`
	ImagesUrl      string      `json:"imagesUrl"      orm:"images_url"       description:"图集 url"`
	ImagesCoverUrl string      `json:"imagesCoverUrl" orm:"images_cover_url" description:"图集封面 url"`
	CreateTime     *gtime.Time `json:"createTime"     orm:"create_time"      description:"创建时间"`
}
