package common

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const PgTimestampFmt = `2006-01-02 15:04:05`

var ErrNoSuchEvent = errors.New("no such event")

type Event struct {
	ID         uint64        `json:"id" db:"id"`
	Title      string        `json:"title" db:"title"`
	StartTime  time.Time     `json:"startTime" db:"start_time"`
	Duration   time.Duration `json:"duration" db:"duration"`
	InviteList string        `json:"inviteList" db:"invite_list"`
	Comment    string        `json:"comment" db:"comment"`
	Created    time.Time     `json:"created" db:"created"`
	Updated    time.Time     `json:"updated" db:"updated"`
}

func (e *Event) ParseEvent(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return err
	}
	return nil
}

type Application interface {
	CreateEvent(ctx context.Context, event *Event) (id uint64, err error)
	ReadEvent(ctx context.Context, id uint64) (event *Event, err error)
	UpdateEvent(ctx context.Context, event *Event, id uint64) (err error)
	DeleteEvent(ctx context.Context, id uint64) (err error)
	ListEvents(ctx context.Context) (events []*Event, err error)
}

type TestApp struct{}

func (t TestApp) ReadEvent(_ context.Context, id uint64) (event *Event, err error) {
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

func (t TestApp) CreateEvent(_ context.Context, event *Event) (uint64, error) {
	if event.ID == 0 {
		return 1, nil
	}
	return 0, io.ErrShortBuffer
}

func (t TestApp) UpdateEvent(_ context.Context, _ *Event, id uint64) (err error) {
	switch id {
	case 0:
		err = ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestApp) DeleteEvent(_ context.Context, id uint64) (err error) {
	switch id {
	case 0:
		err = ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestApp) ListEvents(_ context.Context) ([]*Event, error) {
	result := make([]*Event, 5)
	st, _ := time.Parse(PgTimestampFmt, PgTimestampFmt)
	for i := 0; i < 5; i++ {
		result[i] = &Event{
			ID:         uint64(i),
			Title:      "goga",
			StartTime:  st,
			Duration:   time.Hour,
			InviteList: "sosulki",
			Comment:    "gvozd'",
			Created:    st,
			Updated:    st,
		}
	}
	return result, nil
}
