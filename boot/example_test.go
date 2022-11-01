package boot_test

import (
	"flag"

	"github.com/goyek/goyek/v2"

	"github.com/goyek/x/boot"
)

func ExampleMain() {
	// define a flag used by a task
	msg := flag.String("msg", "hello world", `message to display by "hi" task`)

	// define a task printing the message (configurable via flag)
	goyek.Define(goyek.Task{
		Name:  "hi",
		Usage: "Greetings",
		Action: func(tf *goyek.TF) {
			tf.Log(*msg)
		},
	})

	// run the build pipeline
	boot.Main()
}
