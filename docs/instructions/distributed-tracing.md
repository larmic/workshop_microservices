# Distributed Tracing — Workshop-Notizen

> Trainer-Notizen für Story 7. Kurz gehalten, fokussiert auf das Grund­konzept (Trace-ID + Propagation) und den Marktüberblick, ohne den Workshop in eine OpenTelemetry-Schulung zu kippen.

## 1. Worum geht es?

**Die Analogie:** Sendungsverfolgung über mehrere Lieferdienste hinweg. Ein Paket wird vom Online-Shop an DHL übergeben, DHL gibt es an einen Partner in Polen weiter, der Partner an einen lokalen Kurier. Ohne eine **gemeinsame Sendungsnummer**, die jeder Beteiligte mit­dokumentiert, weiß am Ende niemand, wo das Paket steckt. Genau dieses Prinzip löst Distributed Tracing für verteilte Software-Aufrufe.

**Das Problem:** Im Microservice-Setup ist eine fachliche Operation („buche mir Flug + Hotel + Mietwagen") über mehrere Services verteilt. Jeder Service loggt seine Sicht — aber niemand kann die Logzeilen zu **einer** Operation zusammenführen, außer mit Timestamps und Glück. Sobald Stories 3 (Circuit Breaker), 4 (Bulkhead) oder gar 5/6 (Saga) ins Spiel kommen, ist die Suche nach „warum ist Buchung X gescheitert?" reine Detektivarbeit.

**Die Lösung:** Eine **Trace-ID** wird beim Eingang in das System einmal vergeben (oder vom Aufrufer übernommen) und **bei jedem ausgehenden Aufruf mitgeschickt**. Jeder Service schreibt sie in jede Logzeile. Damit reicht `grep <trace-id>` über die Compose-Logs, um den kompletten Vorgang zu rekonstruieren.

> 🎯 *Demo-Einstieg:* Im Story-6-Setup eine Buchung absenden und dann versuchen, **eine konkrete Anfrage** in den Logs aller vier Services nachzuverfolgen. Schnell wird klar: das ist nicht trivial. Dann in Story 7 dasselbe — und `grep <trace-id>` zeigt alles in einem Block.

---

## 2. Die zentralen Konzepte

| Begriff | Bedeutung |
|---|---|
| **Trace** | Ein kompletter Geschäftsvorgang über alle Services hinweg. Identifiziert durch die **Trace-ID** (16 Byte, üblicherweise 32 Hex-Zeichen). |
| **Span** | Ein einzelner Arbeitsabschnitt innerhalb des Traces — typischerweise ein HTTP-Call, ein DB-Query, eine Funktionsausführung. Hat eine eigene **Span-ID** (8 Byte). |
| **Parent-Span-ID** | Verweis auf die Span, in deren Kontext diese Span entstanden ist. Ergibt einen Span-Baum. |
| **Sampling-Flag** | Bit, das entscheidet, ob dieser Trace ans Backend gesendet wird oder nicht (für Cost-Control bei Volumen-Traffic). |

Ein Trace sieht in einem Tracing-Tool wie Jaeger so aus:

```
Trace abc123…  (1.2 s)
├─ POST /booking/bookings          (booking-service, 1.2 s)
│   ├─ POST /bookings              (flight-service,  120 ms)
│   ├─ POST /bookings              (hotel-service,   180 ms) ✗ 500
│   ├─ POST /bookings              (car-service,     900 ms)
│   └─ DELETE /bookings/{id}       (flight-service,   80 ms)  ← Compensation
```

Im Workshop bauen wir bewusst **nur Trace-IDs**, keinen vollwertigen Span-Baum mit Parent-Child-Verknüpfung — das wäre für 60 Minuten zu viel und ist Tool-spezifisch. Span-IDs werden pro Hop generiert, dienen aber im Workshop-Code nur dazu, das `traceparent`-Format gültig zu halten.

---

## 3. W3C Trace Context — der Standard

Bis ca. 2020 hatte jedes Tracing-Ecosystem sein eigenes Header-Format (Zipkin: `X-B3-*`, Jaeger: `uber-trace-id`, Datadog: `x-datadog-*` …). Das war Inkompatibilitäts­hölle. Seitdem ist **W3C Trace Context** der gemeinsame Standard, den alle großen Vendor unterstützen.

Der zentrale Header heißt `traceparent` und hat genau 55 Zeichen:

```
traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
             │  └─── Trace-ID (16 Bytes = 32 hex) ─┘ └─Span-ID─┘ └Flags
             │
             └── Version (immer 00 derzeit)
```

| Feld | Länge | Bedeutung |
|---|---|---|
| `version` | 2 hex | Aktuell immer `00`. Höhere Versionen müssen abwärtskompatibel ignoriert werden. |
| `trace-id` | 32 hex | 16 Byte zufällig. **Eindeutig für den ganzen Trace** — bleibt über alle Hops gleich. Darf nicht nur Nullen sein. |
| `parent-id` (Span-ID) | 16 hex | 8 Byte zufällig. **Pro Hop neu**. Der Empfänger benutzt diesen Wert als Parent für seine eigene Span. |
| `flags` | 2 hex | Bitmaske. Bit 0 = „sampled" (Trace soll exportiert werden). |

Optional gibt es noch `tracestate` für Vendor-spezifische Daten — im Workshop nicht relevant.

> ⚠️ **Stolperstein:** Die Trace-ID `00000000000000000000000000000000` und Span-ID `0000000000000000` sind ungültig. Wer Trace-Kontexte selbst generiert, muss zumindest sicherstellen, dass `crypto/rand` & Co. genutzt werden — nicht `math/rand` mit Default-Seed.

---

## 4. Pseudocode

Ungefähr das, was die Reference-Implementierung in `services/shared/tracing/tracing.go` macht:

```kotlin
// Beim eingehenden Request — Server-Middleware
fun middleware(req, next) {
    val traceContext = parse(req.header("traceparent")) ?: generateNew()
    val ctx = req.context.with(traceContext)
    next(req.withContext(ctx))
}

// Beim ausgehenden Request — Client-Inject
fun inject(ctx, outboundReq) {
    val incoming = ctx.get(TraceContext)
    val newSpan = randomSpanId()                            // pro Hop neu
    val hopContext = incoming.copy(spanId = newSpan)
    outboundReq.setHeader("traceparent", hopContext.toHeader())
}

// In jeder Logzeile
fun log(ctx, msg) {
    val tc = ctx.get(TraceContext)
    logger.info(msg, "trace_id" to tc.traceId, "span_id" to tc.spanId)
}

// Parser — strikte Format-Validierung
fun parse(header: String?): TraceContext? {
    if (header == null) return null
    val parts = header.split("-")
    if (parts.size != 4) return null
    if (parts[0] != "00") return null                        // nur v0 verstehen
    if (parts[1].length != 32 || !parts[1].isHex) return null
    if (parts[2].length != 16 || !parts[2].isHex) return null
    if (parts[3].length != 2  || !parts[3].isHex) return null
    if (parts[1].all { it == '0' } || parts[2].all { it == '0' }) return null
    return TraceContext(parts[1], parts[2], parts[3])
}
```

> 🎯 *Folie:* Den Pseudocode kurz zeigen, dann auf den echten Go-Code in `services/shared/tracing/tracing.go` verweisen — exakt diese Logik in ca. 100 Zeilen.

---

## 5. Wer erstellt was? — Durchlauf in unserem Stack

Im Workshop-Setup ist die Aufgabenverteilung bewusst asymmetrisch: **Booking** ist der einzige Entry-Point, **Flight / Hotel / Car** sind passive Trace-Empfänger. Konkret:

```
Client ─POST /booking/bookings─►  Booking
                                  ├── erstellt Trace-ID T  (falls keiner reinkommt)
                                  ├── erstellt Span-ID S1  (eigener Server-Span)
                                  │
                                  ├─traceparent: 00-T-S2-01─►  Flight   (loggt mit T, S2)
                                  ├─traceparent: 00-T-S3-01─►  Hotel    (loggt mit T, S3)
                                  └─traceparent: 00-T-S4-01─►  Car      (loggt mit T, S4)

                                  ├── publish event { traceparent: 00-T-S5-01 }
                                                                  │
                                              Flight ◄────────────┘   (loggt mit T, S5)
```

Eine **Trace-ID T**, fünf verschiedene **Span-IDs (S1–S5)** — alle von Booking erstellt. Flight/Hotel/Car erzeugen selbst gar nichts.

### Was Booking erstellt

1. **Trace-ID** — einmalig beim eingehenden Request, falls kein `traceparent` mitkommt (Browser/curl-Fall). Dafür ist die `Middleware`-Variante zuständig (siehe Abschnitt 4).
2. **Span-ID Nr. 1** — für den eigenen Server-Span (die eingehende `POST /booking/bookings`-Verarbeitung).
3. **Span-ID Nr. 2–4** — pro Outbound-Call zu Flight, Hotel, Car (Client-Inject). Trace-ID bleibt konstant, Span-ID ist pro Hop neu.
4. **`traceparent` als Event-Property** — beim Publish eines `CompensationRequested`-Events in Story 6/7.

### Was Flight / Hotel / Car erstellen — nichts

Sie nutzen die `Propagate`-Middleware:

- Lesen den eintreffenden `traceparent`.
- Loggen mit der mitgelieferten Trace-ID + Span-ID.
- **Fehlt der Header → keine `trace_id` im Log**, *kein* Fallback auf eine selbst generierte ID. Absicht: nur der Entry-Point sät Trace-IDs. Der Kontrast macht in den Stories 1–6 (ohne Tracing) sichtbar, dass die Backend-Logs ohne `trace_id` bleiben — der rote Faden entsteht erst, sobald Booking ihn in Story 7 zieht.
- Bei Story 6/7-Compensation-Events: parsen `traceparent` aus der Event-Property, legen den Kontext in ihre Worker-Goroutine, loggen damit weiter.

### Was wäre, wenn ein Downstream-Service selbst weitere Services aufruft?

Aktuell haben Flight/Hotel/Car keine eigenen Downstream-Abhängigkeiten. Sobald sie das hätten (z. B. Flight ruft einen externen Reservierungs-Provider, Hotel ein Loyalty-Backend), würde sich am Pattern **nichts** ändern:

- Die `Propagate`-Middleware hat den eingehenden Trace-Kontext bereits in den Request-Context gelegt.
- Flight nutzt **denselben Client-Inject** wie Booking — generiert pro Hop eine **neue Span-ID** (`S_new`), behält die **Trace-ID T** bei.
- Der externe Provider sieht `traceparent: 00-T-S_new-01`.
- Falls der externe Provider Trace-Context versteht, loggt er ebenfalls mit T. Andernfalls endet der rote Faden dort — aber alles bis dahin bleibt korreliert.

```
Client ─POST /booking/bookings─►  Booking ─traceparent: 00-T-S2-01─►  Flight
                                                                      ├── Propagate übernimmt: ctx = {T, S2}
                                                                      ├── Client-Inject: neue Span-ID S_new
                                                                      └─traceparent: 00-T-S_new-01─►  external Provider
```

> ⚠️ **Wichtig:** Flight darf in diesem Szenario **nicht** versehentlich eine neue Trace-ID erzeugen — sonst zerfällt der Trace genau dort. Die strikte Trennung von `Middleware` (erzeugt, falls keiner kommt — nur am Entry-Point!) und `Propagate` (übernimmt nur — überall sonst) bleibt auch im erweiterten Setup gültig.

**Faustregel:** Jeder Service, der Aufrufe annimmt und weiterleitet, braucht eine Server-Middleware vom Typ *Propagate*, aber niemals *Middleware*. Trace-Initiierung gehört exklusiv an die Systemgrenze (API-Gateway, Public-Service) — nicht in den Downstream-Hop.

---

## 6. Logging-Korrelation — warum strukturiert?

Mit `log.Printf("booking %s failed: %v", id, err)` muss man `trace_id=...` manuell in jeden Format-String einsetzen. Das wird nach drei Wochen vergessen, und ab dann fehlt die ID in Hälfte der Logzeilen.

Mit strukturiertem Logger (`slog` in Go, `structlog` in Python, `logback` JSON-Layout in Java, `pino` in Node) hängt man die ID einmal an den Logger im Request-Kontext — und sie taucht in **jeder** Zeile auf, ohne Format-Strings:

```go
logger := tracing.Logger(ctx)              // hat trace_id schon dran
logger.Info("forward step done", "step", "flight", "duration_ms", 120)
```

Output (JSON):

```json
{"time":"2026-05-12T10:00:01Z","level":"INFO","msg":"forward step done",
 "trace_id":"0af7651916cd43dd8448eb211c80319c","span_id":"b7ad6b7169203331",
 "step":"flight","duration_ms":120}
```

> 💡 Sobald die Logs JSON sind, gibt `docker compose logs -f | jq -c 'select(.trace_id=="abc…")'` einen sauber gefilterten Live-Stream.

---

## 7. Selbst bauen oder Library?

| Aspekt | Selbst | Library (OpenTelemetry) |
|---|---|---|
| Lerneffekt | ⭐⭐⭐⭐⭐ Format und Propagation werden verstanden | ⭐⭐ Auto-Instrumentation, magisch |
| Korrektheit | Header-Parsing strikt machen, Race-Bugs in Context-Propagation möglich | Kampferprobt |
| Span-Baum (Parent/Child) | Müsste man dazubauen | Inklusive |
| Sampling, Batching, Export | Müsste man bauen | Inklusive |
| Backend-Export (Jaeger, Tempo, Datadog…) | Müsste man pro Backend bauen | Inklusive (OTLP) |
| Zeilen Code | ~100 (Header + Middleware + Logger) | 0 + ein paar Zeilen Konfiguration |
| Tool-/Sprach-Bindung | Keine | Pro Sprache eigenes SDK |

**Empfehlung für den Workshop:** Selbst bauen. Die Pflicht-Story-Implementierung ist in jeder Sprache identisch (HTTP-Header lesen, generieren, weiterreichen, in Logs schreiben). OpenTelemetry hingegen ist pro Sprache deutlich unter­schiedlich — Go nutzt Auto-Instrumentation für `net/http`, Java setzt auf Agent + Annotations, Node ist wieder anders.

**Empfehlung für die Praxis:** OpenTelemetry. Die Auto-Instrumentation, das Sampling, das Span-Lifecycle-Management und der OTLP-Export an beliebige Backends sind Arbeit, die niemand neu bauen will.

---

## 8. Marktüberblick

### Standards (Wire-Formate)

| Standard | Status | Wo verbreitet |
|---|---|---|
| **W3C Trace Context** (`traceparent` / `tracestate`) | Aktuell, Default | Überall (OTel, neue Datadog, neue New Relic, etc.) |
| **B3** (`X-B3-TraceId` / `X-B3-SpanId` / …) | Legacy, weit verbreitet | Zipkin, ältere Spring-Cloud-Setups |
| **Jaeger native** (`uber-trace-id`) | Legacy | Ältere Jaeger-only-Setups |
| **Datadog** (`x-datadog-trace-id` / `x-datadog-parent-id`) | Vendor-spezifisch | Ältere Datadog-Agents, mittlerweile auch W3C-fähig |

Moderne SDKs unterstützen meist mehrere Formate parallel (Multi-Propagator), damit man schrittweise migrieren kann.

### Instrumentation-Libraries

| Library / Framework | Beschreibung |
|---|---|
| **OpenTelemetry** | De-facto-Standard. SDK + Auto-Instrumentation für die meisten gängigen Frameworks. Vendor-agnostisch über OTLP-Export. |
| OpenTracing (deprecated) | Vorgänger, in OpenTelemetry aufgegangen. |
| OpenCensus (deprecated) | Googles früherer Ansatz, ebenfalls in OpenTelemetry aufgegangen. |
| Vendor-SDKs (Datadog APM, New Relic Agent, Dynatrace OneAgent, …) | Eigene SDKs der APM-Hersteller. Oft komfortabler, dafür Vendor-Lock-in. |
| Spring Cloud Sleuth / Micrometer Tracing | Java-spezifisch. Sleuth ist abgekündigt, Micrometer Tracing ist der Nachfolger. |

### Backends / Tools

| Tool | Typ | Stärken | Schwächen |
|---|---|---|---|
| **Jaeger** | Open Source | Klassiker, einfach lokal zu betreiben, eigene UI, CNCF-Projekt | Storage skaliert nicht trivial |
| **Grafana Tempo** | Open Source | Sehr günstig (Objekt-Storage als Backend), integriert sich in Grafana/Loki/Mimir | UI nur über Grafana |
| **Zipkin** | Open Source | Älter, einfach, schlank | Funktionsumfang begrenzt im Vergleich |
| **Datadog APM** | SaaS | Sehr ausgereifte UI, Korrelation Logs+Traces+Metrics | Kosten, Vendor-Lock-in |
| **Honeycomb** | SaaS | Fokus auf Observability + high-cardinality Queries | Pricing-Modell |
| **Dynatrace** | SaaS | OneAgent + AI-basiertes Root-Cause | Schwergewichtig, teuer |
| **New Relic** | SaaS | Etabliert, breit | Pricing |
| **Elastic APM** | Self-Hosted / SaaS | Integration in ELK-Stack | Komplex zu betreiben |
| **AWS X-Ray** | SaaS (AWS) | Nahtlos in AWS-Services | Nur AWS-Welt, ältere Architektur |

Für den Workshop (Bonus): **Jaeger All-in-One** als einzelner Container ist die kürzeste Strecke zur ersten Trace-UI.

---

## 9. Sampling — die Cost-Frage

In Produktion sind Traces teuer (Storage, Netzwerk, Backend-Kosten). Niemand tracet jede Anfrage. Übliche Strategien:

| Strategie | Wann entscheiden | Vorteil | Nachteil |
|---|---|---|---|
| **Head-based** (z.B. 1 %) | Beim Eingang | Kein Buffering nötig | Fehler-Traces gehen oft verloren |
| **Tail-based** | Nach Trace-Ende | Fehler/Slow-Traces immer behalten | Backend muss bufferen, mehr Aufwand |
| **Adaptive** | Dynamisch | Reagiert auf Last | Komplex zu betreiben |
| **Always-on für Errors / Slow** | Per Heuristik | Beste Sichtbarkeit der Probleme | Erfordert Tail-based-Stack |

> 💬 *Diskussionsfrage:* Bei wem im Workshop läuft Tracing produktiv? Wie hoch ist eure Sampling-Rate, und wie geht ihr mit den vielen verlorenen erfolgreichen Traces um?

---

## 10. Tracing über Async-Grenzen

Sobald die Kommunikation nicht mehr synchron HTTP ist (Story 6: Compensation-Events), reicht der HTTP-Header nicht mehr. Der Trace-Kontext muss als **Property auf dem Event** mitwandern:

```json
{
  "eventId": "...",
  "sagaId": "...",
  "traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
  "payload": { ... }
}
```

Der Konsument liest die Property, legt den Trace-Kontext in den Goroutine-/Thread-Kontext und loggt damit weiter. Andernfalls bricht der Trace genau an der Bus-Grenze ab — die spätere Stornierungs-Logzeile hat keine Korrelation mehr zur ursprünglichen Buchung.

> ⚠️ **Häufiger Stolperstein:** Bei Kafka/RabbitMQ/SNS-SQS gibt es Message-Header bzw. Properties. **Nicht** den Trace-Kontext in den Payload pressen, sondern in die Message-Header — sonst muss jeder Konsument die Payload-Schema kennen, um den Trace zu propagieren.

---

## 11. Diskussionsfragen

1. **Wer schreibt die Trace-ID heute schon in eure Logs?** Häufige Lücke: das Gateway (Traefik/Nginx) generiert eine, aber die Backend-Services schreiben sie nicht in die App-Logs.
2. **Wie weit reicht euer Trace?** Bis zum Frontend? In den Browser? Über CDN-Hops? Per Mobile-App-SDK?
3. **Sampling: head- oder tail-based?** Und wer entscheidet bei euch die Sampling-Rate?
4. **Korrelation Logs ↔ Traces ↔ Metriken** — habt ihr alle drei im selben Tool? Wechselt der Operator pro Frage das Tool?
5. **Async-Tracing:** Bei wem im Stack laufen Events/Queues? Wandert die Trace-ID dort schon mit?
6. **Was kostet euch Tracing?** Storage, Bandbreite, Backend-Lizenz — habt ihr eine Größenordnung?

---

## 12. Wo der Code liegt (Reference-Implementierung)

| Was | Pfad |
|---|---|
| Tracing-Bibliothek (Parse, Generate, Middleware, Inject, Logger) | `services/shared/tracing/tracing.go` |
| Tests | `services/shared/tracing/tracing_test.go` |
| Server-Middleware verdrahten | `services/booking/story7/main.go`, `services/flight/main.go`, `services/hotel/main.go`, `services/car/main.go` |
| Outbound-Inject + strukturierte Logs | `services/booking/story7/handler/booking.go` |
| Trace-Kontext in Compensation-Events | `services/booking/story7/saga/*.go` und Konsumenten `services/{flight,hotel,car}/handler/compensation.go` |
| Backend-Handler nutzen Trace-Logger | `services/{flight,hotel,car}/handler/*.go` |
