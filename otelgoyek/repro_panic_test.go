package otelgoyek_test

import (
	"context"
	"testing"

	"github.com/goyek/goyek/v3"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/goyek/x/otelgoyek"
)

func TestMiddleware_PanicStatus_Repro(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exp)),
	)

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "panic",
		Action: func(_ *goyek.A) {
			panic("something went wrong")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{"panic"})

	spans := exp.GetSpans()
	var taskSpan *tracetest.SpanStub
	for _, s := range spans {
		if s.Name == "panic" {
			taskSpan = &s
			break
		}
	}

	if taskSpan == nil {
		t.Fatal("task span not found")
	}

	if taskSpan.Status.Code != codes.Error {
		t.Errorf("expected span status Error, got %v", taskSpan.Status.Code)
	}

	expectedMsg := "task panicked: panic"
	if taskSpan.Status.Description != expectedMsg {
		t.Errorf("expected span status description %q, got %q", expectedMsg, taskSpan.Status.Description)
	}
}
