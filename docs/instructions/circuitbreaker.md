# Circuit Breaker — Workshop-Notizen

> Trainer-Notizen für Story 3. Reihenfolge folgt einer typischen Slide-Sequenz; jeder Abschnitt hat einen kleinen Hinweis, was an der Tafel/auf der Folie passieren sollte.

## 1. Worum geht es?

**Die Analogie:** Der Sicherungsautomat im Stromkasten. Wenn zu viel Strom fließt, springt er raus und unterbricht den Stromkreis — bevor das Haus brennt. Nach kurzer Zeit kannst du ihn manuell wieder reinschalten. Der Software-Circuit-Breaker macht dasselbe automatisiert mit Service-Calls.

**Das Problem:** Service A ruft Service B. B ist überlastet/kaputt. A wartet jeweils 30 Sekunden auf Timeout. Die Anfragen stauen sich, A's Threads/Goroutines/Connections gehen aus, A wird ebenfalls langsam, dessen Aufrufer kommen ins Schwitzen — der klassische Cascading-Failure.

**Die Lösung:** A merkt sich, dass B kaputt ist, und kürzt zukünftige Calls sofort ab. So bleibt A schnell und gesund, auch wenn B unten ist. Nach einer Wartezeit testet A, ob B wieder da ist.

> 🎯 *Demo-Einstieg:* Im Workshop-Stack Flight auf "Fehler" stellen, ohne CB (Story 2) auf `/booking/offers` curlen — jeder Call dauert 3s. Mit CB (Story 3) das gleiche — nach den ersten 5 Fehlern sind die Calls instant.

---

## 2. Die drei Zustände

```
                  Erfolg
                ┌─────────┐
                │         │
                ▼         │
   ┌────────────────┐     │
   │     CLOSED     │─────┘
   │  alles normal  │
   └────────┬───────┘
            │
            │ N aufeinanderfolgende Fehler
            ▼
   ┌────────────────┐
   │      OPEN      │◄───────────┐
   │ Calls werden   │            │
   │ sofort         │            │ Probe schlägt fehl
   │ short-circuited│            │
   └────────┬───────┘            │
            │                    │
            │ Timeout abgelaufen │
            ▼                    │
   ┌────────────────┐            │
   │   HALF_OPEN    │────────────┘
   │ ein Probe-Call │
   │   wird zuge-   │
   │     lassen     │────────────┐
   └────────────────┘            │
                                 │ Probe erfolgreich
                                 ▼
                              CLOSED
```

### CLOSED — der Normalbetrieb
- Calls gehen direkt durch.
- Jeder Fehler erhöht einen Zähler.
- Bei `N` aufeinanderfolgenden Fehlern (oder bei einer Fehler-**Quote** über X%) → Übergang nach OPEN.
- Erfolg setzt den Zähler zurück.

### OPEN — die Schutzschaltung
- Calls werden gar nicht erst zum Backend geschickt — sofort `CircuitOpenError`.
- Das Caller-System bleibt schnell und reaktiv, statt in Timeouts zu hängen.
- Nach einer konfigurierten Wartezeit (Workshop: 30s) wird der nächste Call zum Test zugelassen → HALF_OPEN.

### HALF_OPEN — der vorsichtige Test
- **Genau ein** Probe-Call wird zum Backend durchgelassen. Andere parallele Anfragen bekommen weiterhin den `CircuitOpenError`.
- Erfolg → Backend gilt als gesund → CLOSED, Fehlerzähler auf 0.
- Fehler → Backend immer noch kaputt → zurück nach OPEN, Wartezeit beginnt von vorne.

> ⚠️ **Stolperstein für Diskussion:** Ohne den Single-Probe-Schutz würden in HALF_OPEN alle wartenden Anfragen gleichzeitig durchgelassen — ein "Probe-Storm" der das gerade erholende Backend wieder umhaut.

---

## 3. Was muss konfiguriert werden?

| Parameter | Workshop-Wert | Bedeutung |
|---|---|---|
| `failureThreshold` | 5 | Zähl-basiert: nach N Fehlern öffnet der CB. Alternativ: `failureRate` (%) über ein Sliding Window — robuster gegen Mischverkehr. |
| `openTimeout` / `waitDurationInOpenState` | 30s | Wie lange im OPEN-Zustand verharren, bevor ein Probe-Call versucht wird. Zu kurz = unruhig, zu lang = Recovery dauert. |
| `callTimeout` | 3s | Wann ein einzelner Aufruf als "fehlgeschlagen" gilt, weil er hängt. **Wichtig**: Der CB hilft nur, wenn das Aufrufen selbst nicht stundenlang blockiert. |
| Sliding Window | (nicht verwendet) | Bei rate-basierter Variante: Anzahl letzter Calls oder Zeitfenster, über das die Quote berechnet wird. |
| Minimum Calls | (nicht verwendet) | Bei rate-basierter Variante: erst auswerten, wenn z.B. mind. 10 Calls vorlagen. Verhindert dass 1/1 Fehler den CB öffnet. |

> 🎯 *Folie:* Diese Tabelle zeigen und betonen: Es gibt keine universellen Werte. Sie hängen ab von:
> - SLA des Backends ("normal" sind 50ms? 2s?)
> - Geschäftliche Toleranz (lieber kurz blocken oder lieber lange Timeouts hinnehmen?)
> - Traffic-Volumen (bei 1 Req/s ist count-basiert OK, bei 10k Req/s eher rate-basiert)

### Was zählt als "Fehler"?

Wichtige Diskussion, oft unterschätzt:

| Ergebnis | Failure? | Warum |
|---|---|---|
| HTTP 5xx | ✅ ja | Backend selbst ist kaputt |
| Timeout | ✅ ja | Backend antwortet nicht in akzeptabler Zeit |
| Connection refused / DNS-Fehler | ✅ ja | Backend nicht erreichbar |
| HTTP 4xx | ❌ nein (üblicherweise) | Der Aufrufer hat Mist gemacht — Backend ist gesund |
| Cancelled durch Client | ❌ nein | Kein Backend-Problem |

> 💡 Resilience4j hat dafür einen `recordExceptions` / `ignoreExceptions`-Predicate. Spring Cloud Circuit Breaker analog. Im Workshop-Code: bewusst simpel gehalten, alles ≥ 400 zählt — Diskussionspunkt für die Teilnehmer: "Wo seht ihr das anders?"

---

## 4. Pseudocode

```kotlin
fun call(fn) {
    // OPEN → HALF_OPEN, sobald die Wartezeit abgelaufen ist
    if (state == OPEN && now() - openedAt > openTimeout) {
        state = HALF_OPEN
    }

    when (state) {
        OPEN      -> throw CircuitOpenError
        HALF_OPEN -> if (probeInFlight) throw CircuitOpenError
                     else probeInFlight = true     // atomar setzen!
        CLOSED    -> {}                            // einfach durchwinken
    }

    try {
        val result = fn()
        onSuccess()
        return result
    } catch (e) {
        onFailure()
        throw e
    }
}

fun onFailure() {
    failures++
    when {
        state == HALF_OPEN                         -> open()
        state == CLOSED && failures >= threshold   -> open()
    }
}

fun onSuccess() {
    if (state == HALF_OPEN) state = CLOSED
    probeInFlight = false
    failures = 0
}

fun open() {
    state = OPEN
    openedAt = now()
    probeInFlight = false
}
```

> ⚠️ Der `probeInFlight`-Flip muss **atomar** passieren (Compare-and-Set / Mutex). Sonst rutschen unter Last mehrere Aufrufe gleichzeitig als „Probe" durch — der Probe-Storm ist genau das, was HALF_OPEN verhindern soll.

> 🎯 *Folie:* Pseudocode kurz zeigen, dann auf den echten Code in `services/booking/story3/circuitbreaker/circuitbreaker.go` verweisen — exakt diese Logik in ~80 Zeilen Go.

---

## 5. Wo CB einsetzen — und wo nicht

Das ist der **wichtigste pädagogische Übergang** der Story 3.

### Lese-Operationen (GET) — Idealfall

`GET /booking/offers` ist idempotent und ohne Seiteneffekte. Wenn der Flight-Service down ist:
- Mit CB: Antwort enthält Hotels + Cars, `flights: []`, Header `X-Fallback: flight`. Der Nutzer sieht *etwas* — der Browse-Screen funktioniert.
- Ohne CB: ganzer Request scheitert mit 500. Schwarzes Loch.

**Graceful Degradation passt perfekt zu GET.** Workshop-Teilnehmer sollen das selbst sehen: Backend deaktivieren, in Story 2 vs. Story 3 testen.

### Schreib-Operationen (POST) — knifflig

Bei `POST /booking/bookings` will man Flug + Hotel + Auto buchen. Wenn der Flight-Service nach 2 erfolgreichen Calls (Hotel + Auto bereits gebucht) ausfällt — was tun?

Drei Stufen:
1. **Naive Graceful Degradation** (was wir initial hatten): `flight: null`, der Kunde wundert sich warum er kein Flug-Ticket hat aber ein Hotelzimmer und einen Mietwagen. **Schlechte UX.**
2. **Fail-Fast** (was Story 3 jetzt macht): wenn ein CB OPEN ist, sofort `503 Service Unavailable`. Wenn mid-call ein Service kippt, abbrechen. Bessere UX, aber: was ist mit den schon erfolgten Sub-Buchungen? Die bleiben **als Leiche im System**.
3. **Saga-Pattern** (Story 5): explizite Compensation. Schlägt der Auto-Call fehl, wird Hotel und Flug aktiv per `DELETE` zurückgerollt. Echte Atomicity über Service-Grenzen.

> 🎯 *Diese drei Stufen sind die Brücke zu Story 5*. Im Workshop laut sagen: "Der Circuit Breaker schützt eure Services. Atomicity über mehrere Services ist ein **anderes** Problem — und die Lösung heißt Saga."

> ⚠️ **Häufiger Denkfehler:** "Wir machen einfach Retry" reicht nicht. Retry ohne Idempotency-Key kann zu Doppelbuchungen führen. Retry ohne CB lässt das geschwächte Backend nicht atmen. Beides zusammen ist kombinierbar — aber Saga ist die saubere Lösung für alles-oder-nichts.

---

## 6. Selbst implementieren oder Library?

| Aspekt | Selbst | Library |
|---|---|---|
| Lerneffekt | ⭐⭐⭐⭐⭐ State-Machine wird verstanden | ⭐⭐ Bleibt Black Box |
| Korrektheit | Subtile Race-Bugs lauern (Probe-Storm, atomar/Mutex, Goroutine-Lecks) | Kampf­erprobt, von vielen Augen geprüft |
| Features | Was du baust | Sliding Window, Bulkhead-Integration, Metrics, Health-Indicator, Annotations… |
| Zeilen Code | ~80 (Go) | 0 + ein paar Zeilen Konfiguration |

**Empfehlung für die Praxis:** Library nehmen.

**Empfehlung für den Workshop:** Einmal selber bauen (oder den Reference-Code `circuitbreaker/circuitbreaker.go` lesen), damit man weiß, was die Library macht. Danach bei den eigenen Projekten auf Bewährtes setzen.

### Standardlibraries je Stack

| Stack | Library |
|---|---|
| Java (framework-frei) | **Resilience4j** |
| Spring | Spring Cloud Circuit Breaker (Wrapper, intern Resilience4j) |
| Quarkus / MicroProfile | MicroProfile Fault Tolerance (Annotation-basiert) |
| .NET | **Polly** |
| Go | gobreaker, sony/gobreaker, hashicorp/go-circuitbreaker |
| Node.js | opossum |

> 💬 *Diskussionsfrage an die Teilnehmer:* "Welche von denen nutzt ihr in eurem Tagesgeschäft? Was hat euch überrascht?"

---

## 7. Ausblick: CB als Sidecar / Service Mesh

In den vorherigen Abschnitten lebt der CB *im* Service. Alternative: CB ins Service-Mesh verlagern (Envoy / Istio / Linkerd / Consul Connect).

**Vorteile:**
- Sprach-/Framework-unabhängig: Java, Go, Python, Rust — alle bekommen die gleiche Resilience.
- Zentrale Konfiguration via Mesh-Control-Plane.
- Beobachtbarkeit out-of-the-box (Tracing, Metriken).
- Anwendungscode bleibt sauber.

**Nachteile:**
- Komplexität explodiert: zusätzliche Container, Control-Plane, Lernkurve.
- Fallback-Logik (was passiert *wenn* der CB feuert?) muss trotzdem im Anwendungscode stehen — das Mesh kann nicht entscheiden, dass `flights:[]` sinnvoll ist.
- Latenz pro Hop steigt minimal (Sidecar-Proxy).
- Debugging ist anders (Header-Propagation, mTLS, etc.).

> 🔭 *Im Workshop nur kurz erwähnen.* Service-Mesh-Tiefe gehört in einen eigenen Termin. Aber wichtig: die Teilnehmer sollen wissen, dass der CB nicht zwingend Anwendungscode sein muss.

---

## 8. Diskussionsfragen für den Workshop

Zum Abschluss, wenn Zeit bleibt:

1. **Wie testet ihr den CB?** Tipp: Chaos Engineering — Toxiproxy, Gremlin, oder unser Dashboard mit Chaos-Buttons. Realistische Szenarien (5xx, Latenz, Drops) sind Pflicht.
2. **Wie bekommt der Operator mit, dass ein CB OPEN ist?** Metriken (Prometheus), Alerts. Im Reference-Code: Logs + Dashboard. In Produktion: pro CB einen Counter `circuit_state{name=...,state=OPEN}` exportieren.
3. **Was, wenn euer eigener Service der "schwache" ist?** CB ist *outbound*. Inbound ist Bulkhead/Rate-Limiting (Story 4) oder Backpressure.
4. **Granularität: pro Service oder pro Endpoint?** Im Reference-Code pro Service. Aber: was, wenn `/flights` gut antwortet aber `/bookings` kaputt ist? Diskutieren.
5. **Per-Instance vs. per-Service:** Mit dem Dashboard kann man genau einer von drei Replicas auf "Fehler" stellen. Diskussionspunkt: Random-Load-Balancing trifft die kaputte Instanz nur in 1/3 der Fälle — der CB öffnet evtl. gar nicht. Ist das ein Bug oder ein Feature?

---

## 9. Wo der Code liegt (Reference-Implementierung)

| Was | Pfad |
|---|---|
| CB-State-Machine | `services/booking/story3/circuitbreaker/circuitbreaker.go` |
| Drei CBs verdrahten | `services/booking/story3/main.go` |
| GET (Graceful Degradation) | `services/booking/story3/handler/booking.go` (`BookingOffersHandler`, `fetchOffersWithCB`) |
| POST (Fail-Fast + Saga-Hook) | `services/booking/story3/handler/booking.go` (`CreateBookingHandler`) |
| Admin-Endpunkte | `services/booking/story3/handler/admin.go` (`/admin/circuit-state`, `/admin/circuit-events`) |
| Chaos-Steuerung Backend | `services/shared/chaos/chaos.go` |
| Dashboard-Visualisierung | `services/dashboard/static/index.html` + `services/dashboard/handler/{chaos,circuit}.go` |
