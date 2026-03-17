# Story 1: Der erste cloud-native Booking-Service

**Thema:** Twelve-Factor-App, Health-Checks, Externe Konfiguration
**Zeitrahmen:** ca. 90 Minuten

## Kontext

Unsere Reisebuchungsplattform braucht einen zentralen Service, der Buchungsanfragen entgegennimmt und an die spezialisierten Backend-Services weiterleitet. Bevor wir mit der eigentlichen Geschäftslogik beginnen, muss der Service die Grundlagen einer Cloud-nativen Anwendung erfüllen: Er muss überwachbar sein, in verschiedenen Umgebungen laufen können und die Prinzipien der Twelve-Factor-App befolgen.

## User Story

Als **Entwicklungs- und Betriebsteam**
möchte ich **einen cloud-nativen Booking-Service mit Health-Checks und externalisierbarer Konfiguration**,
damit **der Service überwacht werden kann, bei Problemen automatisch neu gestartet wird und derselbe Build ohne Code-Änderungen in verschiedenen Umgebungen läuft**.

## Akzeptanzkriterien

**Grundsetup & Health-Checks:**
- [ ] Ein neues Projekt für den `BookingService` ist angelegt
- [ ] Der Service startet auf einem konfigurierbaren Port
- [ ] Ein Health-Endpoint unter `/health` ist verfügbar
- [ ] Der Health-Endpoint gibt HTTP 200 zurück, wenn der Service gesund ist
- [ ] Der Service loggt auf stdout (nicht in Dateien)

**Externe Konfiguration:**
- [ ] Die URLs der Backend-Services (Flight, Hotel, Car) sind konfigurierbar
- [ ] Timeouts für HTTP-Aufrufe sind konfigurierbar
- [ ] Es existieren Profile für mindestens zwei Umgebungen (z.B. `dev` und `prod`)
- [ ] Konfiguration kann über Umgebungsvariablen überschrieben werden
- [ ] Sensible Daten (falls vorhanden) sind nicht im Code-Repository
- [ ] Die aktuelle Konfiguration ist über einen Info-Endpoint einsehbar (ohne sensible Daten)

**Erste Orchestrierung:**
- [ ] `GET /booking/offers` ruft FlightService (`GET /flights`), HotelService (`GET /hotels`) und CarService (`GET /cars`) auf
- [ ] Die kombinierten Ergebnisse (verfügbare Flüge, Hotels, Mietwagen) werden als JSON zurückgegeben
- [ ] Die konfigurierbaren URLs und Timeouts aus dem Config-Block werden für die Aufrufe genutzt
- [ ] Kein Error-Handling nötig — wenn ein Service nicht erreichbar ist, darf der Aufruf fehlschlagen

## Technische Hinweise

- **Bereitgestellte Services:** FlightService (läuft unter `http://localhost:8081`)
- **Twelve-Factor relevante Aspekte:**
  - II. Dependencies — Abhängigkeiten explizit deklarieren (`pom.xml`, `build.gradle` oder `go.mod`), keine impliziten System-Dependencies
  - III. Config — Konfiguration in Umgebungsvariablen, nicht im Code
  - V. Build, release, run — Strikte Trennung von Build und Run (das Dockerfile ist ein gutes Beispiel)
  - VI. Processes — Der Service ist stateless, kein lokaler State zwischen Requests
  - VII. Port binding — Der Service exportiert HTTP über Port-Binding, kein externer App-Server nötig
  - IX. Disposability — Schnelles Starten und Graceful Shutdown ermöglichen elastisches Skalieren
  - XI. Logs — Logs als Event-Streams auf stdout behandeln
  - XII. Admin processes — Health-Check als Admin-Prozess
- **Konfigurationshierarchie (typisch):**
  1. Default-Werte im Code
  2. Externe Konfigurationsdateien
  3. Umgebungsvariablen (höchste Priorität)
- **Empfohlene Konfigurationsstruktur:**
  ```yaml
  booking:
    services:
      flight:
        url: http://localhost:8081
        timeout: 5000
      hotel:
        url: http://localhost:8082
        timeout: 5000
      car:
        url: http://localhost:8083
        timeout: 5000
  ```

## Bonus (optional)

- Implementiere separate Endpoints für Liveness (`/health/live`) und Readiness (`/health/ready`)
  - Liveness: "Läuft mein Prozess?"
  - Readiness: "Kann ich Traffic verarbeiten?"
- Füge Build-Informationen zum Health-Endpoint hinzu (Version, Build-Zeit)
- Implementiere Hot-Reload der Konfiguration ohne Neustart
- Nutze einen zentralen Config-Server (z.B. Spring Cloud Config)
