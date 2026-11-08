package archtest

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// moduleRoot is the directory holding go.mod (services/). Go runs tests with
// the package directory as working directory.
const moduleRoot = ".."

// modulePath reads the module path from go.mod.
func modulePath(t *testing.T) string {
	t.Helper()
	f, err := os.Open(filepath.Join(moduleRoot, "go.mod"))
	if err != nil {
		t.Fatalf("open go.mod: %v", err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(s.Text()), "module "); ok {
			return strings.TrimSpace(rest)
		}
	}
	t.Fatal("no module line in go.mod")
	return ""
}

// skipDir reports directories that never contain Go code of this module.
func skipDir(rel string, name string) bool {
	if strings.HasPrefix(name, ".") || name == "testdata" || name == "node_modules" {
		return true
	}
	switch rel {
	case "bin", "slides", "traefik", "booking/custom":
		return true
	}
	return false
}

func readFile(t *testing.T, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(moduleRoot, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(b)
}
