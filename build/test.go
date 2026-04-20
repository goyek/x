package main

import (
	"github.com/goyek/goyek/v3"

	"github.com/goyek/x/cmd"
)

var test = goyek.Define(goyek.Task{
	Name:  "test",
	Usage: "go test",
	Action: func(a *goyek.A) {
		cmdLine := "go test -race -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./..."
		a.Log(cmdLine)
		if !cmd.Exec(a, cmdLine) {
			return
		}
		cmdLine = "go tool cover -html=coverage.out -o coverage.html"
		a.Log(cmdLine)
		cmd.Exec(a, cmdLine)
	},
})
