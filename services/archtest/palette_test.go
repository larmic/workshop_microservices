package archtest

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// The palette table lives in services/CLAUDE.md; the single source of truth
// is the :root block of slides/theme.css, the dashboard repeats the values
// inline. TestPalette checks that the roles are present on both sides (see the
// sync-palette skill). TestPaletteNoForeignColours checks the other direction:
// dashboard and pages use no colour outside the palette. dashboard/static/voxel.css
// is the separate easter-egg theme and deliberately not part of either.
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

// Colours beyond the core roles that services/CLAUDE.md documents.
var documentedColours = []colour{
	{"Violett dunkel", 0x4d, 0x27, 0xa8},
	{"Status-Punkt ok (dunkler Kopf)", 0x22, 0xc5, 0x5e},
	{"Status-Punkt Warnung (dunkler Kopf)", 0xf5, 0x9e, 0x0b},
	{"Status-Punkt Fehler (dunkler Kopf)", 0xef, 0x44, 0x44},
	{"Service Flight (Text)", 0x14, 0x56, 0xb3},
	{"Service Flight (Hintergrund)", 0xe3, 0xef, 0xfd},
	{"Service Hotel (Text)", 0x8a, 0x4a, 0x00},
	{"Service Hotel (Hintergrund)", 0xfd, 0xee, 0xde},
}

// Derived tones that exist today: hover states, borders, a gradient stop, the
// print overflow marker and the easter-egg console banner. New colours do not
// belong here; they go into the palette or reuse an existing role.
var toleratedColours = []colour{
	{"Violett Hover", 0x5a, 0x31, 0xbd},
	{"Highlight Hover", 0xc5, 0xed, 0x5e},
	{"Verlauf Seitenhintergrund", 0x1a, 0x0d, 0x6e},
	{"Rahmen ok", 0xb6, 0xe3, 0xc2},
	{"Rahmen Fehler", 0xf3, 0xb8, 0xb8},
	{"Rahmen kühl", 0xd0, 0xd7, 0xe2},
	{"Rahmen kühl dunkel", 0xc5, 0xd2, 0xdd},
	{"Druck Marginalie", 0x8d, 0x8a, 0x9e},
	{"Druck Marginalie dunkel", 0x8f, 0x8a, 0xb8},
	{"Druck Überlauf-Markierung", 0xee, 0x00, 0x00},
	{"Egg Konsole (Text)", 0x1d, 0x28, 0x38},
	{"Egg Konsole (Hintergrund)", 0xf0, 0xa5, 0x31},
}

var (
	// A leading "&" is a numeric character reference (&#128172;), not a colour.
	hexColour = regexp.MustCompile(`(?:^|[^&0-9A-Za-z])#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})\b`)
	rgbColour = regexp.MustCompile(`rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)`)
)

// coloursIn returns every colour of src as "#RRGGBB", translucent rgba()
// variants reduced to their solid tone.
func coloursIn(src string) map[string]bool {
	found := map[string]bool{}
	for _, m := range hexColour.FindAllStringSubmatch(src, -1) {
		h := strings.ToUpper(m[1])
		if len(h) == 3 {
			h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
		}
		found["#"+h] = true
	}
	for _, m := range rgbColour.FindAllStringSubmatch(src, -1) {
		var r, g, b int
		fmt.Sscan(m[1], &r)
		fmt.Sscan(m[2], &g)
		fmt.Sscan(m[3], &b)
		found[fmt.Sprintf("#%02X%02X%02X", r, g, b)] = true
	}
	return found
}

// neutral reports whether a "#RRGGBB" colour is a grey or an off-white: its
// channels differ by at most 16. Neutrals are free to use.
func neutral(hex string) bool {
	var r, g, b int
	fmt.Sscanf(hex, "#%02X%02X%02X", &r, &g, &b)
	return max(r, g, b)-min(r, g, b) <= 16
}

func TestPaletteNoForeignColours(t *testing.T) {
	allowed := map[string]bool{}
	for _, list := range [][]colour{
		{primaryDark, primaryViolet, mutedViolet, highlight, wave},
		statusColours, documentedColours, toleratedColours,
	} {
		for _, c := range list {
			allowed[c.hex()] = true
		}
	}
	for _, file := range []string{
		"dashboard/static/index.html",
		"../pages/landing/index.html",
		"../pages/vorbereitung/index.html",
		"../pages/feedback/index.html",
		"../pages/material/index.html",
	} {
		for hex := range coloursIn(readFile(t, file)) {
			if !allowed[hex] && !neutral(hex) {
				t.Errorf("%s: %s is not in the palette (services/CLAUDE.md)", file, hex)
			}
		}
	}
}
