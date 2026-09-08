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
- Hook: &bdquo;Beim letzten Mal hat ein Prefetcher &uuml;ber einen GET-Link Buchungen storniert. Das soll uns nicht noch einmal passieren.&ldquo;
- Wiedererkennung: dieselbe Story im Dashboard unter Story 2, &bdquo;Story lesen&ldquo;. Dort gibt es bewusst keinen Spickzettel, keine Buttons und keinen Service.
- Time-Box: 3 Minuten Aufgabe stellen, 20 Minuten Teamarbeit (das Badge oben), dann je Team 2 Minuten Vorstellung, anschlie&szlig;end das Quiz. Der ganze Block dauert etwa 55 Minuten.
- Spaltenk&ouml;pfe der Tabelle vorab auf die Flipcharts zeichnen, das spart f&uuml;nf Minuten.
- W&auml;hrend der Teamarbeit herumgehen und nur Fragen stellen: &bdquo;Was passiert beim zweiten Aufruf?&ldquo;, &bdquo;Welcher Code, wenn das Hotel ablehnt?&ldquo;
- Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-02-api-design-session.md</code>. Trainer-Hinweis mit L&ouml;sungsraum: <code>docs/instructions/rest-vs-restful.md</code>.
