package common

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const PgTimestampFmt = `2006-01-02 15:04:05`

var ErrNoSuchEvent = errors.New("no such event")

type Storage interface {
	GetEvent(id uint64) (*Event, error)
	AddEvent(ctx context.Context, event Event) (uint64, error)
	EditEvent(ctx context.Context, id uint64, event Event) error
	RemoveEvent(ctx context.Context, id uint64) error
	ListEvents() ([]Event, error)
}

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