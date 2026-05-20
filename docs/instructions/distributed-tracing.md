# Distributed Tracing — Workshop-Notizen

> Trainer-Notizen für Story 7. Kurz gehalten: Konzepte vorab, alles Weitere am konkreten Beispiel.

## 1. Worum geht es?

**Analogie:** Sendungsverfolgung. Ein Paket wandert vom Online-Shop über DHL, einen Partner in Polen, einen lokalen Kurier — ohne **gemeinsame Sendungsnummer** weiß am Ende niemand, wo es steckt.

**Problem:** Eine fachliche Operation („buche Flug + Hotel + Mietwagen") läuft über vier Services. Jeder loggt seine Sicht — aber niemand kann die Zeilen zu **einer** Operation zusammenführen. Spätestens mit Saga (Stories 5/6) wird Debugging zur Detektivarbeit.

**Lösung:** Eine **Trace-ID**, einmal pro Request vergeben, durch alle Hops weitergereicht und in jede Logzeile geschrieben. `grep <trace-id>` über die Compose-Logs zeigt den kompletten Vorgang.

---

## 2. Konzepte

| Begriff | Bedeutung |
|---|---|
| **Trace** | Ein kompletter Geschäftsvorgang über alle Services hinweg. Eine **Trace-ID** (16 Byte / 32 hex), die für alle Hops identisch bleibt. |
| **Span** | Ein einzelner Arbeitsabschnitt im Trace — HTTP-Call, DB-Query, Funktionsausführung. Eigene **Span-ID** (8 Byte / 16 hex). |
| **Hop** | Service-zu-Service-Übergang. Pro Hop **neue Span-ID**, **gleiche Trace-ID**. |
| **`traceparent`** | Der W3C-Header. Vier Felder mit festen Längen, sprach- und vendor-neutral. |
| **Sampling-Flag** | Bit 0 in `flags`. „Soll dieser Trace exportiert werden?" — Cost-Control. Propagiert über alle Hops. |

### Der `traceparent`-Header

```
traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
             │  └────── Trace-ID (32 hex) ─────┘ └Span-ID (16 hex)┘ └Flags
             │
             └── Version (immer 00)
```

| Feld | Länge | Bedeutung |
|---|---|---|
| `version` | 2 hex | Aktuell immer `00`. Höhere Versionen müssen abwärtskompatibel ignoriert werden. |
| `trace-id` | 32 hex | 16 Byte zufällig. **Bleibt über alle Hops gleich.** Darf nicht nur Nullen sein. |
| `parent-id` (Span-ID) | 16 hex | 8 Byte zufällig. **Pro Hop neu**. Aus Empfänger-Sicht: Span, die ihn aufruft. |
| `flags` | 2 hex | Bitmaske. Bit 0 = „sampled". |

> ⚠️ Nur-Nullen-Trace oder -Span sind **ungültig**. `crypto/rand` benutzen, nicht `math/rand` mit Default-Seed.

Optional: `tracestate` für Vendor-spezifische Daten — im Workshop nicht relevant.

---

## 3. Beispiel: Eine Buchung durch unseren Stack

Im Workshop ist die Aufgabenverteilung **asymmetrisch**: **Booking** ist der einzige Entry-Point, **Flight / Hotel / Car** sind passive Trace-Empfänger.

### Header- und Log-Verlauf

```
Client ──► Booking
            │
            │  Server-Middleware: kein traceparent → generateNew()
            │    T  = 0af7651916cd43dd8448eb211c80319c   (Trace-ID, bleibt überall gleich)
            │    S0 = a1a1a1a1a1a1a1a1                   (initiale Span-ID im Booking-Context)
            │
            │  Booking loggt:  {"trace_id":"0af7…319c", "span_id":"a1a1…a1a1", "msg":"booking start"}
            │
            │  Client-Inject vor Call zu Flight: neue Span-ID
            │    S1 = b7ad6b7169203331
            │
            └─traceparent: 00-0af7…319c-b7ad…3331-01─►  Flight
                                                       │
                                                       │  Server-Middleware (Propagate):
                                                       │    parse → ctx = {T, S1}
                                                       │
                                                       │  Flight loggt:  {"trace_id":"0af7…319c", "span_id":"b7ad…3331", "msg":"flight booked"}
                                                       │
                                                       │  (Wenn Flight weiterruft:)
                                                       │  Client-Inject vor downstream-Call: neue Span-ID
                                                       │    S2 = c3c3c3c3c3c3c3c3
                                                       │
                                                       └─traceparent: 00-0af7…319c-c3c3…c3c3-01─►  external Provider
                                                                                                  │
                                                                                                  │  loggt:  {"trace_id":"0af7…319c", "span_id":"c3c3…c3c3", …}
                                                                                                  ▼
```

### Wer erstellt was?

**Booking** (Entry-Point):
1. **Trace-ID T** — einmalig, falls kein `traceparent` reinkommt (Browser/curl-Fall).
2. **Span-ID S0** — initiale Span im Booking-Context. **Mit dieser ID loggt Booking während der ganzen Anfrage.**
3. **Span-IDs S1, S2, S3** — pro Outbound-Call zu Flight, Hotel, Car. Diese leben **nur im Outbound-Header**, nicht im Booking-Log.
4. **`traceparent` als Event-Property** — beim async Publish in Story 6/7 (siehe Abschnitt 4).

**Flight / Hotel / Car** (Downstream): erstellen **nichts** für eingehende Requests.
- `Propagate`-Middleware übernimmt den eintreffenden `traceparent` in den Request-Context.
- Loggen mit der mitgelieferten `trace_id` + `span_id`.
- **Fehlt der Header → keine `trace_id` im Log**, kein Fallback. Absicht: nur der Entry-Point sät Trace-IDs.

**Wenn ein Downstream selbst weiterruft** (z. B. Flight ruft externen Provider):
- Pattern identisch: `Propagate` übernimmt eingehenden Kontext, `Client-Inject` erzeugt **neue Span-ID** für den Outbound-Hop.
- Trace-ID T wandert weiter, Span-ID wird neu (S2 im Diagramm oben).
- Flight darf in diesem Fall **niemals** eine neue Trace-ID erzeugen — sonst zerfällt der Trace genau dort.

**Faustregel:** Jeder Service, der Aufrufe annimmt und weiterleitet, braucht eine Server-Middleware vom Typ *Propagate*. Eine *Middleware* (die einen neuen Trace bei Bedarf erzeugt) gehört **exklusiv an die Systemgrenze** — API-Gateway oder Public-Service.

### Was wird geloggt?

**Nicht den ganzen `traceparent`-Header**, sondern `trace_id` und `span_id` als **getrennte Felder** im strukturierten Log:

```json
{"time":"2026-05-12T10:00:01Z", "level":"INFO", "msg":"flight booked",
 "trace_id":"0af7651916cd43dd8448eb211c80319c", "span_id":"b7ad6b7169203331",
 "step":"flight", "duration_ms":120}
```

| Service | `trace_id` | `span_id` | Anzahl Logzeilen mit dieser Span |
|---|---|---|---|
| Booking | `0af7…319c` | `a1a1…a1a1` (S0) | viele — alle Booking-Logs |
| Flight | `0af7…319c` | `b7ad…3331` (S1) | viele — alle Flight-Logs |
| Hotel | `0af7…319c` | (S2, eigene) | viele — alle Hotel-Logs |
| Car | `0af7…319c` | (S3, eigene) | viele — alle Car-Logs |
| External Provider | `0af7…319c` | `c3c3…c3c3` | viele — alle Provider-Logs |

Effekt:

- `grep 0af7…319c` über alle Service-Logs → **kompletter Vorgang** zusammenhängend.
- `grep b7ad…3331` → **nur Flight**.
- `docker compose logs | jq -c 'select(.trace_id=="0af7…319c")'` → sauber gefilterter Live-Stream.

### Der Pseudo-Code dazu

```kotlin
// Server-Middleware (Entry-Point: erzeugt, falls keiner kommt)
fun middleware(req, next) {
    val tc = parse(req.header("traceparent")) ?: generateNew()
    next(req.withContext(ctx.with(tc)))
}

// Server-Middleware (Downstream: nur übernehmen, nichts erzeugen)
fun propagate(req, next) {
    val tc = parse(req.header("traceparent"))
    if (tc != null) next(req.withContext(ctx.with(tc))) else next(req)
}

// Client-Inject: pro Outbound-Hop neue Span-ID
fun inject(ctx, outboundReq) {
    val incoming = ctx.get(TraceContext)
    val hop = incoming.copy(spanId = randomSpanId())
    outboundReq.setHeader("traceparent", hop.toHeader())
}

// Logger: trace_id + span_id automatisch in jede Zeile
fun log(ctx, msg, vararg fields) {
    val tc = ctx.get(TraceContext)
    logger.info(msg, "trace_id" to tc.traceId, "span_id" to tc.spanId, *fields)
}
```

Reference-Implementierung: `services/shared/tracing/tracing.go` (~100 Zeilen).

---

## 4. Tracing über Async-Grenzen

Bei der Choreography-Saga (Story 6) reicht der HTTP-Header nicht mehr — der Trace-Kontext muss **als Property auf dem Event** mitwandern:

```json
{
  "eventId": "...",
  "sagaId": "...",
  "traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
  "payload": { ... }
}
```

Der Konsument liest die Property, legt den Trace-Kontext in den Worker-Goroutine-Context und loggt damit weiter. Andernfalls bricht der Trace genau an der Bus-Grenze ab.

> ⚠️ Bei echten Brokern (Kafka, RabbitMQ, SNS/SQS): **Message-Header** benutzen, nicht den Payload. Sonst muss jeder Konsument das Payload-Schema kennen, nur um den Trace zu propagieren.

---

## 5. Strukturiertes Logging — warum?

Mit `log.Printf("booking %s failed: %v", id, err)` muss `trace_id=...` manuell in jeden Format-String. Wird vergessen, ID fehlt in der Hälfte der Logzeilen.

Mit strukturiertem Logger (`slog` in Go, `structlog` Python, `logback` JSON in Java, `pino` Node): ID einmal an den Request-Logger hängen, sie taucht in **jeder** Zeile auf:

```go
logger := tracing.Logger(ctx)              // hat trace_id schon dran
logger.Info("forward step done", "step", "flight", "duration_ms", 120)
```

JSON-Output macht das spätere Filtern mit `jq` trivial.

---

## 6. Selbst bauen oder OpenTelemetry?

| Aspekt | Selbst | OpenTelemetry |
|---|---|---|
| Lerneffekt | ⭐⭐⭐⭐⭐ | ⭐⭐ — Auto-Instrumentation versteckt den Mechanismus |
| Korrektheit | Header-Parsing strikt halten, Race-Bugs möglich | Kampferprobt |
| Span-Baum (Parent/Child) | Müsste man dazubauen | Inklusive |
| Sampling, Batching, Export | Müsste man bauen | Inklusive |
| Backend-Export (Jaeger, Tempo, Datadog…) | Pro Backend selbst | Inklusive (OTLP) |
| Zeilen Code | ~100 | 0 + ein paar Zeilen Konfiguration |
| Sprach-/Tool-Bindung | keine | pro Sprache eigenes SDK |

**Workshop:** selbst bauen — Mechanismus verstehen.
**Praxis:** OpenTelemetry — Auto-Instrumentation, Sampling, OTLP-Export wollt ihr nicht neu erfinden.

---

## 7. Marktüberblick

### Wire-Formate

| Standard | Status | Wo verbreitet |
|---|---|---|
| **W3C Trace Context** (`traceparent` / `tracestate`) | Aktuell, Default | Überall (OTel, neue Datadog, neue New Relic) |
| **B3** (`X-B3-TraceId` / …) | Legacy, weit verbreitet | Zipkin, ältere Spring-Cloud |
| **Jaeger native** (`uber-trace-id`) | Legacy | Ältere Jaeger-only-Setups |
| **Datadog** (`x-datadog-*`) | Vendor-spezifisch, mittlerweile auch W3C | Ältere Datadog-Agents |

### SDKs

| Library | Beschreibung |
|---|---|
| **OpenTelemetry** | De-facto-Standard. SDK + Auto-Instrumentation. Vendor-agnostisch über OTLP. |
| OpenTracing / OpenCensus | Deprecated, in OTel aufgegangen. |
| Vendor-SDKs (Datadog APM, New Relic, Dynatrace, …) | Komfortabel, aber Vendor-Lock-in. |
| Spring Cloud Sleuth / Micrometer Tracing | Java. Sleuth abgekündigt, Micrometer ist Nachfolger. |

### Backends

| Tool | Typ | Stärken |
|---|---|---|
| **Jaeger** | OSS | Klassiker, einfach lokal, CNCF-Projekt |
| **Grafana Tempo** | OSS | Sehr günstig (Object-Storage), Grafana-Integration |
| **Zipkin** | OSS | Älter, schlank |
| **Datadog APM** | SaaS | Reife UI, Logs+Traces+Metrics korreliert |
| **Honeycomb** | SaaS | High-cardinality Queries |
| **Elastic APM** | OSS/SaaS | ELK-Integration |
| **AWS X-Ray** | SaaS | Nahtlos in AWS |

Für den Workshop-Bonus: **Jaeger All-in-One** ist die kürzeste Strecke zur ersten Trace-UI.

---

## 8. Sampling — die Cost-Frage

| Strategie | Wann entscheiden | Vor- / Nachteil |
|---|---|---|
| **Head-based** (z. B. 1 %) | Beim Eingang | + einfach / − Fehler-Traces oft verloren |
| **Tail-based** | Nach Trace-Ende | + Fehler/Slow immer behalten / − Backend muss buffern |
| **Adaptive** | Dynamisch | + folgt Last / − komplex zu betreiben |
| **Always-on für Errors / Slow** | Per Heuristik | + beste Sichtbarkeit / − braucht Tail-based-Stack |

> 💬 *Diskussionsfrage:* Wer im Workshop nutzt Tracing produktiv? Wie hoch ist eure Sampling-Rate?

---

## 9. Diskussionsfragen

1. **Wer schreibt die Trace-ID heute schon in eure Logs?** Häufige Lücke: das Gateway generiert eine, aber die Backends schreiben sie nicht in die App-Logs.
2. **Wie weit reicht euer Trace?** Bis ins Frontend? Über CDN-Hops? Mobile-App-SDK?
3. **Sampling: head- oder tail-based?** Wer entscheidet die Rate?
4. **Korrelation Logs ↔ Traces ↔ Metriken** — alles im selben Tool, oder pro Frage Tool-Wechsel?
5. **Async-Tracing:** Bei wem laufen Events/Queues? Wandert die Trace-ID dort schon mit?
6. **Was kostet euch Tracing?** Storage, Bandbreite, Backend-Lizenz?

---

## 10. Wo der Code liegt

| Was | Pfad |
|---|---|
| Tracing-Bibliothek (Parse, Generate, Middleware, Inject, Logger) | `services/shared/tracing/tracing.go` |
| Tests | `services/shared/tracing/tracing_test.go` |
| Server-Middleware verdrahten | `services/booking/story7/main.go`, `services/{flight,hotel,car}/main.go` |
| Outbound-Inject + strukturierte Logs | `services/booking/story7/handler/booking.go` |
| Trace-Kontext in Compensation-Events | `services/booking/story7/saga/*.go`, `services/{flight,hotel,car}/handler/compensation.go` |
| Backend-Handler nutzen Trace-Logger | `services/{flight,hotel,car}/handler/*.go` |
