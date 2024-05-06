package main

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"mkm.pub/masstasker"
	masstaskerpb "mkm.pub/masstasker/pkg/proto"
)

func (c *Context) processMainGroup(ctx context.Context) error {
	slog.Info("querying main task", "group", c.MainGroup)
	task, err := c.mt.Query(ctx, c.MainGroup, c.OwnFor)
	if err != nil {
		return err
	}

	err = c.mt.RunWithLease(ctx, task, c.processMainTask, masstasker.WithOwnFor(c.OwnFor), masstasker.WithHeartbeat(2*time.Second))
	return c.mt.CommitOrMove(ctx, task, err, c.ErrorGroup)
}

func (c *Context) processMainTask(ctx context.Context, task *masstasker.Task) error {
	var data masstaskerpb.Test
	if err := task.UnmarshalDataTo(&data); err != nil {
		return err
	}
	slog.Info("got task", "task", protojson.Format(&data))

	slog.Info("doing some fake work", "duration", c.FakeWorkDuration)
	time.Sleep(c.FakeWorkDuration)
	slog.Info("work completed, committing")

	return c.mt.CommitOrMove(ctx, task, nil, "")
}
