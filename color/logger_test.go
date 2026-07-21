package color_test

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"

	goyekcolor "github.com/goyek/x/color"
)

func TestCodeLineLoggerColoredMethods(t *testing.T) {
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
			out := &strings.Builder{}
			tc.call(&goyekcolor.CodeLineLogger{}, out)

			got := out.String()
			if !strings.HasPrefix(got, tc.colorStart) {
				t.Errorf("output %q does not start with %q", got, tc.colorStart)
			}
			if !strings.Contains(got, "message\n") {
				t.Errorf("output %q does not contain the log message", got)
			}
			if !strings.HasSuffix(got, ansiReset) {
				t.Errorf("output %q does not end with a color reset", got)
			}
		})
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

	runner(goyek.Input{Output: goyek.SyncWriter(out), Logger: logger})

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

	runner(goyek.Input{Output: goyek.SyncWriter(out), Logger: logger})

	if !strings.Contains(out.String(), wantLocation+": message from action") {
		t.Errorf("output %q does not attribute action log to %s", out.String(), wantLocation)
	}
}

func logFromHelper(a *goyek.A) {
	a.Helper()
	a.Log("message from helper")
}
