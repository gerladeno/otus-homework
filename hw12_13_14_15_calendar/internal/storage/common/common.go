package common

import (
	"context"
	"errors"
	"time"
)

const PgTimeStampFmt = `2006-01-02 15:04:05`

var ErrNoSuchEvent = errors.New("no such event")

type Storage interface {
	GetEvent(id uint64) (*Event, error)
	AddEvent(ctx context.Context, event Event) (uint64, error)
	EditEvent(ctx context.Context, id uint64, event Event) error
	RemoveEvent(ctx context.Context, id uint64) error
	ListEvents() ([]Event, error)
}

type Event struct {
	ID         uint64        `db:"id"`
	Title      string        `db:"title"`
	StartTime  time.Time     `db:"start_time"`
	Duration   time.Duration `db:"duration"`
	InviteList string        `db:"invite_list"`
	Comment    string        `db:"comment"`
	Created    time.Time     `db:"created"`
	Updated    time.Time     `db:"updated"`
}
