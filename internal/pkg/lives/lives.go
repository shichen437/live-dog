package lives

import (
	"net/url"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	builders sync.Map
)

type Live interface {
	GetInfo() (*RoomInfo, error)
	GetPlatform() string
	GetLiveId() int
	GetRefer() string
}

func Register(domain string, b Builder) {
	builders.Store(domain, b)
}

type Builder interface {
	Build(*url.URL, int) (Live, error)
}

func getBuilder(domain string) (Builder, error) {
	builder, ok := builders.Load(domain)
	if !ok {
		return nil, gerror.New("unknown domain")
	}
	return builder.(Builder), nil
}

func New(url *url.URL, liveId int) (live Live, err error) {
	builder, err := getBuilder(url.Hostname())
	if err != nil {
		return nil, gerror.New("not support this domain" + url.Hostname())
	}
	live, err = builder.Build(url, liveId)
	if err != nil {
		return nil, gerror.New("not support this domain" + url.Hostname())
	}
	return
}
