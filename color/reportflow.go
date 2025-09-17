package color

import (
	"time"

	"github.com/fatih/color"
	"github.com/goyek/goyek/v2"
)

// ReportFlow is a middleware which reports the flow execution status with colors.
//
// The format is based on the reports provided by the Go test runner.
func ReportFlow(next goyek.Executor) goyek.Executor {
	return func(in goyek.ExecuteInput) error {
		out := in.Output
		c := color.New(color.Bold)

		from := time.Now()
		if err := next(in); err != nil {
			c = c.Add(color.FgRed)
			c.Fprintf(out, "%v\t%.3fs\n", err, time.Since(from).Seconds())
			return err
		}

		c = c.Add(color.FgGreen)
		c.Fprintf(out, "ok\t%.3fs\n", time.Since(from).Seconds())
		return nil
	}
}
