package color

import (
	"time"

	"github.com/fatih/color"
	"github.com/goyek/goyek/v3"
)

// ReportStatus is a middleware which reports the task run status with colors.
//
// The format is based on the reports provided by the Go test runner.
func ReportStatus(next goyek.Runner) goyek.Runner {
	return func(in goyek.Input) goyek.Result {
		c := color.New(color.FgBlue)

		// report start task
		writeString(in.Output, c.Sprintf("===== TASK  %s\n", in.TaskName))
		start := time.Now()

		// run
		res := next(in)

		// report task end
		c = color.New(color.FgGreen)
		status := "PASS"
		switch res.Status {
		case goyek.StatusFailed:
			c = color.New(color.FgRed)
			status = "FAIL"
		case goyek.StatusSkipped:
			c = color.New(color.FgYellow)
			status = "SKIP"
		case goyek.StatusNotRun:
			status = "NOOP"
		}
		writeString(in.Output, c.Sprintf("----- %s: %s (%.2fs)\n", status, in.TaskName, time.Since(start).Seconds()))

		// report panic if happened
		if res.PanicStack != nil {
			var panicHeader string
			if res.PanicValue != nil {
				panicHeader = c.Sprintf("panic: %v", res.PanicValue)
			} else {
				panicHeader = c.Sprint("panic(nil) or runtime.Goexit() called")
			}
			panicStack := c.Sprintf("%s", res.PanicStack)
			writeString(in.Output, panicHeader+"\n\n"+panicStack)
		}

		return res
	}
}
