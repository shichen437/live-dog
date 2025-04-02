// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaParse is the golang structure of table media_parse for DAO operations like Where/Data.
type MediaParse struct {
	g.Meta         `orm:"table:media_parse, do:true"`
	Id             interface{} // 媒体解析主键 ID
	Platform       interface{} // 平台
	Referer        interface{} // 来源
	Author         interface{} // 作者名称
	AuthorUid      interface{} // 作者 UID
	Desc           interface{} // 媒体描述
	MediaId        interface{} // 媒体 ID
	Type           interface{} // 媒体类型
	VideoUrl       interface{} // 视频 url
	VideoCoverUrl  interface{} // 视频封面 url
	VideoData      interface{} // 视频数据
	MusicUrl       interface{} // 音乐 url
	MusicCoverUrl  interface{} // 音乐封面 url
	ImagesUrl      interface{} // 图集 url
	ImagesCoverUrl interface{} // 图集封面 url
	CreateTime     *gtime.Time // 创建时间
}
