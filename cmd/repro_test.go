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
			Exec(a, "FOO=bar env", ClearEnv(), Stdout(&output))
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := output.String()
	if !strings.Contains(got, "FOO=bar") {
		t.Error("FOO=bar should be present even if ClearEnv() is used")
	}
}

func TestEnv_NilInheritance(t *testing.T) {
	t.Setenv("GOYEK_TEST_VAR", "present")
	a := &goyek.A{}
	cmd := &exec.Cmd{
		Env: nil,
	}

	Env("NEW_VAR", "value")(a, cmd)

	foundInherited := false
	foundNew := false
	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "GOYEK_TEST_VAR=") {
			foundInherited = true
		}
		if e == "NEW_VAR=value" {
			foundNew = true
		}
	}

	if !foundInherited {
		t.Error("expected inherited environment variable to be present")
	}
	if !foundNew {
		t.Error("expected NEW_VAR=value to be present")
	}
}
