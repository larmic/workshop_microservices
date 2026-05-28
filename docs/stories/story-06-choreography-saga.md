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

## Akzeptanzkriterien (Pflicht, Schmalspur-Variante)

- [ ] Bei Fehlschlag eines Saga-Schritts publiziert der Booking-Service ein `CompensationRequested`-Event (statt synchron `DELETE` aufzurufen)
- [ ] Flight-, Hotel- und Car-Service nehmen das Event entgegen, antworten **sofort mit `202 Accepted`** und führen die Stornierung asynchron in einer Goroutine / Worker-Task aus (fire-and-forget aus Sicht des Senders)
- [ ] Booking setzt den Step nach erfolgreichem Event-Dispatch auf `COMPENSATED` und die Saga direkt auf `FAILED`. Booking erwartet **kein Reply** und wartet **nicht** auf den fachlichen Rollback
- [ ] Das Konzept der Idempotenz wird gezeigt: `eventId` wird pro Event eindeutig erzeugt und mitgesendet (persistente Speicherung im Backend ist Bonus)
- [ ] Der Unterschied zwischen Orchestration (Story 5) und Choreography (Story 6) wird im Code-Aufbau erkennbar und in der README des Booking-Service kurz reflektiert

## Akzeptanzkriterien (Bonus, Production-Reife)

Diese Punkte sind bewusst nicht Pflicht. Sie sind das, was die Schmalspur-Variante strukturell *nicht* leistet, und sind Diskussionsstoff für das Recap. Wer mag, kann sie als Vertiefung umsetzen:

- [ ] Reply-Events: Backends bestätigen mit `BookingCancelled` oder `CancellationFailed`
- [ ] Booking konsumiert die Reply-Events und führt den Saga-Status nach (`COMPENSATING` → `FAILED` oder `COMPENSATION_INCOMPLETE`)
- [ ] Timeout-Erkennung im Booking-Service: Bleibt ein erwartetes Reply nach X Sekunden aus, wird die Saga als `STUCK` markiert und ein Alarm-Log geschrieben
- [ ] Persistente Idempotenz: doppelt zugestellte Events werden über die `eventId` im Backend erkannt und ignoriert (Dedup-Tabelle pro `eventId`)
- [ ] Feature-Flag zur Laufzeit zwischen synchroner Kompensation (Story 5) und asynchroner Kompensation (Story 6)

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
- **Fire-and-Forget (Pflicht-Variante):** Booking ist nach „Event raus" fertig. Backend antwortet sofort `202 Accepted` und macht den eigentlichen Rollback asynchron. Booking weiß *nicht*, ob der Rollback erfolgreich war. Das ist die bewusst fragile Schmalspur, die zeigt, was Eventing-ohne-Broker strukturell nicht leistet (Diskussion im Recap).
- **Reply-Pattern (Bonus):** Wer die volle Variante will, baut den Reply-Channel dazu: Backend POSTet später ein `BookingCancelled` / `CancellationFailed` an einen Booking-Endpoint, Booking führt den Saga-Status nach. Plus Timeout-Erkennung: bleibt ein Reply aus, geht die Saga auf `STUCK`. Das ist *kein* Fire-and-Forget mehr, sondern asynchrones Request-Reply über zwei Webhook-Richtungen.
- **Was Booking aufgibt, was bleibt (Pflicht-Variante):**
  - **Aufgegeben:** direkte Verantwortung für die Ausführung der Stornierung. Booking weiß auch *nicht mehr*, ob sie geklappt hat
  - **Bleibt:** Verantwortung für den **Saga-Status** gegenüber dem Kunden (Booking gibt die Endaussage „FAILED" zurück, sobald die Events raus sind), Logging des Event-Dispatches
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
