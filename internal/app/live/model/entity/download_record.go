// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DownloadRecord is the golang structure for table download_record.
type DownloadRecord struct {
	Id         int64       `json:"id"         orm:"id"          description:"主键 ID"`
	Title      string      `json:"title"      orm:"title"       description:"任务名称"`
	TaskId     string      `json:"taskId"     orm:"task_id"     description:"任务 ID"`
	Status     string      `json:"status"     orm:"status"      description:"任务状态"`
	Output     string      `json:"output"     orm:"output"      description:"输出路径"`
	ErrorMsg   string      `json:"errorMsg"   orm:"error_msg"   description:"错误信息"`
	StartTime  *gtime.Time `json:"startTime"  orm:"start_time"  description:"开始时间"`
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"结束时间"`
}
