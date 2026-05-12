# Story 6: Die Saga wird leise – Kompensation via Events

**Thema:** Event-Driven Architecture, Choreography-Saga
**Zeitrahmen:** ca. 60 Minuten

## Kontext

In Story 5 trägt der Booking-Service die volle Verantwortung für die Kompensation: Er ruft synchron `DELETE /bookings/{id}` gegen jeden zuvor erfolgreich aufgerufenen Backend-Service auf, wartet auf jede Antwort und behandelt Fehler in seinem eigenen Code. Damit ist Booking nicht nur Orchestrator des Happy Paths, sondern auch Single Point of Responsibility für jede Stornierung. Skaliert ein Backend träge oder ist es kurz nicht erreichbar, blockiert Booking — und wird selbst zum Bottleneck.

Fachlich ist die Stornierung aber Aufgabe des jeweiligen Backends: Wer eine Buchung anlegen kann, muss sie auch zurücknehmen können — ohne dass ein Orchestrator daneben steht. Lösung: Booking publiziert ein **Event** („Kompensation erforderlich"), die Backend-Services abonnieren und reagieren eigenständig. Wir wechseln damit von **Orchestration** (Story 5) zu **Choreography** — derselbe fachliche Ablauf, andere Verantwortungsverteilung.

## User Story

Als **System**
möchte ich **bei einer fehlgeschlagenen Buchung die Kompensation asynchron an die zuständigen Backend-Services delegieren**,
damit **der Booking-Service nicht für deren Verfügbarkeit haften muss und die Backends ihre eigene Stornierungslogik kapseln können**.

## Akzeptanzkriterien

- [ ] Bei Fehlschlag eines Saga-Schritts publiziert der Booking-Service ein `CompensationRequested`-Event (statt synchron `DELETE` aufzurufen)
- [ ] Flight-, Hotel- und Car-Service abonnieren das Event und führen ihre lokale Stornierung selbst aus
- [ ] Services bestätigen das Ergebnis mit einem Reply-Event (`BookingCancelled` oder `CancellationFailed`)
- [ ] Booking konsumiert die Reply-Events und führt den Saga-Status nach (`COMPENSATING` → `FAILED` oder `COMPENSATION_INCOMPLETE`)
- [ ] **Timeout-Erkennung im Booking-Service:** Bleibt ein erwartetes Reply-Event nach X Sekunden aus, wird die Saga als `STUCK` markiert und ein Alarm-Log geschrieben
- [ ] Idempotenz: Doppelt zugestellte Events werden über die `eventId` erkannt und ignoriert
- [ ] Der Unterschied zwischen Orchestration (Story 5) und Choreography (Story 6) wird im Code-Aufbau erkennbar und in der README des Booking-Service kurz reflektiert

## Technische Hinweise

- **Webhook-basiertes Eventing** (kein separater Message-Broker, siehe `stories.md`): Services registrieren sich beim Event-Publisher
  - `POST /webhooks/subscribe` mit `{ "eventType": "CompensationRequested", "callbackUrl": "http://..." }`
- **Event-Struktur:**
  ```json
  {
    "eventType": "CompensationRequested",
    "eventId": "uuid",
    "timestamp": "2026-05-12T10:30:00Z",
    "sagaId": "S-12345",
    "payload": {
      "service": "FLIGHT",
      "bookingId": "F-7c1a9f"
    }
  }
  ```
- **Reply-Pattern:** Booking ist *nicht* fertig nach „Event raus". Der Kunde will eine Endaussage („gebucht / storniert / steckt fest"). Backend antwortet daher mit einem Reply-Event; Booking konsumiert das und führt den Saga-Status nach.
- **Was Booking aufgibt — und was nicht:**
  - **Aufgegeben:** direkte Verantwortung für die Ausführung der Stornierung, eigene Retry-Schleifen
  - **Bleibt:** Verantwortung für den **Gesamtstatus** der Saga gegenüber dem Kunden, Timeout-Erkennung, Operator-Eskalation
- **Vergleich Orchestration ↔ Choreography:**
  - Orchestration: Wissen zentral, Bug-Lokalisierung einfach, Kopplung höher
  - Choreography: Wissen verteilt, Backends entkoppelt, „verteilter Monolith"-Risiko bei schlechtem Schnitt

## Diskussions-Anker

- Warum nutzen wir Webhooks statt Kafka/RabbitMQ? Was würde sich ändern?
- Was passiert, wenn ein Reply-Event nie ankommt? (Anschluss an Story 5, Frage 5 zu Saga-Beobachtbarkeit)
- Wo wandert das Saga-Wissen jetzt hin — und wann wird Choreography zum verteilten Monolithen?
- Welcher Teil der Story-5-Implementierung fällt komplett weg, welcher bleibt unverändert?

## Bonus (optional)

- **Dead-Letter-Behandlung:** Was passiert mit Events, die nach N Versuchen nicht zugestellt werden konnten?
- **Mischbetrieb über Feature-Flag:** Booking kann zur Laufzeit zwischen synchroner Kompensation (Story 5) und asynchroner Kompensation (Story 6) umschalten — schöner Showcase im Workshop
- **Auch der Happy Path über Events:** Nicht nur Kompensation, sondern auch das Forward-Booking als Event-Choreography — zeigt, wie weit man Choreography treiben kann
