<div class="page">

<p class="kicker">Story 8</p>

<div class="page-head">

## Den roten Faden im Log

<span class="badge">&asymp; 60 min</span>
</div>

<p class="subtitle">Eine Buchung, vier Logs. Danach reicht ein <span class="hl">grep</span>.</p>

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
<li>Jeder ausgehende Call reicht ihn weiter, Forward und Kompensation</li>
<li>Auch die Compensation-Events aus Story 7 tragen ihn, als Event-Property</li>
</ul>
</div>
<div class="story-list">
<h4>Die Logs</h4>
<ul>
<li>Jede Logzeile aller Services tr&auml;gt die <code>trace_id</code></li>
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
- Hook: &bdquo;Erinnert ihr euch an die Saga aus Story 6 und 7? Eine Buchung mit Kompensation, bis zu sechs HTTP-Calls plus Events. Wer wei&szlig; jetzt noch, welche Logzeile zu welchem Vorgang geh&ouml;rt?&ldquo; Dann die Demo: erst grep auf Story 7 (frustrierend), dann auf Story 8 (eine Trace-ID, alles da).
- Wiedererkennung: dieselbe Story (Kontext, User Story, Akzeptanzkriterien) im Dashboard unter Story 8, &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story8/</code> und <code>services/shared/tracing/</code>.
- Wichtig: Booking ist der Entry-Point und erzeugt einen Trace, falls keiner reinkommt. Flight, Hotel, Car sind passive Empf&auml;nger: Sie verl&auml;ngern den eintreffenden Header, erzeugen aber nie selbst einen. In Stories 1 bis 7 bleiben ihre Logs deshalb ohne <code>trace_id</code>, der Kontrast macht den Effekt sichtbar.
- Demo-Drehbuch: <code>POST /booking/bookings</code> mit Hotel auf &bdquo;Fehler&ldquo;. Die Response tr&auml;gt den <code>traceparent</code> zur&uuml;ck. Diese ID per <code>docker compose logs | jq -c 'select(.trace_id=="...")'</code> &uuml;ber alle Services filtern: Forward und Kompensation in einem Block.
- Bonus-Optionen nennen: Jaeger-Container im Compose (sch&ouml;ne UI, aber Sprach- und Tool-abh&auml;ngig), Trace-ID in der Response an den Kunden (Support kann sie nutzen), Sampling-Strategie diskutieren.
- Time-Box 60 min inklusive Demo. Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-08-tracing.md</code>.
