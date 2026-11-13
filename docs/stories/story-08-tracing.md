# Story 8: Den roten Faden im Log

**Thema:** Distributed Tracing
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Im aktuellen Setup zeigen die Service-Logs jede Anfrage isoliert: der
Booking-Service loggt seine Saga, Flight/Hotel/Car loggen ihre lokalen
Buchungen — aber niemand kann ein POST `/booking/bookings` über alle
beteiligten Services hinweg als **eine zusammenhängende Operation**
verfolgen. Korrelation passiert heute nur über Timestamps und Augenmaß.
Das wird mit jeder zusätzlichen Komponente schmerzhafter — und spätestens
bei eurer Saga aus Story 6 (eine Operation, bis zu fünf HTTP-Calls über
drei Backend-Services: drei Buchungen, zwei DELETE) wird Debugging zur
Detektivarbeit.

Eine **Trace-ID** pro eingehender Anfrage, die durch alle nachgelagerten
Aufrufe propagiert und in jede Logzeile geschrieben wird, löst das
Problem auch ohne ein eigenes Tracing-Backend: ein simples
`grep <trace-id>` über die Compose-Logs zeigt den vollständigen Ablauf
einer Buchung.

**Wer initiiert den Trace?** Bewusst nur der Booking-Service als
Entry-Point. Flight, Hotel und Car sind **passive Trace-Empfänger** —
sie übernehmen einen eintreffenden `traceparent`-Header und starten per
HTTP bewusst keinen eigenen. Dadurch entsteht der "rote Faden" erst,
sobald Story 8 ihn im Booking-Service zieht; in Stories 1–7 (ohne
Tracing-Propagation) bleiben die Logs von Flight/Hotel/Car ohne
`trace_id` — der Kontrast macht den Effekt sichtbar. Das ist eine
Workshop-Entscheidung: In OpenTelemetry startet jeder Service ohne
eingehenden Header einen Trace. Tabu ist nur, einen eingehenden zu
überschreiben.

Eine Ausnahme gibt es: Kommt ein Compensation-Event ohne
`traceparent`-Property an, starten die Backends dafür einen neuen
Trace. Das betrifft nur den Event-Bonus und ist dort genau der Effekt,
den man sehen soll.

## User Story

Als **Entwickler:in im Betrieb**
möchte ich **einen einzelnen Geschäftsvorgang über Service-Grenzen
hinweg in den Logs verfolgen können**,
damit **ich Fehlerursachen, Latenzen und Saga-Verläufe gezielt
analysieren kann, ohne Logs nach Zeitstempeln zusammenzupuzzeln**.

## Akzeptanzkriterien

Ausgangspunkt ist euer Booking-Service aus Story 6: orchestrierte Saga
Flight → Hotel → Car, Kompensation per `DELETE /bookings/{id}`.
Events, Circuit Breaker oder Bulkhead braucht ihr nicht.

- [ ] Bei jedem eingehenden Request am Booking-Service wird ein
  Trace-Kontext erzeugt oder aus dem Header übernommen
  (W3C Trace Context: `traceparent`)
- [ ] Die Response trägt den effektiven `traceparent` zurück, damit
  Client, Dashboard und Support die Trace-ID kennen
- [ ] Der Trace-Kontext wird auf jedem ausgehenden HTTP-Call
  weitergereicht: Buchung (`POST /bookings` an Flight, Hotel, Car) und
  Kompensation (`DELETE /bookings/{id}`). Pro Hop eine neue Span-ID, die
  Trace-ID bleibt
- [ ] Jede Logzeile eures Booking-Service enthält die Trace-ID. Flight,
  Hotel und Car übernehmen einen eintreffenden `traceparent` schon und
  loggen ihn
- [ ] Logs sind strukturiert (JSON), damit `trace_id` ein eigenes Feld
  ist und nicht in einem freien Textstring versteckt liegt
- [ ] Ein Filter zeigt den vollständigen Ablauf einer Buchung,
  Buchung **und** Kompensation:
  `docker compose logs --no-log-prefix booking-custom flight hotel car | jq -cR 'fromjson? | select(.trace_id=="<trace-id>")'`
  (ohne jq: `docker compose logs | grep <trace-id>`)
  Beide Befehle laufen im Ordner `services/`. Meldet Compose
  `invalid compose project`, fehlt in eurer `.env` die Zeile `COMPOSE_FILE`
  aus `.env.example`.
- [ ] Die OpenAPI-Dokumentation des Booking-Services beschreibt den
  `traceparent`-Header (Request und Response)

Tipp: Eine Kompensation löst ihr im Dashboard mit „Buchung mit kaputtem
Mietwagen" aus. Die `carId` `C9` ist im Car-Service reserviert, Car
antwortet mit 409, eure Saga storniert Flight und Hotel.

### Story 6 noch nicht fertig?

Die Kriterien gelten für die Calls, die ihr habt, mindestens
`GET /offers` oder ein `POST /bookings` an ein Backend. Die Kompensation
zeigt der Trainer an der Referenz (Toggle „Referenz"). Fehlt der
Service ganz, reicht ein Stub mit einem Endpoint, der Flight aufruft.

## Technische Hinweise

- **Standard:** [W3C Trace Context](https://www.w3.org/TR/trace-context/)
  ist heute das Default-Wire-Format. Der `traceparent`-Header trägt
  Version + Trace-ID + Span-ID + Flags und ist nur 55 Zeichen lang.
- **Selbst bauen für den Workshop:** Parsing, Generierung und Propagation
  sind in gut 120 Zeilen Code abgehandelt — analog zu Story 4/5 lernt
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
    flight POST /bookings           │  Buchung
    hotel  POST /bookings  ✗ failed │
    flight DELETE /bookings/{id}   ─┘  Compensation
  ```
  Mit derselben Trace-ID in allen Logzeilen wird der Vorgang sofort
  sichtbar.
- **Zwei Middlewares, zwei Rollen:** Die Library stellt bewusst zwei
  Einstiegspunkte bereit. `Middleware` erzeugt einen Trace, falls keiner
  reinkommt — das ist die Entry-Point-Variante (Booking-Service).
  `Propagate` übernimmt nur einen bereits vorhandenen Trace und legt
  sonst nichts an — das ist die Downstream-Variante (Flight/Hotel/Car).
  Diese Trennung sorgt dafür, dass die Trace-ID in den
  Downstream-Logs erst auftaucht, wenn der Entry-Point sie aktiv
  propagiert — und nicht zufällig vom Service selbst generiert wurde.

## Bonus (optional)

- **Detektiv-Runde „Die kaputte Buchung" (10 Minuten):** Im Dashboard
  schickt „Rush Hour" fünf Buchungen gleichzeitig, eine davon mit
  `carId` `C9`. Findet ihre Trace-ID in den Logs (Car loggt
  `car not available`) und tragt sie ins Lösungsfeld ein, das Dashboard
  stoppt die Zeit. Danach mit dieser ID den ganzen Ablauf filtern.
- **Kompensation als Event (wie in der Story-7-Referenz):** Statt
  `DELETE` schickt ihr `POST /events/compensation` mit
  `{eventId, sagaId, bookingId, traceparent}` an Flight, Hotel und Car
  (Antwort 202). Die Backends verarbeiten das Event asynchron und lesen
  den Trace **nur aus der Body-Property**. Probiert es einmal ohne
  Property: Die `request`-Zeile trägt eure ID, `received`, `processing`
  und `done` tragen eine neue. Der Trace zerfällt.

- **Trace-Backend einbinden:** Jaeger oder Grafana Tempo als
  All-in-One-Container im Compose ergänzen und OpenTelemetry-SDK
  einbinden. Das ist sprach- und tool-abhängig (Go-OTel ≠ Java-OTel
  ≠ Node-OTel) und für den Workshop bewusst optional. Die schöne UI mit
  Span-Bäumen ist beeindruckend — der Lerneffekt steckt aber im
  selbstgebauten Teil.
- **Sampling-Strategie diskutieren:** In Produktion will man nicht
  jeden Request tracen — fehlerhafte Sagas zu 100 %, erfolgreiche z.B.
  zu 1 %.
- **Consul-Lookups als eigene Events loggen:** Zeigt die
  Service-Discovery-Latenz pro Request.

## Bezug zum bisherigen Workshop

- **Story 1 (Logs):** Logs werden mit Trace-IDs angereichert — die
  „eine Logzeile pro Aktion"-Praxis bekommt damit ihren roten Faden.
- **Story 4/5 (CB, Bulkhead):** In der Referenz tragen CB- und
  Bulkhead-Rejects die `trace_id`. Bei euch wäre das die Stelle, an der
  das Board aus Story 4 und die Stühle aus Story 5 geloggt würden.
- **Story 6 (Saga):** Erst mit Tracing sieht man die Saga als
  zusammenhängende Operation, nicht als Fragmente in mehreren
  Service-Logs.
- **Story 7 (Choreography-Saga):** Story 7 war ein Würfelspiel mit
  Zurufen. Ein Zuruf hat keinen HTTP-Header, die Trace-ID muss als
  Property mitwandern, sonst zerfällt der Vorgang an der Bus-Grenze.
  Das ist Recap-Frage 04 und der Event-Bonus. Die Story-8-Referenz
  zeigt eine Umsetzung.
