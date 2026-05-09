package cmd

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestExec_ClearEnv_InlineEnv(t *testing.T) {
	f := &goyek.Flow{}
	var output strings.Builder
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			// Even with ClearEnv(), inline env should survive
			Exec(a, "FOO=bar env", ClearEnv(), Stdout(&output))
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := output.String()
	if !strings.Contains(got, "FOO=bar") {
		t.Error("FOO=bar should survive ClearEnv() when provided inline")
	}
}

func TestEnv_NilInheritance(t *testing.T) {
	t.Setenv("GOYEK_TEST_VAR", "present")

	c := &exec.Cmd{} // Env is nil
	a := &goyek.A{}

	Env("OTHER_VAR", "value")(a, c)

	foundTestVar := false
	for _, e := range c.Env {
		if strings.HasPrefix(e, "GOYEK_TEST_VAR=") {
			foundTestVar = true
			break
		}
	}

	if !foundTestVar {
		t.Error("Env() should have preserved os.Environ() if cmd.Env was nil")
	}

	foundOtherVar := false
	for _, e := range c.Env {
		if e == "OTHER_VAR=value" {
			foundOtherVar = true
			break
		}
	}
	if !foundOtherVar {
		t.Error("OTHER_VAR=value should be present")
	}
}
