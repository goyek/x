package otelgoyek_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/goyek/goyek/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/goyek/x/otelgoyek"
)

const (
	attrTaskOutput  = "goyek.task.output"
	traceparent     = "00-0102030405060708090a0b0c0d0e0f10-0102030405060708-01"
	spanNameExecute = "Execute"
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
		if string(attr.Key) == attrTaskOutput {
			t.Errorf("found %s attribute even though output capture is disabled: %v", attrTaskOutput, attr.Value.AsString())
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
		if s.Name == spanNameExecute {
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

func TestMiddleware_ExtractsTraceContextFromEnvironment(t *testing.T) {
	useTraceContextPropagator(t)
	t.Setenv("TRACEPARENT", traceparent)
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(_ *goyek.A) {
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	assertEnvParent(t, spans[0])
}

func TestMiddleware_WithPropagator(t *testing.T) {
	usePropagator(t, propagation.NewCompositeTextMapPropagator())
	t.Setenv("TRACEPARENT", traceparent)
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(_ *goyek.A) {
		},
	})
	f.Use(otelgoyek.Middleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithPropagator(propagation.TraceContext{}),
	))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	assertEnvParent(t, spans[0])
}

func TestExecutorMiddleware_ExtractsTraceContextFromEnvironment(t *testing.T) {
	useTraceContextPropagator(t)
	t.Setenv("TRACEPARENT", traceparent)
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(_ *goyek.A) {
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	assertEnvParent(t, spans[0])
}

func TestExecutorMiddleware_WithPropagator(t *testing.T) {
	usePropagator(t, propagation.NewCompositeTextMapPropagator())
	t.Setenv("TRACEPARENT", traceparent)
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "test",
		Action: func(_ *goyek.A) {
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithPropagator(propagation.TraceContext{}),
	))

	_ = f.Execute(context.Background(), []string{"test"})

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	assertEnvParent(t, spans[0])
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
		if string(attr.Key) == attrTaskOutput {
			got = attr.Value.AsString()
			found = true
			break
		}
	}

	if !found {
		t.Errorf("%s attribute not found", attrTaskOutput)
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
		if s.Name == spanNameExecute {
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

func TestExecutorMiddleware_WithDisableOutput_StatusLeak(t *testing.T) {
	exp, tp := setupOTel()

	// We need a custom executor that returns an error with sensitive info
	// because goyek.Flow.Execute returns an error if a task fails.
	mw := otelgoyek.ExecutorMiddleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithDisableOutput(true),
	)

	next := func(_ goyek.ExecuteInput) error {
		return errors.New("sensitive error message")
	}

	executor := mw(next)

	_ = executor(goyek.ExecuteInput{
		Context: context.Background(),
		Tasks:   []string{"test"},
	})

	spans := exp.GetSpans()
	var executeSpan *tracetest.SpanStub
	for _, s := range spans {
		if s.Name == spanNameExecute {
			executeSpan = &s
			break
		}
	}

	if executeSpan == nil {
		t.Fatal("Execute span not found")
	}

	if executeSpan.Status.Code != codes.Error {
		t.Errorf("expected span status Error, got %v", executeSpan.Status.Code)
	}

	if executeSpan.Status.Description == "sensitive error message" {
		t.Errorf("found sensitive error message in span status even though output capture is disabled")
	}
}

func TestMiddleware_DisableOutput_Panic(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "panic",
		Action: func(_ *goyek.A) {
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
				if string(attr.Key) == attrTaskOutput {
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

func useTraceContextPropagator(t *testing.T) {
	t.Helper()
	usePropagator(t, propagation.TraceContext{})
}

func usePropagator(t *testing.T, propagator propagation.TextMapPropagator) {
	t.Helper()
	previous := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagator)
	t.Cleanup(func() {
		otel.SetTextMapPropagator(previous)
	})
}

func assertEnvParent(t *testing.T, span tracetest.SpanStub) {
	t.Helper()
	if got, want := span.Parent.TraceID().String(), "0102030405060708090a0b0c0d0e0f10"; got != want {
		t.Errorf("parent trace ID = %q, want %q", got, want)
	}
	if got, want := span.Parent.SpanID().String(), "0102030405060708"; got != want {
		t.Errorf("parent span ID = %q, want %q", got, want)
	}
	if !span.Parent.IsRemote() {
		t.Error("parent span context is not remote")
	}
}

func TestMiddleware_OutputRace(_ *testing.T) {
	_, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: "race",
		Action: func(a *goyek.A) {
			var wg sync.WaitGroup
			for i := 0; i < 2; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < 1000; j++ {
						fmt.Fprint(a.Output(), "a")
					}
				}()
			}
			wg.Wait()
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{"race"})
}
