package server

import (
	"context"
	"fmt"
	"log"
	"math"
	"slices"
	"sync"
	"time"

	skipcommon "github.com/Workiva/go-datastructures/common"
	"github.com/Workiva/go-datastructures/slice/skip"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	masstasker "mkm.pub/masstasker/pkg/proto"
)

var (
	numTasks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "masstasker_num_tasks",
			Help: "Number of tasks",
		},
		[]string{"group"},
	)
)

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

// server is used to implement masstasker.MassTasker
type server struct {
	masstasker.UnimplementedMassTaskerServer

	sync.Mutex
	next   uint64
	tasks  map[uint64]*masstasker.Task
	labels map[string]map[uint64]struct{}
	groups map[string]*skip.SkipList
}

func New() *server {
	return &server{
		next:   1,
		tasks:  map[uint64]*masstasker.Task{},
		labels: map[string]map[uint64]struct{}{},
		groups: map[string]*skip.SkipList{},
	}
}

func (s *server) delete(i uint64) {
	t := s.tasks[i]
	numTasks.WithLabelValues(t.Group).Dec()

	for k, v := range t.Labels {
		key := labelKey(k, v)
		delete(s.labels[key], i)
	}

	s.groups[t.Group].Delete(&skipNode{t.NotBefore.AsTime(), t.Id})
	delete(s.tasks, i)
}

func (s *server) create(c *masstasker.Task) {
	numTasks.WithLabelValues(c.Group).Inc()

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

func sample[T any](s []T) []T {
	l := slices.Min([]int{1, len(s)})
	return s[:l]
}

func (s *server) Update(ctx context.Context, in *masstasker.UpdateRequest) (*masstasker.UpdateResponse, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("Update: created: %d, deleted: %d, predicates: %d. Sample: created: %v", len(in.Created), len(in.Deleted), len(in.Predicates), sample(in.Created))

	for _, id := range in.Predicates {
		if _, found := s.tasks[id]; !found {
			return nil, status.Errorf(codes.FailedPrecondition, "predicate id %d doesn't exist", id)
		}
	}

	var del []uint64
	for _, d := range in.Deleted {
		switch d := d.Sel.(type) {
		case *masstasker.TaskRef_Id:
			if _, exists := s.tasks[d.Id]; !exists {
				return nil, status.Errorf(codes.NotFound, "task %d doesn't exist", d)
			}
			del = append(del, d.Id)
		case *masstasker.TaskRef_Selector:
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
	return &masstasker.UpdateResponse{CreatedIds: res}, nil
}

func (s *server) Query(ctx context.Context, in *masstasker.QueryRequest) (*masstasker.QueryResponse, error) {
	now := in.Now.AsTime()
	if in.Now == nil {
		now = time.Now()
	}
	start := time.Now()

	for {
		now = now.Add(time.Since(start))
		res, d, err := s.query(in, now)
		if err != nil {
			return nil, err
		}
		if d > 0 {
			log.Printf("Found owned task in group:%q, sleeping %v", in.Group, d)
			if !in.Wait {
				return nil, status.Errorf(codes.NotFound, "cannot find any value visible at %q", now)
			}
			time.Sleep(d)
		} else {
			return res, nil
		}
	}
}

func (s *server) query(in *masstasker.QueryRequest, now time.Time) (*masstasker.QueryResponse, time.Duration, error) {
	s.Lock()
	defer s.Unlock()
	log.Printf("Received: %v", in)

	sk := s.groups[in.Group]
	if sk == nil {
		return nil, 0, status.Errorf(codes.NotFound, "empty group %q", in.Group)
	}

	it := sk.Iter(&skipNode{})
	if it.Value() == nil {
		return nil, 0, status.Errorf(codes.NotFound, "cannot find any value for group %q", in.Group)
	}
	v := it.Value().(*skipNode)
	if v.Compare(&skipNode{ts: now}) > 0 {
		ts := v.ts
		return nil, ts.Sub(now), nil
	}

	id := v.id
	t := s.tasks[id]

	if in.GetOwnFor().AsDuration() > 0 {
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

	return &masstasker.QueryResponse{Task: t}, 0, nil
}

func (s *server) Debug(ctx context.Context, in *masstasker.DebugRequest) (*masstasker.DebugResponse, error) {
	s.Lock()
	defer s.Unlock()

	var tasks []*masstasker.Task
	var numTasks uint64

	limit := uint64(math.MaxUint64)
	if in.Limit > 0 {
		limit = in.Limit
	}
	if group := in.Group; group != "" {
		if sk := s.groups[group]; sk != nil {
			it := sk.Iter(&skipNode{})
			numTasks = sk.Len()
			for it.Next() && limit > 0 {
				id := it.Value().(*skipNode).id
				tasks = append(tasks, s.tasks[id])
				limit--
			}
		}
	} else {
		numTasks = uint64(len(s.tasks))
		for _, t := range s.tasks {
			if limit == 0 {
				break
			}
			tasks = append(tasks, t)
			limit--
		}
	}
	return &masstasker.DebugResponse{Tasks: tasks, NumTasks: numTasks}, nil
}
