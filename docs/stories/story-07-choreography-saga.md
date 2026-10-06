# Story 7: Die Saga wird leise

**Thema:** Event-Driven Architecture, Choreography-Saga, als Rollenspiel statt Code
**Zeitrahmen:** ca. 60 Minuten (25 Minuten Spiel, 20 Minuten Störungen, Recap)

## Kontext

In Story 6 trägt der Booking-Service die volle Verantwortung für die Kompensation: Er ruft synchron `DELETE /bookings/{id}` gegen jeden zuvor erfolgreich aufgerufenen Backend-Service auf, wartet auf jede Antwort und behandelt Fehler in seinem eigenen Code. Damit ist Booking nicht nur Orchestrator des Happy Paths, sondern auch Single Point of Responsibility für jede Stornierung. Ist ein Backend kurz nicht erreichbar, blockiert Booking und wird selbst zum Engpass.

Fachlich ist die Stornierung aber Aufgabe des jeweiligen Backends: Wer eine Buchung anlegen kann, muss sie auch zurücknehmen können, ohne dass ein Orchestrator daneben steht. Lösung: Booking publiziert ein **Event** („Kompensation erforderlich"), die Backend-Services reagieren eigenständig. Wir wechseln von **Orchestration** (Story 6) zu **Choreography**: derselbe fachliche Ablauf, andere Verantwortungsverteilung.

Der Unterschied steckt nicht im Code, der ist gegenüber Story 6 zwei Zeilen groß. Er steckt darin, wer wartet und wer am Ende Bescheid weiß. Deshalb spielen wir die Saga, statt sie zu programmieren.

## User Story

Als **System**
möchte ich **bei einer fehlgeschlagenen Buchung die Kompensation asynchron an die zuständigen Backend-Services delegieren**,
damit **der Booking-Service nicht für deren Verfügbarkeit haften muss und die Backends ihre eigene Stornierungslogik kapseln können**.

---

## Teil 1: Die Saga wird leise (ca. 25 Minuten)

Ihr seid die Services. Der Broker ist ein Tisch mit Zetteln.

### Material

- Fünf Rollenkarten: Kunde, Booking, Flight, Hotel, Car
- Ein Tisch oder ein Flipchart als Broker
- Zettel in zwei Farben: Events und Buchungskarten
- Eine Karteikarte pro Beobachter

### Aufbau

1. Vier Leute sind Booking, Flight, Hotel und Car. Eine Person ist der Kunde.
2. Ein Tisch in der Mitte ist der Broker. Events sind Zettel, die dort abgelegt werden.
3. Jedes Backend hat einen Stapel Buchungskarten, die es ausstellen und wieder einziehen kann.
4. Der Kunde bucht eine Reise: Flug, Hotel, Mietwagen. Hotel lehnt jede zweite Buchung ab.
5. Alle anderen beobachten und notieren, wer wann was wusste.

### Ablauf

1. **Runde 1, Orchestration:** Booking ruft die Backends nacheinander an. Scheitert Hotel, ruft Booking Flight an und lässt stornieren. Booking antwortet dem Kunden erst, wenn alles erledigt ist.
2. **Runde 2, Choreography:** Booking spricht Flight an (Karte), Hotel an (lehnt ab), legt einen Zettel `CompensationRequested` auf den Tisch und antwortet dem Kunden sofort. Flight schaut irgendwann auf den Tisch, holt den Zettel und zieht die Karte ein.
3. Vergleichen: Wer hat gewartet? Wer wusste am Ende, ob storniert wurde?

### Was passieren sollte

- Runde 1: Booking wartet die ganze Zeit und weiß am Ende alles. Scheitert der Storno bei Flight, merkt Booking es sofort.
- Runde 2: Booking wartet nicht und weiß am Ende nicht, ob storniert wurde. Die Runde fühlt sich schneller und lockerer an, aber niemand hat mehr den Überblick. Das ist der Trade-off.

---

## Teil 2: Drei Störungen (ca. 20 Minuten)

Runde 2 noch einmal, jedes Mal mit einem Haken. Vor jeder Störung festlegen: Was passiert, wenn der Tisch ein **Blatt Papier** ist? Was, wenn er ein **Postfach mit Quittung** ist? Das ist der Unterschied zwischen einem Webhook und einem Broker.

| Störung | Beim Blatt Papier | Beim Postfach mit Quittung |
|---|---|---|
| **Flight ist nicht da.** Flight verlässt den Raum, bevor der Zettel liegt. Booking hat dem Kunden schon geantwortet. | Der Zettel liegt ewig, niemand merkt es, der Flug bleibt gebucht. | Der Zettel bleibt sicher liegen, Flight holt ihn beim Zurückkommen ab (Persistenz, Redelivery). |
| **Der Zettel kommt zweimal.** Booking ist unsicher und legt einen zweiten hin. Flight findet beide. | Flight storniert zweimal. Bei einer Rückerstattung wäre das Geld zweimal weg. | Jeder Zettel trägt eine Nummer (`eventId`), Flight merkt sich bearbeitete Nummern (Idempotenz). |
| **Niemand sagt Bescheid.** Flight hat storniert. Der Kunde fragt Booking: „Ist mein Flug jetzt weg?" | Booking kann nur sagen: „Ich habe einen Zettel hingelegt." | Flight legt einen Antwortzettel `BookingCancelled` hin, Booking wartet darauf mit Timeout (Reply-Pattern, `STUCK`-Status). |

Auflösung in Spielsprache: Das Postfach hält den Zettel fest und liefert ihn nach, auch doppelt (Persistenz, Redelivery). Deshalb trägt jeder Zettel eine Nummer, wer sie schon kennt, legt ihn weg (`eventId` als Dedup-Key, Idempotenz). Wer eine Antwort braucht, wartet auf einen Antwortzettel (Reply-Event, Timeout, `STUCK`).

Merksatz: Eventing macht das Problem nicht kleiner. Es macht es leiser.

---

## Recap: Vier Fragen

1. Ohne Broker: Der Zettel kommt nie an. Booking schickt das Event per HTTP-POST, Hotel ist down. Wo ist das Event jetzt?
2. Mit Broker: Der Zettel kommt zweimal. Der Broker liefert, bis das Backend bestätigt. Stirbt es dazwischen, kommt das Event noch einmal. Was, wenn das Backend echten State ändert?
3. Webhooks reichen doch. Kein Broker, keine Infrastruktur. Wozu Kafka, RabbitMQ oder NATS überhaupt?
4. Brauchen wir überhaupt einen Dirigenten? Volle Choreographie hieße: Hotel reagiert auf Flight, Car auf Hotel. Wer merkt es, wenn Hotel das nicht weiß? (Faustregel nach Chris Richardson, *Microservices Patterns*: Choreographie für einfache Sagas ohne Reihenfolge, Orchestrierung sobald es einen Ablauf gibt.)

Bonus zu Frage 3: Was ein Broker kann, was ein Webhook nicht kann. Behalten (Persistenz, Redelivery, Replay), Verteilen (Fan-out, Backpressure, Ordering), Scheitern lassen (Dead-Letter-Queue). Exactly-once gibt es trotzdem nicht, die praktische Näherung ist at-least-once plus Idempotenz.

Antworten und Anekdoten: [questions/story7.md](../questions/story7.md)

---

## Referenz-Implementierung

Es gibt weiterhin eine Referenz unter `services/booking/story7/`: Forward-Pfad wie Story 6, die Kompensation läuft als `CompensationRequested`-Event per Webhook-POST an die Backends, die sofort mit `202 Accepted` antworten und asynchron stornieren. Im Dashboard zeigt das Choreography-Panel unter Story 7 die Saga-Zustände. Demo nach dem Spiel: Hotel auf „Fehler" stellen, `POST /booking/bookings` auslösen, beobachten, wie Booking auf `FAILED` geht, bevor der Rollback im Backend abgeschlossen ist.

Wer das Pattern selbst bauen möchte, findet die Schmalspur-Variante und die Bonus-Kriterien (Reply-Events, Timeout-Erkennung, persistente Idempotenz, Feature-Flag) in der Referenz. Die technischen Hinweise unten gelten dafür unverändert.

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
- Was passiert, wenn ein Reply-Event nie ankommt? (Anschluss an Story 6, Frage 5 zu Saga-Beobachtbarkeit)
- Wo wandert das Saga-Wissen jetzt hin — und wann wird Choreography zum verteilten Monolithen?
- Welcher Teil der Story-6-Implementierung fällt komplett weg, welcher bleibt unverändert?

## Bonus (optional)

- **Dead-Letter-Behandlung:** Was passiert mit Events, die nach N Versuchen nicht zugestellt werden konnten?
- **Mischbetrieb über Feature-Flag:** Booking kann zur Laufzeit zwischen synchroner Kompensation (Story 6) und asynchroner Kompensation (Story 7) umschalten — schöner Showcase im Workshop
- **Auch der Happy Path über Events:** Nicht nur Kompensation, sondern auch das Forward-Booking als Event-Choreography — zeigt, wie weit man Choreography treiben kann
