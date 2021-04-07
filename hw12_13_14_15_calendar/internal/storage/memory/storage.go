package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	mu      sync.Mutex
	events  map[uint64]*common.Event
	counter uint64
	log     *logrus.Logger
}

func New(log *logrus.Logger) *Storage {
	events := make(map[uint64]*common.Event)
	return &Storage{events: events, log: log}
}

func (s *Storage) ReadEvent(_ context.Context, id uint64) (*common.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if elem, ok := s.events[id]; ok {
		return elem, nil
	}
	return nil, common.ErrNoSuchEvent
}

func (s *Storage) CreateEvent(_ context.Context, event *common.Event) (uint64, error) {
	event.Created = time.Now()
	event.Updated = time.Now()
	var id uint64
	s.mu.Lock()
	{
		id = s.counter
		event.ID = id
		s.events[s.counter] = event
		s.counter++
	}
	s.mu.Unlock()
	s.log.Trace("added event ", id)
	return id, nil
}

func (s *Storage) UpdateEvent(_ context.Context, id uint64, event *common.Event) error {
	event.ID = id
	s.mu.Lock()
	{
		event.Created = s.events[id].Created
		event.Updated = time.Now()
		if _, ok := s.events[id]; !ok {
			return common.ErrNoSuchEvent
		}
		s.events[id] = event
	}
	s.mu.Unlock()
	s.log.Trace("modified event ", id)
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id uint64) error {
	if _, ok := s.events[id]; !ok {
		return common.ErrNoSuchEvent
	}
	s.mu.Lock()
	{
		delete(s.events, id)
	}
	s.mu.Unlock()
	s.log.Trace("removed event ", id)
	return nil
}

func (s *Storage) ListEvents(_ context.Context) ([]*common.Event, error) {
	events := make([]*common.Event, 0)
	s.mu.Lock()
	for _, event := range s.events {
		events = append(events, event)
	}
	s.mu.Unlock()
	return events, nil
}
