// Package otelgoyek provides OpenTelemetry instrumentation for goyek.
package otelgoyek

import (
	"fmt"
	"io"
	"strings"

	"github.com/goyek/goyek/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName    = "github.com/goyek/x/otelgoyek"
	instrumentationVersion = "0.2.0"
)

// Middleware returns a [goyek.Middleware] which adds OpenTelemetry tracing instrumentation to task run.
func Middleware(opts ...Option) goyek.Middleware {
	cfg := newConfig(opts)
	tracer := cfg.TracerProvider.Tracer(instrumentationName, trace.WithInstrumentationVersion(instrumentationVersion))
	r := runner{tracer}
	return r.Middleware
}

// ExecutorMiddleware returns a [goyek.ExecutorMiddleware] which adds OpenTelemetry tracing instrumentation to flow execution.
func ExecutorMiddleware(opts ...Option) goyek.ExecutorMiddleware {
	cfg := newConfig(opts)
	tracer := cfg.TracerProvider.Tracer(instrumentationName, trace.WithInstrumentationVersion(instrumentationVersion))
	e := executor{tracer}
	return e.Middleware
}

type executor struct {
	tracer trace.Tracer
}

func (e *executor) Middleware(next goyek.Executor) goyek.Executor {
	return func(in goyek.ExecuteInput) error {
		ctx, span := e.tracer.Start(in.Context, "Execute")
		defer span.End()

		in.Context = ctx

		sb := &strings.Builder{}
		in.Output = io.MultiWriter(in.Output, sb)

		err := next(in)

		span.SetAttributes(attribute.String("goyek.flow.output", sb.String()))
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		return err
	}
}

type runner struct {
	tracer trace.Tracer
}

func (r *runner) Middleware(next goyek.Runner) goyek.Runner {
	return func(in goyek.Input) goyek.Result {
		ctx, span := r.tracer.Start(in.Context, in.TaskName,
			trace.WithAttributes(attribute.String("goyek.task.name", in.TaskName)))
		defer span.End()

		in.Context = ctx

		sb := &strings.Builder{}
		in.Output = io.MultiWriter(in.Output, sb)

		res := next(in)

		span.SetAttributes(attribute.String("goyek.task.output", sb.String()))

		span.SetAttributes(attribute.String("goyek.task.status", res.Status.String()))
		if res.Status == goyek.StatusFailed {
			span.SetStatus(codes.Error, "task failed: "+in.TaskName)
		}

		if res.PanicStack != nil {
			if res.PanicValue != nil {
				span.SetAttributes(attribute.String("goyek.task.panic.value", fmt.Sprint(res.PanicValue)))
			} else {
				span.SetAttributes(attribute.String("goyek.task.panic.value", "panic(nil) or runtime.Goexit() called"))
			}

			span.SetAttributes(attribute.String("goyek.task.panic.stack", fmt.Sprint(res.PanicStack)))
		}

		return res
	}
}
