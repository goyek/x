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

	oldFlow := goyek.DefaultFlow
	defer func() { goyek.DefaultFlow = oldFlow }()
	goyek.DefaultFlow = f
	goyek.Use(mw)

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	if strings.Contains(got, "password") {
		t.Errorf("Secret value was logged: %s", got)
	}
	if !strings.Contains(got, "[MASKED]") {
		t.Errorf("Expected [MASKED] in logs, but got: %s", got)
	}
}
