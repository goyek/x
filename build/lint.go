package main

import (
	"github.com/goyek/goyek/v3"
)

var lint = goyek.Define(goyek.Task{
	Name:  "lint",
	Usage: "golangci-lint run --fix",
	Action: func(a *goyek.A) {
		if !runExec(a, "go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint", runDir(a, dirBuild)) {
			return
		}
		runExec(a, "golangci-lint run --fix", runDir(a, dirRoot))
		runExec(a, "golangci-lint run --fix", runDir(a, dirBuild))
	},
})
