package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mkm.pub/masstasker"
)

type Context struct {
	*CLI

	mt *masstasker.Client
}

type CLI struct {
	ServerAddr     string `name:"server-addr" default:":50053"`
	BootstrapGroup string `name:"bootstrap-group" default:"bootstrap"`
	MainGroup      string `name:"main-group" default:"main"`
	ErrorGroup     string `name:"error-group" default:"error"`

	OwnFor time.Duration `name:"own-for" default:"20s"`

	NumTasks   int `default:"1"`
	NumWorkers int `default:"2"`

	FakeWorkDuration time.Duration `default:"40s"`

	ListenHTTP string           `name:"listen-http" default:":8481" help:"For monitoring"`
	Version    kong.VersionFlag `name:"version" help:"Print version information and quit"`
}

// repeat calls cb repeatedly until it returns an error. It passes ctx through.
func repeat(ctx context.Context, cb func(context.Context) error) error {
	for {
		if err := cb(ctx); err != nil {
			return err
		}
	}
}

func (flags *CLI) Run(c *Context) error {
	go flags.startMonitoringServer()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	mt, err := masstasker.Dial(ctx, flags.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c.mt = mt

	go func() { cancel(fmt.Errorf("bootstrap: %w", repeat(ctx, c.processBootstrap))) }()

	for i := 0; i < c.NumWorkers; i++ {
		go func() { cancel(fmt.Errorf("worker %d: %w", i, repeat(ctx, c.processMainGroup))) }()
	}

	<-ctx.Done()

	return context.Cause(ctx)
}

func (flags *CLI) startMonitoringServer() {
	http.Handle("/metrics", promhttp.Handler())

	slog.Info("HTTP server listening", "addr", flags.ListenHTTP)
	if err := http.ListenAndServe(flags.ListenHTTP, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func main() {
	grpc.EnableTracing = true

	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Description(`MassTasker`),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}),
		kong.DefaultEnvars("MASSTASKER"),
	)

	err := ctx.Run(&Context{CLI: &cli})
	ctx.FatalIfErrorf(err)
}
