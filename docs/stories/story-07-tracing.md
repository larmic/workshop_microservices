# Story 7: Den roten Faden im Log

**Thema:** Distributed Tracing
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Im aktuellen Setup zeigen die Service-Logs jede Anfrage isoliert: der
Booking-Service loggt seine Saga, Flight/Hotel/Car loggen ihre lokalen
Buchungen — aber niemand kann ein POST `/booking/bookings` über alle
beteiligten Services hinweg als **eine zusammenhängende Operation**
verfolgen. Korrelation passiert heute nur über Timestamps und Augenmaß.
Das wird mit jeder zusätzlichen Komponente schmerzhafter — und spätestens
bei der Saga aus Story 5/6 (eine Operation, bis zu sechs HTTP-Calls über
drei Backend-Services) wird Debugging zur Detektivarbeit.

Eine **Trace-ID** pro eingehender Anfrage, die durch alle nachgelagerten
Aufrufe propagiert und in jede Logzeile geschrieben wird, löst das
Problem auch ohne ein eigenes Tracing-Backend: ein simples
`grep <trace-id>` über die Compose-Logs zeigt den vollständigen Ablauf
einer Buchung.

**Wer initiiert den Trace?** Bewusst nur der Booking-Service als
Entry-Point. Flight, Hotel und Car sind **passive Trace-Empfänger** —
sie verlängern einen eintreffenden `traceparent`-Header, erzeugen aber
**niemals selbst** einen neuen. Dadurch entsteht der "rote Faden" erst,
sobald Story 7 ihn im Booking-Service zieht; in Stories 1–6 (ohne
Tracing-Propagation) bleiben die Logs von Flight/Hotel/Car ohne
`trace_id` — der Kontrast macht den Effekt sichtbar. Dieses Muster
spiegelt auch die reale Welt wider: Trace-Initiierung gehört an den
Entry-Point (API-Gateway, Public-Service), nicht in jeden
Downstream-Hop.

## User Story

Als **Entwickler:in im Betrieb**
möchte ich **einen einzelnen Geschäftsvorgang über Service-Grenzen
hinweg in den Logs verfolgen können**,
damit **ich Fehlerursachen, Latenzen und Saga-Verläufe gezielt
analysieren kann, ohne Logs nach Zeitstempeln zusammenzupuzzeln**.

## Akzeptanzkriterien

- [ ] Bei jedem eingehenden Request am Booking-Service wird ein
  Trace-Kontext erzeugt oder aus dem Header übernommen
  (W3C Trace Context: `traceparent`)
- [ ] Der Trace-Kontext wird auf jedem ausgehenden HTTP-Call
  weitergereicht (Flight, Hotel, Car und deren Kompensation)
- [ ] Auch über die **asynchronen Compensation-Events** aus Story 6 wird
  der Trace-Kontext mitgeführt (als Event-Property), sodass die
  Stornierungs-Logzeile dieselbe Trace-ID trägt wie die ursprüngliche
  Buchung
- [ ] Jede Logzeile aller Services enthält die Trace-ID, sodass
  `docker compose logs | grep <trace-id>` den vollständigen Ablauf
  zeigt — Forward-Buchung **und** etwaige Kompensation
- [ ] Logs sind strukturiert (JSON), damit `trace_id` ein eigenes Feld
  ist und nicht in einem freien Textstring versteckt liegt
- [ ] Die OpenAPI-Dokumentation des Booking-Services beschreibt den
  `traceparent`-Header

## Technische Hinweise

- **Standard:** [W3C Trace Context](https://www.w3.org/TR/trace-context/)
  ist heute das Default-Wire-Format. Der `traceparent`-Header trägt
  Version + Trace-ID + Span-ID + Flags und ist nur 55 Zeichen lang.
- **Selbst bauen für den Workshop:** Parsing, Generierung und Propagation
  sind in ca. 100 Zeilen Code abgehandelt — analog zu Story 3/4 lernt
  man den Mechanismus dadurch wirklich. Eine fertige Library (z.B.
  OpenTelemetry) würde den Lerneffekt verstecken und ist außerdem stark
  Sprach- und Tool-abhängig.
- **Span-Konzept:** Pro Outbound-Hop eine neue Span-ID generieren
  (gleiche Trace-ID). Einen vollwertigen Span-Baum mit Parent-Child-
  Verknüpfung zu bauen ist für 60 Minuten zu viel — bewusst weglassen,
  in den Trainer-Notizen wird der Vollumfang erklärt.
- **Strukturiertes Logging:** Mit `slog` (statt `log`) lässt sich
  `trace_id` als eigenes Feld anhängen
  (`slog.String("trace_id", ...)`) — mit `log.Printf` müsste man die ID
  manuell in jeden Format-String einsetzen. JSON-Output macht das spätere
  Filtern (z.B. mit `jq`) trivial.
- **Saga + Tracing:** Eine Saga ergibt naturgemäß einen Vorgang mit
  geschachtelten Aufrufen:
  ```
  POST /booking/bookings           ─┐
    flight POST /bookings           │  Forward
    hotel  POST /bookings  ✗ failed │
    flight DELETE /bookings/{id}   ─┘  Compensation
  ```
  Mit derselben Trace-ID in allen Logzeilen wird der Vorgang sofort
  sichtbar.
- **Async-Grenze (Story 6):** Bei Choreography muss die Trace-ID
  **explizit als Property** auf das Event mitwandern — der HTTP-Header
  geht beim Übergang in die Worker-Goroutine verloren. Im Event-Body
  z.B. ein Feld `"traceparent": "00-..."`.
- **Zwei Middlewares, zwei Rollen:** Die Library stellt bewusst zwei
  Einstiegspunkte bereit. `Middleware` erzeugt einen Trace, falls keiner
  reinkommt — das ist die Entry-Point-Variante (Booking-Service).
  `Propagate` übernimmt nur einen bereits vorhandenen Trace und legt
  sonst nichts an — das ist die Downstream-Variante (Flight/Hotel/Car).
  Diese Trennung sorgt dafür, dass die Trace-ID in den
  Downstream-Logs erst auftaucht, wenn der Entry-Point sie aktiv
  propagiert — und nicht zufällig vom Service selbst generiert wurde.

## Bonus (optional)

- **Trace-Backend einbinden:** Jaeger oder Grafana Tempo als
  All-in-One-Container im Compose ergänzen und OpenTelemetry-SDK
  einbinden. Das ist sprach- und tool-abhängig (Go-OTel ≠ Java-OTel
  ≠ Node-OTel) und für den Workshop bewusst optional. Die schöne UI mit
  Span-Bäumen ist beeindruckend — der Lerneffekt steckt aber im
  selbstgebauten Teil.
- **Saga-Status mit Trace-ID:** Die Response des Booking-Endpunkts
  zurück­liefert die Trace-ID, sodass ein Kunde sie an den Support
  weitergeben kann.
- **Sampling-Strategie diskutieren:** In Produktion will man nicht
  jeden Request tracen — fehlerhafte Sagas zu 100 %, erfolgreiche z.B.
  zu 1 %.
- **Consul-Lookups als eigene Events loggen:** Zeigt die
  Service-Discovery-Latenz pro Request.

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
  nicht-trivial — die Trace-ID muss als Property auf das Event
  mitwandern, sonst zerfällt der Vorgang an der Bus-Grenze. Genau dieser
  Fall wird in Story 7 implementiert.
