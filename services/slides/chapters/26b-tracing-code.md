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
- Ausf&uuml;hrlicher steht derselbe Pseudo-Code im Dashboard unter Story 8, &bdquo;Spickzettel&ldquo;. Wiedererkennung gewollt.
- <strong>1 Entry-Point:</strong> Nur Booking erzeugt einen Trace, falls keiner reinkommt. Flight, Hotel, Car &uuml;bernehmen einen vorhandenen Header und erzeugen <em>niemals</em> selbst einen. Zwei Middlewares in der Referenz: <code>Middleware</code> (erzeugt) und <code>Propagate</code> (&uuml;bernimmt nur). Sonst zerf&auml;llt der Trace genau an der Stelle, wo ein Downstream &bdquo;sicherheitshalber&ldquo; eine neue ID w&uuml;rfelt.
- <strong>2 pro Hop:</strong> Neue Span-ID, gleiche Trace-ID. Die neue Span erscheint im Booking-Log <em>nicht</em>, sie steht nur im Outbound-Header. Erst Flight sieht sie, weil seine Middleware sie aus dem Header zieht. Ein Beispiel-Trace durch den Stack: Booking loggt mit Span <code>a1a1&hellip;</code>, setzt f&uuml;r Flight <code>b7ad&hellip;</code> in den Header, Flight loggt mit <code>b7ad&hellip;</code>. Trace-ID &uuml;berall <code>0af7&hellip;319c</code>. <code>grep 0af7</code> gibt den Vorgang, <code>grep b7ad</code> nur den Flight-Anteil.
- <strong>3 jede Zeile:</strong> Nicht der ganze Header, sondern <code>trace_id</code> und <code>span_id</code> als getrennte Felder im strukturierten Log. Einmal an den Request-Logger h&auml;ngen, dann steht die ID ohne Format-String in jeder Zeile. Filter mit <code>jq</code> werden trivial.
- <strong>4 Async-Grenze:</strong> Der HTTP-Header geht beim &Uuml;bergang in die Worker-Goroutine verloren, deshalb wandert <code>traceparent</code> aktiv als Event-Property mit. Bei echten Brokern geh&ouml;rt er in die Message-Header, nicht in den Payload.
- <strong>parse() strikt halten:</strong> Feste L&auml;ngen pr&uuml;fen, Nur-Nullen verwerfen, unbekannte Versionen als ung&uuml;ltig behandeln. Ein vergifteter Header f&uuml;hrt zu <code>generateNew()</code>, nie zu einem fortgef&uuml;hrten Unsinn.
- Referenz: <code>services/shared/tracing/tracing.go</code>, etwa 100 Zeilen Go mit strenger Format-Validierung und Tests.
- &Uuml;berleitung: &bdquo;Jetzt zieht ihr den Faden selbst.&ldquo;
