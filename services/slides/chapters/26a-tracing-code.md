## Distributed Tracing

<p class="subtitle">W3C Trace Context in Pseudo-Code</p>

<pre class="cheatsheet"><span class="cmd">// W3C Trace Context: 00-&lt;trace-id-32hex&gt;-&lt;span-id-16hex&gt;-&lt;flags-2hex&gt;</span>

<span class="cmd">// Server-Middleware: eingehenden Header &uuml;bernehmen oder neu erzeugen</span>
on incomingRequest(req):
  tc = parse(req.header("traceparent")) ?: generateNew()
  res.header("traceparent") = tc.toHeader()    // zur&uuml;ck zum Client
  ctx = ctx.with(tc)
  next(req.withContext(ctx))

<span class="cmd">// Client-Inject: pro ausgehendem Hop neue Span-ID, gleiche Trace-ID</span>
on outgoingRequest(ctx, req):
  tc = ctx.get(TraceContext)
  hop = tc.copy(spanId = randomSpanId())
  req.setHeader("traceparent", hop.toHeader())

<span class="cmd">// Logging: jede Zeile bekommt trace_id als Feld</span>
log.info("forward step done",
  trace_id: ctx.tc.traceId, span_id: ctx.tc.spanId, step: "flight")

<span class="cmd">// Async-Grenze: Trace-ID als Event-Property mitschicken</span>
event = {
  eventId, sagaId, bookingId,
  traceparent: ctx.tc.toHeader()    // damit der Konsument den Trace fortf&uuml;hrt
}
publish(event)

on receiveCompensationEvent(event):
  tc = parse(event.traceparent) ?: generateNew()
  asyncCtx = bgCtx.with(tc)
  go process(asyncCtx, event)       // Logs der Goroutine tragen die trace_id
</pre>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 7 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- Vier Knackpunkte hervorheben:
  - <strong>parse() strikt halten</strong> &mdash; Format <code>version-trace-span-flags</code> mit festen L&auml;ngen, Nur-Nullen verwerfen, h&ouml;here Versionen <em>abw&auml;rtskompatibel ignorieren</em> (nur <code>00</code> verstehen, alles andere als ung&uuml;ltig behandeln).
  - <strong>Pro Hop neue Span-ID</strong>, aber <em>gleiche</em> Trace-ID &mdash; das ist der Trick, der den Vorgang als zusammenh&auml;ngende Kette identifizierbar macht.
  - <strong>Logger im Kontext</strong> &mdash; einmal an den Request-Logger h&auml;ngen, dann taucht <code>trace_id</code> ohne Format-String in jeder Zeile auf.
  - <strong>Async-Grenze</strong> &mdash; HTTP-Header geht beim &Uuml;bergang in eine Worker-Goroutine verloren, deshalb <code>traceparent</code> aktiv als Event-Property weitergeben.
- Reference-Code: <code>services/shared/tracing/tracing.go</code> &mdash; ca. 100 Zeilen Go mit strenger Format-Validierung.
- Diskussions-Anker: Was, wenn der Aufrufer einen <em>vergifteten</em> Header schickt? Antwort: strikt parsen, im Zweifel <code>generateNew()</code>. Niemals einen ung&uuml;ltigen Trace fortf&uuml;hren.
