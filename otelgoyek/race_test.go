package otelgoyek_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/otelgoyek"
)

func TestMiddleware_OutputRace(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "race",
		Action: func(a *goyek.A) {
			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					fmt.Fprintf(a.Output(), "message %d\n", i)
				}(i)
			}
			wg.Wait()
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{"race"})

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
}

func TestExecutorMiddleware_OutputRace(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "race",
		Action: func(a *goyek.A) {
			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					fmt.Fprintf(a.Output(), "message %d\n", i)
				}(i)
			}
			wg.Wait()
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{"race"})

	spans := exp.GetSpans()
	found := false
	for _, s := range spans {
		if s.Name == "Execute" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Execute span not found")
	}
}
