// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DownloadRecord is the golang structure of table download_record for DAO operations like Where/Data.
type DownloadRecord struct {
	g.Meta     `orm:"table:download_record, do:true"`
	Id         interface{} // 主键 ID
	Title      interface{} // 任务名称
	TaskId     interface{} // 任务 ID
	Status     interface{} // 任务状态
	Output     interface{} // 输出路径
	ErrorMsg   interface{} // 错误信息
	StartTime  *gtime.Time // 开始时间
	UpdateTime *gtime.Time // 结束时间
}
