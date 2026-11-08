package archtest

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// component maps a package directory (relative to services/) to the
// deployable unit it belongs to: every booking story is a unit of its own
// (booking/story1, booking/story3, …), everything else is grouped by its top
// level directory (flight, hotel, car, dashboard, shared).
func component(rel string) string {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if parts[0] == "booking" && len(parts) > 1 {
		return "booking/" + parts[1]
	}
	return parts[0]
}

// allowed reports whether a package of component from may import a package of
// component to. The rules:
//
//   - each service (booking story, flight, hotel, car, dashboard) imports only
//     itself and shared: stories stay independent of each other, the domain
//     services know neither booking nor dashboard nor each other, and the
//     dashboard talks to the services over HTTP only;
//   - shared imports only shared, never a service package.
func allowed(from, to string) bool {
	return to == from || to == "shared"
}

// TestArchitecture checks the import graph of the module against the rules
// in allowed. It parses the import blocks of every Go file (including tests)
// instead of shelling out to go list, so it needs nothing but the sources.
func TestArchitecture(t *testing.T) {
	mod := modulePath(t)
	fset := token.NewFileSet()
	files := 0
	var violations []string

	err := filepath.WalkDir(moduleRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(moduleRoot, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel != "." && skipDir(rel, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		dir := filepath.ToSlash(filepath.Dir(rel))
		if dir == "." || component(dir) == "archtest" {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		files++
		from := component(dir)
		for _, imp := range f.Imports {
			p, _ := strconv.Unquote(imp.Path.Value)
			target, ok := strings.CutPrefix(p, mod+"/")
			if !ok {
				continue // standard library or third party
			}
			if to := component(target); !allowed(from, to) {
				pos := fset.Position(imp.Pos())
				violations = append(violations, rel+":"+strconv.Itoa(pos.Line)+": "+from+" imports "+target)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk sources: %v", err)
	}
	if files == 0 {
		t.Fatal("no Go files found; is the test running inside services/archtest?")
	}

	sort.Strings(violations)
	for _, v := range violations {
		t.Errorf("forbidden import: %s", v)
	}
	if len(violations) > 0 {
		t.Log("Rule: a service (booking/storyN, flight, hotel, car, dashboard) imports only itself and shared; shared imports only shared.")
	}
}

func TestComponent(t *testing.T) {
	cases := map[string]string{
		"booking/story3/handler": "booking/story3",
		"booking/story8":         "booking/story8",
		"flight/handler":         "flight",
		"shared/chaos":           "shared",
		"dashboard":              "dashboard",
	}
	for in, want := range cases {
		if got := component(in); got != want {
			t.Errorf("component(%q) = %q, want %q", in, got, want)
		}
	}
	if allowed("booking/story4", "booking/story3") {
		t.Error("story4 must not import story3")
	}
	if allowed("shared", "flight") {
		t.Error("shared must not import flight")
	}
	if !allowed("hotel", "shared") {
		t.Error("hotel may import shared")
	}
}
