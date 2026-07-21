package color_test

import (
	"testing"

	fatihcolor "github.com/fatih/color"
)

const (
	ansiBlue   = "\x1b[34m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiReset  = "\x1b[0m"
	ansiYellow = "\x1b[33m"
)

func forceColor(t *testing.T) {
	t.Helper()
	t.Setenv("NO_COLOR", "")
	oldNoColor := fatihcolor.NoColor
	fatihcolor.NoColor = false
	t.Cleanup(func() {
		fatihcolor.NoColor = oldNoColor
	})
}
