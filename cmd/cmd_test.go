package cmd

import (
	"context"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestExec_noCommand(t *testing.T) {
	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			if Exec(a, "VAR=val") {
				a.Error("Exec should have failed for empty args")
			}
			if Exec(a, "") {
				a.Error("Exec should have failed for empty string")
			}
		},
	})
	err := f.Execute(context.Background(), []string{"test"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
