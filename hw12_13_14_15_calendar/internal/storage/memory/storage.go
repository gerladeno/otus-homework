package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	mu      sync.RWMutex
	events  map[int64]common.Event
	counter int64
	log     *logrus.Logger
}

func New(log *logrus.Logger) *Storage {
	events := make(map[int64]common.Event)
	return &Storage{events: events, log: log}
}

func (s *Storage) CreateEvent(_ context.Context, event *common.Event) (int64, error) {
	event.Created = time.Now()
	event.Updated = time.Now()
	var id int64
	s.mu.Lock()
	{
		id = s.counter
		event.ID = id
		s.events[s.counter] = *event
		s.counter++
	}
	s.mu.Unlock()
	s.log.Trace("added event ", id)
	return id, nil
}

func (s *Storage) UpdateEvent(_ context.Context, id int64, event *common.Event) error {
	event.ID = id
	s.mu.Lock()
	{
		event.Created = s.events[id].Created
		event.Updated = time.Now()
		if _, ok := s.events[id]; !ok {
			return common.ErrNoSuchEvent
		}
		s.events[id] = *event
	}
	s.mu.Unlock()
	s.log.Trace("modified event ", id)
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int64) error {
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

func (s *Storage) ListEventsByDay(_ context.Context, date time.Time) ([]common.Event, error) {
	return s.listEvents(date, date.AddDate(0, 0, 1))
}

func (s *Storage) ListEventsByWeek(_ context.Context, date time.Time) ([]common.Event, error) {
	return s.listEvents(date, date.AddDate(0, 0, 7))
}

func (s *Storage) ListEventsByMonth(_ context.Context, date time.Time) ([]common.Event, error) {
	return s.listEvents(date, date.AddDate(0, 1, 0))
}

func (s *Storage) listEvents(fromDate, toDate time.Time) ([]common.Event, error) {
	events := make([]common.Event, 0)
	s.mu.RLock()
	for _, event := range s.events {
		if (event.StartTime.After(fromDate) || event.StartTime.Equal(fromDate)) && event.StartTime.Before(toDate) {
			events = append(events, event)
		}
	}
	s.mu.RUnlock()
	return events, nil
}

func (s *Storage) ListEventsToNotify(_ context.Context) (events []common.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.events {
		if time.Until(event.StartTime).Seconds() < float64(event.NotifyTime) && event.NotifyTime != 0 {
			events = append(events, event)
		}
	}
	return events, nil
}
