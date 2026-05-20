package cmd

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestLogging_Security(t *testing.T) {
	sb := &strings.Builder{}

	mw := func(next goyek.Runner) goyek.Runner {
		return func(in goyek.Input) goyek.Result {
			in.Output = io.MultiWriter(in.Output, sb)
			return next(in)
		}
	}

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			Env("SECRET_KEY", "very-sensitive-value")(a, &exec.Cmd{})
			Exec(a, "SECRET=password echo hello")
		},
	})

	oldFlow := goyek.DefaultFlow
	defer func() { goyek.DefaultFlow = oldFlow }()
	goyek.DefaultFlow = f
	goyek.Use(mw)

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	if strings.Contains(got, "very-sensitive-value") {
		t.Errorf("Secret value from Env was logged: %s", got)
	}
	if strings.Contains(got, "SECRET_KEY") {
		t.Errorf("Env key was logged: %s", got)
	}
	if strings.Contains(got, "password") {
		t.Errorf("Secret value from Exec was logged: %s", got)
	}
	if strings.Contains(got, "SECRET=") {
		t.Errorf("Inline secret was logged: %s", got)
	}
}

func TestMask_Security(t *testing.T) {
	sb := &strings.Builder{}
	// Simulated runExec from build/helper.go
	runExec := func(a *goyek.A, cmdLine string) {
		a.Log("Exec: ", Mask(cmdLine))
	}

	mw := func(next goyek.Runner) goyek.Runner {
		return func(in goyek.Input) goyek.Result {
			in.Output = io.MultiWriter(in.Output, sb)
			return next(in)
		}
	}

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			runExec(a, "PASSWORD=secret echo hello")
		},
	})

	oldFlow := goyek.DefaultFlow
	defer func() { goyek.DefaultFlow = oldFlow }()
	goyek.DefaultFlow = f
	goyek.Use(mw)

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	if strings.Contains(got, "secret") {
		t.Errorf("Secret value was logged: %s", got)
	}
	if !strings.Contains(got, "PASSWORD=[MASKED]") {
		t.Errorf("Secret was not masked in logs: %s", got)
	}
}
