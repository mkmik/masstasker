package masstaskertest

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	masstasker "mkm.pub/masstasker/pkg/proto"
	"mkm.pub/masstasker/pkg/server"
)

func SetupTestGRPCServer(t testing.TB) (*grpc.Server, net.Listener) {
	return SetupTestGRPCServerWithContext(t, context.Background())
}

func SetupTestGRPCServerWithContext(t testing.TB, ctx context.Context) (*grpc.Server, net.Listener) {
	s, l, c, err := NewGRPCServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c)
	return s, l
}

func NewGRPCServer(ctx context.Context) (*grpc.Server, net.Listener, func(), error) {
	lis, err := net.Listen("tcp", ":0") // listen on a random free port
	if err != nil {
		return nil, nil, func() {}, err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	masstasker.RegisterMassTaskerServer(grpcServer, server.NewWithContext(ctx))

	// Register cleanup to stop the server and close the listener when the test finishes
	cleanup := func() {
		grpcServer.Stop()
		lis.Close()
	}

	// Start serving gRPC server in a separate goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return grpcServer, lis, cleanup, nil
}

func New(t testing.TB) masstasker.MassTaskerClient {
	return masstasker.NewMassTaskerClient(NewClientConn(t))
}

func NewWithContext(t testing.TB, ctx context.Context) masstasker.MassTaskerClient {
	return masstasker.NewMassTaskerClient(NewClientConnWithContext(t, ctx))
}

func NewClientConn(t testing.TB) *grpc.ClientConn {
	return NewClientConnWithContext(t, context.Background())
}

func NewClientConnWithContext(t testing.TB, ctx context.Context) *grpc.ClientConn {
	_, lis := SetupTestGRPCServerWithContext(t, ctx)
	clientConn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	t.Cleanup(func() { clientConn.Close() })
	return clientConn
}
