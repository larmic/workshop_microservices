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
<h3>Hop</h3>
<p>Der &Uuml;bergang von Service zu Service. Neue Span-ID, gleiche Trace-ID.</p>
<code>Booking &rarr; Flight = 1 Hop</code>
</div>
<div class="card">
<h3>traceparent</h3>
<p>Der W3C-Header. Sprach- und vendor-neutral, immer 55 Zeichen.</p>
<code>00-&lt;trace 32&gt;-&lt;span 16&gt;-&lt;flags 2&gt;</code>
</div>
<div class="card">
<h3>Sampling-Flag</h3>
<p>Bit 0 in den Flags: Soll dieser Trace exportiert werden? Alles ist zu teuer.</p>
<code>flags = 01 &harr; sampled</code>
</div>
</div>

<div class="market">
<h4>Am Markt</h4>
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

Note:
- <strong>Trace:</strong> Analogie Sendungsverfolgung &uuml;ber mehrere Lieferdienste. Ohne gemeinsame Sendungsnummer wei&szlig; am Ende niemand, wo das Paket steckt. Die Trace-ID bleibt &uuml;ber alle Hops gleich, das ist der ganze Punkt.
- <strong>Span:</strong> Jede Span hat eine eigene ID und eine Parent-Span-ID. Daraus entsteht ein Baum: oben die eingehende Anfrage am Booking-Service, darunter pro Outbound-Call eine Span. Jaeger, Tempo, Datadog zeigen das als Wasserfall mit Dauer pro Span. Im Workshop bauen wir bewusst nur die Trace-Korrelation, keinen Span-Baum.
- <strong>Hop:</strong> Pro Hop wird eine frische Span-ID gew&uuml;rfelt, die alte wandert als Parent in den Header. So entstehen pro Buchung vier bis sechs Spans: eingehende Anfrage, drei Outbound-Calls, gegebenenfalls Kompensation.
- <strong>traceparent:</strong> Vier Felder mit festen L&auml;ngen: version (derzeit immer <code>00</code>), trace-id, parent-id (die Span-ID des Aufrufers), flags. Stolperstein: Nur-Nullen sind ung&uuml;ltig, <code>math/rand</code> mit Default-Seed reicht nicht, immer <code>crypto/rand</code>.
- <strong>Sampling-Flag:</strong> In Produktion sind Traces teuer (Storage, Netz, Lizenzen), niemand exportiert 100 Prozent. Der Entry-Point entscheidet, alle nachgelagerten Services folgen dem Bit. Strategien (Head-based, Tail-based, Adaptive) als Reserve im Recap.
- Am Markt: OpenTelemetry ist der De-facto-Standard mit vendor-neutralem OTLP-Export. Jaeger klassisch lokal und CNCF, Storage skaliert nicht trivial. Grafana Tempo g&uuml;nstig auf Objekt-Storage, integriert mit Loki und Mimir. Zipkin der Urahn. Datadog, Honeycomb, Dynatrace als SaaS: komfortabel, Pricing-Thema.
- Wichtig: Wir bauen die Propagation selbst, etwa 100 Zeilen, kein OpenTelemetry. Der Mechanismus ist sprach- und vendor-neutral, die OTel-SDKs sind es nicht. In Produktion nat&uuml;rlich umgekehrt.
- &Uuml;berleitung: &bdquo;Vier Stellen im Code, mehr ist es nicht.&ldquo;
