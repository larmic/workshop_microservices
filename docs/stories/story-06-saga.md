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
- [ ] Jeder Service bietet einen Stornieren/Kompensieren-Endpoint an (`DELETE /bookings/{id}`), idempotent: zweiter Aufruf, gleicher Zustand, wieder `204`
- [ ] Der aktuelle Status einer Buchung kann abgefragt werden (PENDING, COMPLETED, COMPENSATING, FAILED)
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
- **Rückblick auf Story 2:** Die Storno-Diskussion vom Flipchart kommt hier als Kompensation zurück. `DELETE /bookings/{id}` muss idempotent sein, sobald der Orchestrator wiederholt (Bonus: Retry) oder jemand manuell nachstößt (siehe [questions/story2.md](../questions/story2.md), Frage 4)

## Bonus (optional)

- Implementiere parallele Ausführung unabhängiger Schritte
- Füge Timeout-Handling hinzu (Was passiert, wenn ein Service nicht antwortet?)
- Kompensation mit Retry-Logik durchführen (Kompensation darf nicht fehlschlagen) — siehe Diskussionsfragen 1, 2 und 5 in `docs/questions/story6.md`
- Saga-Status persistieren (für Recovery nach Absturz) — kein Saga-spezifisches Thema, sondern generelle Orchestrator-Robustheit, siehe Diskussionsfrage 7 in `docs/questions/story6.md`
