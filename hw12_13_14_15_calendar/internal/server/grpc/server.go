package grpc

import (
	"context"
	"net"
	"strconv"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/grpc/eventspb"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate protoc -I=proto/ proto/events.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.

type RPCServer struct {
	log    *logrus.Logger
	app    common.Application
	port   int
	server *grpc.Server
}

func NewRPCServer(app common.Application, log *logrus.Logger, port int) *RPCServer {
	return &RPCServer{
		log:    log,
		app:    app,
		port:   port,
		server: grpc.NewServer(),
	}
}

func (r *RPCServer) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(r.port))
	if err != nil {
		return err
	}
	reflection.Register(r.server)
	eventspb.RegisterEventsHandlerServer(r.server, r)
	go func() {
		<-ctx.Done()
		r.Stop()
	}()
	if err = r.server.Serve(l); err != nil {
		r.log.Warnf("grpc server failed to start or stopped unexpectedly: %s", err)
	}
	return nil
}

func (r *RPCServer) Stop() {
	r.server.Stop()
}

func (r *RPCServer) UpdateEvent(ctx context.Context, event *eventspb.Event) (*empty.Empty, error) {
	err := r.app.UpdateEvent(ctx, pb2Event(event), event.GetId())
	return &empty.Empty{}, err
}

func (r *RPCServer) DeleteEvent(ctx context.Context, id *eventspb.Id) (*empty.Empty, error) {
	err := r.app.DeleteEvent(ctx, id.GetId())
	return &empty.Empty{}, err
}

func (r *RPCServer) ReadEvent(ctx context.Context, id *eventspb.Id) (*eventspb.Event, error) {
	event, err := r.app.ReadEvent(ctx, id.GetId())
	if err != nil {
		return &eventspb.Event{}, err
	}
	return event2Pb(event), nil
}

func (r *RPCServer) CreateEvent(ctx context.Context, event *eventspb.Event) (*eventspb.Id, error) {
	id, err := r.app.CreateEvent(ctx, pb2Event(event))
	if err != nil {
		return nil, err
	}
	return &eventspb.Id{Id: id}, nil
}

func (r *RPCServer) ListEvents(ctx context.Context, _ *empty.Empty) (*eventspb.Events, error) {
	events, err := r.app.ListEvents(ctx)
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*eventspb.Event, 0, len(events))
	for _, event := range events {
		eventsProto = append(eventsProto, event2Pb(event))
	}
	return &eventspb.Events{Events: eventsProto}, nil
}

func event2Pb(source *common.Event) *eventspb.Event {
	return &eventspb.Event{
		Id:         source.ID,
		Title:      source.Title,
		StartTime:  timestamppb.New(source.StartTime),
		Duration:   durationpb.New(source.Duration),
		InviteList: source.InviteList,
		Comment:    source.Comment,
		Created:    timestamppb.New(source.Created),
		Updated:    timestamppb.New(source.Updated),
	}
}

func pb2Event(source *eventspb.Event) *common.Event {
	return &common.Event{
		ID:         source.GetId(),
		Title:      source.GetTitle(),
		StartTime:  source.GetStartTime().AsTime(),
		Duration:   source.GetDuration().AsDuration(),
		InviteList: source.GetInviteList(),
		Comment:    source.GetComment(),
	}
}
