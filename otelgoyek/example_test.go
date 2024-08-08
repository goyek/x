package otelgoyek_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/goyek/goyek/v2"
	"github.com/goyek/goyek/v2/middleware"
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

	// Add reporting middlewares.
	goyek.UseExecutor(middleware.ReportFlow)
	goyek.Use(middleware.ReportStatus)

	// Add OpenTelemetry instrumentation to flow execution.
	goyek.UseExecutor(otelgoyek.ExecutorMiddleware())

	// Run the tasks.
	goyek.Use(middleware.BufferParallel)
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
		                "TraceID": "e6780df143913b35644a376fd3543ab3",
		                "SpanID": "afc06d19087b33da",
		                "TraceFlags": "01",
		                "TraceState": "",
		                "Remote": false
		        },
		        "Parent": {
		                "TraceID": "e6780df143913b35644a376fd3543ab3",
		                "SpanID": "cade09c93a2ba275",
		                "TraceFlags": "01",
		                "TraceState": "",
		                "Remote": false
		        },
		        "SpanKind": 1,
		        "StartTime": "2024-08-08T23:19:03.304436137+02:00",
		        "EndTime": "2024-08-08T23:19:03.304525337+02:00",
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
		                                "Value": "unknown_service:main"
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
		                "TraceID": "e6780df143913b35644a376fd3543ab3",
		                "SpanID": "cade09c93a2ba275",
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
		        "StartTime": "2024-08-08T23:19:03.304381737+02:00",
		        "EndTime": "2024-08-08T23:19:03.304575637+02:00",
		        "Attributes": [
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
		                                "Value": "unknown_service:main"
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
