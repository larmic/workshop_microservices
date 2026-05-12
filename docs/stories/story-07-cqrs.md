# Story 7: Lesen ist nicht Schreiben – CQRS für die Buchungs-Historie

**Thema:** CQRS (Command Query Responsibility Segregation)
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Der Booking-Service kennt aktuell nur den **Schreibpfad**: Eine Buchung kommt rein, die Saga läuft, am Ende ist das Ergebnis durch. Der Lesepfad — „zeige mir alle meine Reisen" — ist offen.

Den naiven Ansatz, dafür bei jedem `GET` Flight, Hotel und Car parallel abzufragen und zu joinen, will man nicht haben:

- **Langsam** — drei Backend-Calls pro Lese-Anfrage
- **Fragil** — ist ein Backend down, fehlen Daten oder die Anfrage scheitert ganz
- **Falsches Format** — die Backend-APIs liefern ihre internen Sichten, nicht die kundengerechte Historie

Die Lösung ist eines der zentralen Microservice-Patterns: **Lese- und Schreibmodell trennen**. Der Schreibpfad bleibt wie er ist (Saga, Backend-Services, Story 5 + 6). Daneben entsteht ein **Read-Model**, das auf eine völlig andere Anforderung optimiert ist: schnell antworten, denormalisiert, ausfalltolerant.

Das ist der Kern von **CQRS** — und wichtig: CQRS ist **kein Synonym für Eventing**. CQRS heißt nur „Trennung von Command und Query". Wie das Read-Model aktualisiert wird, ist eine sekundäre Entscheidung. In dieser Story nutzen wir die Events aus Story 6 als Mechanismus, weil sie sowieso schon da sind — aber genauso könnten wir das Read-Model aus einer DB-Replica, einem Cache, einem nächtlichen Batch-Job oder einer materialisierten Sicht speisen. Die Architektur-Idee ist unabhängig vom Transport.

**Echtes Beispiel:** Auf Plattformen wie Check24 werden tausendfach mehr Versicherungen *gelesen* als *abgeschlossen*. Beide Lasten mit demselben Modell zu bedienen wäre unwirtschaftlich. CQRS macht das ungleiche Verhältnis architektonisch sichtbar — und erlaubt für jede Seite die passende Persistenz, Caching-Strategie und Skalierung.

## User Story

Als **Kunde**
möchte ich **meine komplette Buchungshistorie schnell und auch dann abrufen können, wenn einzelne Backend-Services nicht erreichbar sind**,
damit **ich eine zuverlässige Übersicht aller meiner Reisen habe**.

## Akzeptanzkriterien

- [ ] Bei erfolgreichem Saga-Abschluss publiziert der Booking-Service ein `BookingCompleted`-Event (mit `customerId`, `sagaId`, Items aus Flight/Hotel/Car)
- [ ] Bei einem Storno wird ein `BookingCancelled`-Event publiziert
- [ ] Ein **Read-Model** im Booking-Service lauscht auf diese Events und pflegt eine denormalisierte Sicht pro Kunde (z. B. `Map<CustomerId, List<BookingSummary>>`)
- [ ] Neuer Endpoint: `GET /customers/{id}/bookings` liefert die aggregierte Historie aus dem Read-Model
- [ ] Der Endpoint funktioniert **auch wenn Flight, Hotel oder Car gerade down sind** (kein Live-Aufruf der Backends)
- [ ] Eventual Consistency wird in der API-Doku explizit erwähnt („Daten können wenige Sekunden hinter dem Schreibstand liegen")
- [ ] Schreibpfad (Saga) und Lesepfad (Read-Model) sind im Code **erkennbar getrennt** — auch wenn sie im selben Service leben

## Technische Hinweise

- **Strikte Trennung im Code:**
  - **Write-Pfad:** `POST /booking/bookings` → Saga → bei Abschluss → Event raus
  - **Read-Pfad:** `GET /customers/{id}/bookings` → liest **nur** Read-Model, **nie** die Backends
- **Read-Model-Speicher:**
  - Für den Workshop reicht ein in-memory `Map<CustomerId, List<BookingSummary>>`
  - In Produktion: separate Tabelle, Read-Replica, Cache oder dedizierte DB (z. B. NoSQL wegen Leseperformance)
- **Pädagogischer Kern:**
  - CQRS ist **nicht** „Events haben". CQRS ist „Lese- und Schreibmodell unabhängig optimieren".
  - Events sind hier ein bequemer **Mechanismus**, nicht der Sinn des Patterns.
  - Vergleichs-Alternativen ohne Events: DB-Replica mit asynchroner Replikation, Cache-Layer vor der primären DB, materialisierte Sicht in der DB selbst.
- **Eventual Consistency:** Zwischen dem Schreibvorgang (Saga fertig) und dem Auftauchen im Read-Model liegen Millisekunden bis Sekunden. Im UI ggf. mit einem „wird in Kürze sichtbar"-Hinweis auffangen.
- **Idempotenz:** Das Read-Model muss doppelte Events tolerieren (gleicher `eventId` → kein zweiter Eintrag).

## Diskussions-Anker

- Was würde sich ändern, wenn das Read-Model **nicht** aus Events, sondern aus einer DB-Replica gespeist würde? Was bliebe gleich?
- Wann lohnt CQRS *nicht*? (Hint: kleine Apps mit ähnlicher Lese- und Schreiblast.)
- Was passiert, wenn das Read-Model nach einem Crash leer ist — wie kommt der Stand zurück? (→ Event-Replay, siehe Bonus)
- Wo gehört das Read-Model fachlich hin: in den Booking-Service oder in einen eigenen `BookingHistoryService`?

## Bonus (optional)

- **Event-Replay:** Read-Model nach Crash aus dem Event-Log neu aufbauen
- **Persistentes Read-Model:** SQLite-Datei oder separate DB-Tabelle statt In-Memory
- **Zweite Sicht:** Ein zusätzliches Read-Model für Statistiken (z. B. „Top-Reiseziele pro Monat") — zeigt, dass mehrere Read-Models pro Schreibstrom existieren können
- **Getrennter Service:** Read-Model in einen eigenen `BookingHistoryService` extrahieren — physische Trennung von Lese- und Schreibpfad
