package taskmaster

import "google.golang.org/protobuf/proto"

//go:generate protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative taskmaster.proto

func (x *Task) Clone() *Task {
	return proto.Clone(x).(*Task)
}
