package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/goyek/goyek/v2"

	"github.com/goyek/x/cmd"
)

var mdlint = goyek.Define(goyek.Task{
	Name:  "mdlint",
	Usage: "markdownlint-cli (uses docker)",
	Action: func(a *goyek.A) {
		if _, err := exec.LookPath("docker"); err != nil {
			a.Skip(err)
		}
		curDir, err := os.Getwd()
		if err != nil {
			a.Fatal(err)
		}
		mdFiles := find(a, ".md")
		if len(mdFiles) == 0 {
			a.Skip("no .md files")
		}
		dockerImage := "ghcr.io/igorshubovych/markdownlint-cli:v0.37.0"
		cmd.Exec(a, "docker run --rm -v '"+curDir+":/workdir' "+dockerImage+" "+strings.Join(mdFiles, " "))
	},
})
