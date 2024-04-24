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
	s, l, err := server.Listen(ctx, ":0")
	if err != nil {
		t.Fatal(err)
	}
	// Register cleanup to stop the server and close the listener when the test finishes
	t.Cleanup(func() {
		s.Stop()
		l.Close()
	})
	go func() {
		if err := s.Serve(l); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return s, l
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
