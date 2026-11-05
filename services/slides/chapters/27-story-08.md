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
<li>Jeder ausgehende Call reicht ihn weiter, Forward und DELETE-Kompensation</li>
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
- Hook: &bdquo;Erinnert ihr euch an eure Saga aus Story 6? Eine Buchung, bis zu f&uuml;nf HTTP-Calls &uuml;ber drei Backends. Wer wei&szlig; noch, welche Logzeile zu welchem Vorgang geh&ouml;rt?&ldquo; Dann die Demo: erst grep auf Story 7 (frustrierend), dann auf Story 8 (eine Trace-ID, alles da).
- Ausgangspunkt ist der eigene Story-6-Service. Events, CB und Bulkhead braucht niemand. Wer mit Story 6 nicht fertig ist: Kriterien f&uuml;r die Calls, die da sind, Kompensation zeigst du an der Referenz.
- Wiedererkennung: dieselbe Story (Kontext, User Story, Akzeptanzkriterien) im Dashboard unter Story 8, &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story8/</code> und <code>services/shared/tracing/</code>. Sie baut auf Story 7 plus CB und Bulkhead auf und zeigt mehr, als die Gruppen brauchen.
- Wichtig: Booking ist der Entry-Point und erzeugt einen Trace, falls keiner reinkommt. Flight, Hotel, Car sind passive Empf&auml;nger: Sie verl&auml;ngern den eintreffenden Header, erzeugen aber nie selbst einen. In Stories 1 bis 7 bleiben ihre Logs deshalb ohne <code>trace_id</code>, der Kontrast macht den Effekt sichtbar. Einzige Ausnahme: Ein Compensation-Event ohne Property startet einen neuen Trace (Bonus).
- Demo-Drehbuch: im Dashboard &bdquo;Buchung mit kaputtem Mietwagen&ldquo; (<code>carId: C9</code>, Car antwortet 409). Die Response tr&auml;gt den <code>traceparent</code> zur&uuml;ck, das Panel zeigt den fertigen Filter: <code>docker compose logs --no-log-prefix &hellip; | jq -cR 'fromjson? | select(.trace_id=="...")'</code>. Forward und Kompensation in einem Block.
- Bonus-Optionen nennen: Detektiv-Runde im Dashboard (Rush Hour, eine von f&uuml;nf Buchungen scheitert, wer findet die Trace-ID zuerst?), Kompensation als Event mit <code>traceparent</code> als Property, Jaeger-Container im Compose.
- Time-Box 60 min inklusive Demo. Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-08-tracing.md</code>.
