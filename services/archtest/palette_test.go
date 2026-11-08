package archtest

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// The palette table lives in services/CLAUDE.md; the single source of truth
// is the :root block of slides/theme.css, the dashboard repeats the values
// inline. This gate checks that the roles are present on both sides (see the
// sync-palette skill). dashboard/static/voxel.css is the separate easter-egg
// theme and deliberately not part of it.
type colour struct {
	role    string
	r, g, b int
}

var (
	primaryDark   = colour{"Primär-Dunkel", 10, 3, 73}             // #0A0349
	primaryViolet = colour{"Primär-Violett", 115, 72, 225}         // #7348E1
	mutedViolet   = colour{"Text gedämpft (Violett)", 91, 78, 147} // #5B4E93
	highlight     = colour{"Highlight", 211, 248, 113}             // #D3F871
	wave          = colour{"Welle", 238, 231, 251}                 // #EEE7FB
)

var statusColours = []colour{
	{"Status closed (Hintergrund)", 0xd4, 0xf4, 0xdd},
	{"Status closed (Text)", 0x1d, 0x7a, 0x3a},
	{"Status open (Hintergrund)", 0xfd, 0xe0, 0xe0},
	{"Status open (Text)", 0xb3, 0x26, 0x1e},
	{"Status halfopen (Hintergrund)", 0xff, 0xf3, 0xcd},
	{"Status halfopen (Text)", 0x8a, 0x6d, 0x00},
}

func (c colour) hex() string { return fmt.Sprintf("#%02X%02X%02X", c.r, c.g, c.b) }

// in reports whether css contains the solid colour as #RRGGBB (any case) or
// as rgb(r, g, b) with arbitrary spacing. Translucent rgba() variants do not
// count: they are derived from the role, not the role itself.
func (c colour) in(css string) bool {
	if strings.Contains(strings.ToUpper(css), c.hex()) {
		return true
	}
	re := regexp.MustCompile(fmt.Sprintf(`rgb\(\s*%d\s*,\s*%d\s*,\s*%d\s*\)`, c.r, c.g, c.b))
	return re.MatchString(css)
}

var rootBlock = regexp.MustCompile(`(?s):root\s*\{(.*?)\}`)

func TestPalette(t *testing.T) {
	theme := readFile(t, "slides/theme.css")
	m := rootBlock.FindStringSubmatch(theme)
	if m == nil {
		t.Fatal("slides/theme.css: no :root block found")
	}
	root := m[1]
	dashboard := readFile(t, "dashboard/static/index.html")

	for _, c := range []colour{primaryDark, primaryViolet, highlight, wave} {
		if !c.in(root) {
			t.Errorf("slides/theme.css :root: %s %s missing", c.role, c.hex())
		}
	}
	if !mutedViolet.in(theme) {
		t.Errorf("slides/theme.css: %s %s missing", mutedViolet.role, mutedViolet.hex())
	}
	for _, c := range append([]colour{primaryDark, primaryViolet, mutedViolet}, statusColours...) {
		if !c.in(dashboard) {
			t.Errorf("dashboard/static/index.html: %s %s missing (palette in services/CLAUDE.md)", c.role, c.hex())
		}
	}
}
