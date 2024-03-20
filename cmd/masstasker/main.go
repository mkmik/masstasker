package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/alecthomas/kong"
	taskmaster "mkm.pub/masstasker/pkg/proto"
	"mkm.pub/masstasker/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// set by goreleaser
var version = "(devel)"

type Context struct {
	*CLI
}

type CLI struct {
	ListenHTTP string `name:"listen-http" default:":8080"`
	ListenRPC  string `name:"listen-rpc" default:":50053"`

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
	taskmaster.RegisterTaskmasterServer(s, server.New())
	reflection.Register(s)

	go func() {
		log.Printf("HTTP server listening on %s", flags.ListenHTTP)
		if err := http.ListenAndServe(flags.ListenHTTP, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Printf("serving on %s", flags.ListenRPC)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Description(`Taskmaster`),
		kong.UsageOnError(),
		kong.Vars{"version": getVersion()},
		kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}),
	)

	err := ctx.Run(&Context{CLI: &cli})
	ctx.FatalIfErrorf(err)
}
