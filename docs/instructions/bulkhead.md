# Bulkhead — Workshop-Notizen

> Trainer-Notizen. Kurz und knapp gehalten — fokussiert auf die typischen Denkfehler, die im Workshop hochkommen.

## 1. Worum geht es?

**Die Analogie:** Wasserdichte Schotten im Schiff. Dringt Wasser in einen Bereich ein, bleiben die anderen trocken — das Schiff sinkt nicht. Genau dieses Prinzip überträgt das Pattern auf Software.

**Das Problem:** Der Booking-Service ruft Hotel, Flight, Car. Hängt **einer** dieser Aufrufe (z. B. Flight wegen Netzwerk-Timeout), saugen sich alle gemeinsamen Ressourcen (Threads, Connections, Goroutines, DB-Pool) am hängenden Service fest. Hotel- und Car-Aufrufe kommen nicht mehr durch — obwohl die Services gesund sind. **Eine kaputte Abhängigkeit reißt den ganzen Service mit.**

**Die Lösung:** Ressourcen werden in **isolierte Pools** unterteilt — pro Downstream-Abhängigkeit. Läuft ein Pool voll, sind die anderen davon unberührt.

---

## 2. Wer wird geschützt?

> Das ist die wichtigste pädagogische Klarstellung — Teilnehmende denken oft, Bulkhead schütze die Downstream-Services.

Bulkhead ist ein **client-seitiges, defensives Pattern**. Geschützt wird:

1. **Der Booking-Service selbst** (der Aufrufer) — damit er nicht alle eigenen Ressourcen an einen einzigen kaputten Downstream verliert.
2. **Die übrigen Downstream-Calls im selben Service** — Hotel und Car laufen weiter, wenn Flight hängt.
3. **Die eigenen Aufrufer** (Frontend, API-Gateway) — die bekommen zumindest noch teilweise funktionierende Antworten.

**Nicht geschützt** werden Hotel, Flight, Car selbst. Wer einen Downstream-Service vor Überlast schützen will, braucht **andere** Patterns, die **dort** sitzen — nicht im Booking-Service.

> 🎯 **Im Workshop klar trennen:**
> - **Resilienz-Patterns im Aufrufer** (Bulkhead, Circuit Breaker, Timeout, Retry) → defensives Verhalten gegenüber unzuverlässigen Abhängigkeiten.
> - **Schutz-Patterns im Aufgerufenen** (Rate Limiting, Backpressure, Load Shedding) → Selbstschutz vor Überlast.

**Die gespiegelte Form:** Pools pro Consumer **im Backend** (Nygard, Azure Architecture Center) sind ebenfalls ein Bulkhead — sie schützen dort das Backend vor einem lauten Client. „Outbound, nicht Inbound" auf der Folie gilt für unseren Pool pro Downstream. Die serverseitige Variante gibt es, wir bauen sie im Workshop nur nicht.

---

## 3. Async ≠ Bulkhead

Verbreiteter Denkfehler: „Ich rufe die drei Services einfach in eigenen Threads / Goroutines async auf — fertig." **Stimmt nicht.**

Async (CompletableFuture, Goroutine, Reactor, Coroutine) sorgt nur dafür, dass der **aufrufende Thread** nicht blockiert. Die Aufrufe konkurrieren aber weiterhin um **gemeinsame Ressourcen**:

- **HTTP-Connection-Pool** — alle Connections wandern in den Flight-Aufruf, Hotel/Car bekommen nichts mehr.
- **Goroutines** — beliebig viele möglich, aber jede hält Speicher, Sockets, evtl. DB-Connections. Bei genug hängenden Goroutines kippt der Prozess (OOM, FD-Limit).
- **Event-Loop** — blockiert nicht, aber Backpressure staut sich in Buffern und Queues.

**Bulkhead verlangt zusätzlich**: harte **Limits pro Abhängigkeit** (Semaphore + dedizierte HTTP-Clients pro Downstream). Async und Bulkhead sind komplementär — beides nötig.

---

## 4. Was passiert beim Limit-Überschreiten?

Wenn der Pool auf 10 begrenzt ist und der 11. Aufruf kommt — drei mögliche Strategien:

| Strategie | Verhalten | Wann sinnvoll |
|---|---|---|
| **Fail-Fast** | Sofort `BulkheadFullException` | **Default in Microservices.** Schnelles Feedback, Fallback/Retry kann greifen. |
| **Bounded Queue** | bis N warten, dann Fehler | Bursty Last mit erträglichem Mittelwert |
| **Wait + Timeout** | kurz warten, dann ablehnen | Mittelweg |

**Empfehlung:** Fail-Fast. Im Fehlerfall füllt sich eine Queue ohnehin sofort, und der Client-Timeout schlägt eh zu — lieber sofort scheitern, damit Circuit Breaker und Fallbacks sauber greifen können.

---

## 5. Pool-Sizing — Little's Law

> Kontraintuitiv: **Ein Service mit höherer normaler Latenz braucht mehr Slots, nicht weniger.** Dimensioniert wird auf die normale Latenz. Läuft der Pool voll, weil das Backend gerade degradiert, ist das das Signal (Reject, Fallback, Alarm), kein Grund, den Pool zu vergrößern.

```
benötigte Slots = Durchsatz × Latenz
L = λ × W
```

**Beispiel bei 50 req/s:**

| Service | normale Latenz | Slots |
|---|---|---|
| Flight | 20 ms | 1 |
| Car | 500 ms | 25 |

**Heuristik:**
1. Aus Monitoring (P99-Latenz, Peak-Throughput) rechnen.
2. 1,5×–2× Puffer für Lastspitzen.
3. Summe aller Pool-Maxima muss unter den Server-Limits bleiben (Connections, FDs, Speicher).
4. Mit Lasttests gegenprüfen.

> ⚠️ **Stolperstein:** Zu kleiner Pool für einen langsamen Service erzeugt selbstgemachte `BulkheadFullException`s — der Downstream ist gesund, nur etwas langsam. Symptom sieht aus wie Ausfall, ist aber Fehlkonfiguration.

**Pool-Größe und Timeout zusammen denken:** Aggressives Timeout + ausreichender Pool. Bei langem Timeout binden hängende Slots den Pool dauerhaft.

---

## 6. Umsetzungsvarianten

| Variante | Beschreibung | Trade-off |
|---|---|---|
| **Thread-Pool** | Eigener Thread-Pool pro Downstream (Hystrix, Resilience4j `ThreadPoolBulkhead`) | Echte Isolation inkl. Timeout, mehr Overhead |
| **Semaphore** | Zähler begrenzt gleichzeitige Aufrufe (Resilience4j `Bulkhead`, Go `chan struct{}`) | Leichtgewichtig, kein eigenes Timeout-Verhalten |
| **Prozess/Pod** | Kritische Funktionen in eigene Services/Pods auslagern | Maximale Isolation, hoher Betriebs-Aufwand |
| **Sidecar/Mesh** | Envoy/Istio begrenzt pro Upstream-Cluster über `max_requests` bzw. `max_connections` | Sprachunabhängig, aber Konfiguration außerhalb des Codes |

**Am Markt:** Resilience4j ist Standard im Java-Umfeld (Semaphore- und ThreadPool-Variante), Polly im .NET-Lager, MicroProfile per `@Bulkhead`-Annotation, go-resiliency für Go.

---

## 7. Pseudocode

Semaphore-basierte Variante — minimaler Kern, illustriert das Pattern:

```kotlin
class Bulkhead(val maxConcurrent: Int) {
    private var inProgress = 0   // atomar / unter Mutex!

    fun call(fn) {
        // Slot reservieren — Fail-Fast wenn voll
        acquireOrThrow()
        try {
            return fn()
        } finally {
            release()           // Slot **immer** freigeben (auch bei Exception)
        }
    }

    fun acquireOrThrow() {
        synchronized(this) {
            if (inProgress >= maxConcurrent) {
                throw BulkheadFullError    // sofort ablehnen
            }
            inProgress++
        }
    }

    fun release() {
        synchronized(this) { inProgress-- }
    }
}
```

> ⚠️ Inkrement und Check müssen **atomar** zusammen passieren (Compare-and-Set, Mutex, Semaphore-Primitive). Sonst rutschen unter Last mehr Aufrufe gleichzeitig durch als erlaubt — die Isolation ist dahin.

> ⚠️ Das `release()` **muss** im `finally` stehen. Vergessen → Slot-Lecks → Pool füllt sich über die Zeit, Bulkhead öffnet nie wieder.

### Ein Bulkhead pro Downstream

```kotlin
val flightBulkhead = Bulkhead(maxConcurrent = 10)
val hotelBulkhead  = Bulkhead(maxConcurrent = 10)
val carBulkhead    = Bulkhead(maxConcurrent = 25)   // höhere normale Latenz → mehr Slots

fun getOffers(req): Offers {
    val flights = flightBulkhead.call { flightClient.search(req) }
    val hotels  = hotelBulkhead.call  { hotelClient.search(req)  }
    val cars    = carBulkhead.call    { carClient.search(req)    }
    return Offers(flights, hotels, cars)
}
```

> 💡 **Wichtig:** Ein **separater** Bulkhead pro Downstream. Ein gemeinsamer Pool für alle Calls würde das ganze Pattern ad absurdum führen — Flight könnte alle Slots belegen und Hotel/Car aushungern.

### Variante mit Wait + Timeout

```kotlin
fun acquireOrThrow(maxWait: Duration) {
    val deadline = now() + maxWait
    synchronized(this) {
        while (inProgress >= maxConcurrent) {
            val remaining = deadline - now()
            if (remaining <= 0) throw BulkheadFullError
            wait(remaining)        // condition variable
        }
        inProgress++
    }
}
```

Nur sinnvoll, wenn kurze Spitzen geglättet werden sollen — bei dauerhafter Überlast bringt das Warten nichts und kostet Latenz.

---

## 8. Zusammenspiel mit anderen Patterns

Reihenfolge auf dem Aufrufpfad:

```
Request ─▶ Bulkhead ─▶ Timeout ─▶ Circuit Breaker ─▶ Retry ─▶ Downstream
```

- **Bulkhead** begrenzt, wie viele Ressourcen ein Downstream maximal binden darf.
- **Timeout** begrenzt, wie lange ein einzelner Aufruf hängen darf.
- **Circuit Breaker** kappt Aufrufe vollständig, wenn der Downstream wiederholt fehlschlägt.
- **Retry** versucht es bei transienten Fehlern erneut (mit Backoff, idempotent!).

**Bulkhead ≠ Circuit Breaker:** Der CB sagt „Backend ist krank, ich versuche es eine Weile nicht" — Reaktion auf die **Fehlerrate**. Der Bulkhead sagt „ich verbrenne maximal N Slots für dieses Backend" — Reaktion auf **Ressourcen-Druck**. Das Szenario, das nur der Bulkhead löst: Das Backend antwortet **langsam, aber fehlerfrei** — der CB sieht keine Fehler und bleibt geschlossen (Closed).

---

## 9. Diskussionsfragen

1. Wie groß dimensioniert ihr eure Pools? Habt ihr eine Heuristik oder messt ihr?
2. Fail-Fast vs. Queue — wo zieht ihr die Grenze?
3. Thread-Pool-Bulkhead vs. Semaphore — welche Variante passt zu eurem Stack?
4. Wann lohnt sich Bulkhead auf Pod-Ebene (eigene Services) statt nur In-Process?
5. Was sagt euer Monitoring, wenn ein Bulkhead voll läuft? Alerting? Dashboard?
