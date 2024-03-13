package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"net/http"

	skipcommon "github.com/Workiva/go-datastructures/common"
	"github.com/Workiva/go-datastructures/slice/skip"
	"github.com/alecthomas/kong"
	taskmaster "github.com/mkmik/taskmaster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func labelKey(key, value string) string {
	return fmt.Sprintf("%s/%s", key, value)
}

type skipNode struct {
	ts time.Time
	id uint64
}

// Compare compares skipNodes with timestamp order.
func (a *skipNode) Compare(b skipcommon.Comparator) int {
	ta := a.ts
	tb := b.(*skipNode).ts

	if c := ta.Compare(tb); c != 0 {
		return c
	}

	aId, bId := a.id, b.(*skipNode).id
	if aId > bId {
		return 1
	} else if aId < bId {
		return -1
	}

	return 0
}

// server is used to implement taskmaster.Taskmaster
type server struct {
	taskmaster.UnimplementedTaskmasterServer

	sync.Mutex
	next   uint64
	tasks  map[uint64]*taskmaster.Task
	labels map[string]map[uint64]struct{}
	groups map[string]*skip.SkipList
}

func newServer() *server {
	return &server{
		next:   1,
		tasks:  map[uint64]*taskmaster.Task{},
		labels: map[string]map[uint64]struct{}{},
		groups: map[string]*skip.SkipList{},
	}
}

func (s *server) delete(i uint64) {
	t := s.tasks[i]
	for k, v := range t.Labels {
		key := labelKey(k, v)
		delete(s.labels[key], i)
	}

	s.groups[t.Group].Delete(&skipNode{t.NotBefore.AsTime(), t.Id})
	delete(s.tasks, i)
}

func (s *server) create(c *taskmaster.Task) {
	for k, v := range c.Labels {
		key := labelKey(k, v)
		if s.labels[key] == nil {
			s.labels[key] = map[uint64]struct{}{}
		}
		s.labels[key][c.Id] = struct{}{}
	}

	s.tasks[c.Id] = c
	if s.groups[c.Group] == nil {
		s.groups[c.Group] = skip.New(uint(10))
	}
	s.groups[c.Group].Insert(&skipNode{c.NotBefore.AsTime(), c.Id})
}

func (s *server) Update(ctx context.Context, in *taskmaster.UpdateRequest) (*taskmaster.UpdateResponse, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("Update: %v", in)

	for _, id := range in.Predicates {
		if _, found := s.tasks[id]; !found {
			return nil, status.Errorf(codes.FailedPrecondition, "predicate id %d doesn't exist", id)
		}
	}

	var del []uint64
	for _, d := range in.Deleted {
		switch d := d.Sel.(type) {
		case *taskmaster.TaskRef_Id:
			if _, exists := s.tasks[d.Id]; !exists {
				return nil, fmt.Errorf("task %d doesn't exist", d)
			}
			del = append(del, d.Id)
		case *taskmaster.TaskRef_Selector:
			for k, v := range d.Selector.Labels {
				key := labelKey(k, v)
				for i := range s.labels[key] {
					del = append(del, i)
				}
			}
		case nil:
			return nil, fmt.Errorf("one must be set")
		default:
			return nil, fmt.Errorf("unhandled case %T", d)
		}
	}

	for _, i := range del {
		s.delete(i)
	}

	res := make([]uint64, len(in.Created))
	for i, c := range in.Created {
		c.Id = s.next
		s.next++

		res[i] = c.Id

		if c.NotBefore == nil {
			c.NotBefore = timestamppb.Now()
		}

		s.create(c)
	}
	return &taskmaster.UpdateResponse{CreatedIds: res}, nil
}

func (s *server) Query(ctx context.Context, in *taskmaster.QueryRequest) (*taskmaster.QueryResponse, error) {
	for {
		now := in.Now
		if now == nil {
			now = timestamppb.Now()
		}
		res, d, err := s.query(in, now)
		if err != nil {
			return nil, err
		}
		if d > 0 {
			if !in.Wait {
				return nil, status.Errorf(codes.NotFound, "cannot find any value visible at %q", now)
			}
			time.Sleep(d)
		} else {
			return res, nil
		}
	}
}

func (s *server) query(in *taskmaster.QueryRequest, nowpb *timestamppb.Timestamp) (*taskmaster.QueryResponse, time.Duration, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("Received: %v", in)

	sk := s.groups[in.Group]
	if sk == nil {
		return nil, 0, status.Errorf(codes.NotFound, "empty group %q", in.Group)
	}

	it := sk.Iter(&skipNode{})
	v := it.Value().(*skipNode)
	if v == nil {
		return nil, 0, status.Errorf(codes.NotFound, "cannot find any value for group %q", in.Group)
	}
	now := nowpb.AsTime()
	if v.Compare(&skipNode{ts: now}) > 0 {
		ts := v.ts
		return nil, ts.Sub(now), nil
	}

	id := v.id
	t := s.tasks[id]

	if in.OwnFor != nil {
		d := in.OwnFor.AsDuration()
		nt := t.Clone()
		tp := timestamppb.New(time.Now().Add(d))
		nt.NotBefore = tp

		s.delete(nt.Id)
		nt.Id = s.next
		s.next++
		s.create(nt)
		t = nt
	}

	return &taskmaster.QueryResponse{Task: t}, 0, nil
}

func (s *server) Debug(ctx context.Context, in *taskmaster.DebugRequest) (*taskmaster.DebugResponse, error) {
	s.Lock()
	defer s.Unlock()

	for _, t := range s.tasks {
		log.Printf("%v", t)
	}
	return &taskmaster.DebugResponse{}, nil
}

func (flags *CLI) Run(*Context) error {
	lis, err := net.Listen("tcp", flags.ListenRPC)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpc.EnableTracing = true
	s := grpc.NewServer()
	taskmaster.RegisterTaskmasterServer(s, newServer())
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

func getVersion() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	// otherwise fallback to the version set by goreleaser
	return version
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
