package sqlstorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	db      *sqlx.DB
	log     *logrus.Logger
	mu      sync.RWMutex
	events  map[uint64]common.Event
	counter uint64
}

func New(log *logrus.Logger, dsn string) (*Storage, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	var counter uint64
	err = db.Get(&counter, "SELECT max(id) + 1 from events")
	if err != nil {
		log.Warn("can't get max id, most likely something is wrong, data loss possible. Continue with id = 0")
	}
	events := make(map[uint64]common.Event)
	tmp, err := getEvents(context.Background(), log, db)
	if err != nil {
		log.Warn("failed to get a list of events, proceeding with empty list, probable data loss")
	}
	for _, event := range tmp {
		events[event.ID] = event
	}
	return &Storage{db: db, log: log, events: events, counter: counter}, nil
}

func (s *Storage) GetEvent(id uint64) (*common.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if elem, ok := s.events[id]; ok {
		return &elem, nil
	}
	return nil, common.ErrNoSuchEvent
}

func (s *Storage) AddEvent(ctx context.Context, event common.Event) (uint64, error) {
	s.mu.Lock()
	id := s.counter
	s.mu.Unlock()
	event.ID = id
	event.Created = time.Now()
	event.Updated = time.Now()
	query := fmt.Sprintf(`
INSERT INTO events (id, title, start_time, duration, invite_list, comment) VALUES (%d, '%s', '%s', %d, '%s', '%s')
`, id, event.Title, event.StartTime.Format(common.PgTimeStampFmt), int(event.Duration.Seconds()), event.InviteList, event.Comment)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		s.log.Warn("failed to add event ", id)
		return 0, err
	}
	s.mu.Lock()
	s.events[s.counter] = event
	s.counter++
	s.mu.Unlock()
	s.log.Trace("added event ", id)
	return id, nil
}

func (s *Storage) EditEvent(ctx context.Context, id uint64, event common.Event) error {
	event.ID = id
	s.mu.Lock()
	event.Created = s.events[id].Created
	event.Updated = time.Now()
	query := fmt.Sprintf(`
UPDATE events SET (title, start_time, duration, invite_list, comment, created) = ('%s', '%s', %d, '%s', '%s', '%s')
WHERE id = %d
`, event.Title, event.StartTime.Format(common.PgTimeStampFmt), int(event.Duration.Seconds()), event.InviteList, event.Comment, event.Created.Format(common.PgTimeStampFmt), id)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		s.log.Warn("failed to edit event ", id)
		return err
	}
	if _, ok := s.events[id]; !ok {
		return common.ErrNoSuchEvent
	}
	s.events[id] = event
	s.mu.Unlock()
	s.log.Trace("modified event ", id)
	return nil
}

func (s *Storage) RemoveEvent(ctx context.Context, id uint64) error {
	if _, ok := s.events[id]; !ok {
		return common.ErrNoSuchEvent
	}
	query := fmt.Sprintf(`DELETE FROM events WHERE id = %d`, id)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		s.log.Warn("failed to remove event ", id)
		return err
	}
	s.mu.Lock()
	delete(s.events, id)
	s.mu.Unlock()
	s.log.Trace("removed event ", id)
	return nil
}

func (s *Storage) ListEvents() ([]common.Event, error) {
	events := make([]common.Event, 0)
	s.mu.Lock()
	for _, event := range s.events {
		events = append(events, event)
	}
	s.mu.Unlock()
	return events, nil
}

func getEvents(ctx context.Context, log *logrus.Logger, db *sqlx.DB) ([]common.Event, error) {
	query := `SELECT * from events`
	result := make([]common.Event, 0)
	rows, err := db.QueryxContext(ctx, query)
	defer func() {
		if err := rows.Close(); err != nil {
			log.Warn("err closing rows: ", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	var event common.Event
	for rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	return result, nil
}
