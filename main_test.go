package main

import (
	"context"
	"log"
	"net"
	"testing"

	taskmaster "github.com/mkmik/taskmaster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupTestGRPCServer(t *testing.T) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":0") // listen on a random free port
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	taskmaster.RegisterTaskmasterServer(grpcServer, newServer())

	// Register cleanup to stop the server and close the listener when the test finishes
	t.Cleanup(func() {
		grpcServer.Stop()
		lis.Close()
	})

	// Start serving gRPC server in a separate goroutine
	go func() {
		//		t.Logf("Listening on %v", lis)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return grpcServer, lis
}

func createTestClient(t *testing.T) taskmaster.TaskmasterClient {
	_, lis := setupTestGRPCServer(t)
	clientConn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}

	t.Cleanup(func() { clientConn.Close() })
	client := taskmaster.NewTaskmasterClient(clientConn)

	return client
}

func TestConnect(t *testing.T) {
	client := createTestClient(t)

	_, err := client.Debug(context.Background(), &taskmaster.DebugRequest{})
	if err != nil {
		t.Errorf("Failed to call Debug method: %v", err)
	}
}
