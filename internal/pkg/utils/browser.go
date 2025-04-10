package utils

import (
	"fmt"
	"math/rand"
	"time"
)

type BrowserType string

const (
	Chrome  BrowserType = "chrome"
	Firefox BrowserType = "firefox"
	Edge    BrowserType = "edge"
	Safari  BrowserType = "safari"

	WinPlatform = "Win32"
	MacPlatform = "MacIntel"
)

func GenFingerPrint(browserType BrowserType) string {
	switch browserType {
	case Safari:
		return genFingerPrint(MacPlatform)
	default:
		return genFingerPrint(WinPlatform)
	}
}

func genFingerPrint(platform string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 随机生成各种尺寸参数
	innerWidth := r.Intn(897) + 1024
	innerHeight := r.Intn(313) + 768
	outerWidth := innerWidth + r.Intn(9) + 24
	outerHeight := innerHeight + r.Intn(16) + 75
	screenX := 0

	// 随机选择 screenY 值
	screenY := 0
	if r.Intn(2) == 1 {
		screenY = 30
	}

	sizeWidth := r.Intn(897) + 1024
	sizeHeight := r.Intn(313) + 768
	availWidth := r.Intn(641) + 1280
	availHeight := r.Intn(281) + 800

	// 构建指纹字符串
	fingerprint := fmt.Sprintf("%d|%d|%d|%d|%d|%d|0|0|%d|%d|%d|%d|%d|%d|24|24|%s",
		innerWidth, innerHeight, outerWidth, outerHeight,
		screenX, screenY, sizeWidth, sizeHeight,
		availWidth, availHeight, innerWidth, innerHeight, platform)

	return fingerprint
}
