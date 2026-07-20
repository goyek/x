package color_test

import (
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/goyek/goyek/v3"

	goyekcolor "github.com/goyek/x/color"
)

func TestReportFlowConcurrentRecordsAreAtomic(t *testing.T) {
	forceColor(t)

	const failedTask = "fail"
	wantErr := errors.New("flow failed")
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	executor := goyekcolor.ReportFlow(func(in goyek.ExecuteInput) error {
		entered <- struct{}{}
		<-release
		if in.Tasks[0] == failedTask {
			return wantErr
		}
		return nil
	})
	out := &recordingWriter{}

	var wg sync.WaitGroup
	for _, task := range []string{"pass", failedTask} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := executor(goyek.ExecuteInput{Output: out, Tasks: []string{task}})
			if task == failedTask && !errors.Is(err, wantErr) {
				t.Errorf("got error %v, want %v", err, wantErr)
			}
			if task == "pass" && err != nil {
				t.Errorf("got unexpected error: %v", err)
			}
		}()
	}
	<-entered
	<-entered
	close(release)
	wg.Wait()

	writes := out.snapshot()
	if len(writes) != 2 {
		t.Fatalf("got %d output writes, want two atomic records: %q", len(writes), writes)
	}
	var gotPass, gotFail bool
	for _, write := range writes {
		if !strings.HasSuffix(write, "s\n\x1b[22;0m") {
			t.Errorf("write is not a complete colored record: %q", write)
		}
		switch {
		case strings.HasPrefix(write, "\x1b[1;32mok\t"):
			gotPass = true
		case strings.HasPrefix(write, "\x1b[1;31mflow failed\t"):
			gotFail = true
		default:
			t.Errorf("unexpected flow record: %q", write)
		}
	}
	if !gotPass || !gotFail {
		t.Errorf("got pass=%t and fail=%t, want both records", gotPass, gotFail)
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
