package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	taskmaster "github.com/mkmik/taskmaster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

// server is used to implement taskmaster.Taskmaster
type server struct {
	taskmaster.UnimplementedTaskmasterServer

	sync.Mutex
	next   int64
	tasks  map[int64]taskmaster.Task
	labels map[string]map[int64]struct{}
}

func labelKey(key, value string) string {
	return fmt.Sprintf("%s/%s", key, value)
}

func newServer() *server {
	return &server{
		tasks:  map[int64]taskmaster.Task{},
		labels: map[string]map[int64]struct{}{},
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
		delete(s.tasks, i)
	}

	res := make([]int64, len(in.Created))
	for i, c := range in.Created {
		res[i] = s.next
		c.Id = s.next

		for k, v := range c.Labels {
			key := labelKey(k, v)
			if s.labels[key] == nil {
				s.labels[key] = map[int64]struct{}{}
			}
			s.labels[key][c.Id] = struct{}{}
		}

		s.tasks[c.Id] = *c

		s.next++
	}
	return &taskmaster.UpdateResponse{CreatedIds: res}, nil
}

func (s *server) Query(ctx context.Context, in *taskmaster.QueryRequest) (*taskmaster.QueryResponse, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("Received: %v", in)

	var res []*taskmaster.Task
	for k, v := range in.Selector.Labels {
		key := labelKey(k, v)
		for i := range s.labels[key] {
			t := s.tasks[i]
			res = append(res, &t)
		}
	}

	return &taskmaster.QueryResponse{Tasks: res}, nil
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
