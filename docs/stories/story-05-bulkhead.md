# Story 5: Isolation ist Stärke

**Thema:** Bulkhead Pattern
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Wenn der HotelService extrem langsam antwortet, könnten alle verfügbaren Threads des Booking-Service blockiert werden. Dadurch wären auch Anfragen an den FlightService betroffen, obwohl dieser einwandfrei funktioniert. Das Bulkhead-Pattern isoliert die Ressourcen.

## User Story

Als **Betriebsteam**
möchte ich **dass Probleme mit einem Backend-Service nicht die Aufrufe an andere Backend-Services beeinträchtigen**,
damit **ein langsamer oder fehlerhafter Service nicht das gesamte System blockiert**.

## Akzeptanzkriterien

- [ ] Jeder Backend-Service (Flight, Hotel, Car) hat einen eigenen Thread-Pool
- [ ] Die Thread-Pool-Größe ist pro Service konfigurierbar
- [ ] Bei erschöpftem Thread-Pool werden neue Anfragen sofort abgelehnt (nicht gequeued)
- [ ] Ein langsamer HotelService blockiert nicht die Aufrufe an FlightService
- [ ] Metriken über Thread-Pool-Auslastung sind verfügbar
- [ ] Abgelehnte Anfragen werden mit passendem HTTP-Status (503 Service Unavailable) beantwortet

## Technische Hinweise

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

## Bonus (optional)

- Implementiere Adaptive Bulkhead, der sich an die Last anpasst
- Füge Queuing mit begrenzter Warteschlange hinzu
