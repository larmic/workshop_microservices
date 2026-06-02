# Backlog Kickoff-Workshop · 2026-05-30

Priorisierte, umsetzbare TODOs aus dem Feedback. Herleitung und Einordnung in
[feedback.md](feedback.md). Status: `[ ]` offen · `[~]` in Arbeit · `[x]` erledigt.
Pfade relativ zur Repo-Wurzel.

Reihenfolge der Blöcke nach Wirkung: D und B (Didaktik, Pattern-Verständnis) haben
den größten Hebel, A (Quick-Wins) ist schnell erledigt.

## A · Quick-Wins (risikoarm, klar umsetzbar)

- [ ] **A1 · OpenAPI: `format: uuid` entfernen** · Priorität: Hoch
  Betrifft `services/booking/story1..7/api/openapi.yaml` (id-Felder von Flight,
  Hotel, Car) und die Domain-Service-Specs. Die Handler liefern Prefix-IDs
  (`F-…`, `H-…`, `C-…`, `B-…`), keine UUIDs. Empfehlung: nur `type: string`, die
  sprechenden IDs behalten. Beispiele aktualisieren.

- [ ] **A2 · Bulkhead `inFlight` zu `inProgress` umbenennen** · Priorität: Mittel
  Konsistent über: `services/booking/story4..7/bulkhead/bulkhead.go` (Struct-Feld
  und Variablen), Slides `services/slides/chapters/17a-bulkhead-code.md` und
  `17-bulkhead.md`, `services/dashboard/static/index.html` (Label und Metrik),
  `docs/instructions/bulkhead.md`. Vermeidet Verwechslung mit dem Flight-Service.
  Achtung: Admin-Contract (C2) gibt `inFlight` als JSON-Feld aus, gemeinsam ändern.

- [ ] **A3 · "Am Markt" im Web sichtbar machen** · Priorität: Mittel
  `services/slides/theme.css:670-678`: `.market-row` aus dem Fragment-Gating lösen,
  damit die Framework-Übersicht auch im statischen Web erscheint. Trade-off notieren:
  live als Reveal vs. im Web immer sichtbar. Betrifft alle Kapitel mit `.market-row`
  (u. a. `20-saga.md`, `11-service-discovery.md`, `14-circuit-breaker.md`,
  `17-bulkhead.md`, `23-choreography-saga.md`, `26-tracing.md`).

- [ ] **A4 · MicroProfile LRA bei Saga "Am Markt" ergänzen** · Priorität: Niedrig
  `services/slides/chapters/20-saga.md:50-58`: Chip `MicroProfile LRA` hinzufügen
  (passt als Java/Jakarta-Variante neben Axon).

- [ ] **A5 · Setup-Schrift vergrößern** · Priorität: Mittel
  `services/slides/theme.css:510-513`: `.box.setup pre` von `0.85em` auf `1em` bis
  `1.2em`, Zeilenhöhe ggf. auf 1.6. Der Code wird abgetippt, soll gut lesbar sein.

- [ ] **A6 · Kontrast-Pass (WCAG-AA)** · Priorität: Hoch
  `services/slides/theme.css`: Subtitle (`:110`), Links (`:129`), `.cols h3` (`:208`)
  und Inline-`code` (`:293-298`, `:621`) auf mindestens 4.5:1 anheben. Insbesondere
  `code`-Override für dunkle `.box` ergänzen (heute dunkles Lila auf dunklem Grund).
  Danach Palette-Sync mit dem Dashboard prüfen (Skill `sync-palette`,
  `services/CLAUDE.md` Farbpalette).

## B · Slides- und Inhalts-Klärungen (Pattern-Verständnis)

- [ ] **B1 · 12-Factor zu Microservices verbrücken** · Priorität: Hoch
  In `services/slides/chapters/06-12-faktor-app.md`, `docs/themen.md:72-82` und
  `docs/stories/story-01-cloud-native-booking-service.md` klarstellen: 12-Factor ist
  cloud-native und breiter als Microservices, notwendig aber nicht hinreichend.
  Microservices ergänzen Service-Schnitt, unabhängige Deploybarkeit, Resilienz.

- [ ] **B2 · Abgrenzungs-Folie CB vs. Bulkhead vs. Rate-Limit** · Priorität: Hoch
  Neue Folie oder Exkurs vor dem Bulkhead-Hands-on: drei Patterns nebeneinander
  (zu viele Fehler vs. zu viel Parallelität vs. zu viel Durchsatz), warum sie
  komplementär sind, warum der Bulkhead bewusst keine Fallback-Antwort liefert.
  Inhalt existiert in `docs/questions/story4.md` und `docs/instructions/bulkhead.md`,
  gehört auf die Slides.

- [ ] **B3 · Bulkhead `rejected` erklären** · Priorität: Mittel
  Auf der Bulkhead-Folie (`services/slides/chapters/17-bulkhead.md`) klarstellen:
  `rejected` = sofortige Ablehnung bei vollem Pool (Fail-Fast), Antwort 503 mit
  Header `X-Bulkhead-Full`, kein Queueing. Ggf. Kommentar in
  `booking/story4/bulkhead/bulkhead.go:61-64`.

- [ ] **B4 · Bulkhead-Demo: Aufruf-Reihenfolge sichtbar machen** · Priorität: Mittel
  Im Demo-Hinweis erklären, dass die Aufrufe sequenziell Flight, Hotel, Car laufen
  (`booking/story4/handler/booking.go:73-78`) und das die Rejection-Kaskade im
  Dashboard formt. Kommentar im Handler ergänzen.

- [ ] **B5 · Recap Story 3: CB nicht in `/health`** · Priorität: Mittel
  Antwort explizit auf die Recap-Folie: `/health` ist nur Liveness, die Probe macht
  z. B. K8s-Readiness, der CB-State gehört auf `/admin/circuit-state`. Der Code macht
  es bereits richtig (`shared/handler/health.go`, `booking/story3/main.go`).

- [ ] **B6 · CB auf POST diskutieren** · Priorität: Mittel
  Auf einer Folie thematisieren, warum CB auf schreibende POST-Operationen heikel ist
  (mögliche Teilausführung, Seiteneffekt im Half-Open, Idempotenz-Annahme). Bezug zu
  `booking/story3/handler/booking.go`. Verknüpft mit E2.

- [ ] **B7 · "Gilt das auch im Monolithen?" pro Pattern beantworten** · Priorität: Hoch · *Entscheidung Lars*
  Kurze "Monolith vs. Microservices"-Notiz je Story (Circuit Breaker, Bulkhead,
  Saga, Tracing): welches Problem löst das Pattern, gilt es auch im Monolithen, was
  wird im verteilten System schlimmer. Adressiert den Kernbefund.

- [ ] **B8 · "Warum nur 1 DB pro Service" motivieren** · Priorität: Mittel
  Diskussions-Anker in `docs/themen.md` (Block 5, bei `:82`): unabhängige
  Schema-Evolution, differenzierte Skalierung, Transaktionsgrenzen als
  Service-Grenzen, kein Shared-State.

## C · Neue Doku-Seiten

- [ ] **C1 · `docs/troubleshooting.md`** · Priorität: Hoch
  Windows/WSL2: localhost und Traefik-Timing ("wenn nicht erreichbar, kurz warten
  und neu laden"), Docker-Desktop-Networking, Port-Konflikte (80, 8080, 8500),
  Image-Pull. GitHub-Zertifikat in WSL (Gast-WLAN, Proxy, Zertifikat hinterlegen).
  Verlinken aus `docs/vorbereitung.md`.

- [ ] **C2 · `docs/admin-contract.md`** · Priorität: Hoch
  Exakter Dashboard-Vertrag mit JSON-Schemata, damit Custom-Services in jeder
  Sprache das Dashboard füllen: `GET /health`; ab Story 3 `GET /admin/circuit-state`;
  ab Story 4 `GET /admin/bulkhead-state` plus `POST /admin/bulkhead-reset`; ab
  Story 5 `GET /admin/sagas` plus `POST /admin/sagas-reset`. Felder aus
  `booking/story3/circuitbreaker/circuitbreaker.go`,
  `booking/story4/bulkhead/bulkhead.go`, `booking/story5/saga/saga.go` übernehmen.
  In den Stories und der Vorbereitung verlinken.

- [ ] **C3 · OpenAPI-/Story-Hinweise im Dashboard prominenter** · Priorität: Niedrig
  Basis vorhanden (`dashboard/static/index.html:1300-1310` verlinkt
  `/api/{service}/openapi`). Pro Story den passenden Spec-Link als sichtbaren Hinweis
  zur jeweiligen Aufgabe anbieten.

## D · Didaktik & Curriculum · *Entscheidung Lars*

- [ ] **D1 · Motivation pro Pattern schärfen** · Priorität: Hoch
  Vor jeder Story kurz: welches konkrete Problem, warum schmerzt es ohne das Pattern.
  Mehr "warum" statt "wie". Kernbefund.

- [ ] **D2 · Nicht-codende Aufgabe ergänzen** · Priorität: Hoch
  Whiteboard- oder DDD-Warm-up: "Wie würdet ihr das System schneiden?" oder
  "Wie würdet ihr dieses Problem lösen?" vor der Lösung. Eigenes Denken anregen,
  dann mit der Referenz vergleichen.

- [ ] **D3 · Stories teilweise "auf der Tonspur"** · Priorität: Mittel
  Nicht alle 7 voll implementieren. Kandidaten zum Durchsprechen statt Bauen:
  Story 7 (Tracing), Teile von Story 6 (Choreography). In `docs/themen.md` und den
  Story-Headern als Pflicht vs. Vertiefung kennzeichnen.

- [ ] **D4 · Theorie/Praxis über beide Tage ausbalancieren** · Priorität: Mittel
  Tag 2 entzerren (`docs/themen.md:26-39`): Micro-Lectures zwischen die Stories
  (z. B. "Warum Saga?" vor Story 5, "Orchestration vs. Choreography" vor Story 6).

- [ ] **D5 · Voraussetzungen schärfen, ggf. Vor-Workshop** · Priorität: Hoch
  In `docs/vorbereitung.md:5-10` als harte Vorbedingung formulieren: Teilnehmende
  müssen ein Greenfield-Projekt in ihrer Sprache aufsetzen können, das per Dockerfile
  ein Image baut und einen Health-Endpoint bereitstellt. Optionaler Vor-Workshop
  "Was ist OpenAPI? Wie geht Docker?". Setup-Dauer realistisch einplanen.

- [ ] **D6 · Vorbereitungs- und Repo-Strategie** · Priorität: Mittel
  Material aktiv vorab teilen (Repo ist bereits öffentlich), Selbstlern-Charakter
  ausbauen, ggf. zwei Termine mit einer Woche Abstand zum echten Nacharbeiten.

## E · Architektur- und Design-Entscheidungen · *Entscheidung Lars*

- [ ] **E1 · Story 5: Saga-Schema im API-Vertrag überdenken** · Priorität: Mittel
  Erfolg liefert `Booking`, Fehlerfall (503) liefert `SagaFailure` (= Saga)
  (`booking/story5/api/openapi.yaml:114-126`). Option: Fehlerfall ebenfalls als
  `Booking` mit Status liefern und das `Saga`-Schema aus dem öffentlichen Vertrag
  nehmen, oder die Trennung bewusst dokumentieren (Diskussions-Anker für die Story).

- [ ] **E2 · CB-Scope in Story 3** · Priorität: Mittel
  Prüfen, ob der Circuit Breaker in Story 3 auf GET beschränkt wird und die
  POST-Resilienz erst in Story 5 (Saga) sauber gelöst wird. Verknüpft mit B6.
