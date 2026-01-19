# Workshop Stories - Reisebuchungsplattform

Diese User Stories bilden den praktischen Teil des Microservices-Workshops. Alle Stories bauen auf einer gemeinsamen Domäne auf: Eine Reisebuchungsplattform, bei der Kunden Flüge, Hotels und Mietwagen buchen können.

## Workshop-Setup

**Bereitgestellte Backend-Services (Docker-Container vom Trainer):**
- `FlightService` - Verwaltet Flugbuchungen
- `HotelService` - Verwaltet Hotelbuchungen
- `CarService` - Verwaltet Mietwagenbuchungen

**Was die Teilnehmer bauen:**
- Die Buchungs-/Bestellstrecke (Orchestrierung der Backend-Services)

**Technische Rahmenbedingungen:**
- Framework: Frei wählbar (Spring Boot, Quarkus, Micronaut, o.ä.)
- Service Discovery: Consul (Docker-Container wird bereitgestellt)
- Messaging: REST-Webhooks (kein separater Message-Broker erforderlich)

---

# Tag 1: Grundlagen & Kommunikation

---

## Story 1: Der erste Booking-Service

**Thema:** Twelve-Factor-App, Health-Checks
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Unsere Reisebuchungsplattform braucht einen zentralen Service, der Buchungsanfragen entgegennimmt und an die spezialisierten Backend-Services weiterleitet. Bevor wir mit der eigentlichen Geschäftslogik beginnen, muss der Service die Grundlagen einer Cloud-nativen Anwendung erfüllen.

### User Story

Als **Betriebsteam**
möchte ich **den Gesundheitszustand des Booking-Service jederzeit abfragen können**,
damit **ich sicherstellen kann, dass der Service einsatzbereit ist und bei Problemen automatisch neu gestartet werden kann**.

### Akzeptanzkriterien

- [ ] Ein neues Projekt für den `BookingService` ist angelegt
- [ ] Der Service startet auf einem konfigurierbaren Port
- [ ] Ein Health-Endpoint unter `/health` oder `/actuator/health` ist verfügbar
- [ ] Der Health-Endpoint gibt HTTP 200 zurück, wenn der Service gesund ist
- [ ] Der Health-Endpoint prüft die Erreichbarkeit mindestens eines Backend-Services (z.B. FlightService)
- [ ] Der Service loggt auf stdout (nicht in Dateien)
- [ ] Konfiguration erfolgt über Umgebungsvariablen oder externe Config-Dateien (nicht hardcoded)

### Technische Hinweise

- **Bereitgestellte Services:** FlightService (läuft unter `http://localhost:8081`)
- **Twelve-Factor relevante Aspekte:**
  - III. Config - Konfiguration in Umgebungsvariablen
  - XI. Logs - Logs als Event-Streams behandeln
  - XII. Admin processes - Health-Check als Admin-Prozess
- **Empfohlene Libraries:**
  - Spring Boot: `spring-boot-starter-actuator`
  - Quarkus: `quarkus-smallrye-health`
  - Micronaut: `micronaut-management`

### Bonus (optional)

- Implementiere separate Endpoints für Liveness (`/health/live`) und Readiness (`/health/ready`)
- Füge Build-Informationen zum Health-Endpoint hinzu (Version, Build-Zeit)

---

## Story 2: Konfiguration externalisieren

**Thema:** External Configuration
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Der Booking-Service soll in verschiedenen Umgebungen (Entwicklung, Test, Produktion) laufen. Die URLs der Backend-Services, Timeouts und andere Parameter unterscheiden sich je nach Umgebung. Hardcodierte Werte sind keine Option.

### User Story

Als **Entwickler**
möchte ich **alle Konfigurationsparameter des Booking-Service externalisieren können**,
damit **der gleiche Build in verschiedenen Umgebungen ohne Code-Änderungen deployed werden kann**.

### Akzeptanzkriterien

- [ ] Die URLs der Backend-Services (Flight, Hotel, Car) sind konfigurierbar
- [ ] Timeouts für HTTP-Aufrufe sind konfigurierbar
- [ ] Es existieren Profile für mindestens zwei Umgebungen (z.B. `dev` und `prod`)
- [ ] Konfiguration kann über Umgebungsvariablen überschrieben werden
- [ ] Sensible Daten (falls vorhanden) sind nicht im Code-Repository
- [ ] Die aktuelle Konfiguration ist über einen Info-Endpoint einsehbar (ohne sensible Daten)

### Technische Hinweise

- **Konfigurationshierarchie (typisch):**
  1. Default-Werte im Code
  2. Externe Konfigurationsdateien
  3. Umgebungsvariablen (höchste Priorität)
- **Empfohlene Struktur:**
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
- **Twelve-Factor Aspekt:** III. Config - "Store config in the environment"

### Bonus (optional)

- Implementiere Hot-Reload der Konfiguration ohne Neustart
- Nutze einen zentralen Config-Server (z.B. Spring Cloud Config)

---

## Story 3: Services dynamisch finden

**Thema:** Service Discovery / Service Registry
**Zeitrahmen:** ca. 60 Minuten

### Kontext

In einer dynamischen Container-Umgebung ändern sich IP-Adressen und Ports der Services ständig. Der Booking-Service soll die Backend-Services nicht mehr über statische URLs ansprechen, sondern über eine Service Registry dynamisch finden.

### User Story

Als **Booking-Service**
möchte ich **die Backend-Services über ihren logischen Namen finden können**,
damit **ich nicht von statischen IP-Adressen oder Ports abhängig bin und Services elastisch skaliert werden können**.

### Akzeptanzkriterien

- [ ] Der Booking-Service registriert sich bei Consul
- [ ] Der Booking-Service findet FlightService, HotelService und CarService über Consul
- [ ] Die statischen URLs aus Story 2 werden durch Service-Namen ersetzt
- [ ] Der Health-Check wird bei Consul registriert
- [ ] Bei Ausfall eines Service-Instanz wird automatisch eine andere verwendet (falls verfügbar)
- [ ] Der Service deregistriert sich beim Herunterfahren

### Technische Hinweise

- **Consul UI:** `http://localhost:8500` (wird vom Trainer bereitgestellt)
- **Service-Namen:**
  - `flight-service`
  - `hotel-service`
  - `car-service`
- **Empfohlene Libraries:**
  - Spring Boot: `spring-cloud-starter-consul-discovery`
  - Quarkus: `quarkus-consul-config`
  - Alternativ: Direkte Consul HTTP API
- **Consul Health-Check-Typen:** HTTP, TCP, TTL

### Bonus (optional)

- Implementiere Client-Side Load Balancing bei mehreren Instanzen
- Nutze Consul Key-Value Store für zusätzliche Konfiguration

---

## Story 4: Ein Gateway für alle

**Thema:** API-Gateway
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Externe Clients (Web-App, Mobile-App) sollen nicht direkt mit den einzelnen Services kommunizieren. Ein API-Gateway dient als zentraler Einstiegspunkt, übernimmt Routing und kann übergreifende Aspekte wie Rate-Limiting implementieren.

### User Story

Als **Frontend-Entwickler**
möchte ich **alle Backend-Services über einen einzigen Einstiegspunkt ansprechen können**,
damit **ich mich nicht um die interne Service-Topologie kümmern muss und zentrale Policies (Rate-Limiting, CORS) einheitlich angewendet werden**.

### Akzeptanzkriterien

- [ ] Ein API-Gateway ist implementiert und läuft auf Port 8080
- [ ] Anfragen an `/api/flights/**` werden an den FlightService weitergeleitet
- [ ] Anfragen an `/api/hotels/**` werden an den HotelService weitergeleitet
- [ ] Anfragen an `/api/cars/**` werden an den CarService weitergeleitet
- [ ] Anfragen an `/api/bookings/**` werden an den BookingService weitergeleitet
- [ ] Das Gateway nutzt Service Discovery (keine hardcoded URLs)
- [ ] CORS ist für Frontend-Clients konfiguriert
- [ ] Basis Rate-Limiting ist implementiert (z.B. max. 100 Requests/Minute pro Client)

### Technische Hinweise

- **Empfohlene Technologien:**
  - Spring Cloud Gateway
  - Kong (Docker)
  - Traefik (Docker)
  - Eigene Implementierung mit Reverse-Proxy-Pattern
- **Routing-Konfiguration (Beispiel Spring Cloud Gateway):**
  ```yaml
  spring:
    cloud:
      gateway:
        routes:
          - id: flight-service
            uri: lb://flight-service
            predicates:
              - Path=/api/flights/**
  ```

### Bonus (optional)

- Implementiere Request-/Response-Logging im Gateway
- Füge einen `/health`-Aggregator hinzu, der den Status aller Services zusammenfasst

---

# Tag 2: Resilience & Events

---

## Story 5: Wenn der Flug ausfällt

**Thema:** Circuit Breaker
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Der FlightService ist zeitweise nicht erreichbar (Netzwerkprobleme, Überlast, Deployment). Der Booking-Service soll in diesem Fall nicht ebenfalls ausfallen, sondern graceful degradieren und dem Nutzer eine sinnvolle Alternative bieten.

### User Story

Als **Kunde**
möchte ich **auch bei Ausfall des Flugbuchungssystems eine Teilbuchung (Hotel + Mietwagen) durchführen können**,
damit **meine Reiseplanung nicht komplett blockiert wird**.

### Akzeptanzkriterien

- [ ] Ein Circuit Breaker ist um die FlightService-Aufrufe implementiert
- [ ] Nach 5 aufeinanderfolgenden Fehlern öffnet sich der Circuit
- [ ] Bei offenem Circuit wird ein Fallback ausgeführt (z.B. "Flugbuchung derzeit nicht verfügbar")
- [ ] Der Circuit schließt sich nach 30 Sekunden wieder (Half-Open State)
- [ ] Der aktuelle Circuit-Status ist über einen Endpoint abfragbar
- [ ] Timeouts sind konfiguriert (max. 3 Sekunden Wartezeit)
- [ ] Metriken über Circuit-Zustandswechsel werden geloggt

### Technische Hinweise

- **Empfohlene Libraries:**
  - Resilience4j (empfohlen, framework-agnostisch)
  - Spring Cloud Circuit Breaker
  - MicroProfile Fault Tolerance (Quarkus)
- **Circuit Breaker Zustände:**
  - CLOSED: Normale Funktion, Fehler werden gezählt
  - OPEN: Calls werden sofort mit Fallback beantwortet
  - HALF_OPEN: Testweise einzelne Calls durchlassen
- **Konfigurationsbeispiel (Resilience4j):**
  ```yaml
  resilience4j:
    circuitbreaker:
      instances:
        flightService:
          failureRateThreshold: 50
          waitDurationInOpenState: 30000
          slidingWindowSize: 10
  ```

### Bonus (optional)

- Implementiere Circuit Breaker für alle Backend-Services
- Visualisiere die Circuit-States in einem Dashboard

---

## Story 6: Isolation ist Stärke

**Thema:** Bulkhead Pattern
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Wenn der HotelService extrem langsam antwortet, könnten alle verfügbaren Threads des Booking-Service blockiert werden. Dadurch wären auch Anfragen an den FlightService betroffen, obwohl dieser einwandfrei funktioniert. Das Bulkhead-Pattern isoliert die Ressourcen.

### User Story

Als **Betriebsteam**
möchte ich **dass Probleme mit einem Backend-Service nicht die Aufrufe an andere Backend-Services beeinträchtigen**,
damit **ein langsamer oder fehlerhafter Service nicht das gesamte System blockiert**.

### Akzeptanzkriterien

- [ ] Jeder Backend-Service (Flight, Hotel, Car) hat einen eigenen Thread-Pool
- [ ] Die Thread-Pool-Größe ist pro Service konfigurierbar
- [ ] Bei erschöpftem Thread-Pool werden neue Anfragen sofort abgelehnt (nicht gequeued)
- [ ] Ein langsamer HotelService blockiert nicht die Aufrufe an FlightService
- [ ] Metriken über Thread-Pool-Auslastung sind verfügbar
- [ ] Abgelehnte Anfragen werden mit passendem HTTP-Status (503 Service Unavailable) beantwortet

### Technische Hinweise

- **Empfohlene Libraries:**
  - Resilience4j Bulkhead
  - MicroProfile Fault Tolerance `@Bulkhead`
- **Bulkhead-Typen:**
  - Semaphore-basiert: Begrenzt parallele Aufrufe
  - Thread-Pool-basiert: Eigener Thread-Pool pro Service
- **Konfigurationsbeispiel (Resilience4j):**
  ```yaml
  resilience4j:
    bulkhead:
      instances:
        hotelService:
          maxConcurrentCalls: 10
          maxWaitDuration: 100ms
  ```
- **Kombination mit Circuit Breaker:** Bulkhead und Circuit Breaker ergänzen sich

### Bonus (optional)

- Implementiere Adaptive Bulkhead, der sich an die Last anpasst
- Füge Queuing mit begrenzter Warteschlange hinzu

---

## Story 7: Alles oder nichts - aber richtig

**Thema:** Saga Pattern
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Eine Reisebuchung umfasst Flug, Hotel und Mietwagen. Wenn die Hotelbuchung fehlschlägt, nachdem der Flug bereits gebucht wurde, muss der Flug storniert werden. Klassische Datenbank-Transaktionen funktionieren nicht über Service-Grenzen hinweg.

### User Story

Als **Kunde**
möchte ich **eine Komplettbuchung (Flug + Hotel + Mietwagen) durchführen, die entweder vollständig erfolgreich ist oder komplett zurückgerollt wird**,
damit **ich nicht mit einer unvollständigen Buchung dastehe**.

### Akzeptanzkriterien

- [ ] Eine Buchungsanfrage für Flug, Hotel und Mietwagen kann gestellt werden
- [ ] Die Services werden nacheinander aufgerufen (Orchestration-Saga)
- [ ] Bei Fehler in einem Schritt werden alle vorherigen Schritte kompensiert (Rollback)
- [ ] Jeder Service bietet einen Stornieren/Kompensieren-Endpoint an
- [ ] Der Saga-Status wird persistiert (für Recovery nach Absturz)
- [ ] Der aktuelle Status einer Buchung kann abgefragt werden (PENDING, COMPLETED, COMPENSATING, FAILED)
- [ ] Kompensation wird mit Retry-Logik durchgeführt (Kompensation darf nicht fehlschlagen)

### Technische Hinweise

- **Saga-Varianten:**
  - **Choreography:** Services kommunizieren über Events (dezentral)
  - **Orchestration:** Ein Koordinator steuert den Ablauf (zentral) - empfohlen für diesen Workshop
- **Saga-Schritte für Reisebuchung:**
  1. Flug buchen
  2. Hotel buchen
  3. Mietwagen buchen
  - Kompensation: Mietwagen stornieren → Hotel stornieren → Flug stornieren
- **Webhook-Callbacks:** Services rufen bei Statusänderung einen Callback-URL auf
- **API-Endpunkte (Backend-Services):**
  - `POST /bookings` - Buchung erstellen
  - `DELETE /bookings/{id}` - Buchung stornieren (Kompensation)

### Bonus (optional)

- Implementiere parallele Ausführung unabhängiger Schritte
- Füge Timeout-Handling hinzu (Was passiert, wenn ein Service nicht antwortet?)

---

## Story 8: Events erzählen Geschichten

**Thema:** Event-Driven Architecture, CQRS
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Nach erfolgreicher Buchung sollen verschiedene Aktionen ausgelöst werden: Bestätigungs-E-Mail versenden, Statistiken aktualisieren, Partner benachrichtigen. Diese Aktionen sollen die Buchung nicht verlangsamen. Zusätzlich soll eine optimierte Leseansicht für Buchungsübersichten entstehen.

### User Story

Als **System**
möchte ich **bei erfolgreicher Buchung ein Event veröffentlichen, das von interessierten Services verarbeitet werden kann**,
damit **die Buchung schnell abgeschlossen wird und nachgelagerte Prozesse asynchron laufen können**.

### Akzeptanzkriterien

- [ ] Nach erfolgreicher Buchung wird ein `BookingCompleted`-Event veröffentlicht
- [ ] Ein Notification-Service empfängt das Event und loggt eine Bestätigungs-Nachricht
- [ ] Ein Analytics-Service empfängt das Event und aktualisiert Buchungsstatistiken
- [ ] Events werden über REST-Webhooks an registrierte Subscriber verteilt
- [ ] Ein separates Read-Model für Buchungsübersichten existiert (CQRS)
- [ ] Das Read-Model wird durch Events aktualisiert (eventual consistency)
- [ ] Die Buchungs-API antwortet sofort, ohne auf Subscriber zu warten

### Technische Hinweise

- **Event-Struktur (Beispiel):**
  ```json
  {
    "eventType": "BookingCompleted",
    "eventId": "uuid",
    "timestamp": "2024-01-15T10:30:00Z",
    "payload": {
      "bookingId": "B-12345",
      "customerId": "C-999",
      "totalAmount": 1250.00,
      "items": ["FLIGHT", "HOTEL", "CAR"]
    }
  }
  ```
- **Webhook-Registrierung:**
  - Services registrieren sich mit ihrer Callback-URL
  - `POST /webhooks/subscribe` mit `{ "eventType": "BookingCompleted", "callbackUrl": "http://..." }`
- **CQRS Read-Model:**
  - Denormalisierte Ansicht, optimiert für Abfragen
  - Wird durch Events aktualisiert, nicht durch direkte DB-Schreiboperationen
- **Eventual Consistency:** Das Read-Model kann kurzzeitig veraltet sein

### Bonus (optional)

- Implementiere Event-Replay zum Neuaufbau des Read-Models
- Füge Idempotenz-Handling für doppelte Events hinzu

---

# Optionale Stories (bei Zeit)

---

## Story 9: Mobile First

**Thema:** Backends for Frontends (BFF)
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Die mobile App benötigt andere Daten als die Web-Anwendung: kleinere Payloads, zusammengefasste Informationen, optimiert für langsame Verbindungen. Ein spezialisiertes Backend für das mobile Frontend löst dieses Problem.

### User Story

Als **Mobile-App-Entwickler**
möchte ich **ein Backend, das genau die Daten liefert, die meine App braucht**,
damit **ich keine überflüssigen Daten übertragen muss und die App auch bei schlechter Verbindung performant bleibt**.

### Akzeptanzkriterien

- [ ] Ein Mobile-BFF-Service ist implementiert
- [ ] Das BFF aggregiert Daten aus mehreren Backend-Services
- [ ] Die Payload ist kompakter als die der einzelnen Services
- [ ] Nur für Mobile relevante Felder werden zurückgegeben
- [ ] Das BFF cached häufig angefragte, selten ändernde Daten
- [ ] GraphQL oder spezialisierte REST-Endpoints sind verfügbar

### Technische Hinweise

- **BFF vs. API Gateway:**
  - API Gateway: Routing, übergreifende Policies
  - BFF: Frontend-spezifische Aggregation und Transformation
- **Payload-Optimierung:**
  - Felder auswählen (nur benötigte)
  - Daten zusammenfassen (z.B. "Reise" statt separate Flight/Hotel/Car)
  - Pagination für Listen
- **Empfohlene Technologien:**
  - GraphQL für flexible Abfragen
  - REST mit Sparse Fieldsets

### Bonus (optional)

- Implementiere ein separates BFF für die Web-Anwendung
- Füge Offline-Fähigkeit durch intelligentes Caching hinzu

---

## Story 10: Ohne Unterbrechung

**Thema:** Downtimeless Deployment
**Zeitrahmen:** ca. 60 Minuten

### Kontext

Der Booking-Service soll aktualisiert werden können, ohne dass laufende Buchungen unterbrochen werden. Kunden sollen von Updates nichts mitbekommen.

### User Story

Als **Betriebsteam**
möchte ich **neue Versionen des Booking-Service deployen können, ohne dass es zu Ausfallzeiten kommt**,
damit **Kunden jederzeit buchen können und wir häufiger deployen können**.

### Akzeptanzkriterien

- [ ] Mindestens 2 Instanzen des Booking-Service laufen
- [ ] Rolling Update: Neue Version wird schrittweise ausgerollt
- [ ] Laufende Requests werden abgeschlossen bevor eine Instanz stoppt (Graceful Shutdown)
- [ ] Der Health-Check unterscheidet zwischen Liveness und Readiness
- [ ] Neue Instanzen erhalten erst Traffic, wenn sie ready sind
- [ ] Bei fehlgeschlagenem Deployment wird automatisch zurückgerollt
- [ ] Während des Deployments sind alle Funktionen verfügbar

### Technische Hinweise

- **Graceful Shutdown:**
  - Keine neuen Requests annehmen
  - Laufende Requests abschließen
  - Ressourcen freigeben
  - Aus Service Registry deregistrieren
- **Readiness vs. Liveness:**
  - Liveness: "Lebt der Prozess?" → Neustart bei Failure
  - Readiness: "Kann der Service Requests verarbeiten?" → Kein Traffic bei Failure
- **Rolling Update Strategie:**
  - maxUnavailable: 0 (nie weniger als die gewünschte Anzahl)
  - maxSurge: 1 (eine zusätzliche Instanz während Update)
- **Testing:** Mit Last-Generator während des Deployments testen

### Bonus (optional)

- Implementiere Blue-Green Deployment als Alternative
- Füge Canary Releases hinzu (nur 10% Traffic auf neue Version)

---

# Anhang

## Service-API-Dokumentation

### FlightService (Port 8081)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | `/flights` | Liste verfügbarer Flüge |
| GET | `/flights/{id}` | Flugdetails |
| POST | `/bookings` | Flug buchen |
| GET | `/bookings/{id}` | Buchungsdetails |
| DELETE | `/bookings/{id}` | Buchung stornieren |
| GET | `/health` | Health-Check |

### HotelService (Port 8082)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | `/hotels` | Liste verfügbarer Hotels |
| GET | `/hotels/{id}` | Hoteldetails |
| POST | `/bookings` | Hotel buchen |
| GET | `/bookings/{id}` | Buchungsdetails |
| DELETE | `/bookings/{id}` | Buchung stornieren |
| GET | `/health` | Health-Check |

### CarService (Port 8083)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | `/cars` | Liste verfügbarer Mietwagen |
| GET | `/cars/{id}` | Mietwagendetails |
| POST | `/bookings` | Mietwagen buchen |
| GET | `/bookings/{id}` | Buchungsdetails |
| DELETE | `/bookings/{id}` | Buchung stornieren |
| GET | `/health` | Health-Check |

## Story-Abhängigkeiten

```
Story 1 (Grundsetup)
    │
    ├── Story 2 (Configuration)
    │       │
    │       └── Story 3 (Service Discovery)
    │               │
    │               └── Story 4 (API Gateway)
    │
    └── Story 5 (Circuit Breaker)
            │
            ├── Story 6 (Bulkhead)
            │
            └── Story 7 (Saga)
                    │
                    └── Story 8 (Events/CQRS)

Optional:
    Story 4 → Story 9 (BFF)
    Story 1 + Story 3 → Story 10 (Downtimeless)
```

## Glossar

| Begriff | Erklärung |
|---------|-----------|
| Circuit Breaker | Schutzmechanismus, der bei wiederholten Fehlern weitere Aufrufe unterbricht |
| Bulkhead | Isolierung von Ressourcen, um Kaskadenausfälle zu verhindern |
| Saga | Pattern für verteilte Transaktionen mit Kompensationslogik |
| CQRS | Trennung von Lese- und Schreibmodellen |
| BFF | Backend for Frontend - spezialisiertes Backend pro Client-Typ |
| Eventual Consistency | Daten werden irgendwann konsistent, nicht sofort |
| Graceful Shutdown | Kontrolliertes Herunterfahren ohne Abbruch laufender Operationen |
