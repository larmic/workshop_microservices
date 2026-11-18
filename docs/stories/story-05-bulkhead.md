# Story 5: Der geteilte Pool

**Thema:** Bulkhead Pattern, als Rollenspiel mit Chips statt Code
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

Booking kann nur eine feste Zahl Requests gleichzeitig bearbeiten, im Spiel neun. Jeder Request holt sich bei Flight, Hotel und Car je einen Chip, das ist seine Antwort. Wird Hotel langsam, bleiben die Requests dort stehen, bis alle neun Plätze voll sind und nichts mehr geht, obwohl Flight und Car gesund sind. Dann bekommt jeder Service ein eigenes Kontingent von drei Plätzen: Ist das von Hotel voll, gibt Booking sofort einen Default-Chip, und alle anderen laufen weiter durch.

### Material

- Chips in vier Farben, je 20 Stück: Flight, Hotel, Car und Default (die leere Hotelliste)
- Drei Service-Schilder A4 hoch auf Karton: Hotel doppelseitig („Hotel · schnell", hinten „Hotel · langsam"), Flight („Flight · schnell"), Car („Car · schnell"). Doppelseitig heißt Rückseite bedrucken oder zwei Blätter zusammenkleben
- Das Trainer-Skript (ein Blatt Papier, beidseitig): vorn Ablauf, Regeln und STOPP-Fragen, hinten Durchlauf mit 14 Personen
- Optional: ein Schild für Booking („9 parallel · Bulkhead: Hotel höchstens 3“) und zwei Kreppband-Felder auf dem Boden („bei Booking“ und „wartet bei Hotel“), damit Booking sieht statt zählt

Außer den Chips kein Einkauf, kein 3D-Druck, kein A3. Gesamtliste: [`docs/orga/materialliste.md`](../orga/materialliste.md), Druckvorlagen und Trainer-Skript unter <https://larmic.github.io/workshop_microservices/material/>.

### Rollen

1. **Der Trainer** spielt keine Rolle. Er sagt die Stufen an, gibt Hotel die Anweisung, langsam zu werden, stoppt im richtigen Moment und leitet die Gespräche.
2. **Booking** ist eine Person. Sie lässt die Requests einzeln herein, achtet darauf, dass höchstens neun gleichzeitig in Arbeit sind und alle richtig stehen, und gibt den Default-Chip aus. Ab Stufe 3 achtet sie zusätzlich auf das Kontingent: Bei Hotel warten höchstens drei.
3. **Flight, Hotel und Car** halten je ihr Schild und geben ihre Chips aus, einen pro Request. Ist Hotel langsam, gibt Hotel bis zur Ansage des Trainers nichts aus.
4. **Alle anderen** sind Requests auf `GET /booking/offers` und stehen in einer Schlange vor Booking.

Die Zahl der parallelen Plätze richtet sich nach der Gruppe: Es braucht einen Request mehr als Plätze, damit in Stufe 2 einer an der Tür stehen bleibt.

| Personen | parallel | Kontingent pro Service | Requests |
|---|---|---|---|
| 8–10 | 3 | 1 | 4–6 |
| 11–13 | 6 | 2 | 7–9 |
| ab 14 | 9 | 3 | 10 und mehr |

Bei 13 Personen geht auch neun parallel, wenn der Trainer selbst den zehnten Request spielt oder eine Person Flight und Car zusammen hält. Unter acht Personen nicht spielbar. **Ab 18 Personen** zwei Gruppen parallel, die Chips werden geteilt (je zehn pro Farbe).

### Aufbau

- Booking steht vor der Schlange, Flight, Hotel und Car stehen nebeneinander in dieser Reihenfolge, so dass man von Booking aus nacheinander hingeht. Alle Schilder zeigen „schnell".
- Auf jedem Schild liegt ein Chip der passenden Farbe, damit die Zuordnung klar ist. Die Default-Chips hält Booking.

### Die Regeln

Für jeden Request, immer in dieser Reihenfolge:

1. **Booking lässt den Nächsten herein**, sobald der Vorige fertig ist oder wartet. Sind alle Plätze belegt, wartet der Request an der Tür und bekommt nichts, auch keinen Flight-Chip.
2. **Chips abholen: Flight, dann Hotel, dann Car.** Gibt Hotel keinen Chip, bleibt der Request bei Hotel stehen und kommt nicht weiter zu Car.
3. **Mit allen drei Chips** ist der Request fertig: Er behält die Chips und stellt sich hinten wieder an.

**Zusatzregel ab Stufe 3** (Bulkhead): Bei Hotel dürfen höchstens drei gleichzeitig warten. Wer als Vierter kommt, wartet nicht, sondern bekommt von Booking sofort einen Default-Chip, holt Car und ist fertig. Das ist die Antwort `hotels: []` mit `X-Fallback` und `X-Bulkhead-Full`. Flight und Car haben genauso ein Kontingent von drei, es wird nie voll, weil beide sofort antworten.

### Ablauf

Die Zahlen gelten für neun parallele Plätze und zehn Requests (R1 bis R10).

1. **Aufbau und Ansage** (3 min): Rollen, Reihenfolge Flight, Hotel, Car, „Booking schafft neun gleichzeitig“.
2. **Stufe 1: ohne Bulkhead, alles schnell** (2 min): Zwei Requests laufen durch und gehen mit drei Chips. So lernen alle den Ablauf.
3. **Stufe 2: ohne Bulkhead, Hotel langsam** (4 min): Hotel dreht sein Schild auf „langsam". Ansage: „Hotel antwortet noch, nur langsam, etwa 2 Sekunden. In 2 Sekunden kommen Hunderte Requests. Im Spiel heißt das: Hotel gibt bis zum STOPP nichts mehr.“ R1 bis R9 holen Flight und bleiben bei Hotel stehen. R10 steht an der Tür.
4. **STOPP 1** (5 min): Standbild, nichts bewegt sich mehr.
   - „Ihr bei Hotel, was habt ihr in der Hand?" Nur Flight. Wartet auf Hotel.
   - „R10, was hast du?" Nichts, nicht einmal Flight, und Flight ist gesund.
   - „Car, wie viele Chips hast du ausgegeben?" Keinen, obwohl Car gesund ist.
   - „Was ist kaputt?" Nichts. Hotel ist nur langsam, ein Fehler entsteht nirgends.

   Dann alle Chips zurück an die Services, alle stellen sich neu an. **Der Trainer sagt den Bulkhead an**: „Dieselben neun Plätze, aber jeder Service bekommt drei davon.“ Zusatzregel erklären, Hotel dreht auf „schnell".
5. **Stufe 3: mit Bulkhead, alles schnell** (2 min): Zwei Requests laufen durch wie in Stufe 1. Kurz fragen: „Was hat sich geändert?" Nichts. Solange alles läuft, kostet der Bulkhead nichts.
6. **Stufe 4: mit Bulkhead, Hotel langsam** (3 min): Hotel dreht wieder auf „langsam". Vorhersage per Handzeichen: Wie viele bekommen eine Antwort? R1 bis R3 bleiben bei Hotel stehen, ab R4 gibt es Flight, Default und Car, und wer fertig ist, stellt sich wieder an. Weiterlaufen lassen, bis jeder einmal durch ist.
7. **STOPP 2** (6 min): Mit Stufe 2 vergleichen (Tabelle unten).
   - „Wie viele stehen fest, wie viele sind frei?" Stufe 2: alle fest. Stufe 4: drei fest, alle anderen fertig und wieder in der Schlange.
   - „Was hat die Grenze gekostet?" Weniger Hotels. Verfügbarkeit vor Vollständigkeit.
   - „Hilft es, Hotel alle neun Plätze zu geben?" Nein, das ist Stufe 2.
   - „Wem gehört die Grenze?" Booking. Hotel weiß nichts davon. Fünf Booking-Instanzen hätten je drei Hotel-Plätze, Hotel sieht fünfzehn.

   Zum Schluss: „2 Sekunden sind um.“ Hotel gibt den drei Wartenden ihre Chips, das Kontingent ist wieder frei. Hotel war nie kaputt.

Tipp zur Moderation: Vorhersage per Handzeichen, erst dann auflösen.

### Was passieren sollte

| | Stufe 2 | Stufe 4 | was es heißt |
|---|---|---|---|
| warten bei Hotel | 9 | 3 | das Kontingent begrenzt den Schaden |
| an der Tür | 1 | 0 | Booking nimmt wieder an |
| fertig, mit Hotel | 0 | 0 | Hotel bleibt langsam |
| fertig, mit Default | 0 | 7 | Fallback `hotels: []` |
| Car-Chips ausgegeben | 0 | 7 | Car war in Stufe 2 gesund und trotzdem ausgehungert |
| frei für den nächsten Request | 0 | 7 | Verfügbarkeit vor Vollständigkeit |
| Plätze bei Booking | 9 | 9 | gleich viele, nur aufgeteilt |

Der Bulkhead tauscht Vollständigkeit gegen Verfügbarkeit. Kommt der Einwand „Gebt Hotel doch mehr Plätze": Bei langsamem Hotel bräuchte es viel mehr, mehr Plätze verschieben das Problem nur (Rechnung unter „Zum Nachlesen“). Wer an der Tür steht, ist genau die Queue aus Recap 03.

### Durchlauf mit 14 Personen

Booking, Flight, Hotel, Car und die Requests R1 bis R10. Nach jeder Stufe gehen alle Chips zurück an die Services.

**Stufe 1** (ohne Bulkhead, alles schnell):

| Request | Flight | Hotel | Car | Ergebnis | belegt (von 9) |
|---|---|---|---|---|---|
| R1 | ✓ | ✓ | ✓ | fertig, drei Chips | 1, dann 0 |
| R2 | ✓ | ✓ | ✓ | fertig, drei Chips | 1, dann 0 |

**Stufe 2** (ohne Bulkhead, Hotel langsam):

| Request | Flight | Hotel | Car | Ergebnis | belegt (von 9) |
|---|---|---|---|---|---|
| R1 | ✓ | wartet | – | steht bei Hotel | 1 |
| R2–R9 | ✓ | wartet | – | stehen bei Hotel | 2 … 9 |
| R10 | – | – | – | kommt nicht herein, nicht einmal Flight | 9, voll |

**Stufe 3** (mit Bulkhead, alles schnell):

| Request | Flight | Hotel | Car | Ergebnis | bei Hotel (von 3) |
|---|---|---|---|---|---|
| R1 | ✓ | ✓ | ✓ | fertig, drei Chips | 1, dann 0 |
| R2 | ✓ | ✓ | ✓ | fertig, drei Chips | 1, dann 0 |

**Stufe 4** (mit Bulkhead, Hotel langsam):

| Request | Flight | Hotel | Car | Ergebnis | bei Hotel (von 3) |
|---|---|---|---|---|---|
| R1–R3 | ✓ | wartet | – | stehen bei Hotel | 1 … 3, voll |
| R4 | ✓ | Default von Booking | ✓ | fertig ohne Hotel, stellt sich wieder an | 3 |
| R5–R10 | ✓ | Default von Booking | ✓ | fertig ohne Hotel | 3 |

In Stufe 4 sind sieben Personen gebunden (drei Requests bei Hotel, die drei Services, Booking), die übrigen sieben Requests sind frei. Für Stufe 4 allein würden vier Requests reichen, die zehn braucht Stufe 2.

**Chips pro Stufe:** Flight 2 / 9 / 2 / 10, Hotel 2 / 0 / 2 / 0, Car 2 / 0 / 2 / 7, Default 0 / 0 / 0 / 7. Ohne Rückgabe nach jeder Stufe reichen 20 Flight-Chips nicht.

### Maßstab zur Referenz

- **Ein Platz bei Booking** = ein Thread oder Slot. Die neun gemeinsamen Plätze in Stufe 2 stehen für den Thread- oder Connection-Pool eines klassischen Servers.
- **Das Kontingent pro Service** = der Bulkhead, `max[service]`. Im Spiel drei, in der Referenz zehn (`MaxConcurrent`).
- **An der Tür warten** = die Queue des gemeinsamen Pools.
- **Default-Chip von Booking** = `fallback()` mit `hotels: []`, HTTP 200.
- **Schild „langsam"** = Modus „Langsam" des Hotel-Service (2 s), unter dem Timeout von 3 s.
- **Ein Request nach dem anderen**: In echt kommen Hunderte pro Sekunde, deshalb wirkt das langsame Hotel im Spiel wie stehen geblieben.

### Bewusste Vereinfachungen

Offen ansagen:

- **Stufe 2 ist kein Abbild der Referenz.** Die neun gemeinsamen Plätze stehen für den Pool eines klassischen Servers. Die Go-Referenz hat keinen solchen Pool und keine Queue, dort gibt es nur die Bulkheads aus Stufe 3 und 4.
- **Stufe 4 entspricht der Referenz:** eigenes Kontingent, keine Warteschlange, sofortiger Fallback. Der Client bekommt HTTP 200 mit `hotels: []` und den Headern `X-Fallback: hotel` und `X-Bulkhead-Full: hotel`. Eine 503 gibt es bei `GET /booking/offers` nicht.
- Flight und Car haben ein Kontingent wie Hotel, es läuft im Spiel nie voll, weil beide sofort antworten.
- Wer sich hinten anstellt, spielt den nächsten Request, keinen Retry.
- Der Circuit Breaker spielt nicht mit. Langsam ist kein Fehler, er bliebe CLOSED. Die Frage dazu kommt im Recap.
- **Nicht im Spiel und nicht in der Referenz:** eine Warteschlange vor dem Hotel-Kontingent, ein zweites Booking. Beides kommt im Recap als Gedankenspiel vor.

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

Zehnmal langsamer heißt zehnmal mehr Slots, bei gleichem Verkehr. Dimensioniert wird aber auf die normale Latenz: Ist der Pool voll, weil das Backend degradiert, ist das das Signal, kein Grund zum Vergrößern. Genau das hat Stufe 4 gezeigt: Das Hotel-Kontingent war voll, und das war richtig so. Wer den Pool dagegen von vornherein klein macht, „um das Backend zu schonen", lehnt ab, was das Backend hätte schaffen können. In Produktion leitet sich die Zahl aus drei Faktoren ab: Connection-Pool des HTTP-Clients, Backend-Kapazität (Replicas × maxConcurrent ≤ Kapazität) und Little's Law aus gemessener Latenz und Ziel-Durchsatz.

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
