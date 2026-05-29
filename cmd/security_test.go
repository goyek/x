package cmd

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestLogging_Security(t *testing.T) {
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
			Env("SECRET_KEY", "very-sensitive-value")(a, &exec.Cmd{})
			Exec(a, "SECRET=password echo hello")
		},
	})

	oldFlow := goyek.DefaultFlow
	defer func() { goyek.DefaultFlow = oldFlow }()
	goyek.DefaultFlow = f
	goyek.Use(mw)

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
	if strings.Contains(got, "SECRET=") {
		t.Errorf("Inline secret was logged: %s", got)
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		name    string
		cmdLine string
		want    string
	}{
		{
			name:    "no env vars",
			cmdLine: "echo hello",
			want:    "echo hello",
		},
		{
			name:    "one env var",
			cmdLine: "FOO=bar echo hello",
			want:    "FOO=[MASKED] echo hello",
		},
		{
			name:    "multiple env vars",
			cmdLine: "FOO=bar BAZ=qux echo hello",
			want:    "FOO=[MASKED] BAZ=[MASKED] echo hello",
		},
		{
			name:    "env var with space",
			cmdLine: "FOO='bar baz' echo hello",
			want:    "FOO=[MASKED] echo hello",
		},
		{
			name:    "argument with space",
			cmdLine: "echo 'hello world'",
			want:    "echo 'hello world'",
		},
		{
			name:    "malformed",
			cmdLine: "echo 'hello",
			want:    "echo 'hello",
		},
		{
			name:    "empty",
			cmdLine: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.cmdLine); got != tt.want {
				t.Errorf("Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
