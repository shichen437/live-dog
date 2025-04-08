// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AuthorInfoHistoryDao is the data access object for table author_info_history.
type AuthorInfoHistoryDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns AuthorInfoHistoryColumns // columns contains all the column names of Table for convenient usage.
}

// AuthorInfoHistoryColumns defines and stores column names for table author_info_history.
type AuthorInfoHistoryColumns struct {
	Id                 string // 主键 ID
	AuthorId           string // 作者 ID
	LastFollowerCount  string // 上次粉丝数量
	LastFollowingCount string // 上次关注数量
	Num                string // 涨粉数量
	Day                string // 记录日期
	CreateTime         string // 创建时间
	UpdateTime         string // 修改时间
}

// authorInfoHistoryColumns holds the columns for table author_info_history.
var authorInfoHistoryColumns = AuthorInfoHistoryColumns{
	Id:                 "id",
	AuthorId:           "author_id",
	LastFollowerCount:  "last_follower_count",
	LastFollowingCount: "last_following_count",
	Num:                "num",
	Day:                "day",
	CreateTime:         "create_time",
	UpdateTime:         "update_time",
}

// NewAuthorInfoHistoryDao creates and returns a new DAO object for table data access.
func NewAuthorInfoHistoryDao() *AuthorInfoHistoryDao {
	return &AuthorInfoHistoryDao{
		group:   "default",
		table:   "author_info_history",
		columns: authorInfoHistoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AuthorInfoHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AuthorInfoHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AuthorInfoHistoryDao) Columns() AuthorInfoHistoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AuthorInfoHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AuthorInfoHistoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AuthorInfoHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
