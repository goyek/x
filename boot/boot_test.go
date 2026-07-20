package boot

import (
	"flag"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestUsageUsesCurrentGoyekOutput(t *testing.T) {
	originalCommandLine := flag.CommandLine
	originalOutput := goyek.Output()
	t.Cleanup(func() {
		flag.CommandLine = originalCommandLine
		goyek.SetOutput(originalOutput)
	})

	staleOutput := &strings.Builder{}
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	flag.CommandLine.SetOutput(staleOutput)
	flag.CommandLine.Bool("sample", false, "sample flag")

	currentOutput := &strings.Builder{}
	goyek.SetOutput(currentOutput)

	usage()

	if got := staleOutput.String(); got != "" {
		t.Fatalf("stale flag output received usage text: %q", got)
	}
	got := currentOutput.String()
	for _, want := range []string{
		"Usage of build: [tasks] [flags] [--] [args]\n",
		"Tasks:\n",
		"Flags:\n",
		"  -sample\n",
		"    \tsample flag\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("usage output does not contain %q:\n%s", want, got)
		}
	}
}
