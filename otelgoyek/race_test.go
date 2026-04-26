package otelgoyek

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestLimitWriter_Race(_ *testing.T) {
	lw := &limitWriter{
		sb:    &strings.Builder{},
		limit: 10000,
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = fmt.Fprintf(lw, "goroutine %d, line %d\n", id, j)
			}
		}(i)
	}
	wg.Wait()
}
