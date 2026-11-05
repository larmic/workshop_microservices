# Workshop-Fragen: Distributed Tracing (Story 8)

Fragen zur Recap-Folie 28. Ziel: Tracing nicht als „ID in die Logs
schreiben" abhaken, sondern verstehen, was der rote Faden leistet, was
ihm fehlt und wo er reißt. Die Teilnehmenden haben ihren eigenen
Story-6-Service instrumentiert (synchrone Saga, Kompensation per
`DELETE`), keine Events.

---

## Frage 1 — Wir haben Trace-IDs propagiert, aber keinen Baum gebaut. Was fehlt uns? (Recap 01)

**Antwort:** Ein Trace ist der ganze Vorgang, eine Span ein einzelnes
Wegstück mit eigener ID und Parent-ID. Erst die Parent-Beziehung ergibt
den Baum, den Jaeger, Tempo oder Datadog als Wasserfall zeigen: Dauer
pro Span, kritischer Pfad, Wartezeit zwischen den Aufrufen. Unsere
Library verwirft die alte Span-ID bei jedem Hop, deshalb gibt es nur
die Korrelation über die Trace-ID.

Was wir damit können: alle Zeilen einer Buchung finden. Was nicht:
sehen, welcher Aufruf wie lange gedauert hat und wer auf wen gewartet
hat.

**Spicy:** Logs mit Trace-ID sind 80 Prozent des Nutzens für
20 Prozent des Aufwands. Die restlichen 20 Prozent sind genau das,
wofür OpenTelemetry da ist.

---

## Frage 2 — Flight, Hotel, Car würfeln nie selbst eine ID. Warum nicht sicherheitshalber? (Recap 02)

**Antwort:** Trace-Initiierung gehört an die Systemgrenze (API-Gateway,
Public-Service). Würde jedes Backend „sicherheitshalber" eine ID
erzeugen, starten Aufrufe ohne Header einen eigenen Trace, obwohl sie
Teil eines größeren Vorgangs sind. Im Tool sieht das nach viel
Aktivität aus, beim Debuggen sind es tausend Mini-Traces.

In der Referenz: `tracing.Middleware` erzeugt (Booking),
`tracing.Propagate` übernimmt nur (Flight, Hotel, Car). Deshalb sind
die Backend-Logs in Story 1 bis 7 ohne `trace_id`, der Kontrast macht
den Effekt sichtbar.

**Die eine Ausnahme:** Ein Compensation-Event ohne
`traceparent`-Property startet in den Backends einen neuen Trace
(`contextFromEvent`). Das ist bewusst so, damit man im Event-Bonus sieht,
wie der Trace zerfällt.

---

## Frage 3 — Warum ist ein Format-String mit eingebauter Trace-ID nicht gut genug? (Recap 03)

**Antwort:** Mit einem Format-String muss die ID von Hand in jede
Logzeile. Nach drei Wochen ist sie in der Hälfte der Zeilen vergessen.
Ein strukturierter Logger (slog, structlog, Logback JSON, pino) hängt
sie einmal an den Request-Logger, danach steht sie überall, und zwar
als eigenes Feld, nach dem `jq` filtern kann.

**Spicy:** „Wir loggen schon mit Trace-ID" ist oft die Antwort von
Teams mit Format-Strings. Ist die ID an einer Stelle Feld und an drei
anderen freier Text, ist die Korrelation Theater.

---

## Frage 4 — Ein Zettel hat keinen HTTP-Header. Wie kommt der Trace mit? (Recap 04)

**Frage:** Im Story-7-Rollenspiel lag die Stornierung als Zettel auf dem
Tisch. Hätten wir das programmiert, wäre es ein Event gewesen. Wo steht
dann die Trace-ID?

**Antwort:** Der Kontext muss aktiv mitwandern: beim Publish als
Property auf das Event, beim Konsumenten parsen, in den Worker-Kontext
legen, mit dem Logger fortführen. Bei echten Brokern (Kafka, RabbitMQ,
SNS/SQS) gehört er in die Message-Header, nicht in den Payload, sonst
muss jeder Consumer das Schema kennen, nur um den Trace weiterzureichen.

Wer den Bonus gemacht hat, kann es zeigen: Ohne Property trägt die
`request`-Zeile im Backend noch die richtige ID, `received`,
`processing` und `done` tragen eine neue.

**Spicy:** Wer nur HTTP-Header propagiert, hat einen Trace, der genau
da abreißt, wo es spannend wird: an der Bus-Grenze.

---

## Frage 5 — Hundert Zeilen reichen. Wozu dann OpenTelemetry? (Recap 05)

**Antwort:** Parsing und Propagation sind in hundert Zeilen erledigt,
lehrreich und sprachneutral. Span-Lebenszyklus, Auto-Instrumentierung
für HTTP-Clients, DB-Treiber und Broker, Sampling und Export an
beliebige Backends sind die Arbeit, die niemand neu bauen will. Im
Workshop selbst bauen, in Produktion OpenTelemetry.

---

## Reserve

- **Sampling:** Head-based (etwa 1 Prozent) entscheidet am Eingang,
  Fehler-Traces gehen dabei verloren. Tail-based entscheidet am Ende,
  das Backend muss puffern. Pragmatisch: head-based plus immer bei
  Fehlern.
- **Korrelation:** Metrik-Spike, Exemplar, konkreter Trace, dazu die
  Logs derselben Trace-ID. Wer diese Kette nicht hat, wechselt pro
  Frage drei Tools.
- **Detektiv-Runde:** Wie habt ihr die kaputte Buchung gefunden? Fast
  immer über die Fehlerzeile von Car und dann die `trace_id`. Genau so
  sucht man in Produktion: vom Symptom zur ID, von der ID zum Ablauf.
