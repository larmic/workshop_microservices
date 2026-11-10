# Story 5: Der geteilte Pool

**Thema:** Bulkhead Pattern, als Rollenspiel mit Chips, Trennwänden, Takt und Uhr statt Code
**Zeitrahmen:** ca. 60 Minuten (25 Minuten Spiel, 10 Minuten Demo mit der Referenz, 15 Minuten Rechnen, 10 Minuten Recap)

## Kontext

Antwortet der HotelService extrem langsam, binden seine Aufrufe die Ressourcen des Booking-Service, in einem klassischen Server den gemeinsamen Thread- oder Connection-Pool. Dann warten auch Anfragen an den FlightService, obwohl dieser einwandfrei funktioniert, bis sie beim Aufrufer in den Timeout laufen. Das Bulkhead-Pattern gibt jedem Backend einen eigenen Pool: Läuft einer voll, bleiben die anderen trocken.

Der Code dafür sind sieben Zeilen, deshalb tippen wir ihn nicht ab. Die Arbeit steckt in der Zahl darin: Wie groß muss der Pool sein?

## User Story

Als **Betriebsteam**
möchte ich, **dass Probleme mit einem Backend-Service nicht die Aufrufe an andere Backend-Services beeinträchtigen**,
damit **ein langsamer oder fehlerhafter Service nicht das gesamte System blockiert**.

---

## Teil 1: Das Rollenspiel (ca. 25 Minuten)

Gleiche Schlange und gleiche Zustandskarte wie in Story 4. Neu: ein Takt, Chips im Bulkhead-Tablett, und Hotel ist langsam, nicht kaputt.

### Material

- Sechs Chips (Münzen, Pokerchips oder Zettel)
- Ein Bulkhead-Tablett (3D-gedruckt): drei Spalten Flight, Hotel und Car mit je zwei Chip-Mulden und zwei steckbaren Trennwänden. Ohne Wände ist es ein gemeinsamer Pool mit sechs Chips (Probelauf, Runde 1), mit Wänden hat jeder Service seine zwei Chips (Runde 2). Alternativ: ein großer Becher, dann drei Becher, beschriftet mit F, H und C
- Ein A4-Blatt als Uhr mit acht Feldern, nummeriert 1 bis 8
- Zwanzig Suchkarten, nummeriert 1 bis 20, als Stapel am Anfang der Schlange
- 20 Karten `hotels: []` für die Fallback-Antworten in Runde 2, gebraucht werden 14, der Rest ist Reserve (oder Klebepunkte, mit denen Booking die Suchkarte markiert)
- Die Zustandskarten CLOSED / OPEN / HALF_OPEN und die Fehlerkarten 1 bis 5 aus Story 4, ein Satz pro Tisch

Gesamtliste und Druckvorlagen: [`docs/orga/materialliste.md`](../orga/materialliste.md), druckfertig unter <https://larmic.github.io/workshop_microservices/material/>.

Spielbar ab sieben Personen: Takt, Booking, Hotel, Flight + Car und mindestens drei Suchanfragen. Ab neun Personen spielen Flight und Car getrennt, ab zwölf Personen sind zwei Leute bei Booking (eine am Tablett, eine an den Stapeln).

**Variante ab 14 Personen:** zwei Tische mit eigenem Material und gleicher Ansage. Am Ende vergleichen beide Tische ihre Zahlen.

### Rollen

1. **Takt** ist eine Person aus der Gruppe. Sie zählt laut 1 bis 8, immer wieder von vorn, und klatscht erst, wenn Booking und Hotel „fertig" sagen, etwa alle drei Sekunden. Bei jedem Klatscher nimmt die vorderste Suchanfrage die nächste **Suchkarte** vom Stapel (1 bis 20), gibt sie Booking und stellt sich **sofort hinten wieder an**. Die Runde endet, wenn der Stapel leer ist, also nach dem 20. Klatscher.
2. **Booking** ruft für jede Suchkarte Flight, Hotel und Car nacheinander auf. Jeder Aufruf braucht einen Chip aus dem Tablett. Wartende Karten liegen in einer Reihe bei Booking, die Personen stellen sich trotzdem hinten an. Dazu liegen bei Booking die Zustandskarte CLOSED, die Fehlerkarten 1 bis 5, die Karten `hotels: []` und der Timeout: zwölf Takte ab dem Hotel-Aufruf, das Warten auf einen Chip zählt nicht mit. Läuft der Timeout ab, gelten die Regeln aus Story 4: nächste Fehlerkarte hinlegen, bei der 5 OPEN, ein Erfolg legt alle zurück. Fertige Suchen legt Booking auf zwei Stapel: komplett und `hotels: []`.
3. **Hotel** legt Karte und Chip auf das Feld der Uhr, dessen Zahl Takt gerade gerufen hat. Kommt die Zahl wieder (nach acht Takten), gibt Hotel zuerst beides an Booking zurück, dann kommt die neue Karte. **Flight und Car** antworten sofort: Sie tippen den Chip nur an, er bleibt bei Booking.
4. **Alle anderen** sind Suchanfragen auf `GET /booking/offers`.

### Ablauf

1. **Probelauf:** Tablett ohne Wände, sechs Chips, Hotel gibt schon beim nächsten Takt zurück. Drei Takte, damit alle die Handgriffe kennen. Danach Stapel und Uhr zurücksetzen.
2. **Runde 1, ohne Wände, ein Pool mit sechs Chips:** Alle Aufrufe holen ihren Chip aus demselben Pool. Ist kein Chip frei, wartet die Karte in der Reihe. Kommt eine Karte von Hotel zurück, behält sie ihren Chip für Car. Erst danach geht der Chip an die älteste wartende Karte, also die mit der kleinsten Nummer. Frage vor der Runde: Wie viele Karten liegen nach 20 Takten ohne Chip in der Reihe?
3. **Runde 2, mit Wänden, je zwei Chips für Flight, Hotel und Car:** Vorher Stapel und Uhr zurücksetzen und die beiden Wände einstecken. Jeder Aufruf holt seinen Chip aus der Spalte seines Backends. Ist die Hotel-Spalte leer, legt Booking sofort die Karte `hotels: []` auf die Suchkarte und ruft Car auf. Niemand wartet. Frage vor der Runde: Wie viele der 20 Suchen bekommen einen Hotel-Chip?
4. **Nach Runde 2:** Steht die Zustandskarte noch auf CLOSED?

Gezählt wird nach dem 20. Klatscher: was in der Reihe bei Booking liegt, was auf der Uhr liegt und was auf den beiden Stapeln.

Tipp zur Moderation: Vorhersage per Handzeichen, erst dann auflösen. Wer am nächsten lag, spielt in der nächsten Runde Hotel.

### Was passieren sollte

- Probelauf: Hotel ist schnell, bei Hotel liegt nie mehr als ein Chip. Es wartet niemand.
- Runde 1: Vier, ab Suche 7 wartet jede schon auf Flight. Am Ende liegen sechs Karten auf der Uhr (11 bis 16), vier in der Reihe (17 bis 20), zehn sind komplett. Flight ist gesund und bekommt trotzdem keinen Chip (Recap 01).
- Runde 2: Sechs, die anderen 14 sofort. Einen Hotel-Chip bekommen die Suchen 1, 2, 9, 10, 17 und 18. Nach Takt 20 sind vier davon komplett, 17 und 18 liegen noch auf der Uhr, 14 liegen auf dem Stapel `hotels: []`. Zwei Chips bei acht Takten heißen: im Dauerlauf bekommt jede vierte Suche Hotels, Little's Law zum Anfassen. Flight und Car laufen ungestört. Das ist die Isolation.
- Ehrlicher Vergleich: Runde 1 hat nach 20 Takten mehr komplette Antworten (10 gegen 4), weil Hotel dort alle sechs Chips nehmen durfte, auf Kosten von Flight und Car. Dafür kommt jede Suche ab Nummer 7 später. Der Bulkhead tauscht Vollständigkeit gegen Verfügbarkeit. Kommt der Einwand „Gebt Hotel in Runde 2 doch vier Chips": Bei acht Takten Latenz bräuchte Hotel acht, mehr Chips verschieben das Problem nur (Brücke zu Teil 3). Die Reihe aus Runde 1 ist genau die Queue aus Recap 04.
- Die Zustandskarte: Ja, acht Takte liegen unter zwölf, es gibt keine Fehlerkarte, und Rejects erreichen den Breaker nie. In der Referenz sitzt der Bulkhead vor dem Circuit Breaker.

### Spickzettel Takt für Takt

Ein frei gewordener Chip kann im selben Takt mehrere Aufrufe bedienen, weil Flight und Car ihn sofort zurückgeben. Beispiel Takt 9 in Runde 1: Karte 1 kommt von Hotel und ruft Car, dann geht derselbe Chip an Karte 7, die Flight und Hotel ruft.

**Runde 1** (ohne Wände, sechs Chips)

| Takt | Was passiert | Uhr | Reihe | komplett |
|---|---|---|---|---|
| 1–6 | Karten 1 bis 6 holen je einen Chip, Flight, dann Hotel. Danach ist der Pool leer. | 6 | 0 | 0 |
| 7–8 | Karten 7 und 8 warten auf Flight. | 6 | 2 | 0 |
| 9–14 | Hotel gibt Karten 1 bis 6 zurück, sie rufen Car und sind komplett. Der Chip geht an Karten 7 bis 12. Karten 9 bis 14 kommen neu in die Reihe. | 6 | 2 | 6 |
| 15–16 | Felder 7 und 8 sind leer (in Takt 7 und 8 kam nichts auf die Uhr). Karten 15 und 16 warten zusätzlich. | 6 | 4 | 6 |
| 17–20 | Hotel gibt Karten 7 bis 10 zurück, der Chip geht an Karten 13 bis 16. | 6 | 4 | 10 |

**Runde 2** (mit Wänden, je zwei Chips)

| Takt | Was passiert | Uhr | komplett | `hotels: []` |
|---|---|---|---|---|
| 1–2 | Karten 1 und 2 holen die beiden Hotel-Chips. | 2 | 0 | 0 |
| 3–8 | Hotel-Spalte leer, Karten 3 bis 8 bekommen sofort `hotels: []`. | 2 | 0 | 6 |
| 9–10 | Hotel gibt Karten 1 und 2 zurück, Karten 9 und 10 bekommen die Chips. | 2 | 2 | 6 |
| 11–16 | Karten 11 bis 16 bekommen sofort `hotels: []`. | 2 | 2 | 12 |
| 17–18 | Hotel gibt Karten 9 und 10 zurück, Karten 17 und 18 bekommen die Chips. | 2 | 4 | 12 |
| 19–20 | Karten 19 und 20 bekommen sofort `hotels: []`. | 2 | 4 | 14 |

### Maßstab und Vereinfachungen

- **Maßstab:** 1 Chip steht für 5 Slots, zwei Chips je Backend sind also `maxConcurrent` 10 aus der Referenz. 8 Takte zu 12 Takten entsprechen 2 s Latenz („Langsam") zu 3 s Timeout, 1 Takt ist also 0,25 s. 1 Suchkarte steht für 5 Suchen. Eine Karte pro Takt sind damit 20 Hotel-Aufrufe pro Sekunde, wie in Teil 3.
- **Runde 1 ist kein Abbild der Referenz.** Der gemeinsame Pool ohne Wände steht für den Thread- oder Connection-Pool eines klassischen Servers. Die Go-Referenz hat keinen solchen Pool, dort gibt es nur die drei Bulkheads aus Runde 2.
- **Runde 2 entspricht der Referenz:** eigener Pool je Backend, keine Warteschlange, sofortiger Fallback. Der Client bekommt HTTP 200 mit `hotels: []` sowie die Header `X-Fallback: hotel` und `X-Bulkhead-Full: hotel`. Eine 503 gibt es bei `GET /booking/offers` nicht.
- **Brücke zu Teil 3:** Im Probelauf ist L = 1 Chip, zwei Chips je Backend hätten mit Luft gereicht. Bei langsamem Hotel ist L = 8 Chips, also 40 Slots, der Pool hat 10. Die leere Hotel-Spalte ist dann das Signal, kein Grund, den Pool zu vergrößern.
- **Nicht im Spiel und nicht in der Referenz:** Chips nachkippen (non-blocking), eine Warteschlange vor der Hotel-Spalte, ein zweites Booking. Diese Varianten kommen im Recap als Gedankenspiel vor.

---

## Teil 2: Die Demo mit der Referenz (ca. 10 Minuten)

Vorher im Dashboard alle Services auf „Normal" stellen. Steht Flight noch aus der Story-4-Demo auf „Fehler", verfälscht das den Burst. Dann unter Story 5 Hotel auf „Langsam" stellen (2 s) und den Burst-Knopf drücken. Das Dashboard schickt 20 parallele `GET /booking/offers` an die Referenz (`POST /api/booking-ref-story5/burst`). Alle 20 Antworten sind HTTP 200, „failed" bleibt bei 0. Maßgeblich ist der Zähler „hotel rej": Er steht bei etwa zehn. Der Zähler „X-Bulkhead-Full" zählt jede Antwort mit irgendeinem Reject und kann deshalb höher liegen: Flight-Rejects sind möglich, weil alle 20 Suchen zuerst Flight rufen. Flight und Car laufen weiter, der Hotel-Breaker bleibt CLOSED. Danach Hotel zurück auf „Normal".

---

## Teil 3: Warum ausgerechnet zehn? (ca. 15 Minuten)

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

Zehnmal langsamer heißt zehnmal mehr Slots, bei gleichem Verkehr. Dimensioniert wird aber auf die normale Latenz: Ist der Pool voll, weil das Backend degradiert, ist das das Signal, kein Grund zum Vergrößern. Genau das hat Runde 2 gezeigt. Wer den Pool dagegen von vornherein klein macht, „um das Backend zu schonen", lehnt ab, was das Backend hätte schaffen können. In Produktion leitet sich die Zahl aus drei Faktoren ab: Connection-Pool des HTTP-Clients, Backend-Kapazität (Replicas × maxConcurrent ≤ Kapazität) und Little's Law aus gemessener Latenz und Ziel-Durchsatz. Im Spiel sind es 2 Chips, im Code 10, die Rechnung ist dieselbe.

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

Es gibt weiterhin eine Referenz unter `services/booking/story5/` (Go, Semaphore via `chan struct{}`, ein Bulkhead pro Backend, `maxConcurrent` 10, Fail-Fast ohne Warteschlange). Ist ein Pool voll, gibt `GET /booking/offers` trotzdem HTTP 200 zurück, mit leerer Liste für das Backend und den Headern `X-Fallback` und `X-Bulkhead-Full`. Der Bulkhead sitzt vor dem Circuit Breaker, Rejects zählen deshalb nicht als Fehler. Im Dashboard zeigt das Bulkhead-Panel unter Story 5 die Zähler pro Backend, der Burst-Knopf (`POST /api/booking-ref-story5/burst`) schickt 20 parallele Suchen und zählt die Antworten mit `X-Bulkhead-Full`.

Wer das Pattern selbst bauen möchte, findet die Bausteine in der Referenz und in `docs/instructions/bulkhead.md`:

- Resilience4j Bulkhead, MicroProfile `@Bulkhead`, Polly, go-resiliency
- Semaphore-basiert (begrenzt parallele Aufrufe) oder Thread-Pool-basiert (eigener Pool pro Service)
- Kombination mit dem Circuit Breaker aus Story 4: Bulkhead und Circuit Breaker ergänzen sich
