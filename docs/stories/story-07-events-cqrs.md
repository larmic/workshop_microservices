# Story 7: Events erzählen Geschichten

**Thema:** Event-Driven Architecture, CQRS
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Nach erfolgreicher Buchung sollen verschiedene Aktionen ausgelöst werden: Bestätigungs-E-Mail versenden, Statistiken aktualisieren, Partner benachrichtigen. Diese Aktionen sollen die Buchung nicht verlangsamen. Zusätzlich soll eine optimierte Leseansicht für Buchungsübersichten entstehen.

## User Story

Als **System**
möchte ich **bei erfolgreicher Buchung ein Event veröffentlichen, das von interessierten Services verarbeitet werden kann**,
damit **die Buchung schnell abgeschlossen wird und nachgelagerte Prozesse asynchron laufen können**.

## Akzeptanzkriterien

- [ ] Nach erfolgreicher Buchung wird ein `BookingCompleted`-Event veröffentlicht
- [ ] Ein Notification-Service empfängt das Event und loggt eine Bestätigungs-Nachricht
- [ ] Ein Analytics-Service empfängt das Event und aktualisiert Buchungsstatistiken
- [ ] Events werden über REST-Webhooks an registrierte Subscriber verteilt
- [ ] Ein separates Read-Model für Buchungsübersichten existiert (CQRS)
- [ ] Das Read-Model wird durch Events aktualisiert (eventual consistency)
- [ ] Die Buchungs-API antwortet sofort, ohne auf Subscriber zu warten

## Technische Hinweise

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

## Bonus (optional)

- Implementiere Event-Replay zum Neuaufbau des Read-Models
- Füge Idempotenz-Handling für doppelte Events hinzu
