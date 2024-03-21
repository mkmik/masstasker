package masstasker_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mkm.pub/masstasker"
	"mkm.pub/masstasker/pkg/masstaskertest"
)

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
