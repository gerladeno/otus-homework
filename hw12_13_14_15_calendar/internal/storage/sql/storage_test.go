package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	// run ONLY on empty DB
	// goose -dir internal/storage/sql/migrations postgres "user=calendar password=calendar dbname=postgres sslmode=disable" down
	// goose -dir internal/storage/sql/migrations postgres "user=calendar password=calendar dbname=postgres sslmode=disable" up
	t.Skip()
	log := logrus.New()
	events, err := New(context.Background(), log, "host=localhost port=5432 user=calendar password=calendar dbname=postgres sslmode=disable")
	require.NoError(t, err)
	tt, err := time.Parse(common.PgTimestampFmt, "2020-01-01 00:00:00")
	require.NoError(t, err)

	id, err := events.CreateEvent(context.Background(), &common.Event{
		Title:       "First",
		StartTime:   tt,
		Duration:    0,
		Owner:       11,
		Description: "first",
		NotifyTime:  100,
	})
	require.NoError(t, err)
	require.Equal(t, id, int64(1))

	id, err = events.CreateEvent(context.Background(), &common.Event{
		Title:       "Second",
		StartTime:   tt,
		Duration:    0,
		Owner:       11,
		Description: "second",
		NotifyTime:  100,
	})
	require.NoError(t, err)
	require.Equal(t, id, int64(2))

	test, _ := events.ListEventsByDay(context.Background(), tt)
	require.Len(t, test, 2)

	err = events.UpdateEvent(context.Background(), 1, &common.Event{
		Title:       "First edited",
		StartTime:   tt,
		Duration:    5,
		Owner:       15,
		Description: "First edited",
		NotifyTime:  1000,
	})
	require.NoError(t, err)

	err = events.DeleteEvent(context.Background(), 2)
	require.NoError(t, err)

	elems, err := events.ListEventsByDay(context.Background(), tt)
	require.Len(t, elems, 1)
	require.NoError(t, err)
	require.Equal(t, elems[0].Title, "First edited")
	require.Equal(t, elems[0].StartTime, tt)
	require.Equal(t, elems[0].Duration, int64(5))
	require.Equal(t, elems[0].Owner, int64(15))
	require.Equal(t, elems[0].Description, "First edited")
	require.Equal(t, elems[0].NotifyTime, int32(1000))
	require.True(t, elems[0].Created.Before(elems[0].Updated))

	id, err = events.CreateEvent(context.Background(), &common.Event{})
	require.NoError(t, err)
	require.Equal(t, id, int64(3))

	_, err = events.CreateEvent(context.Background(), &common.Event{StartTime: tt.AddDate(0, 0, 6)})
	require.NoError(t, err)
	_, err = events.CreateEvent(context.Background(), &common.Event{StartTime: tt.AddDate(0, 0, 20)})
	require.NoError(t, err)
	_, err = events.CreateEvent(context.Background(), &common.Event{StartTime: tt.AddDate(0, 0, 60)})
	require.NoError(t, err)

	test, err = events.ListEventsByDay(context.Background(), tt)
	require.NoError(t, err)
	require.Len(t, test, 1)
	test, err = events.ListEventsByWeek(context.Background(), tt)
	require.NoError(t, err)
	require.Len(t, test, 2)
	test, err = events.ListEventsByMonth(context.Background(), tt)
	require.NoError(t, err)
	require.Len(t, test, 3)
}
