package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	skipcommon "github.com/Workiva/go-datastructures/common"
	"github.com/Workiva/go-datastructures/slice/skip"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	taskmaster "github.com/mkmik/taskmaster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Flags struct {
	ListenRPC string
}

func (f *Flags) Bind(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}
	fs.StringVar(&f.ListenRPC, "listen-rpc", ":50053", "RPC listen address")
}

func labelKey(key, value string) string {
	return fmt.Sprintf("%s/%s", key, value)
}

type skipNode struct {
	ts timestamp.Timestamp
	id int64
}

// Compare compares skipNodes with reverse timestamp order.
func (a skipNode) Compare(b skipcommon.Comparator) int {
	ta := a.ts
	tb := b.(skipNode).ts

	aSecs, bSecs := ta.Seconds, tb.Seconds
	if aSecs < bSecs {
		return 1
	} else if aSecs > bSecs {
		return -1
	}

	aNanos, bNanos := ta.Nanos, tb.Nanos
	if aNanos < bNanos {
		return 1
	} else if aNanos > bNanos {
		return -1
	}

	aId, bId := a.id, b.(skipNode).id
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
	next   int64
	tasks  map[int64]taskmaster.Task
	labels map[string]map[int64]struct{}
	groups map[string]*skip.SkipList
}

func newServer() *server {
	return &server{
		tasks:  map[int64]taskmaster.Task{},
		labels: map[string]map[int64]struct{}{},
		groups: map[string]*skip.SkipList{},
	}
}

func (s *server) Update(ctx context.Context, in *taskmaster.UpdateRequest) (*taskmaster.UpdateResponse, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("Update: %v", in)

	var del []int64
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
		t := s.tasks[i]
		for k, v := range t.Labels {
			key := labelKey(k, v)
			delete(s.labels[key], i)
		}

		s.groups[t.Group].Delete(skipNode{*t.NotBefore, t.Id})
		delete(s.tasks, i)
	}

	res := make([]int64, len(in.Created))
	for i, c := range in.Created {
		c.Id = s.next
		s.next++

		res[i] = c.Id

		if c.NotBefore == nil {
			c.NotBefore = ptypes.TimestampNow()
		}

		for k, v := range c.Labels {
			key := labelKey(k, v)
			if s.labels[key] == nil {
				s.labels[key] = map[int64]struct{}{}
			}
			s.labels[key][c.Id] = struct{}{}
		}

		s.tasks[c.Id] = *c
		if s.groups[c.Group] == nil {
			s.groups[c.Group] = skip.New(uint(10))
		}
		s.groups[c.Group].Insert(skipNode{*c.NotBefore, c.Id})
	}
	return &taskmaster.UpdateResponse{CreatedIds: res}, nil
}

func (s *server) Query(ctx context.Context, in *taskmaster.QueryRequest) (*taskmaster.QueryResponse, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("Received: %v", in)
	if in.Now == nil {
		in.Now = ptypes.TimestampNow()
	}

	sk := s.groups[in.Group]
	if sk == nil {
		return nil, status.Errorf(codes.NotFound, "empty group %q", in.Group)
	}

	it := sk.Iter(skipNode{*in.Now, 0})
	v := it.Value()
	if v == nil {
		return nil, status.Errorf(codes.NotFound, "cannot find any value visible at %q", in.Now)
	}
	id := v.(skipNode).id
	t := s.tasks[id]
	return &taskmaster.QueryResponse{Task: &t}, nil
}

func (s *server) Debug(ctx context.Context, in *taskmaster.DebugRequest) (*taskmaster.DebugResponse, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("debug:")

	for _, t := range s.tasks {
		log.Printf("%v", t)
	}
	return &taskmaster.DebugResponse{}, nil
}

func mainE(flags Flags) error {
	lis, err := net.Listen("tcp", flags.ListenRPC)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	taskmaster.RegisterTaskmasterServer(s, newServer())
	reflection.Register(s)

	log.Printf("serving on %s", flags.ListenRPC)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func main() {
	var flags Flags
	flags.Bind(nil)
	flag.Parse()

	if err := mainE(flags); err != nil {
		log.Fatal(err)
	}
}