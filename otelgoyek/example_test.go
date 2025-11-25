package otelgoyek_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/goyek/v3/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/goyek/x/otelgoyek"
)

func run(ctx context.Context, w io.Writer, tasks []string) (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	go func() {
		<-ctx.Done()
		fmt.Fprintln(w, "first interrupt, graceful stop")
		stop()
	}()

	// Setup OpenTelemetry tracing pipeline.
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return err
	}
	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(traceExporter))
	defer func() {
		err = errors.Join(err, tracerProvider.Shutdown(context.Background()))
	}()
	otel.SetTracerProvider(tracerProvider)

	// Define a task printing a message and set it as the default task.
	hi := goyek.Define(goyek.Task{
		Name:  "hi",
		Usage: "Greetings",
		Action: func(a *goyek.A) {
			a.Log("Hello world!")
		},
	})
	goyek.SetDefault(hi)

	// Add OpenTelemetry instrumentation to task run.
	goyek.Use(otelgoyek.Middleware())

	// Add common middlewares.
	goyek.UseExecutor(middleware.ReportFlow)
	goyek.Use(middleware.ReportStatus)
	goyek.Use(middleware.BufferParallel)

	// Add OpenTelemetry instrumentation to flow execution.
	goyek.UseExecutor(otelgoyek.ExecutorMiddleware())

	// Run the tasks.
	goyek.SetOutput(w)
	return goyek.Execute(ctx, tasks)
}

func Example() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args[1:]); err != nil {
		os.Exit(1)
	}

	/*
		$ go run .
		===== TASK  hi
		      main.go:45: Hello world!
		----- PASS: hi (0.00s)
		ok      0.000s
		{
		        "Name": "hi",
		        "SpanContext": {
		                "TraceID": "c6e535658cb43357c4188eb87cbcd844",
		                "SpanID": "33afad8f183e2955",
		                "TraceFlags": "01",
		                "TraceState": "",
		                "Remote": false
		        },
		        "Parent": {
		                "TraceID": "c6e535658cb43357c4188eb87cbcd844",
		                "SpanID": "6ea000e83a2029a9",
		                "TraceFlags": "01",
		                "TraceState": "",
		                "Remote": false
		        },
		        "SpanKind": 1,
		        "StartTime": "2024-08-09T09:02:58.883480983+02:00",
		        "EndTime": "2024-08-09T09:02:58.883500285+02:00",
		        "Attributes": [
		                {
		                        "Key": "goyek.task.name",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "hi"
		                        }
		                },
		                {
		                        "Key": "goyek.task.output",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "      main.go:45: Hello world!\n"
		                        }
		                },
		                {
		                        "Key": "goyek.task.status",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "PASS"
		                        }
		                }
		        ],
		        "Events": null,
		        "Links": null,
		        "Status": {
		                "Code": "Unset",
		                "Description": ""
		        },
		        "DroppedAttributes": 0,
		        "DroppedEvents": 0,
		        "DroppedLinks": 0,
		        "ChildSpanCount": 0,
		        "Resource": [
		                {
		                        "Key": "service.name",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "unknown_service:x"
		                        }
		                },
		                {
		                        "Key": "telemetry.sdk.language",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "go"
		                        }
		                },
		                {
		                        "Key": "telemetry.sdk.name",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "opentelemetry"
		                        }
		                },
		                {
		                        "Key": "telemetry.sdk.version",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "1.28.0"
		                        }
		                }
		        ],
		        "InstrumentationLibrary": {
		                "Name": "github.com/goyek/x/otelgoyek",
		                "Version": "0.2.0",
		                "SchemaURL": ""
		        }
		}
		{
		        "Name": "Execute",
		        "SpanContext": {
		                "TraceID": "c6e535658cb43357c4188eb87cbcd844",
		                "SpanID": "6ea000e83a2029a9",
		                "TraceFlags": "01",
		                "TraceState": "",
		                "Remote": false
		        },
		        "Parent": {
		                "TraceID": "00000000000000000000000000000000",
		                "SpanID": "0000000000000000",
		                "TraceFlags": "00",
		                "TraceState": "",
		                "Remote": false
		        },
		        "SpanKind": 1,
		        "StartTime": "2024-08-09T09:02:58.88343108+02:00",
		        "EndTime": "2024-08-09T09:02:58.883508185+02:00",
		        "Attributes": [
		                {
		                        "Key": "goyek.flow.tasks",
		                        "Value": {
		                                "Type": "STRINGSLICE",
		                                "Value": []
		                        }
		                },
		                {
		                        "Key": "goyek.flow.skip_tasks",
		                        "Value": {
		                                "Type": "STRINGSLICE",
		                                "Value": []
		                        }
		                },
		                {
		                        "Key": "goyek.flow.no_deps",
		                        "Value": {
		                                "Type": "BOOL",
		                                "Value": false
		                        }
		                },
		                {
		                        "Key": "goyek.flow.output",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "===== TASK  hi\n      main.go:45: Hello world!\n----- PASS: hi (0.00s)\nok\t0.000s\n"
		                        }
		                }
		        ],
		        "Events": null,
		        "Links": null,
		        "Status": {
		                "Code": "Unset",
		                "Description": ""
		        },
		        "DroppedAttributes": 0,
		        "DroppedEvents": 0,
		        "DroppedLinks": 0,
		        "ChildSpanCount": 1,
		        "Resource": [
		                {
		                        "Key": "service.name",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "unknown_service:x"
		                        }
		                },
		                {
		                        "Key": "telemetry.sdk.language",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "go"
		                        }
		                },
		                {
		                        "Key": "telemetry.sdk.name",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "opentelemetry"
		                        }
		                },
		                {
		                        "Key": "telemetry.sdk.version",
		                        "Value": {
		                                "Type": "STRING",
		                                "Value": "1.28.0"
		                        }
		                }
		        ],
		        "InstrumentationLibrary": {
		                "Name": "github.com/goyek/x/otelgoyek",
		                "Version": "0.2.0",
		                "SchemaURL": ""
		        }
		}
	*/
}
