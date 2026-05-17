package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestRunExec_Security(t *testing.T) {
	sb := &strings.Builder{}

	mw := func(next goyek.Runner) goyek.Runner {
		return func(in goyek.Input) goyek.Result {
			in.Output = io.MultiWriter(in.Output, sb)
			return next(in)
		}
	}

	f := &goyek.Flow{}
	f.Use(mw)
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			runExec(a, "SECRET=password echo hello")
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	fmt.Println("Captured output:", got)
	if strings.Contains(got, "password") {
		t.Errorf("Secret value from runExec was logged: %s", got)
	}
	if !strings.Contains(got, "SECRET=[MASKED]") {
		t.Errorf("Secret value from runExec was NOT masked: %s", got)
	}
}
