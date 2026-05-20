## Story 7

<p class="subtitle">Den roten Faden im Log <span class="time-badge">&asymp; 60 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

Im aktuellen Setup zeigen die Service-Logs jede Anfrage isoliert: Booking loggt seine Saga, Flight/Hotel/Car loggen ihre lokalen Buchungen &mdash; aber niemand kann eine Buchung als <em>eine</em> Operation &uuml;ber alle Services hinweg verfolgen. Korrelation passiert &uuml;ber Timestamps und Augenma&szlig;.

Eine <strong>Trace-ID</strong> pro eingehender Anfrage, &uuml;ber alle nachgelagerten Aufrufe propagiert und in jede Logzeile geschrieben, l&ouml;st das Problem auch ohne ein Tracing-Backend: <code>grep &lt;trace-id&gt;</code> zeigt den vollst&auml;ndigen Ablauf.

#### User Story

Als <em>Entwickler:in im Betrieb</em> m&ouml;chte ich <em>einen einzelnen Gesch&auml;ftsvorgang &uuml;ber Service-Grenzen hinweg in den Logs verfolgen k&ouml;nnen</em>, damit <em>ich Fehlerursachen, Latenzen und Saga-Verl&auml;ufe gezielt analysieren kann, ohne Logs zusammenzupuzzeln</em>.

</div>

</div>
<div>

<div class="story-card">

#### Akzeptanzkriterien

- Booking erzeugt oder &uuml;bernimmt bei jedem Request einen W3C-<code>traceparent</code>
- Trace-Kontext wird auf <em>jeden</em> ausgehenden HTTP-Call weitergereicht
- Trace-Kontext wandert auch &uuml;ber die <strong>Compensation-Events</strong> aus Story 6 (als Event-Property)
- Jede Logzeile aller Services tr&auml;gt die <code>trace_id</code>
- Logs sind strukturiert (JSON) &mdash; <code>trace_id</code> ist ein eigenes Feld
- OpenAPI dokumentiert den <code>traceparent</code>-Header

</div>

</div>
</div>

Note:
- Hook: &bdquo;Erinnert ihr euch an die Saga aus Story 5/6? Eine Buchung mit Kompensation: bis zu sechs HTTP-Calls plus Events. Wer wei&szlig; jetzt noch, welche Logzeile zu welchem Vorgang geh&ouml;rt? Genau &mdash; niemand.&ldquo; Dann die Demo: erst grep auf Story 5/6 (frustrierend), dann auf Story 7 (eine Trace-ID, alles da).
- Wiedererkennung: dieselbe Karte (Kontext / User Story / Akzeptanzkriterien) im Dashboard unter Story 7 &rarr; &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story7/</code> + <code>services/shared/tracing/</code>.
- Wichtig: <strong>Booking ist der Entry-Point</strong>, erzeugt einen Trace, falls keiner reinkommt. Flight/Hotel/Car sind <em>passive Empf&auml;nger</em> &mdash; sie verl&auml;ngern den eintreffenden Header, erzeugen aber niemals selbst. Spiegelt die reale Welt: Trace-Initiierung am Entry-Point (API-Gateway, Public-Service), nicht in jedem Downstream-Hop.
- Demo-Drehbuch: <code>POST /booking/bookings</code> mit Hotel auf &bdquo;Fehler&ldquo; &rarr; Dashboard zeigt die Trace-ID nach jedem Call. Diese ID per <code>docker compose logs | jq -c 'select(.trace_id=="...")'</code> &uuml;ber alle Services filtern &rarr; Forward + Kompensation in einem zusammenh&auml;ngenden Block.
- Bonus-Optionen explizit nennen: Jaeger-Container im Compose erg&auml;nzen (sch&ouml;ne UI, aber Sprach-/Tool-abh&auml;ngig), Sampling-Strategie diskutieren, Trace-ID in der Response an den Kunden zur&uuml;ckgeben (Support kann sie nutzen).
- Time-Box 60 min. Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-07-tracing.md</code>.
