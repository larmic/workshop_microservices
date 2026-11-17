<div class="page">

<p class="kicker">Distributed Tracing</p>

<div class="page-head">

## Eine ID, vier Services

<img class="head-figure" src="./assets/tracing.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Trace</h3>
<p>Der ganze Gesch&auml;ftsvorgang. Eine ID f&uuml;r alle Services entlang einer Anfrage.</p>
<code>Trace-ID: 16 Byte, 32 hex</code>
</div>
<div class="card">
<h3>Span</h3>
<p>Ein Wegst&uuml;ck im Trace: ein HTTP-Call, ein DB-Query, eine Funktion.</p>
<code>Span-ID: 8 Byte, 16 hex</code>
</div>
<div class="card">
<h3>Weitergabe</h3>
<p>Der &Uuml;bergang von Service zu Service. Neue Span-ID, gleiche Trace-ID.</p>
<code>Booking &rarr; Flight = neue Span-ID</code>
</div>
<div class="card">
<h3>traceparent</h3>
<p>Der W3C-Header. Sprach- und vendor-neutral, 55 Zeichen.</p>
<code>00-&lt;trace 32&gt;-&lt;parent 16&gt;-&lt;flags 2&gt;</code>
</div>
<div class="card">
<h3>Sampling-Flag</h3>
<p>Letztes Feld im Header: Wird dieser Vorgang ans Tracing-Tool geschickt? Alle aufzuzeichnen w&auml;re zu teuer, also nur eine Stichprobe. F&uuml;r die Logs egal.</p>
<code>01 = aufzeichnen &middot; 00 = nur weiterreichen</code>
</div>
<div class="card market-card">
<h3>Am Markt</h3>
<div class="pills">
<span class="pill brand">OpenTelemetry</span>
<span class="pill">Jaeger</span>
<span class="pill">Grafana Tempo</span>
<span class="pill">Zipkin</span>
<span class="pill">Datadog APM</span>
<span class="pill">Honeycomb</span>
</div>
</div>
</div>

</div>

</div>

Note:
- <strong>Trace:</strong> Sendungsverfolgung &uuml;ber mehrere Lieferdienste. Die Trace-ID bleibt &uuml;berall gleich.
- <strong>Span:</strong> Eigene ID plus Parent ergibt im SDK den Wasserfall. Wir bauen keinen Span-Baum.
- <strong>Weitergabe:</strong> Pro ausgehendem Aufruf eine frische Span-ID als parent-id in den Header. Der Empf&auml;nger sieht sie als Parent. Im Netzwerk-Jargon hei&szlig;t das oft &bdquo;Hop&ldquo;, ist aber kein Feld und kein Z&auml;hler.
- <strong>traceparent:</strong> parent-id ist die Span-ID dieses Aufrufs. Nur-Nullen ung&uuml;ltig, <code>crypto/rand</code>: IDs sollen unvorhersagbar sein.
- <strong>Sampling-Flag:</strong> Der Entry-Point entscheidet, alle anderen folgen, damit ein Trace ganz oder gar nicht aufgezeichnet wird. Gesampelt hei&szlig;t in der Stichprobe. Bei <code>00</code> misst niemand, der Header wird nur weitergereicht. Die Referenz setzt immer <code>01</code>.
- Am Markt: OpenTelemetry ist Standard, dazu Jaeger, Tempo, SaaS. Wir bauen die Propagation selbst, gut 120 Zeilen.
- &Uuml;berleitung: &bdquo;Vier Stellen im Code, mehr ist es nicht.&ldquo;
