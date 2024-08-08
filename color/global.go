// Package color contains goyek features which additionally
// have colors.
//
// Set NO_COLOR environment variable to a non-empty string
// or use the NoColor function to prevent colorizing the output.
package color

import (
	"os"

	"github.com/fatih/color"
)

func init() {
	color.NoColor = os.Getenv("NO_COLOR") != ""
}

// NoColor prevents colorizing the output.
func NoColor() {
	color.NoColor = true
}
