package graphviz_test

import (
	"strings"
	"testing"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/graphviz"
)

func TestDraw_NilWriter(t *testing.T) {
	err := graphviz.Draw(nil, &goyek.Flow{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDraw_NilFlow(t *testing.T) {
	err := graphviz.Draw(&strings.Builder{}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDraw_EmptyFlow(t *testing.T) {
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
}

func TestDraw_SingleTask(t *testing.T) {
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
}

func TestDraw_TaskWithUsage(t *testing.T) {
	f := &goyek.Flow{}
	f.Define(goyek.Task{Name: "build", Usage: "compiles the code"})
	sb := &strings.Builder{}
	err := graphviz.Draw(sb, f)
	if err != nil {
		t.Fatal(err)
	}
	got := sb.String()
	if !strings.Contains(got, "  \"build\" [tooltip=\"compiles the code\"];\n") {
		t.Errorf("missing tooltip, got:\n%s", got)
	}
}

func TestDraw_TaskWithDependencies(t *testing.T) {
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
}

func TestDraw_DiamondDependency(t *testing.T) {
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
}
