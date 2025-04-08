// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AuthorInfo is the golang structure of table author_info for DAO operations like Where/Data.
type AuthorInfo struct {
	g.Meta         `orm:"table:author_info, do:true"`
	Id             interface{} // 主键 ID
	UniqueId       interface{} // 加密 UID
	Platform       interface{} // 作者平台
	Nickname       interface{} // 作者昵称
	Signature      interface{} // 签名
	AvatarUrl      interface{} // 头像url
	Ip             interface{} // 作者 IP
	Refer          interface{} // 来源
	FollowerCount  interface{} // 粉丝数量
	FollowingCount interface{} // 关注数量
	CreateTime     *gtime.Time // 创建时间
	UpdateTime     *gtime.Time // 修改时间
}
