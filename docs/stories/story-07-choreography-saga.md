# Story 7: Die Saga wird leise

**Thema:** Event-Driven Architecture, Choreography-Saga, als Würfelspiel statt Code
**Zeitrahmen:** ca. 60 Minuten: 10 Minuten Einstieg (Konzeptfolien), 25 Minuten Spiel inklusive Kassensturz, 5 Minuten Demo mit der Referenz, 20 Minuten Recap. Die ≈ 25 Minuten auf der Folie sind das Spiel mit Kassensturz.

## Kontext

In Story 6 trägt der Booking-Service die volle Verantwortung für die Kompensation: Er ruft synchron `DELETE /bookings/{id}` gegen jeden zuvor erfolgreich aufgerufenen Backend-Service auf, wartet auf jede Antwort und behandelt Fehler in seinem eigenen Code. Damit ist Booking nicht nur Orchestrator des Happy Paths, sondern auch Single Point of Responsibility für jede Stornierung. Ist ein Backend kurz nicht erreichbar, muss Booking den Fehler selbst behandeln und meldet `COMPENSATION_FAILED`.

Fachlich ist die Stornierung aber Aufgabe des jeweiligen Backends: Wer eine Buchung anlegen kann, muss sie auch zurücknehmen können. Deshalb ruft Booking den gebuchten Backends das Storno nur noch als **Event** zu (`CompensationRequested`) und wartet nur, bis das Backend es angenommen hat (`202 Accepted`). Stornieren muss das Backend selbst. Booking hängt damit nicht mehr an der Dauer des Stornos, an der Erreichbarkeit aber weiter: Antwortet ein Backend nicht, wartet Booking wie in Story 6 bis zum Timeout (3 Sekunden pro Zuruf). Der Preis: Fehler werden still.

Ehrlich gesagt ist das ein **Mittelweg** und keine volle Choreografie: Booking dirigiert weiter den Hinweg und entscheidet, wer ein Storno bekommt. Storniert wird per Zuruf. Der Unterschied zu Story 6 steckt nicht im Code, der ist klein. Er steckt darin, wer wartet und wer am Ende Bescheid weiß. Deshalb würfeln wir die Saga, statt sie zu programmieren.

## User Story

Als **System**
möchte ich **bei einer fehlgeschlagenen Buchung die Kompensation asynchron an die zuständigen Backend-Services delegieren**,
damit **der Booking-Service nicht für deren Verfügbarkeit haften muss und die Backends ihre eigene Stornierungslogik kapseln können**.

---

## Teil 1: Das Würfelspiel „Retry oder Storno“ (ca. 25 Minuten inklusive Kassensturz)

Booking und Kunde wollen 4 Reisen buchen und haben dafür 6 Versuche. Storniert wird per Zuruf.

### Material (pro Tisch)

- Ein Würfel
- 13 Pokerchips für den Kunden, Farbe egal
- Ein Zettel und ein Stift für Booking

Nichts zu drucken, nichts vorzubereiten. Gesamtliste: [`docs/orga/materialliste.md`](../orga/materialliste.md).

### Tische und Rollen

Tische mit vier bis sieben Personen, am besten fünf. Die Rollen bleiben das ganze Spiel gleich:

- **Kunde und Booking** spielen zusammen. Booking würfelt, ruft, bezahlt die Services und führt den Zettel. Der Kunde gibt Booking das Reisegeld.
- **Flight, Hotel und Car** würfeln nur, sie entscheiden nichts. Was bei ihnen liegen bleibt, behalten sie.

Zu viert übernimmt Booking auch Car. Ab sechs Personen führt jemand die Strichliste und trägt am Ende an der Tafel ein. Ab 8 Personen spielen zwei Tische, jeder mit eigenem Würfel und eigenen Chips. Mehr Tische sind besser: Bei vier Tischen taucht fast sicher ein stilles Storno auf.

### Die Regel

Ziel: 4 Reisen buchen. Höchstens 6 Reisen dürfen starten, gescheiterte zählen mit.

Vor dem Spiel:

1. **Kunde und Booking** spielen zusammen. **Flight, Hotel, Car** würfeln nur und behalten, was bei ihnen liegt.
2. Der Kunde hat 13 Chips. **Eine Reise kostet 3**, je einen für Flight, Hotel und Car.
3. Kunde und Booking machen aus, wie viele Chips der Kunde Booking als **Retry-Budget** gibt. Booking legt sie als eigenen Stapel zur Seite. Sie kommen nie zurück.
4. Kann der Kunde Booking keine 3 Chips mehr geben, ist das Spiel vorbei.

Jede Reise:

1. Der Kunde gibt Booking 3 Chips. Damit startet die Reise.
2. Booking fragt der Reihe nach **Flight, Hotel, Car** und würfelt einmal: Bei 4–6 buchen alle drei. Bei 1 sagt Flight nein, bei 2 Hotel, bei 3 Car. Nur wer vor dem Nein gebucht hat, bekommt einen Chip.
3. Wer nein sagt, würfelt den Grund: Bei 1–3 ist er nicht erreichbar (`503`). Booking zahlt ihm einen Chip aus dem Retry-Budget, er bucht doch. Wer nach ihm kommt, bucht ohne neuen Wurf, die Reise klappt.
4. Bei 4–6 ist er ausgebucht, oder das Budget ist leer: Die Reise scheitert. Booking ruft den schon Gebuchten rückwärts „Storno!“ zu und sagt sofort „Storniert“.
5. Erst dann würfelt jeder Gerufene: Bei 1–4 hat er es gehört und gibt den Chip an den Kunden zurück. Bei 5–6 hat er es überhört und behält ihn. Nach jeder Reise gibt Booking dem Kunden zurück, was vom Reisegeld übrig ist.

Gewonnen ist, sobald 4 Reisen geklappt haben. Verloren ist nach 6 Reisen mit weniger Erfolgen oder wenn der Kunde Booking keine 3 Chips mehr geben kann.

Ergänzungen zum Ansagen:

- Chips vorher laut abzählen lassen. Die Balance hängt an 13 Chips und 6 Reisen.
- Reisegeld und Retry-Budget liegen bei Booking getrennt. Beim Retry zahlt Booking aus dem Budget, der Chip aus dem Reisegeld für diesen Service geht nach der Reise an den Kunden zurück.
- Booking ruft nur die schon Gebuchten. Sagt Flight nein, gibt es nichts zu stornieren.
- „Storniert“ kommt vor dem Würfeln der Gerufenen, nie danach. Wer warten will, hat das Thema verstanden: Stoff für die Nachbesprechung.
- Die 6 Reisen ausspielen, auch wenn das Ziel schon verloren ist. Sonst fehlen stille Stornos.

### Die eine Entscheidung: Wie viel Budget?

Jeder Budget-Chip fehlt dem Kunden beim Bezahlen, er hilft nur bei 503, also bei jedem zweiten Nein, und unbenutzt verfällt er. Im Spiel selbst gibt es keine Falle: Hat Booking Budget, setzt es den Chip bei 503 immer ein.

| Retry-Budget | Kunde und Booking gewinnen | Spiele mit stillem Storno |
|---|---|---|
| 0 | 32 % | 59 % |
| 1 | 54 % | 43 % |
| 2 | 44 % | 32 % |
| 3 | 17 % | 27 % |

Simuliert mit 200.000 Spielen und exakt nachgerechnet. Zu wenig Resilienz verliert, zu viel auch. Mit 14 Chips oder 7 Reisen wird es deutlich leichter, mit 12 Chips deutlich schwerer.

### Vorwärts oder rückwärts

Scheitert ein Schritt einer Saga, hat der Orchestrator zwei Wege:

- **Rückwärts (Backward Recovery):** das schon Gebuchte kompensieren, in umgekehrter Reihenfolge. Das ist Story 6.
- **Vorwärts (Forward Recovery):** den gescheiterten Schritt wiederholen oder einen anderen Anbieter nehmen. Die Saga läuft weiter, nichts wird storniert.

Vorwärts lohnt nur bei Fehlern, die vorübergehen (`503`, Timeout), nicht bei einem fachlichen Nein (ausgebucht). Bei einem Timeout ist zudem offen, ob der Schritt schon gebucht hat. Ein Retry braucht dann Idempotenz, sonst entsteht eine Doppelbuchung. Unsere Referenz kennt nur rückwärts: Jeder Schritt läuft einmal durch den Circuit Breaker, scheitert er, wird kompensiert.

### Kassensturz

Booking liest seine Striche vor, dann zählen Flight, Hotel und Car offen ihre Chips. Jeder Service sollte genau so viele Chips haben, wie Reisen geklappt haben (Budget-Chips eingerechnet). Jeder Chip darüber ist ein stilles Storno: Booking hat „storniert“ gemeldet, der Kunde hat sein Geld nicht zurück. Der Trainer sammelt pro Tisch Budget, gewonnen oder verloren und stille Stornos an der Tafel.

### Was passieren sollte

- Tische mit Budget 0 verlieren meist, weil jedes 503 zur Kompensation wird. Tische mit Budget 3 verlieren, weil Chips fürs Buchen fehlen.
- An den meisten Tischen bleibt mindestens ein Chip bei einem Service liegen, obwohl Booking „storniert“ gemeldet hat. Genau das tut die Referenz: `dispatch-failed` im Log, Step trotzdem `COMPENSATED` (Recap 01).
- Die doppelte Zustellung (Broker) wird nicht gespielt, sondern im Recap (Frage 2) als Gedankenexperiment aufgelöst.

Bewusste Vereinfachungen: Die Würfelquoten sind Zeitraffer, in echt geht eins von tausend Stornos verloren, nicht eins von dreien. Car hat im Spiel keinen Circuit Breaker. In der Referenz öffnet der Car-Breaker nach fünf Fehlschlägen, danach startet 30 Sekunden lang keine Saga. Die Referenz wiederholt keine Vorwärtsschritte, das Retry-Budget ist eine Erweiterung im Spiel.

---

## Teil 2: Die Demo mit der Referenz (ca. 5 Minuten)

Die Demo kommt direkt nach dem Kassensturz (Note der Spielfolie). Die Chips bleiben bis zum Recap liegen.

Eine Saga genügt, das Bild im Dashboard ist dasselbe wie in Story 6. Im Dashboard unter Story 7: Car auf „Fehler“ stellen, dann `POST /booking/bookings` auslösen. Die Saga endet mit `FAILED`, Hotel und Flight stehen auf `COMPENSATED`. Neu sind zwei Dinge:

- Der Hinweis am Ergebnis: „Kompensation gemeldet, Ankunft ungeprüft“.
- Im Log von `booking-ref-story7` läuft die Kompensation rückwärts, erst Hotel, dann Flight, jeweils mit `phase=publishing` und `phase=dispatched`. Mehr als „abgeschickt“ weiß Booking nicht.

Das stille Storno aus dem Spiel zeigt die Demo nicht: Ist Flight down, scheitert schon die Buchung, und es gibt nichts zu stornieren. Dass Booking auch ohne 202 `COMPENSATED` setzt, zeigt der Code (`compensate` in `services/booking/story7/handler/booking.go`).

---

## Recap: Vier Fragen (ca. 20 Minuten)

1. Ein Chip bleibt liegen. Booking schickt das Storno per HTTP-POST, Flight ist seit der Buchung down. Wo ist das Event jetzt?
2. Das Storno kommt zweimal. Ein Broker hält das Storno fest und liefert, bis Flight bestätigt. Stirbt Flight dazwischen, kommt das Event noch einmal. Was, wenn das Backend echten State ändert?
3. Booking meldet zu früh „storniert“. Booking soll erst „storniert“ sagen, wenn alle Chips zurück sind. Was muss zurückkommen, und was gilt, wenn es nie kommt?
4. Brauchen wir überhaupt einen Dirigenten? Volle Choreografie hieße: Hotel reagiert auf Flight, Car auf Hotel. Wer merkt es, wenn Hotel das nicht weiß? (Faustregel nach Chris Richardson, *Microservices Patterns*: Choreografie für einfache Sagas ohne Reihenfolge, Orchestrierung sobald es einen Ablauf gibt.)

Die Fragen 1 bis 3 steigen auf, vom Zuruf über den Broker bis zur Antwort zurück:

| Stufe | Technik | In der Referenz |
|---|---|---|
| Zuruf (gespielt) | Webhook, `POST /events/compensation`, `202 Accepted` | ja |
| Broker (Frage 2) | Persistenz und Redelivery (at-least-once), Dedup über `eventId` | nein, die Backends loggen nur |
| Antwort zurück (Frage 3) | Reply-Event, Timeout, Saga-Status `STUCK` | nein |

Merksatz: Eventing macht das Problem nicht kleiner. Es macht es leiser.

Bonus zu Frage 2: Was ein Broker kann, was ein Webhook nicht kann. Behalten (Persistenz, Redelivery, Replay), Verteilen (Fan-out, Backpressure, Ordering), Scheitern lassen (Dead-Letter-Queue). Exactly-once-Zustellung gibt es trotzdem nicht, at-least-once plus Idempotenz ergibt exactly-once-Wirkung.

Antworten und Anekdoten: [questions/story7.md](../questions/story7.md)

---

## Referenz-Implementierung

Unter `services/booking/story7/` liegt die Referenz: Der Hinweg ist derselbe wie in Story 6 (Flight, Hotel, Car nacheinander, je mit Circuit Breaker). Scheitert ein Schritt, schickt Booking jedem Schritt mit Status `BOOKED` ein `CompensationRequested`-Event, rückwärts und einzeln per `POST /events/compensation`. Booking wartet jeweils auf `202 Accepted`, höchstens 3 Sekunden (HTTP-Client-Timeout). Kommt keine 202, loggt Booking `dispatch-failed` und setzt den Step trotzdem auf `COMPENSATED`. Es gibt genau einen Versuch, kein Retry und keine Nachlieferung. Die Backends antworten sofort mit 202 und loggen die Stornierung danach in einer Goroutine, sie deduplizieren nicht.

Im Dashboard zeigt das Choreography-Panel unter Story 7 die Saga-Zustände. Reply-Events, Timeout-Erkennung (`STUCK`) und persistente Idempotenz gibt es in der Referenz **nicht**. Sie sind Erweiterungen für alle, die das Pattern selbst weiterbauen wollen.

## Technische Hinweise

- **Webhook statt Broker:** Booking löst die Backends wie in Story 3 über Consul auf und schickt das Event direkt per `POST {backend}/events/compensation`. Es gibt keinen separaten Message-Broker und keine Subscription, Booking entscheidet selbst, wer ein Event bekommt.
- **Event-Struktur** (flacher Body):
  ```json
  {
    "eventId": "ev-3f9a1c2b",
    "sagaId": "B-4e2d7a10",
    "bookingId": "F-7c1a9f"
  }
  ```
  `eventId` ist `ev-` plus acht Hex-Zeichen und wäre der Dedup-Key. Das Backend antwortet mit `202 Accepted`, bei fehlenden Feldern mit `400`.
- **Wartet nur auf die Annahme (Referenz-Variante):** Booking wartet nacheinander auf jede 202, aber nicht auf die eigentliche Stornierung. Ob sie geklappt hat, weiß Booking nicht. Das ist die bewusst schmale Variante, die zeigt, was Eventing ohne Broker strukturell nicht leistet (Diskussion im Recap).
- **Reply-Pattern (Bonus):** Wer die volle Variante will, baut den Rückkanal dazu: Das Backend POSTet später ein `BookingCancelled` oder `CancellationFailed` an einen Booking-Endpoint, Booking führt den Saga-Status nach. Dazu Timeout-Erkennung: Bleibt ein Reply aus, geht die Saga auf `STUCK`. Das ist asynchrones Request-Reply über zwei Webhook-Richtungen.
- **Was Booking aufgibt, was bleibt (Referenz-Variante):**
  - **Aufgegeben:** die Ausführung der Stornierung. Booking weiß nicht mehr, ob sie geklappt hat
  - **Bleibt:** der **Saga-Status** gegenüber dem Kunden (Booking antwortet mit `FAILED`, sobald die Events angenommen oder gescheitert sind) und das Logging des Event-Dispatches
- **Vergleich Orchestration ↔ Choreography:**
  - Orchestration: Wissen zentral, Bug-Lokalisierung einfach, Kopplung höher
  - Choreography: Wissen verteilt, Backends entkoppelt, „verteilter Monolith“-Risiko bei schlechtem Schnitt

## Diskussions-Anker

- Warum nutzen wir Webhooks statt Kafka/RabbitMQ? Was würde sich ändern?
- Was passiert, wenn ein Reply-Event nie ankommt? (Anschluss an Story 6, Frage 5 zu Saga-Beobachtbarkeit)
- Wo wandert das Saga-Wissen jetzt hin, und wann wird Choreography zum verteilten Monolithen?
- Welcher Teil der Story-6-Implementierung fällt komplett weg, welcher bleibt unverändert?

## Bonus (optional)

- **Dead-Letter-Behandlung:** Was passiert mit Events, die nach N Versuchen nicht zugestellt werden konnten?
- **Idempotenter Empfänger:** Das Backend merkt sich verarbeitete `eventId`s in derselben Transaktion wie das Storno und quittiert Wiederholungen nur.
- **Auch der Happy Path über Events:** Nicht nur Kompensation, sondern auch das Buchen als Event-Choreography. Das zeigt, wie weit man Choreography treiben kann.
