# Services — Hinweise

Reference-Implementierung des Microservices-Workshops. Ergänzt die übergeordnete
`../CLAUDE.md` um services-spezifische Konventionen.

## Farbpalette

Einheitliche Palette für **Slides** (`slides/`) und **Dashboard** (`dashboard/`).
Single Source of Truth: CSS-Variablen in `slides/theme.css` (`:root`).
Im Dashboard sind die Werte (historisch) inline in `dashboard/static/index.html`
hinterlegt — beim Anpassen beide Stellen synchron halten.

Ausnahme: `dashboard/static/retro.css` ist ein bewusst separates Easteregg-Theme
(SCUMM-Retro-Look, Aktivierung per Konami-Code) und nimmt an der
Palette-Synchronisation nicht teil.

| Rolle | RGB | Hex | Verwendung |
|---|---|---|---|
| Primär-Dunkel | `rgb( 10,   3,  73)` | `#0A0349` | Box-Hintergrund (Slides), Header / Headings / Card-Akzent (Dashboard) |
| Primär-Violett | `rgb(115,  72, 225)` | `#7348E1` | Box-Rahmen, Links, Buttons, Tab-Aktiv, Progress, Highlights im UI |
| Highlight (Akzent) | `rgb(211, 248, 113)` | `#D3F871` | `.hl` / `<mark>` auf Slides |
| Welle (Slides) | `rgb(238, 231, 251)` | `#EEE7FB` | Dekorative Welle am unteren Slide-Rand |
| Page-Hintergrund (Dashboard) | — | `#f0f2f5` | Body |
| Slide-Hintergrund | — | `#ffffff` | Standard-Slide |
| Slide-Text | — | `#000000` | Standard-Textfarbe auf Slides |
| Text (Dashboard) | — | `#333` | Standard |
| Text gedämpft | — | `#555` / `#888` / `#aaa` | Sekundär / Tertiär |

Status-Farben (Dashboard, Circuit-Breaker-Badges etc.):

| Status | Hintergrund | Text |
|---|---|---|
| closed / ok | `#d4f4dd` | `#1d7a3a` |
| open / Fehler | `#fde0e0` | `#b3261e` |
| halfopen / Warnung | `#fff3cd` | `#8a6d00` |

## Audience-Hinweis

Beim Editieren von Slides oder Dashboard gilt die Workshop-Zielgruppe der
Wurzel-`CLAUDE.md`: Software-Architekt:innen / Tech-Leads. Inhalte werden auf
Deutsch verfasst, Code-Identifier auf Englisch.
