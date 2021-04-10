package grpc

import (
	"context"
	"io"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/grpc/eventspb"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const testPort = 3005

func TestRPCServer(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute) //nolint:govet
	r := NewRPCServer(common.TestApp{}, logrus.New(), testPort)
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
	events, err := client.ListEvents(ctx, &empty.Empty{})
	require.NoError(t, err)

	st, _ := time.Parse(common.PgTimestampFmt, common.PgTimestampFmt)
	for i, event := range events.Events {
		tmp := pb2Event(event)
		require.Equal(t, tmp, &common.Event{
			ID:         uint64(i),
			Title:      "goga",
			StartTime:  st,
			Duration:   time.Hour,
			InviteList: "sosulki",
			Comment:    "gvozd'",
		})
	}

	testEvent := common.Event{
		ID:         1,
		Title:      "goga",
		StartTime:  st,
		Duration:   time.Hour,
		InviteList: "sosulki",
		Comment:    "gvozd'",
	}

	_, err = client.CreateEvent(ctx, event2Pb(&testEvent))
	require.Error(t, err, io.ErrShortBuffer)
	testEvent.ID = 0
	id, err := client.CreateEvent(ctx, event2Pb(&testEvent))
	require.NoError(t, err)
	require.True(t, id.GetId() == 1)

	_, err = client.UpdateEvent(ctx, event2Pb(&testEvent))
	require.Error(t, err, io.ErrShortBuffer)
	testEvent.ID = 0
	_, err = client.UpdateEvent(ctx, event2Pb(&testEvent))
	require.Error(t, err, common.ErrNoSuchEvent)
	testEvent.ID = 2
	_, err = client.UpdateEvent(ctx, event2Pb(&testEvent))
	require.NoError(t, err)

	_, err = client.DeleteEvent(ctx, &eventspb.Id{Id: 0})
	require.Error(t, err, common.ErrNoSuchEvent)
	_, err = client.DeleteEvent(ctx, &eventspb.Id{Id: 1})
	require.Error(t, err, io.ErrShortBuffer)
	_, err = client.DeleteEvent(ctx, &eventspb.Id{Id: 2})
	require.NoError(t, err)

	_, err = client.ReadEvent(ctx, &eventspb.Id{Id: 0})
	require.Error(t, err, common.ErrNoSuchEvent)
	_, err = client.ReadEvent(ctx, &eventspb.Id{Id: 1})
	require.Error(t, err, io.ErrShortBuffer)
	event, err := client.ReadEvent(ctx, &eventspb.Id{Id: 2})
	require.NoError(t, err)
	require.Equal(t, pb2Event(event), &common.Event{})

	r.Stop()
	wg.Wait()
}

func StartClient() (eventspb.EventsHandlerClient, *grpc.ClientConn, error) {
	cc, err := grpc.Dial("localhost:"+strconv.Itoa(testPort), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := eventspb.NewEventsHandlerClient(cc)
	return client, cc, nil
}
