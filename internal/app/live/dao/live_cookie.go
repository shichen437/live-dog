// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/shichen437/live-dog/internal/app/live/dao/internal"
)

// internalLiveCookieDao is internal type for wrapping internal DAO implements.
type internalLiveCookieDao = *internal.LiveCookieDao

// liveCookieDao is the data access object for table live_cookie.
// You can define custom methods on it to extend its functionality as you wish.
type liveCookieDao struct {
	internalLiveCookieDao
}

var (
	// LiveCookie is globally public accessible object for table live_cookie operations.
	LiveCookie = liveCookieDao{
		internal.NewLiveCookieDao(),
	}
)

// Fill with you ideas below.
