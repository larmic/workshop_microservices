# Story 4: Wenn der Flug ausfällt

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

- Implementiere Circuit Breaker für alle Backend-Services
- Visualisiere die Circuit-States in einem Dashboard
