package masstasker

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	masstasker "mkm.pub/masstasker/pkg/proto"
)

type Task = masstasker.Task

// A Client is a high-level client to the MassTasker API.
//
// It operates directly on masstasker.Task protobuf messages. The main advantage of using this high-level API
// over the low-level gRPC one is that all Client's method operate on Task structures and you rarely have to
// care about IDs.
//
// Tasks in MassTasker server are immutable. You can't update a task. What you do is you delete a task and atomically
// create a new one. This high-level API lets you express this by simply mutating the in-memory client task and issuing an update
// request. After the call, the task will have the new ID.
type Client struct {
	conn *grpc.ClientConn
	RPC  masstasker.MassTaskerClient
}

func (c *Client) Close() {
	c.conn.Close()
}

// Dial connects to the gRPC service and returns a Client.
//
// Common opts: grpc.WithTransportCredentials(insecure.NewCredentials())
func Dial(ctx context.Context, addr string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}
	return Connect(conn), nil
}

// Connect creates a connection from an underlying grpc client connection.
func Connect(conn *grpc.ClientConn) *Client {
	return &Client{
		conn: conn,
		RPC:  masstasker.NewMassTaskerClient(conn),
	}
}

// Create creates the given tasks.
// The IDs are allocated upon task creation and are not known until the server responds.
// This method updates the task IDs in the arguments with the IDs returned by the server.
func (c *Client) Create(ctx context.Context, tasks ...*masstasker.Task) error {
	return c.ComplexUpdate(ctx, tasks, nil, nil)
}

func (c *Client) Update(ctx context.Context, tasks ...*masstasker.Task) error {
	return c.ComplexUpdate(ctx, tasks, tasks, nil)
}

func (c *Client) Delete(ctx context.Context, tasks ...*masstasker.Task) error {
	return c.ComplexUpdate(ctx, nil, tasks, nil)
}

func (c *Client) ComplexUpdate(ctx context.Context, create []*masstasker.Task, delete []*masstasker.Task, predicates []*masstasker.Task) error {
	res, err := c.RPC.Update(ctx, &masstasker.UpdateRequest{
		Created:    create,
		Deleted:    taskRefs(delete),
		Predicates: taskIDs(predicates),
	})
	if err != nil {
		return err
	}
	for i, id := range res.GetCreatedIds() {
		create[i].Id = id
	}
	return nil
}

type QueryOpt func(*masstasker.QueryRequest)

func NonBlocking(req *masstasker.QueryRequest) {
	req.Wait = false
}

func (c *Client) Query(ctx context.Context, group string, ownFor time.Duration, opts ...QueryOpt) (*masstasker.Task, error) {
	req := &masstasker.QueryRequest{
		Group:  group,
		Now:    timestamppb.Now(),
		OwnFor: durationpb.New(ownFor),
		Wait:   true,
	}
	for _, o := range opts {
		o(req)
	}
	res, err := c.RPC.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Task, nil
}

func taskRefs(tasks []*masstasker.Task) []*masstasker.TaskRef {
	var refs []*masstasker.TaskRef
	for _, t := range tasks {
		refs = append(refs, &masstasker.TaskRef{
			Sel: &masstasker.TaskRef_Id{Id: t.Id},
		})
	}
	return refs
}

func taskIDs(tasks []*masstasker.Task) []uint64 {
	var ids []uint64
	for _, t := range tasks {
		ids = append(ids, t.Id)
	}
	return ids
}
