package server

import (
	"context"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	api "mkm.pub/masstasker"
	masstaskerpb "mkm.pub/masstasker/pkg/proto"
)

type ListenOpt func(*listenOpts)

type listenOpts struct {
	bootstrapTaskGroup string
	grpcOptions        []grpc.ServerOption
}

func WithBootstrapTaskGroup(group string) ListenOpt {
	return func(o *listenOpts) {
		o.bootstrapTaskGroup = group
	}
}

func WithMaxMsgSize(size int) ListenOpt {
	return func(o *listenOpts) {
		o.grpcOptions = append(o.grpcOptions,
			grpc.MaxRecvMsgSize(size),
			grpc.MaxSendMsgSize(size),
		)
	}
}

// Create a new MassTasker GRPC server and not start it yet
func Listen(ctx context.Context, addr string, opts ...ListenOpt) (*grpc.Server, net.Listener, error) {
	l, err := net.Listen("tcp", addr) // listen on a random free port
	if err != nil {
		return nil, nil, err
	}

	var options listenOpts
	for _, o := range opts {
		o(&options)
	}

	s := grpc.NewServer(options.grpcOptions...)
	masstaskerpb.RegisterMassTaskerServer(s, NewWithContext(ctx))
	reflection.Register(s)

	if group := options.bootstrapTaskGroup; group != "" {
		go insertBootstrapTask(l.Addr().String(), group)
	}

	return s, l, nil
}

// Create a new MassTasker GRPC server
func ListenAndServe(ctx context.Context, addr string, opts ...ListenOpt) error {
	s, l, err := Listen(ctx, addr, opts...)
	if err != nil {
		return err
	}
	defer s.Stop()
	defer l.Close()
	return s.Serve(l)
}

func insertBootstrapTask(addr string, group string) {
	ctx := context.Background()
	mt, err := api.Dial(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	if err := mt.Create(ctx, &masstaskerpb.Task{Group: group}); err != nil {
		log.Fatal(err)
	}
	slog.Info("created bootstrap task", "group", group, "addr", addr)
}
