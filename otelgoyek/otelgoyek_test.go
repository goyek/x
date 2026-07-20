package otelgoyek_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
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
	attrTaskOutput            = "goyek.task.output"
	traceparent               = "00-0102030405060708090a0b0c0d0e0f10-0102030405060708-01"
	spanNameExecute           = "Execute"
	taskNameTest              = "test"
	taskNamePanic             = "panic"
	concurrentWriters         = 16
	concurrentWritesPerWriter = 128
)

func TestMiddleware_WithDisableOutput(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: taskNameTest,
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "secret message")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithDisableOutput(true)))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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

func TestMiddleware_PanicStatus(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: taskNamePanic,
		Action: func(_ *goyek.A) {
			panic("something went wrong")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{taskNamePanic})

	spans := exp.GetSpans()
	if len(spans) == 0 {
		t.Fatal("no spans recorded")
	}

	for _, span := range spans {
		if span.Name == taskNamePanic {
			if span.Status.Code != codes.Error {
				t.Errorf("expected span status Error for panicking task, got %v", span.Status.Code)
			}
		}
	}
}

func TestMiddleware_PanicTruncation(t *testing.T) {
	exp, tp := setupOTel()

	limit := 10
	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: taskNamePanic,
		Action: func(_ *goyek.A) {
			panic("1234567890ABCDE")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithOutputLimit(limit)))

	_ = f.Execute(context.Background(), []string{taskNamePanic})

	spans := exp.GetSpans()
	var taskSpan *tracetest.SpanStub
	for _, s := range spans {
		if s.Name == taskNamePanic {
			taskSpan = &s
			break
		}
	}

	if taskSpan == nil {
		t.Fatal("task span not found")
	}

	for _, attr := range taskSpan.Attributes {
		if string(attr.Key) == "goyek.task.panic.value" {
			got := attr.Value.AsString()
			if len(got) > limit {
				t.Errorf("panic value not truncated: got length %d, limit %d", len(got), limit)
			}
			if got != "1234567890" {
				t.Errorf("expected truncated panic value '1234567890', got %q", got)
			}
		}
		if string(attr.Key) == "goyek.task.panic.stack" {
			got := attr.Value.AsString()
			if len(got) > limit {
				t.Errorf("panic stack not truncated: got length %d, limit %d", len(got), limit)
			}
		}
	}
}

func TestExecutorMiddleware_ErrorTruncation(t *testing.T) {
	exp, tp := setupOTel()

	limit := 5
	mw := otelgoyek.ExecutorMiddleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithOutputLimit(limit),
	)

	next := func(_ goyek.ExecuteInput) error {
		return errors.New("1234567890")
	}

	executor := mw(next)

	_ = executor(goyek.ExecuteInput{
		Context: context.Background(),
		Tasks:   []string{taskNameTest},
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

	got := executeSpan.Status.Description
	if len(got) > limit {
		t.Errorf("error status description not truncated: got length %d, limit %d", len(got), limit)
	}
	if got != "12345" {
		t.Errorf("expected truncated error message '12345', got %q", got)
	}
}

func TestExecutorMiddleware_WithDisableOutput(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: taskNameTest,
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "secret flow message")
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithDisableOutput(true)))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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
		Name: taskNameTest,
		Action: func(_ *goyek.A) {
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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
		Name: taskNameTest,
		Action: func(_ *goyek.A) {
		},
	})
	f.Use(otelgoyek.Middleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithPropagator(propagation.TraceContext{}),
	))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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
		Name: taskNameTest,
		Action: func(_ *goyek.A) {
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp)))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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
		Name: taskNameTest,
		Action: func(_ *goyek.A) {
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithPropagator(propagation.TraceContext{}),
	))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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
		Name: taskNameTest,
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "1234567890")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithOutputLimit(5)))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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

func TestMiddleware_CapturesConcurrentOutput(t *testing.T) {
	exp, tp := setupOTel()

	var writeErr error
	output := &strings.Builder{}
	f := &goyek.Flow{}
	f.SetOutput(goyek.SyncWriter(output))
	f.Define(goyek.Task{
		Name: taskNameTest,
		Action: func(a *goyek.A) {
			writeErr = writeConcurrentRecords(a.Output(), "task")
		},
	})
	f.Use(otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp)))

	if err := f.Execute(context.Background(), []string{taskNameTest}); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if writeErr != nil {
		t.Fatalf("writing task output: %v", writeErr)
	}

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	captured := attributeValue(t, spans[0], attrTaskOutput)
	if destination := output.String(); captured != destination {
		t.Fatalf("captured output does not match destination output\ncaptured:    %q\ndestination: %q", captured, destination)
	}
	assertConcurrentRecords(t, captured, "task")
}

func TestMiddleware_CapturesOutputWithNilDestination(t *testing.T) {
	exp, tp := setupOTel()
	runner := otelgoyek.Middleware(otelgoyek.WithTracerProvider(tp))(goyek.NewRunner(func(a *goyek.A) {
		fmt.Fprint(a.Output(), "task output")
	}))

	result := runner(goyek.Input{Context: context.Background(), TaskName: taskNameTest})
	if result.Status != goyek.StatusPassed {
		t.Fatalf("runner() status = %s, want PASS", result.Status)
	}

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if got := attributeValue(t, spans[0], attrTaskOutput); got != "task output" {
		t.Errorf("captured output = %q, want %q", got, "task output")
	}
}

func TestExecutorMiddleware_WithOutputLimit(t *testing.T) {
	exp, tp := setupOTel()

	f := &goyek.Flow{}
	f.Define(goyek.Task{
		Name: taskNameTest,
		Action: func(a *goyek.A) {
			fmt.Fprint(a.Output(), "1234567890")
		},
	})
	f.UseExecutor(otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp), otelgoyek.WithOutputLimit(3)))

	_ = f.Execute(context.Background(), []string{taskNameTest})

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

func TestExecutorMiddleware_CapturesConcurrentOutput(t *testing.T) {
	exp, tp := setupOTel()
	output := &strings.Builder{}
	mw := otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp))
	executor := mw(func(in goyek.ExecuteInput) error {
		return writeConcurrentRecords(in.Output, "flow")
	})

	if err := executor(goyek.ExecuteInput{
		Context: context.Background(),
		Tasks:   []string{taskNameTest},
		Output:  goyek.SyncWriter(output),
	}); err != nil {
		t.Fatalf("executor() error = %v", err)
	}

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	captured := attributeValue(t, spans[0], "goyek.flow.output")
	if destination := output.String(); captured != destination {
		t.Fatalf("captured output does not match destination output\ncaptured:    %q\ndestination: %q", captured, destination)
	}
	assertConcurrentRecords(t, captured, "flow")
}

func TestExecutorMiddleware_CapturesOutputWithNilDestination(t *testing.T) {
	exp, tp := setupOTel()
	executor := otelgoyek.ExecutorMiddleware(otelgoyek.WithTracerProvider(tp))(func(in goyek.ExecuteInput) error {
		_, err := io.WriteString(in.Output, "flow output")
		return err
	})

	if err := executor(goyek.ExecuteInput{Context: context.Background(), Tasks: []string{taskNameTest}}); err != nil {
		t.Fatalf("executor() error = %v", err)
	}

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if got := attributeValue(t, spans[0], "goyek.flow.output"); got != "flow output" {
		t.Errorf("captured output = %q, want %q", got, "flow output")
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
		Tasks:   []string{taskNameTest},
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
		Name: taskNamePanic,
		Action: func(_ *goyek.A) {
			panic("sensitive info")
		},
	})
	f.Use(otelgoyek.Middleware(
		otelgoyek.WithTracerProvider(tp),
		otelgoyek.WithDisableOutput(true),
	))

	_ = f.Execute(context.Background(), []string{taskNamePanic})

	spans := exp.GetSpans()
	if len(spans) == 0 {
		t.Fatal("no spans recorded")
	}

	for _, span := range spans {
		if span.Name == taskNamePanic {
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

func writeConcurrentRecords(w io.Writer, prefix string) error {
	start := make(chan struct{})
	var ready sync.WaitGroup
	var writers sync.WaitGroup
	errCh := make(chan error, concurrentWriters)
	ready.Add(concurrentWriters)
	writers.Add(concurrentWriters)

	for writerID := range concurrentWriters {
		record := fmt.Sprintf("%s-%02d\n", prefix, writerID)
		go func() {
			defer writers.Done()
			ready.Done()
			<-start
			for range concurrentWritesPerWriter {
				n, err := io.WriteString(w, record)
				if err != nil {
					errCh <- err
					return
				}
				if n != len(record) {
					errCh <- io.ErrShortWrite
					return
				}
			}
		}()
	}

	ready.Wait()
	close(start)
	writers.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func assertConcurrentRecords(t *testing.T, output, prefix string) {
	t.Helper()

	counts := make(map[string]int, concurrentWriters)
	for record := range strings.Lines(output) {
		counts[record]++
	}
	if got, want := len(counts), concurrentWriters; got != want {
		t.Errorf("unique record count = %d, want %d", got, want)
	}
	for writerID := range concurrentWriters {
		record := fmt.Sprintf("%s-%02d\n", prefix, writerID)
		if got, want := counts[record], concurrentWritesPerWriter; got != want {
			t.Errorf("record %q count = %d, want %d", record, got, want)
		}
	}
}

func attributeValue(t *testing.T, span tracetest.SpanStub, key string) string {
	t.Helper()
	for _, attr := range span.Attributes {
		if string(attr.Key) == key {
			return attr.Value.AsString()
		}
	}
	t.Fatalf("%s attribute not found", key)
	return ""
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
