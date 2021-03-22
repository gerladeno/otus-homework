package sqlstorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	// run ONLY on empty DB
	t.Skip()
	log := logrus.New()
	events, err := New(log, "host=localhost port=5432 user=calendar password=calendar dbname=postgres sslmode=disable")
	require.NoError(t, err)
	tt, err := time.Parse(common.PgTimestampFmt, "2020-01-01 00:00:00")
	require.NoError(t, err)
	id, err := events.AddEvent(context.Background(), common.Event{
		Title:      "First",
		StartTime:  tt,
		Duration:   0,
		InviteList: "blabla",
		Comment:    "first",
	})
	require.NoError(t, err)
	require.Equal(t, id, uint64(0))

	id, err = events.AddEvent(context.Background(), common.Event{
		Title:      "Second",
		StartTime:  tt,
		Duration:   0,
		InviteList: "blablablabla",
		Comment:    "Second",
	})
	require.NoError(t, err)
	require.Equal(t, id, uint64(1))

	test, _ := events.ListEvents()
	require.Len(t, test, 2)

	err = events.EditEvent(context.Background(), 0, common.Event{
		Title:      "First edited",
		StartTime:  tt,
		Duration:   0,
		InviteList: "blabla edited",
		Comment:    "First edited",
	})
	require.NoError(t, err)

	err = events.RemoveEvent(context.Background(), 1)
	require.NoError(t, err)

	require.Len(t, events.events, 1)
	elem, err := events.GetEvent(1)
	require.True(t, errors.Is(err, common.ErrNoSuchEvent))
	elem, err = events.GetEvent(0)
	require.NoError(t, err)
	require.Equal(t, elem.Title, "First edited")
	require.Equal(t, elem.StartTime, tt)
	require.Equal(t, elem.InviteList, "blabla edited")
	require.Equal(t, elem.Comment, "First edited")
	require.True(t, elem.Created.Before(elem.Updated))

	id, err = events.AddEvent(context.Background(), common.Event{})
	require.NoError(t, err)
	require.Equal(t, id, uint64(2))
}
