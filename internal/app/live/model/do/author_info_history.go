// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AuthorInfoHistory is the golang structure of table author_info_history for DAO operations like Where/Data.
type AuthorInfoHistory struct {
	g.Meta             `orm:"table:author_info_history, do:true"`
	Id                 interface{} // 主键 ID
	AuthorId           interface{} // 作者 ID
	LastFollowerCount  interface{} // 上次粉丝数量
	LastFollowingCount interface{} // 上次关注数量
	Num                interface{} // 涨粉数量
	Day                interface{} // 记录日期
	CreateTime         *gtime.Time // 创建时间
	UpdateTime         *gtime.Time // 修改时间
}
