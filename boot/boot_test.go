package boot

import (
	"flag"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantTasks []string
		wantArgs  []string
		wantV     bool
		wantMsg   string
		wantErr   string
	}{
		{
			name: "no arguments",
		},
		{
			name:      "tasks only",
			args:      []string{"all", "diff"},
			wantTasks: []string{"all", "diff"},
		},
		{
			name:      "tasks before flag",
			args:      []string{"all", "diff", "-v"},
			wantTasks: []string{"all", "diff"},
			wantV:     true,
		},
		{
			name:     "task after flag",
			args:     []string{"-v", "all"},
			wantArgs: []string{"all"},
			wantV:    true,
			wantErr:  "unexpected arguments: [all]",
		},
		{
			name:     "multiple tasks after flag",
			args:     []string{"-v", "all", "diff"},
			wantArgs: []string{"all", "diff"},
			wantV:    true,
			wantErr:  "unexpected arguments: [all diff]",
		},
		{
			name:      "arguments after separator",
			args:      []string{"all", "-v", "--", "first", "-second"},
			wantTasks: []string{"all"},
			wantArgs:  []string{"first", "-second"},
			wantV:     true,
		},
		{
			name:      "separator without arguments",
			args:      []string{"all", "--"},
			wantTasks: []string{"all"},
			wantArgs:  []string{},
		},
		{
			name:     "separator before tasks",
			args:     []string{"--", "all", "diff"},
			wantArgs: []string{"all", "diff"},
		},
		{
			name:     "task between flag and separator",
			args:     []string{"-v", "all", "--", "first"},
			wantArgs: []string{"all"},
			wantV:    true,
			wantErr:  "unexpected arguments: [all]",
		},
		{
			name:      "separate flag value before separator",
			args:      []string{"all", "-msg", "hello", "--", "first"},
			wantTasks: []string{"all"},
			wantArgs:  []string{"first"},
			wantMsg:   "hello",
		},
		{
			name:    "separate separator cannot be a flag value",
			args:    []string{"-msg", "--", "all"},
			wantErr: "flag needs an argument: -msg",
		},
		{
			name:    "separator as inline flag value",
			args:    []string{"-msg=--"},
			wantMsg: "--",
		},
		{
			name:     "inline separator value and arguments",
			args:     []string{"-msg=--", "--", "all"},
			wantArgs: []string{"all"},
			wantMsg:  "--",
		},
		{
			name:    "unknown flag",
			args:    []string{"-unknown"},
			wantErr: "flag provided but not defined: -unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := flag.NewFlagSet("build", flag.ContinueOnError)
			flags.SetOutput(io.Discard)
			verbose := flags.Bool("v", false, "verbose output")
			message := flags.String("msg", "", "message")

			gotTasks, err := parseArgs(flags, tt.args)
			if !slices.Equal(gotTasks, tt.wantTasks) {
				t.Errorf("parseArgs(%v) tasks = %v, want %v", tt.args, gotTasks, tt.wantTasks)
			}
			if !slices.Equal(flags.Args(), tt.wantArgs) {
				t.Errorf("parseArgs(%v) flag args = %v, want %v", tt.args, flags.Args(), tt.wantArgs)
			}
			if *verbose != tt.wantV {
				t.Errorf("parseArgs(%v) -v = %t, want %t", tt.args, *verbose, tt.wantV)
			}
			if *message != tt.wantMsg {
				t.Errorf("parseArgs(%v) -msg = %q, want %q", tt.args, *message, tt.wantMsg)
			}
			if gotErr := errorString(err); gotErr != tt.wantErr {
				t.Errorf("parseArgs(%v) error = %q, want %q", tt.args, gotErr, tt.wantErr)
			}
		})
	}
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

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
