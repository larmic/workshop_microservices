<div class="page">

<p class="kicker">Story 2 &middot; Design-Session</p>

<div class="page-head">

## Erst denken, dann tippen

<span class="badge design">&asymp; 55 min <small>davon 20 Flipchart</small></span>
</div>

<p class="subtitle">Storno und Umbuchung am Flipchart. <span class="hl">Kein Code</span>, kein Deployment, kein Port.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Produktverantwortliche</dd>
<dt>m&ouml;chte ich</dt>
<dd>eine API f&uuml;r Storno und Umbuchung, die sich an die HTTP-Regeln h&auml;lt,</dd>
<dt>damit</dt>
<dd>Wiederholungen nichts doppelt buchen, Fehler sofort auffallen und die API lange h&auml;lt.</dd>
</dl>
<div class="story-list">
<h4>Aufgabe</h4>
<ul>
<li>Teams zu 3 bis 4 Personen, ein Flipchart pro Team</li>
<li>Eine Kundin storniert die ganze Reise, der Grund bleibt nachvollziehbar</li>
<li>Ein Kunde bucht nur den Flug um, Hotel und Mietwagen bleiben</li>
<li>Der Client wiederholt Anfragen bei Timeout</li>
</ul>
</div>
<div class="story-list">
<h4>Ergebnis</h4>
<ul>
<li>Eine Tabelle, eine Zeile pro Endpoint</li>
<li>Methode und Pfad, Request-Kern, Antwort und Status-Code</li>
<li>Idempotent: ja oder nein</li>
<li>Begr&uuml;ndung, auch f&uuml;r bewusste Abweichungen</li>
</ul>
</div>
</div>

</div>

</div>

Note:
- Beim Herumgehen fragen: Was antwortet eure API, wenn ein Backend beim Storno oder Umbuchen scheitert, etwa das Hotel ablehnt oder gar nicht antwortet? Unerwartete 5xx muss man nicht einzeln dokumentieren, einen bewusst gelieferten 503 schon. So macht es unsere Referenz ab Story 4 in der <code>openapi.yaml</code>. Darauf achten, dass niemand ein 200 mit Fehler-Body in die Tabelle schreibt.
