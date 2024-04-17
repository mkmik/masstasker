package masstasker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mkm.pub/masstasker"
	"mkm.pub/masstasker/pkg/masstaskertest"
	masstaskerpb "mkm.pub/masstasker/pkg/proto"
)

type Test = masstaskerpb.Test

func must[T any](v T, err error) func(tb testing.TB) T {
	return func(t testing.TB) T {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}

func testClient(t testing.TB) *masstasker.Client {
	return masstasker.Connect(masstaskertest.NewClientConn(t))
}

func TestCreate(t *testing.T) {
	client := testClient(t)
	t1 := &masstasker.Task{}
	t2 := &masstasker.Task{}
	if err := client.Create(context.Background(), t1, t2); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(1); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := t2.GetId(), uint64(2); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}

func TestDelete(t *testing.T) {
	client := testClient(t)
	t1 := &masstasker.Task{}
	if err := client.Create(context.Background(), t1); err != nil {
		t.Fatal(err)
	}
	if err := client.Delete(context.Background(), t1); err != nil {
		t.Fatal(err)
	}
	err := client.Delete(context.Background(), t1)
	if got, want := status.Code(err), codes.NotFound; got != want {
		t.Fatalf("Delete returned unexpected error code: got %v, want %v", got, want)
	}
}

func TestUpdate(t *testing.T) {
	client := testClient(t)
	t1 := &masstasker.Task{Group: "g1"}
	if err := client.Create(context.Background(), t1); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(1); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	// "move" the task in another group
	t1.Group = "g2"
	if err := client.Update(context.Background(), t1); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(2); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}

func TestQuery(t *testing.T) {
	const (
		testGroup1 = "testGroup1"
		testGroup2 = "testGroup2"

		foo1 = "foo1"
		foo2 = "foo2"
	)
	client := testClient(t)
	t1 := must(masstasker.NewTask(testGroup1, &Test{Foo: foo1}))(t)
	t2 := must(masstasker.NewTask(testGroup2, &Test{Foo: foo2}))(t)
	if err := client.Create(context.Background(), t1, t2); err != nil {
		t.Fatal(err)
	}
	if got, want := t1.GetId(), uint64(1); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := t2.GetId(), uint64(2); got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	oldT1ID := t1.Id
	work, err := client.Query(context.Background(), testGroup1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := work.Id, oldT1ID; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	work, err = client.Query(context.Background(), testGroup1, 1*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	// owning means changing the task, changing the task means replacing the task.
	if got, dontWant := work.Id, oldT1ID; got == dontWant {
		t.Errorf("got: %d, dontWant: %d", got, dontWant)
	}
}

func TestSavedError(t *testing.T) {
	testCases := []struct {
		err   error
		saved bool
	}{
		{
			fmt.Errorf("some error"),
			false,
		},
		{
			nil,
			false,
		},
		{
			masstasker.NewSavedError(fmt.Errorf("some error")),
			true,
		},
		{
			fmt.Errorf("wrapped %w", masstasker.NewSavedError(fmt.Errorf("some error"))),
			true,
		},
		{
			fmt.Errorf("wrapped %w", fmt.Errorf("again %w", masstasker.NewSavedError(fmt.Errorf("some error")))),
			true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if got, want := masstasker.IsSavedError(tc.err), tc.saved; got != want {
				t.Errorf("got: %v, want: %v", got, want)
			}
		})
	}
}
