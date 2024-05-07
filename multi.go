package masstasker

import (
	"context"
	"time"

	"github.com/jonboulle/clockwork"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MultiOp func(ctx context.Context, req *multiRequest)
type multiRequest struct {
	create []*Task
	delete []*Task
	pred   []*Task
}

// The Create op will create new tasks when executed. The allocated IDs will be set in the provided
// task when the op is executed.
func Create(tasks ...*Task) MultiOp {
	return func(ctx context.Context, req *multiRequest) {
		req.create = append(req.create, tasks...)
	}
}

// The delete op will delete tasks when executed.
func Delete(tasks ...*Task) MultiOp {
	return func(ctx context.Context, req *multiRequest) {
		req.delete = append(req.delete, tasks...)
	}
}

// The pred op records the IDs of the provided tasks and check that tasks with such IDs exist.
func Pred(tasks ...*Task) MultiOp {
	return func(ctx context.Context, req *multiRequest) {
		req.pred = append(req.pred, tasks...)
	}
}

// Update is an operation that combines a Delete with a Create of the same task.
//
// Tasks are immutable and cannot be updated in place. By deleting the task and contextually
// re-create a new task with the updated content the ID changes and other workers can notice
// that the task has been updated.
func Update(tasks ...*Task) MultiOp {
	return func(ctx context.Context, req *multiRequest) {
		req.create = append(req.create, tasks...)
		req.delete = append(req.delete, tasks...)
	}
}

// The Move operation is an Update where all the Group fields are updated to the provided value.
func Move(group string, tasks ...*Task) MultiOp {
	return func(ctx context.Context, req *multiRequest) {
		for _, t := range tasks {
			t.Group = group
		}
		Update(tasks...)(ctx, req)
	}
}

// The Disown operation resets NotBefore before updating.
func Disown(tasks ...*Task) MultiOp {
	return func(ctx context.Context, req *multiRequest) {
		for _, t := range tasks {
			t.NotBefore = nil
		}
		Update(tasks...)(ctx, req)
	}
}

// Reown is just a "sugar" for a) set the NotBefore field in every task by now+ownFor b) issue an Update
// Sometimes it makes the intent of the client code clearer.
func Reown(ownFor time.Duration, tasks ...*Task) MultiOp {
	return func(ctx context.Context, req *multiRequest) {
		clock := clockwork.FromContext(ctx)
		for _, t := range tasks {
			t.NotBefore = timestamppb.New(clock.Now().Add(ownFor))
		}
		Update(tasks...)(ctx, req)
	}
}

// DeleteOrMove produces either a Move op or a Delete op, depending on whether
// the error is a saved error
func DeleteOrMove(errorGroup string, tasks ...*Task) func(err error) MultiOp {
	return func(err error) MultiOp {
		if IsSavedError(err) {
			// some error messages may contain bad utf-8 data. If we want to save the error messages we must clean it up
			filtered, _, err := transform.String(runes.ReplaceIllFormed(), err.Error())
			for _, t := range tasks {
				if err != nil {
					t.Error = err.Error()
				} else {
					t.Error = filtered
				}
			}
			return Move(errorGroup, tasks...)
		} else {
			return Delete(tasks...)
		}
	}
}

// Do performs a transactional update of the masstasker state by performing all the provided operations atomically.
func (mt *Client) Do(ctx context.Context, ops ...MultiOp) error {
	var r multiRequest
	for _, op := range ops {
		op(ctx, &r)
	}
	return mt.ComplexUpdate(ctx, r.create, r.delete, r.pred)
}
