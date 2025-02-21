package file_manage

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "github.com/shichen437/live-dog/api/v1/live"
	"github.com/shichen437/live-dog/internal/app/live/model"
	"github.com/shichen437/live-dog/internal/app/live/service"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

func init() {
	service.RegisterFileManage(New())
}

func New() *sFileManage {
	return &sFileManage{}
}

type sFileManage struct {
}

func (f *sFileManage) List(ctx context.Context, req *v1.GetFileInfoListReq) (res *v1.GetFileInfoListRes, err error) {
	res = &v1.GetFileInfoListRes{}
	base, err := filepath.Abs(utils.Output)
	utils.WriteErrLogT(ctx, err, "无效目录")
	absPath, err := filepath.Abs(filepath.Join(base, req.Path))
	utils.WriteErrLogT(ctx, err, "无效路径")
	if !strings.HasPrefix(absPath, base) {
		return nil, utils.TError(ctx, "无效路径")
	}
	files, err := os.ReadDir(absPath)
	utils.WriteErrLogT(ctx, err, "无效路径")
	if len(files) == 0 || err != nil {
		return
	}
	var list []*model.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil || isHiddenFile(info) {
			continue
		}
		if req.Filename != "" && !matchPattern(file.Name(), req.Filename) {
			continue
		}
		list = append(list, &model.FileInfo{
			Filename:     file.Name(),
			IsFolder:     file.IsDir(),
			Size:         info.Size(),
			LastModified: info.ModTime().Local().UnixMilli(),
		})
	}
	res.Rows = list
	return
}

func (f *sFileManage) Delete(ctx context.Context, req *v1.DeleteFileInfoReq) (res *v1.DeleteFileInfoRes, err error) {
	res = &v1.DeleteFileInfoRes{}
	if len(req.Filenames) == 0 {
		return
	}
	base, err := filepath.Abs(utils.Output)
	utils.WriteErrLogT(ctx, err, "无效目录")
	absPath, err := filepath.Abs(filepath.Join(base, req.Path))
	utils.WriteErrLogT(ctx, err, "无效路径")
	if !strings.HasPrefix(absPath, base) {
		return nil, utils.TError(ctx, "无效路径")
	}
	for _, filename := range req.Filenames {
		err = os.RemoveAll(filepath.Join(absPath, filename))
		utils.WriteErrLogT(ctx, err, "删除失败")
	}
	return
}

func (f *sFileManage) Play(ctx context.Context, req *v1.GetFilePlayReq) (res *v1.GetFilePlayRes, err error) {
	if req.Path == "" {
		err = gerror.New("路径不能为空")
		return
	}
	filepath := gfile.Join(utils.GetOutputPath(), req.Path)
	if !gfile.Exists(filepath) {
		err = gerror.New("文件不存在")
		return
	}
	fileInfo, err := os.Stat(filepath)
	if err != nil || fileInfo == nil || fileInfo.IsDir() {
		return
	}

	contentType := "video/mp4"
	ext := path.Ext(filepath)
	switch ext {
	case ".flv":
		contentType = "video/x-flv"
	case ".aac":
		contentType = "audio/aac"
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	}

	r := g.RequestFromCtx(ctx)
	w := r.Response
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", gconv.String(fileInfo.Size()))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "no-cache")

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 流式传输文件
	if _, err = io.Copy(w.Writer, file); err != nil {
		return nil, err
	}

	return
}

func isHiddenFile(file fs.FileInfo) bool {
	if file.IsDir() {
		return false
	}
	return strings.HasPrefix(file.Name(), ".")
}

func matchPattern(filename, pattern string) bool {
	return strings.Contains(filename, pattern)
}
