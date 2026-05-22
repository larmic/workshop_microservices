# TODO &mdash; Slides

> Offene Punkte für den Online-Gang. Aktuell vor allem **Bildmaterial** zur Auflockerung der text­lastigen Slides.

---

## ✅ Bereits eingesetzte Bilder

- `07-vorbereitung.md` &mdash; Dashboard-Screenshot (`assets/dashboard-ui.png`) als `<img class="dashboard-image">`
- `11-service-discovery.md` &mdash; Registry-Diagramm (`assets/service_discovery.png`)
- `14-circuit-breaker.md` &mdash; Sicherungskasten (`assets/circuitbreaker.png`)
- `17-bulkhead.md` &mdash; Schiffsquerschnitt mit Schotten (`assets/bulkhead.png`)
- `20-saga.md` &mdash; Domino-Reihe (`assets/saga.png`)
- `23-choreography-saga.md` &mdash; synchronisierte Figuren (`assets/choreography_saga.png`)
- `26-tracing.md` &mdash; Paket-Sendung (`assets/tracing.png`)
- `29-zusammenfassung.md` &mdash; Journey-Map (`assets/zusammenfassung.png`)

Pro Pattern-Block gibt es jetzt:
- Eine **Hero-Slide** an Position 0 des Vertikal-Stacks (`{n}-titel.md`) — nur Titel + Subtitle + Bild bei **Opacity 0.40** (Bild ist klar erkennbar, dominiert aber nicht).
- Die **Karten-Slide** + alle **Sub-Slides** mit demselben Bild als **Watermark bei Opacity 0.18** (factor-row und chips überlagern die subtile Bild-Textur).

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
