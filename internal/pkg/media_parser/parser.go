package media_parser

import (
	"context"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/shichen437/live-dog/internal/pkg/utils"
)

var (
	parserBuilders sync.Map
)

type MeidaParser interface {
	ParseURL(ctx context.Context) (*MediaInfo, error)
}

func Register(domain string, b MediaParserBuilder) {
	parserBuilders.Store(domain, b)
}

type MediaParserBuilder interface {
	Build(string) (MeidaParser, error)
}

func getBuilder(domain string) (MediaParserBuilder, error) {
	builder, ok := parserBuilders.Load(domain)
	if !ok {
		return nil, gerror.New("unknown domain")
	}
	return builder.(MediaParserBuilder), nil
}

func NewParser(url string) (MeidaParser, error) {
	p, reqUrl := parseSourceUrl(url)
	if p == "" {
		return nil, gerror.New("不支持的平台链接: " + url)
	}
	builder, err := getBuilder(p)
	if err != nil {
		return nil, gerror.New("not support this domain: " + url)
	}
	downloader, err := builder.Build(reqUrl)
	if err != nil {
		return nil, gerror.New("failed to build downloader for domain: " + url)
	}
	return downloader, nil
}

func parseSourceUrl(shareURL string) (platform, reqUrl string) {
	url := utils.FindFirstMatch(shareURL, BaseReg)
	if shareURL == "" || url == "" {
		return "", ""
	}

	for _, platform := range platformSet {
		if strings.Contains(url, platform) {
			return platform, url
		}
	}
	return "", ""
}
