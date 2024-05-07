package server_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"mkm.pub/masstasker/pkg/masstaskertest"
	masstasker "mkm.pub/masstasker/pkg/proto"
)

func TestConnect(t *testing.T) {
	client := masstaskertest.New(t)

	_, err := client.Debug(context.Background(), &masstasker.DebugRequest{})
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

	client := masstaskertest.New(t)
	ctx := context.Background()

	_, err := client.Update(ctx, &masstasker.UpdateRequest{
		Created: []*masstasker.Task{
			{
				Group: testGroup,
				Data: must(anypb.New(&masstasker.Test{
					Foo: "demo",
				}))(t),
			},
		},
	})
	if err != nil {
		t.Errorf("Failed to call Debug method: %v", err)
	}

	dbg, err := client.Debug(ctx, &masstasker.DebugRequest{})
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

	res, err := client.Query(ctx, &masstasker.QueryRequest{
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

	_, err = client.Query(ctx, &masstasker.QueryRequest{
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

	_, err = client.Update(ctx, &masstasker.UpdateRequest{
		Deleted: []*masstasker.TaskRef{{
			Sel: &masstasker.TaskRef_Id{Id: res.Task.Id},
		}},
	})
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// now we should have exactly zero tasks
	dbg, err = client.Debug(ctx, &masstasker.DebugRequest{})
	if err != nil {
		t.Fatalf("Failed to call Debug method: %v", err)
	}
	if got, want := len(dbg.Tasks), 0; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	// and if we use that task ID as precondition it will fail
	_, err = client.Update(ctx, &masstasker.UpdateRequest{
		Created: []*masstasker.Task{
			{
				Group: testGroup,
				Data: must(anypb.New(&masstasker.Test{
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
