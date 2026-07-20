package color_test

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/goyek/goyek/v3"

	goyekcolor "github.com/goyek/x/color"
)

func TestReportStatusWritesAtomicColoredRecords(t *testing.T) {
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
			out := &recordingWriter{}
			runner := goyekcolor.ReportStatus(func(goyek.Input) goyek.Result { return tc.result })
			got := runner(goyek.Input{Output: out, TaskName: "task"})
			if !reflect.DeepEqual(got, tc.result) {
				t.Fatalf("got result %#v, want %#v", got, tc.result)
			}

			writes := out.snapshot()
			if len(writes) != 2 {
				t.Fatalf("got %d output writes, want two atomic records: %q", len(writes), writes)
			}
			if writes[0] != ansiBlue+"===== TASK  task\n"+ansiReset {
				t.Errorf("unexpected task-start record: %q", writes[0])
			}
			wantStatusPrefix := tc.colorStart + "----- " + tc.status + ": task ("
			if !strings.HasPrefix(writes[1], wantStatusPrefix) || !strings.HasSuffix(writes[1], "s)\n"+ansiReset) {
				t.Errorf("unexpected task-status record: %q", writes[1])
			}
		})
	}
}

func TestReportStatusWritesPanicAsOneRecord(t *testing.T) {
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
			out := &recordingWriter{}
			result := goyek.Result{
				Status:     goyek.StatusFailed,
				PanicValue: tc.panicValue,
				PanicStack: []byte("stack\n"),
			}
			runner := goyekcolor.ReportStatus(func(goyek.Input) goyek.Result { return result })
			runner(goyek.Input{Output: out, TaskName: "task"})

			writes := out.snapshot()
			if len(writes) != 3 {
				t.Fatalf("got %d output writes, want start, status, and panic records: %q", len(writes), writes)
			}
			want := ansiRed + tc.wantHeader + ansiReset + "\n\n" + ansiRed + "stack\n" + ansiReset
			if writes[2] != want {
				t.Errorf("got panic record %q, want %q", writes[2], want)
			}
		})
	}
}

func TestReportStatusConcurrentRecordsAreAtomic(t *testing.T) {
	forceColor(t)

	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	runner := goyekcolor.ReportStatus(func(goyek.Input) goyek.Result {
		entered <- struct{}{}
		<-release
		return goyek.Result{Status: goyek.StatusPassed}
	})
	out := &recordingWriter{}
	var wg sync.WaitGroup
	for i := range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runner(goyek.Input{Output: out, TaskName: fmt.Sprintf("task-%d", i)})
		}()
	}
	<-entered
	<-entered
	close(release)
	wg.Wait()

	writes := out.snapshot()
	if len(writes) != 4 {
		t.Fatalf("got %d output writes, want four atomic records: %q", len(writes), writes)
	}
	for _, write := range writes {
		if !strings.HasPrefix(write, "\x1b[") || !strings.HasSuffix(write, ansiReset) {
			t.Errorf("write is not a complete colored record: %q", write)
		}
	}
}
