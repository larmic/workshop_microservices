# Story 3: Wenn der Flug ausfällt

**Thema:** Circuit Breaker
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Der FlightService ist zeitweise nicht erreichbar (Netzwerkprobleme, Überlast, Deployment). Der Booking-Service soll in diesem Fall nicht ebenfalls ausfallen, sondern graceful degradieren und dem Nutzer eine sinnvolle Alternative bieten.

## User Story

Als **Kunde**
möchte ich **auch bei Ausfall des Flugbuchungssystems eine Teilbuchung (Hotel + Mietwagen) durchführen können**,
damit **meine Reiseplanung nicht komplett blockiert wird**.

## Akzeptanzkriterien

- [ ] Ein Circuit Breaker ist um die FlightService-Aufrufe implementiert
- [ ] Nach 5 aufeinanderfolgenden Fehlern öffnet sich der Circuit
- [ ] Bei offenem Circuit wird ein Fallback ausgeführt (z.B. "Flugbuchung derzeit nicht verfügbar")
- [ ] Der Circuit schließt sich nach 30 Sekunden wieder (Half-Open State)
- [ ] Der aktuelle Circuit-Status ist über einen Endpoint abfragbar
- [ ] Timeouts sind konfiguriert (max. 3 Sekunden Wartezeit)
- [ ] Metriken über Circuit-Zustandswechsel werden geloggt

## Technische Hinweise

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

## Bonus (optional)

- Implementiere Circuit Breaker für **alle** Backend-Services (Flight, Hotel, Car) — jeder mit eigenem State.
- Visualisiere die Circuit-States in einem Dashboard.

## Reference-Implementierung

Unter `services/booking/story3/` liegt eine vollständige Go-Reference, die direkt beide Bonus-Punkte abdeckt:

- **Selbstgebauter Circuit Breaker** (keine Library) — `services/booking/story3/circuitbreaker/circuitbreaker.go`. Drei Zustände, lazy Übergang von OPEN nach HALF_OPEN beim nächsten Call, atomarer Probe-Slot gegen Probe-Storms in HALF_OPEN.
- **Drei separate CBs** für Flight, Hotel und Car. Jeder Backend-Aufruf läuft durch seinen eigenen Breaker; bei Fallback wird der Teilbereich leer zurückgegeben (`flights: []` bzw. `flight: null`).
- **Detailliertes Logging**: jede Zustandsänderung, jede SHORT-CIRCUIT-Entscheidung mit Restzeit, jeder Fehlerzähler-Inkrement und jeder Probe-Versuch wird mit Begründung im Log ausgegeben — so ist im Workshop nachvollziehbar, *warum* der Breaker gerade so handelt.
- **Dashboard-Visualisierung** (`services/dashboard/`): pro CB eine Karte mit Badge (CLOSED/OPEN/HALF_OPEN), Countdown bis HALF_OPEN, Failure-Counter, Anzahl Fallbacks. Aktualisierung im 1-Sekunden-Takt.
- **Chaos-Steuerung**: Backend-Services (Flight, Hotel, Car) bieten `/admin/chaos` mit den Modi `normal`, `slow`, `fail`. Das Dashboard erlaubt das Schalten pro Service oder pro Replica, sodass HALF_OPEN-Recovery live demonstriert werden kann.

### Workshop-Drehbuch

1. Stack starten: `docker compose -f services/docker-compose.yml -f services/docker-compose.infra.yml -f services/docker-compose.reference.yml up -d`
2. Dashboard öffnen (`http://localhost/dashboard`). Alle drei CBs zeigen CLOSED.
3. Flight im Dashboard auf "Fehler" stellen. Fünfmal `curl localhost:8087/booking/offers` schicken — der Flight-CB wechselt nach OPEN, Antwort enthält `flights: []` und Header `X-Circuit-Open: flight`.
4. Logs des `booking-story3` Containers ansehen (`docker compose logs -f booking-story3`) — jede Entscheidung des Breakers ist sichtbar.
5. Flight zurück auf "Normal". Nach 30 Sekunden geht der CB beim nächsten Call automatisch in HALF_OPEN, der Probe-Call schließt den Breaker wieder.
