package masstasker_test

import (
	"fmt"
	"log"

	"mkm.pub/masstasker"
	mypb "mkm.pub/masstasker/pkg/proto"
)

func ExampleTask_data() {
	task, err := masstasker.NewTask("my_group", &mypb.Test{Foo: "bar"})
	if err != nil {
		log.Fatal(err)
	}

	// ...

	var data mypb.Test
	if err := task.UnmarshalDataTo(&data); err != nil {
		log.Fatal(err)
	}

	fmt.Println(data.Foo)
	// Output:
	// bar
}
