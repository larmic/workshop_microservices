<div class="page">

<p class="kicker">Story 1</p>

<div class="page-head">

## Erst lauff&auml;hig, dann smart

<span class="badge">&asymp; 90 min</span>
</div>

<p class="subtitle">Laptop, Cloud, Raspberry Pi. Ein Build, <span class="hl">&uuml;berall gleich</span>.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Betriebsteam</dd>
<dt>m&ouml;chte ich</dt>
<dd>einen Service, der in beliebigen Umgebungen ohne Code-&Auml;nderung l&auml;uft und sich automatisch &uuml;berwachen l&auml;sst,</dd>
<dt>damit</dt>
<dd>Probleme fr&uuml;h erkannt werden und Deployment-Pfade nicht spezifisch sein m&uuml;ssen.</dd>
</dl>
<div class="story-list">
<h4>Der Service kann</h4>
<ul>
<li><code>GET /health</code> liefert 200, wenn der Service gesund ist</li>
<li><code>GET /info</code> zeigt die aktuelle Konfiguration</li>
<li><code>GET /openapi</code> liefert die OpenAPI-Spec aus, <a href="https://github.com/larmic/workshop_microservices/blob/main/services/booking/story1/api/openapi.yaml" target="_blank" rel="noopener">Vorlage</a></li>
<li><code>GET /booking/offers</code> aggregiert Flight, Hotel und Car</li>
</ul>
</div>
<div class="story-list">
<h4>Betrieb</h4>
<ul>
<li>Logs gehen auf <code>stdout</code></li>
<li>Der Service liest die Backend-URLs aus Umgebungsvariablen</li>
</ul>
<h4>Verpackung</h4>
<ul>
<li><code>Dockerfile</code> im Service-Verzeichnis, Port <code>8080</code></li>
<li><code>CUSTOM_BOOKING_PATH</code> in <code>services/.env</code> zeigt darauf, damit <code>make docker-up-hub</code> mitbaut</li>
</ul>
</div>
</div>

<div class="story-foot">
<strong>Einmalig, und nicht eure Aufgabe:</strong> Die Werte der Backend-URLs liefert <code>docker-compose.custom.yml</code>, ihr lest sie nur, ihr setzt sie nicht. Details in <code>docs/custom-setup.md</code>.
</div>

</div>

</div>

Note:
- Hook zum Einsteigen vorlesen: &bdquo;Alle dachten, der Service l&auml;uft, bis ein Kunde anrief und fragte, warum seit drei Stunden nichts mehr geht. Vertrauen ist gut, ein Health-Endpoint ist besser.&ldquo;
- Wiedererkennung: dieselbe Story (Kontext / User Story / Akzeptanzkriterien) findet ihr im Dashboard unter &bdquo;Story lesen&ldquo;.
- Sprache und Framework sind frei (Go, Java, Quarkus, Node, &hellip;). Die Referenz unter `services/booking/story1/` ist nur ein Go-Beispiel.
- Kein Error-Handling in Story 1, das ist Absicht. Wir bauen das Skelett; Resilienz kommt in Stories 4 bis 6.
- Der Setup-Hinweis unten steht absichtlich nur hier in Story 1. Ab Story 3 erweitern wir denselben Service iterativ, das Custom-Setup bleibt unver&auml;ndert.
- Time-Box 90 min inkl. Setup. Dashboard `http://localhost` zeigt den Story-1-Modus inkl. Spickzettel mit Pseudo-Code.
- Vollst&auml;ndige Aufgabenbeschreibung: `docs/stories/story-01-cloud-native-booking-service.md`.
