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
		Action: func(a *goyek.A) {
			cmd.Exec(a, "git diff --exit-code")

			sb := &strings.Builder{}
			out := io.MultiWriter(a.Output(), sb)
			cmd.Exec(a, "git status --porcelain", cmd.Stdout(out))
			if sb.Len() > 0 {
				a.Error("git status --porcelain returned output")
			}
		},
	})
}
