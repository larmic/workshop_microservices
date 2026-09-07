# Schriften der Slides

Die Folien nutzen **Poppins** (Indian Type Foundry, Jonny Pinhorn, Ninad Kale)
als einzige Schriftfamilie, selbst gehostet, damit Checkout, CI, GitHub Pages,
Docker und der Beamer offline identisch rendern.

| Datei | Schnitt | Verwendung |
|---|---|---|
| `poppins-400.woff2` | Regular | Fließtext |
| `poppins-400-italic.woff2` | Italic | Lead-Sätze, Zitate |
| `poppins-500.woff2` | Medium | Claims, Rollen |
| `poppins-600.woff2` | SemiBold | Kicker, Pillen |
| `poppins-700.woff2` | Bold | Kartentitel |
| `poppins-800.woff2` | ExtraBold | Folientitel, Kapitelnummern |

Dazu **Sacramento** (Astigmatic) als verbundene Schreibschrift für genau eine
Stelle: den „roten Faden" auf der Feedback-Folie (`sacramento-400.woff2`).

Quelle: Google Fonts (`fonts.googleapis.com/css2?family=Poppins` bzw.
`Sacramento`), Subset `latin`, geladen am 2026-09-07. Lizenz: SIL Open Font
License 1.1, siehe `OFL-poppins.txt` und `OFL-sacramento.txt`. Eingebunden
über `@font-face` in `../../theme.css`.

Die Hausschrift Faktum ist kommerziell lizenziert und liegt bewusst nicht im
Repo (siehe `internal/`, gitignored).
