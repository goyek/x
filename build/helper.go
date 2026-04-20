package main

import (
	"github.com/goyek/goyek/v3"

	"github.com/goyek/x/cmd"
)

func runExec(a *goyek.A, cmdLine string, opts ...cmd.Option) bool {
	a.Helper()
	a.Log("Exec: ", cmdLine)
	return cmd.Exec(a, cmdLine, opts...)
}

func runDir(a *goyek.A, s string) cmd.Option {
	a.Helper()
	a.Log("Work dir: ", s)
	return cmd.Dir(s)
}
