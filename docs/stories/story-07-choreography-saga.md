# Story 7: Die Saga wird leise

**Thema:** Event-Driven Architecture, Choreography-Saga, als Rollenspiel statt Code
**Zeitrahmen:** ca. 60 Minuten: 10 Minuten Einstieg (Konzeptfolien), 20 Minuten Spiel, 5 Minuten Nachbereitung „Die Leiter“, 10 Minuten Demo mit der Referenz, 15 Minuten Recap. Spiel und Leiter zusammen sind die ≈ 25 Minuten auf der Folie.

## Kontext

In Story 6 trägt der Booking-Service die volle Verantwortung für die Kompensation: Er ruft synchron `DELETE /bookings/{id}` gegen jeden zuvor erfolgreich aufgerufenen Backend-Service auf, wartet auf jede Antwort und behandelt Fehler in seinem eigenen Code. Damit ist Booking nicht nur Orchestrator des Happy Paths, sondern auch Single Point of Responsibility für jede Stornierung. Ist ein Backend kurz nicht erreichbar, muss Booking den Fehler selbst behandeln und meldet `COMPENSATION_FAILED`.

Fachlich ist die Stornierung aber Aufgabe des jeweiligen Backends: Wer eine Buchung anlegen kann, muss sie auch zurücknehmen können. Deshalb ruft Booking den gebuchten Backends das Storno nur noch als **Event** zu (`CompensationRequested`) und wartet nur, bis das Backend es angenommen hat (`202 Accepted`). Stornieren muss das Backend selbst. Booking hängt damit nicht mehr an der Dauer des Stornos, an der Erreichbarkeit aber weiter: Antwortet ein Backend nicht, wartet Booking wie in Story 6 bis zum Timeout (3 Sekunden pro Zuruf). Der Preis: Fehler werden still.

Ehrlich gesagt ist das ein **Mittelweg** und keine volle Choreografie: Booking dirigiert weiter den Hinweg und entscheidet, wer ein Storno bekommt. Storniert wird per Zuruf. Der Unterschied zu Story 6 steckt nicht im Code, der ist klein. Er steckt darin, wer wartet und wer am Ende Bescheid weiß. Deshalb spielen wir die Saga, statt sie zu programmieren.

## User Story

Als **System**
möchte ich **bei einer fehlgeschlagenen Buchung die Kompensation asynchron an die zuständigen Backend-Services delegieren**,
damit **der Booking-Service nicht für deren Verfügbarkeit haften muss und die Backends ihre eigene Stornierungslogik kapseln können**.

---

## Teil 1: Das Rollenspiel (ca. 20 Minuten)

Ihr seid die Kund:innen. Car sagt immer nein. Der Becher zeigt, was wirklich storniert wurde.

### Material (pro Tisch)

- Sechs verdeckte Listen (je drei für Hotel und Flight, eine pro Runde) oder je eine Liste mit drei Spalten, auf Klemmbrett oder gefaltetem Blatt
- Ein Becher pro Reise, mit Reisenummer beschriftet (1, 2, 3, …)
- Gummibärchen in zwei Schälchen, eins für Hotel, eins für Flight: etwa sieben pro Reise und Tisch, plus Naschreserve (Alternative: Büroklammern, etwa bei Allergien oder veganer Ernährung)
- Ein Korb oder eine Schachtel für Runde 3
- Pro Reise zwei vorbereitete Zettel, fortlaufend nummeriert in Storno-Reihenfolge: E1 „Hotel · Reise 1“, E2 „Flight · Reise 1“, E3 „Hotel · Reise 2“, E4 „Flight · Reise 2“, …
- Eine Strichliste für Booking
- Ein Briefing-Kärtchen für Flight (siehe unten)

Spielbar ab sechs Personen: Booking, Hotel, Flight, Car und mindestens zwei Kund:innen. Ab 9 Personen an einem Tisch spielt eine eigene Person den Korb, darunter übernimmt Car. Ab 12 Personen spielen zwei Tische mit eigenen Services und identischer Ansage. Jeder Tisch hat dann mindestens zwei Reisen, das reicht, weil das Briefing auf Reise 2 setzt.

### Rollen

1. **Booking** hat eine Strichliste und bucht eine Reise nach der anderen bei Flight, Hotel und Car. Sagt Car nein, storniert Booking rückwärts, erst Hotel, dann Flight. Danach sagt Booking „storniert“.
2. **Hotel und Flight** buchen auf einer verdeckten Liste. Stornieren heißt: an den Fingern bis fünf zählen, durchstreichen, ein Gummibärchen in den Becher der Reise. Jedes Mal, auch wenn die Reise schon durchgestrichen ist. Flight folgt den geheimen Ansagen des Trainers.
3. **Car** hat heute keine Autos und sagt immer „nein“. In Runde 3 ist Car auch der Korb.
4. **Alle anderen** haben je eine Reise. Die Becher stehen nummeriert auf dem Tisch, angestellt wird nach Reisenummer. Die nächste Reise kommt dran, wenn die vorige „storniert“ gehört hat. Storniert niemand mehr und ist der Korb leer: **Kassensturz**, in jeden Becher gehören zwei.

Ergänzungen zum Ansagen:

- Die Kund:innen stellen sich in Reihenfolge der Reisenummern bei Booking an, sonst passen die vornummerierten Zettel aus Runde 3 nicht. Die Becher stehen nummeriert in einer Reihe auf dem Tisch und bleiben die ganze Runde dort, auch wenn die Kund:in schon „storniert“ gehört hat. Booking bedient immer nur eine Reise.
- Weil Car immer ablehnt, scheitert jede Reise. Hotel und Flight sind dann schon gebucht und müssen storniert werden. Deshalb gehören nach jeder Runde genau zwei Gummibärchen in jeden Becher.
- Das Zählen bis fünf gibt dem Storno eine spürbare Dauer, in allen Runden gleich. In Runde 1 wartet die Schlange diese Zeit hörbar mit, in Runde 2 und 3 nicht.
- „Jedes Mal“ ist wörtlich gemeint: Hotel und Flight legen bei jeder Storno-Ansage ein Gummibärchen hin, auch wenn die Reise auf ihrer Liste schon durchgestrichen ist. Menschen deduplizieren von selbst, unsere Backends nicht.
- Ist Car auch der Korb, sagt Car zwischendurch weiter „nein“ zu neuen Buchungen. Eine Buchungsfrage von Booking geht vor, danach liest der Korb weiter.
- Der Korb beginnt mit dem Zählen erst, wenn das Gummibärchen liegt, damit das Zählen des Backends nicht mit dem Timeout kollidiert. Reagiert ein Backend gar nicht (es zählt nicht an den Fingern), zählt der Korb direkt nach dem Vorlesen.
- Die Runde endet erst, wenn Hotel und Flight fertig sind und der Korb leer ist. Wer früher zählt, übersieht Gummibärchen, die noch unterwegs sind. Genau das ist der Lernpunkt aus Runde 3: „storniert“ kam vor dem Gummibärchen, das ist Eventual Consistency.
- Kassensturz: Alle Kund:innen zählen laut ihren Becher. Dann deckt Flight die Liste der Runde auf (in Runde 2 ist Reise 2 nicht durchgestrichen). Danach wird geleert (oder gegessen), und die nächste Runde beginnt mit leeren Bechern und der nächsten Liste.

### Ablauf

Runden einzeln ansagen, erst tippen lassen, dann spielen, dann Kassensturz, dann das Fragment aufdecken. Bei zwei Tischen jedes Fragment erst aufdecken, wenn beide Tische gezählt haben.

1. **Runde 1, Anruf (wie Story 6):** „Hotel, storniere Reise 1!“ Booking wartet auf „erledigt“, dann dasselbe mit Flight. Das Backend sagt „erledigt“ erst, wenn das Gummibärchen im Becher liegt. Frage: Wer wartet?
2. **Runde 2, Zuruf (unsere Referenz):** Das Backend ruft sofort „gehört!“ und storniert danach, Booking ruft gleich das nächste Backend (erst Hotel, dann Flight) und bedient danach die nächste Reise. Bleibt es still, zählt Booking laut bis drei, macht einen Strich und meldet trotzdem „storniert“. Frage: Stimmt jeder Becher?
3. **Runde 3, Korb:** Booking wirft pro Storno einen Zettel in den Korb (erst Hotel, dann Flight) und meldet sofort „storniert“. Der Korb liest die Zettel der Reihe nach vor („E4: Flight, storniere Reise 2!“), das Backend storniert und sagt „quittiert“. Fehlt das nach dem Gummibärchen, zählt der Korb laut bis fünf und liest denselben Zettel noch einmal. Erst nach „quittiert“ kommt der nächste Zettel dran. Frage: Stimmt jeder Becher?

Zur Vorhersage: Dass etwas schiefgehen kann, verrät die Rollenzeile („Flight folgt den geheimen Ansagen“). Wer nachfragt, bekommt als Hilfe: „Angenommen, Flight hört ein Storno nicht.“ Die Nummer auf dem Zettel regelt im Spiel nur die Reihenfolge, Flight merkt sie sich nicht. Was eine gemerkte Nummer bringen würde, zeigt erst die Leiter.

### Geheimes Briefing

Nur Flight bekommt eine Ansage, als Kärtchen, bei zwei Tischen für beide Flights identisch. Hotel spielt alle Runden normal.

| Runde | Ansage an Flight |
|---|---|
| 1 | Spiel ganz normal. |
| 2 | Bei Reise 2 reagierst du auf „storniere“ nicht, kein „gehört“, kein Gummibärchen. Buchen geht normal. |
| 3 | Bei Reise 2 legst du das Gummibärchen hin, sagst aber nicht „quittiert“. Liest der Korb noch einmal vor, machst du weiter wie immer. |

### Was passieren sollte

- **Runde 1:** Jeder Becher hat zwei. Die ganze Schlange wartet, bis beide Backends ihr Gummibärchen hingelegt haben, aber Booking weiß am Ende alles. Wäre ein Storno gescheitert, hätte Booking es sofort gehört.
- **Runde 2:** Die Schlange läuft spürbar schneller, denn das Zählen bis fünf passiert jetzt, während Booking schon weitermacht. Beim Kassensturz hat Reise 2 nur ein Gummibärchen, der Flug ist noch gebucht. Booking hat es sogar gemerkt (ein Strich), aber trotzdem „storniert“ gemeldet. Genau das tut die Referenz: `dispatch-failed` im Log, Step trotzdem `COMPENSATED` (Recap 01).
- **Runde 3:** Keiner fehlt, einer hat drei. Der Korb hat den Zettel E4 wiederholt, weil das „quittiert“ fehlte, und Flight hat zweimal storniert (Recap 02). Bei einem Storno ist das harmlos, bei einer Rückerstattung wäre das Geld zweimal weg.
- **Runde 2 gegen Runde 3:** Ohne Korb geht ein Storno verloren, mit Korb kommt es doppelt. Was ein Broker mehr kann als ein Webhook, ist Recap 03.

Bewusste Vereinfachungen: Im Spiel hat Car keinen Circuit Breaker. In der Referenz streicht niemand etwas durch: Die Backends haben keinen Buchungsbestand und loggen nur `received`, `processing` und `done`. In der Referenz öffnet er nach fünf Fehlschlägen, danach startet 30 Sekunden lang keine Saga. Die Schlange ist ebenfalls vereinfacht, denn die Referenz bedient Anfragen parallel, nur die Schritte innerhalb einer Saga laufen nacheinander. Der Korb (Runde 3) ist nicht die Referenz, sondern die erste Bonus-Stufe.

---

## Teil 2: Nachbereitung „Die Leiter“ (ca. 5 Minuten)

Booking haftet nicht mehr für das Stornieren selbst, aber weiter für Saga-Status und Antwort an die Kund:in. Was fehlt, damit jeder Becher stimmt und Booking es weiß? Die Folie deckt Stufe für Stufe auf. Gespielt werden nur Anruf, Zuruf und Korb, Nummer und Antwortzettel sind reine Auflösung.

| Stufe | Im Spiel | Technik | In der Referenz |
|---|---|---|---|
| Anruf | Runde 1, Booking wartet auf „erledigt“ | synchroner Aufruf, `DELETE /bookings/{id}` | ja, in Story 6 |
| Zuruf | Runde 2, Booking wartet nur auf „gehört“ | Webhook, `POST /events/compensation`, `202 Accepted` | ja, Story 7 |
| Korb | Runde 3, liest vor bis „quittiert“ | Broker mit Persistenz und Redelivery (at-least-once) | nein |
| Nummer | nicht gespielt: Wer sich die Zettelnummer merkt und sie wiedererkennt, quittiert nur, ohne Gummibärchen | Dedup über `eventId`, Idempotenz | nein, die Backends loggen nur |
| Antwortzettel | nicht gespielt: Booking meldet „storniert“ erst nach zwei Antwortzetteln, sonst zählt Booking bis zehn und meldet `STUCK` | Reply-Event, Timeout, Saga-Status `STUCK` | nein |

Merksatz: Eventing macht das Problem nicht kleiner. Es macht es leiser.

---

## Teil 3: Die Demo mit der Referenz (ca. 10 Minuten)

Die Demo kommt nach der Leiter (Hinweis in der Note von Folie 24a).

Im Dashboard unter Story 7: Car auf „Fehler“ stellen, dann `POST /booking/bookings` auslösen. Die Saga endet mit `FAILED`, Hotel und Flight stehen auf `COMPENSATED`. Im Log von `booking-ref-story7` sieht man die Kompensation rückwärts, erst Hotel, dann Flight, jeweils mit `phase=publishing` und `phase=dispatched`.

Nach fünf Klicks ist der Car-Breaker offen. Danach startet 30 Sekunden lang keine Saga, der Pre-Check antwortet sofort mit 503. Den Strich aus Runde 2 zeigt die Demo nicht: Ist Flight down, scheitert schon die Buchung, und es gibt nichts zu stornieren. Dass Booking auch ohne 202 `COMPENSATED` setzt, zeigt der Code (`compensate` in `services/booking/story7/handler/booking.go`).

---

## Recap: Vier Fragen (ca. 15 Minuten)

1. Runde 2: Ein Gummibärchen fehlt. Booking schickt das Storno per HTTP-POST, Flight ist seit der Buchung down. Wo ist das Event jetzt?
2. Runde 3: Ein Becher hat drei. Der Broker liefert, bis das Backend bestätigt. Stirbt es dazwischen, kommt das Event noch einmal. Was, wenn das Backend echten State ändert?
3. Runde 2 gegen Runde 3: Webhooks reichen doch. Kein Broker, keine Infrastruktur. Wozu Kafka, RabbitMQ oder NATS überhaupt?
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
