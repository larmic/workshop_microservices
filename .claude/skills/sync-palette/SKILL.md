---
name: sync-palette
description: Prüft, ob die Workshop-Farbpalette in slides/theme.css und dashboard/static/index.html synchron ist. Aufrufen, wenn Farben in Slides oder Dashboard geändert wurden oder wenn eine Konsistenzprüfung der visuellen Identität gewünscht ist.
---

# sync-palette

Die Single Source of Truth der Workshop-Farbpalette sind die CSS-Variablen in
`services/slides/theme.css` (`:root`). Im Dashboard sind dieselben Werte historisch
**inline** in `services/dashboard/static/index.html` hinterlegt. Beide Stellen müssen
synchron gehalten werden (siehe `services/CLAUDE.md`).

## Referenz-Palette

| Rolle              | RGB                  | Hex       |
|--------------------|----------------------|-----------|
| Primär-Dunkel      | `rgb(10, 3, 73)`     | `#0A0349` |
| Primär-Violett     | `rgb(115, 72, 225)`  | `#7348E1` |
| Highlight (Akzent) | `rgb(211, 248, 113)` | `#D3F871` |
| Welle (Slides)     | `rgb(238, 231, 251)` | `#EEE7FB` |

## Ablauf

1. Lies `services/slides/theme.css` und extrahiere alle Werte aus dem `:root`-Block
   (`--slide-text`, `--wave`, `--box-bg`, `--box-border`, `--highlight`, …).
   Normalisiere zu Hex (Großschreibung, ohne Leerzeichen).
2. Lies `services/dashboard/static/index.html` und extrahiere alle Farbwerte aus dem
   `<style>`-Block. Sammle Hex (`#0A0349`) und `rgb(...)`-Notationen,
   normalisiere ebenfalls zu Hex.
3. Vergleiche pro Rolle:
   - **Primär-Dunkel** → muss in beiden Quellen mindestens einmal vorkommen.
   - **Primär-Violett** → muss in beiden Quellen mindestens einmal vorkommen.
   - **Highlight** → muss in `theme.css` vorkommen (Dashboard nutzt es selten,
     kein Pflichtmatch).
4. Melde Drift mit Datei + Zeile:
   - Fehlende Rolle in einer Quelle.
   - Abweichende Schreibweise (`#0a0349` vs. `#0A0349`) → Hinweis,
     kein Fehler.
   - Unbekannte zusätzliche Farben → als Info auflisten, damit klar wird,
     ob sie in die Tabelle in `services/CLAUDE.md` aufgenommen werden müssen.
5. Schlage konkrete Edits vor (alte → neue Zeile), führe sie aber nicht ohne
   Rückfrage aus.

## Status-Farben (Dashboard-only, nicht in theme.css)

Diese gelten nur fürs Dashboard und müssen **nicht** in den Slides existieren,
sind aber Teil der `services/CLAUDE.md`-Tabelle und sollten bei
Dashboard-Refactorings stabil bleiben:

| Status            | Hintergrund | Text      |
|-------------------|-------------|-----------|
| closed / ok       | `#d4f4dd`   | `#1d7a3a` |
| open / Fehler     | `#fde0e0`   | `#b3261e` |
| halfopen / Warnung| `#fff3cd`   | `#8a6d00` |

Wenn die Status-Farben fehlen oder verändert sind, in `services/dashboard/static/index.html`
melden, aber **nicht** an `services/slides/theme.css` synchronisieren.

## Ausgabeformat

Liefere einen kurzen Bericht in dieser Form:

```
Farbpalette-Sync
================
✓ Primär-Dunkel       #0A0349   theme.css:23, index.html:17
✓ Primär-Violett      #7348E1   theme.css:24, index.html:65
✗ Highlight           #D3F871   theme.css:26, index.html: FEHLT
ℹ Status-closed       #d4f4dd   index.html:142

Empfehlung: …
```

## Audience-Hinweis

Antworten und Vorschläge auf Deutsch verfassen (Workshop-Zielgruppe:
Architekt:innen / Tech-Leads, siehe `services/CLAUDE.md`).
