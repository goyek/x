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
		t.Error("FOO=bar should be present even if ClearEnv() is used, as it was specified in the command line")
	}
	// Check that other environment variables are NOT present (e.g. PATH)
	if strings.Contains(got, "PATH=") {
		t.Error("PATH should not be present when ClearEnv() is used")
	}
}

func TestExec_Env_InlineEnv_Precedence(t *testing.T) {
	f := &goyek.Flow{}
	var output strings.Builder
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			// Inline FOO=inline should win over Env("FOO", "option")
			Exec(a, "FOO=inline env", Env("FOO", "option"), Stdout(&output))
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := output.String()
	if !strings.Contains(got, "FOO=inline") {
		t.Errorf("FOO=inline should be present, got: %s", got)
	}
	if strings.Contains(got, "FOO=option") {
		t.Error("FOO=option should have been overridden by inline variable")
	}
}

func TestEnv_Inheritance(t *testing.T) {
	t.Setenv("GOYEK_INHERITED", "true")

	a := &goyek.A{}
	c := &exec.Cmd{} // Env is nil

	Env("FOO", "bar")(a, c)

	foundInherited := false
	foundNew := false
	for _, e := range c.Env {
		if e == "GOYEK_INHERITED=true" {
			foundInherited = true
		}
		if e == "FOO=bar" {
			foundNew = true
		}
	}

	if !foundInherited {
		t.Error("expected inherited environment to be preserved by Env option")
	}
	if !foundNew {
		t.Error("expected new environment variable to be set by Env option")
	}
}

func TestEnv_NilInheritance(t *testing.T) {
	// This tests the case where cmd.Env is nil and we call Env option.
	// It should initialize cmd.Env with os.Environ().
	t.Setenv("GOYEK_TEST_VAR", "present")
	a := &goyek.A{}
	cmd := &exec.Cmd{
		Env: nil,
	}

	Env("OTHER_VAR", "value")(a, cmd)

	if cmd.Env == nil {
		t.Fatal("expected Env not to be nil")
	}

	foundInherited := false
	foundOther := false
	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "GOYEK_TEST_VAR=") {
			foundInherited = true
		}
		if e == "OTHER_VAR=value" {
			foundOther = true
		}
	}
	if !foundInherited {
		t.Error("Inherited environment missing")
	}
	if !foundOther {
		t.Error("New environment variable missing")
	}
}
