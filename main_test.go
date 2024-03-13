package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	taskmaster "github.com/mkmik/taskmaster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
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

func must[T any](v T, err error) func(tb testing.TB) T {
	return func(t testing.TB) T {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}

func TestSimpleWorkflow(t *testing.T) {
	const (
		testGroup        = "test1"
		testID    uint64 = 1
	)

	client := createTestClient(t)
	ctx := context.Background()

	_, err := client.Update(ctx, &taskmaster.UpdateRequest{
		Created: []*taskmaster.Task{
			{
				Group: testGroup,
				Data: must(anypb.New(&taskmaster.Test{
					Foo: "demo",
				}))(t),
			},
		},
	})
	if err != nil {
		t.Errorf("Failed to call Debug method: %v", err)
	}

	dbg, err := client.Debug(ctx, &taskmaster.DebugRequest{})
	if err != nil {
		t.Fatalf("Failed to call Debug method: %v", err)
	}
	if got, want := len(dbg.Tasks), 1; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	task := dbg.Tasks[0]
	if got, want := task.Id, testID; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
	if got, want := task.Group, testGroup; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	res, err := client.Query(ctx, &taskmaster.QueryRequest{
		Group:  testGroup,
		OwnFor: durationpb.New(1 * time.Hour),
		Wait:   false,
	})
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	t.Logf("query result: %v", res)

	if got, dontWant := res.Task.Id, testID; got == dontWant {
		t.Fatalf("got: %v but should be different", got)
	}

	_, err = client.Query(ctx, &taskmaster.QueryRequest{
		Group:  testGroup,
		OwnFor: durationpb.New(1 * time.Hour),
		Wait:   false,
	})
	if err == nil {
		t.Fatalf("Expecting error")
	}
	if got, want := status.Code(err), codes.NotFound; got != want {
		t.Fatalf("Query returned unexpected error code: got %v, want %v", got, want)
	}

	_, err = client.Update(ctx, &taskmaster.UpdateRequest{
		Deleted: []*taskmaster.TaskRef{{
			Sel: &taskmaster.TaskRef_Id{Id: res.Task.Id},
		}},
	})
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// now we should have exactly zero tasks
	dbg, err = client.Debug(ctx, &taskmaster.DebugRequest{})
	if err != nil {
		t.Fatalf("Failed to call Debug method: %v", err)
	}
	if got, want := len(dbg.Tasks), 0; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	// and if we use that task ID as precondition it will fail
	_, err = client.Update(ctx, &taskmaster.UpdateRequest{
		Created: []*taskmaster.Task{
			{
				Group: testGroup,
				Data: must(anypb.New(&taskmaster.Test{
					Foo: "impossible task",
				}))(t),
			},
		},
		Predicates: []uint64{res.Task.Id},
	})
	if got, want := status.Code(err), codes.FailedPrecondition; got != want {
		t.Fatalf("Query returned unexpected error code: got %v, want %v", got, want)
	}
}
