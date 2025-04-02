package utils

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/mod/modfile"
)

// GetDefaultFFmpegPath 获取默认ffmpeg路径
func GetDefaultFFmpegPath() (string, error) {
	path, err := exec.LookPath("ffmpeg")
	if errors.Is(err, exec.ErrDot) {
		path, err = exec.LookPath("./ffmpeg")
	}
	return path, err
}

func GetOutputPath() string {
	return Output
}

func GetDownloadPath() string {
	return Output + "/download"
}

func GetTempDownloadPath() string {
	return "./resource/temp/download"
}

func GetFileUploadPath() string {
	return Upload
}

// 读取go.mod module
func GetGodModule() string {
	goModFilePathData, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	modFile, err := modfile.Parse("go.mod", goModFilePathData, nil)
	if err != nil {
		return ""
	}
	return modFile.Module.Mod.String()
}

// 格式化go文件
func FmtGoFile(path string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "go fmt "+path)
	} else {
		cmd = exec.Command("go", "fmt ", path)
	}
	if err := cmd.Start(); err != nil { // 运行命令
		log.Fatal(err)
	}
}

func IsTimeRange(st, et string) bool {
	startTime, endTime, ok := validAndGetTime(st, et)
	if !ok {
		return false
	}
	// 获取当前时间
	now := time.Now()
	currentTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	if startTime.After(endTime) {
		// 如果结束时间在开始时间之前，则认为是跨天
		return (currentTime.After(startTime) && currentTime.Before(endTime.AddDate(0, 0, 1))) ||
			(currentTime.Before(startTime) && currentTime.Before(endTime))
	} else {
		// 正常情况
		return currentTime.After(startTime) && currentTime.Before(endTime)
	}
}

func IsWithinCustomTimes(st, et string, step int) bool {
	startTime, endTime, ok := validAndGetTime(st, et)
	if !ok {
		return true
	}
	diff := startTime.Sub(endTime).Abs()
	return diff <= time.Duration(step)*time.Minute
}

func validAndGetTime(st string, et string) (time.Time, time.Time, bool) {
	if st == "" || et == "" {
		return time.Time{}, time.Time{}, false
	}
	sArr := strings.Split(st, ":")
	eArr := strings.Split(et, ":")
	if len(sArr) != 2 || len(eArr) != 2 {
		return time.Time{}, time.Time{}, false
	}
	sh, err1 := strconv.Atoi(sArr[0])
	sm, err2 := strconv.Atoi(sArr[1])
	eh, err3 := strconv.Atoi(eArr[0])
	em, err4 := strconv.Atoi(eArr[1])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return time.Time{}, time.Time{}, false
	}
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), sh, sm, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), eh, em, 0, 0, now.Location())
	return startTime, endTime, true
}

func IsDocker() bool {
	p, err := g.Cfg().Get(gctx.GetInitCtx(), "project.platform")
	return err == nil && p != nil && p.String() == "docker"
}
