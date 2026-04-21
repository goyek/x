package otelgoyek_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/otelgoyek"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestMiddleware_WithDisableOutput(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exp)),
	)

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "secret message")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithDisableOutput(true)))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	for _, attr := range spans[0].Attributes {
		if string(attr.Key) == "goyek.task.output" {
			t.Errorf("found goyek.task.output attribute even though output capture is disabled: %v", attr.Value.AsString())
		}
	}
}

func TestExecutorMiddleware_WithDisableOutput(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exp)),
	)

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "secret flow message")
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithDisableOutput(true)))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
	// Usually 1 for Execute, but if we don't use runner middleware, it's just 1.
	// Actually Execute span will be there.
	var executeSpan *tracetest.SpanStub
	for _, s := range spans {
		if s.Name == "Execute" {
			executeSpan = &s
			break
		}
	}

	if executeSpan == nil {
		t.Fatal("Execute span not found")
	}

	for _, attr := range executeSpan.Attributes {
		if string(attr.Key) == "goyek.flow.output" {
			t.Errorf("found goyek.flow.output attribute even though output capture is disabled: %v", attr.Value.AsString())
		}
	}
}
