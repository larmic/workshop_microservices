# Kuratierung Kickoff-Workshop · 2026-05-30

Nacharbeit vom 2026-06-02, nach Umsetzung der Quick-Wins A1 bis A6. Diese Datei
beantwortet vier Fragen zum Kurs: Sind die Stories vollständig? Bietet die
Themen-Reihenfolge einen roten Faden? Welche offenen Aufgaben lohnen sich als
Artefakt, welche sollten erst mit Lars oder im Team entschieden werden? Sie baut
auf [feedback.md](feedback.md) auf und mündet in einen neuen Backlog-Block F in
[backlog.md](backlog.md). Technische Punkte sind gegen die Codebasis verifiziert.

Status-Werte wie in der feedback.md: **bestätigt** · **teils** · **korrekt
umgesetzt** · **Einschätzung**.

## 1. Sind die Stories vollständig?

Kurzantwort: für ihr Lernziel ja. Das Set (Cloud-native Grundgerüst, Service
Discovery, Circuit Breaker, Bulkhead, Saga in zwei Varianten, Distributed Tracing)
ist in sich geschlossen, und jede Story baut auf der vorigen auf. Eine neue
Pflicht-Story ist nicht nötig.

Eine inhaltliche Lücke ist aber real und konkret benennbar:

- **"1 DB pro Service" wird behauptet, sogar als umgesetzt, ist aber nirgends
  vorhanden.** `docs/themen.md:82` sagt wörtlich, Health-Checks, externe
  Konfiguration und „1 DB pro Service" würden „in Story 1 konkret umgesetzt".
  Story 1 implementiert aber keine Datenbank
  (`docs/stories/story-01-cloud-native-booking-service.md`: nur Health, Config,
  Orchestrierung), und auch sonst hält die Referenz allen Zustand in-memory: die
  Domain-Services über hartkodierte Slices (`services/flight/handler/flight.go`,
  analog hotel/car), der Saga-Store als Map
  (`services/booking/story5/saga/saga.go`). *bestätigt.* → **F2**, verwandt mit **B8**
- Das ist eine **bewusste Reduktion** (Fokus auf Kommunikations- und
  Resilienz-Pattern, nicht auf Persistenz) und soll so bleiben. Falsch ist nur,
  dass `themen.md` sie als umgesetzt verkauft und das Prinzip nicht motiviert.
  Empfehlung: in-memory beibehalten, die Reduktion sichtbar machen, das Prinzip
  motivieren. Keine Persistenz-Story bauen.

Daneben gibt es Themen, die der Kurs **bewusst weglässt** und die nur als Ausblick
gehören, nicht als Lücke. Sie sind heute in den Trainer-Notes je Pattern verstreut
begründet, aber nirgends gebündelt, und genau diese Punkte lösen die
wiederkehrenden "Warum nicht X?"-Diskussionen aus:

- In-memory statt Datenbank (siehe oben)
- HTTP-POST-Events statt Message-Broker in Story 6 (`docs/questions/story6.md`
  begründet das bereits)
- Selbstgebaute Patterns statt Bibliotheken: Circuit Breaker, Bulkhead, Tracing
  statt Resilience4j oder OpenTelemetry (Lerneffekt vor Produktionsnähe, siehe
  `docs/instructions/*`)
- CQRS nur als Vortrag, ohne eigene Hands-on-Story

Klar außerhalb des Zwei-Tage-Rahmens und höchstens als Randnotiz: Security und
mTLS zwischen Services, API-Versionierung, Contract-Testing, Service Mesh, Event
Sourcing, Metrics/Prometheus, Deployment-Strategien als Hands-on. Diese **nicht**
als Vollständigkeits-Lücke führen.

→ Neuer Block-F-Vorschlag: **F1** (Reduktionen bündeln), **F2** (DB-Reduktion
sichtbar machen plus Diskussions-Anker).

## 2. Bietet die Themen-Reihenfolge einen roten Faden?

Ja, einen starken und problemgetriebenen. Jede Story löst ein Problem, das die
vorige aufwirft:

1. Story 1 → 2: Die Backend-URLs sind hartkodiert. Was bei mehreren Instanzen?
   → Service Discovery.
2. Story 2 → 3: Discovery findet Instanzen, aber eine hängt. → Circuit Breaker.
3. Story 3 → 4: Der CB fängt Fehler ab, aber langsame, fehlerfreie Calls saugen
   alle Ressourcen. → Bulkhead.
4. Story 4 → 5: Ein Bulkhead-Reject heißt Teilausfall. Wer räumt auf? → Saga.
5. Story 5 → 6: Synchrone Kompensation ist fragil. → Event-basierte Choreography.
6. Story 6 → CQRS: Wie fragt man eine über Services verteilte Buchung ab?
   → Lese-Modell.
7. → Story 7: Bei so vielen Services und Events, wie debuggt man?
   → Distributed Tracing.

Die Makro-Cluster ergeben eine saubere Dramaturgie: Grundlagen, Kommunikation,
Resilienz, Daten und Konsistenz (Story 5, 6, CQRS), Observability (Story 7),
Deployment und Kultur. **Eine Umsortierung ist nicht nötig.**

Die eigentliche Schwäche ist nicht die Reihenfolge, sondern ihre Sichtbarkeit:

- **Der rote Faden steckt im Inhalt, wird Teilnehmenden aber nicht explizit
  gezeigt.** Die Brücken oben stehen so in keiner Story und auf keiner Folie. Das
  ist exakt die buildbare Seite des Kernbefunds (Motivation pro Pattern, **D1**):
  ein Satz "bisher / Problem jetzt / Lösung" je Story macht das Warum sichtbar.
  *Einschätzung.* → **F3**, **F4**
- Die Story-`## Kontext`-Abschnitte leisten das ansatzweise (Story 5 erklärt etwa,
  warum klassische DB-Transaktionen über Service-Grenzen scheitern), aber
  uneinheitlich und nicht als durchgehende Kette.
- Bestätigt wird zudem das schon erfasste Rhythmus-Problem: Tag 2 reiht vier
  Stories hintereinander, bevor wieder Theorie kommt. Das schwächt nicht den
  Faden, aber die Aufnahme. → **D4**

## 3. und 4. Welche Aufgaben als Artefakt, welche als Diskussion?

Die Trennung ist im Backlog bereits angelegt: A, B und C sind umsetzbare
Artefakte, D und E sind mit *Entscheidung Lars* markiert. Wichtig ist die
Konsequenz daraus: **Die höchsthebeligen Punkte (D1, D2, B7) sind keine
Code-Aufgaben.** Das bestätigt den Kernbefund. Man kann viel "bauen" und trotzdem
das eigentliche Problem (Motivation, eigenes Denken) nicht lösen.

| ID | Kurz | Einordnung | Begründung |
|---|---|---|---|
| B1 | 12-Factor-Brücke | Artefakt jetzt | klarer Slide- oder Doc-Edit |
| B2 | Abgrenzung CB/Bulkhead/Rate-Limit | Artefakt jetzt | Inhalt existiert, muss auf die Folie |
| B3 | `rejected` erklären | Artefakt jetzt | Folie plus Code-Kommentar |
| B4 | Aufruf-Reihenfolge | Artefakt jetzt | Folie plus Code-Kommentar |
| B5 | Recap CB nicht in `/health` | Artefakt jetzt | Code stimmt, nur Folie fehlt |
| B6 | CB auf POST | Diskussion | inhaltliche Tiefe, verknüpft E2 |
| B7 | Monolith vs. Microservice je Pattern | Diskussion | *Entscheidung Lars*, Kernbefund |
| B8 | "1 DB pro Service" motivieren | Artefakt jetzt | Slide-Anker, Tiefe siehe F2 |
| C1 | `troubleshooting.md` | Artefakt jetzt | entlastet Teilnehmende stark |
| C2 | `admin-contract.md` | Artefakt jetzt | öffnet den Custom-Service-Pfad |
| C3 | Dashboard-Story-Links | Artefakt jetzt | kleiner Frontend-Hinweis |
| D1 | Motivation pro Pattern | Diskussion | Kernbefund, buildbare Seite = F3 |
| D2 | Nicht-codende DDD-Aufgabe | Diskussion | curriculares Format |
| D3 | Stories auf der Tonspur | Diskussion | welche, wie kennzeichnen |
| D4 | Theorie/Praxis-Balance | Diskussion | Tagesdramaturgie |
| D5 | Voraussetzungen, Vor-Workshop | Diskussion | organisatorisch |
| D6 | Repo- und Termin-Strategie | Diskussion | organisatorisch |
| D7 | Admin-Endpoints entkoppeln | Diskussion | *Entscheidung Lars*, offen |
| E1 | Saga-Schema im Vertrag | Diskussion | API-Design-Entscheidung |
| E2 | CB-Scope GET/POST | Diskussion | Design-Entscheidung |
| F1 | Reduktionen bündeln | Artefakt jetzt | Folie plus `themen.md`-Notiz |
| F2 | DB-Reduktion sichtbar plus Anker | Diskussion | *Entscheidung Lars* (Tiefe) |
| F3 | Inter-Story-Brücken | Artefakt jetzt | buildbare Seite von D1 |
| F4 | Lernpfad-Übersichtsfolie | Artefakt jetzt | Orientierung, vor allem Tag 2 |
| D8 | REST-vs.-RESTful-Einheit (Nachtrag) | entschieden, Ausgestaltung offen | Theorie-Input plus Flipchart-API-Design, erste nicht-codende Aufgabe (D2) |

## 5. Neue Punkte (Block F)

Aus den Befunden 1 und 2 leiten sich vier neue, im Backlog als Block F geführte
TODOs ab. Begründung hier, umsetzbare Form in [backlog.md](backlog.md).

- **F1 · Reduktionen bündeln.** Eine "Was dieser Workshop bewusst nicht
  behandelt"-Übersicht (Folie plus Notiz in `docs/themen.md`) beantwortet die
  wiederkehrenden "Warum nicht X?"-Fragen an einer Stelle, statt verstreut in den
  Trainer-Notes.
- **F2 · Database-per-Service.** Korrigiert die irreführende Aussage in
  `docs/themen.md:82` und motiviert das Prinzip, ohne eine Persistenz-Story zu
  bauen. Die Tiefe (nur Diskussions-Anker oder doch ein optionaler Bonus) ist
  *Entscheidung Lars*. Erweitert **B8**.
- **F3 · Inter-Story-Brücken.** Macht den roten Faden je Story sichtbar
  ("bisher / Problem jetzt / Lösung", auf Slide und Story-Header). Ist die
  konkrete, buildbare Umsetzung von **D1** und stützt **B7** und **D4**.
- **F4 · Lernpfad-Übersicht.** Eine Folie mit der Problem-Pattern-Kette (Story 1
  bis 7) zur Orientierung, besonders in der dichten Story-Folge an Tag 2.

## Empfehlung

1. **Stories sind für ihren Zweck vollständig.** Keine neue Pflicht-Story. Die
   einzige echte "behauptet, nicht gelebt"-Stelle ist Database-per-Service:
   in-memory beibehalten, die Reduktion sichtbar machen (F1, F2), das Prinzip
   motivieren (B8) und die falsche Aussage in `themen.md:82` korrigieren.
2. **Höchster buildbarer Hebel ist der rote Faden:** F3 plus F4 machen die
   vorhandene Dramaturgie für Teilnehmende sichtbar und operationalisieren den
   Kernbefund (D1), ohne auf die großen curricularen Entscheidungen zu warten.
3. **Jetzt als Artefakt umsetzen:** C1, C2, F1, F3, F4 sowie die Slide-Klärungen
   B1 bis B5. Alle risikoarm, hoher Nutzen.
4. **Vorher mit Lars oder im Team entscheiden:** D1 bis D7, E1, E2, B6, B7 und die
   Persistenz-Tiefe (F2). Hier liegt der größte Hebel, aber kein Code löst ihn.
5. **Sequenz:** die Artefakte zeitnah, parallel eine kurze
   Curriculum-Entscheidungsrunde für den D- und E-Block vor dem nächsten Termin.

## Nachtrag 2026-06-04

Nachgereichter Punkt: Viele Microservices sprechen über REST miteinander, die
RESTful-Prinzipien sind aber oft unbekannt oder werden falsch angewendet
(`GET /getUser?id=1423` statt `GET /users/1423`, Idempotenz von `DELETE`).
Daraus wird eine eigene Einheit: kurzer Theorie-Input "REST vs. RESTful" plus
eine Flipchart-Übung, in der Teams gemeinsam eine API für ein Problem aus der
Reise-Domäne designen. Die Einheit ist entschieden, offen ist die Ausgestaltung.

Das ergänzt die Vollständigkeits-Analyse aus Abschnitt 1: RESTful-API-Design
war dort nicht erfasst und ist abzugrenzen von der API-Versionierung, die
bewusst außerhalb des Zwei-Tage-Rahmens bleibt. RESTful ist Grundlagen-Stoff
(wie Services sauber miteinander sprechen), Versionierung ist Betriebs-Detail
und bleibt Randnotiz. Geführt wird der Punkt als **D8** im Backlog (Begründung
dort und in [feedback.md](feedback.md)). Inhaltlich ist es die erste konkrete
nicht-codende Aufgabe (Teil von **D2**), stützt die Theorie/Praxis-Balance an
Tag 1 (**D4**) und berührt die Idempotenz-Diskussion beim CB-POST (**B6**/**E2**).
