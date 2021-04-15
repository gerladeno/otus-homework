package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const PgTimestampFmt = `2006-01-02 15:04:05`

var ErrNoSuchEvent = errors.New("no such event")

type Notification struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	EventTime time.Time `json:"event_time"`
	Owner     int64     `json:"owner"`
}

func (n *Notification) Encode() ([]byte, error) {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(n)
	return b.Bytes(), err
}

func (n *Notification) String() string {
	return fmt.Sprintf("id: %d, title: %s, time: %s, onwer: %d", n.ID, n.Title, n.EventTime.Format(PgTimestampFmt), n.Owner)
}

type Event struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	StartTime   time.Time `json:"startTime" db:"start_time"`
	Duration    int64     `json:"duration" db:"duration"`
	Description string    `json:"description" db:"description"`
	Owner       int64     `json:"owner" db:"owner"`
	NotifyTime  int32     `json:"notifyTime" db:"notify_time"`
	Created     time.Time `json:"created" db:"created"`
	Updated     time.Time `json:"updated" db:"updated"`
}

func (e *Event) Notification() *Notification {
	return &Notification{
		ID:        e.ID,
		Title:     e.Title,
		EventTime: e.StartTime,
		Owner:     e.Owner,
	}
}

func (e *Event) ParseEvent(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return err
	}
	return nil
}

type Application interface {
	CreateEvent(ctx context.Context, event *Event) (id int64, err error)
	UpdateEvent(ctx context.Context, id int64, event *Event) (err error)
	DeleteEvent(ctx context.Context, id int64) (err error)
	ListEventsByDay(ctx context.Context, date time.Time) (events []Event, err error)
	ListEventsByWeek(ctx context.Context, date time.Time) (events []Event, err error)
	ListEventsByMonth(ctx context.Context, date time.Time) (events []Event, err error)
}

type TestApp struct{}

func (t TestApp) ReadEvent(_ context.Context, id int64) (event *Event, err error) {
	switch id {
	case 0:
		err = ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
		event = &Event{}
	}
	return event, err
}

func (t TestApp) CreateEvent(_ context.Context, event *Event) (int64, error) {
	if event.ID == 0 {
		return 1, nil
	}
	return 0, io.ErrShortBuffer
}

func (t TestApp) UpdateEvent(_ context.Context, id int64, _ *Event) (err error) {
	switch id {
	case 0:
		err = ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestApp) DeleteEvent(_ context.Context, id int64) (err error) {
	switch id {
	case 0:
		err = ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestApp) ListEventsByDay(_ context.Context, date time.Time) ([]Event, error) {
	return t.listEvents(date, 5)
}

func (t TestApp) ListEventsByWeek(_ context.Context, date time.Time) ([]Event, error) {
	return t.listEvents(date, 15)
}

func (t TestApp) ListEventsByMonth(_ context.Context, date time.Time) ([]Event, error) {
	return t.listEvents(date, 50)
}

func (t TestApp) listEvents(dateTime time.Time, cnt int) ([]Event, error) {
	result := make([]Event, cnt)
	for i := 0; i < cnt; i++ {
		result[i] = Event{
			ID:          int64(i),
			Title:       "goga",
			StartTime:   dateTime,
			Duration:    60 * 60,
			Description: "description",
			Owner:       int64(i * 2),
			NotifyTime:  int32(i * 10),
			Created:     dateTime,
			Updated:     dateTime,
		}
	}
	return result, nil
}
