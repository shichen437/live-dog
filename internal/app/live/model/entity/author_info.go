// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AuthorInfo is the golang structure for table author_info.
type AuthorInfo struct {
	Id             int         `json:"id"             orm:"id"              description:"主键 ID"`
	UniqueId       string      `json:"uniqueId"       orm:"unique_id"       description:"加密 UID"`
	Platform       string      `json:"platform"       orm:"platform"        description:"作者平台"`
	Nickname       string      `json:"nickname"       orm:"nickname"        description:"作者昵称"`
	Signature      string      `json:"signature"      orm:"signature"       description:"签名"`
	AvatarUrl      string      `json:"avatarUrl"      orm:"avatar_url"      description:"头像url"`
	Ip             string      `json:"ip"             orm:"ip"              description:"作者 IP"`
	Refer          string      `json:"refer"          orm:"refer"           description:"来源"`
	FollowerCount  int64       `json:"followerCount"  orm:"follower_count"  description:"粉丝数量"`
	FollowingCount int         `json:"followingCount" orm:"following_count" description:"关注数量"`
	CreateTime     *gtime.Time `json:"createTime"     orm:"create_time"     description:"创建时间"`
	UpdateTime     *gtime.Time `json:"updateTime"     orm:"update_time"     description:"修改时间"`
}
