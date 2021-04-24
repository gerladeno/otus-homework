package app

import (
	"context"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/sirupsen/logrus"
)

type App struct {
	log     *logrus.Logger
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, event *common.Event) (id int64, err error)
	UpdateEvent(ctx context.Context, id int64, event *common.Event) (err error)
	DeleteEvent(ctx context.Context, id int64) (err error)
	ListEventsByDay(ctx context.Context, date time.Time) (events []common.Event, err error)
	ListEventsByWeek(ctx context.Context, date time.Time) (events []common.Event, err error)
	ListEventsByMonth(ctx context.Context, date time.Time) (events []common.Event, err error)
	ListEventsToNotify(ctx context.Context) (events []common.Event, err error)
}

func New(log *logrus.Logger, storage Storage) *App {
	return &App{log: log, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event *common.Event) (id int64, err error) {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, id int64, event *common.Event) (err error) {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int64) (err error) {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) ListEventsByDay(ctx context.Context, date time.Time) (events []common.Event, err error) {
	return a.storage.ListEventsByDay(ctx, date)
}

func (a *App) ListEventsByWeek(ctx context.Context, date time.Time) (events []common.Event, err error) {
	return a.storage.ListEventsByWeek(ctx, date)
}

func (a *App) ListEventsByMonth(ctx context.Context, date time.Time) (events []common.Event, err error) {
	return a.storage.ListEventsByMonth(ctx, date)
}
