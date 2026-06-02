package main

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestRunExec_Masking(t *testing.T) {
	sb := &strings.Builder{}

	// Middleware to capture output
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
			runExec(a, "SECRET=password echo hello")
		},
	})
	f.Use(mw)

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	if strings.Contains(got, "password") {
		t.Errorf("Secret 'password' found in logs: %s", got)
	}
	if !strings.Contains(got, "SECRET=[MASKED]") {
		t.Errorf("Masked secret not found in logs: %s", got)
	}
}
