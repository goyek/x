// Package cmd offers functions for running programs in a Shell-like way.
package cmd

import (
	"io"
	"os"
	"os/exec"

	"github.com/goyek/goyek/v2"
	"github.com/mattn/go-shellwords"
)

// Option configures the command.
type Option func(a *goyek.A, cmd *exec.Cmd)

// Exec runs the command.
// It calls a.Error[f] and returns false in case of any problems.
// Example usage:
//
//	cmd.Exec(a, "FOO=foo BAR=baz ./foo --bar=baz", cmd.Dir("pkg"))
func Exec(a *goyek.A, cmdLine string, opts ...Option) bool {
	a.Helper()

	envs, args, err := shellwords.ParseWithEnvs(cmdLine)
	if err != nil {
		a.Error("parse command line: ", err)
		return false
	}

	cmd := exec.CommandContext(a.Context(), args[0], args[1:]...) //nolint:gosec // it is a convenient function to run programs
	cmd.Stdin = os.Stdin
	cmd.Stdout = a.Output()
	cmd.Stderr = a.Output()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs...)
	for _, opt := range opts {
		opt(a, cmd)
	}

	a.Log("Exec: ", cmdLine)
	if err := cmd.Run(); err != nil {
		a.Error(err)
		return false
	}
	return true
}

// Dir is an option to set the working directory.
func Dir(s string) Option {
	return func(a *goyek.A, cmd *exec.Cmd) {
		a.Helper()
		a.Log("Work dir: ", s)
		cmd.Dir = s
	}
}

// Env is an option to set an environment variable.
func Env(k, v string) Option {
	return func(a *goyek.A, cmd *exec.Cmd) {
		a.Helper()
		env := k + "=" + v
		a.Log("Env: ", env)
		cmd.Env = append(cmd.Env, env)
	}
}

// Stdin is an option to set the standard input.
func Stdin(r io.Reader) Option {
	return func(a *goyek.A, cmd *exec.Cmd) {
		cmd.Stdin = r
	}
}

// Stdout is an option to set the standard output.
func Stdout(w io.Writer) Option {
	return func(a *goyek.A, cmd *exec.Cmd) {
		cmd.Stdout = w
	}
}

// Stderr is an option to set the standard error.
func Stderr(w io.Writer) Option {
	return func(a *goyek.A, cmd *exec.Cmd) {
		cmd.Stderr = w
	}
}
