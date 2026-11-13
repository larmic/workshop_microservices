# Services — Hinweise

Reference-Implementierung des Microservices-Workshops. Ergänzt die übergeordnete
`../CLAUDE.md` um services-spezifische Konventionen.

## Farbpalette

Einheitliche Palette für **Slides** (`slides/`) und **Dashboard** (`dashboard/`).
Single Source of Truth: CSS-Variablen in `slides/theme.css` (`:root`).
Im Dashboard sind die Werte (historisch) inline in `dashboard/static/index.html`
hinterlegt — beim Anpassen beide Stellen synchron halten.

Ausnahme: `dashboard/static/voxel.css` ist ein bewusst separates Easteregg-Theme
(Voxel-Look „Inselwelt“ mit den Tokens von larmic.de, Aktivierung per
Konami-Code) und nimmt an der Palette-Synchronisation nicht teil.

| Rolle | RGB | Hex | Verwendung |
|---|---|---|---|
| Primär-Dunkel | `rgb( 10,   3,  73)` | `#0A0349` | Box-Hintergrund (Slides), Header / Headings / Card-Akzent (Dashboard) |
| Primär-Violett | `rgb(115,  72, 225)` | `#7348E1` | Box-Rahmen, Links, Buttons, Tab-Aktiv, Progress, Highlights im UI |
| Highlight (Akzent) | `rgb(211, 248, 113)` | `#D3F871` | `.hl` / `<mark>` auf Slides |
| Welle (Slides) | `rgb(238, 231, 251)` | `#EEE7FB` | Dekorative Welle am unteren Slide-Rand |
| Page-Hintergrund (Dashboard) | — | `#f0f2f5` | Body |
| Slide-Hintergrund | — | `#ffffff` | Standard-Slide |
| Slide-Text | `rgb( 10,   3,  73)` | `#0A0349` | Standard-Textfarbe auf Slides (Primär-Dunkel, kein Schwarz) |
| Text (Dashboard) | — | `#333` | Standard |
| Text gedämpft | — | `#555` / `#888` / `#aaa` | Sekundär / Tertiär |
| Violett dunkel | — | `#4d27a8` | Inline-Code und dunkle Link-Variante auf Slides |
| Text gedämpft (Violett) | — | `#5B4E93` | Kleine Labels, Badges, Kicker (Slides: Agenda, Merksätze; Dashboard: statt `#888` für kleinen Text, ≥ 6:1 auf hellem Grund) |

Einzelfarben nur für den Eisberg auf Folie 08b (kein Teil der Palette):
`#3A1F9E` mittlere Schicht zwischen Violett und Primär-Dunkel, `#E5484D`
rote Boje am Tiefenmesser.

Status-Farben (Dashboard, Circuit-Breaker-Badges, Health- und Startup-Symbole auf hellem Grund):

| Status | Hintergrund | Text |
|---|---|---|
| closed / ok | `#d4f4dd` | `#1d7a3a` |
| open / Fehler | `#fde0e0` | `#b3261e` |
| halfopen / Warnung | `#fff3cd` | `#8a6d00` |

Status-Punkte im dunklen Dashboard-Kopf (Health-Chip): `#22c55e` ok,
`#f59e0b` Warnung, `#ef4444` Fehler. Die Status-Textfarben wären dort zu
dunkel, deshalb diese drei kräftigeren Töne, nur an dieser Stelle.

Service-Farben (Chips im Dashboard, Karten und Schilder der Druckseite):

| Service | Text | Hintergrund |
|---|---|---|
| Flight | `#1456b3` | `#e3effd` |
| Hotel | `#8a4a00` | `#fdeede` |
| Car | `#1d7a3a` | `#e6f4ea` |

Das Palette-Gate (`scripts/gates.sh palette`) prüft beide Richtungen: Die
Rollen kommen in Slides und Dashboard vor, und Dashboard und `pages/` nutzen
keine Farbe außerhalb dieser Tabellen. Grau- und Fast-Weiß-Töne sind frei.
Abgeleitete Töne (Hover, Rahmen) stehen als Ausnahmeliste in
`archtest/palette_test.go`; neue Farben gehören in die Palette, nicht dorthin.

## Schriften (Slides)

Die Slides nutzen **Poppins** (SIL OFL 1.1), selbst gehostet
unter `slides/assets/fonts/` und per `@font-face` in `slides/theme.css`
eingebunden. Einzige Ausnahme ist **Sacramento** (SIL OFL 1.1), eine
Schreibschrift nur für den „roten Faden“ auf der Feedback-Folie. Damit rendern Checkout, CI, GitHub Pages, Docker und Beamer
identisch, auch offline. Titel laufen in ExtraBold (800), Fließtext in
Regular (400). Keine System- oder Lizenzschriften (Futura, Faktum) im Stack;
Details in `slides/assets/fonts/README.md`. Die Druckseite
`pages/material/` nutzt ebenfalls Poppins (Kopie unter `pages/material/fonts/`),
damit die Bogen auf jedem Rechner gleich umbrechen. Dashboard und die übrigen
Pages nutzen die Systemschrift.

## Audience-Hinweis

Beim Editieren von Slides oder Dashboard gilt die Workshop-Zielgruppe der
Wurzel-`CLAUDE.md`: Software-Architekt:innen / Tech-Leads. Inhalte werden auf
Deutsch verfasst, Code-Identifier auf Englisch.
