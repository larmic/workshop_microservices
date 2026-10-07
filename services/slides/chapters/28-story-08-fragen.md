<div class="page">

<p class="kicker">Story 8 &middot; Recap</p>

## F&uuml;nf Fragen an euch

<div class="page-body">

<div class="numlist recap compact dense">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Trace ohne Span-Baum</p>
<h3>Wir haben Trace-IDs propagiert, aber keinen Baum gebaut. Was <span class="hl">fehlt</span> uns?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">80 % f&uuml;r 20 %</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Nur Booking erzeugt einen Trace</p>
<h3>Flight, Hotel, Car w&uuml;rfeln <span class="hl">nie</span> selbst eine ID. Warum nicht sicherheitshalber?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Entry-Point only</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">JSON statt Format-String</p>
<h3>Warum ist <code>log.Printf</code> mit eingebauter Trace-ID nicht <span class="hl">gut genug</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">trace_id als Feld</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Die Async-Grenze</p>
<h3>Beim Compensation-Event ist der HTTP-Header <span class="hl">weg</span>. Wie kommt der Trace mit?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">traceparent als Property</code>
</div>
<div class="numlist-row">
<div class="numlist-num">05</div>
<div>
<p class="numlist-label">Hundert Zeilen reichen</p>
<h3>Wozu dann <span class="hl">OpenTelemetry</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">OTel = SDK + Export</code>
</div>
</div>

</div>

</div>

Note:
- Alle f&uuml;nf Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein. Einstieg aus <code>docs/themen.md</code>: &bdquo;Wo h&auml;tte euch Tracing in Stories 4 bis 7 schon geholfen?&ldquo;
- <strong>Span-Baum, meine Antwort:</strong> Ein Trace ist der komplette Vorgang, eine Span ein einzelnes Wegst&uuml;ck mit eigener ID und Parent-ID. Daraus entsteht der Baum, den Jaeger, Tempo oder Datadog als Wasserfall zeigen: Dauer pro Span, kritischer Pfad, Wartezeit zwischen Spans. Im Workshop bewusst weggelassen. <strong>Spicy:</strong> Logs mit Trace-ID sind 80 Prozent des Nutzens f&uuml;r 20 Prozent des Aufwands. Die restlichen 20 Prozent sind genau das, wof&uuml;r OpenTelemetry da ist. Selbst bauen w&auml;re Folklore.
- <strong>Entry-Point, meine Antwort:</strong> Trace-Initiierung geh&ouml;rt an die Systemgrenze (API-Gateway, Public-Service), nicht in jeden Downstream-Hop. W&uuml;rde Flight selbst eine ID erzeugen, starten Aufrufe ohne Header einen neuen Trace, obwohl sie Teil eines gr&ouml;&szlig;eren Vorgangs sind. In der Referenz: <code>Middleware</code> erzeugt, <code>Propagate</code> &uuml;bernimmt nur. <strong>Spicy:</strong> Wer in jedem Service &bdquo;sicherheitshalber&ldquo; eine ID w&uuml;rfelt, baut sich tausende Mini-Traces. Im Tool sieht das nach viel Aktivit&auml;t aus, beim Debuggen ist es n-fache Detektivarbeit.
- <strong>Format-String, meine Antwort:</strong> Mit <code>log.Printf</code> muss die ID in jeden Format-String von Hand. Nach drei Wochen vergessen, ab dann fehlt sie in der H&auml;lfte der Zeilen. Ein strukturierter Logger (slog, structlog, logback JSON, pino) h&auml;ngt sie einmal am Request-Logger an, dann steht sie &uuml;berall. <strong>Spicy:</strong> &bdquo;Wir loggen schon mit Trace-ID&ldquo; ist die Antwort von Teams mit Format-Strings. Ist die ID an einer Stelle ein Feld und an drei anderen freier Text, ist die Korrelation Theater.
- <strong>Async-Grenze, meine Antwort:</strong> Der Kontext muss aktiv als Event-Property mitwandern, beim Konsumenten parsen, in den Worker-Kontext legen, mit dem Logger fortf&uuml;hren. Bei echten Brokern (Kafka, RabbitMQ, SNS/SQS) in die Message-Header, nicht in den Payload, sonst muss jeder Consumer das Schema kennen, nur um den Trace weiterzureichen. <strong>Spicy:</strong> Wer nur HTTP-Header propagiert, hat einen Trace, der genau da abrei&szlig;t, wo es spannend wird: an der Bus-Grenze.
- <strong>OpenTelemetry, meine Antwort:</strong> Parsing und Propagation sind in 100 Zeilen erledigt, lehrreich und sprachneutral. Span-Lifecycle, Auto-Instrumentation f&uuml;r HTTP-Clients, DB-Treiber und Broker, Sampling, OTLP-Export an beliebige Backends: Das ist die Arbeit, die niemand neu bauen will. Im Workshop selbst bauen, in Produktion OpenTelemetry. <strong>Spicy:</strong> Wer in Produktion Tracing selbst schreibt, bezahlt mit Lebenszeit f&uuml;r Multi-Propagator-Support, Tail-based Sampling und Exemplars, die das Team nie zur&uuml;ckbekommt.
- <strong>Reserve, Sampling:</strong> Head-based (etwa 1 Prozent) entscheidet am Eingang, einfach, aber Fehler-Traces gehen verloren. Tail-based entscheidet am Ende, Fehler und langsame Traces bleiben immer, Backend muss puffern. Pragmatischer Default: head-based plus always-on bei Fehlern. Ohne Plan: 100 Prozent, und das Tracing kostet mehr als der Stack. 0,01 Prozent, und beim Kundenproblem fehlt genau sein Trace.
- <strong>Reserve, Korrelation:</strong> Logs, Traces und Metriken im selben Tool? P99-Spike in der RED-Metrik, Klick aufs Exemplar, konkreter Trace, Klick auf &bdquo;Logs&ldquo;, alle Zeilen derselben Trace-ID. Wer die Kette nicht hat, wechselt drei Tools pro Frage. Grafana (Loki, Tempo, Mimir), Datadog, Honeycomb k&ouml;nnen das.
- Grundlage: <code>docs/instructions/distributed-tracing.md</code>, Abschnitte 9 und 10. Eine eigene <code>docs/questions/story8.md</code> gibt es bisher nicht.
