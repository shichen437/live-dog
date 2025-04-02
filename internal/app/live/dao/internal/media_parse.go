// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaParseDao is the data access object for table media_parse.
type MediaParseDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns MediaParseColumns // columns contains all the column names of Table for convenient usage.
}

// MediaParseColumns defines and stores column names for table media_parse.
type MediaParseColumns struct {
	Id             string // 媒体解析主键 ID
	Platform       string // 平台
	Referer        string // 来源
	Author         string // 作者名称
	AuthorUid      string // 作者 UID
	Desc           string // 媒体描述
	MediaId        string // 媒体 ID
	Type           string // 媒体类型
	VideoUrl       string // 视频 url
	VideoCoverUrl  string // 视频封面 url
	VideoData      string // 视频数据
	MusicUrl       string // 音乐 url
	MusicCoverUrl  string // 音乐封面 url
	ImagesUrl      string // 图集 url
	ImagesCoverUrl string // 图集封面 url
	CreateTime     string // 创建时间
}

// mediaParseColumns holds the columns for table media_parse.
var mediaParseColumns = MediaParseColumns{
	Id:             "id",
	Platform:       "platform",
	Referer:        "referer",
	Author:         "author",
	AuthorUid:      "author_uid",
	Desc:           "desc",
	MediaId:        "media_id",
	Type:           "type",
	VideoUrl:       "video_url",
	VideoCoverUrl:  "video_cover_url",
	VideoData:      "video_data",
	MusicUrl:       "music_url",
	MusicCoverUrl:  "music_cover_url",
	ImagesUrl:      "images_url",
	ImagesCoverUrl: "images_cover_url",
	CreateTime:     "create_time",
}

// NewMediaParseDao creates and returns a new DAO object for table data access.
func NewMediaParseDao() *MediaParseDao {
	return &MediaParseDao{
		group:   "default",
		table:   "media_parse",
		columns: mediaParseColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MediaParseDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MediaParseDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MediaParseDao) Columns() MediaParseColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MediaParseDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MediaParseDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MediaParseDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
