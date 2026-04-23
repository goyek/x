package otelgoyek_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/goyek/goyek/v3"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/goyek/x/otelgoyek"
)

func TestMiddleware_WithDisableOutput(t *testing.T) {
	exp, tp := setupOTel()

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
	exp, tp := setupOTel()

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

func TestMiddleware_WithOutputLimit(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "1234567890")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithOutputLimit(5)))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	var got string
	found := false
	for _, attr := range spans[0].Attributes {
		if string(attr.Key) == "goyek.task.output" {
			got = attr.Value.AsString()
			found = true
			break
		}
	}

	if !found {
		t.Error("goyek.task.output attribute not found")
	}
	if got != "12345" {
		t.Errorf("expected truncated output '12345', got %q", got)
	}
}

func TestExecutorMiddleware_WithOutputLimit(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "1234567890")
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithOutputLimit(3)))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
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

	var got string
	found := false
	for _, attr := range executeSpan.Attributes {
		if string(attr.Key) == "goyek.flow.output" {
			got = attr.Value.AsString()
			found = true
			break
		}
	}

	if !found {
		t.Error("goyek.flow.output attribute not found")
	}
	if got != "123" {
		t.Errorf("expected truncated output '123', got %q", got)
	}
}

func TestMiddleware_DisableOutput_Panic(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "panic",
		Action: func(a *goyek.A) {
			panic("sensitive info")
		},
	})
	f.Use(otelgoyek.Middleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithDisableOutput(true),
	))

	_ = f.Execute(context.Background(), []string{"panic"})

	spans := exp.GetSpans()
	if len(spans) == 0 {
		t.Fatal("no spans recorded")
	}

	for _, span := range spans {
		if span.Name == "panic" {
			for _, attr := range span.Attributes {
				if string(attr.Key) == "goyek.task.panic.value" {
					t.Errorf("panic value recorded even though output is disabled: %v", attr.Value.AsString())
				}
				if string(attr.Key) == "goyek.task.panic.stack" {
					t.Errorf("panic stack recorded even though output is disabled")
				}
				if string(attr.Key) == "goyek.task.output" {
					t.Errorf("output recorded even though output is disabled")
				}
			}
		}
	}
}

func setupOTel() (*tracetest.InMemoryExporter, *trace.TracerProvider) {
	exp := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exp)),
	)
	return exp, tp
}
