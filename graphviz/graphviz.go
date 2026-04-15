// Package graphviz visualizes a dependency graph with registered tasks.
package graphviz

import (
	"errors"
	"fmt"
	"io"

	"github.com/goyek/goyek/v3"
)

// Option configures the Graphviz drawing.
type Option interface {
	apply(*config)
}

type config struct {
}

// Draw visualizes a dependency graph with registered tasks in DOT format.
func Draw(w io.Writer, flow *goyek.Flow, opts ...Option) error {
	if w == nil {
		return errors.New("nil writer")
	}
	if flow == nil {
		return errors.New("nil flow")
	}

	cfg := &config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	if _, err := io.WriteString(w, "digraph G {\n"); err != nil {
		return err
	}

	for _, task := range flow.Tasks() {
		if usage := task.Usage(); usage != "" {
			if _, err := fmt.Fprintf(w, "  %q [tooltip=%q];\n", task.Name(), usage); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "  %q;\n", task.Name()); err != nil {
				return err
			}
		}
		for _, dep := range task.Deps() {
			if _, err := fmt.Fprintf(w, "  %q -> %q;\n", task.Name(), dep.Name()); err != nil {
				return err
			}
		}
	}

	if _, err := io.WriteString(w, "}\n"); err != nil {
		return err
	}

	return nil
}
