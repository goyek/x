package otelgoyek

import (
	"strings"
	"sync"
	"testing"
)

func TestLimitWriter_ConcurrentWriteAndString(t *testing.T) {
	const writers = 10
	const message = "message"

	w := &limitWriter{limit: writers * len(message)}
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(writers + 1)

	go func() {
		defer wg.Done()
		<-start
		for range writers {
			_ = w.String()
		}
	}()

	for range writers {
		go func() {
			defer wg.Done()
			<-start
			if n, err := w.Write([]byte(message)); err != nil || n != len(message) {
				t.Errorf("Write() = (%d, %v), want (%d, nil)", n, err, len(message))
			}
		}()
	}

	close(start)
	wg.Wait()

	if got := strings.Count(w.String(), message); got != writers {
		t.Errorf("got %d messages, want %d", got, writers)
	}
}
