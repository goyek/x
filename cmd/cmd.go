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
type Option func(tf *goyek.TF, cmd *exec.Cmd)

// Exec runs
// It calls tf.Error[f] and returns false in case of any problems.
// Example usage:
//
//	cmd.Exec(tf, "FOO=foo BAR=baz ./foo --bar=baz", cmd.Dir("pkg"))
func Exec(tf *goyek.TF, cmdLine string, opts ...Option) bool {
	tf.Helper()

	envs, args, err := shellwords.ParseWithEnvs(cmdLine)
	if err != nil {
		tf.Error("parse command line: ", err)
		return false
	}

	cmd := tf.Cmd(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs...)
	for _, opt := range opts {
		opt(tf, cmd)
	}

	tf.Log("Exec: ", cmdLine)
	if err := cmd.Run(); err != nil {
		tf.Error(err)
		return false
	}
	return true
}

// Dir is an option to set the working directory.
func Dir(s string) Option {
	return func(tf *goyek.TF, cmd *exec.Cmd) {
		tf.Helper()
		tf.Log("Work dir: ", s)
		cmd.Dir = s
	}
}

// Env is an option to set an environment variable.
func Env(k, v string) Option {
	return func(tf *goyek.TF, cmd *exec.Cmd) {
		tf.Helper()
		env := k + "=" + v
		tf.Log("Env: ", env)
		cmd.Env = append(cmd.Env, env)
	}
}

// Stdin is an option to set the standard input.
func Stdin(r io.Reader) Option {
	return func(tf *goyek.TF, cmd *exec.Cmd) {
		cmd.Stdin = r
	}
}

// Stdout is an option to set the standard output.
func Stdout(w io.Writer) Option {
	return func(tf *goyek.TF, cmd *exec.Cmd) {
		cmd.Stdout = w
	}
}

// Stderr is an option to set the standard error.
func Stderr(w io.Writer) Option {
	return func(tf *goyek.TF, cmd *exec.Cmd) {
		cmd.Stderr = w
	}
}
