// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/shichen437/live-dog/internal/app/admin/dao/internal"
)

// internalPushMessageLogDao is internal type for wrapping internal DAO implements.
type internalPushMessageLogDao = *internal.PushMessageLogDao

// pushMessageLogDao is the data access object for table push_message_log.
// You can define custom methods on it to extend its functionality as you wish.
type pushMessageLogDao struct {
	internalPushMessageLogDao
}

var (
	// PushMessageLog is globally public accessible object for table push_message_log operations.
	PushMessageLog = pushMessageLogDao{
		internal.NewPushMessageLogDao(),
	}
)

// Fill with you ideas below.
