package cmd

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestUnsetEnv(t *testing.T) {
	a := &goyek.A{}
	cmd := &exec.Cmd{
		Env: []string{"FOO=bar", "BAZ=qux"},
	}

	UnsetEnv("FOO")(a, cmd)

	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "FOO=") {
			t.Errorf("expected FOO to be unset, but got: %s", e)
		}
	}
	found := false
	for _, e := range cmd.Env {
		if e == "BAZ=qux" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected BAZ=qux to be preserved")
	}
}

func TestUnsetEnv_NoValue(t *testing.T) {
	a := &goyek.A{}
	cmd := &exec.Cmd{
		Env: []string{"FOO", "BAR=baz"},
	}

	UnsetEnv("FOO")(a, cmd)

	for _, e := range cmd.Env {
		if e == "FOO" {
			t.Errorf("expected FOO to be unset")
		}
	}
	found := false
	for _, e := range cmd.Env {
		if e == "BAR=baz" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected BAR=baz to be preserved")
	}
}

func TestUnsetEnv_Nil(t *testing.T) {
	t.Setenv("GOYEK_TEST_VAR", "present")
	a := &goyek.A{}
	cmd := &exec.Cmd{
		Env: nil,
	}

	UnsetEnv("GOYEK_TEST_VAR")(a, cmd)

	if cmd.Env == nil {
		t.Error("expected Env not to be nil")
	}
	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "GOYEK_TEST_VAR=") {
			t.Errorf("expected GOYEK_TEST_VAR to be unset, but got: %s", e)
		}
	}
}

func TestEnv_Inheritance(t *testing.T) {
	t.Setenv("GOYEK_TEST_VAR", "present")
	a := &goyek.A{}
	cmd := &exec.Cmd{}

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
		t.Error("inherited environment lost")
	}
	if !foundNew {
		t.Error("NEW_VAR=value not found")
	}
}

func TestClearEnv(t *testing.T) {
	a := &goyek.A{}
	cmd := &exec.Cmd{
		Env: []string{"FOO=bar", "BAZ=qux"},
	}

	ClearEnv()(a, cmd)

	if cmd.Env == nil || len(cmd.Env) != 0 {
		t.Errorf("expected empty env, but got: %v", cmd.Env)
	}
}

func TestExec_ClearEnv(t *testing.T) {
	t.Setenv("GOYEK_TEST_VAR", "present")

	f := &goyek.Flow{}
	var output strings.Builder
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			Exec(a, "env", ClearEnv(), Stdout(&output))
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := output.String()
	if strings.Contains(got, "GOYEK_TEST_VAR=present") {
		t.Error("GOYEK_TEST_VAR should not be present in a cleared environment")
	}
}

func TestExec_ClearEnv_WithEnv(t *testing.T) {
	t.Setenv("GOYEK_HIDDEN", "secret")
	f := &goyek.Flow{}
	var output strings.Builder
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			Exec(a, "env", ClearEnv(), Env("NEW_VAR", "value"), Stdout(&output))
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := output.String()
	if !strings.Contains(got, "NEW_VAR=value") {
		t.Error("NEW_VAR=value should be present")
	}

	if strings.Contains(got, "GOYEK_HIDDEN=") {
		t.Error("GOYEK_HIDDEN should not be present")
	}
}

func TestExec_ClearEnv_WithInlineEnv(t *testing.T) {
	t.Setenv("GOYEK_HIDDEN", "secret")
	f := &goyek.Flow{}
	var output strings.Builder
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			Exec(a, "INLINE_VAR=value env", ClearEnv(), Stdout(&output))
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := output.String()
	if !strings.Contains(got, "INLINE_VAR=value") {
		t.Error("INLINE_VAR=value should be present even if environment was cleared")
	}

	if strings.Contains(got, "GOYEK_HIDDEN=") {
		t.Error("GOYEK_HIDDEN should not be present")
	}
}

func TestExec_UnsetEnv(t *testing.T) {
	t.Setenv("GOYEK_TEST_VAR", "present")
	t.Setenv("ANOTHER_VAR", "stay")

	f := &goyek.Flow{}
	var output strings.Builder
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			Exec(a, "env", UnsetEnv("GOYEK_TEST_VAR"), Stdout(&output))
		},
	})

	_ = f.Execute(context.Background(), []string{"test"})

	got := output.String()
	if strings.Contains(got, "GOYEK_TEST_VAR=present") {
		t.Error("GOYEK_TEST_VAR should have been unset")
	}
	if !strings.Contains(got, "ANOTHER_VAR=stay") {
		t.Error("ANOTHER_VAR should still be present")
	}
}

func TestExec_NoCommand(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Exec did not panic when no command was provided")
		}
	}()
	Exec(&goyek.A{}, "")
}

func TestExec_EnvOnly(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Exec did not panic when only env vars were provided")
		}
	}()
	Exec(&goyek.A{}, "FOO=bar")
}
