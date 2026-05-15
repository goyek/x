package main

import (
	"strings"

	"github.com/goyek/goyek/v3"
	"github.com/mattn/go-shellwords"

	"github.com/goyek/x/cmd"
)

func runExec(a *goyek.A, cmdLine string, opts ...cmd.Option) bool {
	a.Helper()

	msg := cmdLine
	envs, args, err := shellwords.ParseWithEnvs(cmdLine)
	if err == nil && len(envs) > 0 {
		var sb strings.Builder
		for _, env := range envs {
			if split := strings.SplitN(env, "=", 2); len(split) == 2 {
				sb.WriteString(split[0])
				sb.WriteString("=[MASKED] ")
			}
		}
		sb.WriteString(strings.Join(args, " "))
		msg = sb.String()
	}

	a.Log("Exec: ", msg)
	return cmd.Exec(a, cmdLine, opts...)
}

func runDir(a *goyek.A, s string) cmd.Option {
	a.Helper()
	a.Log("Work dir: ", s)
	return cmd.Dir(s)
}
