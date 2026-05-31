package cmd

import (
	"strings"

	"github.com/mattn/go-shellwords"
)

// Mask replaces the values of leading environment variable assignments with [MASKED].
func Mask(cmdLine string) string {
	envs, args, err := shellwords.ParseWithEnvs(cmdLine)
	if err != nil || len(envs) == 0 {
		return cmdLine
	}

	var sb strings.Builder
	for _, env := range envs {
		k, _, found := strings.Cut(env, "=")
		if found {
			sb.WriteString(k)
			sb.WriteString("=[MASKED] ")
		} else {
			sb.WriteString(env)
			sb.WriteString(" ")
		}
	}

	for i, arg := range args {
		if i > 0 {
			sb.WriteString(" ")
		}
		if strings.Contains(arg, " ") {
			sb.WriteString("'")
			sb.WriteString(arg)
			sb.WriteString("'")
		} else {
			sb.WriteString(arg)
		}
	}

	return sb.String()
}
