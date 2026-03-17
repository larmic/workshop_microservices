# Story 6: Alles oder nichts - aber richtig

**Thema:** Saga Pattern
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Eine Reisebuchung umfasst Flug, Hotel und Mietwagen. Wenn die Hotelbuchung fehlschlägt, nachdem der Flug bereits gebucht wurde, muss der Flug storniert werden. Klassische Datenbank-Transaktionen funktionieren nicht über Service-Grenzen hinweg.

## User Story

Als **Kunde**
möchte ich **eine Komplettbuchung (Flug + Hotel + Mietwagen) durchführen, die entweder vollständig erfolgreich ist oder komplett zurückgerollt wird**,
damit **ich nicht mit einer unvollständigen Buchung dastehe**.

## Akzeptanzkriterien

- [ ] Eine Buchungsanfrage für Flug, Hotel und Mietwagen kann gestellt werden
- [ ] Die Services werden nacheinander aufgerufen (Orchestration-Saga)
- [ ] Bei Fehler in einem Schritt werden alle vorherigen Schritte kompensiert (Rollback)
- [ ] Jeder Service bietet einen Stornieren/Kompensieren-Endpoint an
- [ ] Der Saga-Status wird persistiert (für Recovery nach Absturz)
- [ ] Der aktuelle Status einer Buchung kann abgefragt werden (PENDING, COMPLETED, COMPENSATING, FAILED)
- [ ] Kompensation wird mit Retry-Logik durchgeführt (Kompensation darf nicht fehlschlagen)

## Technische Hinweise

- **Saga-Varianten:**
  - **Choreography:** Services kommunizieren über Events (dezentral)
  - **Orchestration:** Ein Koordinator steuert den Ablauf (zentral) - empfohlen für diesen Workshop
- **Saga-Schritte für Reisebuchung:**
  1. Flug buchen
  2. Hotel buchen
  3. Mietwagen buchen
  - Kompensation: Mietwagen stornieren → Hotel stornieren → Flug stornieren
- **Webhook-Callbacks:** Services rufen bei Statusänderung einen Callback-URL auf
- **API-Endpunkte (Backend-Services):**
  - `POST /bookings` - Buchung erstellen
  - `DELETE /bookings/{id}` - Buchung stornieren (Kompensation)

## Bonus (optional)

- Implementiere parallele Ausführung unabhängiger Schritte
- Füge Timeout-Handling hinzu (Was passiert, wenn ein Service nicht antwortet?)
