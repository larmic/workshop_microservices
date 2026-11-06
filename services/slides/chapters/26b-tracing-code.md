<div class="page">

<p class="kicker">Distributed Tracing</p>

## Vier Stellen im Code

<p class="subtitle">In Pseudo-Code, wie im Spickzettel.</p>

<div class="page-body codebody">

<div class="codeblock">
<div class="k">on incomingRequest(req):                       <span class="dim">// 1 Entry-Point: Booking</span></div>
<div>  tc  = parse(req.header("traceparent")) ?: <span class="hi">generateNew()</span></div>
<div>  ctx = ctx.with(tc)</div>
<div class="gap k">on outgoingRequest(ctx, req):                  <span class="dim">// 2 pro Hop</span></div>
<div>  hop = ctx.tc.copy(spanId = <span class="hi">randomSpanId()</span>)  <span class="dim">// neue Span, gleiche Trace</span></div>
<div>  req.setHeader("traceparent", hop.toHeader())</div>
<div class="gap k">log.info("flight booked", <span class="hi">trace_id</span>: ctx.tc.traceId)   <span class="dim">// 3 jede Zeile</span></div>
<div class="gap k">event = { eventId, sagaId, <span class="hi">traceparent</span>: ctx.tc.toHeader() }   <span class="dim">// 4 Async-Grenze</span></div>
<div>on receiveEvent(event):</div>
<div>  ctx = bgCtx.with(parse(event.traceparent) ?: generateNew())</div>
</div>

<p class="codenote">Erzeugen, weiterreichen, loggen, &uuml;ber den Bus retten. Hundert Zeilen, keine Library.</p>

</div>

</div>

Note:
- Derselbe Pseudo-Code steht im Dashboard unter Story 8, &bdquo;Spickzettel&ldquo;.
- <strong>1 Entry-Point:</strong> Booking erzeugt einen Trace, falls keiner kommt (<code>Middleware</code>). Flight, Hotel, Car &uuml;bernehmen per HTTP nur (<code>Propagate</code>), Workshop-Entscheidung. Tabu: eine eingehende ID &uuml;berschreiben.
- <strong>2 pro Hop:</strong> Neue Span-ID, gleiche Trace-ID. Die neue Span steht nur im Outbound-Header, erst Flight loggt sie. <code>grep 0af7</code> gibt den Vorgang, <code>grep b7ad</code> nur Flight.
- <strong>3 jede Zeile:</strong> <code>trace_id</code> und <code>span_id</code> als eigene Felder, einmal an den Request-Logger h&auml;ngen.
- <strong>4 Async-Grenze:</strong> <code>traceparent</code> wandert als Event-Property mit. Bei echten Brokern in die Message-Header.
- <code>parse()</code> strikt: feste L&auml;ngen, nicht nur Nullen, sonst <code>generateNew()</code>.
- Referenz: <code>services/shared/tracing/tracing.go</code>, ~100 Zeilen.
- &Uuml;berleitung: &bdquo;Jetzt zieht ihr den Faden selbst.&ldquo;
