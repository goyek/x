package graphviz_test

import (
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/graphviz"
)

func TestDraw(t *testing.T) {
	t.Run("nil writer", func(t *testing.T) {
		err := graphviz.Draw(nil, &goyek.Flow{})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("nil flow", func(t *testing.T) {
		err := graphviz.Draw(&strings.Builder{}, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("empty flow", func(t *testing.T) {
		f := &goyek.Flow{}
		sb := &strings.Builder{}
		err := graphviz.Draw(sb, f)
		if err != nil {
			t.Fatal(err)
		}
		got := sb.String()
		want := "digraph G {\n}\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("single task", func(t *testing.T) {
		f := &goyek.Flow{}
		f.Define(goyek.Task{Name: "build"})
		sb := &strings.Builder{}
		err := graphviz.Draw(sb, f)
		if err != nil {
			t.Fatal(err)
		}
		got := sb.String()
		if !strings.Contains(got, "  \"build\";\n") {
			t.Errorf("missing node declaration, got:\n%s", got)
		}
	})

	t.Run("task with dependencies", func(t *testing.T) {
		f := &goyek.Flow{}
		test := f.Define(goyek.Task{Name: "test"})
		lint := f.Define(goyek.Task{Name: "lint"})
		f.Define(goyek.Task{
			Name: "all",
			Deps: goyek.Deps{test, lint},
		})
		sb := &strings.Builder{}
		err := graphviz.Draw(sb, f)
		if err != nil {
			t.Fatal(err)
		}
		got := sb.String()
		if !strings.Contains(got, "  \"all\" -> \"test\";\n") {
			t.Errorf("missing edge all -> test, got:\n%s", got)
		}
		if !strings.Contains(got, "  \"all\" -> \"lint\";\n") {
			t.Errorf("missing edge all -> lint, got:\n%s", got)
		}
	})

	t.Run("diamond dependency", func(t *testing.T) {
		f := &goyek.Flow{}
		d := f.Define(goyek.Task{Name: "D"})
		b := f.Define(goyek.Task{Name: "B", Deps: goyek.Deps{d}})
		c := f.Define(goyek.Task{Name: "C", Deps: goyek.Deps{d}})
		f.Define(goyek.Task{Name: "A", Deps: goyek.Deps{b, c}})

		sb := &strings.Builder{}
		err := graphviz.Draw(sb, f)
		if err != nil {
			t.Fatal(err)
		}
		got := sb.String()
		edges := []string{
			"  \"A\" -> \"B\";\n",
			"  \"A\" -> \"C\";\n",
			"  \"B\" -> \"D\";\n",
			"  \"C\" -> \"D\";\n",
		}
		for _, edge := range edges {
			if !strings.Contains(got, edge) {
				t.Errorf("missing edge %q, got:\n%s", edge, got)
			}
		}
	})
}
