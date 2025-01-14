// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/shichen437/live-dog/internal/app/live/dao/internal"
)

// internalPushChannelEmailDao is internal type for wrapping internal DAO implements.
type internalPushChannelEmailDao = *internal.PushChannelEmailDao

// pushChannelEmailDao is the data access object for table push_channel_email.
// You can define custom methods on it to extend its functionality as you wish.
type pushChannelEmailDao struct {
	internalPushChannelEmailDao
}

var (
	// PushChannelEmail is globally public accessible object for table push_channel_email operations.
	PushChannelEmail = pushChannelEmailDao{
		internal.NewPushChannelEmailDao(),
	}
)

// Fill with you ideas below.
