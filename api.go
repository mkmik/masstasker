package masstasker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	masstasker "mkm.pub/masstasker/pkg/proto"
)

const (
	emptyGroupRetryTime = 5 * time.Second
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
	clock := clockwork.FromContext(ctx)

	for _, c := range create {
		if c.NotBefore == nil {
			c.NotBefore = timestamppb.New(clock.Now())
		}
	}

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

// Move is just a "sugar" for a) update the Group field in every task b) issue an Update to to move the task to another group.
// Sometimes it makes the intent of the client code clearer.
func (c *Client) Move(ctx context.Context, targetGroup string, tasks ...*masstasker.Task) error {
	for _, t := range tasks {
		t.Group = targetGroup
	}
	return c.Update(ctx, tasks...)
}

// Disown is just a "sugar" for a) reset the NotBefore field in every task b) issue an Update
// Sometimes it makes the intent of the client code clearer.
func (c *Client) Disown(ctx context.Context, tasks ...*masstasker.Task) error {
	for _, t := range tasks {
		t.NotBefore = nil
	}
	return c.Update(ctx, tasks...)
}

// Reown is just a "sugar" for a) set the NotBefore field in every task by now+ownFor b) issue an Update
// Sometimes it makes the intent of the client code clearer.
func (c *Client) Reown(ctx context.Context, ownFor time.Duration, tasks ...*masstasker.Task) error {
	clock := clockwork.FromContext(ctx)
	for _, t := range tasks {
		t.NotBefore = timestamppb.New(clock.Now().Add(ownFor))
	}
	return c.Update(ctx, tasks...)
}

type QueryOpt func(*masstasker.QueryRequest)

func NonBlocking(v bool) QueryOpt {
	return func(req *masstasker.QueryRequest) {
		req.Wait = !v
	}
}

func (c *Client) Query(ctx context.Context, group string, ownFor time.Duration, opts ...QueryOpt) (*masstasker.Task, error) {
	clock := clockwork.FromContext(ctx)

	req := &masstasker.QueryRequest{
		Group:  group,
		Now:    timestamppb.New(clock.Now()),
		OwnFor: durationpb.New(ownFor),
		Wait:   true,
	}
	for _, o := range opts {
		o(req)
	}

	for {
		res, err := c.RPC.Query(ctx, req)
		if code := status.Code(err); req.Wait && (code == codes.NotFound || code == codes.DeadlineExceeded) {
			clock.Sleep(emptyGroupRetryTime)
			req.Now = timestamppb.New(clock.Now())

			continue
		}
		if err != nil {
			return nil, err
		}
		return res.Task, nil
	}
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

// NewTask is sugar for creating a Task in a given group and the given proto into the Data field.
func NewTask(group string, data proto.Message) (*masstasker.Task, error) {
	task := &masstasker.Task{Group: group}
	if data != nil {
		if err := task.MarshalFrom(data); err != nil {
			return nil, err
		}
	}
	return task, nil
}

// NewSavedError creates a SavedError
func NewSavedError(err error) error {
	return SavedError{err: err}
}

// A SavedError wraps an error and can be used to signal that an error is not meant to be retried indefinitely
// but instead the task should be moved to an error group.
type SavedError struct {
	err error
}

func (r SavedError) Error() string {
	return r.Unwrap().Error()
}

func (r SavedError) Unwrap() error {
	return r.err
}

func IsSavedError(err error) bool {
	var serr SavedError
	return errors.As(err, &serr)
}

type leaseOptions struct {
	heartbeat time.Duration
	ownFor    time.Duration
}

type LeaseOption func(*leaseOptions)

func WithHeartbeat(d time.Duration) LeaseOption {
	return func(o *leaseOptions) {
		o.heartbeat = d
	}
}

func WithOwnFor(d time.Duration) LeaseOption {
	return func(o *leaseOptions) {
		o.ownFor = d
	}
}

// RunWithLease runs fn while in the background "renewing a lease" on the task by periodically updating bumping the NotBefore
func (mt *Client) RunWithLease(ctx context.Context, task *masstasker.Task, fn func(context.Context, *masstasker.Task) error, opts ...LeaseOption) error {
	clock := clockwork.FromContext(ctx)

	opt := leaseOptions{
		heartbeat: 10 * time.Second,
		ownFor:    30 * time.Second,
	}
	for _, o := range opts {
		o(&opt)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	nt := task.Clone()
	errCh := make(chan error, 1)
	go func() { errCh <- fn(ctx, nt) }()

	for {
		select {
		case err := <-errCh:
			return err
		case <-clock.After(opt.heartbeat):
			if err := mt.Reown(ctx, opt.ownFor, task); err != nil {
				return fmt.Errorf("error renewing task lease: %w", err)
			}
		}
	}
}

// CommitOrMove deletes the task if err is nil, or moves the task to an errorGroup if the error is a "saved error".
// Otherwise it "disowns" the task so that another worker can pick it up as soon as they are not busy with older tasks.
// Any other error that is not a "saved error" is returned by this function after the task is "disowned"
// If errorGroup is empty, then it doesn't move saved errors to the error group but simply returns them as any other error.
func (mt *Client) CommitOrMove(ctx context.Context, task *masstasker.Task, err error, errorGroup string) error {
	if err == nil {
		// mark the task as done
		return mt.Delete(ctx, task)
	}
	if errorGroup != "" {
		if IsSavedError(err) {
			// some error messages may contain bad utf-8 data. If we want to save the error messages we must clean it up
			filtered, _, err := transform.String(runes.ReplaceIllFormed(), err.Error())
			if err != nil {
				task.Error = err.Error()
			} else {
				task.Error = filtered
			}
			return mt.Move(ctx, errorGroup, task)
		}
	}
	// for other errors let's first disown the task so that other workers can pick it up
	// as soon as other tasks in front of the queue are processed,
	// and then return the error so that this worker can crash.

	// This case catches transient errors (which we want to be transparently retried)
	// and misconfigurations where we want to notice that the worker is crashing.

	_ = mt.Disown(ctx, task) // error ignored intentionally since resetting NotBefore is just an optimization

	return err
}
