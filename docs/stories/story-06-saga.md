# Story 6: Alles oder nichts

**Thema:** Saga Pattern
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Eine Reisebuchung umfasst Flug, Hotel und Mietwagen. Wenn die Hotelbuchung fehlschlägt, nachdem der Flug bereits gebucht wurde, muss der Flug storniert werden. Klassische Datenbank-Transaktionen funktionieren nicht über Service-Grenzen hinweg.

## User Story

Als **Kunde**
möchte ich **eine Komplettbuchung (Flug + Hotel + Mietwagen) durchführen, die am Ende entweder vollständig steht oder nichts mehr gebucht ist**,
damit **ich nicht mit einer halben Reise dastehe**.

## Akzeptanzkriterien

- [ ] Eine Buchungsanfrage für Flug, Hotel und Mietwagen kann gestellt werden
- [ ] Die Services werden nacheinander aufgerufen (Orchestration-Saga)
- [ ] Bei Fehler in einem Schritt werden alle vorherigen Schritte kompensiert, in umgekehrter Reihenfolge
- [ ] Kompensiert wird über `DELETE /bookings/{id}`. Flight, Hotel und Car bieten das schon an, idempotent: zweiter Aufruf, gleicher Zustand, wieder `204`
- [ ] Der aktuelle Status einer Buchung kann abgefragt werden (PENDING, COMPLETED, COMPENSATING, FAILED), in der Referenz über `GET /booking/bookings/{sagaId}`
- [ ] Optional: Retry-Logik für die Kompensation, Timeout-Handling, Persistenz für Recovery

## Technische Hinweise

- **Orchestration:** Ihr baut sie: Booking ruft nacheinander auf und kompensiert selbst. Das Storno per Event statt per Aufruf spielen wir in Story 7.
- **Saga-Schritte für Reisebuchung:**
  1. Flug buchen
  2. Hotel buchen
  3. Mietwagen buchen
  - Kompensation: Mietwagen stornieren → Hotel stornieren → Flug stornieren
- **API-Endpunkte (Backend-Services):**
  - `POST /bookings` - Buchung erstellen
  - `DELETE /bookings/{id}` - Buchung stornieren (Kompensation)
- **Fehler auslösen:** Im Dashboard einen Service auf „Fehler“ stellen (Chaos-Modus in der Service-Karte) oder mit `carId` `C9` buchen: Die ist im Car-Service reserviert, Car antwortet mit 409.
- **Antwort im Fehlerfall:** Scheitert die Saga, antwortet `POST /booking/bookings` mit einem Fehler-Status und der Saga im Body. Die Referenz nimmt `503`.
- **Anzeige im Dashboard:** Damit die Saga-Karte auch bei „Mein Service“ etwas zeigt, fragt das Dashboard zwei Endpunkte an eurem Service ab:
  - `GET /admin/sagas` liefert die Liste, neueste zuerst: `[ { sagaId, status, customerName, createdAt, failedAt, steps: [ { service, bookingId, status } ] } ]`. `service` ist `flight`, `hotel` oder `car`, der Schritt-Status `BOOKED`, `COMPENSATED` oder `COMPENSATION_FAILED`.
  - `POST /admin/sagas-reset` leert die Liste und antwortet mit `204`.
- **Rückblick auf Story 2:** Die Storno-Diskussion vom Flipchart kommt hier als Kompensation zurück. `DELETE /bookings/{id}` muss idempotent sein, sobald der Orchestrator wiederholt (Bonus: Retry) oder jemand manuell nachstößt

### Story 3 noch nicht fertig?

Die Saga ersetzt den `POST /booking/bookings` aus Story 3, ihr braucht ihn also nicht als Vorlage. Fehlt die Service Discovery, ruft Flight, Hotel und Car über die Env-URLs aus Story 1 auf. Die Kriterien dieser Story gelten trotzdem.

## Bonus (optional)

- Implementiere parallele Ausführung unabhängiger Schritte
- Füge Timeout-Handling hinzu (Was passiert, wenn ein Service nicht antwortet?)
- Kompensation mit Retry-Logik durchführen (Kompensation darf nicht fehlschlagen) — Thema im Recap
- Saga-Status persistieren (für Recovery nach Absturz) — kein Saga-spezifisches Thema, sondern generelle Orchestrator-Robustheit, ebenfalls Thema im Recap
