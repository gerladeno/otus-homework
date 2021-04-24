package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("methods", func(t *testing.T) {
		log := logrus.New()
		events := New(log)
		tt, err := time.Parse(common.PgTimestampFmt, "2020-01-01 00:00:00")
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		id, err := events.CreateEvent(ctx, &common.Event{
			Title:       "First",
			StartTime:   tt,
			Duration:    0,
			Owner:       11,
			Description: "first",
			NotifyTime:  100,
		})
		require.NoError(t, err)
		require.Equal(t, id, int64(0))

		id, err = events.CreateEvent(ctx, &common.Event{
			Title:       "Second",
			StartTime:   tt,
			Duration:    0,
			Owner:       11,
			Description: "second",
			NotifyTime:  100,
		})
		require.NoError(t, err)
		require.Equal(t, id, int64(1))

		test, _ := events.ListEventsByDay(ctx, tt)
		require.Len(t, test, 2)

		err = events.UpdateEvent(ctx, 0, &common.Event{
			Title:       "First edited",
			StartTime:   tt,
			Duration:    5,
			Owner:       15,
			Description: "First edited",
			NotifyTime:  1000,
		})
		require.NoError(t, err)

		err = events.DeleteEvent(ctx, 1)
		require.NoError(t, err)

		require.Len(t, events.events, 1)
		elems, err := events.ListEventsByDay(ctx, tt)
		require.NoError(t, err)
		require.Equal(t, elems[0].Title, "First edited")
		require.Equal(t, elems[0].StartTime, tt)
		require.Equal(t, elems[0].Duration, int64(5))
		require.Equal(t, elems[0].Owner, int64(15))
		require.Equal(t, elems[0].Description, "First edited")
		require.Equal(t, elems[0].NotifyTime, int32(1000))
		require.True(t, elems[0].Created.Before(elems[0].Updated))

		id, err = events.CreateEvent(ctx, &common.Event{})
		require.NoError(t, err)
		require.Equal(t, id, int64(2))

		_, err = events.CreateEvent(ctx, &common.Event{StartTime: tt.AddDate(0, 0, 6)})
		require.NoError(t, err)
		_, err = events.CreateEvent(ctx, &common.Event{StartTime: tt.AddDate(0, 0, 20)})
		require.NoError(t, err)
		_, err = events.CreateEvent(ctx, &common.Event{StartTime: tt.AddDate(0, 0, 60)})
		require.NoError(t, err)

		test, err = events.ListEventsByDay(ctx, tt)
		require.NoError(t, err)
		require.Len(t, test, 1)
		test, err = events.ListEventsByWeek(ctx, tt)
		require.NoError(t, err)
		require.Len(t, test, 2)
		test, err = events.ListEventsByMonth(ctx, tt)
		require.NoError(t, err)
		require.Len(t, test, 3)
	})
	t.Run("concurrent", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		l := 100
		log := logrus.New()
		events := New(log)
		var wg sync.WaitGroup
		for i := 0; i < l; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := events.CreateEvent(ctx, &common.Event{})
				require.NoError(t, err)
			}()
		}
		wg.Wait()
		require.Len(t, events.events, l)
	})
}
