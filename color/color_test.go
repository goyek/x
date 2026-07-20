package color_test

import (
	"sync"
	"testing"

	fatihcolor "github.com/fatih/color"
)

const (
	ansiBlue   = "\x1b[34m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiReset  = "\x1b[0m"
	ansiYellow = "\x1b[33m"
)

type recordingWriter struct {
	mu     sync.Mutex
	writes []string
}

func (w *recordingWriter) Write(p []byte) (int, error) {
	w.record(string(p))
	return len(p), nil
}

func (w *recordingWriter) WriteString(s string) (int, error) {
	w.record(s)
	return len(s), nil
}

func (w *recordingWriter) record(s string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writes = append(w.writes, s)
}

func (w *recordingWriter) snapshot() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.writes...)
}

func forceColor(t *testing.T) {
	t.Helper()
	t.Setenv("NO_COLOR", "")
	oldNoColor := fatihcolor.NoColor
	fatihcolor.NoColor = false
	t.Cleanup(func() {
		fatihcolor.NoColor = oldNoColor
	})
}
