package color_test

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/goyek/goyek/v3"

	goyekcolor "github.com/goyek/x/color"
)

func TestCodeLineLoggerColoredMethodsWriteOneRecord(t *testing.T) {
	forceColor(t)

	tests := []struct {
		name       string
		colorStart string
		call       func(*goyekcolor.CodeLineLogger, io.Writer)
	}{
		{name: "Error", colorStart: ansiRed, call: func(l *goyekcolor.CodeLineLogger, w io.Writer) { l.Error(w, "message") }},
		{name: "Errorf", colorStart: ansiRed, call: func(l *goyekcolor.CodeLineLogger, w io.Writer) { l.Errorf(w, "%s", "message") }},
		{name: "Fatal", colorStart: ansiRed, call: func(l *goyekcolor.CodeLineLogger, w io.Writer) { l.Fatal(w, "message") }},
		{name: "Fatalf", colorStart: ansiRed, call: func(l *goyekcolor.CodeLineLogger, w io.Writer) { l.Fatalf(w, "%s", "message") }},
		{name: "Skip", colorStart: ansiYellow, call: func(l *goyekcolor.CodeLineLogger, w io.Writer) { l.Skip(w, "message") }},
		{name: "Skipf", colorStart: ansiYellow, call: func(l *goyekcolor.CodeLineLogger, w io.Writer) { l.Skipf(w, "%s", "message") }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &recordingWriter{}
			tc.call(&goyekcolor.CodeLineLogger{}, out)

			writes := out.snapshot()
			if len(writes) != 1 {
				t.Fatalf("got %d output writes, want one atomic record: %q", len(writes), writes)
			}
			if !strings.HasPrefix(writes[0], tc.colorStart) {
				t.Errorf("output %q does not start with %q", writes[0], tc.colorStart)
			}
			if !strings.Contains(writes[0], "message\n") {
				t.Errorf("output %q does not contain the log message", writes[0])
			}
			if !strings.HasSuffix(writes[0], ansiReset) {
				t.Errorf("output %q does not end with a color reset", writes[0])
			}
		})
	}
}

func TestCodeLineLoggerConcurrentColoredRecordsAreAtomic(t *testing.T) {
	forceColor(t)

	const recordCount = 32
	logger := &goyekcolor.CodeLineLogger{}
	out := &recordingWriter{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := range recordCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			logger.Errorf(out, "record-%02d", i)
		}()
	}
	close(start)
	wg.Wait()

	writes := out.snapshot()
	if len(writes) != recordCount {
		t.Fatalf("got %d output writes, want %d atomic records", len(writes), recordCount)
	}
	joined := strings.Join(writes, "")
	for i := range recordCount {
		message := fmt.Sprintf("record-%02d", i)
		if strings.Count(joined, message) != 1 {
			t.Errorf("message %q does not occur exactly once in %q", message, joined)
		}
	}
	for _, write := range writes {
		if !strings.HasPrefix(write, ansiRed) || !strings.HasSuffix(write, ansiReset) {
			t.Errorf("write is not a complete colored record: %q", write)
		}
	}
}

func TestCodeLineLoggerHelperAttribution(t *testing.T) {
	logger := &goyekcolor.CodeLineLogger{}
	out := &strings.Builder{}
	var wantLocation string
	runner := goyek.NewRunner(func(a *goyek.A) {
		_, file, line, _ := runtime.Caller(0)
		wantLocation = fmt.Sprintf("%s:%d", filepath.Base(file), line+3)

		logFromHelper(a)
	})

	runner(goyek.Input{Output: out, Logger: logger})

	if !strings.Contains(out.String(), wantLocation+": message from helper") {
		t.Errorf("output %q does not attribute helper log to %s", out.String(), wantLocation)
	}
}

func TestCodeLineLoggerHelperInActionAttribution(t *testing.T) {
	logger := &goyekcolor.CodeLineLogger{}
	out := &strings.Builder{}
	var wantLocation string
	runner := goyek.NewRunner(func(a *goyek.A) {
		a.Helper()
		_, file, line, _ := runtime.Caller(0)
		wantLocation = fmt.Sprintf("%s:%d", filepath.Base(file), line+3)

		a.Log("message from action")
	})

	runner(goyek.Input{Output: out, Logger: logger})

	if !strings.Contains(out.String(), wantLocation+": message from action") {
		t.Errorf("output %q does not attribute action log to %s", out.String(), wantLocation)
	}
}

func logFromHelper(a *goyek.A) {
	a.Helper()
	a.Log("message from helper")
}
