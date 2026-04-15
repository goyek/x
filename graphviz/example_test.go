package graphviz_test

import (
	"os"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/graphviz"
)

func ExampleDraw() {
	f := &goyek.Flow{}
	test := f.Define(goyek.Task{
		Name:  "test",
		Usage: "run tests",
	})
	f.Define(goyek.Task{
		Name:  "all",
		Usage: "run all tasks",
		Deps:  goyek.Deps{test},
	})

	if err := graphviz.Draw(os.Stdout, f); err != nil {
		panic(err)
	}

	// Output:
	// digraph G {
	//   "all";
	//   "all" -> "test";
	//   "test";
	// }
}
