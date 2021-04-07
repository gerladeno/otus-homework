package memorystorage

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("CRUD", func(t *testing.T) {
		log := logrus.New()
		events := New(log)
		tt, err := time.Parse(common.PgTimestampFmt, "2020-01-01 00:00:00")
		require.NoError(t, err)
		id, err := events.CreateEvent(context.Background(), &common.Event{
			Title:      "First",
			StartTime:  tt,
			Duration:   0,
			InviteList: "blabla",
			Comment:    "first",
		})
		require.NoError(t, err)
		require.Equal(t, id, uint64(0))

		id, err = events.CreateEvent(context.Background(), &common.Event{
			Title:      "Second",
			StartTime:  tt,
			Duration:   0,
			InviteList: "blablablabla",
			Comment:    "Second",
		})
		require.NoError(t, err)
		require.Equal(t, id, uint64(1))

		test, _ := events.ListEvents(context.Background())
		require.Len(t, test, 2)

		err = events.UpdateEvent(context.Background(), 0, &common.Event{
			Title:      "First edited",
			StartTime:  tt,
			Duration:   0,
			InviteList: "blabla edited",
			Comment:    "First edited",
		})
		require.NoError(t, err)

		err = events.DeleteEvent(context.Background(), 1)
		require.NoError(t, err)

		require.Len(t, events.events, 1)
		_, err = events.ReadEvent(context.Background(), 1)
		require.True(t, errors.Is(err, common.ErrNoSuchEvent))
		elem, err := events.ReadEvent(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, elem.Title, "First edited")
		require.Equal(t, elem.StartTime, tt)
		require.Equal(t, elem.InviteList, "blabla edited")
		require.Equal(t, elem.Comment, "First edited")
		require.True(t, elem.Created.Before(elem.Updated))

		id, err = events.CreateEvent(context.Background(), &common.Event{})
		require.NoError(t, err)
		require.Equal(t, id, uint64(2))
	})
	t.Run("concurrent", func(t *testing.T) {
		l := 100
		log := logrus.New()
		events := New(log)
		var wg sync.WaitGroup
		for i := 0; i < l; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := events.CreateEvent(context.Background(), &common.Event{})
				require.NoError(t, err)
			}()
		}
		wg.Wait()
		require.Len(t, events.events, l)
	})
}
