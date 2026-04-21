package otelgoyek

import (
	"go.opentelemetry.io/otel"
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
	DisableOutput  bool
}

func newConfig(opts []Option) *config {
	c := &config{
		TracerProvider: otel.GetTracerProvider(),
		DisableOutput:  false,
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

// WithDisableOutput specifies if the output should not be captured in the span attributes.
// This is useful for security reasons, to avoid sensitive data exposure.
func WithDisableOutput(disable bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.DisableOutput = disable
	})
}
