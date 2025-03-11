package boot_test

import (
	"flag"

	"github.com/goyek/goyek/v2"

	"github.com/goyek/x/boot"
)

var msg = flag.String("msg", "hello world", `message to display by "hi" task`)

var _ = goyek.Define(goyek.Task{
	Name:  "hi",
	Usage: "Greetings",
	Action: func(a *goyek.A) {
		a.Log(*msg)
	},
})

func ExampleMain() {
	boot.Main()
}
