package masstasker_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mkm.pub/masstasker"
	masstaskerpb "mkm.pub/masstasker/pkg/proto"
)

func TestDo(t *testing.T) {
	const (
		testGroup  = "test_group"
		otherGroup = "other_group"
	)

	ctx := context.Background()
	mt := testClient(t)

	taskA, err := masstasker.NewTask(testGroup, &masstaskerpb.Test{Foo: "A"})
	if err != nil {
		t.Fatal(err)
	}

	taskB, err := masstasker.NewTask(testGroup, &masstaskerpb.Test{Foo: "A"})
	if err != nil {
		t.Fatal(err)
	}

	// Test simple create
	if err := mt.Do(ctx, masstasker.Create(taskA, taskB)); err != nil {
		t.Fatal(err)
	}

	taskC, err := masstasker.NewTask(testGroup, &masstaskerpb.Test{Foo: "C"})
	if err != nil {
		t.Fatal(err)
	}

	// Test composite delete + create
	if err := mt.Do(ctx, masstasker.Delete(taskA, taskB), masstasker.Create(taskC)); err != nil {
		t.Fatal(err)
	}

	res, err := mt.RPC.Debug(ctx, &masstaskerpb.DebugRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := res.NumTasks, uint64(1); got != want {
		t.Fatalf("got: %d, want: %d", got, want)
	}
	var data masstaskerpb.Test
	if err := res.Tasks[0].UnmarshalDataTo(&data); err != nil {
		t.Fatal(err)
	}
	if got, want := data.Foo, "C"; got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}

	// Test update (sugar for delete+create)
	taskC.Data = taskA.Data
	if got, want := status.Code(mt.Do(ctx, masstasker.Update(taskC), masstasker.Pred(taskA))), codes.FailedPrecondition; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	if err := mt.Do(ctx, masstasker.Update(taskC)); err != nil {
		t.Fatal(err)
	}

	res, err = mt.RPC.Debug(ctx, &masstaskerpb.DebugRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := res.NumTasks, uint64(1); got != want {
		t.Fatalf("got: %d, want: %d", got, want)
	}
	if err := res.Tasks[0].UnmarshalDataTo(&data); err != nil {
		t.Fatal(err)
	}
	if got, want := data.Foo, "A"; got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}

	// Test move
	if got, want := taskC.Group, testGroup; got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
	if err := mt.Do(ctx, masstasker.Move(otherGroup, taskC)); err != nil {
		t.Fatal(err)
	}
	if got, want := taskC.Group, otherGroup; got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
	res, err = mt.RPC.Debug(ctx, &masstaskerpb.DebugRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := res.NumTasks, uint64(1); got != want {
		t.Fatalf("got: %d, want: %d", got, want)
	}
	if got, want := res.Tasks[0].Group, otherGroup; got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestDoCreate(t *testing.T) {
	mt := testClient(t)
	t1 := &masstasker.Task{}
	t2 := &masstasker.Task{}
	if err := mt.Do(context.Background(), masstasker.Create(t1, t2)); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(1); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := t2.GetId(), uint64(2); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}

func TestDoDelete(t *testing.T) {
	ctx := context.Background()
	mt := testClient(t)
	t1 := &masstasker.Task{}
	if err := mt.Do(ctx, masstasker.Create(t1)); err != nil {
		t.Fatal(err)
	}
	if err := mt.Do(ctx, masstasker.Delete(t1)); err != nil {
		t.Fatal(err)
	}
	err := mt.Do(ctx, masstasker.Delete(t1))
	if got, want := status.Code(err), codes.NotFound; got != want {
		t.Fatalf("Delete returned unexpected error code: got %v, want %v", got, want)
	}
}

func TestDoUpdate(t *testing.T) {
	ctx := context.Background()
	mt := testClient(t)
	t1 := &masstasker.Task{Group: "g1"}
	if err := mt.Do(ctx, masstasker.Create(t1)); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(1); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	// "move" the task in another group
	t1.Group = "g2"
	if err := mt.Do(ctx, masstasker.Update(t1)); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(2); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}

func TestDoQuery(t *testing.T) {
	const (
		testGroup1 = "testGroup1"
		testGroup2 = "testGroup2"

		foo1 = "foo1"
		foo2 = "foo2"
	)
	ctx := context.Background()
	mt := testClient(t)
	t1 := must(masstasker.NewTask(testGroup1, &Test{Foo: foo1}))(t)
	t2 := must(masstasker.NewTask(testGroup2, &Test{Foo: foo2}))(t)
	if err := mt.Do(ctx, masstasker.Create(t1, t2)); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(1); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := t2.GetId(), uint64(2); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	oldT1ID := t1.Id
	work, err := mt.Query(ctx, testGroup1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := work.Id, oldT1ID; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	work, err = mt.Query(ctx, testGroup1, 1*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	// owning means changing the task, changing the task means replacing the task.
	if got, dontWant := work.Id, oldT1ID; got == dontWant {
		t.Errorf("got: %d, dontWant: %d", got, dontWant)
	}
}

func TestDoReown(t *testing.T) {
	const (
		testGroup = "test_group"
		ownFor    = 1 * time.Minute
	)

	clock := clockwork.NewFakeClock()
	ctx := clockwork.AddToContext(context.Background(), clock)

	mt := testClientWithContext(t, ctx)
	task, err := masstasker.NewTask(testGroup, &masstaskerpb.Test{Foo: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if err := mt.Create(ctx, task); err != nil {
		t.Fatal(err)
	}

	if err := mt.Do(ctx, masstasker.Reown(ownFor, task)); err != nil {
		t.Fatal(err)
	}
	t1 := task.NotBefore.AsTime()
	t.Logf("%s", t1)

	clock.Advance(time.Hour)

	if err := mt.Do(ctx, masstasker.Reown(ownFor, task)); err != nil {
		t.Fatal(err)
	}
	t2 := task.NotBefore.AsTime()
	t.Logf("%s", t2)

	if got, want := t2.Sub(t1), time.Hour; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
