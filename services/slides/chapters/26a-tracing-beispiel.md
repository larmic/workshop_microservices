## Distributed Tracing

<p class="subtitle">Eine Buchung durch unseren Stack</p>

<pre class="cheatsheet"><span class="cmd">Booking (Entry-Point) erzeugt:</span>
  T  = 0af7651916cd43dd8448eb211c80319c   &larr; Trace-ID &mdash; bleibt &uuml;berall gleich
  S0 = a1a1a1a1a1a1a1a1                   &larr; Booking-Span (im Booking-Log)

<span class="cmd">Outbound Booking &rarr; Flight (Client-Inject):</span>
  S1 = b7ad6b7169203331                   &larr; neu f&uuml;r diesen Hop, <em>nur im Header</em>
  traceparent: 00-0af7&hellip;319c-b7ad&hellip;3331-01

<span class="cmd">Flight (Propagate) &uuml;bernimmt:</span>
  ctx = { T, S1 }                          &larr; S1 wird Flights aktive Span

<span class="cmd">Falls Flight selbst weiterruft &mdash; Outbound zu externem Provider:</span>
  S2 = c3c3c3c3c3c3c3c3                    &larr; neu f&uuml;r diesen Hop
  traceparent: 00-0af7&hellip;319c-c3c3&hellip;c3c3-01

<span class="cmd">Logs (strukturiert, trace_id + span_id als Felder):</span>
  Booking : { trace_id: 0af7&hellip;319c, span_id: a1a1&hellip;a1a1, msg: "booking start" }
  Flight  : { trace_id: 0af7&hellip;319c, span_id: b7ad&hellip;3331, msg: "flight booked" }
  External: { trace_id: 0af7&hellip;319c, span_id: c3c3&hellip;c3c3, msg: "reservation" }
</pre>

<p class="quote"><code>grep 0af7&hellip;319c</code> &rarr; gesamter Vorgang. <code>grep b7ad&hellip;3331</code> &rarr; nur Flight.</p>

Note:
- <strong>Zentrale Aussage</strong>: <em>eine</em> Trace-ID bleibt durch die ganze Kette, <em>pro Hop</em> wird eine neue Span-ID gew&uuml;rfelt. Die Asymmetrie zwischen Entry-Point (Booking erzeugt) und Downstream (Flight/Hotel/Car &uuml;bernehmen) ist Absicht und spiegelt die reale Welt &mdash; Trace-Initiierung geh&ouml;rt an die Systemgrenze.
- <strong>Stolperstein erkl&auml;ren</strong>: S1 erscheint im Booking-Log <em>nicht</em>. Sie wird zwar von Booking erzeugt, aber nur als Outbound-Header gesetzt. Booking loggt weiter mit seinem initialen S0. Erst Flight sieht S1 in seinen Logs &mdash; weil seine Server-Middleware sie aus dem Header gezogen hat.
- <strong>Downstream-of-Downstream</strong>: Wenn Flight selbst weiterruft, wendet er <em>denselben</em> Client-Inject an wie Booking &mdash; neue Span-ID, gleiche Trace-ID. Wichtig: Flight darf <em>niemals</em> eine neue Trace-ID erzeugen, sonst zerf&auml;llt der Trace genau dort. Die strikte Trennung von <code>Middleware</code> (erzeugt, falls keiner kommt) und <code>Propagate</code> (&uuml;bernimmt nur) ist daf&uuml;r das Werkzeug.
- <strong>Was geloggt wird</strong>: nicht der ganze <code>traceparent</code>-Header, sondern <code>trace_id</code> und <code>span_id</code> als <em>getrennte Felder</em> im strukturierten Log. Das macht <code>jq</code>/<code>grep</code>-Filter trivial.
- &Uuml;berleitung: Jetzt zeigen wir den konkreten Pseudo-Code, der das umsetzt &mdash; Server-Middleware, Client-Inject, Logging, Async-Grenze.
