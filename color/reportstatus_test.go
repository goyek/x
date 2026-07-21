package color_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"

	goyekcolor "github.com/goyek/x/color"
)

func TestReportStatusColoredRecords(t *testing.T) {
	forceColor(t)

	tests := []struct {
		name       string
		result     goyek.Result
		colorStart string
		status     string
	}{
		{name: "passed", result: goyek.Result{Status: goyek.StatusPassed}, colorStart: ansiGreen, status: "PASS"},
		{name: "failed", result: goyek.Result{Status: goyek.StatusFailed}, colorStart: ansiRed, status: "FAIL"},
		{name: "skipped", result: goyek.Result{Status: goyek.StatusSkipped}, colorStart: ansiYellow, status: "SKIP"},
		{name: "not run", result: goyek.Result{Status: goyek.StatusNotRun}, colorStart: ansiGreen, status: "NOOP"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &strings.Builder{}
			runner := goyekcolor.ReportStatus(func(goyek.Input) goyek.Result { return tc.result })
			got := runner(goyek.Input{Output: out, TaskName: "task"})
			if !reflect.DeepEqual(got, tc.result) {
				t.Fatalf("got result %#v, want %#v", got, tc.result)
			}

			output := out.String()
			startRecord := ansiBlue + "===== TASK  task\n" + ansiReset
			if !strings.HasPrefix(output, startRecord) {
				t.Errorf("unexpected task-start record: %q", output)
			}
			wantStatusPrefix := tc.colorStart + "----- " + tc.status + ": task ("
			statusRecord := strings.TrimPrefix(output, startRecord)
			if !strings.HasPrefix(statusRecord, wantStatusPrefix) || !strings.HasSuffix(statusRecord, "s)\n"+ansiReset) {
				t.Errorf("unexpected task-status record: %q", statusRecord)
			}
		})
	}
}

func TestReportStatusNilOutput(t *testing.T) {
	var gotOutput io.Writer
	runner := goyekcolor.ReportStatus(func(in goyek.Input) goyek.Result {
		gotOutput = in.Output
		return goyek.Result{Status: goyek.StatusPassed}
	})

	result := runner(goyek.Input{TaskName: "task"})

	if result.Status != goyek.StatusPassed {
		t.Fatalf("got status %v, want %v", result.Status, goyek.StatusPassed)
	}
	if gotOutput != io.Discard {
		t.Fatalf("next runner received %T output, want io.Discard", gotOutput)
	}
}

func TestReportStatusPanic(t *testing.T) {
	forceColor(t)

	tests := []struct {
		name       string
		panicValue interface{}
		wantHeader string
	}{
		{name: "value", panicValue: "boom", wantHeader: "panic: boom"},
		{name: "nil value", wantHeader: "panic(nil) or runtime.Goexit() called"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &strings.Builder{}
			result := goyek.Result{
				Status:     goyek.StatusFailed,
				PanicValue: tc.panicValue,
				PanicStack: []byte("stack\n"),
			}
			runner := goyekcolor.ReportStatus(func(goyek.Input) goyek.Result { return result })
			runner(goyek.Input{Output: out, TaskName: "task"})

			want := ansiRed + tc.wantHeader + ansiReset + "\n\n" + ansiRed + "stack\n" + ansiReset
			if got := out.String(); !strings.HasSuffix(got, want) {
				t.Errorf("output %q does not end with panic report %q", got, want)
			}
		})
	}
}
