package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestRunExec_Security(t *testing.T) {
	var buf bytes.Buffer

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test-runexec",
		Action: func(a *goyek.A) {
			runExec(a, "SECRET_VAL=sensitive-password echo hello")
		},
	})

	f.Use(func(next goyek.Runner) goyek.Runner {
		return func(in goyek.Input) goyek.Result {
			in.Output = &buf
			return next(in)
		}
	})

	_ = f.Execute(context.Background(), []string{"test-runexec"})

	got := buf.String()
	if strings.Contains(got, "sensitive-password") {
		t.Errorf("Secret value was leaked in log: %q", got)
	}
	if !strings.Contains(got, "SECRET_VAL=[MASKED]") {
		t.Errorf("Expected logs to contain masked environment variable, but got: %q", got)
	}
}
