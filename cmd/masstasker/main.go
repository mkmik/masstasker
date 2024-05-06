package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/binarylog"
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
	grpc.EnableTracing = true
	sink, err := binarylog.NewTempFileSink()
	if err != nil {
		return err
	}
	binarylog.SetSink(sink)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		slog.Info("HTTP server listening", "addr", flags.ListenHTTP)
		if err := http.ListenAndServe(flags.ListenHTTP, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Printf("serving on %s", flags.ListenRPC)
	return server.ListenAndServe(context.Background(), flags.ListenRPC, server.WithBootstrapTaskGroup(flags.BootstrapTaskGroup))
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
