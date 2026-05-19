# Booking Service — Story 7: Distributed Tracing

## Was diese Story zeigt

Story 7 baut auf der Choreography-Saga aus Story 6 auf und ergänzt sie um
**Distributed Tracing**. Alle vier Services (Booking, Flight, Hotel, Car)
verteilen einen gemeinsamen Trace-Kontext nach W3C-Trace-Context-Standard
und schreiben die Trace-ID in jede strukturierte Logzeile.

Das macht einen kompletten Buchungsvorgang über alle Services hinweg
nachverfolgbar — auch über die asynchrone Compensation-Grenze, an der
HTTP-Header verloren gehen würden.

## Was technisch neu ist gegenüber Story 6

| Aspekt | Story 6 | Story 7 |
|---|---|---|
| **Logging** | `log.Printf` mit Format-String | `slog`-JSON-Handler, jede Zeile mit `trace_id` |
| **Trace-Kontext eingehend** | Nicht vorhanden | `tracing.Middleware` liest `traceparent` oder generiert |
| **Trace-Kontext ausgehend** | Nicht propagiert | `tracing.Inject` setzt `traceparent` auf jeden Outbound-Call |
| **Async-Grenze (Compensation)** | Kein Trace-Bezug | `traceparent` als Event-Property im Body |
| **Saga aus Logs rekonstruieren** | `grep <saga-id>` (nur Booking-Side) | `grep <trace-id>` über alle Services |

## Demo

Eine Buchung absenden:

```bash
curl -X POST http://localhost/api/booking-ref-story7/booking/bookings \
  -H "Content-Type: application/json" \
  -d '{"customerName":"Demo","flightId":"LH123","hotelId":"H1","carId":"C1"}'
```

Trace-ID aus den Logs ablesen und filtern:

```bash
docker compose logs booking-ref-story7 flight hotel car | jq -c 'select(.trace_id=="<trace-id>")'
```

Oder mit eigenem `traceparent` (z.B. vom Frontend / API-Gateway):

```bash
curl -X POST http://localhost/api/booking-ref-story7/booking/bookings \
  -H "Content-Type: application/json" \
  -H "traceparent: 00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-bbbbbbbbbbbbbbbb-01" \
  -d '{"customerName":"Demo","flightId":"LH123","hotelId":"H1","carId":"C1"}'
```

Die Response trägt den `traceparent` zurück (effektive Trace-ID), und alle
beteiligten Services loggen mit `trace_id=aaaaaaaa...`.

## Async-Grenze: warum `traceparent` im Event-Body

Bei der Compensation aus Story 6 schickt Booking ein `CompensationRequested`-
Event per HTTP-POST, das Backend antwortet sofort mit `202 Accepted` und
verarbeitet in einer Goroutine. Beim Goroutine-Start ist der ursprüngliche
HTTP-Request bereits zu Ende — der Trace-Kontext muss daher **explizit als
Property im Event-Body** mitwandern:

```json
{
  "eventId": "ev-...",
  "sagaId": "B-...",
  "bookingId": "F-...",
  "traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
}
```

Der Compensation-Handler liest die Property, baut den Trace-Kontext daraus
und übergibt ihn an die Goroutine — sodass die `phase=processing`- und
`phase=done`-Logzeilen dieselbe `trace_id` tragen wie die ursprüngliche
Buchung.

## Hintergrund zur Implementierung

- **Geteilte Bibliothek:** `services/shared/tracing/` (Parse, New,
  Middleware, Inject, Logger) — ~100 Zeilen, sprachunabhängiger Stil.
- **Pflicht-Pfad** der Story ist tool-agnostisch: keine OTel-Library, kein
  Jaeger/Tempo. `grep <trace-id>` ist das Demo-Werkzeug.
- **Bonus-Pfad** (in `docs/instructions/distributed-tracing.md`): Wie man
  OpenTelemetry + Jaeger einbindet, Sampling betreibt und die richtigen
  Markt-Tools wählt.
