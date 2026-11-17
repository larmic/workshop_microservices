<div class="page">

<p class="kicker">Distributed Tracing</p>

## Vier Stellen im Code

<p class="subtitle">In Pseudo-Code, wie im Spickzettel.</p>

<div class="page-body codebody">

<div class="codeblock">
<div class="k">on incomingRequest(req):                       <span class="dim">// 1 Entry-Point: Booking</span></div>
<div>  tc  = parse(req.header("traceparent")) ?: <span class="hi">generateNew()</span></div>
<div>  res.setHeader("traceparent", tc.toHeader())  <span class="dim">// zur&uuml;ck zum Client</span></div>
<div>  ctx = ctx.with(tc)</div>
<div class="dim">  // traceparent: 00-0af7651916cd43dd8448eb211c80319c-00f067aa0ba902b7-01</div>
<div class="gap k">on outgoingRequest(ctx, req):                  <span class="dim">// 2 pro Hop</span></div>
<div>  hop = ctx.tc.copy(spanId = <span class="hi">randomSpanId()</span>)  <span class="dim">// neue Span, gleiche Trace</span></div>
<div>  req.setHeader("traceparent", hop.toHeader())</div>
<div class="dim">  // traceparent an Flight: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01</div>
<div class="gap k">log.info("flight booked", <span class="hi">trace_id</span>: ctx.tc.traceId)   <span class="dim">// 3 jede Zeile</span></div>
<div class="dim">  // {"msg":"flight booked","trace_id":"0af7651916cd43dd8448eb211c80319c"}</div>
<div class="gap k">event = { eventId, sagaId, <span class="hi">traceparent</span>: ctx.tc.toHeader() }   <span class="dim">// 4 Async-Grenze (Bonus)</span></div>
<div>on receiveEvent(event):</div>
<div>  ctx = bgCtx.with(parse(event.traceparent) ?: generateNew())</div>
</div>

<p class="codenote">Erzeugen, weiterreichen, loggen, &uuml;ber den Bus retten. Gut 120 Zeilen, keine Library.</p>

</div>

</div>

Note:
- Derselbe Ablauf steht im Dashboard unter Story 8, &bdquo;Spickzettel&ldquo;. Dort ist die Kompensation (<code>DELETE</code>) ein eigener Schritt, hier steckt sie in Stelle 2.
- <strong>1 Entry-Point:</strong> Booking erzeugt einen Trace, falls keiner kommt (<code>Middleware</code>). Flight, Hotel, Car &uuml;bernehmen per HTTP nur (<code>Propagate</code>), Workshop-Entscheidung. Tabu: eine eingehende ID &uuml;berschreiben.
- <strong>Beispielwerte:</strong> Trace-ID <code>0af7&hellip;319c</code> bleibt &uuml;berall gleich. Booking hat die Span <code>00f0&hellip;02b7</code>, der Aufruf an Flight bekommt <code>b7ad&hellip;3331</code>.
- <strong>2 pro Hop:</strong> Neue Span-ID, gleiche Trace-ID. Die neue Span steht nur im Outbound-Header, erst Flight loggt sie. <code>grep 0af7</code> gibt den Vorgang, <code>grep b7ad</code> nur Flight.
- <strong>3 jede Zeile:</strong> <code>trace_id</code> und <code>span_id</code> als eigene Felder, einmal an den Request-Logger h&auml;ngen.
- <strong>4 Async-Grenze:</strong> <code>traceparent</code> wandert als Event-Property mit. Bei echten Brokern in die Message-Header.
- <code>parse()</code> strikt: feste L&auml;ngen, nicht nur Nullen, sonst <code>generateNew()</code>.
- Referenz: <code>services/shared/tracing/tracing.go</code>, gut 120 Zeilen.
- &Uuml;berleitung: &bdquo;Jetzt zieht ihr den Faden selbst.&ldquo;
