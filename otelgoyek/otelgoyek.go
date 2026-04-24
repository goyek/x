// Package otelgoyek provides OpenTelemetry instrumentation for goyek.
//
// The instrumentation extracts context from environment variables using the
// global OpenTelemetry text map propagator and the envcar carrier before
// starting spans. Configure the global propagator, for example with
// propagation.TraceContext, to continue traces passed through the environment.
package otelgoyek

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goyek/goyek/v3"
	"go.opentelemetry.io/contrib/propagators/envcar"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName    = "github.com/goyek/x/otelgoyek"
	instrumentationVersion = "0.2.0"
)

// Middleware returns a [goyek.Middleware] which adds
// OpenTelemetry tracing instrumentation to task run.
func Middleware(opts ...Option) goyek.Middleware {
	cfg := newConfig(opts)
	tracer := cfg.TracerProvider.Tracer(instrumentationName, trace.WithInstrumentationVersion(instrumentationVersion))
	r := runner{
		tracer:        tracer,
		disableOutput: cfg.DisableOutput,
		outputLimit:   cfg.OutputLimit,
	}
	return r.Middleware
}

// ExecutorMiddleware returns a [goyek.ExecutorMiddleware] which adds
// OpenTelemetry tracing instrumentation to flow execution.
func ExecutorMiddleware(opts ...Option) goyek.ExecutorMiddleware {
	cfg := newConfig(opts)
	tracer := cfg.TracerProvider.Tracer(instrumentationName, trace.WithInstrumentationVersion(instrumentationVersion))
	e := executor{
		tracer:        tracer,
		disableOutput: cfg.DisableOutput,
		outputLimit:   cfg.OutputLimit,
	}
	return e.Middleware
}

type executor struct {
	tracer        trace.Tracer
	disableOutput bool
	outputLimit   int
}

func (e *executor) Middleware(next goyek.Executor) goyek.Executor {
	return func(in goyek.ExecuteInput) error {
		// ExecutorMiddleware creates the flow root span, so it also needs to
		// extract environment context before task spans inherit from it.
		ctx := extractContextFromEnv(in.Context)
		ctx, span := e.tracer.Start(ctx, "Execute", trace.WithAttributes(
			attribute.StringSlice("goyek.flow.tasks", in.Tasks),
			attribute.StringSlice("goyek.flow.skip_tasks", in.SkipTasks),
			attribute.Bool("goyek.flow.no_deps", in.NoDeps),
		))
		defer span.End()

		in.Context = ctx

		var sb *strings.Builder
		if !e.disableOutput {
			sb = &strings.Builder{}
			in.Output = io.MultiWriter(in.Output, &limitWriter{sb: sb, limit: e.outputLimit})
		}

		err := next(in)

		if !e.disableOutput {
			span.SetAttributes(attribute.String("goyek.flow.output", sb.String()))
		}
		if err != nil {
			msg := err.Error()
			if e.disableOutput {
				msg = "flow execution failed"
			}
			span.SetStatus(codes.Error, msg)
		}
		return err
	}
}

type runner struct {
	tracer        trace.Tracer
	disableOutput bool
	outputLimit   int
}

func (r *runner) Middleware(next goyek.Runner) goyek.Runner {
	return func(in goyek.Input) goyek.Result {
		// Middleware can be used without ExecutorMiddleware, so task spans need
		// to extract environment context before starting their own root span.
		ctx := extractContextFromEnv(in.Context)
		ctx, span := r.tracer.Start(ctx, in.TaskName, trace.WithAttributes(
			attribute.String("goyek.task.name", in.TaskName),
		))
		defer span.End()

		in.Context = ctx

		var sb *strings.Builder
		if !r.disableOutput {
			sb = &strings.Builder{}
			in.Output = io.MultiWriter(in.Output, &limitWriter{sb: sb, limit: r.outputLimit})
		}

		res := next(in)

		if !r.disableOutput {
			span.SetAttributes(attribute.String("goyek.task.output", sb.String()))
		}

		span.SetAttributes(attribute.String("goyek.task.status", res.Status.String()))
		if res.Status == goyek.StatusFailed {
			span.SetStatus(codes.Error, "task failed: "+in.TaskName)
		}

		if res.PanicStack != nil && !r.disableOutput {
			if res.PanicValue != nil {
				span.SetAttributes(attribute.String("goyek.task.panic.value", fmt.Sprint(res.PanicValue)))
			} else {
				span.SetAttributes(attribute.String("goyek.task.panic.value", "panic(nil) or runtime.Goexit() called"))
			}

			span.SetAttributes(attribute.String("goyek.task.panic.stack", string(res.PanicStack)))
		}

		return res
	}
}

type limitWriter struct {
	sb    *strings.Builder
	limit int
}

func (w *limitWriter) Write(p []byte) (int, error) {
	if w.limit <= 0 {
		return len(p), nil
	}
	available := w.limit - w.sb.Len()
	if available <= 0 {
		return len(p), nil
	}
	if len(p) > available {
		n, err := w.sb.Write(p[:available])
		if err != nil {
			return n, err
		}
		return len(p), nil
	}
	return w.sb.Write(p)
}

func extractContextFromEnv(ctx context.Context) context.Context {
	if trace.SpanContextFromContext(ctx).IsValid() {
		return ctx
	}
	return otel.GetTextMapPropagator().Extract(ctx, &envcar.Carrier{})
}
