package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	mapi "mkm.pub/masstasker"
	masstasker "mkm.pub/masstasker/pkg/proto"
	"mkm.pub/masstasker/pkg/server"
)

// set by goreleaser
var version = "(devel)"

type Context struct {
	*CLI
}

type CLI struct {
	ListenHTTP string `name:"listen-http" default:":8080"`
	ListenRPC  string `name:"listen-rpc" default:":50053"`

	BootstrapTaskGroup string `name:"bootstrap-task-group" help:"Insert an empty task once in the given group at first start"`

	Version kong.VersionFlag `name:"version" help:"Print version information and quit"`
}

func getVersion() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	// otherwise fallback to the version set by goreleaser
	return version
}

func (flags *CLI) Run(*Context) error {
	lis, err := net.Listen("tcp", flags.ListenRPC)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpc.EnableTracing = true
	s := grpc.NewServer()
	masstasker.RegisterMassTaskerServer(s, server.New())
	reflection.Register(s)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		slog.Info("HTTP server listening", "addr", flags.ListenHTTP)
		if err := http.ListenAndServe(flags.ListenHTTP, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	if flags.BootstrapTaskGroup != "" {
		go insertBootstrapTask(flags.ListenRPC, flags.BootstrapTaskGroup)
	}

	log.Printf("serving on %s", flags.ListenRPC)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func insertBootstrapTask(addr string, group string) {
	ctx := context.Background()
	mt, err := mapi.Dial(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	if err := mt.Create(ctx, &masstasker.Task{Group: group}); err != nil {
		log.Fatal(err)
	}
	slog.Info("created bootstrap task", "group", group)
}

var (
	timeFormat = "2006-01-02T15:04:05.000000000-07:00"
)

func init() {
	location, _ := time.LoadLocation("Local")
	if location.String() == "UTC" {
		timeFormat = "2006-01-02T15:04:05.000000000Z"
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				// the timestamp in go1.22.1 is padded with zeros. The docs say it should be millisecond precision but it isn't.
				// In any case, let's format the time the way we want
				a.Value = slog.StringValue(a.Value.Time().Format(timeFormat))
			}
			return a
		},
	})))

	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Description(`MassTasker`),
		kong.UsageOnError(),
		kong.Vars{"version": getVersion()},
		kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}),
		kong.DefaultEnvars("MASSTASKER"),
	)

	err := ctx.Run(&Context{CLI: &cli})
	ctx.FatalIfErrorf(err)
}
