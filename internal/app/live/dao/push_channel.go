// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/shichen437/live-dog/internal/app/live/dao/internal"
)

// internalPushChannelDao is internal type for wrapping internal DAO implements.
type internalPushChannelDao = *internal.PushChannelDao

// pushChannelDao is the data access object for table push_channel.
// You can define custom methods on it to extend its functionality as you wish.
type pushChannelDao struct {
	internalPushChannelDao
}

var (
	// PushChannel is globally public accessible object for table push_channel operations.
	PushChannel = pushChannelDao{
		internal.NewPushChannelDao(),
	}
)

// Fill with you ideas below.
