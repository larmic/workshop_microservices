# Feedback Kickoff-Workshop · 2026-05-30

Erster durchgeführter Kickoff. Diese Datei fasst das Feedback sachlich und
anonymisiert zusammen, ordnet jede Beobachtung ein und verweist auf die
zugehörige Backlog-ID in [backlog.md](backlog.md). Die technischen Punkte wurden
gegen die Codebasis verifiziert.

Status-Werte: **bestätigt** (gegen Code geprüft, trifft zu) · **teils** (trifft
mit Einschränkung zu) · **korrekt umgesetzt** (Code macht es bereits richtig, nur
die Vermittlung fehlt) · **Einschätzung** (didaktische oder organisatorische
Bewertung ohne Code-Bezug).

## Was gut lief

- Das **Dashboard** kam exzellent an. Es ermöglicht anspruchsvolle Demos
  (Circuit Breaker, Bulkhead, Saga) ohne viel Overhead und war aus Sicht der
  Teilnehmenden das Highlight.
- Der Gesamteindruck war positiv. Der Kurs hat inhaltlich getragen, das hier
  gesammelte Feedback betrifft die Schärfung, nicht die Substanz.

## Kernbefund

Das Feedback konvergiert auf einen Punkt, der wichtiger ist als jeder Einzel-Bug:
Der Workshop ist aktuell stark im Muster "Pattern erklären, dann nachbauen".
Dadurch kommen **eigenes Denken** und die **Motivation pro Pattern** ("warum ist
das relevant? gilt das nicht auch im Monolithen?") zu kurz, und das **Setup**
frisst Zeit, die der Theorie und Reflexion fehlt. Viele der folgenden Einzelpunkte
sind Symptome dieses Befunds. Der größte Hebel liegt daher in den Blöcken
*Didaktik & Curriculum* (B7, B8, D1 bis D6), nicht in den schnellen Code-Fixes.

## Didaktik & Curriculum

- **Eigenes Denken kommt zu kurz.** Es wurde eher ein Pattern erklärt und dann
  umgesetzt. Dadurch war nicht immer klar, warum das jeweilige Pattern relevant ist.
  *Einschätzung.* Größter Hebel des gesamten Feedbacks. → **D1**
- **Aufgaben ohne Coding fehlen.** Anregung aus einem DDD-Kurs: erst gemeinsam am
  Whiteboard ein Problem durchdenken (Schnitt, Bounded Contexts, Eventflüsse), dann
  implementieren. Heute ist die Domäne (Flug, Hotel, Auto) und der Schnitt fest
  vorgegeben. *bestätigt* (keine solchen Übungen vorhanden). → **D2**
- **Zu viele Stories, alle zum Implementieren.** Alle 7 Stories sind Pflicht mit
  Akzeptanzkriterien. Vorschlag: einige nur "auf der Tonspur" (durchsprechen statt
  bauen), damit für die implementierten die Motivation und Zeit bleibt.
  *bestätigt*. → **D3**
- **Tag 1 zu viel Programmieren, zu wenig Theorie.** Tag 2 reiht vier Stories
  hintereinander, bevor wieder Theorie kommt (`docs/themen.md:13-39`). Balance über
  beide Tage gewünscht. *bestätigt*. → **D4**
- **Anforderungen höher als kommuniziert.** Die Hürde war für einen Teil der Gruppe
  hoch. Klarer machen, was vorausgesetzt wird, ggf. ein Vor-Workshop
  ("Was ist OpenAPI?", "Wie geht Docker?"). *Einschätzung.* → **D5**
- **Repo vorher freigeben.** Anregung, das Material vorab zu teilen (senkt die
  Hürde, hilft Mitnahme, unterstützt den Selbstlern-Charakter), ggf. zwei Termine
  mit einer Woche Abstand. Das Repo ist bereits öffentlich; offen ist die aktive
  Vorab-Weitergabe und das Format. *Einschätzung.* → **D6**

## Setup & Voraussetzungen

- **Nicht alle hatten Docker / Docker Compose.** Heute harte Voraussetzung
  (`docs/vorbereitung.md:5-10`), aber ohne Vorab-Absicherung. *bestätigt*. → **D5**
- **Setup hat lange gedauert.** Es wurde unterschätzt, wie lange eine Gruppe
  braucht, um lokal ein Projekt in ihrer Sprache aufzusetzen. Vorschlag als harte
  Vorbedingung: "muss ein Greenfield-Projekt aufsetzen können, das per Dockerfile
  ein Image baut und einen Health-Endpoint bereitstellt". *Einschätzung.* → **D5**
- **Dashboard via localhost nicht erreichbar (nur Windows).** Alle Container liefen,
  Traefik zeigte die Route korrekt, trotzdem kein Zugriff. Erst nach mehrmaligem
  Starten erfolgreich. Die Traefik-Konfiguration selbst ist unauffällig
  (`services/traefik/dynamic.yml`), die Ursache liegt vermutlich im
  WSL2/Docker-Desktop-Networking bzw. im Timing der Route-Registrierung. Es gibt
  bisher kein Troubleshooting dazu (nur CRLF/WSL2-Hinweise in
  `vorbereitung.md:141-145`). *bestätigt* (Doku-Lücke). → **C1**
- **GitHub-Zertifikatsproblem über WSL.** In der WSL war kein Zertifikat hinterlegt,
  Workaround war das Gast-WLAN. Kein Hinweis in der Doku. *bestätigt*. → **C1**

## Slides & Kontrast

- **Blaue Schrift auf dunklem Untergrund schwer lesbar.** Inline-`code` hat
  `color: #4d27a8` ohne Override für dunkle Boxen (`theme.css:293-298`); in einer
  dunklen `.box` ergibt das dunkles Lila auf dunklem Grund. *bestätigt*. → **A6**
- **Kontraste generell zu niedrig.** Violett `rgb(115,72,225)` für Subtitle, Links
  und `.cols h3` (`theme.css:110, 129, 208`) liegt auf Weiß unter WCAG-AA.
  *bestätigt*. → **A6**
- **Vorbereitungs-Folie: Skripte deutlich größer.** Der abzutippende Setup-Code ist
  `font-size: 0.85em` (ca. 30px), trotz Kommentar "bewusst groß"
  (`theme.css:506-513`). *bestätigt*. → **A5**
- **"Am Markt" im Web nicht auffindbar.** Die Folien mit Framework-Unterstützung
  pro Pattern existieren in den Markdown-Dateien, sind im Web aber unsichtbar.
  Ursache: `.market-row` startet mit `opacity:0` und wird erst sichtbar, wenn das
  letzte Fragment (`<span class="show-all">`) durchgeklickt ist
  (`theme.css:670-678`, z. B. `chapters/20-saga.md:46-63`). Wer im Web nicht alle
  Fragmente durchklickt, sieht es nie. Wunsch: öffentlich zeigen. *bestätigt*. → **A3**

## Pattern-Inhalte (Slides & Trainer-Doku)

- **12-Factor wirkt cloud-native, nicht microservices-spezifisch.** Stimmt: die 12
  Faktoren sind breiter als Microservices. Der Brückensatz fehlt
  ("12-Factor ist notwendig, aber nicht hinreichend für Microservices";
  `chapters/06-12-faktor-app.md`, `docs/themen.md:72-82`). *bestätigt*. → **B1**
- **"Gilt das Pattern nicht auch im Monolithen?"** Kam wiederholt auf. Die Antwort
  fehlt pro Pattern auf der Folie. *teils* (in Trainer-Notes angedeutet, nicht in
  den Stories). → **B7**
- **"Warum nur 1 DB pro Service?"** Wird als Fakt genannt, aber nicht motiviert
  (`docs/themen.md:82`). *teils*. → **B8**
- **Viel Diskussion vor dem Bulkhead-Start.** Warum kein generisches Rate-Limit?
  Warum genügt der Circuit Breaker nicht? Warum liefert der Bulkhead keine Antwort
  wie der CB? Die Antworten existieren in `docs/questions/story4.md` und
  `docs/instructions/bulkhead.md`, aber **nicht auf den Slides**, daher die offene
  Diskussion. *bestätigt*. → **B2**
- **Exkurs Rate-Limit vor dem Bulkhead.** Anregung, vorab kurz zu erklären, was ein
  Rate-Limit ist und worin der Unterschied liegt. *Einschätzung.* → **B2**
- **"rejected" beim Bulkhead unklar.** Aus dem Codebeispiel wird nicht klar, wofür
  es steht. Tatsächlich: Fail-Fast, sofortige Ablehnung bei vollem Pool, Antwort
  503 mit Header `X-Bulkhead-Full` (`booking/story4/bulkhead/bulkhead.go:61-64`).
  *bestätigt* (Erklärung gehört auf die Folie). → **B3**
- **Aufruf-Reihenfolge relevant für die Demo.** Für den Bulkhead ist es wichtig, in
  welcher Reihenfolge die Services aufgerufen werden, damit das Dashboard Sinn
  ergibt. Bestätigt: die Aufrufe sind sequenziell Flight, Hotel, Car
  (`booking/story4/handler/booking.go:73-78`), das beeinflusst die Rejection-Kaskade.
  *bestätigt*. → **B4**
- **Circuit Breaker im Health-Endpoint (Recap Story 3).** Frage: ist das gut, und
  wer führt dann die Probe aus? Vermutete Antwort "nein" ist korrekt. Der Code macht
  es bereits richtig: `/health` ist nur Liveness (`shared/handler/health.go`), der
  CB-State liegt auf `/admin/circuit-state`. Die Probe käme z. B. von einer
  K8s-Readiness. *korrekt umgesetzt* (nur die Recap-Antwort gehört explizit auf die
  Folie). → **B5**
- **Circuit Breaker beim POST sinnvoll?** Antwort eher nein. Bestätigt: in Story 3
  umschließt der CB den POST-Pfad (`booking/story3/handler/booking.go`). Gehört in
  den Folien diskutiert (Idempotenz, Seiteneffekt im Half-Open). *bestätigt*. → **B6**, **E2**

## OpenAPI & Code

- **IDs als UUID deklariert, aber keine UUIDs geliefert.** Bestätigt und
  durchgängig: `Flight.id`, `Hotel.id`, `Car.id` sind in allen Booking-Specs
  (`booking/story1..7/api/openapi.yaml`) und in den Domain-Service-Specs als
  `format: uuid` deklariert. Die Handler erzeugen aber Prefix-IDs
  (`flight/handler/flight.go:161-165` liefert `F-…`, analog `H-…`, `C-…`, `B-…`),
  Testdaten sind `LH123`, `H1`, `C1`. Empfehlung: `format: uuid` entfernen, die
  sprechenden Prefix-IDs sind didaktisch sogar besser. *bestätigt*. → **A1**
- **`inFlight` kollidiert mit dem Flight-Service.** Der Bulkhead-State heißt
  `inFlight`, was bei der Diskussion über den Flight-Service verwirrt. Kommt im
  Go-Code (`booking/story4..7/bulkhead/bulkhead.go`), in den Slides
  (`chapters/17-bulkhead.md`, `17a-bulkhead-code.md`), im Dashboard und in den
  Trainer-Notes vor. Vorschlag `inProgress`. *bestätigt*. → **A2**
- **OpenAPIs teils fehlerhaft, als Story-Hinweis im Dashboard hinterlegen.** Die
  Spec-Links pro Service existieren bereits im Dashboard
  (`dashboard/static/index.html:1300-1310`), könnten pro Story prominenter als
  Hinweis verlinkt werden. *teils*. → **C3**
- **MicroProfile LRA fehlt bei Saga "Am Markt".** Die Liste enthält Temporal,
  Cadence, Camunda, Axon, Step Functions, Eventuate, MassTransit, NServiceBus, aber
  nicht MicroProfile LRA (`chapters/20-saga.md:50-58`). *bestätigt*. → **A4**
- **Story 5 liefert Saga-Objekte zurück.** Fachlich unlogisch, brauchen wir das
  Saga-Schema im Vertrag? Bestätigt mit Einschränkung: bei Erfolg (201) kommt ein
  `Booking` zurück, nur im Fehlerfall (503) ein `SagaFailure` (= Saga)
  (`booking/story5/api/openapi.yaml:114-126`). Offen: Fehlerfall auch als `Booking`
  mit Status liefern, oder die Trennung bewusst dokumentieren. *teils*. → **E1**

## Dashboard & Admin-Contract

- **Ohne Admin-Endpoints zeigt das Dashboard für Custom-Services nichts.** Bestätigt:
  Das Dashboard pollt `/admin/circuit-state`, `/admin/bulkhead-state`,
  `/admin/sagas` (plus Reset-POSTs) und `/health`
  (`dashboard/handler/{circuit,bulkhead,saga}.go`). Wer seinen Service in einer
  anderen Sprache baut und diese Endpoints nicht implementiert, sieht im Dashboard
  keinen State. Dieser Vertrag ist nirgends dokumentiert. *bestätigt*. → **C2**

## Abgleich-Checkliste (jeder Punkt einer Backlog-ID zugeordnet)

| Feedback-Punkt | Backlog-ID |
|---|---|
| 12-Factor eher cloud-native, MS-Bezug fehlt | B1 |
| Vorbereitungs-Folie: Skripte größer | A5 |
| Nicht alle hatten Docker/Compose | D5 |
| Anforderung höher, ggf. Vor-Workshop | D5 |
| Dashboard via localhost nicht erreichbar (Windows) | C1 |
| Vorbereitung: eigenes Projekt mit Health-Endpoint | D5 |
| OpenAPI-UUIDs falsch | A1 |
| Blaue Schrift auf dunkel schwer lesbar | A6 |
| Kontraste der Folien anpassen | A6 |
| Recap Story 3: CB in Health, wer probt? | B5 |
| Projekt-Setup generell Problem, laufender Webserver als Vorbedingung | D5 |
| Bulkhead `inFlight` zu `inProgress` | A2 |
| Bulkhead: warum kein Rate-Limit, warum CB nicht genug, warum keine Antwort | B2 |
| Bulkhead `rejected` unklar | B3 |
| Exkurs Rate-Limit vor Bulkhead | B2 |
| Aufruf-Reihenfolge relevant fürs Dashboard | B4 |
| Admin-Services fehlen, Dashboard zeigt nichts | C2 |
| OpenAPIs als Story-Hinweis im Dashboard | C3 |
| Dashboard kam exzellent an (positiv) | (Was gut lief) |
| CB beim POST sinnvoll? In Folien betrachten | B6, E2 |
| Story 5 liefert Saga-Objekte, Schema-Frage | E1 |
| "Am Markt" öffentlich im Web zeigen | A3 |
| Saga "Am Markt" plus MicroProfile LRA | A4 |
| Repo vorher freigeben, zwei Termine, Selbstlern-Charakter | D6 |
| Tag 1 zu viel Programmieren, zu wenig Theorie | D4 |
| GitHub-Zertifikat über WSL | C1 |
| Setup-Dauer unterschätzt | D5 |
| Zu viele Stories, manche auf der Tonspur | D3 |
| Motivation pro Pattern, eigenes Denken zu kurz | D1 |
| Warum nur 1 DB pro Service | B8 |
| Patterns auch im Monolithen? | B7 |
| Aufgaben ohne Coding (DDD-Whiteboard) | D2 |
