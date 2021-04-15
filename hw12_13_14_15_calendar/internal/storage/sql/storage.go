package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	db      *sqlx.DB
	log     *logrus.Logger
	counter int64
	mu      sync.RWMutex
}

func New(ctx context.Context, log *logrus.Logger, dsn string) (*Storage, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	var counter sql.NullInt64
	err = db.Get(&counter, "SELECT max(id) + 1 from events")
	if err != nil {
		return nil, err
	}
	return &Storage{db: db, log: log, counter: counter.Int64}, nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *common.Event) (int64, error) {
	s.mu.RLock()
	id := s.counter
	s.mu.RUnlock()
	event.ID = id
	event.Created = time.Now()
	event.Updated = time.Now()
	query := fmt.Sprintf(`
INSERT INTO events (id, title, start_time, duration, description, owner, notify_time) VALUES (%d, '%s', '%s', %d, '%s', '%d', '%d')
`, id, event.Title, event.StartTime.Format(common.PgTimestampFmt), event.Duration, event.Description, event.Owner, event.NotifyTime)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	s.mu.Lock()
	s.counter++
	s.mu.Unlock()
	s.log.Trace("added event ", id)
	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id int64, event *common.Event) error {
	event.ID = id
	event.Updated = time.Now()
	query := fmt.Sprintf(`
UPDATE events SET (title, start_time, duration, description, owner, notify_time, updated) = ('%s', '%s', %d, '%s', '%d', '%d', '%s')
WHERE id = %d
`, event.Title, event.StartTime.Format(common.PgTimestampFmt), event.Duration, event.Description, event.Owner, event.NotifyTime, time.Now().Format(common.PgTimestampFmt), id)
	res, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return common.ErrNoSuchEvent
	}
	s.log.Trace("modified event ", id)
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM events WHERE id = %d`, id)
	res, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return common.ErrNoSuchEvent
	}
	s.log.Trace("removed event ", id)
	return nil
}

func (s *Storage) ListEventsByDay(ctx context.Context, date time.Time) ([]common.Event, error) {
	return s.listEvents(ctx, date, date.AddDate(0, 0, 1))
}

func (s *Storage) ListEventsByWeek(ctx context.Context, date time.Time) ([]common.Event, error) {
	return s.listEvents(ctx, date, date.AddDate(0, 0, 7))
}

func (s *Storage) ListEventsByMonth(ctx context.Context, date time.Time) ([]common.Event, error) {
	return s.listEvents(ctx, date, date.AddDate(0, 1, 0))
}

func (s *Storage) listEvents(ctx context.Context, fromDate, toDate time.Time) ([]common.Event, error) {
	query := fmt.Sprintf(`
SELECT *
FROM events
WHERE start_time >= timestamp '%s'
  AND start_time < timestamp '%s'
`, fromDate.Format(common.PgTimestampFmt), toDate.Format(common.PgTimestampFmt))
	result := make([]common.Event, 0)
	rows, err := s.db.QueryxContext(ctx, query)
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warn("err closing rows: ", err)
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

func (s *Storage) ListEventsToNotify(ctx context.Context) (events []common.Event, err error) {
	query := `
SELECT *
FROM events
WHERE EXTRACT(EPOCH FROM start_time) - EXTRACT(EPOCH FROM NOW()) < notify_time;
`
	rows, err := s.db.QueryxContext(ctx, query)
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warn("err closing rows: ", err)
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
		events = append(events, event)
	}
	return events, nil
}
