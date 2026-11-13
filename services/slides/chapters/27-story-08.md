<div class="page">

<p class="kicker">Story 8 &middot; Hands-on</p>

<div class="page-head">

## Den roten Faden im Log

<span class="badge">&asymp; 60 min</span>
</div>

<p class="subtitle">Euer Service aus Story 6, eine Buchung, vier Logs. Danach reicht ein <span class="hl">grep</span>.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Entwickler:in im Betrieb</dd>
<dt>m&ouml;chte ich</dt>
<dd>einen einzelnen Gesch&auml;ftsvorgang &uuml;ber Service-Grenzen hinweg in den Logs verfolgen,</dd>
<dt>damit</dt>
<dd>ich Fehler, Latenzen und Saga-Verl&auml;ufe analysieren kann, ohne Zeitstempel zu puzzeln.</dd>
</dl>
<div class="story-list">
<h4>Der Trace</h4>
<ul>
<li>Booking erzeugt oder &uuml;bernimmt bei jedem Request einen W3C-<code>traceparent</code></li>
<li>Jeder ausgehende Call reicht ihn weiter, Buchung und DELETE-Kompensation</li>
<li>Die Response gibt den <code>traceparent</code> zur&uuml;ck</li>
</ul>
</div>
<div class="story-list">
<h4>Die Logs</h4>
<ul>
<li>Jede Logzeile eures Booking-Service tr&auml;gt die <code>trace_id</code>, Flight, Hotel und Car k&ouml;nnen das schon</li>
<li>Strukturiert als JSON, <code>trace_id</code> ist ein eigenes Feld</li>
</ul>
<h4>Sichtbarkeit</h4>
<ul>
<li>OpenAPI dokumentiert den <code>traceparent</code>-Header</li>
<li><code>docker compose logs | grep &lt;trace-id&gt;</code> zeigt den ganzen Vorgang</li>
</ul>
</div>
</div>

</div>

</div>

Note:
- Hook: &bdquo;Eure Saga aus Story 6: eine Buchung, bis zu f&uuml;nf HTTP-Calls. Welche Logzeile geh&ouml;rt zu welchem Vorgang?&ldquo; Demo: grep auf Story 7, dann auf Story 8.
- Gegeben: Flight, Hotel, Car loggen eine eingehende <code>trace_id</code> schon. Aufgabe: der eigene Story-6-Service, ohne Events, CB, Bulkhead. Nicht fertig? Kompensation an der Referenz zeigen. Die storniert per Event wie in Story 7, der <code>traceparent</code> wandert als Property mit, die Trace-ID bleibt trotzdem dieselbe.
- Downstreams &uuml;bernehmen, starten per HTTP bewusst keinen, Tabu ist &uuml;berschreiben.
- Kompensation ausl&ouml;sen: &bdquo;Buchung mit kaputtem Mietwagen&ldquo; (<code>C9</code>, 409). Panel zeigt den <code>jq</code>-Filter. Nicht f&uuml;nfmal schnell klicken, sonst &ouml;ffnet der Car-Breaker (503).
- Bonus: Detektiv-Runde (Rush Hour, wer findet die kaputte Buchung?), Kompensation als Event, Jaeger.
- Referenz: <code>services/booking/story8/</code>, <code>services/shared/tracing/</code>.
- Time-Box 60 min inkl. Demo. Details: <code>docs/stories/story-08-tracing.md</code>.
