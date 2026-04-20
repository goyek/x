// Build is the build pipeline for this repository.
package main

import (
	"os"

	"github.com/goyek/goyek/v3"

	"github.com/goyek/x/boot"
	"github.com/goyek/x/cmd"
)

// Directories used in repository.
const (
	dirRoot  = "."
	dirBuild = "build"
)

func main() {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}
	goyek.SetDefault(all)
	boot.Main()
}

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
