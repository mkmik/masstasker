package masstasker

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// deprecated: use [Task.MarshalDataFrom]
func (t *Task) MarshalFrom(data proto.Message) error {
	return t.MarshalDataFrom(data)
}

func (t *Task) MarshalDataFrom(data proto.Message) error {
	if t.Data == nil {
		t.Data = &anypb.Any{}
	}
	return t.Data.MarshalFrom(data)
}

// deprecated: use [Task.UnmarshalDataTo]
func (t *Task) UnmarshalTo(dst proto.Message) error {
	return t.UnmarshalDataTo(dst)
}

func (t *Task) UnmarshalDataTo(dst proto.Message) error {
	return t.Data.UnmarshalTo(dst)
}
