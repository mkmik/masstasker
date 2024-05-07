package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"mkm.pub/masstasker"
	masstaskerpb "mkm.pub/masstasker/pkg/proto"
)

func (c *Context) processBootstrap(ctx context.Context) error {
	bootstrap, err := c.mt.Query(ctx, c.BootstrapGroup, c.OwnFor)
	if err != nil {
		return err
	}
	slog.Info("got bootstrap task")

	time.Sleep(2 * time.Second)
	var tasks []*masstasker.Task
	for i := 0; i < c.NumTasks; i++ {
		t, err := masstasker.NewTask(c.MainGroup, &masstaskerpb.Test{
			Foo: fmt.Sprint(i),
		})
		if err != nil {
			return err
		}
		tasks = append(tasks, t)
	}

	slog.Info("Creating tasks", "n", len(tasks))
	if err := c.mt.Do(ctx, masstasker.Delete(bootstrap), masstasker.Create(tasks...)); err != nil {
		return err
	}
	return nil
}
