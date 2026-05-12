# Workshop-Themen

Diese Themensammlung bildet den praktischen Teil des Workshops. Die Aufgaben werden am Beispiel einer Reisebuchung (Hotel, Flug und Auto) durchgeführt.

## Grundlagen

| Thema | Beschreibung |
|-------|--------------|
| Twelve-Factor-App | Kurzer Überblick über die 12 Faktoren für Cloud-native Anwendungen |
| Health-Checks | Ein Must-Have für jede Microservice-Architektur |
| 1 DB pro Service | Datenbank-Isolation als Grundprinzip |
| External Configuration | Externalisierte Konfiguration für flexible Deployments |

## Resilience Patterns

| Thema | Beschreibung |
|-------|--------------|
| Circuit Breaker | Schutz vor Kaskadenfehlern |
| Bulkhead Pattern | Isolation von Fehlern durch Ressourcen-Trennung |
| Saga Pattern | Verteilte Transaktionen ohne 2-Phase-Commit |

## Kommunikation & Routing

| Thema | Beschreibung |
|-------|--------------|
| API-Gateway | Zentraler Einstiegspunkt für alle Client-Anfragen |
| Service Discovery / Service Registry | Dynamische Service-Lokalisierung |
| Backends for Frontends (BFF) | Spezialisierte Backend-Services pro Frontend-Typ |

## Daten & Events

| Thema | Beschreibung |
|-------|--------------|
| Eventsourcing / Event-Driven Architecture | Ereignisbasierte Architektur |
| CQRS | Command Query Responsibility Segregation |

### CQRS im Detail

CQRS (Command Query Responsibility Segregation) heißt schlicht: **Trennung von Command und Query** — Schreib- und Lesemodell werden unabhängig voneinander entworfen und optimiert. Mehr ist erstmal nicht dahinter.

**Wichtig:** CQRS ist **kein Synonym für Eventing**. Wie das Lesemodell aktualisiert wird, ist eine sekundäre Entscheidung — möglich sind Events, eine DB-Replica, ein Cache, ein nächtlicher Batch-Job oder eine materialisierte Sicht in der DB. Die Architektur-Idee ist unabhängig vom Transport.

#### Beispiel aus der Buchungsplattform

Der Booking-Service kennt aktuell nur den **Schreibpfad**: Eine Buchung kommt rein, die Saga läuft, am Ende ist das Ergebnis durch. Der Lesepfad — „zeige mir alle meine Reisen" — ist offen.

Der naive Ansatz, dafür bei jedem `GET` Flight, Hotel und Car parallel abzufragen und zu joinen, hat drei Probleme:

- **Langsam** — drei Backend-Calls pro Lese-Anfrage
- **Fragil** — ist ein Backend down, fehlen Daten oder die Anfrage scheitert ganz
- **Falsches Format** — die Backend-APIs liefern ihre internen Sichten, nicht die kundengerechte Historie

Die Lösung: **Lese- und Schreibmodell trennen**. Der Schreibpfad bleibt wie er ist (Saga). Daneben entsteht ein **Read-Model**, das auf eine völlig andere Anforderung optimiert ist: schnell antworten, denormalisiert, ausfalltolerant.

#### Mehrwert

- **Ausfalltoleranz beim Lesen** — `GET /customers/{id}/bookings` funktioniert auch wenn Flight, Hotel oder Car gerade down sind. Das Read-Model ist eigenständig, kein Live-Aufruf der Backends. Ein anderer Resilienz-Mechanismus als Circuit Breaker: Datenredundanz statt Fehler-Isolation.
- **Format-Entkopplung** — Die Sicht für den Kunden („meine Reisen") muss nicht der Sicht der Backends entsprechen. Read-Model ist denormalisiert und kundenzentriert.
- **Unabhängige Skalierung** — Auf Plattformen wie Check24 werden tausendfach mehr Versicherungen *gelesen* als *abgeschlossen*. Beide Lasten mit demselben Modell zu bedienen wäre unwirtschaftlich. CQRS macht das ungleiche Verhältnis architektonisch sichtbar — und erlaubt für jede Seite die passende Persistenz, Caching-Strategie und Skalierung.

#### Mechanismen für Read-Model-Updates

Die Wahl ist bewusst plural — keiner dieser Wege ist „CQRS":

- **Events** (z. B. die aus der Choreography-Saga): Bei `BookingCompleted` aktualisiert das Read-Model seine Sicht. Bequem, wenn Events ohnehin existieren.
- **DB-Replica mit asynchroner Replikation**: Read-Modell liegt auf einer separaten Datenbank, die der Schreibseite folgt.
- **Cache-Layer** vor der primären DB.
- **Materialisierte Sicht** in der DB selbst.
- **Batch-Job**, der das Read-Model periodisch neu aufbaut.

#### Trade-off: Eventual Consistency

Zwischen dem Schreibvorgang (Saga fertig) und dem Auftauchen im Read-Model liegen Millisekunden bis Sekunden. Im UI ggf. mit einem „wird in Kürze sichtbar"-Hinweis auffangen. Wer starke Konsistenz braucht, sollte CQRS nicht oder nur mit synchroner Aktualisierung einsetzen.

#### Diskussions-Anker

- Was würde sich ändern, wenn das Read-Model **nicht** aus Events, sondern aus einer DB-Replica gespeist würde? Was bliebe gleich?
- Wann lohnt CQRS *nicht*? (Hinweis: kleine Apps mit ähnlicher Lese- und Schreiblast.)
- Was passiert, wenn das Read-Model nach einem Crash leer ist — wie kommt der Stand zurück? (→ Event-Replay)
- Wo gehört das Read-Model fachlich hin: in den Booking-Service oder in einen eigenen `BookingHistoryService`?

## Deployment & Betrieb

| Thema | Beschreibung |
|-------|--------------|
| Downtimeless Deployment | Zero-Downtime Deployments |
| API First Ansatz | API-Design vor Implementierung |

## Kultur & Organisation

| Thema | Beschreibung |
|-------|--------------|
| Microservices ohne DevOps? | Begriffsklärung und warum DevOps essentiell ist |
| "You build it, you run it" | Verantwortung und Betriebsmodelle |
| Cloud Native | Ist Cloud Native eine Voraussetzung für Microservices? |
