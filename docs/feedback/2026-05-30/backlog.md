# Backlog Kickoff-Workshop · 2026-05-30

Priorisierte, umsetzbare TODOs aus dem Feedback. Herleitung und Einordnung in
[feedback.md](feedback.md). Status: `[ ]` offen · `[~]` in Arbeit · `[x]` erledigt.
Pfade relativ zur Repo-Wurzel.

Reihenfolge der Blöcke nach Wirkung: D und B (Didaktik, Pattern-Verständnis) haben
den größten Hebel, A (Quick-Wins) ist schnell erledigt.

## A · Quick-Wins (risikoarm, klar umsetzbar)

**Status: A1 bis A6 umgesetzt am 2026-06-02** (go build/vet/test grün). Hinweise je Punkt.

- [x] **A1 · OpenAPI: `format: uuid` entfernen** · Priorität: Hoch
  Betrifft `services/booking/story1..7/api/openapi.yaml` (id-Felder von Flight,
  Hotel, Car) und die Domain-Service-Specs. Die Handler liefern Prefix-IDs
  (`F-…`, `H-…`, `C-…`, `B-…`), keine UUIDs. Empfehlung: nur `type: string`, die
  sprechenden IDs behalten. Beispiele aktualisieren.

- [x] **A2 · Bulkhead `inFlight` zu `inProgress` umbenennen** · Priorität: Mittel
  Konsistent über: `services/booking/story4..7/bulkhead/bulkhead.go` (Struct-Feld
  und Variablen), Slides `services/slides/chapters/17a-bulkhead-code.md` und
  `17-bulkhead.md`, `services/dashboard/static/index.html` (Label und Metrik),
  `docs/instructions/bulkhead.md`. Vermeidet Verwechslung mit dem Flight-Service.
  Achtung: Admin-Contract (C2) gibt das JSON-Feld aus, gemeinsam geändert (heißt
  jetzt `inProgress`), Go-Struct via gofmt neu ausgerichtet, Dashboard-Label
  "in flight" zu "in progress".
  Bewusst NICHT umbenannt: der Circuit Breaker nutzt intern `probeInFlight`
  (`booking/story3..7/circuitbreaker/circuitbreaker.go`, eine Probe-Anfrage im
  HALF_OPEN). Das Feld ist nicht in Spec, Dashboard oder den CB-Slides sichtbar, und
  "probe in flight" ist idiomatisch. Falls dennoch konsistent gewünscht: separater
  Mini-Rename, offen.

- [x] **A3 · "Am Markt" im Web sichtbar machen** · Priorität: Mittel
  `services/slides/theme.css:670-678`: `.market-row` aus dem Fragment-Gating lösen,
  damit die Framework-Übersicht auch im statischen Web erscheint. Trade-off notieren:
  live als Reveal vs. im Web immer sichtbar. Betrifft alle Kapitel mit `.market-row`
  (u. a. `20-saga.md`, `11-service-discovery.md`, `14-circuit-breaker.md`,
  `17-bulkhead.md`, `23-choreography-saga.md`, `26-tracing.md`).

- [x] **A4 · MicroProfile LRA bei Saga "Am Markt" ergänzen** · Priorität: Niedrig
  `services/slides/chapters/20-saga.md:50-58`: Chip `MicroProfile LRA` hinzufügen
  (passt als Java/Jakarta-Variante neben Axon).

- [x] **A5 · Setup-Schrift vergrößern** · Priorität: Mittel
  `services/slides/theme.css:510-513`: `.box.setup pre` von `0.85em` auf `1em` bis
  `1.2em`, Zeilenhöhe ggf. auf 1.6. Der Code wird abgetippt, soll gut lesbar sein.

- [x] **A6 · Kontrast-Pass (WCAG-AA)** · Priorität: Hoch
  `services/slides/theme.css`: Subtitle (`:110`), Links (`:129`), `.cols h3` (`:208`)
  und Inline-`code` (`:293-298`, `:621`) auf mindestens 4.5:1 anheben. Insbesondere
  `code`-Override für dunkle `.box` ergänzen (heute dunkles Lila auf dunklem Grund).
  Danach Palette-Sync mit dem Dashboard prüfen (Skill `sync-palette`,
  `services/CLAUDE.md` Farbpalette).
  Umgesetzt 2026-06-02: Inline-`code` auf dunklem Grund (`.box`, `.chapter-slide`)
  erhält helle Schrift (`#e0e6f0`), das war der eigentliche "blaue Schrift auf
  dunkel"-Fall (ca. 1.9:1). Der breitere Paletten-Kontrast (Subtitle/Links/`.cols h3`)
  bleibt bewusst offen: die Akzentfarbe `rgb(115,72,225)` erreicht auf Weiß bereits
  ca. 5.6:1, eine Änderung wäre eine Brand-Entscheidung samt Dashboard-Sync und kein
  reiner Quick-Win.

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

- [x] **C1 · `docs/troubleshooting.md`** · Priorität: Hoch
  Windows/WSL2: localhost und Traefik-Timing ("wenn nicht erreichbar, kurz warten
  und neu laden"), Docker-Desktop-Networking, Port-Konflikte (80, 8080, 8500),
  Image-Pull. GitHub-Zertifikat in WSL (Gast-WLAN, Proxy, Zertifikat hinterlegen).
  Verlinken aus `docs/vorbereitung.md`.
  Umgesetzt 2026-06-02: `docs/troubleshooting.md` symptom-orientiert angelegt, jeder
  Eintrag mit Marker Geltung (macOS/Linux/Windows vs. nur Windows/WSL2) und Status
  (verifiziert vs. berichtet). Der crossplattforme Kern (Erstdiagnose, Dashboard
  lädt nicht, Port-Konflikte inkl. 8084 und 8085-8091, Image-Pull, Reset) wurde
  lokal auf macOS gegen den laufenden Stack getestet. Die Windows/WSL2-Punkte
  (localhost trotz laufender Container, WSL-Zertifikat) sind aus dem Feedback
  abgeleitet, mit Links zu offiziellen Docs versehen und klar als "berichtet, nicht
  unter Windows reproduziert" markiert. Verlinkt aus `docs/vorbereitung.md` (Abschnitt 2
  und Windows-Hinweise). Offen: Verifikation der Windows-Punkte durch einen
  Windows-Teilnehmer.

- [ ] **C2 · `docs/admin-contract.md`** · Priorität: Hoch
  Exakter Dashboard-Vertrag mit JSON-Schemata, damit Custom-Services in jeder
  Sprache das Dashboard füllen: `GET /health`; ab Story 3 `GET /admin/circuit-state`;
  ab Story 4 `GET /admin/bulkhead-state` plus `POST /admin/bulkhead-reset`; ab
  Story 5 `GET /admin/sagas` plus `POST /admin/sagas-reset`. Felder aus
  `booking/story3/circuitbreaker/circuitbreaker.go`,
  `booking/story4/bulkhead/bulkhead.go`, `booking/story5/saga/saga.go` übernehmen.
  In den Stories und der Vorbereitung verlinken. Spannungsfeld dazu siehe D7
  (Aufwand vs. Lernfokus). Hinweis: das Bulkhead-Feld heißt nach A2 jetzt `inProgress`.

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

- [x] **D5 · Voraussetzungen schärfen, ggf. Vor-Workshop** · Priorität: Hoch
  In `docs/vorbereitung.md:5-10` als harte Vorbedingung formulieren: Teilnehmende
  müssen ein Greenfield-Projekt in ihrer Sprache aufsetzen können, das per Dockerfile
  ein Image baut und einen Health-Endpoint bereitstellt. Optionaler Vor-Workshop
  "Was ist OpenAPI? Wie geht Docker?". Setup-Dauer realistisch einplanen.
  Umgesetzt 2026-06-03: Pflicht-Vorbereitung als Abschnitt 0 in `docs/vorbereitung.md`
  verankert (Greenfield-Mini-Service mit Akzeptanzkriterien: `GET /health` liefert 200,
  Dockerfile, Image baut, Container lauscht auf 8080; dazu Stack-Smoke-Test und
  Firmennetz-Test gegen die WSL2-/Zertifikatsprobleme aus dem Kickoff). Voraussetzungen
  geschärft (Docker als hart markiert, Greenfield-Fähigkeit explizit). IntraHub-Text und
  E-Mail-Einladung als Vorlagen in `docs/orga/` angelegt, beide benennen die
  Pflicht-Vorbereitung (ca. 1 bis 2 h) und verlinken das Repo. `docs/themen.md` Block 4
  von "Setup (45 Min)" zu "Vorbereitung verifizieren (30 Min)" umgewidmet, Folien
  entsprechend angepasst (`services/slides/chapters/07-vorbereitung.md`, `03-agenda.md`).
  Bewusst NICHT umgesetzt: der separate Vor-Workshop bzw. Setup-Check-Termin
  (Entscheidung Lars). Stattdessen Troubleshooting-Verweis plus Kontaktangebot vorab.
  Ergänzung: Schwierigkeitsgrad im Intranet-Text von "Mittel" auf "Anspruchsvoll"
  angehoben (mündliches Feedback nach dem Kickoff: gut zu folgen, wenn man in
  Architektur schon bewandert ist), mit Einordnung, die Microservices-Vorwissen
  explizit ausnimmt, um die richtige Zielgruppe nicht abzuschrecken.

- [~] **D6 · Vorbereitungs- und Repo-Strategie** · Priorität: Mittel
  Material aktiv vorab teilen (Repo ist bereits öffentlich), Selbstlern-Charakter
  ausbauen, ggf. zwei Termine mit einer Woche Abstand zum echten Nacharbeiten.
  Teilweise umgesetzt 2026-06-03: Repo-Link aktiv in IntraHub-Text und
  E-Mail-Einladung platziert (`docs/orga/`), die Pflicht-Hausaufgabe (D5) stützt den
  Selbstlern-Charakter. Ergänzung: `docs/vorbereitung.md` und `docs/troubleshooting.md`
  sind jetzt als GitHub-Pages-Seite öffentlich teilbar
  (https://larmic.github.io/workshop_microservices/vorbereitung/, client-seitig
  gerendert aus den Markdown-Quellen), Einladung und IntraHub verlinken dorthin.
  Offen bleibt das Zwei-Termine-Modell mit einer Woche Abstand,
  bewusst zurückgestellt: das Format bleibt vorerst 2 Tage am Stück.

- [ ] **D7 · Aufwand der Admin-Endpoints vom Pattern-Lernen entkoppeln** · Priorität: Hoch · *Entscheidung Lars*
  Die Admin-Endpoints (siehe C2) sind wertvoll fürs Dashboard, ihre Implementierung
  kostet aber Zeit und vermischt zwei Ziele: das Pattern lernen vs. das Dashboard
  bedienen. Noch keine Lösung. Mögliche Richtungen (offen, noch nichts entschieden):
  - Admin- bzw. Instrumentierungs-Layer als fertige, einbindbare Vorlage je Sprache
    anbieten, damit Teilnehmende ihn nicht selbst bauen.
  - Admin-Endpoints als optional bzw. "auf der Tonspur" deklarieren (Dashboard zeigt
    den State dann nur für die Referenz-Implementierung live).
  - Den State minimal halten (z. B. nur `/health` plus ein generisches
    `/admin/state`), damit der Bau-Aufwand klein bleibt.
  Bezug: C2 (Vertrag), D1 und D3 (Lernfokus).

## E · Architektur- und Design-Entscheidungen · *Entscheidung Lars*

- [ ] **E1 · Story 5: Saga-Schema im API-Vertrag überdenken** · Priorität: Mittel
  Erfolg liefert `Booking`, Fehlerfall (503) liefert `SagaFailure` (= Saga)
  (`booking/story5/api/openapi.yaml:114-126`). Option: Fehlerfall ebenfalls als
  `Booking` mit Status liefern und das `Saga`-Schema aus dem öffentlichen Vertrag
  nehmen, oder die Trennung bewusst dokumentieren (Diskussions-Anker für die Story).

- [ ] **E2 · CB-Scope in Story 3** · Priorität: Mittel
  Prüfen, ob der Circuit Breaker in Story 3 auf GET beschränkt wird und die
  POST-Resilienz erst in Story 5 (Saga) sauber gelöst wird. Verknüpft mit B6.

## F · Vollständigkeit & roter Faden (aus der Kuratierung vom 2026-06-02)

Herleitung und Einordnung in [kuratierung.md](kuratierung.md).

- [ ] **F1 · Reduktionen-Übersicht "Was wir bewusst nicht behandeln"** · Priorität: Hoch
  Neue Folie plus Notiz in `docs/themen.md`: in-memory statt DB, HTTP-Events statt
  Broker (Story 6), selbstgebaut statt Resilience4j/OpenTelemetry, CQRS nur als
  Vortrag, je mit Ein-Satz-Begründung. Bündelt, was heute verstreut in
  `docs/instructions/*` steht, und beantwortet die wiederkehrenden
  "Warum nicht X?"-Fragen an einer Stelle.

- [ ] **F2 · Database-per-Service als bewusste Reduktion sichtbar machen** · Priorität: Mittel · *Entscheidung Lars*
  `docs/themen.md:82` behauptet, „1 DB pro Service" werde in Story 1 umgesetzt; das
  stimmt nicht (Story 1 hat keine DB, der Zustand ist überall in-memory:
  `services/flight/handler/flight.go`, `services/booking/story5/saga/saga.go`).
  Aussage korrigieren und das Prinzip motivieren (Schema-Evolution, differenzierte
  Skalierung, Transaktionsgrenzen), die Reduktion explizit kennzeichnen. Tiefe
  offen: nur Diskussions-Anker ("Was änderte sich mit je eigener DB: Saga-Persistenz,
  Recovery, Migrationen?") oder optionaler Bonus. Keine neue Pflicht-Story.
  Erweitert B8.

- [ ] **F3 · Inter-Story-Brücken explizit machen** · Priorität: Hoch
  Pro Story eine kurze "bisher / Problem jetzt / Lösung"-Brücke auf der
  Einstiegs-Folie und im Story-Header (`docs/stories/*`, die vorhandenen
  `## Kontext`-Abschnitte vereinheitlichen). Macht den vorhandenen roten Faden
  sichtbar. Ist die konkrete, buildbare Umsetzung von D1 und stützt B7 und D4.

- [ ] **F4 · Lernpfad-/Roter-Faden-Übersichtsfolie** · Priorität: Mittel
  Eine Übersichtsfolie mit der Problem-Pattern-Kette (Story 1 bis 7), die zeigt,
  welches Problem jede Story löst. Orientierung für Teilnehmende, besonders in der
  dichten Story-Folge an Tag 2 (`docs/themen.md`).
