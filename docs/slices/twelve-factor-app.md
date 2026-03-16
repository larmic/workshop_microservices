# The Twelve-Factor App

---

## Warum Twelve-Factor?

- **Cloud-native** — für moderne Plattformen gebaut
- **Portabel** — läuft überall gleich
- **Skalierbar** — horizontal ohne Umbau

---

## Die 12 Faktoren — Übersicht

1. Codebase
2. Dependencies
3. Config
4. Backing Services
5. Build, Release, Run
6. Processes
7. Port Binding
8. Concurrency
9. Disposability
10. Dev/Prod Parity
11. Logs
12. Admin Processes

---

## I. Codebase

- **Ein** Repository — **viele** Deploys
- Gleicher Code für Dev, Staging, Prod
- Unterschiede nur durch Konfiguration

---

## II. Dependencies

- Alle Abhängigkeiten **explizit** deklarieren
- Keine impliziten System-Dependencies

```xml
<!-- pom.xml -->
<dependency>
    <groupId>io.javalin</groupId>
    <artifactId>javalin</artifactId>
    <version>6.6.0</version>
</dependency>
```

---

## III. Config

- Konfiguration in **Umgebungsvariablen**
- Niemals URLs, Credentials o.ä. im Code

```yaml
booking:
  services:
    flight:
      url: ${FLIGHT_SERVICE_URL:http://localhost:8081}
      timeout: ${FLIGHT_TIMEOUT:5000}
```

---

## IV. Backing Services

- Datenbanken, Queues, APIs = **angehängte Ressourcen**
- Austauschbar ohne Code-Änderung
- FlightService, HotelService, CarService → austauschbare Backing Services

---

## V. Build, Release, Run

- **Build** — Code kompilieren, Dependencies auflösen
- **Release** — Build + Konfiguration kombinieren
- **Run** — Release in Umgebung starten

```dockerfile
FROM ghcr.io/graalvm/jdk-community:25
COPY target/app.jar /app.jar
CMD ["java", "-jar", "/app.jar"]
```

---

## VI. Processes

- Services sind **stateless**
- Kein lokaler State zwischen Requests
- Session-Daten → Backing Service (Redis, DB)

---

## VII. Port Binding

- Service exportiert HTTP über **eigenen Port**
- Kein externer Application Server nötig
- `BookingService` bindet sich selbst an Port 8080

---

## VIII. Concurrency

- Skalierung über **Prozesse**, nicht Threads
- Mehr Last → mehr Instanzen
- Horizontale Skalierung statt vertikaler

---

## IX. Disposability

- **Schneller Start** — Sekunden, nicht Minuten
- **Graceful Shutdown** — laufende Requests abschließen
- Ermöglicht elastisches Skalieren und schnelle Deploys

---

## X. Dev/Prod Parity

- Dev, Staging, Prod **möglichst identisch**
- Gleiche Backing Services (nicht SQLite in Dev, Postgres in Prod)
- Kleine Gaps: Zeit, Personal, Tools

---

## XI. Logs

- Logs = **Event-Streams**
- Immer auf `stdout` schreiben
- Plattform kümmert sich um Aggregation

---

## XII. Admin Processes

- Einmal-Tasks als **eigene Prozesse** ausführen
- Gleicher Code, gleiche Konfiguration
- Health-Checks als Beispiel: `/health`

---

## Workshop-Bezug: Story 1

In Story 1 setzen wir diese Faktoren um:

- **II. Dependencies** — explizite Deklaration
- **III. Config** — Service-URLs per Umgebungsvariable
- **V. Build, Release, Run** — Dockerfile
- **VI. Processes** — Stateless BookingService
- **VII. Port Binding** — Self-contained HTTP
- **IX. Disposability** — Fast Startup
- **XI. Logs** — stdout
- **XII. Admin Processes** — Health-Endpoint

---

## Diskussion

- Welche Faktoren setzt ihr in euren Projekten bereits um?
- Wo seht ihr die größten Herausforderungen?
- Welche Faktoren sind in eurem Kontext besonders relevant?
