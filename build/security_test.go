package main

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestRunExec_Masking(t *testing.T) {
	var sb strings.Builder
	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			runExec(a, "SECRET=password echo hello")
		},
	})

	// Use a middleware to capture the output
	mw := func(next goyek.Runner) goyek.Runner {
		return func(in goyek.Input) goyek.Result {
			in.Output = io.MultiWriter(in.Output, &sb)
			return next(in)
		}
	}
	f.Use(mw)

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	if strings.Contains(got, "password") {
		t.Errorf("Secret value was logged: %s", got)
	}
	if !strings.Contains(got, "SECRET=[MASKED]") {
		t.Errorf("Secret was not masked: %s", got)
	}
}
