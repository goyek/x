package main

import (
	"github.com/goyek/goyek/v3"
)

var mod = goyek.Define(goyek.Task{
	Name:  "mod",
	Usage: "go mod tidy",
	Action: func(a *goyek.A) {
		runExec(a, "go mod tidy", runDir(a, dirRoot))
		runExec(a, "go mod tidy", runDir(a, dirBuild))
	},
})
