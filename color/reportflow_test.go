package color_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"

	goyekcolor "github.com/goyek/x/color"
)

func TestReportFlow(t *testing.T) {
	forceColor(t)

	wantErr := errors.New("flow failed")
	tests := []struct {
		name       string
		err        error
		wantPrefix string
	}{
		{name: "pass", wantPrefix: "\x1b[1;32mok\t"},
		{name: "fail", err: wantErr, wantPrefix: "\x1b[1;31mflow failed\t"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &strings.Builder{}
			executor := goyekcolor.ReportFlow(func(goyek.ExecuteInput) error { return tc.err })

			if err := executor(goyek.ExecuteInput{Output: out}); !errors.Is(err, tc.err) {
				t.Fatalf("got error %v, want %v", err, tc.err)
			}
			got := out.String()
			if !strings.HasPrefix(got, tc.wantPrefix) || !strings.HasSuffix(got, "s\n\x1b[22;0m") {
				t.Errorf("unexpected flow output: %q", got)
			}
		})
	}
}

func TestReportFlowNilOutput(t *testing.T) {
	var gotOutput io.Writer
	executor := goyekcolor.ReportFlow(func(in goyek.ExecuteInput) error {
		gotOutput = in.Output
		return nil
	})

	if err := executor(goyek.ExecuteInput{}); err != nil {
		t.Fatalf("executor returned error: %v", err)
	}
	if gotOutput != io.Discard {
		t.Fatalf("next executor received %T output, want io.Discard", gotOutput)
	}
}
