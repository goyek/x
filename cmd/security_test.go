package cmd

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestEnvLogging(t *testing.T) {
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
		},
	})

	// Use the middleware to capture output
	// goyek.Use is global, but Flow might have its own?
	// Actually goyek.Use affects DefaultFlow.
	// For a custom Flow, we might need another way or just use DefaultFlow.

	oldFlow := goyek.DefaultFlow
	defer func() { goyek.DefaultFlow = oldFlow }()
	goyek.DefaultFlow = f
	goyek.Use(mw)

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	if strings.Contains(got, "very-sensitive-value") {
		t.Errorf("Secret value was logged: %s", got)
	}
	if !strings.Contains(got, "SECRET_KEY") {
		t.Errorf("Env key was not logged: %s", got)
	}
}
