package otelgoyek

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Option configures the instrumentation.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (fn optionFunc) apply(cfg *config) {
	fn(cfg)
}

type config struct {
	TracerProvider trace.TracerProvider
	Propagator     propagation.TextMapPropagator
	DisableOutput  bool
	OutputLimit    int
}

func newConfig(opts []Option) *config {
	c := &config{
		TracerProvider: otel.GetTracerProvider(),
		Propagator:     otel.GetTextMapPropagator(),
		DisableOutput:  false,
		OutputLimit:    1024 * 1024, //nolint:mnd // 1 MiB
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	})
}

// WithPropagator specifies a text map propagator to extract context from the environment variables.
// If none is specified, the global propagator is used.
func WithPropagator(propagator propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		if propagator != nil {
			cfg.Propagator = propagator
		}
	})
}

// WithDisableOutput specifies if the output should not be captured in the span attributes.
// This is useful for security reasons, to avoid sensitive data exposure.
func WithDisableOutput(disable bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.DisableOutput = disable
	})
}

// WithOutputLimit specifies the maximum number of bytes of output to capture in the span attributes.
// The default is 1 MiB.
func WithOutputLimit(limit int) Option {
	return optionFunc(func(cfg *config) {
		cfg.OutputLimit = limit
	})
}
