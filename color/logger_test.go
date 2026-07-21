package color_test

import (
	"context"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"

	goyekcolor "github.com/goyek/x/color"
)

func TestCodeLineLogger(t *testing.T) {
	flow := &goyek.Flow{}
	out := &strings.Builder{}
	flow.SetOutput(out)
	flow.SetLogger(&goyekcolor.CodeLineLogger{})
	flow.Define(goyek.Task{
		Name: "task",
		Action: func(a *goyek.A) {
			a.Log("message")
			helperFn(a)
			a.Cleanup(func() {
				a.Log("cleanup")
			})
		},
	})

	_ = flow.Execute(context.Background(), []string{"task"})

	for _, want := range []string{
		"      logger_test.go:21: message",
		"      logger_test.go:22: message from helper",
		"      logger_test.go:24: cleanup",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output %q does not contain %q", out.String(), want)
		}
	}
}

func TestCodeLineLogger_helper_in_action(t *testing.T) {
	flow := &goyek.Flow{}
	out := &strings.Builder{}
	flow.SetOutput(out)
	flow.SetLogger(&goyekcolor.CodeLineLogger{})
	flow.Define(goyek.Task{
		Name: "task",
		Action: func(a *goyek.A) {
			a.Helper()
			a.Log("message")
		},
	})

	_ = flow.Execute(context.Background(), []string{"task"})

	want := "      logger_test.go:51: message"
	if !strings.Contains(out.String(), want) {
		t.Errorf("output %q does not contain %q", out.String(), want)
	}
}

func helperFn(a *goyek.A) {
	a.Helper()
	a.Log("message from helper")
}
