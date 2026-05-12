# Story 9: Ohne Unterbrechung

**Thema:** Downtimeless Deployment
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Der Booking-Service soll aktualisiert werden können, ohne dass laufende Buchungen unterbrochen werden. Kunden sollen von Updates nichts mitbekommen.

## User Story

Als **Betriebsteam**
möchte ich **neue Versionen des Booking-Service deployen können, ohne dass es zu Ausfallzeiten kommt**,
damit **Kunden jederzeit buchen können und wir häufiger deployen können**.

## Akzeptanzkriterien

- [ ] Mindestens 2 Instanzen des Booking-Service laufen
- [ ] Rolling Update: Neue Version wird schrittweise ausgerollt
- [ ] Laufende Requests werden abgeschlossen bevor eine Instanz stoppt (Graceful Shutdown)
- [ ] Der Health-Check unterscheidet zwischen Liveness und Readiness
- [ ] Neue Instanzen erhalten erst Traffic, wenn sie ready sind
- [ ] Bei fehlgeschlagenem Deployment wird automatisch zurückgerollt
- [ ] Während des Deployments sind alle Funktionen verfügbar

## Technische Hinweise

- **Graceful Shutdown:**
  - Keine neuen Requests annehmen
  - Laufende Requests abschließen
  - Ressourcen freigeben
  - Aus Service Registry deregistrieren
- **Readiness vs. Liveness:**
  - Liveness: "Lebt der Prozess?" → Neustart bei Failure
  - Readiness: "Kann der Service Requests verarbeiten?" → Kein Traffic bei Failure
- **Rolling Update Strategie:**
  - maxUnavailable: 0 (nie weniger als die gewünschte Anzahl)
  - maxSurge: 1 (eine zusätzliche Instanz während Update)
- **Testing:** Mit Last-Generator während des Deployments testen

## Bonus (optional)

- Implementiere Blue-Green Deployment als Alternative
- Füge Canary Releases hinzu (nur 10% Traffic auf neue Version)
