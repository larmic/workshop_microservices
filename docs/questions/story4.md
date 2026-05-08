# Workshop-Fragen: Bulkhead (Story 4)

Provokante Fragen zum Bulkhead-Pattern und zu seiner Wechselwirkung mit
dem Circuit Breaker aus Story 3. Ziel: das Pattern als gezielte
Antwort auf ein konkretes Problem (Ressourcen-Erschöpfung) verstehen,
nicht als „CB+Bulkhead = doppelt sicher".

---

## Frage 1 — Was bringt Bulkhead, was Circuit Breaker nicht schon kann?

**Frage:** Wir haben in Story 3 schon einen Circuit Breaker. Der schaltet
ab, wenn ein Backend krank ist. Wozu jetzt noch ein Bulkhead?

**Antwort:** Die beiden Patterns lösen verschiedene Probleme:

- **Circuit Breaker** = "Backend ist krank, ich versuche es eine Weile
  gar nicht mehr." Reaktion auf **Fehler-Rate**.
- **Bulkhead** = "Ich verbrenne maximal N gleichzeitige Threads/Slots
  für dieses Backend, egal wie viel Last reinkommt." Reaktion auf
  **Ressourcen-Druck**, nicht auf Fehler.

Klassisches Killer-Szenario, das **nur Bulkhead** löst: Hotel antwortet
in 2 s — fehlerfrei, der CB bleibt CLOSED. Trotzdem laufen unter Last
beliebig viele Threads in Hotel-Calls auf, der Booking-Service hat keine
Threads mehr für Flight oder Car frei. Der CB sieht keinen Grund zu
reagieren. Bulkhead hingegen kappt nach 10 parallelen Hotel-Calls und
lässt Flight/Car ungestört durch.

```
                 ┌────────────────────┐
   Last  ─────►  │  Booking-Service   │
                 │                    │
                 │  ┌──Bulkhead Hotel ─┴─►  Hotel
                 │  │     (max 10)
                 │  ├──Bulkhead Flight ──►  Flight
                 │  │     (max 10)
                 │  └──Bulkhead Car ───►   Car
                 │        (max 10)
                 └────────────────────┘
```

---

## Frage 2 — Werden bei einem Burst von 20 die Backends auch wirklich nur mit 10 gleichzeitigen Aufrufen belastet?

**Frage:** Wenn ich 20 parallele Aufrufe gegen `/booking/offers` schicke
und der Bulkhead `maxConcurrent=10` hat — wie viele Aufrufe sehen
Flight, Hotel und Car wirklich?

**Antwort:** Maximal 10 gleichzeitig — pro Backend. Genau das ist der
Sinn des Patterns. Die anderen 10 werden im Booking-Service abgewiesen,
bevor sie das Backend überhaupt erreichen, und der Aufrufer bekommt
einen Fallback (leeres Array). Der Backend-Service wird also vor Last
geschützt — er sieht von den 20 Client-Requests nur die ersten 10 als
echte Anfragen.

Wichtig: das Limit ist **gleichzeitig**, nicht **insgesamt**. Wenn ein
Slot frei wird (Call ist fertig), darf der nächste rein.

---

## Frage 3 — Warum sind nach einem Burst von 20 mit allen Backends auf "langsam" auf Hotel z.B. nur 7 rejected, nicht 10?

**Frage:** Erwartung wäre `calls=20, rejected=10`. Tatsächlich: `calls=20,
rejected=7`. Warum ist die Zahl unsauber?

**Antwort:** Wegen der **sequenziellen Aufruf-Reihenfolge** (Flight →
Hotel → Car) und Timing-Jitter. Skizze:

```
t=0      ┌─ 20 Goroutinen treffen FLIGHT-Bulkhead
         ├─► 10 bekommen Slot   ──► Flight läuft 2 s
         └─► 10 abgewiesen      ──► gehen sofort zu HOTEL

t≈0      die 10 abgewiesenen treffen HOTEL-Bulkhead
         └─► alle 10 bekommen Slot ──► Hotel läuft 2 s
                                        (Hotel jetzt voll: 10/10)

t≈2000   Flight-Slots werden frei.
         Die 10 erfolgreichen Flight-Goroutinen wollen jetzt zu HOTEL.
         GLEICHZEITIG werden die ersten Hotel-Slots frei (laufen seit t=0).

         Race im 50–200 ms breiten Zeitfenster:
           ─►  3 Goroutinen erwischen einen frei werdenden Slot
           ─►  7 Goroutinen kommen zu früh, Pool noch voll → REJECTED
```

Bei **Flight** ist die Zahl sauber `rejected=10`, weil dort alle 20 zur
gleichen Zeit (t=0) ankommen — keine Streuung, klare Kante.

**Take-away:** Bulkhead garantiert nur die obere Schranke "max 10
parallel". Die genaue Anzahl der Rejects ist timing-abhängig und in
realen Systemen nie deterministisch.

**Wenn man eine saubere Demo will:** nur **Hotel** auf langsam stellen
(Flight/Car normal). Dann sind alle 20 Goroutinen praktisch zeitgleich
bei Hotel und du siehst stabil `rejected ≈ 10`.

### Beobachtung über viele Bursts hinweg

Bei wiederholten Bursts (z.B. 42 Bursts à 20 = 840 Calls je Service)
zeigt sich der Effekt sehr deutlich als **Kaskade**:

```
   flight:   in flight 0/10   calls 840   rejected 341   (~40 %)
   hotel:    in flight 0/10   calls 840   rejected 189   (~22 %)
   car:      in flight 0/10   calls 840   rejected   8   (~ 1 %)
```

Mit jedem Service in der Aufrufkette **sinkt die Reject-Rate**. Das
ist kein Zufall, das ist die Konsequenz des Patterns:

- **Flight** sieht den Burst in voller Wucht. 20 Goroutinen schlagen
  zeitgleich auf, 10 bekommen einen Slot, 10 werden sofort abgewiesen.
- **Hotel** sieht denselben Burst — aber **zeitlich entzerrt**.
  Die 10 von Flight abgewiesenen Goroutinen sind sofort da, die 10
  erfolgreichen kommen 50–200 ms später nach. Hotel hat nur einen Teil
  der Last gleichzeitig zu schultern.
- **Car** sieht das, was vom Burst übrig ist — fast geglättet. Die
  Streuung ist groß genug, dass meist genug Slots frei sind.

**Take-away in einer Zeile:** Eine Bulkhead-Kette ist gleichzeitig
ein **Traffic-Shaper**. Jede Stufe schützt nicht nur sich selbst,
sondern auch alle nachgelagerten Services vor ungeglätteten Bursts —
ein kostenloser Nebeneffekt, der in Produktion oft die spürbarste
Wirkung hat.

**Spicy:** Wer im Aggregator nur das letzte Backend (Car) instrumentiert
und „die Bulkheads funktionieren ja, fast keine Rejects" beobachtet,
hat den Punkt verpasst — Car ist nicht ungestresst, weil sein Bulkhead
gut konfiguriert ist, sondern weil Flight und Hotel die Last vorher
aufgefangen haben.

---

## Frage 4 — Wenn der Bulkhead einen Aufruf abweist, läuft die Anfrage trotzdem zu Hotel und Car weiter. Ist das Teil des Bulkhead-Patterns?

**Frage:** Ein Bulkhead-Reject auf Flight führt nicht zum Abbruch der
gesamten Anfrage — es geht weiter zu Hotel und Car. Ist das, was
Bulkhead macht?

**Antwort:** **Nein.** Streng genommen sind das zwei orthogonale Dinge:

1. **Was Bulkhead macht:** "Pool ist voll → sofort ablehnen, kein
   Queueing." Punkt. Was der Aufrufer mit der Ablehnung tut, ist nicht
   Sache des Patterns.
2. **Was unser Code zusätzlich macht:** Fail-Soft im
   `/booking/offers`-Handler. Bei einem Fehler (egal ob CB-Open,
   Bulkhead-Full oder echter Backend-Fehler) wird ein leeres Array
   zurückgegeben und das nächste Backend trotzdem aufgerufen. Diese
   Logik existiert seit Story 3 für den Circuit Breaker — wir haben den
   Bulkhead-Reject einfach in dieselbe Behandlung einsortiert.

```
fetchOffersWithBHCB(...)
   │
   └─► bh.Execute        ─► Bulkhead voll? → ErrBulkheadFull
          │                                          │
          └─► cb.Execute    ─► CB OPEN? → ErrCircuitOpen
                 │                                   │
                 └─► HTTP-Call                       │
                       ↓                             ↓
                       err? ────────► markFallback() + emptyJSONArray
                                       (── Strategie des Aggregators ──)
```

Im Schreibpfad (`POST /booking/bookings`) ist die Strategie eine andere:
**Fail-Fast** — wenn ein CB OPEN ist oder ein Call fehlschlägt, bricht
die ganze Buchung ab. Eine Teilbuchung wäre für den Kunden schlimmer als
ein 503.

**Take-away:** Bulkhead ist ein **Mechanismus**. Was nach einem Reject
passiert (Fail-Soft, Fail-Fast, Retry mit anderem Backend, …), ist eine
**Designentscheidung des Aufrufers** — die hängt davon ab, ob die
Operation idempotent ist, ob Teilergebnisse sinnvoll sind, usw.

---

## Frage 5 — Im Dashboard sehe ich, dass der Circuit Breaker auf OPEN geht, während der Bulkhead reagiert. Ist das richtig so?

**Frage:** Wenn ich einen Burst mit langsamem Backend mache, wird
manchmal auch der CB OPEN — gleichzeitig mit Bulkhead-Rejects. Beeinflusst
der Bulkhead den CB?

**Antwort:** **Nein, der Bulkhead beeinflusst den CB nicht** — und das
ist Absicht. Im Code:

```go
err := bh.Execute(ctx, func(ctx context.Context) error {
    return cb.Execute(ctx, func(ctx context.Context) error {
        // echter Backend-Call
    })
})
```

Wenn der Bulkhead voll ist, gibt er sofort `ErrBulkheadFull` zurück —
**ohne** `cb.Execute` aufzurufen. Der CB sieht von dem Reject nichts,
seine Counter (`failureCount`, `totalCalls`) bleiben sauber.

**Begründung:**

- Ein Bulkhead-Reject sagt "ich schütze mich selbst", nicht "das Backend
  ist krank". Würden wir das als Failure ans CB melden, würde der CB
  irgendwann fälschlich auf OPEN gehen — obwohl das Backend in Wahrheit
  gesund antwortet. Wir würden uns selbst die Tür zumachen.
- Beide Patterns sind unabhängig: CB = Backend-Gesundheit. Bulkhead =
  Eigenschutz vor Ressourcen-Erschöpfung.

**Wenn der CB im Demo trotzdem OPEN ist, sind die wahrscheinlichen
Ursachen:**

1. **Latenz > 3000 ms** im Chaos-Slider → der `httpClient.Timeout` von
   3 s schlägt zu → das ist ein echter Failure (Timeout) → nach 5
   Timeouts geht der CB auf OPEN.
2. Backend stand im Modus **"Fehler"** statt "Langsam" → 500 zurück →
   echte Failures → CB öffnet nach 5.
3. **Reste aus einem vorherigen Lauf**: der CB bleibt nach OPEN noch
   30 s offen, danach HALF_OPEN. Er räumt sich nicht durch einen Reset
   des Bulkheads auf.

Kontrollierte Reproduktion:

| Latenz | Bulkhead reagiert? | CB reagiert?                         |
|--------|--------------------|--------------------------------------|
| 2000 ms | ja (bei Burst)    | nein — Backend antwortet rechtzeitig |
| 3500 ms | ja                | ja — `httpClient`-Timeout = Failure  |

**Take-away:** CB und Bulkhead sind **komplementär, nicht alternativ**.
Sie können gleichzeitig feuern — und das ist gut so, weil sie auf
verschiedene Symptome reagieren.

---

## Frage 6 — Warum stellt der Bulkhead Anfragen nicht in eine Schlange, sondern lehnt sofort ab?

**Frage:** Wäre es nicht freundlicher, wartende Aufrufe in eine kurze
Queue zu stellen, statt sie sofort mit Fallback abzufertigen?

**Antwort:** Eine begrenzte Queue wäre eine Variante (siehe
Resilience4j: `maxWaitDuration`). In unserer Implementierung haben wir
**bewusst nicht gequeued**, aus zwei Gründen:

1. **Queueing tarnt das Problem.** Wenn der Bulkhead voll ist, ist das
   Backend bereits am Limit. Eine Queue verschiebt das Problem nur in
   die Zukunft und erhöht die End-to-End-Latency. Beim Workshop sieht
   man am sofortigen Reject klar: "wir sind über der Kapazität".
2. **Backpressure-Signal nach oben.** Ein sofortiger Reject (mit 503
   bzw. Fallback) sagt dem Aufrufer: "lass mich in Ruhe, ich bin voll."
   Eine Queue absorbiert das Signal und kaschiert es.

**Take-away:** Queue oder kein Queue ist eine bewusste Entscheidung.
Default für die Demo: **kein Queue**, sofort sichtbar im Counter
`rejected`.

---

## Frage 7 — Wir haben `maxConcurrent=10` hartkodiert. Wo kommt diese Zahl her?

**Frage:** Warum ausgerechnet 10? Wäre 5 sicherer? Oder 100, damit weniger
abgewiesen wird?

**Antwort:** 10 ist eine **Workshop-Default-Zahl**, die fürs Demo gut
funktioniert (20er-Burst → klar sichtbarer Effekt). In Produktion
orientiert man die Zahl an konkreten Begrenzungen weiter unten:

- **Connection-Pool des HTTP-Clients.** Wenn der Booking-Service nur
  20 gleichzeitige Verbindungen zu Hotel halten kann, ist
  `maxConcurrent > 20` sinnlos — die Calls würden auf der TCP-Schicht
  sowieso warten.
- **Thread-Pool des Backends.** Hotel kann z.B. 50 Anfragen parallel
  bedienen. Wenn 5 Booking-Replicas mit `maxConcurrent=20` darauf
  zugreifen, sind das schon 100 — Hotel kippt.
  **Faustregel:** `Booking-Replicas × maxConcurrent ≤ Backend-Kapazität`.
- **Latenz × Throughput-Ziel.** Little's Law: `concurrency = latency
  × throughput`. Bei 50 ms Antwortzeit und 200 req/s Ziel ist die
  nötige Concurrency = 10. Mehr ist Verschwendung, weniger
  würgt die Last.

```
   maxConcurrent zu klein:     unnötige Rejects, Backend hat noch Luft
   maxConcurrent zu groß:      schützt nichts mehr, Bulkhead wird zur Folklore
   genau richtig:              max Last, die Backend + Pool sicher tragen
```

**Spicy Take-away:** Wer `maxConcurrent` aus dem Bauch heraus setzt
(„10 klingt gut"), hat den Bulkhead nicht implementiert, sondern
dekoriert. Das Limit gehört aus **gemessener Backend-Kapazität**
abgeleitet, nicht aus Beispiel-Code übernommen.

---

## Frage 8 — Bulkhead schützt den Aufrufer vor sich selbst. Aber wer schützt das Backend?

**Frage:** Wenn Booking-Service 5 Replicas hat, jede mit
`maxConcurrent=10` für Hotel, dann darf Hotel gleichzeitig 50 Calls
sehen. Ein Bulkhead pro Client schützt das Backend nicht wirklich, oder?

**Antwort:** **Korrekt.** Client-side Bulkhead schützt **den Client**,
nicht das Backend. Beispiel:

```
   Booking-Replica A (BH 10) ──┐
   Booking-Replica B (BH 10) ──┤
   Booking-Replica C (BH 10) ──┼───►  Hotel sieht bis zu 50 parallel
   Booking-Replica D (BH 10) ──┤
   Booking-Replica E (BH 10) ──┘

   Hotel bekommt davon NICHTS mit, dass es Bulkheads gibt.
   Wenn Hotel intern nur 30 Threads hat, wird's bei 50 eng.
```

Komplementäre Patterns auf der Backend-Seite:

1. **Server-side Rate Limiting** — Hotel selbst lehnt nach 30
   gleichzeitigen Calls ab. Schützt Hotel vor jedem Aufrufer (auch vor
   böswilligen).
2. **API-Gateway / Service Mesh** — zentraler Punkt, an dem über alle
   Aufrufer hinweg ein Limit greift.
3. **Backpressure / 429 + Retry-After** — Hotel kommuniziert
   Überlast aktiv, Aufrufer reagieren darauf.

**Spicy Take-away:** Bulkhead allein ist eine **halbierte Lösung**.
Sie macht den eigenen Service stabil, aber das geschützte Backend
braucht ergänzende Mechanismen. Wer nur den Client schützt und glaubt,
das Backend sei auch gerettet, hat das Pattern falsch verstanden.

---

## Sammelthemen für die Diskussion

- Wie würdet ihr `maxConcurrent` in einer realen Anwendung festlegen?
  (Hint: es geht um Thread-Pool-Größe / Connection-Pool-Größe der
  abhängigen Ressourcen, nicht um eine geratene Magic-Number.)
- Wann wäre es sinnvoll, dass ein Bulkhead-Reject **doch** als
  CB-Failure gewertet wird? (Hint: praktisch nie — aber Diskussion über
  Edge-Cases wie "ich kenne die Backend-Kapazität exakt".)
- Was wäre die Saga-Lösung für `/booking/bookings`, wenn Hotel mitten in
  der Buchung den Bulkhead-Reject bekommt? (Brücke zu Story 5.)
