# Story 5: Der geteilte Pool

**Thema:** Bulkhead Pattern, als Rollenspiel mit zwei Stühlen und einer Wand statt Code
**Zeitrahmen:** ca. 45 Minuten (25 Minuten Spiel, 10 Minuten Demo mit der Referenz, 10 Minuten Recap)

## Kontext

Antwortet der HotelService extrem langsam, binden seine Aufrufe die Ressourcen des Booking-Service, in einem klassischen Server den gemeinsamen Thread- oder Connection-Pool. Dann warten auch Anfragen an den FlightService, obwohl dieser einwandfrei funktioniert, bis sie beim Aufrufer in den Timeout laufen. Das Bulkhead-Pattern gibt jedem Backend einen eigenen Pool: Läuft einer voll, bleiben die anderen trocken.

Der Code dafür sind sieben Zeilen, deshalb tippen wir ihn nicht ab. Die Arbeit steckt in der Zahl darin: Wie groß muss der Pool sein?

## User Story

Als **Betriebsteam**
möchte ich, **dass Probleme mit einem Backend-Service nicht die Aufrufe an andere Backend-Services beeinträchtigen**,
damit **ein langsamer oder fehlerhafter Service nicht das gesamte System blockiert**.

---

## Teil 1: Das Rollenspiel „Der geteilte Pool" (ca. 25 Minuten)

Zwei Stühle sind die Threads von Booking. Wer sitzt, hat gerade einen Aufruf laufen. Hotel wird langsam, die Stühle laufen voll, und Requests stehen dahinter, obwohl Flight und Car gesund sind. Dann schiebt der Trainer eine Wand zwischen die Stühle.

### Material

- Ein Kartenstapel 01–27 (A7 hoch, Karton 200–250 g, 8 pro A4-Bogen, 4 Bögen). Auf jeder Karte stehen nur die große Nummer und „Steh auf bei XX", violetter Rahmen. Die STOPP-Karten 04, 15 und 27 sind dunkel. Notfalls von Hand auf Karteikarten
- Drei Service-Schilder A4 hoch auf Karton: Hotel doppelseitig („Hotel · schnell", hinten „Hotel · langsam"), Flight („Flight · schnell"), Car („Car · schnell"). Doppelseitig heißt Rückseite bedrucken oder zwei Blätter zusammenkleben
- Das Trainer-Skript (ein Blatt Papier, beidseitig): vorn Ablauf, Regeln und STOPP-Fragen, hinten Kartentabelle und Durchlauf mit 10 Personen
- Aus dem Raum: zwei Stühle, ein Tisch für alle abgelegten Karten (ein Haufen), eine „Wand" (Pinnwand, Flipchart oder ein hochkant gestellter Tisch)

Kein Einkauf, kein 3D-Druck, kein A3. Gesamtliste: [`docs/orga/materialliste.md`](../orga/materialliste.md), Druckvorlagen und Trainer-Skript unter <https://larmic.github.io/workshop_microservices/material/>.

### Rollen

1. **Der Trainer** spielt keine Rolle. Er baut auf, erklärt die Regeln, leitet die drei STOPPs und schiebt die Wand.
2. **Booking** ist eine Person. Sie hält den Stapel 01–27, ruft die Nummer der obersten Karte und gibt die Karte dem Nächsten in der Schlange. Bei einer STOPP-Karte ruft sie „STOPP" und gibt die Karte dem Trainer. Sonst liest Booking nichts vor.
3. **Hotel, Flight und Car** halten je ihr Schild. Hotel steht bei den Stühlen, Flight und Car daneben. Jeder Service zählt mit, an den Fingern oder per Strichliste: Hotel zählt, wer sich auf einen Stuhl setzt, Flight und Car zählen, wer sie abklatscht.
4. **Alle anderen** sind Requests auf `GET /booking/offers` und stehen in einer Schlange vor Booking.

| Personen | Besetzung |
|---|---|
| 7 | Booking und 6 Requests. Die drei Schilder lehnen an Wand oder Stuhl, der Trainer zählt |
| 8 | dazu Hotel als Person (zählt), Flight und Car lehnen, der Trainer zählt dort |
| 9 | eine Person hält Flight und Car zusammen |
| 10 | Normalfall: Booking, Hotel, Flight, Car, 6 Requests |
| 11–13 | mehr Requests, sonst alles gleich |

Unter sieben Personen nicht spielbar, Runde 1 braucht sechs Requests. **Ab 14 Personen** zwei Gruppen parallel, jede mit eigenem Satz (Karten, Schilder, zwei Stühle, Wand).

### Aufbau

- Zwei Stühle nebeneinander, Sitzfläche zu Booking, damit wer sitzt, die Nummern hört.
- Hotel steht mit Schild bei den Stühlen, Flight und Car daneben. Alle Schilder zeigen „schnell".
- Ein Tisch, auf dem alle abgelegten Karten auf einem Haufen landen.
- Die Wand liegt bereit, steht aber noch nicht zwischen den Stühlen.

### Die Regeln

Bei jeder Nummer, immer in dieser Reihenfolge:

1. **Booking ruft die Nummer.**
2. **Wer sitzt** und diese Nummer auf seiner Karte hat, steht auf, klatscht Car ab, legt die Karte auf den Tisch und stellt sich hinten in die Schlange.
3. **Wer hinter den Stühlen wartet**, setzt sich auf den frei gewordenen Stuhl (wer zuerst kam, zuerst) und klatscht dabei Flight ab.
4. **Booking gibt die Karte dem Nächsten in der Schlange.** Ist ein Stuhl frei: Flight abklatschen, hinsetzen. Sonst hinter die Stühle stellen und warten.

Links oder rechts ist egal. Hotel zählt jeden, der sich setzt.

**Zusatzregel ab Runde 2** (nach der Wand): Der Request klatscht sofort Flight ab. Ist der Stuhl auf der Hotel-Seite frei, setzt er sich (Hotel zählt). Ist er besetzt, wartet er **nicht**, sondern klatscht sofort Car ab, legt die Karte auf den Tisch und stellt sich hinten an. Das ist die Antwort `hotels: []` mit `X-Fallback` und `X-Bulkhead-Full`. Wer auf dem Hotel-Stuhl sitzt, klatscht beim Aufstehen Car ab. Der Stuhl auf der Flight-/Car-Seite bleibt frei: Flight und Car sind schnell, ihr Abteil wird nie voll.

### Ablauf

1. **Aufbau und Ansage** (3 min): Stühle, Services, Regeln.
2. **Probe 01–03** (2 min): Alle Schilder auf „schnell". Es sitzt immer nur einer, der zweite Stuhl bleibt leer.
3. **STOPP 1, Karte 04** (2 min): „Wie viele Stühle waren gleichzeitig besetzt?" Einer: Ankünfte × Dauer = 1 × 1 = 1. Danach dreht Hotel sein Schild auf „langsam": „Ab jetzt sitzt man vier Nummern." Vorhersage per Handzeichen: Was passiert? Die Zähler fangen bei 0 an, die Probe zählt nicht.
4. **Runde 1, 05–14** (4 min): Nach 06 sind beide Stühle voll, ab 07 stehen Requests dahinter, die Reihe wächst. Flight und Car stehen herum.
5. **STOPP 2, Karte 15** (5 min): Standbild, zwei sitzen (09, 10), vier stehen dahinter (11–14).
   - „Ihr hinten, worauf wartet ihr?" Auf einen Stuhl. Zuerst wollt ihr nur Flight, und Flight ist gesund.
   - „Flight, Car, wie viele kamen bei euch an?" Flight 6, Car 4.
   - „Hat der Circuit Breaker etwas gemerkt?" Nein: 2 s liegen unter 3 s Timeout, langsam ist kein Fehler, CLOSED.
   - „Karte 14, wann stehst du auf?" Bei 26, so spät käme die Antwort.

   Dann alle Karten auf den Tisch, alle stellen sich neu an. **Der Trainer schiebt die Wand zwischen die Stühle.** Hotel stellt sich auf die eine Seite, Flight und Car auf die andere: „Dieselben zwei Stühle, nur mit einer Wand." Zusatzregel erklären. Vorhersage: Wie viele der nächsten elf bekommen Hotels? Die Zähler fangen neu an, die Werte aus Runde 1 merken.
6. **Runde 2, 16–26** (3 min): Auf dem Hotel-Stuhl sitzen nacheinander 16, 20 und 24. Alle anderen klatschen Flight und Car ab und sind sofort fertig. Niemand steht an.
7. **STOPP 3, Karte 27** (6 min): Zähler abfragen und mit Runde 1 vergleichen (Tabelle unten).
   - „Welche Runde war besser?" Ehrlich: Runde 1 hatte mehr Hotels, in Runde 2 wartet niemand, und mehr sind beantwortet.
   - „Wie viele Stühle hätte Hotel gebraucht?" Vier: Jede Nummer kommt einer dazu, jeder bleibt vier. Hotel hat einen, deshalb bekommt jeder vierte Hotels (16, 20, 24).
   - „Steht der Breaker noch auf CLOSED?" Ja, ein besetzter Stuhl ist kein Fehler.
   - „Wem gehören die Stühle?" Booking. Hotel weiß nichts davon. Die Wand schützt Booking, nicht Hotel, fünf Booking-Instanzen hätten fünf Hotel-Stühle.

Tipp zur Moderation: Vorhersage per Handzeichen, erst dann auflösen.

### Was passieren sollte

Gezählt wird bei STOPP 3 über die Services:

| | Runde 1 | Runde 2 | was es heißt |
|---|---|---|---|
| Hotel: bei mir gesessen | 6 | 3 | Runde 1 lieferte mehr Hotels |
| Flight: abgeklatscht | 6 | 11 | Flight war in Runde 1 gesund und trotzdem ausgehungert |
| Car: abgeklatscht = beantwortet | 4 | 10 | Verfügbarkeit vor Vollständigkeit |
| davon komplett mit Hotels | 4 | 2 | ehrlicher Tausch |
| davon sofort ohne Hotels | 0 | 8 | Fallback `hotels: []` |
| stehen hinter den Stühlen | 4 | 0 | die Queue ist weg |
| Breaker | CLOSED | CLOSED | langsam ist kein Fehler |
| Stühle bei Booking | 2 | 2 | gleich viele, nur durch eine Wand getrennt |

Der Bulkhead tauscht Vollständigkeit gegen Verfügbarkeit. Kommt der Einwand „Gebt Hotel doch beide Stühle": Bei langsamem Hotel bräuchte es vier, mehr Stühle verschieben das Problem nur (Rechnung unter „Zum Nachlesen“). Die Reihe hinter den Stühlen ist genau die Queue aus Recap 04.

### Durchlauf mit 10 Personen

Booking, Hotel, Flight, Car und die Requests P1 bis P6. **Probe:** 01 an P1, P1 setzt sich. Bei 02 steht P1 auf (Car), P2 bekommt 02 und setzt sich. Bei 03 genauso mit P3. Bei 04 steht P3 auf, Booking ruft STOPP. Die Schlange ist danach P4, P5, P6, P1, P2, P3.

**Runde 1** (Zähler ab 05):

| Nr. | Karte an | Was passiert | Stuhl | Stuhl | dahinter |
|---|---|---|---|---|---|
| 05 | P4 | setzt sich | P4 (05) | – | – |
| 06 | P5 | setzt sich, beide voll | P4 | P5 (06) | – |
| 07 | P6 | kein Stuhl frei: dahinter | P4 | P5 | P6 (07) |
| 08 | P1 | dahinter | P4 | P5 | P6, P1 (08) |
| 09 | P2 | P4 (05) steht auf, P6 (07) setzt sich, P2 dahinter | P6 | P5 | P1, P2 (09) |
| 10 | P3 | P5 (06) steht auf, P1 (08) setzt sich, P3 dahinter | P6 | P1 | P2, P3 (10) |
| 11 | P4 | niemand steht auf, P4 dahinter | P6 | P1 | P2, P3, P4 (11) |
| 12 | P5 | dahinter | P6 | P1 | P2 … P5 (12) |
| 13 | P6 | P6 (07) steht auf und stellt sich an, P2 (09) setzt sich, P6 bekommt 13 und geht dahinter | P2 | P1 | P3 … P6 (13) |
| 14 | P1 | P1 (08) steht auf, P3 (10) setzt sich, P1 bekommt 14 und geht dahinter | P2 | P3 | P4, P5, P6, P1 (14) |
| 15 | STOPP 2 | Booking gibt die Karte dem Trainer | P2 (09) | P3 (10) | 11, 12, 13, 14 |

Zähler Runde 1: Hotel 6 (05–10 haben sich gesetzt), Flight 6, Car 4 (05–08 fertig).

**Runde 2** (neu angestellt, Zähler auf 0):

| Nr. | Was passiert | Hotel-Stuhl |
|---|---|---|
| 16 | Flight, setzt sich auf den Hotel-Stuhl | 16 |
| 17–19 | Flight, Hotel besetzt: sofort Car, Tisch, anstellen | 16 |
| 20 | 16 steht auf (Car), 20 setzt sich | 20 |
| 21–23 | sofort fertig | 20 |
| 24 | 20 steht auf, 24 setzt sich | 24 |
| 25–26 | sofort fertig | 24 |
| 27 | STOPP 3, Booking gibt die Karte dem Trainer | 24 |

Zähler Runde 2: Hotel 3 (16, 20, 24), Flight 11, Car 10 (24 sitzt noch). Wer welche Karte bekommt, hängt davon ab, wie sich die Leute anstellen, die Kartenzahlen nicht.

### Die Karten

| Karte | Phase | Text |
|---|---|---|
| 01 | Probe | Steh auf bei 02 |
| 02 | Probe | Steh auf bei 03 |
| 03 | Probe | Steh auf bei 04 |
| 04 | – | STOPP 1 |
| 05 | Runde 1 | Steh auf bei 09 |
| 06 | Runde 1 | Steh auf bei 10 |
| 07 | Runde 1 | Steh auf bei 13 |
| 08 | Runde 1 | Steh auf bei 14 |
| 09 | Runde 1 | Steh auf bei 17 |
| 10 | Runde 1 | Steh auf bei 18 |
| 11 | Runde 1 | Steh auf bei 21 |
| 12 | Runde 1 | Steh auf bei 22 |
| 13 | Runde 1 | Steh auf bei 25 |
| 14 | Runde 1 | Steh auf bei 26 |
| 15 | – | STOPP 2 |
| 16–26 | Runde 2 | Steh auf bei Nummer + 4 (16 → 20, 17 → 21, … 26 → 30) |
| 27 | – | STOPP 3 |

Warum die Zahlen stimmen: In der Probe ist Hotel normal, man sitzt eine Nummer. In Runde 1 und 2 ist Hotel langsam, man sitzt vier Nummern ab dem Hinsetzen. Runde 1 läuft deterministisch: 07 wartet, setzt sich bei 09 und steht bei 13 auf. Die Karten 09–14 nennen Zahlen hinter STOPP 2, das ist egal, der STOPP räumt auf. Sie zeigen, wann die Antwort käme: Karte 14 erst bei 26, also 12 Nummern oder 6 s nach dem Ankommen. In Runde 2 braucht nur die Zahl, wer den Hotel-Stuhl bekommt. 25 und 26 nennen 29 und 30, sie sitzen nie.

### Maßstab zur Referenz

- **1 Nummer = 0,5 s.** Hotel normal (0,5 s) ist eine Nummer, „Langsam" (2 s) vier Nummern, der Timeout (3 s) sechs Nummern.
- **Pro Nummer kommt ein Request**, also 20 Requests pro Sekunde wie unter „Zum Nachlesen“.
- **1 Stuhl = 10 Slots** (`MaxConcurrent` 10). Runde 1 ist ein gemeinsamer Pool aus zwei Stühlen (20 Slots). In Runde 2 hat Hotel einen Stuhl, also genau die 10 Slots der Referenz, Flight und Car das Abteil auf der anderen Seite.
- **Stuhl** = Thread oder Slot, **Wand** = Bulkhead, **Abteil** = `max[service]`, **hinter den Stühlen stehen** = Queue, **sofort fertig ohne Hotels** = `fallback()` mit `hotels: []`, **Schild „langsam"** = Modus „Langsam" des Hotel-Service.

### Bewusste Vereinfachungen

Offen ansagen:

- **Runde 1 ist kein Abbild der Referenz.** Die zwei Stühle ohne Wand stehen für den Thread- oder Connection-Pool eines klassischen Servers. Die Go-Referenz hat keinen solchen Pool und keine Queue, dort gibt es nur die Bulkheads aus Runde 2.
- **Runde 2 entspricht der Referenz:** eigenes Abteil, keine Warteschlange, sofortiger Fallback. Der Client bekommt HTTP 200 mit `hotels: []` und den Headern `X-Fallback: hotel` und `X-Bulkhead-Full: hotel`. Eine 503 gibt es bei `GET /booking/offers` nicht.
- In der Referenz haben Flight und Car je ein eigenes Abteil mit 10 Slots. Im Spiel teilen sie sich den zweiten Stuhl, er bleibt leer, weil beide schnell sind. Flight und Car abzuklatschen kostet keine Zeit.
- Wer sich hinten anstellt, spielt den nächsten Request, keinen Retry.
- Der Breaker liegt nicht im Spiel, er wird nur gefragt. In der Referenz sitzt der Bulkhead vor dem Circuit Breaker, Rejects erreichen ihn nie.
- **Nicht im Spiel und nicht in der Referenz:** eine Warteschlange vor dem Hotel-Stuhl, ein zweites Booking. Beides kommt im Recap als Gedankenspiel vor.

---

## Teil 2: Die Demo mit der Referenz (ca. 10 Minuten)

Vorher im Dashboard alle Services auf „Normal" stellen. Steht Flight noch aus der Story-4-Demo auf „Fehler", verfälscht das den Burst. Dann unter Story 5 Hotel auf „Langsam" stellen (2 s) und den Burst-Knopf drücken. Das Dashboard schickt 20 parallele `GET /booking/offers` an die Referenz (`POST /api/booking-ref-story5/burst`). Alle 20 Antworten sind HTTP 200, die Ergebniszeile zeigt „0 fail“. Maßgeblich ist der Zähler „hotel rej": Er steht bei etwa zehn. Der Zähler „X-Bulkhead-Full" zählt jede Antwort mit irgendeinem Reject und kann deshalb höher liegen: Flight-Rejects sind möglich, weil alle 20 Suchen zuerst Flight rufen. Flight und Car laufen weiter. Der Hotel-Breaker merkt nichts, 2 s liegen unter dem Timeout von 3 s (das Story-5-Panel zeigt ihn nicht). Vor einem zweiten Burst „Reset“ drücken, sonst zählt „hotel rej“ weiter. Danach Hotel zurück auf „Normal".

---

## Zum Nachlesen: Warum ausgerechnet zehn?

Little's Law: `L = λ × W`

- `L` = gleichzeitig laufende Aufrufe, also benötigte Slots
- `λ` = Aufrufe pro Sekunde
- `W` = Dauer eines Aufrufs

Der Booking-Service schickt 20 Hotel-Aufrufe pro Sekunde. Legt euch vorher fest: Braucht ein langsames Backend mehr, weniger oder gleich viele Slots?

| Hotel antwortet in | λ | W | L = Slots |
|---|---|---|---|
| 50 ms | 20 / s | 0,05 s | 1 |
| 500 ms | 20 / s | 0,5 s | 10 |
| 2 s („Langsam") | 20 / s | 2 s | 40 |
| 3 s Timeout | 20 / s | 3 s | 60 |

Zehnmal langsamer heißt zehnmal mehr Slots, bei gleichem Verkehr. Dimensioniert wird aber auf die normale Latenz: Ist der Pool voll, weil das Backend degradiert, ist das das Signal, kein Grund zum Vergrößern. Genau das hat Runde 2 gezeigt: Der Hotel-Stuhl war voll, und das war richtig so. Wer den Pool dagegen von vornherein klein macht, „um das Backend zu schonen", lehnt ab, was das Backend hätte schaffen können. In Produktion leitet sich die Zahl aus drei Faktoren ab: Connection-Pool des HTTP-Clients, Backend-Kapazität (Replicas × maxConcurrent ≤ Kapazität) und Little's Law aus gemessener Latenz und Ziel-Durchsatz. Im Spiel: normal 1 Stuhl, langsam 4, Hotel hat 1 Stuhl = 10 Slots. Die Rechnung ist dieselbe.

---

## Recap: Vier Fragen

1. Wir haben doch einen Circuit Breaker. Was kann der Bulkhead, was der CB nicht kann?
2. Wir sind non-blocking, Threads sind bei uns billig. Brauchen wir den Bulkhead dann überhaupt noch?
3. Warum keine Warteschlange? Wäre kurz warten nicht freundlicher?
4. Fünf Booking-Instanzen mit je zehn Slots: 50 Aufrufe bei Hotel. Wen schützt der Bulkhead dann?

Bonus zu Frage 2: Wo ist die Grenze bei non-blocking? Threads sind billig, knapp werden Speicher, Dateideskriptoren und der Downstream selbst. Die Grenze verschwindet nicht, sie wird nur vom Betriebssystem gesetzt statt von euch.

Antworten und Anekdoten: [questions/story5.md](../questions/story5.md)

---

## Referenz-Implementierung

Es gibt weiterhin eine Referenz unter `services/booking/story5/` (Go, Semaphore via `chan struct{}`, ein Bulkhead pro Backend, `maxConcurrent` 10, Fail-Fast ohne Warteschlange). Ist ein Pool voll, gibt `GET /booking/offers` trotzdem HTTP 200 zurück, mit leerer Liste für das Backend und den Headern `X-Fallback` und `X-Bulkhead-Full`. Der Bulkhead sitzt vor dem Circuit Breaker, Rejects zählen deshalb nicht als Fehler. Im Dashboard zeigt das Bulkhead-Panel unter Story 5 die Zähler pro Backend, der Burst-Knopf (`POST /api/booking-ref-story5/burst`) schickt 20 parallele Suchen und zählt die Antworten mit `X-Bulkhead-Full`.

Wer das Pattern selbst bauen möchte, findet die Bausteine in der Referenz und in `docs/instructions/bulkhead.md`:

- Resilience4j Bulkhead, MicroProfile `@Bulkhead`, Polly, go-resiliency
- Semaphore-basiert (begrenzt parallele Aufrufe) oder Thread-Pool-basiert (eigener Pool pro Service)
- Kombination mit dem Circuit Breaker aus Story 4: Bulkhead und Circuit Breaker ergänzen sich
