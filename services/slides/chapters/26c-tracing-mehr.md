<!-- .slide: data-background-image="./assets/tracing.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Tracing kann mehr

<p class="subtitle">&hellip; als wir hier nutzen</p>

<div class="box">

- **Vollwertiger Span-Baum** &mdash; Parent/Child-Verkn&uuml;pfung pro Hop, Span-Dauer pro Schritt, Wartezeit zwischen Spans
- **OpenTelemetry SDK** &mdash; Auto-Instrumentation f&uuml;r HTTP-Clients, DB-Treiber, Message-Broker
- **OTLP-Export** &mdash; vendor-neutraler Push an Jaeger / Tempo / Datadog / Honeycomb / &hellip;
- **Sampling-Strategien** &mdash; Head-based (1&nbsp;%), Tail-based (Fehler immer behalten), Adaptive
- **Korrelation Logs &harr; Traces &harr; Metriken** &mdash; ein Tool, ein Klick vom Log zur Span zur RED-Metrik
- **Frontend-Tracing** &mdash; RUM-SDKs ziehen den Trace bis in den Browser, &uuml;ber das Backend hinweg
- **Continuous Profiling** &mdash; Pyroscope / Parca / Datadog Profiler erg&auml;nzen Spans mit Flamegraphs
- **Exemplars** &mdash; Prometheus-Metriken verlinken direkt auf die zugeh&ouml;rige Trace
- ...

</div>

<p class="quote">Im Workshop selbst gebaut. In Produktion: <span class="hl">OpenTelemetry</span>.</p>

Note:
- Span-Baum: in der Workshop-Variante haben wir bewusst nur Trace-IDs gebaut, keinen vollwertigen Parent/Child-Baum. F&uuml;r 60 Min zu viel. In Jaeger / Tempo / Datadog wird daraus eine baumartige Visualisierung mit Spans pro HTTP-Call, DB-Query, etc. &mdash; das ist der eigentliche &bdquo;Wow-Effekt&ldquo;, den die SDKs liefern.
- OpenTelemetry-Empfehlung in einer Zeile: <em>&bdquo;Auto-Instrumentation, Sampling, Span-Lifecycle, OTLP-Export an beliebige Backends &mdash; Arbeit, die niemand neu bauen will.&ldquo;</em>
- Sampling-Tabelle:
  <pre>head-based     : Entscheidung beim Eingang &mdash; kein Buffering, aber Fehler-Traces gehen verloren
tail-based     : Entscheidung am Trace-Ende &mdash; Fehler/Slow immer behalten, Backend muss buffern
adaptive       : Rate folgt Last &mdash; komplex zu betreiben
always-on errs : Heuristik, beste Sichtbarkeit der Probleme</pre>
- Korrelation: das eigentliche Versprechen moderner Observability-Plattformen. Ein langsamer Endpunkt &rarr; RED-Metrik zeigt P99-Spike &rarr; ein Klick auf Exemplar &rarr; konkreter Trace &rarr; in der Trace-Span auf &bdquo;Logs anzeigen&ldquo; &rarr; alle Logs derselben Trace-ID. Wer diese Kette nicht hat, wechselt drei Tools pro Frage.
- Continuous Profiling ist die n&auml;chste Schicht &uuml;ber Traces &mdash; statt &bdquo;wie lange dauert der Call&ldquo; sieht man <em>wo</em> die CPU-Zeit verbrannt wird. Pyroscope (Grafana), Parca, Datadog Profiler.
- Wichtigste Take-aways:
  <pre>Selbst gebaut: lehrreich, sprach-neutral, 100 Zeilen.
OpenTelemetry: produktiv, vendor-neutral, sprach-spezifisch.
Beides parallel macht keinen Sinn &mdash; eins von beiden, mit Plan.</pre>
