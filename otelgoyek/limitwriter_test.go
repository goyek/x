package otelgoyek

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestLimitWriter_ConcurrentWriteAndString(t *testing.T) {
	const (
		writers         = 16
		writesPerWriter = 256
	)

	w := &limitWriter{limit: writers * writesPerWriter * 4}
	start := make(chan struct{})
	done := make(chan struct{})
	observedLength := make(chan int, 1)
	var ready sync.WaitGroup
	var writerWG sync.WaitGroup
	var readerWG sync.WaitGroup
	ready.Add(writers + 1)
	writerWG.Add(writers)
	readerWG.Add(1)

	go func() {
		defer readerWG.Done()
		ready.Done()
		<-start
		maxLength := 0
		for {
			select {
			case <-done:
				observedLength <- max(maxLength, len(w.String()))
				return
			default:
				maxLength = max(maxLength, len(w.String()))
				runtime.Gosched()
			}
		}
	}()

	for writerID := range writers {
		record := fmt.Appendf(nil, "%02d:\n", writerID)
		go func() {
			defer writerWG.Done()
			ready.Done()
			<-start
			for range writesPerWriter {
				if n, err := w.Write(record); err != nil || n != len(record) {
					t.Errorf("Write() = (%d, %v), want (%d, nil)", n, err, len(record))
					return
				}
			}
		}()
	}

	ready.Wait()
	close(start)
	writerWG.Wait()
	close(done)
	readerWG.Wait()

	if got := <-observedLength; got == 0 {
		t.Error("String() did not observe any concurrently written output")
	}
	output := w.String()
	if got, want := len(output), writers*writesPerWriter*4; got != want {
		t.Errorf("captured output length = %d, want %d", got, want)
	}
	for writerID := range writers {
		record := fmt.Sprintf("%02d:\n", writerID)
		if got := strings.Count(output, record); got != writesPerWriter {
			t.Errorf("record %q count = %d, want %d", record, got, writesPerWriter)
		}
	}
}
