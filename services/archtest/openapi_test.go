package archtest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Every service with a spec must expose these operations (see the
// openapi-handler-check skill): they are used by Consul, the dashboard and
// Swagger UI.
var requiredOperations = []string{"GET /health", "GET /info", "GET /openapi"}

var httpMethods = map[string]bool{
	"get": true, "post": true, "put": true, "patch": true,
	"delete": true, "head": true, "options": true,
}

// TestOpenAPI compares, per service, the operations declared in
// api/openapi.yaml with the routes registered via mux.HandleFunc / mux.Handle
// in main.go and handler/*.go. Drift in either direction fails: the reference
// implementation is teaching material, participants must be able to trust
// the spec.
func TestOpenAPI(t *testing.T) {
	services := specServices(t)
	if len(services) < 10 {
		t.Fatalf("found only %d services with api/openapi.yaml (%v), expected flight, hotel, car and the booking stories", len(services), services)
	}
	for _, svc := range services {
		t.Run(svc, func(t *testing.T) {
			spec := specOperations(t, filepath.Join(svc, "api", "openapi.yaml"))
			routes := registeredRoutes(t, svc)

			for _, op := range sortedKeys(spec) {
				if !routes[op] {
					t.Errorf("%s: declared in api/openapi.yaml, but no route is registered", op)
				}
			}
			for _, op := range sortedKeys(routes) {
				if !spec[op] {
					t.Errorf("%s: route is registered, but missing in api/openapi.yaml", op)
				}
			}
			for _, op := range requiredOperations {
				if !spec[op] {
					t.Errorf("%s: required endpoint missing", op)
				}
			}
		})
	}
}

// specServices returns every directory (relative to services/) that holds an
// api/openapi.yaml next to a main.go.
func specServices(t *testing.T) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(moduleRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(moduleRoot, path)
		rel = filepath.ToSlash(rel)
		if !d.IsDir() || rel == "." {
			return nil
		}
		if skipDir(rel, d.Name()) {
			return filepath.SkipDir
		}
		if exists(filepath.Join(path, "api", "openapi.yaml")) && exists(filepath.Join(path, "main.go")) {
			out = append(out, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk services: %v", err)
	}
	return out
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// specOperations extracts "METHOD /path" for every operation under the
// top-level "paths:" key. The module has no dependencies, so instead of a
// YAML library this reads the (block-style) structure by indentation: path
// keys are the first indentation level below "paths:", HTTP methods the
// first level below a path key.
func specOperations(t *testing.T, rel string) map[string]bool {
	t.Helper()
	ops := map[string]bool{}
	inPaths := false
	pathIndent, opIndent := -1, -1
	current := ""

	for i, line := range strings.Split(readFile(t, rel), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if strings.HasPrefix(line, "\t") {
			t.Fatalf("%s:%d: tabs are not allowed for YAML indentation", rel, i+1)
		}
		if indent == 0 {
			inPaths = trimmed == "paths:"
			continue
		}
		if !inPaths {
			continue
		}
		key, _, hasColon := strings.Cut(trimmed, ":")
		key = strings.Trim(key, `'"`)
		if pathIndent < 0 {
			pathIndent = indent
		}
		switch {
		case indent == pathIndent:
			if !hasColon || !strings.HasPrefix(key, "/") {
				t.Fatalf("%s:%d: expected a path key below paths:, got %q", rel, i+1, trimmed)
			}
			current, opIndent = key, -1
		case indent > pathIndent:
			if opIndent < 0 {
				opIndent = indent
			}
			if indent == opIndent && hasColon && httpMethods[key] {
				ops[strings.ToUpper(key)+" "+normalizePath(current)] = true
			}
		default:
			t.Fatalf("%s:%d: unexpected indentation below paths:", rel, i+1)
		}
	}
	if len(ops) == 0 {
		t.Fatalf("%s: no operations found below paths:", rel)
	}
	return ops
}

// registeredRoutes collects the string patterns passed to HandleFunc/Handle
// in main.go and handler/*.go. Patterns must carry a method
// ("GET /health", net/http 1.22+), otherwise the spec cannot be compared.
func registeredRoutes(t *testing.T, svc string) map[string]bool {
	t.Helper()
	files := []string{filepath.Join(moduleRoot, svc, "main.go")}
	handlerFiles, _ := filepath.Glob(filepath.Join(moduleRoot, svc, "handler", "*.go"))
	files = append(files, handlerFiles...)

	routes := map[string]bool{}
	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || (sel.Sel.Name != "HandleFunc" && sel.Sel.Name != "Handle") {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				t.Errorf("%s: route pattern is not a string literal, cannot compare with the spec", fset.Position(call.Pos()))
				return true
			}
			pattern, _ := strconv.Unquote(lit.Value)
			method, path, ok := strings.Cut(pattern, " ")
			if !ok || !httpMethods[strings.ToLower(method)] {
				t.Errorf("%s: route %q has no HTTP method (use \"GET /path\")", fset.Position(call.Pos()), pattern)
				return true
			}
			routes[method+" "+normalizePath(strings.TrimSpace(path))] = true
			return true
		})
	}
	return routes
}

// normalizePath drops a trailing slash and the "..." wildcard suffix of
// net/http patterns; path parameters stay as {name} in both worlds.
func normalizePath(p string) string {
	p = strings.ReplaceAll(p, "...}", "}")
	if len(p) > 1 {
		p = strings.TrimSuffix(p, "/")
	}
	return p
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
