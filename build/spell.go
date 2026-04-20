package main

import (
	"strings"

	"github.com/goyek/goyek/v3"

	"github.com/goyek/x/cmd"
)

var spell = goyek.Define(goyek.Task{
	Name:  "spell",
	Usage: "misspell",
	Action: func(a *goyek.A) {
		a.Log("go install misspell")
		if !cmd.Exec(a, "go install github.com/client9/misspell/cmd/misspell", cmd.Dir(dirBuild)) {
			return
		}
		mdFiles := find(a, ".md")
		if len(mdFiles) == 0 {
			a.Skip("no .md files")
		}
		cmdLine := "misspell -error -locale=US -w " + strings.Join(mdFiles, " ")
		a.Log(cmdLine)
		cmd.Exec(a, cmdLine)
	},
})
