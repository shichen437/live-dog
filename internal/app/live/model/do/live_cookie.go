// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveCookie is the golang structure of table live_cookie for DAO operations like Where/Data.
type LiveCookie struct {
	g.Meta     `orm:"table:live_cookie, do:true"`
	Id         interface{} // ID
	Platform   interface{} // 平台
	Cookie     interface{} // cookie
	Remark     interface{} // 备注
	CreateTime *gtime.Time // 创建时间
	ActionTime *gtime.Time // 更新时间
}
