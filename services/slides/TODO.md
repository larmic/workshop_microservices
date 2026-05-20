# TODO &mdash; Slides

> Offene Punkte für den Online-Gang. Aktuell vor allem **Bildmaterial** zur Auflockerung der text­lastigen Slides.

## Bild-Empfehlungen

Aktuell gibt es nur zwei Bilder (`assets/intro.png`, `assets/way_to_microservices.png`). Bei 47 Slides ist das wenig. Die folgenden Stellen würden besonders profitieren — sortiert nach Priorität.

### Priorität 1 &mdash; sofortiger Wirkungseffekt

| Slide | Bild-Idee | Such-Begriffe |
|---|---|---|
| `17-bulkhead.md` | **Schiffsquerschnitt mit Schotten** ⭐ — die literale Metapher | `ship cross section bulkhead compartments illustration`, `watertight compartment vector` |
| `26c-tracing-mehr.md` | **Jaeger/Tempo-Waterfall-Mockup** — zeigt was OpenTelemetry-UIs liefern | `jaeger trace waterfall screenshot`, `distributed tracing UI` |

> ✅ `07-vorbereitung.md` &mdash; Dashboard-Screenshot ist als `assets/dashboard-ui.png` integriert.

### Priorität 2 &mdash; gute Metaphern-Anker

| Slide | Bild-Idee | Such-Begriffe |
|---|---|---|
| `14-circuit-breaker.md` | Sicherungskasten / Hardware-CB | `electrical circuit breaker fuse box vector`, `circuit breaker icon` |
| `11-service-discovery.md` | Telefonbuch / Wegweiser / Service-Registry | `phonebook illustration`, `service registry diagram`, undraw.co `directory` |
| `20-saga.md` | Domino-Reihe rückwärts / Wikinger-Saga | `domino chain illustration`, `viking longship`, `narrative thread` |
| `23-choreography-saga.md` | Tänzer ohne Choreograph / Orchester vs. Ensemble | `dance choreography vector`, `orchestra conductor`, `jazz ensemble` |
| `26-tracing.md` | Paketverfolgung / Sendungs-QR-Code | `package tracking illustration`, undraw.co `shipping`, `parcel tracking` |

### Priorität 3 &mdash; Bonus für Closing

| Slide | Bild-Idee | Such-Begriffe |
|---|---|---|
| `29-zusammenfassung.md` | Reisekarte / Zielgerade / Schiff im Hafen | `journey map illustration`, `seven stops map`, undraw.co `destination` |
| `06-12-faktor-app.md` | Heroku-Bezug / 12-Heroku-Originale | Vorsicht Lizenz — eher abstrahieren statt Logo |

## Bild-Quellen mit Lizenz-Notizen

| Quelle | Eigenschaften | Lizenz |
|---|---|---|
| [unDraw](https://undraw.co) | anpassbare Farbe, viele Illustrationen | freie Nutzung, keine Attribution |
| [Storyset](https://storyset.com) | animierbar, vielseitig | Attribution erforderlich |
| [Humaaans](https://humaaans.com) | modulare Figuren | CC0 |
| [Open Doodles](https://opendoodles.com) | handgezeichnet | CC0 |
| [DrawKit](https://drawkit.com) | Vektor, hochwertig | frei mit Attribution |
| [unsplash.com](https://unsplash.com) | Foto-Stock | freie Nutzung |
| [Excalidraw](https://excalidraw.com) | selbst zeichnen (Schotten, Saga-Sequenz) | export PNG/SVG |

## Visueller Stil-Hinweis

Workshop-Palette (siehe `theme.css` und Wurzel-`CLAUDE.md`):

| Farbe | Hex | Verwendung |
|---|---|---|
| Primär-Dunkel | `#0A0349` | Box-Hintergrund, Kapitel-Slides |
| Primär-Violett | `#7348E1` | Rahmen, Links, Akzente |
| Highlight | `#D3F871` | `.hl`, Pull-Quotes |
| Welle | `#EEE7FB` | Dekoration am Slide-Rand |

**Empfehlung:** Bilder am besten in **Lila/Cyan-Tönen** oder **monochrom dunkelblau** &mdash; passt zur Palette. Stark farbige Stock-Fotos wirken fremd. Bei unDraw die Primärfarbe auf `#7348E1` setzen — fügt sich nahtlos ein.

## Sonstige offene Punkte

- [ ] **`docs/questions/story7.md`** anlegen — die anderen Stories haben eine eigene Fragen-Datei, Story 7 (Tracing) nicht. Aktuell sind die Recap-Notes aus `docs/instructions/distributed-tracing.md` (Abschnitt 9 + 10) synthetisiert. Wenn ihr die Diskussion ähnlich tief wie bei Stories 1–6 vorbereiten wollt, lohnt sich eine eigene Datei.
- [ ] **Sprecher-Notizen Story 1 Recap** prüfen — wurden im finalen Polish-Pass auf 6 Boxen erweitert. Falls Zeit knapp, eine der neuen Boxen (Config oder OpenAPI) ggf. wieder rausnehmen.
- [ ] **Foto vom Trainer** für Kapitel 1 / Closing — falls gewünscht persönliche Note.

## Wenn ihr nichts mehr macht

Das Deck ist ohne weitere Bilder vortragbar. Die zwei Closing-Slides + drei Kapitel-Trenner sorgen schon für genug visuelle Variation. Bilder sind die Kür, nicht die Pflicht.
