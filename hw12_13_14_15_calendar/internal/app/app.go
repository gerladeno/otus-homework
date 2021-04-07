package app

import (
	"context"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/sirupsen/logrus"
)

type App struct {
	log     *logrus.Logger
	storage Storage
}

type Storage interface {
	ReadEvent(ctx context.Context, id uint64) (*common.Event, error)
	CreateEvent(ctx context.Context, event *common.Event) (uint64, error)
	UpdateEvent(ctx context.Context, id uint64, event *common.Event) error
	DeleteEvent(ctx context.Context, id uint64) error
	ListEvents(ctx context.Context) ([]*common.Event, error)
}

func New(log *logrus.Logger, storage Storage) *App {
	return &App{log: log, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event *common.Event) (id uint64, err error) {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) ReadEvent(ctx context.Context, id uint64) (event *common.Event, err error) {
	return a.storage.ReadEvent(ctx, id)
}

func (a *App) UpdateEvent(ctx context.Context, event *common.Event, id uint64) (err error) {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id uint64) (err error) {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) ListEvents(ctx context.Context) (events []*common.Event, err error) {
	return a.storage.ListEvents(ctx)
}
