# Workshop-Skript

Moderationsleitfaden für den Microservices-Workshop am Beispiel der Reisebuchungsplattform (Flug, Hotel, Auto). Knappe Stichpunkte, Diskussions-Anker und Verweise auf die Hands-on-Stories.

Ergänzende Dokumente: [Vorbereitung](vorbereitung.md) · [Idea-Sammlung](idea.md) · [Fragen](questions)

---

## 0. Überblick & Zeitleiste

Grobe Aufteilung über zwei Tage. Zeiten sind Richtwerte — bitte ans Tempo der Gruppe anpassen.

### Tag 1 (≈ 8 h inkl. Pausen)

| Block | Thema | Zeit |
|-------|-------|------|
| 1 | Ankommen & Motivation | 30 Min |
| 2 | Vortrag: Architektur & Kommunikation (Geschichte) | 45 Min |
| ☕ | Kaffeepause | 15 Min |
| 3 | Vortrag: 12 Factor App | 30 Min |
| 4 | Setup der Entwicklungsumgebung | 45 Min |
| 5 | Diskussion: „MS = Monolithen mit Netzwerkproblemen" | 30 Min |
| 🍽 | Mittagspause | 60 Min |
| 6a | Hands-on: Story 1 (Cloud-native Setup) | 90 Min |
| ☕ | Kaffeepause | 15 Min |
| 6b | Hands-on: Story 2 (Service Discovery) | 60 Min |

### Tag 2 (≈ 8 h inkl. Pausen)

| Block | Thema | Zeit |
|-------|-------|------|
| 6c | Hands-on: Story 3 (Circuit Breaker) | 60 Min |
| 6d | Hands-on: Story 4 (Bulkhead) | 60 Min |
| ☕ | Kaffeepause | 15 Min |
| 6e | Hands-on: Story 5 (Saga, Orchestration) | 60 Min |
| 🍽 | Mittagspause | 60 Min |
| 6f | Hands-on: Story 6 (Choreography-Saga) | 60 Min |
| 7 | Vortrag & Diskussion: CQRS | 30 Min |
| 8 | Hands-on: Story 7 (Distributed Tracing) | 60 Min |
| ☕ | Kaffeepause | 15 Min |
| 9 | Diskussion: BFF | 20 Min |
| 10 | Diskussion: Downtimeless Deployment | 20 Min |
| 11 | Abschluss & Kulturwandel | 30 Min |

---

## 1. Ankommen & Motivation (30 Min)

- Begrüßung, kurze Vorstellungsrunde
- Erwartungen der Teilnehmer einsammeln (Flipchart / Miro)
- Leitfrage: **„Brauchen wir Microservices überhaupt?"**
  - Erst Monolith, MS nur wenn wirklich sinnvoll
  - Workshop zeigt, *wie* man sie macht — nicht, *dass* man sie braucht
- Verweis: [idea.md](idea.md) — „Brauchen wir Microservices?"

---

## 2. Vortrag: Architektur & Kommunikation (45 Min)

Die Geschichte der verteilten Systeme als roter Faden — kein technisches Deep-Dive, sondern Einordnung.

- **Monolith** — alles in einem Prozess, einer DB, einem Deployment
- **SOA** — unternehmensweite Service-Orientierung, oft mit ESB
- **SOAP / WS-\*** — XML, WSDL, schwergewichtig
- **REST + JSON** — leichtgewichtig, ressourcenorientiert, der heutige Default
- **Microservices** — kleine, unabhängig deploybare Services pro Bounded Context
- **Modulith** — modularisierter Monolith als pragmatischer Mittelweg
- **SCS (Self-Contained Systems)** — vertikale Schnitte inkl. UI, gröberer Schnitt als MS
- **Rolle von Docker** — warum MS in der Breite erst durch Container praktikabel wurden

**Diskussions-Anker:**
- Was unterscheidet SOA von Microservices wirklich? (Hinweis: Scope — unternehmensweit vs. Team-Implementierung)
- Wann SCS, wann MS? Wann reicht ein Modulith?

---

## 3. Vortrag: Bedingungen an einen Microservice — 12 Factor App (30 Min)

- 12 Faktoren kurz durchgehen, jeweils mit Praxisbezug zum Booking-Service
- Schwerpunkte, die im Workshop praktisch werden:
  - **Codebase** — ein Service, ein Repo
  - **Config** — Externe Konfiguration (kommt in Story 1)
  - **Backing Services** — DB, Consul als angehängte Ressourcen
  - **Build, Release, Run** — saubere Trennung
  - **Disposability** — schneller Start, sauberer Shutdown (Health-Checks)
  - **Dev/Prod Parity** — Docker hilft hier massiv
- Hinweis: Health-Checks, externe Konfiguration und „1 DB pro Service" werden in Story 1 konkret umgesetzt

---

## 4. Setup der Entwicklungsumgebung (45 Min)

Alle Teilnehmer am Start halten — bei Problemen sofort einsammeln und parallel helfen.

- Repository clonen
- Docker Compose: Backend-Services (Flight, Hotel, Car) + Consul + API-Gateway starten
- Health-Endpoints aller Services prüfen (`/health`)
- Consul-UI öffnen, registrierte Services sehen
- API-Gateway erreichbar machen (Port 8080)
- **OpenAPI-Schwenk (API First):** APIs der Backend-Services in Swagger UI / OpenAPI-Datei anschauen — Vertrag *vor* Implementierung
- Verweis: [vorbereitung.md](vorbereitung.md) für Schritt-für-Schritt-Anleitung

---

## 5. Kurz-Vortrag & Diskussion: „MS sind Monolithen mit Netzwerkproblemen" (30 Min)

These aufwerfen und mit der Gruppe sezieren.

**Pro-Argumente (das Netzwerk macht's komplex):**
- Partielle Ausfälle → Circuit Breaker, Bulkhead
- Keine verteilten Transaktionen → Saga
- Auffindbarkeit → Service Discovery
- Nachvollziehbarkeit → Distributed Tracing
- Konsistenz → Eventual Consistency, CQRS

**Aber auch on top:**
- Deployment-Topologie, Versionierung
- Team-Schnitte (Conway's Law)
- Polyglot Persistence
- Observability als eigenständige Disziplin

**Diskussions-Anker:** Welche Probleme hätte man im Monolithen auch? Welche entstehen erst durch das Netzwerk?

**Überleitung:** Deshalb dreht sich der praktische Teil (Stories 1–6, 10) stark um Kommunikation und Resilienz.

---

## 6. Hands-on: Stories 1–6

Detail-Anleitungen in den Story-Dateien. Pro Story: kurze Einleitung, Teilnehmer arbeiten selbständig, am Ende gemeinsamer Recap mit Diskussion.

### 6a. Story 1 — [Der erste cloud-native Booking-Service](stories/story-01-cloud-native-booking-service.md) (90 Min)
- Lernpointe: 12 Factor in der Praxis, Health-Checks, externe Konfiguration
- Recap-Frage: Welche Faktoren sind in eurer Umsetzung wirklich erfüllt?

### 6b. Story 2 — [Services dynamisch finden](stories/story-02-service-discovery.md) (60 Min)
- Lernpointe: Service Discovery mit Consul ersetzt hartkodierte URLs
- Recap-Frage: Was passiert, wenn Consul kurz weg ist?

### 6c. Story 3 — [Wenn der Flug ausfällt](stories/story-03-circuit-breaker.md) (60 Min)
- Lernpointe: Circuit Breaker, graceful Degradation
- Recap-Frage: Was ist ein gutes Default-Verhalten im „Open"-State?

### 6d. Story 4 — [Isolation ist Stärke](stories/story-04-bulkhead.md) (60 Min)
- Lernpointe: Ressourcen-Isolation, getrennte Thread-Pools / Connection-Pools
- Recap-Frage: Wo macht Bulkhead in eurer Architektur sonst noch Sinn?

### 6e. Story 5 — [Alles oder nichts – aber richtig](stories/story-05-saga.md) (60 Min)
- Lernpointe: Orchestration-Saga mit synchroner Kompensation
- Recap-Frage: Wer kennt das Endergebnis bei dieser Variante?

### 6f. Story 6 — [Die Saga wird leise](stories/story-06-choreography-saga.md) (60 Min)
- Lernpointe: Choreography-Saga via Events, Wissen verteilt sich
- Recap-Frage: Wann Orchestration, wann Choreography?

---

## 7. Kurz-Vortrag & Diskussion: CQRS (30 Min)

CQRS (Command Query Responsibility Segregation) heißt schlicht: **Trennung von Command und Query** — Schreib- und Lesemodell werden unabhängig voneinander entworfen und optimiert. Mehr ist erstmal nicht dahinter.

**Wichtig:** CQRS ist **kein Synonym für Eventsourcing** und **kein Synonym für Event-Driven Architecture**. Wie das Lesemodell aktualisiert wird, ist eine sekundäre Entscheidung — möglich sind Events, eine DB-Replica, ein Cache, ein nächtlicher Batch-Job oder eine materialisierte Sicht in der DB. Die Architektur-Idee ist unabhängig vom Transport.

**Warum wird CQRS so oft mit DDD genannt?** Historisch (Greg Young, ~2010) — CQRS, DDD und Eventsourcing wurden gemeinsam popularisiert. Notwendig ist diese Kombination nicht.

### Beispiel aus der Buchungsplattform

Der Booking-Service kennt aktuell nur den **Schreibpfad**: Eine Buchung kommt rein, die Saga läuft, am Ende ist das Ergebnis durch. Der Lesepfad — „zeige mir alle meine Reisen" — ist offen.

Der naive Ansatz, dafür bei jedem `GET` Flight, Hotel und Car parallel abzufragen und zu joinen, hat drei Probleme:

- **Langsam** — drei Backend-Calls pro Lese-Anfrage
- **Fragil** — ist ein Backend down, fehlen Daten oder die Anfrage scheitert ganz
- **Falsches Format** — die Backend-APIs liefern ihre internen Sichten, nicht die kundengerechte Historie

Die Lösung: **Lese- und Schreibmodell trennen**. Der Schreibpfad bleibt wie er ist (Saga). Daneben entsteht ein **Read-Model**, das auf eine völlig andere Anforderung optimiert ist: schnell antworten, denormalisiert, ausfalltolerant.

### Mehrwert

- **Ausfalltoleranz beim Lesen** — `GET /customers/{id}/bookings` funktioniert auch wenn Flight, Hotel oder Car gerade down sind. Das Read-Model ist eigenständig, kein Live-Aufruf der Backends. Ein anderer Resilienz-Mechanismus als Circuit Breaker: Datenredundanz statt Fehler-Isolation.
- **Format-Entkopplung** — Die Sicht für den Kunden („meine Reisen") muss nicht der Sicht der Backends entsprechen. Read-Model ist denormalisiert und kundenzentriert.
- **Unabhängige Skalierung** — Auf Plattformen wie Check24 werden tausendfach mehr Versicherungen *gelesen* als *abgeschlossen*. Beide Lasten mit demselben Modell zu bedienen wäre unwirtschaftlich. CQRS macht das ungleiche Verhältnis architektonisch sichtbar — und erlaubt für jede Seite die passende Persistenz, Caching-Strategie und Skalierung.

### Mechanismen für Read-Model-Updates

Die Wahl ist bewusst plural — keiner dieser Wege ist „CQRS":

- **Events** (z. B. die aus der Choreography-Saga): Bei `BookingCompleted` aktualisiert das Read-Model seine Sicht. Bequem, wenn Events ohnehin existieren.
- **DB-Replica mit asynchroner Replikation**: Read-Modell liegt auf einer separaten Datenbank, die der Schreibseite folgt.
- **Cache-Layer** vor der primären DB.
- **Materialisierte Sicht** in der DB selbst.
- **Batch-Job**, der das Read-Model periodisch neu aufbaut.

### Trade-off: Eventual Consistency

Zwischen dem Schreibvorgang (Saga fertig) und dem Auftauchen im Read-Model liegen Millisekunden bis Sekunden. Im UI ggf. mit einem „wird in Kürze sichtbar"-Hinweis auffangen. Wer starke Konsistenz braucht, sollte CQRS nicht oder nur mit synchroner Aktualisierung einsetzen.

### Diskussions-Anker

- Was würde sich ändern, wenn das Read-Model **nicht** aus Events, sondern aus einer DB-Replica gespeist würde? Was bliebe gleich?
- Wann lohnt CQRS *nicht*? (Hinweis: kleine Apps mit ähnlicher Lese- und Schreiblast.)
- Was passiert, wenn das Read-Model nach einem Crash leer ist — wie kommt der Stand zurück? (→ Event-Replay)
- Wo gehört das Read-Model fachlich hin: in den Booking-Service oder in einen eigenen `BookingHistoryService`?

---

## 8. Hands-on: Story 7 — [Den roten Faden im Log](stories/story-07-tracing.md) (60 Min)

- Lernpointe: Distributed Tracing macht Geschäftsvorgänge über Service-Grenzen hinweg sichtbar
- Recap-Frage: Wo hätte euch Tracing schon in Stories 3–6 geholfen?

---

## 9. Kurz-Vortrag & Diskussion: Backends for Frontends (20 Min)

Konzept: **spezialisierte Backends pro Client-Typ** (Web, Mobile, Smart-TV, IoT). Jedes BFF schneidet Daten, Payload-Größe und Latenz auf seinen Client zu.

**Im Workshop bewusst kein Hands-on:** Der Booking-Service erfüllt diese Rolle in unserem Beispiel faktisch schon — er aggregiert Flight, Hotel und Car für einen Web-/Mobile-Client. Eine echte BFF-Aufspaltung wäre konstruiert.

**Wann sinnvoll?**
- Mehrere Clients mit deutlich unterschiedlichen Anforderungen (z. B. Mobile braucht weniger Daten, optimierte Bilder, andere Auth)
- Latenzkritische Clients (TV, IoT)
- Unterschiedliche Release-Zyklen der Clients

**Wann nicht?**
- Ein Client / ein primärer Use Case → unnötige Code-Duplizierung
- Kleines Team — Pflege von n BFFs wird teuer
- Frontend und Backend ohnehin im selben Team — Modulith oder einfacher Aggregator reicht

---

## 10. Kurz-Vortrag & Diskussion: Downtimeless Deployment (20 Min)

Ist Zero-Downtime-Deployment „MS-spezifisch"? — **Nein**, Monolithen haben das gleiche Problem (Blue/Green, Canary). Für MS bekommt es eine eigene Note:

- **Rolling Updates pro Service** — pro Service unabhängig deploybar
- **API-Kompatibilität N & N+1** — alte und neue Version müssen koexistieren (Tolerant Reader, additive Änderungen, Deprecation-Phasen)
- **Datenbankmigrationen ohne Lock** — Expand/Contract-Pattern statt destruktive Migrations
- **Feature Toggles** — Deployment vs. Release entkoppeln

**Im Workshop bewusst kein Hands-on** — eher Pflicht-Wissen für Betrieb, weniger ein eigenes Pattern, das man „programmiert".

---

## 11. Abschluss & Kulturwandel (30 Min)

### Recap

- Was nehmt ihr aus den zwei Tagen mit?
- Welche Patterns würdet ihr morgen in eurer Architektur ansprechen?

### Diskussion: Kulturwandel durch Microservices

Microservices sind genauso eine Organisations- wie eine Technik-Frage.

- **Microservices ohne DevOps?** — Spoiler: nein. Wer keine Pipeline, kein Monitoring, kein Ownership-Modell hat, baut sich mit MS einen Operations-Albtraum.
- **„You build it, you run it"** — Was heißt das konkret? Arbeiten wir 24/7? Brauchen wir Wartungsteams? Wie steht ihr zu Bereitschaft?
- **Muss man „Cloud Native" sein?** — Cloud Native ist ein Toolkit (Container, Orchestrierung, declarative APIs), kein Zwang. Aber on-prem MS ohne Cloud-Native-Toolchain ist deutlich teurer.

### Grenzfrage

Wo liegt die Grenze: **Monolith → Modulith → Microservices**? Wer fühlt sich auf welcher Stufe heute zu Hause? Was wäre der Trigger, eine Stufe weiterzugehen?

### Feedback

- Was war hilfreich, was zu schnell, was zu langsam?
- Welche Themen würdet ihr stärker, welche schwächer gewichten?
