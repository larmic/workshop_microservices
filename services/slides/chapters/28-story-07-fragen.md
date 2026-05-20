## Story 7 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Trace vs. Span</h3>
<p>Wir haben Trace-IDs propagiert, aber <span class="hl">keinen</span> Span-Baum gebaut. Was fehlt uns?</p>
<code>Span = ein Hop</code>
<aside class="notes"><strong>Meine Antwort:</strong> Ein <strong>Trace</strong> ist der komplette Gesch&auml;ftsvorgang, identifiziert &uuml;ber die Trace-ID (16 Byte, 32 hex). Eine <strong>Span</strong> ist ein einzelner Arbeitsabschnitt im Trace &mdash; typischerweise ein HTTP-Call, ein DB-Query, eine Funktionsausf&uuml;hrung &mdash; mit eigener Span-ID (8 Byte) und Parent-Span-ID. Daraus entsteht der Span-Baum, den Jaeger / Tempo / Datadog als &bdquo;Wasserfall&ldquo; visualisieren. Im Workshop bewusst weggelassen: f&uuml;r 60 Min zu viel, Tool-spezifisch.<br><strong>Spicy:</strong> Ohne Span-Baum hat man &bdquo;Logs mit Trace-ID&ldquo; &mdash; das ist 80&nbsp;% des Nutzens f&uuml;r 20&nbsp;% des Aufwands. Die restlichen 20&nbsp;% (Visualisierung, P99-pro-Span, kritischer Pfad) sind genau das, wof&uuml;r OpenTelemetry da ist. Selbst bauen w&auml;re Folklore.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Wer initiiert den Trace?</h3>
<p>Nur Booking. Flight/Hotel/Car erzeugen <span class="hl">nie</span> selbst eine. Warum nicht?</p>
<code>Entry-Point only</code>
<aside class="notes"><strong>Meine Antwort:</strong> Trace-Initiierung geh&ouml;rt an den Entry-Point (API-Gateway, Public-Service) &mdash; nicht in jeden Downstream-Hop. W&uuml;rde Flight selbst eine Trace-ID erzeugen, w&uuml;rden Aufrufe ohne propagierten Header einen <em>neuen</em> Trace starten &mdash; obwohl sie eigentlich Teil eines &uuml;bergeordneten Vorgangs sind. Im Code: zwei Middlewares mit unterschiedlicher Rolle. <code>Middleware</code> erzeugt einen Trace, falls keiner kommt (Entry-Point). <code>Propagate</code> &uuml;bernimmt nur einen vorhandenen Trace, sonst nichts (Downstream).<br><strong>Spicy:</strong> Wer in jedem Service &bdquo;sicherheitshalber&ldquo; eine Trace-ID generiert, baut sich tausende Mini-Traces. Im Tracing-Tool sieht das aus wie hohe Aktivit&auml;t, beim Debuggen ist es n-fache Detektivarbeit.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Strukturiertes Logging</h3>
<p>JSON statt Format-String. Warum ist <code>log.Printf</code> mit eingebauter Trace-ID nicht <span class="hl">gut genug</span>?</p>
<code>trace_id als Feld</code>
<aside class="notes"><strong>Meine Antwort:</strong> Mit <code>log.Printf("booking %s failed: %v", id, err)</code> muss man <code>trace_id=...</code> manuell in jeden Format-String einsetzen. Nach drei Wochen vergessen, ab dann fehlt die ID in der H&auml;lfte der Logzeilen. Strukturierter Logger (slog, structlog, logback JSON, pino) h&auml;ngt die ID einmal am Request-Kontext-Logger &mdash; sie taucht in <em>jeder</em> Zeile auf, ohne Format-String. Konsequenz: <code>docker compose logs | jq -c 'select(.trace_id=="abc...")'</code> gibt einen sauber gefilterten Live-Stream.<br><strong>Spicy:</strong> &bdquo;Wir loggen schon mit Trace-ID&ldquo; ist die Antwort von Teams, die Format-Strings benutzen. Ein Blick in die Logs reicht meist: ist die ID an einer Stelle ein eigenes Feld, an drei anderen nur im freien Text? Dann ist Logging-Korrelation Theater.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Async-Grenze</h3>
<p>Bei den Compensation-Events aus Story 6: der HTTP-Header ist <span class="hl">weg</span>. Wie kommt der Trace mit?</p>
<code>traceparent als Property</code>
<aside class="notes"><strong>Meine Antwort:</strong> Der Trace-Kontext muss <em>aktiv</em> als Event-Property mitwandern. Beim Konsumenten parsen, in den Worker-Goroutine-Kontext legen, mit dem Logger fortf&uuml;hren. Bei echten Brokern (Kafka / RabbitMQ / SNS-SQS): in die <strong>Message-Header</strong>, nicht in den Payload &mdash; sonst muss jeder Consumer das Payload-Schema kennen, nur um den Trace zu propagieren. In unserer Webhook-Variante ist der Trace im Event-Body, weil wir keine Header haben.<br><strong>Spicy:</strong> Async-Tracing ist die unausgesprochene Pflicht jeder Eventing-Architektur. Wer nur HTTP-Header propagiert und Events als &bdquo;ist halt async&ldquo; behandelt, hat einen Trace, der genau da abreisst, wo es spannend wird &mdash; an der Bus-Grenze.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Selbst bauen oder OpenTelemetry?</h3>
<p>100 Zeilen Code reichen f&uuml;r Propagation. Wozu OpenTelemetry?</p>
<code>OTel = SDK + Export</code>
<aside class="notes"><strong>Meine Antwort:</strong> Header-Parsing und Propagation sind in ca. 100 Zeilen abgehandelt &mdash; lehrreich, sprach-neutral. <em>Aber</em>: Span-Lifecycle, Auto-Instrumentation f&uuml;r HTTP-Clients / DB-Treiber / Broker, Sampling-Strategien, OTLP-Export an beliebige Backends &mdash; das ist die eigentliche Arbeit, die niemand neu bauen will. Empfehlung: <strong>Im Workshop selbst bauen</strong>, damit man wei&szlig;, was die Library macht. <strong>In Produktion OpenTelemetry</strong>.<br><strong>Spicy:</strong> Wer in Produktion Tracing selbst implementiert, hat OpenTelemetry-Komitee-Arbeit nicht verfolgt. Multi-Propagator-Support, Tail-based Sampling, kontextbasiertes Logging, Exemplars f&uuml;r Prometheus &mdash; das alles ist Lebenszeit, die das eigene Team nicht bekommt zur&uuml;ck.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Sampling</h3>
<p>Jeden Request tracen ist <span class="hl">teuer</span>. Head- oder Tail-based?</p>
<code>Fehler immer behalten</code>
<aside class="notes"><strong>Meine Antwort:</strong> In Produktion sind Traces teuer (Storage, Netzwerk, Backend-Lizenzen) &mdash; niemand tracet 100&nbsp;%. <strong>Head-based</strong> (z.&nbsp;B. 1&nbsp;%): Entscheidung beim Eingang &mdash; einfach, aber Fehler-Traces gehen oft verloren. <strong>Tail-based</strong>: Entscheidung am Trace-Ende, Fehler/Slow-Traces werden immer behalten &mdash; teurer im Backend (Buffering n&ouml;tig). <strong>Adaptive</strong>: Rate folgt Last &mdash; komplex zu betreiben. Pragmatischer Default: head-based 1&nbsp;% + &bdquo;Always-on f&uuml;r Errors&ldquo; per Heuristik.<br><strong>Spicy:</strong> Sampling-Rate ohne Plan ist Self-Sabotage. 100&nbsp;% &rarr; Tracing-Backend kostet mehr als der Stack selbst. 0,01&nbsp;% &rarr; wenn ein Kunde ein Problem hat, ist genau sein Trace nicht da.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">7</span> Logs &harr; Traces &harr; Metriken</h3>
<p>Habt ihr alle drei im <span class="hl">selben Tool</span> &mdash; oder wechselt ihr pro Frage?</p>
<code>One-Click-Korrelation</code>
<aside class="notes"><strong>Meine Antwort:</strong> Das ist das eigentliche Versprechen moderner Observability-Plattformen: ein langsamer Endpunkt &rarr; RED-Metrik zeigt P99-Spike &rarr; ein Klick auf das <em>Exemplar</em> &rarr; konkreter Trace &rarr; in der Trace-Span auf &bdquo;Logs&ldquo; &rarr; alle Logs derselben Trace-ID. Wer diese Kette nicht hat, wechselt drei Tools pro Frage &mdash; und meist verliert man eine Verkn&uuml;pfung an einer der &Uuml;berg&auml;nge. Konkrete Stacks, die das k&ouml;nnen: Grafana (Loki + Tempo + Mimir + Exemplars), Datadog, Honeycomb, Dynatrace.<br><strong>Spicy:</strong> &bdquo;Wir haben Tracing&ldquo; ist nicht dasselbe wie &bdquo;wir nutzen Tracing&ldquo;. Wer im Trace-Tool landet und dann manuell in Kibana die <code>trace_id</code> per Copy-Paste eingibt, hat die Korrelation noch vor sich.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Diese Diskussionspunkte basieren auf <code>docs/instructions/distributed-tracing.md</code> (Abschnitte 10 und 9). Eine eigene <code>docs/questions/story7.md</code> existiert (bisher) nicht.
</aside>
