<div class="page">

<p class="kicker">Story 2</p>

<div class="page-head">

## Erst denken, dann tippen

<span class="badge">20 min Teamarbeit</span>
</div>

<p class="subtitle">Storno und Umbuchung am Flipchart. Kein Code, kein Deployment, <span class="hl">kein Port</span>.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Produktverantwortliche</dd>
<dt>m&ouml;chte ich</dt>
<dd>eine API f&uuml;r Storno und Umbuchung, die HTTP so nutzt, wie es gemeint ist,</dd>
<dt>damit</dt>
<dd>Retries sicher sind, Fehler sichtbar werden und wir sie nicht in drei Monaten umbauen.</dd>
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
- Beim Herumgehen fragen: Was antwortet eure API, wenn ein Backend beim Storno oder Umbuchen scheitert, etwa das Hotel ablehnt oder gar nicht antwortet? Ein 5xx kann immer passieren und geh&ouml;rt nicht in die API-Doku. Darauf achten, dass niemand ein 200 mit Fehler-Body in die Tabelle schreibt.
