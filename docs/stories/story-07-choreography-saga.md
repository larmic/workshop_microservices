# Story 7: Die Saga wird leise

**Thema:** Event-Driven Architecture, Choreography-Saga, als Würfelspiel statt Code
**Zeitrahmen:** ca. 60 Minuten: 10 Minuten Einstieg (Konzeptfolien), 20 Minuten Spiel, 5 Minuten Nachbereitung „Die Leiter“, 10 Minuten Demo mit der Referenz, 15 Minuten Recap. Spiel und Leiter zusammen sind die ≈ 25 Minuten auf der Folie.

## Kontext

In Story 6 trägt der Booking-Service die volle Verantwortung für die Kompensation: Er ruft synchron `DELETE /bookings/{id}` gegen jeden zuvor erfolgreich aufgerufenen Backend-Service auf, wartet auf jede Antwort und behandelt Fehler in seinem eigenen Code. Damit ist Booking nicht nur Orchestrator des Happy Paths, sondern auch Single Point of Responsibility für jede Stornierung. Ist ein Backend kurz nicht erreichbar, muss Booking den Fehler selbst behandeln und meldet `COMPENSATION_FAILED`.

Fachlich ist die Stornierung aber Aufgabe des jeweiligen Backends: Wer eine Buchung anlegen kann, muss sie auch zurücknehmen können. Deshalb ruft Booking den gebuchten Backends das Storno nur noch als **Event** zu (`CompensationRequested`) und wartet nur, bis das Backend es angenommen hat (`202 Accepted`). Stornieren muss das Backend selbst. Booking hängt damit nicht mehr an der Dauer des Stornos, an der Erreichbarkeit aber weiter: Antwortet ein Backend nicht, wartet Booking wie in Story 6 bis zum Timeout (3 Sekunden pro Zuruf). Der Preis: Fehler werden still.

Ehrlich gesagt ist das ein **Mittelweg** und keine volle Choreografie: Booking dirigiert weiter den Hinweg und entscheidet, wer ein Storno bekommt. Storniert wird per Zuruf. Der Unterschied zu Story 6 steckt nicht im Code, der ist klein. Er steckt darin, wer wartet und wer am Ende Bescheid weiß. Deshalb würfeln wir die Saga, statt sie zu programmieren.

## User Story

Als **System**
möchte ich **bei einer fehlgeschlagenen Buchung die Kompensation asynchron an die zuständigen Backend-Services delegieren**,
damit **der Booking-Service nicht für deren Verfügbarkeit haften muss und die Backends ihre eigene Stornierungslogik kapseln können**.

---

## Teil 1: Das Würfelspiel „Retry oder Storno“ (ca. 20 Minuten)

Booking und Kunde wollen 4 Reisen buchen und haben dafür 6 Versuche. Storniert wird per Zuruf.

### Material (pro Tisch)

- Ein Würfel
- 13 Pokerchips für den Kunden, Farbe egal
- Ein Zettel und ein Stift für Booking

Nichts zu drucken, nichts vorzubereiten. Gesamtliste: [`docs/orga/materialliste.md`](../orga/materialliste.md).

### Tische und Rollen

Tische mit vier bis sieben Personen, am besten fünf. Die Rollen bleiben das ganze Spiel gleich:

- **Booking und Kunde** sind ein Team. Booking würfelt, ruft und führt den Zettel. Der Kunde zahlt.
- **Flight, Hotel und Car** würfeln nur, sie entscheiden nichts. Was bei ihnen liegen bleibt, behalten sie.

Zu viert übernimmt Booking auch Car. Ab sechs Personen führt jemand die Strichliste und trägt am Ende an der Tafel ein. Ab 8 Personen spielen zwei Tische, jeder mit eigenem Würfel und eigenen Chips. Mehr Tische sind besser: Bei vier Tischen taucht fast sicher ein stilles Storno auf.

### Die Regel

1. Vor der ersten Reise gibt der Kunde beliebig viele seiner 13 Chips als **Retry-Budget** an Booking. Sie kommen nie zurück.
2. Booking würfelt pro Reise: Bei 4–6 buchen alle. Bei 1–3 sagt ein Service nein (1 Flight, 2 Hotel, 3 Car).
3. Gebucht wird in der Reihenfolge Flight, Hotel, Car. Jedes Ja kostet den Kunden einen Chip an diesen Service.
4. Wer nein sagt, würfelt: 1–3 „503“, 4–6 „ausgebucht“.
5. Bei 503 darf Booking statt des Kunden einen Budget-Chip an den Service zahlen. Der sagt ja, und die Reise läuft weiter.
6. Sonst ruft Booking rückwärts nur die schon Gebuchten („Storno Hotel! Storno Flight!“), sagt sofort „Storniert“ und macht einen Strich.
7. Erst dann würfelt jeder Gerufene: Bei 1–4 geht der Chip zurück an den Kunden, bei 5–6 hat er nichts gehört, behält den Chip und schweigt.
8. Hat der Kunde vor einer Reise weniger als 3 Chips, ist Schluss.
9. Das Team gewinnt, wenn es in höchstens 6 Versuchen 4 Reisen bucht, sonst gewinnen die Downstreams.

Ergänzungen zum Ansagen:

- Chips vorher laut abzählen lassen. Die Balance hängt an 13 Chips und 6 Reisen.
- „Statt des Kunden“ betonen: Bei einem Retry zahlt nur das Budget, nicht zusätzlich der Kunde.
- Booking ruft nur die schon Gebuchten. Sagt Flight nein, gibt es nichts zu stornieren.
- „Storniert“ kommt vor dem Würfeln der Gerufenen, nie danach. Wer warten will, hat das Thema verstanden: Stoff für die Nachbesprechung.
- Die 6 Reisen ausspielen, auch wenn das Ziel schon verloren ist. Sonst fehlen stille Stornos.

### Die eine Entscheidung: Wie viel Budget?

Jeder Budget-Chip fehlt dem Kunden beim Bezahlen, er hilft nur bei 503, also bei jedem zweiten Nein, und unbenutzt verfällt er. Im Spiel selbst gibt es keine Falle: Hat Booking Budget, setzt es den Chip bei 503 immer ein.

| Retry-Budget | Team gewinnt | Spiele mit stillem Storno |
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

Booking liest seine Striche vor, dann zählen Flight, Hotel und Car offen ihre Chips. Jeder Service sollte genau so viele Chips haben, wie Reisen durch sind (Budget-Chips eingerechnet). Jeder Chip darüber ist ein stilles Storno: Booking hat „storniert“ gemeldet, der Kunde hat sein Geld nicht zurück. Der Trainer sammelt pro Tisch Budget, gewonnen oder verloren und stille Stornos an der Tafel.

### Was passieren sollte

- Tische mit Budget 0 verlieren meist, weil jedes 503 zur Kompensation wird. Tische mit Budget 3 verlieren, weil Chips fürs Buchen fehlen.
- An den meisten Tischen bleibt mindestens ein Chip bei einem Service liegen, obwohl Booking „storniert“ gemeldet hat. Genau das tut die Referenz: `dispatch-failed` im Log, Step trotzdem `COMPENSATED` (Recap 01).
- Die doppelte Zustellung (Korb) wird nicht gespielt, sondern auf der Leiter als Gedankenexperiment aufgelöst.

Bewusste Vereinfachungen: Die Würfelquoten sind Zeitraffer, in echt geht eins von tausend Stornos verloren, nicht eins von dreien. Car hat im Spiel keinen Circuit Breaker. In der Referenz öffnet der Car-Breaker nach fünf Fehlschlägen, danach startet 30 Sekunden lang keine Saga. Die Referenz wiederholt keine Vorwärtsschritte, das Retry-Budget ist eine Erweiterung im Spiel.

---

## Teil 2: Nachbereitung „Die Leiter“ (ca. 5 Minuten)

Booking storniert nicht mehr selbst, antwortet aber weiter dem Kunden. Was fehlt, damit kein Chip liegen bleibt und Booking es weiß? Die Folie deckt Stufe für Stufe auf. Gespielt wurde nur der Zuruf, die übrigen Stufen sind Gedankenexperimente an den liegen gebliebenen Chips.

| Stufe | Im Spiel | Technik | In der Referenz |
|---|---|---|---|
| Anruf | nicht gespielt, kennen alle aus Story 6: Booking wartet, bis der Chip zurück ist | synchroner Aufruf, `DELETE /bookings/{id}` | ja, in Story 6 |
| Zuruf | gespielt: Booking ruft, meldet sofort „storniert“, bei 5–6 bleibt der Chip liegen | Webhook, `POST /events/compensation`, `202 Accepted` | ja, Story 7 |
| Korb | Gedankenexperiment: Der Korb stellt zu, bis „erledigt“ kommt. Geht das „erledigt“ verloren, zahlt der Service zweimal | Broker mit Persistenz und Redelivery (at-least-once) | nein |
| Nummer | Gedankenexperiment: Wer sich die Storno-Nummer merkt, quittiert nur und zahlt nicht noch einmal | Dedup über `eventId`, Idempotenz | nein, die Backends loggen nur |
| Antwort zurück | Gedankenexperiment: Booking meldet „storniert“ erst, wenn alle Chips zurück sind, sonst nach einem Timeout `STUCK` | Reply-Event, Timeout, Saga-Status `STUCK` | nein |

Merksatz: Eventing macht das Problem nicht kleiner. Es macht es leiser.

---

## Teil 3: Die Demo mit der Referenz (ca. 10 Minuten)

Die Demo kommt nach der Leiter (Hinweis in der Note von Folie 24a).

Im Dashboard unter Story 7: Car auf „Fehler“ stellen, dann `POST /booking/bookings` auslösen. Die Saga endet mit `FAILED`, Hotel und Flight stehen auf `COMPENSATED`. Im Log von `booking-ref-story7` sieht man die Kompensation rückwärts, erst Hotel, dann Flight, jeweils mit `phase=publishing` und `phase=dispatched`.

Nach fünf Klicks ist der Car-Breaker offen. Danach startet 30 Sekunden lang keine Saga, der Pre-Check antwortet sofort mit 503. Das stille Storno aus dem Spiel zeigt die Demo nicht: Ist Flight down, scheitert schon die Buchung, und es gibt nichts zu stornieren. Dass Booking auch ohne 202 `COMPENSATED` setzt, zeigt der Code (`compensate` in `services/booking/story7/handler/booking.go`).

---

## Recap: Vier Fragen (ca. 15 Minuten)

1. Ein Chip bleibt liegen. Booking schickt das Storno per HTTP-POST, Flight ist seit der Buchung down. Wo ist das Event jetzt?
2. Der Korb liefert doppelt. Der Broker liefert, bis das Backend bestätigt. Stirbt es dazwischen, kommt das Event noch einmal. Was, wenn das Backend echten State ändert?
3. Webhooks reichen doch. Kein Broker, keine Infrastruktur. Wozu Kafka, RabbitMQ oder NATS überhaupt?
4. Brauchen wir überhaupt einen Dirigenten? Volle Choreographie hieße: Hotel reagiert auf Flight, Car auf Hotel. Wer merkt es, wenn Hotel das nicht weiß? (Faustregel nach Chris Richardson, *Microservices Patterns*: Choreographie für einfache Sagas ohne Reihenfolge, Orchestrierung sobald es einen Ablauf gibt.)

Bonus zu Frage 3: Was ein Broker kann, was ein Webhook nicht kann. Behalten (Persistenz, Redelivery, Replay), Verteilen (Fan-out, Backpressure, Ordering), Scheitern lassen (Dead-Letter-Queue). Exactly-once-Zustellung gibt es trotzdem nicht, at-least-once plus Idempotenz ergibt exactly-once-Wirkung.

Antworten und Anekdoten: [questions/story7.md](../questions/story7.md)

---

## Referenz-Implementierung

Unter `services/booking/story7/` liegt die Referenz: Der Hinweg ist derselbe wie in Story 6 (Flight, Hotel, Car nacheinander, je mit Circuit Breaker). Scheitert ein Schritt, schickt Booking jedem Schritt mit Status `BOOKED` ein `CompensationRequested`-Event, rückwärts und einzeln per `POST /events/compensation`. Booking wartet jeweils auf `202 Accepted`, höchstens 3 Sekunden (HTTP-Client-Timeout). Kommt keine 202, loggt Booking `dispatch-failed` und setzt den Step trotzdem auf `COMPENSATED`. Es gibt genau einen Versuch, kein Retry und keine Nachlieferung. Die Backends antworten sofort mit 202 und loggen die Stornierung danach in einer Goroutine, sie deduplizieren nicht.

Im Dashboard zeigt das Choreography-Panel unter Story 7 die Saga-Zustände. Reply-Events, Timeout-Erkennung (`STUCK`) und persistente Idempotenz gibt es in der Referenz **nicht**. Sie sind die Bonus-Kriterien für alle, die das Pattern selbst weiterbauen wollen.

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
- **Wartet nur auf die Annahme (Pflicht-Variante):** Booking wartet nacheinander auf jede 202, aber nicht auf die eigentliche Stornierung. Ob sie geklappt hat, weiß Booking nicht. Das ist die bewusst schmale Variante, die zeigt, was Eventing ohne Broker strukturell nicht leistet (Diskussion im Recap).
- **Reply-Pattern (Bonus):** Wer die volle Variante will, baut den Rückkanal dazu: Das Backend POSTet später ein `BookingCancelled` oder `CancellationFailed` an einen Booking-Endpoint, Booking führt den Saga-Status nach. Dazu Timeout-Erkennung: Bleibt ein Reply aus, geht die Saga auf `STUCK`. Das ist asynchrones Request-Reply über zwei Webhook-Richtungen.
- **Was Booking aufgibt, was bleibt (Pflicht-Variante):**
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
- **Auch der Happy Path über Events:** Nicht nur Kompensation, sondern auch das Forward-Booking als Event-Choreography. Das zeigt, wie weit man Choreography treiben kann.
