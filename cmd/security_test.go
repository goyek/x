package cmd

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestLogging_Security(t *testing.T) {
	sb := &strings.Builder{}

	f := &goyek.Flow{}
	f.SetOutput(sb)
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			Env("SECRET_KEY", "very-sensitive-value")(a, &exec.Cmd{})
			Exec(a, "SECRET=password echo hello")
		},
	})

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
	if !strings.Contains(got, "SECRET=[MASKED]") {
		t.Errorf("Inline secret was not masked: %s", got)
	}
}

func TestRunExec_NoLeak(t *testing.T) {
	sb := &strings.Builder{}

	f := &goyek.Flow{}
	f.SetOutput(sb)
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			cmdLine := "SECRET=password echo hello"
			Exec(a, cmdLine)
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := sb.String()
	if strings.Contains(got, "password") {
		t.Errorf("Secret value was logged: %s", got)
	}
	if !strings.Contains(got, "SECRET=[MASKED]") {
		t.Errorf("Secret was not masked: %s", got)
	}
}
