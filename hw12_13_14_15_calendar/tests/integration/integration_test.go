// build +integration

package integration_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	internalgrpc "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/grpc/eventsv1"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/tests/integration"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const retries = 3

type CalendarSuite struct {
	ctx context.Context
	suite.Suite
	clientHTTP     *integration.CalendarHTTPApi
	clientGRPC     eventsv1.EventsHandlerClient
	idsToDelete    map[int64]struct{}
	notificationCh chan *common.Notification
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

func (s *CalendarSuite) SetupSuite() {
	calendarHost := os.Getenv("CALENDAR_HOST")
	if calendarHost == "" {
		calendarHost = "127.0.0.1"
	}
	var (
		err      error
		connGRPC *grpc.ClientConn
	)
	for i := 0; i < retries; i++ {
		connGRPC, err = grpc.Dial(calendarHost+":3005", grpc.WithInsecure())
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	s.Require().NoError(err)
	connHTTP := http.Client{}
	connHTTP.Transport = http.DefaultTransport
	connHTTP.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	s.clientHTTP = &integration.CalendarHTTPApi{ConnHTTP: &connHTTP, Host: "http://" + calendarHost + ":8888"}
	s.clientGRPC = eventsv1.NewEventsHandlerClient(connGRPC)
	s.notificationCh = make(chan *common.Notification, 100)
	s.idsToDelete = make(map[int64]struct{})
	s.wg.Add(1)
	go s.serveNotifyHTTP()
}

func (s *CalendarSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
}

func (s *CalendarSuite) TearDownTest() {
	for id := range s.idsToDelete {
		_, err := s.clientGRPC.DeleteEvent(s.ctx, &eventsv1.DeleteEventRequest{Id: id})
		s.Require().NoError(err)
		delete(s.idsToDelete, id)
	}
	s.cancel()
}

func (s *CalendarSuite) TearDownSuite() {
	s.wg.Wait()
}

func (s *CalendarSuite) TestAddAndListEventsGRPC() {
	t, err := time.Parse(common.PgTimestampFmt, "2001-01-01 00:00:00")
	s.Require().NoError(err)
	ids := make(map[int64]struct{})
	event := common.Event{
		Title:       "Very Important Event",
		StartTime:   t,
		Duration:    10,
		Description: "shitty description",
		Owner:       1,
		NotifyTime:  600000000,
	}
	id, err := s.clientGRPC.CreateEvent(s.ctx, &eventsv1.CreateEventRequest{Event: internalgrpc.Event2Pb(event)})
	s.Require().NoError(err)
	ids[id.GetId()] = struct{}{}
	toUpdate := id.GetId()

	event.Owner = 2
	event.StartTime = event.StartTime.Add(5 * 24 * time.Hour)
	id, err = s.clientGRPC.CreateEvent(s.ctx, &eventsv1.CreateEventRequest{Event: internalgrpc.Event2Pb(event)})
	s.Require().NoError(err)
	ids[id.GetId()] = struct{}{}

	event.Owner = 3
	event.StartTime = event.StartTime.Add(20 * 24 * time.Hour)
	id, err = s.clientGRPC.CreateEvent(s.ctx, &eventsv1.CreateEventRequest{Event: internalgrpc.Event2Pb(event)})
	s.Require().NoError(err)
	ids[id.GetId()] = struct{}{}

	event.Owner = 4
	event.StartTime = event.StartTime.Add(20 * 24 * time.Hour)
	id, err = s.clientGRPC.CreateEvent(s.ctx, &eventsv1.CreateEventRequest{Event: internalgrpc.Event2Pb(event)})
	s.Require().NoError(err)
	ids[id.GetId()] = struct{}{}

	event.Owner = 100
	event.StartTime = t
	_, err = s.clientGRPC.UpdateEvent(s.ctx, &eventsv1.UpdateEventRequest{Event: internalgrpc.Event2Pb(event), Id: toUpdate})
	s.Require().NoError(err)

	events, err := s.clientGRPC.ListEventsByDay(s.ctx, &eventsv1.ListEventsRequest{FromDate: timestamppb.New(t)})
	s.Require().NoError(err)
	s.Require().Len(events.GetEvents(), 1)
	s.Require().Equal(events.GetEvents()[0].GetOwner(), int64(100))

	events, err = s.clientGRPC.ListEventsByWeek(s.ctx, &eventsv1.ListEventsRequest{FromDate: timestamppb.New(t)})
	s.Require().NoError(err)
	s.Require().Len(events.GetEvents(), 2)

	events, err = s.clientGRPC.ListEventsByMonth(s.ctx, &eventsv1.ListEventsRequest{FromDate: timestamppb.New(t)})
	s.Require().NoError(err)
	s.Require().Len(events.GetEvents(), 3)

	for id := range ids {
		s.idsToDelete[id] = struct{}{}
	}
	err = func() error {
		timer := time.NewTimer(70 * time.Second)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-timer.C
			cancel()
		}()
		for {
			select {
			case notification := <-s.notificationCh:
				delete(ids, notification.ID)
				if len(ids) == 0 {
					return nil
				}
			case <-ctx.Done():
				return fmt.Errorf("failed to notify on about %d events", len(ids))
			}
		}
	}()
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestAddAndListEventsHTTP() {
	t, err := time.Parse(common.PgTimestampFmt, "2001-01-01 00:00:00")
	s.Require().NoError(err)
	ids := make(map[int64]struct{})
	event := common.Event{
		Title:       "Very Important Event",
		StartTime:   t,
		Duration:    10,
		Description: "shitty description",
		Owner:       1,
		NotifyTime:  600000000,
	}
	id, code := s.clientHTTP.CreateEvent(s.ctx, event)
	s.Require().Equal(code, http.StatusOK)
	ids[id] = struct{}{}
	toUpdate := id

	event.Owner = 2
	event.StartTime = event.StartTime.Add(5 * 24 * time.Hour)
	id, code = s.clientHTTP.CreateEvent(s.ctx, event)
	s.Require().Equal(code, http.StatusOK)
	ids[id] = struct{}{}

	event.Owner = 3
	event.StartTime = event.StartTime.Add(20 * 24 * time.Hour)
	id, code = s.clientHTTP.CreateEvent(s.ctx, event)
	s.Require().Equal(code, http.StatusOK)
	ids[id] = struct{}{}

	event.Owner = 4
	event.StartTime = event.StartTime.Add(20 * 24 * time.Hour)
	id, code = s.clientHTTP.CreateEvent(s.ctx, event)
	s.Require().Equal(code, http.StatusOK)
	ids[id] = struct{}{}

	event.Owner = 100
	event.StartTime = t
	code = s.clientHTTP.UpdateEvent(s.ctx, event, toUpdate)
	s.Require().Equal(code, http.StatusOK)

	events, code := s.clientHTTP.ListEventsByDay(s.ctx, t.Format("2006-01-02"))
	s.Require().Equal(code, http.StatusOK)
	s.Require().Len(events, 1)
	s.Require().Equal(events[0].Owner, int64(100))

	events, code = s.clientHTTP.ListEventsByWeek(s.ctx, t.Format("2006-01-02"))
	s.Require().Equal(code, http.StatusOK)
	s.Require().Len(events, 2)

	events, code = s.clientHTTP.ListEventsByMonth(s.ctx, t.Format("2006-01-02"))
	s.Require().Equal(code, http.StatusOK)
	s.Require().Len(events, 3)

	for id := range ids {
		s.idsToDelete[id] = struct{}{}
	}
	err = func() error {
		timer := time.NewTimer(70 * time.Second)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-timer.C
			cancel()
		}()
		for {
			select {
			case notification := <-s.notificationCh:
				delete(ids, notification.ID)
				if len(ids) == 0 {
					return nil
				}
			case <-ctx.Done():
				return fmt.Errorf("failed to notify on about %d events", len(ids))
			}
		}
	}()
	s.Require().NoError(err)
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}

func (s *CalendarSuite) serveNotifyHTTP() {
	defer s.wg.Done()
	notifyHandler := func(w http.ResponseWriter, r *http.Request) {
		n := common.Notification{}
		err := json.NewDecoder(r.Body).Decode(&n)
		s.Require().NoError(err)
		s.notificationCh <- &n
		_, err = fmt.Fprintf(w, "ok")
		s.Require().NoError(err)
	}
	http.HandleFunc("/notify", notifyHandler)
	go func() {
		err := http.ListenAndServe(":3002", nil)
		s.Require().NoError(err)
	}()
	<-s.ctx.Done()
}
