// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AuthorInfoHistory is the golang structure for table author_info_history.
type AuthorInfoHistory struct {
	Id                 int64       `json:"id"                 orm:"id"                   description:"主键 ID"`
	AuthorId           int         `json:"authorId"           orm:"author_id"            description:"作者 ID"`
	LastFollowerCount  int64       `json:"lastFollowerCount"  orm:"last_follower_count"  description:"上次粉丝数量"`
	LastFollowingCount int         `json:"lastFollowingCount" orm:"last_following_count" description:"上次关注数量"`
	Num                int         `json:"num"                orm:"num"                  description:"涨粉数量"`
	Day                string      `json:"day"                orm:"day"                  description:"记录日期"`
	CreateTime         *gtime.Time `json:"createTime"         orm:"create_time"          description:"创建时间"`
	UpdateTime         *gtime.Time `json:"updateTime"         orm:"update_time"          description:"修改时间"`
}
