package recorders

import (
	"context"
	"sync"
	"time"

	"github.com/shichen437/live-dog/internal/pkg/events"
	"github.com/shichen437/live-dog/internal/pkg/interfaces"
	"github.com/shichen437/live-dog/internal/pkg/lives"
	"github.com/shichen437/live-dog/internal/pkg/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func NewManager(ctx context.Context) Manager {
	rm := &manager{
		savers: make(map[int]Recorder),
	}
	utils.GetGlobal(ctx).RecorderManager = rm
	return rm
}

type Manager interface {
	interfaces.Module
	AddRecorder(ctx context.Context, live lives.Live) error
	RemoveRecorder(ctx context.Context, liveId int) error
	RestartRecorder(ctx context.Context, liveId lives.Live) error
	GetRecorder(ctx context.Context, liveId int) (Recorder, error)
	HasRecorder(ctx context.Context, liveId int) bool
}

var (
	newRecorder = NewRecorder
)

type manager struct {
	lock   sync.RWMutex
	savers map[int]Recorder
}

func (m *manager) registryListener(ctx context.Context, ed events.Dispatcher) {
	ed.AddEventListener("LiveStart", events.NewEventListener(func(event *events.Event) {
		live := event.Object.(lives.Live)
		err := m.AddRecorder(ctx, live)
		if err != nil {
			g.Log().Error(ctx, "failed to add recorder")
			return
		}
		liveStartBiz(ctx, live.GetLiveId())
	}))

	ed.AddEventListener("NameChanged", events.NewEventListener(func(event *events.Event) {
		live := event.Object.(lives.Live)
		m.updateName(ctx, live)
	}))

	removeEvtListener := events.NewEventListener(func(event *events.Event) {
		live := event.Object.(lives.Live)
		if !m.HasRecorder(ctx, live.GetLiveId()) {
			return
		}
		err := m.RemoveRecorder(ctx, live.GetLiveId())
		if err != nil {
			g.Log().Error(ctx, "failed to remove recorder")
			return
		}
		liveEndBiz(ctx, live.GetLiveId())
	})
	ed.AddEventListener("LiveEnd", removeEvtListener)
	ed.AddEventListener("ListenStop", removeEvtListener)
}

func (m *manager) Start(ctx context.Context) error {
	global := utils.GetGlobal(ctx)
	if len(global.Lives) > 0 {
		global.WaitGroup.Add(1)
	}
	m.registryListener(ctx, global.EventDispatcher.(events.Dispatcher))
	g.Log().Info(ctx, "RecorderManager Started!")
	return nil
}

func (m *manager) Close(ctx context.Context) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for id, recorder := range m.savers {
		recorder.Close()
		delete(m.savers, id)
	}
	global := utils.GetGlobal(ctx)
	global.WaitGroup.Done()
	g.Log().Info(ctx, "RecorderManager Closed!")
}

func (m *manager) AddRecorder(ctx context.Context, live lives.Live) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.savers[live.GetLiveId()]; ok {
		return gerror.New("this live has a recorder")
	}
	recorder, err := newRecorder(ctx, live)
	if err != nil {
		return err
	}
	m.savers[live.GetLiveId()] = recorder

	return recorder.Start(ctx)
}

func (m *manager) RestartRecorder(ctx context.Context, live lives.Live) error {
	if err := m.RemoveRecorder(ctx, live.GetLiveId()); err != nil {
		return err
	}
	time.Sleep(20 * time.Second)
	if err := m.AddRecorder(ctx, live); err != nil {
		return err
	}
	return nil
}

func (m *manager) RemoveRecorder(ctx context.Context, liveId int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	recorder, ok := m.savers[liveId]
	if !ok {
		return gerror.New("this live has not a recorder")
	}
	recorder.Close()
	delete(m.savers, liveId)
	return nil
}

func (m *manager) GetRecorder(ctx context.Context, liveId int) (Recorder, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r, ok := m.savers[liveId]
	if !ok {
		return nil, gerror.New("this live has a recorder")
	}
	return r, nil
}

func (m *manager) HasRecorder(ctx context.Context, liveId int) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.savers[liveId]
	return ok
}
