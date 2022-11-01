package cmd_test

import (
	"io"
	"strings"

	"github.com/goyek/goyek/v2"

	"github.com/goyek/x/cmd"
)

func ExampleExec() {
	goyek.Define(goyek.Task{
		Name:  "diff",
		Usage: "git diff",
		Action: func(tf *goyek.TF) {
			cmd.Exec(tf, "git diff --exit-code")

			sb := &strings.Builder{}
			out := io.MultiWriter(tf.Output(), sb)
			cmd.Exec(tf, "git status --porcelain", cmd.Stdout(out))
			if sb.Len() > 0 {
				tf.Error("git status --porcelain returned output")
			}
		},
	})
}
