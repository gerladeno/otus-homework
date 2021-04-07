package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
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
	counter uint64
	mu      sync.RWMutex
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
	var counter sql.NullInt64
	err = db.Get(&counter, "SELECT max(id) + 1 from events")
	if err != nil {
		log.Warn("can't get max id, most likely something is wrong, data loss possible. Continue with id = 0")
	}
	return &Storage{db: db, log: log, counter: uint64(counter.Int64)}, nil
}

func (s *Storage) ReadEvent(ctx context.Context, id uint64) (*common.Event, error) {
	query := fmt.Sprintf(`SELECT * from events WHERE id = %d`, id)
	result := make([]common.Event, 0)
	rows, err := s.db.QueryxContext(ctx, query)
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warn("err closing rows: ", err)
		}
	}()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrNoSuchEvent
		}
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
	if len(result) == 0 {
		return nil, common.ErrNoSuchEvent
	}
	return &result[0], nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *common.Event) (uint64, error) {
	s.mu.RLock()
	id := s.counter
	s.mu.RUnlock()
	event.ID = id
	event.Created = time.Now()
	event.Updated = time.Now()
	query := fmt.Sprintf(`
INSERT INTO events (id, title, start_time, duration, invite_list, comment) VALUES (%d, '%s', '%s', %d, '%s', '%s')
`, id, event.Title, event.StartTime.Format(common.PgTimestampFmt), int(event.Duration.Seconds()), event.InviteList, event.Comment)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		s.log.Warn("failed to add event ", id)
		return 0, err
	}
	s.mu.Lock()
	s.counter++
	s.mu.Unlock()
	s.log.Trace("added event ", id)
	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uint64, event *common.Event) error {
	event.ID = id
	event.Updated = time.Now()
	query := fmt.Sprintf(`
UPDATE events SET (title, start_time, duration, invite_list, comment, created) = ('%s', '%s', %d, '%s', '%s', '%s')
WHERE id = %d
`, event.Title, event.StartTime.Format(common.PgTimestampFmt), int(event.Duration.Seconds()), event.InviteList, event.Comment, event.Created.Format(common.PgTimestampFmt), id)
	res, err := s.db.ExecContext(ctx, query)
	if err != nil {
		s.log.Warn("failed to edit event ", id)
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		s.log.Warn("failed to add event ", id)
		return err
	}
	if n == 0 {
		return common.ErrNoSuchEvent
	}
	s.log.Trace("modified event ", id)
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM events WHERE id = %d`, id)
	res, err := s.db.ExecContext(ctx, query)
	if err != nil {
		s.log.Warn("failed to remove event ", id)
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		s.log.Warn("failed to add event ", id)
		return err
	}
	if n == 0 {
		return common.ErrNoSuchEvent
	}
	s.log.Trace("removed event ", id)
	return nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]*common.Event, error) {
	query := `SELECT * from events`
	result := make([]*common.Event, 0)
	rows, err := s.db.QueryxContext(ctx, query)
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warn("err closing rows: ", err)
		}
	}()
	if err != nil {
		s.log.Warn("failed to get a list of events")
		return nil, err
	}
	var event common.Event
	for rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			s.log.Warn("failed to get a list of events")
			return nil, err
		}
		result = append(result, &event)
	}
	return result, nil
}
