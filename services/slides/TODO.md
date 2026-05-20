# TODO &mdash; Slides

> Offene Punkte für den Online-Gang. Aktuell vor allem **Bildmaterial** zur Auflockerung der text­lastigen Slides.

---

## ✅ Bereits eingesetzte Bilder

- `07-vorbereitung.md` &mdash; Dashboard-Screenshot (`assets/dashboard-ui.png`) als `<img class="dashboard-image">`
- `11-service-discovery.md` &mdash; Registry-Diagramm (`assets/service_discovery.png`)
- `14-circuit-breaker.md` &mdash; Sicherungskasten (`assets/circuitbreaker.png`)
- `17-bulkhead.md` &mdash; Schiffsquerschnitt mit Schotten (`assets/bulkhead.png`)
- `20-saga.md` &mdash; Domino-Reihe (`assets/saga.png`)
- `26-tracing.md` &mdash; Paket-Sendung (`assets/tracing.png`)

Pro Pattern-Block gibt es jetzt:
- Eine **Hero-Slide** an Position 0 des Vertikal-Stacks (`{n}-titel.md`) — nur Titel + Subtitle + Bild bei **Opacity 0.40** (Bild ist klar erkennbar, dominiert aber nicht).
- Die **Karten-Slide** + alle **Sub-Slides** mit demselben Bild als **Watermark bei Opacity 0.18** (factor-row und chips überlagern die subtile Bild-Textur).

Choreography-Hero (`23-titel.md`) hat aktuell noch kein Bild — kommt sobald `choreography.png` generiert ist.

---

## Noch offene AI-Bild-Prompts

**Jeder Prompt ist self-contained.** Block kopieren, in deinen Generator einfügen (Midjourney, DALL-E 3, Flux, Stable Diffusion XL, Imagen 3, Meshy, …) und absenden. Alle Vorgaben (Farb-Hex-Codes, Aspect Ratio, Negativ-Liste) sind im Prompt-Text enthalten.

Workshop-Palette zur Erinnerung — die Hex-Codes erscheinen in jedem Prompt:

| Farbe | Hex | Rolle |
|---|---|---|
| Primär-Dunkel | `#0A0349` | Outlines, dominante Flächen |
| Primär-Lila | `#7348E1` | Akzente, Verbindungen |
| Highlight | `#D3F871` | kleine Status-/Warn-Tupfer |
| Welle | `#EEE7FB` | dezente Hintergrund-Flächen |
| Background | `#FFFFFF` | Slide-Hintergrund |

---

### 1. Tracing Waterfall (Slide `26c-tracing-mehr.md`)

**Was es zeigt:** Mockup einer Tracing-UI à la Jaeger / Grafana Tempo — ein „Wasserfall" verschachtelter Spans.

```
A clean UI mockup of a distributed tracing dashboard panel in modern minimal flat design, in the visual style of Jaeger or Grafana Tempo but explicitly without any product branding or logos. The panel has a header bar at the top reading "Trace · 1.2s" with a subtle ID badge. Below the header is a vertical waterfall of nested colored bars representing HTTP spans. From top to bottom: one wide top-level span occupying the full panel width labeled "POST booking 1.2s"; indented below it three parallel-stacked spans of varying widths labeled "flight 120ms", "hotel 180ms" (with a small red error dot at its right end), and "car 900ms"; and a final shorter indented span at the bottom labeled "DELETE flight 80ms" representing a compensation step. On the right side a thin sidebar shows span details with key-value pairs. Subtle drop shadow under the main card. White background hex #FFFFFF, light gray hex #F5F6F8 for the inner content area. Aspect ratio 16:9, landscape orientation, output size 1792x1024 pixels. Color palette: dark blue-purple hex #0A0349 for the header bar and dark text, light purple hex #7348E1 for the span bars in varying opacity for hierarchy, lime green hex #D3F871 for healthy status accents, a muted red tone for the single error indicator. Negative prompt: no real product names or logos, no Jaeger or Grafana branding, no people, not photorealistic, no heavy shadows.
```

---

### 2. Choreography (Slide `23-choreography-saga.md`)

**Was es zeigt:** Abstrakte Figuren in synchronisierter Bewegung ohne Dirigenten — sie kommunizieren über Events.

```
An abstract flat vector illustration in modern minimal friendly style, showing four stylized human figure silhouettes arranged in a loose circle on a white background. The figures are simple geometric shapes with no facial features, with arms raised in motion as if dancing or performing in sync. Smooth curved arrows flow between them in a counter-clockwise pattern, each arrow carrying a small lime-green dot-shaped event icon along its path. The arrangement explicitly shows synchronized movement without a central conductor — each figure responds to the previous one through the event arrows. White background hex #FFFFFF. Aspect ratio 16:9, landscape orientation, output size 1792x1024 pixels. Color palette: dark blue-purple hex #0A0349 for the figure silhouettes, light purple hex #7348E1 for the connecting arrows, lime green hex #D3F871 for the event dots along the arrows. Negative prompt: no realistic faces, no central figure or conductor, no text, no letters, no words, no logos, not photorealistic, no heavy shadows.
```

---

### 3. Closing — Journey Map (Slide `29-zusammenfassung.md`)

**Was es zeigt:** Reise-Karte mit 7 Stationen — Sinnbild für die 7 Workshop-Stories.

```
A modern flat illustration in friendly minimal vector style, like an illustrated board-game adventure map, showing a curving path or trail on a stylized landscape map viewed slightly from above as if from a hot-air balloon. The path winds gently through seven milestone pins from the left (start) to the right (destination). Each pin is a small flag or location marker, subtly numbered 1 through 7. Around the path are abstract hills, geometric trees, and simple cloud shapes suggesting an illustrated adventure landscape. White background hex #FFFFFF with subtle pale lavender hex #EEE7FB clouds. Aspect ratio 16:9, landscape orientation, output size 1792x1024 pixels. Color palette: dark blue-purple hex #0A0349 for the winding path and pin outlines, light purple hex #7348E1 for the pin flags and accent markers, lime green hex #D3F871 for the dots on completed milestones, pale lavender hex #EEE7FB for the clouds and distant hills. Negative prompt: no real geography, no recognizable landmarks, no people, no text labels except small numbers 1 through 7 on the pins, not photorealistic, no heavy shadows.
```

---

### Provider-spezifische Tipps (optional)

- **Midjourney v6** — am Ende des Prompts anhängen: `--ar 16:9 --v 6 --style raw --no text, letters, people, photorealistic`
- **Stable Diffusion XL** — den „Negative prompt:"-Teil hinten herausziehen und in den Negative-Slot kopieren; den Rest in den Positive-Slot.
- **DALL-E 3 / Imagen 3 / Flux** — Prompt 1:1 einfügen, keine Sonderbehandlung nötig.
- Wenn das erste Ergebnis zu kantig oder zu langweilig ist, am Stil-Anfang ergänzen: `in the style of undraw.co illustrations`, `friendly geometric shapes`, oder `subtle gradient fills`.

### Nach dem Generieren

1. Bild als PNG mit weißem (oder transparentem) Hintergrund ablegen unter `services/slides/assets/`. Empfohlene Namen: `tracing-waterfall.png`, `choreography.png`, `journey-map.png`.
2. Im jeweiligen Slide-Markdown einbinden &mdash; **zwei Varianten**:

   **Variante A &mdash; als Watermark-Hintergrund** (wie auf den anderen Pattern-Slides). Ganz oben in der Datei als Slide-Kommentar:
   ```markdown
   <!-- .slide: data-background-image="./assets/choreography.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->
   ```
   Opacity bei Bedarf anpassen: `0.10`&ndash;`0.15` für sehr dezente Textur, `0.20`&ndash;`0.30` für klarer sichtbare Bilder. Geeignet für `choreography.png` (analoges Pattern wie Bulkhead/Saga).

   **Variante B &mdash; als sichtbares Bild im Slide** (für UI-Mockups oder Hero-Bilder am Closing). Im Markdown an passender Stelle:
   ```html
   <img class="dashboard-image" src="./assets/tracing-waterfall.png" alt="..."/>
   ```
   Geeignet für `tracing-waterfall.png` (UI-Mockup mit Detail-Inhalt, der bei niedriger Opacity unleserlich würde) und `journey-map.png` (Closing-Highlight, soll sichtbar sein).

---

## Falls du Bilder lieber selbst suchst (statt KI)

| Quelle | Eigenschaften | Lizenz |
|---|---|---|
| [unDraw](https://undraw.co) | anpassbare Farbe, viele Illustrationen | freie Nutzung, keine Attribution |
| [Storyset](https://storyset.com) | animierbar, vielseitig | Attribution erforderlich |
| [Humaaans](https://humaaans.com) | modulare Figuren | CC0 |
| [Open Doodles](https://opendoodles.com) | handgezeichnet | CC0 |
| [DrawKit](https://drawkit.com) | Vektor, hochwertig | frei mit Attribution |
| [unsplash.com](https://unsplash.com) | Foto-Stock | freie Nutzung |
| [Excalidraw](https://excalidraw.com) | selbst zeichnen (Schotten, Saga-Sequenz) | export PNG/SVG |

Bei unDraw die Primärfarbe auf `#7348E1` setzen — fügt sich nahtlos in die Slide-Palette ein.

---

## Sonstige offene Punkte

- [ ] **`docs/questions/story7.md`** anlegen — die anderen Stories haben eine eigene Fragen-Datei, Story 7 (Tracing) nicht. Aktuell sind die Recap-Notes aus `docs/instructions/distributed-tracing.md` (Abschnitt 9 + 10) synthetisiert. Wenn ihr die Diskussion ähnlich tief wie bei Stories 1–6 vorbereiten wollt, lohnt sich eine eigene Datei.
- [ ] **Sprecher-Notizen Story 1 Recap** prüfen — wurden im finalen Polish-Pass auf 6 Boxen erweitert. Falls Zeit knapp, eine der neuen Boxen (Config oder OpenAPI) ggf. wieder rausnehmen.
- [ ] **Foto vom Trainer** für Kapitel 1 / Closing — falls gewünscht persönliche Note.

## Wenn ihr nichts mehr macht

Das Deck ist ohne weitere Bilder vortragbar. Die zwei Closing-Slides + drei Kapitel-Trenner sorgen schon für genug visuelle Variation. Bilder sind die Kür, nicht die Pflicht.
