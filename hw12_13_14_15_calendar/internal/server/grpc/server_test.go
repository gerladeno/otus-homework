package grpc

import (
	"context"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/grpc/events_v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testPort = 3005

func TestRPCServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r := NewRPCServer(common.TestApp{}, logrus.New(), "tcp", testPort)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := r.Start(ctx)
		require.NoError(t, err)
	}()

	client, cc, err := StartClient()
	defer func() {
		err := cc.Close()
		if err != nil {
			panic(err)
		}
	}()
	require.NoError(t, err)
	st, _ := time.Parse(common.PgTimestampFmt, common.PgTimestampFmt)
	events, err := client.ListEventsByDay(ctx, &events_v1.ListEventsRequest{FromDate: timestamppb.New(st)})
	require.NoError(t, err)

	for i, event := range events.Events {
		tmp := pb2Event(event)
		require.Equal(t, tmp, &common.Event{
			ID:          int64(i),
			Title:       "goga",
			StartTime:   st,
			Duration:    60 * 60,
			Description: "description",
			Owner:       int64(i * 2),
			NotifyTime:  int32(i * 10),
		})
	}

	events, err = client.ListEventsByWeek(ctx, &events_v1.ListEventsRequest{FromDate: timestamppb.New(st)})
	require.NoError(t, err)

	for i, event := range events.Events {
		tmp := pb2Event(event)
		require.Equal(t, tmp, &common.Event{
			ID:          int64(i),
			Title:       "goga",
			StartTime:   st,
			Duration:    60 * 60,
			Description: "description",
			Owner:       int64(i * 2),
			NotifyTime:  int32(i * 10),
		})
	}

	events, err = client.ListEventsByMonth(ctx, &events_v1.ListEventsRequest{FromDate: timestamppb.New(st)})
	require.NoError(t, err)

	for i, event := range events.Events {
		tmp := pb2Event(event)
		require.Equal(t, tmp, &common.Event{
			ID:          int64(i),
			Title:       "goga",
			StartTime:   st,
			Duration:    60 * 60,
			Description: "description",
			Owner:       int64(i * 2),
			NotifyTime:  int32(i * 10),
		})
	}

	testEvent := common.Event{
		ID:          1,
		Title:       "goga",
		StartTime:   st,
		Duration:    60 * 60,
		Description: "description",
		Owner:       2,
		NotifyTime:  10,
	}

	_, err = client.CreateEvent(ctx, &events_v1.CreateEventRequest{Event: event2Pb(testEvent)})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), io.ErrShortBuffer.Error()))
	testEvent.ID = 0
	id, err := client.CreateEvent(ctx, &events_v1.CreateEventRequest{Event: event2Pb(testEvent)})
	require.NoError(t, err)
	require.True(t, id.GetId() == 1)

	_, err = client.UpdateEvent(ctx, &events_v1.UpdateEventRequest{Event: event2Pb(testEvent), Id: 1})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), io.ErrShortBuffer.Error()))
	testEvent.ID = 0
	_, err = client.UpdateEvent(ctx, &events_v1.UpdateEventRequest{Event: event2Pb(testEvent), Id: 0})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), common.ErrNoSuchEvent.Error()))
	testEvent.ID = 2
	_, err = client.UpdateEvent(ctx, &events_v1.UpdateEventRequest{Event: event2Pb(testEvent), Id: 2})
	require.NoError(t, err)

	_, err = client.DeleteEvent(ctx, &events_v1.DeleteEventRequest{Id: 0})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), common.ErrNoSuchEvent.Error()))
	_, err = client.DeleteEvent(ctx, &events_v1.DeleteEventRequest{Id: 1})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), io.ErrShortBuffer.Error()))
	_, err = client.DeleteEvent(ctx, &events_v1.DeleteEventRequest{Id: 2})
	require.NoError(t, err)

	r.Stop()
	wg.Wait()
}

func StartClient() (events_v1.EventsHandlerClient, *grpc.ClientConn, error) {
	cc, err := grpc.Dial("localhost:"+strconv.Itoa(testPort), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := events_v1.NewEventsHandlerClient(cc)
	return client, cc, nil
}
