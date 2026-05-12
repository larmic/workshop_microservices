# Story 10: Den roten Faden im Log

**Thema:** Distributed Tracing
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Im aktuellen Setup zeigen die Service-Logs jede Anfrage isoliert: der
Booking-Service loggt seine Saga, Flight/Hotel/Car loggen ihre lokalen
Buchungen — aber niemand kann ein POST `/booking/bookings` über alle
beteiligten Services hinweg als **eine zusammenhängende Operation**
verfolgen. Korrelation passiert heute nur über Timestamps und Augenmaß.
Das wird mit jeder zusätzlichen Komponente schmerzhafter.

Ein **Trace ID** pro eingehender Anfrage, der durch alle nachgelagerten
Aufrufe propagiert wird, macht den Ablauf in Logs und Tools wie Jaeger
oder Tempo sichtbar — gerade die Saga aus Story 5 ist dafür ein
Lehrbuchbeispiel: eine Operation, sechs HTTP-Calls (3 Forward + bis zu 3
Compensation), drei Services.

## User Story

Als **Entwickler:in im Betrieb**
möchte ich **einen einzelnen Geschäftsvorgang über Service-Grenzen
hinweg in Logs und Traces verfolgen können**,
damit **ich Fehlerursachen, Latenzen und Saga-Verläufe gezielt
analysieren kann, ohne Logs nach Zeitstempeln zusammenzupuzzeln**.

## Akzeptanzkriterien

- [ ] Bei jedem eingehenden Request am Booking-Service wird ein
  Trace-Kontext erzeugt oder aus dem Header übernommen
  (W3C Trace Context: `traceparent`)
- [ ] Der Trace-Kontext wird auf jedem ausgehenden HTTP-Call
  weitergereicht (Flight, Hotel, Car und deren Kompensation)
- [ ] Jede Logzeile aller Services enthält die Trace-ID, sodass
  `grep <trace-id>` den vollständigen Ablauf zeigt
- [ ] Die Saga aus Story 5 ist als ein Trace mit einer Span pro
  Forward- und Compensation-Schritt sichtbar
- [ ] Ein Trace-Backend (Jaeger oder Tempo) ist im
  `docker-compose` integriert und über das Gateway erreichbar
- [ ] Die OpenAPI-Dokumentation aller Services beschreibt den
  `traceparent`-Header

## Technische Hinweise

- **Standard:** [W3C Trace Context](https://www.w3.org/TR/trace-context/)
  ist heute das Default-Wire-Format. `traceparent`-Header trägt
  Trace-ID + Span-ID + Flags.
- **Library:** [OpenTelemetry für Go](https://opentelemetry.io/docs/languages/go/)
  bietet Auto-Instrumentation für `net/http` (sowohl Server als auch
  Client). Damit ist die Propagation in beide Richtungen ohne manuelles
  Header-Handling abgedeckt.
- **Trace-Backend für den Workshop:**
  - **Jaeger** — älterer Klassiker, eigenständig, eigene UI
  - **Grafana Tempo** — passt zu Grafana/Loki-Stack
  - Für 60 Minuten reicht Jaeger als All-in-One-Container
- **Log-Korrelation:** Die Trace-ID muss zusätzlich in jede Logzeile.
  Mit `slog` (statt `log`) lässt sich das strukturiert beifügen
  (`slog.String("trace_id", ...)`) — mit `log.Printf` muss man die ID
  manuell in jedes Format-String einsetzen.
- **Saga + Tracing:** Eine Saga ergibt naturgemäß einen Trace mit
  geschachtelten Spans:
  ```
  span: POST /booking/bookings           ─┐
    span: flight POST /bookings           │  Forward
    span: hotel POST /bookings  ✗  failed │
    span: flight DELETE /bookings/{id}   ─┘  Compensation
  ```

## Bonus (optional)

- Traces nicht nur über HTTP, sondern auch über Consul-Lookups
  (zeigt Service-Discovery-Latenz pro Request)
- Sampling-Strategie diskutieren: bei produktiver Last nicht jeder
  Request, sondern z.B. nur fehlerhafte Sagas zu 100 %, erfolgreiche
  zu 1 %
- Connection zum Saga-Status-Endpoint aus Story 5: enthält die
  Saga-Antwort die Trace-ID, kann der Kunde sie an den Support
  weiterreichen

## Bezug zum bisherigen Workshop

- **Story 1 (Logs):** Logs werden mit Trace-IDs angereichert — die
  „eine Logzeile pro Aktion"-Praxis bekommt damit ihren roten Faden.
- **Story 3/4 (CB, Bulkhead):** Reaktionszeiten und Reject-Rates pro
  Backend werden im Trace sichtbar — kein Rätselraten mehr, ob ein
  langsamer Request vom CB, vom Bulkhead oder vom Backend selbst kommt.
- **Story 5 (Saga):** Erst mit Tracing sieht man die Saga als
  zusammenhängende Operation, nicht als Fragmente in mehreren
  Service-Logs.
- **Story 6 (Choreography-Saga):** Tracing über asynchrone Events ist
  nicht-trivial — Trace-ID muss als Property auf das Event mitwandern,
  sonst zerfällt der Trace an der Bus-Grenze.
