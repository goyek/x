// Package boot contains an extension of goyek.Main which additionally
// defines flags and configures the flow in a convenient way.
package boot

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/goyek/v3/middleware"

	"github.com/goyek/x/color"
)

const exitCodeInvalid = 2

// Reusable flags used by the build pipeline.
var (
	v       = flag.Bool("v", false, "print all tasks as they are run")
	dryRun  = flag.Bool("dry-run", false, "print all tasks without executing actions")
	longRun = flag.Duration("long-run", time.Minute, "print when a task takes longer")
	noDeps  = flag.Bool("no-deps", false, "do not process dependencies")
	skip    = flag.String("skip", "", "skip processing the `comma-separated tasks`")
	noColor = flag.Bool("no-color", false, "disable colorizing output")
)

// Main is an extension of goyek.Main which additionally
// defines flags and uses the most useful middlewares.
func Main() {
	tasks, args := goyek.SplitTasks(os.Args[1:])
	flag.CommandLine.SetOutput(goyek.Output())
	flag.Usage = usage
	if err := flag.CommandLine.Parse(args); err != nil {
		fmt.Fprintln(goyek.Output(), err)
		os.Exit(exitCodeInvalid)
	}

	if *dryRun {
		*v = true // needed to report the task status
	}

	goyek.UseExecutor(color.ReportFlow)

	if *dryRun {
		goyek.Use(middleware.DryRun)
	}
	goyek.Use(color.ReportStatus)
	if *v {
		goyek.Use(middleware.BufferParallel)
	} else {
		goyek.Use(middleware.SilentNonFailed)
	}
	if *longRun > 0 {
		goyek.Use(middleware.ReportLongRun(*longRun))
	}
	if *noColor {
		color.NoColor()
	}

	var opts []goyek.Option
	if *noDeps {
		opts = append(opts, goyek.NoDeps())
	}
	if *skip != "" {
		skippedTasks := strings.Split(*skip, ",")
		opts = append(opts, goyek.Skip(skippedTasks...))
	}

	goyek.SetUsage(usage)
	goyek.SetLogger(&color.CodeLineLogger{})
	goyek.Main(tasks, opts...)
}

func usage() {
	fmt.Println("Usage of build: [tasks] [flags] [--] [args]")
	goyek.Print()
	fmt.Println("Flags:")
	flag.PrintDefaults()
}
