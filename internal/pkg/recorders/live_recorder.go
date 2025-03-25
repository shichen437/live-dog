package recorders

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shichen437/live-dog/internal/pkg/events"
	"github.com/shichen437/live-dog/internal/pkg/lives"
	"github.com/shichen437/live-dog/internal/pkg/parser"
	ffmpeg_parser "github.com/shichen437/live-dog/internal/pkg/parser/ffmpeg"
	"github.com/shichen437/live-dog/internal/pkg/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	begin uint32 = iota
	pending
	running
	stopped
)

type Recorder interface {
	Start(ctx context.Context) error
	Close()
	StartTime() time.Time
	GetStatus() (map[string]string, error)
}

type recorder struct {
	Live lives.Live

	ed         events.Dispatcher
	startTime  time.Time
	parser     parser.Parser
	parserLock *sync.RWMutex

	stop  chan struct{}
	state uint32
}

func NewRecorder(ctx context.Context, live lives.Live) (Recorder, error) {
	global := utils.GetGlobal(ctx)
	return &recorder{
		Live:       live,
		startTime:  time.Now(),
		ed:         global.EventDispatcher.(events.Dispatcher),
		state:      begin,
		stop:       make(chan struct{}),
		parserLock: new(sync.RWMutex),
	}, nil
}

func (r *recorder) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&r.state, begin, pending) {
		return nil
	}
	go r.run(ctx)
	r.ed.DispatchEvent(events.NewEvent("RecorderStart", r.Live))
	atomic.CompareAndSwapUint32(&r.state, pending, running)
	return nil
}

func (r *recorder) Close() {
	if !atomic.CompareAndSwapUint32(&r.state, running, stopped) {
		return
	}
	close(r.stop)
	if p := r.getParser(); p != nil {
		if err := p.Stop(); err != nil {
			g.Log().Error(gctx.New(), "failed to end recorder")
		}
	}
	r.ed.DispatchEvent(events.NewEvent("RecorderStop", r.Live))
}

func (r *recorder) StartTime() time.Time {
	return r.startTime
}

func (r *recorder) GetStatus() (map[string]string, error) {
	statusP, ok := r.getParser().(parser.StatusParser)
	if !ok {
		return nil, gerror.New("parser does not implement StatusParser")
	}
	return statusP.Status()
}

func (r *recorder) tryRecord(ctx context.Context) {
	var streamInfos []*lives.StreamUrlInfo
	var err error
	info, _ := r.Live.GetInfo()
	streamInfos = info.StreamInfos
	if len(streamInfos) == 0 {
		g.Log().Error(ctx, "failed to get stream info", err)
		time.Sleep(5 * time.Second)
		return
	}
	streamInfo := streamInfos[0]

	fileName, outputPath, err := r.getOutPathAndFilename(ctx, info)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = os.MkdirAll(outputPath, os.ModePerm); err != nil {
		g.Log().Error(ctx, "failed to create output path", err)
		return
	}
	parserCfg := map[string]string{
		"timeout_in_us": strconv.Itoa(100000000),
	}
	p, err := newParser(parserCfg)
	if err != nil {
		g.Log().Error(ctx, "failed to init parse", err)
		return
	}
	r.setAndCloseParser(p)
	r.startTime = time.Now()
	r.parser.ParseLiveStream(ctx, streamInfo, fileName, r.Live.GetRefer())
	removeEmptyFile(fileName)
}

func newParser(cfg map[string]string) (parser.Parser, error) {
	return parser.New(ffmpeg_parser.Name, cfg)
}

func removeEmptyFile(file string) {
	if stat, err := os.Stat(file); err == nil && stat.Size() == 0 {
		os.Remove(file)
	}
}

func (r *recorder) run(ctx context.Context) {
	for {
		select {
		case <-r.stop:
			return
		default:
			r.tryRecord(ctx)
		}
	}
}

func (r *recorder) getParser() parser.Parser {
	r.parserLock.RLock()
	defer r.parserLock.RUnlock()
	return r.parser
}

func (r *recorder) setAndCloseParser(p parser.Parser) {
	r.parserLock.Lock()
	defer r.parserLock.Unlock()
	if r.parser != nil {
		if err := r.parser.Stop(); err != nil {
			g.Log().Error(gctx.New(), "failed to end recorder", err)
		}
	}
	r.parser = p
}

func (r *recorder) getOutPathAndFilename(ctx context.Context, info *lives.RoomInfo) (string, string, error) {
	global := utils.GetGlobal(ctx)
	m, ok := global.ModelsMap[r.Live.GetLiveId()]
	var format string
	if ok && m.LiveManage.Format != "" {
		format = m.LiveManage.Format
	} else {
		format = "flv"
	}
	buf := new(bytes.Buffer)
	outTmpl := utils.GetOutputPathTemplate()
	err := outTmpl.Execute(buf, info)
	if err != nil {
		return "", "", gerror.New("failed to get outputPath template")
	}
	outputPath := buf.String()
	filenameTmpl := utils.GetFilenameTemplate(outputPath, format)
	buf.Reset()
	err = filenameTmpl.Execute(buf, info)
	if err != nil {
		return "", "", gerror.New("failed to get filename template")
	}
	filename := buf.String()
	return filename, outputPath, nil
}
