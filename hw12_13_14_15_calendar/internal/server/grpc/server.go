package grpc

import (
	"context"
	"net"
	"strconv"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/grpc/events_v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate protoc -I=proto/ proto/events_v1.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.

type RPCServer struct {
	log     *logrus.Logger
	app     common.Application
	network string
	port    int
	server  *grpc.Server
}

func NewRPCServer(app common.Application, log *logrus.Logger, network string, port int) *RPCServer {
	return &RPCServer{
		log:     log,
		network: network,
		app:     app,
		port:    port,
		server:  grpc.NewServer(),
	}
}

func (r *RPCServer) Start(ctx context.Context) error {
	l, err := net.Listen(r.network, ":"+strconv.Itoa(r.port))
	if err != nil {
		return err
	}
	reflection.Register(r.server)
	events_v1.RegisterEventsHandlerServer(r.server, r)
	go func() {
		<-ctx.Done()
		r.Stop()
	}()
	if err = r.server.Serve(l); err != nil {
		return err
	}
	return nil
}

func (r *RPCServer) Stop() {
	r.server.Stop()
}

func (r *RPCServer) UpdateEvent(ctx context.Context, request *events_v1.UpdateEventRequest) (*events_v1.UpdateEventResponse, error) {
	err := r.app.UpdateEvent(ctx, request.GetId(), pb2Event(request.GetEvent()))
	return &events_v1.UpdateEventResponse{Ok: err != nil}, err
}

func (r *RPCServer) DeleteEvent(ctx context.Context, id *events_v1.DeleteEventRequest) (*events_v1.DeleteEventResponse, error) {
	err := r.app.DeleteEvent(ctx, id.GetId())
	return &events_v1.DeleteEventResponse{Ok: err != nil}, err
}

func (r *RPCServer) CreateEvent(ctx context.Context, event *events_v1.CreateEventRequest) (*events_v1.CreateEventResponse, error) {
	id, err := r.app.CreateEvent(ctx, pb2Event(event.GetEvent()))
	if err != nil {
		return nil, err
	}
	return &events_v1.CreateEventResponse{Id: id}, nil
}

func (r *RPCServer) ListEventsByDay(ctx context.Context, date *events_v1.ListEventsRequest) (*events_v1.ListEventsResponse, error) {
	eventsList, err := r.app.ListEventsByDay(ctx, date.GetFromDate().AsTime())
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*events_v1.Event, 0, len(eventsList))
	for _, event := range eventsList {
		eventsProto = append(eventsProto, event2Pb(event))
	}
	return &events_v1.ListEventsResponse{Events: eventsProto}, nil
}

func (r *RPCServer) ListEventsByWeek(ctx context.Context, date *events_v1.ListEventsRequest) (*events_v1.ListEventsResponse, error) {
	eventsList, err := r.app.ListEventsByWeek(ctx, date.GetFromDate().AsTime())
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*events_v1.Event, 0, len(eventsList))
	for _, event := range eventsList {
		eventsProto = append(eventsProto, event2Pb(event))
	}
	return &events_v1.ListEventsResponse{Events: eventsProto}, nil
}

func (r *RPCServer) ListEventsByMonth(ctx context.Context, date *events_v1.ListEventsRequest) (*events_v1.ListEventsResponse, error) {
	eventsList, err := r.app.ListEventsByMonth(ctx, date.GetFromDate().AsTime())
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*events_v1.Event, 0, len(eventsList))
	for _, event := range eventsList {
		eventsProto = append(eventsProto, event2Pb(event))
	}
	return &events_v1.ListEventsResponse{Events: eventsProto}, nil
}

func event2Pb(source common.Event) *events_v1.Event {
	return &events_v1.Event{
		Id:          source.ID,
		Title:       source.Title,
		StartTime:   timestamppb.New(source.StartTime),
		Duration:    source.Duration,
		Description: source.Description,
		Owner:       source.Owner,
		NotifyTime:  source.NotifyTime,
		Created:     timestamppb.New(source.Created),
		Updated:     timestamppb.New(source.Updated),
	}
}

func pb2Event(source *events_v1.Event) *common.Event {
	return &common.Event{
		ID:          source.GetId(),
		Title:       source.GetTitle(),
		StartTime:   source.GetStartTime().AsTime(),
		Duration:    source.GetDuration(),
		Description: source.GetDescription(),
		Owner:       source.GetOwner(),
		NotifyTime:  source.GetNotifyTime(),
	}
}
