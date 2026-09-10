# Story 5: Der geteilte Pool

**Thema:** Bulkhead Pattern, als Übung mit Bechern und Chips statt Code
**Zeitrahmen:** ca. 60 Minuten (25 Minuten Spiel, 20 Minuten Rechnen, Recap)

## Kontext

Wenn der HotelService extrem langsam antwortet, blockieren seine Aufrufe alle Slots des Booking-Service. Dadurch scheitern auch Anfragen an den FlightService, obwohl dieser einwandfrei funktioniert. Das Bulkhead-Pattern gibt jedem Backend einen eigenen Pool: Läuft einer voll, bleiben die anderen trocken.

Der Code dafür sind sieben Zeilen, deshalb tippen wir ihn nicht ab. Die Arbeit steckt in der Zahl darin: Wie groß muss der Pool sein?

## User Story

Als **Betriebsteam**
möchte ich, **dass Probleme mit einem Backend-Service nicht die Aufrufe an andere Backend-Services beeinträchtigen**,
damit **ein langsamer oder fehlerhafter Service nicht das gesamte System blockiert**.

---

## Teil 1: Der geteilte Pool (ca. 25 Minuten)

Ihr seid die Requests. Der Becher ist der Pool.

### Material

- Ein großer Becher, drei kleine Becher
- Zehn Chips (Münzen, Pokerchips oder Zettel)
- Karteikarten für die Vorhersagen, eine Uhr mit Sekundenanzeige

### Aufbau

1. Zehn Chips in einem Becher, der Pool des Booking-Service.
2. Drei Leute sind Flight, Hotel und Car.
3. Alle anderen sind Requests und stellen sich an.
4. Wer bedient werden will, braucht einen Chip. Zurück gibt es ihn erst, wenn das Backend fertig ist.
5. Flight und Car antworten sofort. Hotel braucht zwanzig Sekunden.

### Ablauf

1. Jedes Team schreibt vorher auf: Was passiert, wenn ein Drittel der Requests zu Hotel will?
2. **Runde 1:** ein Becher für alle drei Backends.
3. **Runde 2:** drei Becher mit vier, drei und drei Chips.
4. Karten aufdecken, vergleichen.

### Was passieren sollte

- Runde 1: Nach wenigen Sekunden liegen alle zehn Chips bei Hotel. Flight- und Car-Requests scheitern, obwohl beide sofort antworten könnten.
- Runde 2: Hotel-Requests scheitern weiter, drei Chips reichen nicht. Flight und Car laufen ungestört. Das ist die Isolation.

---

## Teil 2: Warum ausgerechnet zehn? (ca. 20 Minuten)

Little's Law: `L = λ × W`

- `L` = gleichzeitig laufende Aufrufe, also benötigte Slots
- `λ` = Aufrufe pro Sekunde
- `W` = Dauer eines Aufrufs

Der Booking-Service schickt 20 Hotel-Aufrufe pro Sekunde. Legt euch vorher fest: Braucht ein langsames Backend mehr, weniger oder gleich viele Slots?

| Hotel antwortet in | λ | W | L = Slots |
|---|---|---|---|
| 50 ms | 20 / s | 0,05 s | 1 |
| 500 ms | 20 / s | 0,5 s | 10 |
| 3 s Timeout | 20 / s | 3 s | 60 |

Zehnmal langsamer heißt zehnmal mehr Slots, bei gleichem Verkehr. Wer den Pool klein macht, „um das Backend zu schonen", lehnt ab, was das Backend hätte schaffen können. In Produktion leitet sich die Zahl aus drei Faktoren ab: Connection-Pool des HTTP-Clients, Backend-Kapazität (Replicas × maxConcurrent ≤ Kapazität) und Little's Law aus gemessener Latenz und Ziel-Durchsatz.

---

## Recap: Fünf Fragen

1. Wir haben doch einen Circuit Breaker. Was kann der Bulkhead, was der CB nicht kann?
2. Ein Rate Limit erlaubt hundert Aufrufe pro Sekunde. Und wenn alle hundert gleichzeitig hängen?
3. Wir sind non-blocking, Threads sind bei uns billig. Brauchen wir den Bulkhead dann überhaupt noch?
4. Warum keine Warteschlange? Wäre kurz warten nicht freundlicher?
5. Fünf Replicas mal zehn Slots sind fünfzig Aufrufe. Wen schützt der Bulkhead dann?

Bonus zu Frage 3: Wo ist die Grenze bei non-blocking? Threads sind billig, knapp werden Speicher, Dateideskriptoren und das Backend selbst. Die Grenze verschwindet nicht, sie wird nur vom Betriebssystem gesetzt statt von euch.

Antworten und Anekdoten: [questions/story5.md](../questions/story5.md)

---

## Referenz-Implementierung

Es gibt weiterhin eine Referenz unter `services/booking/story5/` (Go, Semaphore via `chan struct{}`, ein Bulkhead pro Backend, Fail-Fast mit `503`). Im Dashboard zeigt das Bulkhead-Panel unter Story 5 die Zähler pro Backend. Demo nach dem Spiel: Hotel auf „Langsam" stellen, `POST /admin/burst` drücken, die Rejects im Hotel-Pool beobachten, während Flight und Car weiterlaufen.

Wer das Pattern selbst bauen möchte, findet die Bausteine in der Referenz und in `docs/instructions/bulkhead.md`:

- Resilience4j Bulkhead, MicroProfile `@Bulkhead`, Polly, go-resiliency
- Semaphore-basiert (begrenzt parallele Aufrufe) oder Thread-Pool-basiert (eigener Pool pro Service)
- Kombination mit dem Circuit Breaker aus Story 4: Bulkhead und Circuit Breaker ergänzen sich
