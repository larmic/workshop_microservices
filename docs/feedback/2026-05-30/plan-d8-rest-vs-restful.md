# Plan D8 · Neue Story 2 "Design-Session: REST vs. RESTful" mit Umnummerierung 2..7 zu 3..8

Planungsstand: 2026-09-06, freigegeben (Lars). Ergänzt 2026-09-07: die
Diskussionsphase wird ein Quiz "RESTful oder nicht?" mit drei festen Beispielen
(Lars). Umsetzung noch nicht begonnen.
Backlog-Eintrag: [D8 in backlog.md](backlog.md). Herleitung in
[feedback.md](feedback.md) und [kuratierung.md](kuratierung.md). Pfade relativ
zur Repo-Wurzel. Die Story-Nummern in diesem Dokument sind die **neuen**
Nummern, sofern nicht anders gesagt.

## Kontext

Kickoff-Feedback (Nachtrag 2026-06-04): Viele Services sprechen REST, die
RESTful-Prinzipien sind aber unbekannt oder werden falsch angewendet
(`GET /getUser?id=1423`, unklare Idempotenz). Im Workshop gibt es dazu kein
Material, REST ist nur als "heutiger Default" genannt (`docs/themen.md:60`,
Slide `05-architektur-kommunikation.md` ist eine reine Bild-Slide). Entschieden
(D8): kurzer Theorie-Input plus nicht-codende Flipchart-Übung in Teams.

Entscheidungen vom 2026-09-06 (Lars):

- **Übung = Variante 1:** Teams designen am Flipchart eine API für Storno und
  Umbuchung aus der Reise-Domäne. Der Review der eigenen Workshop-API
  (RPC-artige Admin-Endpoints) ist nur ein Rückgriff im Recap der Bulkhead-Story.
- **Verortung = neue Story 2** zwischen der heutigen Story 1 und Story 2. Die
  heutigen Stories 2..7 werden zu 3..8, **vollständig inklusive Technik**
  (Verzeichnisse, Import-Pfade, Ports, Traefik, Makefile, Compose, Workflows,
  Docker-Tags, Dashboard, Slides, Docs). Story 2 ist die erste Story ohne Code.
- **Agenda-Autorität = Slides** (`services/slides/chapters/03-agenda.md`): Tag 1
  Vormittag endet mit Story 1, Nachmittag Story 2 (neu), Story 3 (Discovery),
  Story 4 (Circuit Breaker). `docs/themen.md` wird an diese Aufteilung
  angepasst (Tag 2 wird entlastet).
- **Name:** "Design-Session: REST vs. RESTful".
- **Hero-Slide** wie bei den Pattern-Blöcken, zunächst ohne Bild, Bild-Prompt in
  `services/slides/TODO.md`.

Warum genau zwischen Story 1 und 2: Story 1 endet mit `GET /booking/offers`,
die Discovery-Story führt mit `POST /booking/bookings` den ersten schreibenden
Endpoint ein (heute `docs/stories/story-02-service-discovery.md:40-49`). Die
Teams designen also unmittelbar vor dem ersten eigenen Schreib-Endpoint selbst
eine Ressource. Die Saga-Story bringt später `GET /booking/bookings/{id}` und
die Kompensation (Storno), dort zahlt die Storno-Diskussion ein zweites Mal
ein. Die Einheit löst zugleich einen Teil von **D2** (erste nicht-codende
Aufgabe) ein und stützt **D4** (Theorie/Praxis-Balance).

Verworfene Alternativen: (a) unnummerierte Einheit ohne Umbau, wie die
Vortragskapitel 05, 08, 29, 30 in den Slides; (b) nur didaktische
Umnummerierung bei unveränderten technischen Identifiern (hätte eine
Mapping-Tabelle "Story 3 liegt in `story2/`" gebraucht). Beides zugunsten
eines konsistenten Story-Begriffs verworfen.

## Teil A · Didaktik der neuen Story 2 (ca. 55 Min)

### Ablauf

| Phase | Zeit | Inhalt |
|---|---|---|
| Theorie-Input | 10 bis 12 Min | REST vs. RESTful anhand der Backlog-Beispiele: Ressourcen statt Verben, HTTP-Methode trägt Semantik, Idempotenz (GET/PUT/DELETE ja, POST nein, macht Retries sicher), GET ohne Seiteneffekte, Status-Codes statt 200 mit Fehler-Body, Sub-Ressourcen und Filter. HATEOAS/Richardson nur Ausblick |
| Aufgabe stellen | 3 Min | Szenario: "Eine Kundin will ihre Reise stornieren, eine andere möchte nur den Flug umbuchen." Teams zu 3 bis 4 Personen, Flipchart |
| Teamarbeit | 20 Min | Ergebnis als Endpoint-Tabelle: Methode, Pfad, Request-Kern, Antwort und Status-Code, idempotent ja/nein, Begründung |
| Vorstellung und Vergleich | 8 bis 10 Min | Je Team 2 Min. Trainer sammelt Varianten: Storno als `DELETE /bookings/{id}` vs. `POST /bookings/{id}/cancellation` vs. `PATCH` auf Status; Umbuchung als `PUT` vs. Sub-Ressource vs. Storno plus Neubuchung |
| Quiz "RESTful oder nicht?" | 8 bis 10 Min | Drei Endpoints nacheinander auf Slides (festgelegt Lars, 2026-09-07). Pro Beispiel: Endpoint zeigen, Handzeichen abfragen, eine Person aus der Minderheit begründen lassen, dann auflösen. Dramaturgie: eindeutiges Nein, trügerisches Ja, trügerisches Nein |
| Brücke | 2 Min | "In Story 3 baut ihr `POST /booking/bookings`. Nehmt die Regeln mit." |

Lösungsraum bewusst offen, der Trainer moderiert Trade-offs: Idempotenz bei
Retries, Teilausführung, Nachvollziehbarkeit des Storno-Grunds, doppeltes
`DELETE`, wann RPC-artige Endpoints legitim sind.

### Die drei Quiz-Beispiele (fest, Reihenfolge verbindlich)

**Quiz 1 · Eindeutig nicht RESTful**

```
GET /booking/cancelBooking?id=4711
```

Erwartung: fast alle sagen "nicht RESTful". Auflösung: Verb im Pfad, Identität
als Query-Parameter, GET mit Seiteneffekt. Anekdote aus dem Netz: 2005 hat der
Google Web Accelerator Links vorgeladen; bei 37signals Backpack waren
"Löschen"-Links einfache GET-Links, der Prefetcher hat Nutzern Daten gelöscht.
Seitdem ist "GET verändert nichts" Selbstschutz, keine Stilfrage. Besser:
`DELETE /booking/bookings/4711` oder Quiz 2.

**Quiz 2 · Sieht falsch aus, ist RESTful**

```
POST /booking/bookings/4711/cancellation
```

Erwartung: viele sagen "nein, da steht eine Aktion drin". Auflösung:
"cancellation" ist ein Substantiv, also eine Ressource. Der Client legt eine
Stornierung an, bekommt `201 Created`, kann sie später per GET nachlesen (Grund,
Zeitpunkt, Gebühr). Oft besser als nacktes DELETE, weil die Buchung als Historie
bleibt. Direkter Bezug zur Flipchart-Aufgabe, dort taucht diese Variante auf.
Anker aus dem Netz: Stripe modelliert Rückerstattungen als eigene Ressource
(`POST /v1/refunds`).

**Quiz 3 · Sieht richtig aus, ist nicht RESTful**

```
POST /booking/bookings
→ 200 OK
{ "status": "error", "message": "hotel not available" }
```

Erwartung: die meisten sagen "ja, sauber". Auflösung: Methode und Pfad stimmen,
der Status-Code lügt. HTTP sagt Erfolg, der Body sagt Fehler. Load Balancer,
Monitoring, Retries und der Circuit Breaker aus Story 4 sehen nur die 200 und
halten den Service für gesund. Richtig: `409 Conflict` oder
`422 Unprocessable Content` mit dem Fehler im Body. Beste Brücke in die
Resilience-Themen: Status-Codes sind Infrastruktur, keine Kosmetik.

Quellenhinweis für die Slides: Martin Fowler, "Richardson Maturity Model"
(zugänglichster Einstieg); Roy Fielding, "REST APIs must be hypertext-driven"
(2008, strengste Lesart). Für das Quiz reicht Fowler.

Reserve für die Sprecher-Notizen (falls Zeit bleibt oder ein Team genau das
gebaut hat):

- `POST /users/1423/delete` (Methode trägt keine Semantik)
- `POST /flights/search` mit komplexem Filter-Body (Grauzone: Suche als
  Ressource vs. `GET /flights?from=BRE&...` mit URL-Längen-Grenze)
- `PUT /bookings/1423` mit Teil-Objekt (Grauzone: PUT ersetzt ganz, PATCH
  ändert teilweise; was passiert mit weggelassenen Feldern?)
- `POST /admin/bulkhead-reset` aus der eigenen Referenz-Implementierung
  (bewusst RPC-artig fürs Dashboard, Beispiel für "Abweichen mit Grund")
- Optional ein Beispiel direkt aus den Flipcharts der Teams aufgreifen

### Artefakte

**A1 · Story-Dokument `docs/stories/story-02-api-design-session.md`**
Muster: `story-01-*.md` (Narrativ aus Auftraggeber-Sicht, Zeitrahmen, Thema,
Ziel, Aufgaben). Abweichung: Abschnitt "Aufgaben" beschreibt die
Flipchart-Arbeit und das Ergebnisformat, kein Code, kein Akzeptanzkriterium mit
Endpoint. Ergänzt einen Hinweis "Diese Story hat keine Referenz-Implementierung
und keinen Port".

**A2 · Trainer-Hinweis `docs/instructions/rest-vs-restful.md`**
Aufbau analog `docs/instructions/circuitbreaker.md` (nummerierte Abschnitte):

1. Worum geht es (Kernbefund Kickoff)
2. REST vs. RESTful: die sechs Prinzipien mit Beispielen
3. Ablauf der Design-Session (Tabelle oben, Teamgröße, Material)
4. Szenario und Leitfragen für die Teams
5. Lösungsraum: typische Varianten für Storno und Umbuchung mit Trade-offs
6. Häufige Fehler und wie man sie am Flipchart anspricht
7. Quiz "RESTful oder nicht?": die drei Beispiele mit erwarteter Verteilung
   im Raum, Auflösung, Anekdote (Google Web Accelerator, Stripe) und
   Moderationshinweisen (Handzeichen, Minderheit begründen lassen, dann
   auflösen); Reserve-Beispiele mit Grauzonen-Argumenten
8. Diskussionsfragen (Idempotenz bei Retries, Bezug zu CB auf POST in der
   CB-Story, Kompensation in der Saga-Story)
9. Bezug zum Workshop-Code: wo die Referenz bewusst von RESTful abweicht
   (`POST /admin/bulkhead-reset`, `POST /admin/sagas-reset`, `POST /admin/chaos`),
   als Reserve-Beispiel im Quiz und als Rückgriff im Recap der Bulkhead-Story

**A3 · Fragen-Datei `docs/questions/story2.md`**
Kurz, im Muster der anderen Fragen-Dateien (Frage, Antwort, Spicy). 3 bis 4
Fragen aus Abschnitt 7 des Trainer-Hinweises. Die heutigen `story2..6.md`
werden vorher zu `story3..7.md` verschoben (Teil B).

**A4 · Slides (`services/slides/chapters/`)**
Neuer vertikaler Stack nach `10-story-01-fragen.md`. Dateinamen mit
Buchstaben-Suffix, damit die bestehenden Präfixe 11..30 nicht angefasst werden
(Präzedenz: `06a`, `08a`, `13a`):

- `10a-titel.md` · Hero "REST vs. RESTful" (Muster `11-titel.md`, ohne Bild)
- `10b-rest-vs-restful.md` · Prinzipien als `.factor-row` mit 6 `fragment`-Karten,
  je Gegenbeispiel/Beispiel und `<aside class="notes">`
- `10c-rest-vs-restful-beispiele.md` · Tabelle "so nicht / so" (Muster
  `14a-circuit-breaker-code.md`)
- `10d-story-02.md` · Story-Card (Muster `09-story-01.md`): Kontext, User Story,
  rechts "Aufgabe und Ergebnisformat", `time-badge` "≈ 55 min"
- `10e-story-02-fragen.md` · Recap mit Design-Varianten als `.recap-grid`
  (Muster `10-story-01-fragen.md`)
- `10f-quiz-1.md`, `10g-quiz-2.md`, `10h-quiz-3.md` · Quiz-Slides, je ein
  Beispiel groß in einer `.box` (bei Quiz 3 inklusive Response-Zeile), Titel
  "Quiz 1/3" bis "Quiz 3/3", Subtitle "RESTful oder nicht?". Die Auflösung
  (Urteil, ein Satz Begründung, bessere Variante) als `fragment`, damit sie
  erst nach den Handzeichen erscheint. Sprecher-Notizen mit erwarteter
  Verteilung, Anekdote (Quiz 1: Google Web Accelerator 2005, Quiz 2: Stripe
  Refunds, Quiz 3: Brücke zum Circuit Breaker in Story 4) und den
  Reserve-Beispielen. Die drei Beispiele stehen oben unter "Die drei
  Quiz-Beispiele".

Einbindung in `services/slides/index.html` nach Zeile 40 als verschachteltes
`<section>` (Titel, Theorie, Beispiele) gefolgt von der Story-Card, den
Quiz-Slides als eigenem vertikalen Stack und der Fragen-Slide, analog
zum Discovery-Block in Zeilen 41-48.

`03-agenda.md`: Tag 1 Nachmittag wird "Story 2 · Design-Session", "Story 3 ·
Service Discovery", "Story 4 · Circuit Breaker"; Tag 2 entsprechend 5..8.
`08-monolithen-mit-netzwerkproblemen.md` Notes: Story-Verweise anpassen.
`05-architektur-kommunikation.md` Notes: ein Satz Vorverweis auf Story 2.
`services/slides/TODO.md`: Bild-Prompt für die Hero-Slide ergänzen (Stil und
Palette wie die vorhandenen Prompts), offener Punkt `story7.md` wird zu
`story8.md`.

Konventionen: Umlaute als HTML-Entities, Notes per `Note:` oder
`<aside class="notes">`, kein Em-Dash im Fließtext.

**A5 · Moderationsleitfaden `docs/themen.md`**

- Zeitleiste Tag 1 an die Slides-Agenda anpassen: Vormittag Blöcke 1 bis 5 plus
  Story 1 (90 Min), Nachmittag Story 2 (55 Min), Story 3 (60), Story 4 (60).
  Tag 2: Story 5, 6, 7, CQRS, Story 8, BFF, Deployment, Abschluss. Summen
  nachrechnen und im Umsetzungs-Vermerk notieren.
- Block 2 (`:60`): Vorverweis "vertieft in Story 2".
- Abschnitt 6 "Hands-on: Stories 1–6" wird "Stories 1–7", neuer Unterabschnitt
  6b "Story 2 · Design-Session: REST vs. RESTful (55 Min)" mit Lernpointe,
  Recap-Frage und Verweis auf den Trainer-Hinweis. Folgende Unterabschnitte
  6c..6g, Abschnitt 8 wird "Story 8".
- Überleitung in Block 5 (`:127`) und Recap-Frage in Abschnitt 8 (`:213`,
  "Stories 3–6") auf neue Nummern anpassen.

**A6 · Brücken in Nachbar-Stories**

- `docs/stories/story-03-service-discovery.md` (ex 02), Aufgabe "Buchung
  durchführen": Satz "Design-Regeln aus Story 2 anwenden".
- Slide-Notes der Discovery-Story-Card: gleiche Brücke.
- Saga-Story (Docs und Slide-Notes): Rückverweis auf die Storno-Diskussion.

**A7 · Dashboard-Inhalt für Story 2**
Story-Section nur mit `story-info` (Kontext, User Story, Aufgabe und
Ergebnisformat), ohne Cheatsheet, Setup, Chaos oder API-Links. Ein Hinweis-Satz
"Diese Story findet am Flipchart statt". Technische Details in Teil B.

**A8 · Backlog und Feedback-Ordner**

- `backlog.md`: D8 auf `[x]` mit Umsetzungsnotiz (Story 2, Variante 1,
  Umnummerierung), D2 auf `[~]` ("erste nicht-codende Aufgabe umgesetzt,
  DDD-/Schnitt-Warm-up offen"), D4 Notiz, Agenda-Angleichung Slides/themen.md
  als erledigt vermerken.
- `docs/feedback/README.md`: Konvention um optionale `plan-<id>-*.md` pro
  Backlog-Punkt ergänzen (erledigt mit diesem Dokument).
- Historische Feedback-Texte (`feedback.md`, `kuratierung.md`, erledigte
  Backlog-Einträge) behalten die alten Story-Nummern, sie dokumentieren den
  Stand des Kickoffs. Am Kopf von `backlog.md` ein Hinweis auf die
  Umnummerierung vom 2026-09.

## Teil B · Technische Umnummerierung 2..7 zu 3..8

Ergebnis: Code-Stories bleiben sieben (1, 3 bis 8), Story 2 hat kein
Verzeichnis, keinen Port, kein Image. Port-Formel `8084+N` bleibt, 8086 bleibt
frei. Der ganze Umbau geht in **einen PR**, weil `deploy-pages.yml` die Slides
bei jedem Merge sofort veröffentlicht. Offene Renovate-PRs gegen `build.yml`
und `docker.yml` vorher mergen oder rebasen.

### B1 · Umbenennungen (absteigend, `git mv`)

Reihenfolge 7→8, 6→7, … 2→3, sonst kollidieren Zielnamen:

- `services/booking/story{N}` → `story{N+1}`
- `docs/stories/story-0{N}-*.md` → `story-0{N+1}-*.md`
- `docs/questions/story{2..6}.md` → `story{3..7}.md` (`story8.md` fehlt danach
  weiterhin, war schon als `story7.md` offen; Eintrag in `slides/TODO.md`)
- `services/slides/chapters/*-story-0{N}*.md` → `*-story-0{N+1}*.md`, die
  Zahlenpräfixe `12-`, `13-`, … bleiben (sind Sequenz, keine Story-Nummer)

### B2 · Textersetzung in einem Durchlauf

Ein Perl-Skript (außerhalb des Repos) mit Callback (`N ≥ 2 → N+1`), damit keine
Kaskade "6→7→8" entsteht. Dateiliste aus `git ls-files` minus Ausschlüsse:
`docs/feedback/`, `services/.env.example`, `services/booking/custom/`,
Binärdateien (`*.jar *.png *.svg *.woff2`), `go.sum`, `.DS_Store`.

Ersetzungsmuster (jeweils mit `(?![0-9])`-Schutz, obwohl es heute keine
zweistelligen Story-Nummern gibt):

- `story[2-7]` in Pfaden, Identifiern, Klassen (`booking-ref-story7`,
  `booking/story7`, `story7.md`, `.story7-hint` inkl. `retro.css:205`)
- `Story[2-7]` (Go-Variablen `bookingRefStory7URL`), `STORY[2-7]_URL`
- `Story [2-7]` in Prosa (Rest siehe Manuell-Liste)
- `story-0[2-7]` in Links und `index.html`
- Dashboard-Attribute `data-story`, `data-story-node`, `data-story-helpers`,
  `data-source-toggle`, Aufrufe `showStory(N)`, `setBookingMode(N,…)`,
  `bookingUrl(N,…)`, IDs `sN-offers|booking|result`
- Consul-Name `"booking-N"` in `main.go` und `example: booking-N` in den
  OpenAPI-Specs (story6/7 enthalten heute inkonsistent `booking-5`, bei der
  Gelegenheit korrigieren), `version: N.0.0`, `servers.url …booking-ref-storyN`
- Compose-Ports `"8086:8080"` bis `"8091:8080"` → jeweils +1 (nur diese Form,
  Basis-Ports 8080 bis 8085 und 8099 nicht anfassen)
- Danach `gofmt -l services`

Zusätzliche Fundorte über die erste Erkundung hinaus: `services/load-balancing.md`,
`services/booking/story6|7/README.md`, `services/shared/tracing/tracing.go:3`,
`services/hotel|flight|car/api/openapi.yaml` ("Story 5/6" in Beschreibungen),
`.claude/agents/didactic-reviewer.md` (teils veraltete Beispiele),
`.claude/agents/docker-compose-validator.md` (veraltete Ports 18080 ff.),
`.claude/skills/openapi-handler-check/SKILL.md`.

### B3 · Manuell-Liste (nach dem Skriptlauf gegenlesen)

Das Skript verschiebt in Bereichen nur die erste Zahl ("Story 5/6" würde
"Story 6/6"). Review-Grep vor und nach dem Lauf:

```
grep -rnE 'Stor(y|ies) [0-9] ?(/|–|-|&ndash;|bis|und|,) ?[0-9]|Stories [0-9]|sieben|alle 7|8085 bis|1808[0-9]' --exclude-dir=.git --exclude-dir=feedback .
```

Bekannte Stellen (heutige Zeilennummern): `docs/themen.md:127,131,213`;
`docs/stories/story-07-tracing.md:14,27,67,124`;
`docs/stories/story-03-circuit-breaker.md:65` (`localhost:8087` → `8088`);
`docs/questions/story1.md:171`, `story2.md:122,277`, `story5.md:322`;
`docs/instructions/distributed-tracing.md:88`; Slides `08-monolithen…:19`,
`09-story-01.md:52`, `13a:10,13`, `19a:10`, `26-titel`/`26-tracing` ("Story 5
und 6"), `27-story-07:42`, `29-zusammenfassung.md` ("sieben" → "acht",
Tabelle); `services/dashboard/static/index.html:1923` ("Story 5/6");
`docs/troubleshooting.md:82` (Ports `8085 bis 8091` → `8085, 8087 bis 8092,
8086 frei`); `docs/vorbereitung.md:15` ("alle 7 Stories" → 8);
`README.md:42,51`, `CLAUDE.md:31,45`, `services/README.md:21,31-32,42,61`
(Hinweis "story2 hat keinen Code").

`services/.env.example` ist für Claude nicht lesbar (Deny-Regel für `.env.*`).
Lars prüft selbst mit `grep -n story services/.env.example` (erwartet: nur
`CUSTOM_BOOKING_PATH`, keine Story-Ports).

### B4 · Dashboard (`services/dashboard/`)

- `main.go`: kein Eintrag für Story 2. Nur Shift 2..7 → 3..8 in Env-Namen,
  Variablen, Map `bookingURLs` und Proxy-Routen. `health_overview.go:66-70`
  iteriert die Map, Story 2 erscheint damit automatisch weder im
  Startup-Overlay noch im Health-Panel. Optional, separater Commit: `saga6-*`
  Endpunkte und JS (`fetchSaga6State`, `resetSagas6`, `triggerSaga6`, IDs
  `saga6-*`) zu `saga7-*`, sonst bleibt die 6 als Fossil.
- `static/index.html`:
  - Stepper (`:1358-1372`): acht Buttons, Button 2 mit Modifier
    `class="stepper-node design"` und eigenem CSS (gestrichelter Rahmen) neben
    `.stepper-node` (`:818-840`), damit sichtbar ist, dass nichts deployt wird.
  - Neue `<section class="story-section" data-story="2" hidden>` nach Story 1
    (endet `:1467`): nur `story-helpers` mit `details.story-info` (Kontext,
    User Story, Aufgabe und Ergebnisformat), optional zweites `details`
    "Ablauf und Zeitbox". Kein `source-toggle`, kein Cheatsheet, kein Setup,
    keine Chaos- oder Reset-Buttons.
  - `STORY_META` (`:3002-3017`): Eintrag 2 neu, Rest shiften. `STORY_COUNT = 8`.
  - `STORY_KEY` (`:3019`) auf `workshop-dashboard.currentStory.v2`, sonst landen
    wiederkehrende Browser mit gespeichertem `3` auf der falschen Story.
  - `BookingMode`-Literal (`:2073`): Keys 1, 3..8. Gates (`:2091-2094`):
    `storyNum === 4/5/6/7`.
  - Link-Chips "Service APIs" (`:1315-1321`) und Swagger-URLS in
    `docker-compose.infra.yml:40-46`: kein Eintrag für Story 2.
  - Einstellige Regexe `booking-ref-story(\d)` (`:2152`) und `story\d`
    (`:2079,2081`) reichen für 8.

### B5 · Build, Compose, Traefik, CI, Docker Hub

- `services/Makefile` (`:1,11,25-41`), `docker-compose.reference.yml` (sieben
  Services), `traefik/dynamic.yml` (Router `:35-81`, Middlewares `:152-180`,
  Services `:218-246`): mechanischer Shift.
- `.github/workflows/build.yml:89-189`: Jobnamen und Pfade shiften.
  `docker.yml`: Filter-Pfade `:53-58` und Matrix `:147` →
  `[1, 3, 4, 5, 6, 7, 8]`. `dockerhub-description.yml:28` Short-Description
  "story1, story3..story8, custom". `.github/dockerhub/booking.md:16-21`
  Tag-Tabelle plus Absatz "Nummerierung seit 2026-09: Story 2 ist eine
  Design-Session ohne Image; ältere Anleitungen meinen mit `storyN` den heutigen
  Tag `story(N+1)`".
- Docker-Hub-Tags: `story3..story7` werden beim ersten Push mit neuem Inhalt
  überschrieben (Tag-Semantik ändert sich). Nur `story2` bliebe verwaist mit
  altem Consul-Stand. Empfehlung: `story2` nach dem ersten grünen Lauf auf
  Docker Hub löschen (Hub-UI), Historie steht in `booking.md`.
- `docs/troubleshooting.md`: Punkt "nach Repo-Update `docker compose … pull`
  bzw. `make docker-up-hub`", weil lokal gecachte `story3..7`-Images sonst den
  alten Inhalt tragen.
- `.claude/skills/new-booking-story/SKILL.md` komplett synchronisieren: heute
  `booking-storyN` statt `booking-ref-storyN`, `BOOKING_STORYN_URL` statt
  `BOOKING_REF_STORYN_URL`, beschreibt ein nicht existierendes `http/`-Verzeichnis,
  alte Actions-Versionen. Port-Regel ergänzen "8084+N, 8086 unbelegt".
- `docs/feedback/README.md`: Absatz "Ordner vor 2026-09 verwenden die alte
  Nummerierung (alt N entspricht neu N+1 für N ≥ 2)". `docs/feedback/2026-05-30/*`
  selbst bleibt unverändert.

### B6 · Aufwand und Risiken

| Block | Aufwand |
|---|---|
| Skript und Manuell-Liste vorbereiten, Trockenlauf | 1 h |
| `git mv`, Perl-Lauf, gofmt | 0,5 h |
| Manuelle Prosa (docs, slides, themen, README/CLAUDE, .claude, SKILL-Sync) | 2 h |
| Dashboard (Section, META/COUNT/KEY, BookingMode, Gates, CSS, optional saga6→7) | 1 bis 1,5 h |
| CI und Docker Hub (Matrix, Filter, Beschreibung, Tag-Löschung nach Lauf) | 0,5 h plus Workflow-Laufzeit |
| Verifikation | 1,5 h |
| Teil A Didaktik (Story-Doc, Trainer-Hinweis, Fragen, 8 Slides inkl. drei Quiz-Slides, themen.md, Dashboard-Text) | 3,5 bis 4,5 h |
| **Summe** | **10,5 bis 12,5 h** |

Risiken: Doppelverschiebungen in Prosa (durch Single-Pass plus Manuell-Liste
abgefangen, Review-Pflicht); Tag-Semantik auf Docker Hub ändert sich für
Teilnehmende mit alten lokalen Images oder Forks; Pages-Deploy bei Merge,
darum ein PR; Renovate-Kollisionen in den Workflows; `services/.env.example`
und IDE-Run-Configs außerhalb des Skripts; `docs/questions/story8.md` fehlt
weiterhin.

### Commit-Reihenfolge im PR

1. Umbenennungen und mechanischer Shift (B1, B2) plus gofmt, Build grün
2. Manuell-Liste und Dashboard-Technik (B3, B4, B5 ohne Docker-Hub-Löschung)
3. Optional `saga6` → `saga7`
4. Teil A: neue Story-2-Inhalte (Docs, Slides, Dashboard-Section, themen.md)
5. Backlog, Plan-Vermerk, README-Hinweise (A8)

## Verifikation (Reihenfolge, alles vom Repo-Root, kein `cd`)

1. `git status --short`: Renames als `R`, sonst nur erwartete Änderungen.
2. `go build -C services ./...`, `go vet -C services ./...`,
   `go test -C services ./...`, `gofmt -l services` leer.
3. Compose: `docker compose --project-directory services -f services/docker-compose.yml
   -f services/docker-compose.infra.yml -f services/docker-compose.reference.yml
   -f services/docker-compose.custom.yml config --quiet`, dazu Port-Duplikate per
   `config --format json | jq` prüfen (8086 darf fehlen). Agent
   `docker-compose-validator` laufen lassen.
4. Traefik: `grep -oE 'booking-ref-story[0-9]' services/traefik/dynamic.yml | sort | uniq -c`,
   jeder Name genau sechsmal, alle Namen in Compose vorhanden.
5. Slides: alle `data-markdown`-Pfade aus `index.html` existieren
   (`test -f`-Schleife), dann `python3 -m http.server -d services/slides 8000`:
   Agenda, Kapitel 10a bis 10h, Quiz-Fragments erscheinen erst auf Klick,
   Story-Karten 2 bis 8, Zusammenfassungstabelle, Notes (Taste S), keine 404
   in der Konsole.
6. Dashboard (`make -C services docker-up`): Startup-Overlay listet
   `booking-ref-story1,3..8`; Stepper 1 bis 8; Story 2 zeigt nur die Karte;
   CB-, Bulkhead-, Saga-Panels pollen Story 4/5/6/7; Tracing-Button (8) liefert
   `trace_id`; Custom-Toggle funktioniert; localStorage-Key ist `.v2`.
7. Workflow-YAML parsen (`actionlint` oder `python3 -c 'import yaml…'`),
   Matrix und Jobnamen sichtprüfen.
8. Restsuche muss leer sein (außer bewusst neuen Story-2-Dateien):
   `grep -rnE 'booking-ref-story2\b|booking/story2\b|STORY2_URL|"8086:8080"|Story7URL|story-07-|questions/story7|saga6' --exclude-dir=.git --exclude-dir=feedback .`
   plus die Manuell-Liste aus B3 erneut.
9. Skill `openapi-handler-check`, Agent `didactic-reviewer` (Story-Kommentare
   im Go-Code), Zeitleiste Tag 1 und Tag 2 nachrechnen.
10. Nach Merge: `build.yml`, `docker.yml`, `deploy-pages.yml` grün, Docker Hub
    zeigt `story1`, `story3..story8`, dann `story2` löschen, Pages-Slides
    gegenlesen.
