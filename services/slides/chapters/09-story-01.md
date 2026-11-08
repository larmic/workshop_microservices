<div class="page">

<p class="kicker">Story 1 &middot; Hands-on</p>

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
<li><code>GET /info</code> zeigt die aktuelle Konfiguration</li>
<li><code>GET /openapi</code> liefert die OpenAPI-Spec aus, <a href="https://github.com/larmic/workshop_microservices/blob/main/services/booking/story1/api/openapi.yaml" target="_blank" rel="noopener">Vorlage</a></li>
<li><code>GET /booking/offers</code> aggregiert Flight, Hotel und Car</li>
</ul>
</div>
<div class="story-list">
<h4>Betrieb</h4>
<ul>
<li>Logs gehen auf <code>stdout</code></li>
<li>Backend-URLs aus Umgebungsvariablen: <code>FLIGHT_SERVICE_URL</code>, <code>HOTEL_SERVICE_URL</code>, <code>CAR_SERVICE_URL</code></li>
</ul>
</div>
</div>

<ol class="story-path" aria-label="Ablauf">
<li class="done"><span class="sp-step">Hausaufgabe</span><span class="sp-what"><code>/health</code> &middot; Dockerfile &middot; Port 8080</span></li>
<li class="done"><span class="sp-step">Start (0)</span><span class="sp-what">eingeklinkt und getestet</span></li>
<li class="now"><span class="sp-step">Story 1</span><span class="sp-what">Logs, Env-Vars, <code>/info</code>, <code>/openapi</code>, Offers</span></li>
</ol>

</div>

</div>

Note:
- Wiedererkennung: dieselbe Story (Kontext / User Story / Akzeptanzkriterien) findet ihr im Dashboard unter &bdquo;Aufgabe&ldquo;.
- Sprache und Framework sind frei (Go, Java, Quarkus, Node, &hellip;). Die Referenz unter `services/booking/story1/` ist nur ein Go-Beispiel.
- Kein Error-Handling in Story 1, das ist Absicht. Wir bauen das Skelett; Resilienz kommt in Story 4 und 5, Konsistenz in 6.
- Achtung beim Einklinken: Ohne eigenen Pfad l&auml;uft die Kotlin-Beispiell&ouml;sung, dann ist auch alles gr&uuml;n. Sie meldet sich in <code>docker logs booking-custom</code> mit einem gro&szlig;en Kasten &bdquo;BEISPIEL-L&Ouml;SUNG&ldquo;.
- Ab Story 3 erweitern wir denselben Service, das Setup bleibt unver&auml;ndert. Time-Box 90 min, Spickzettel im Dashboard zugeklappt unter der Aufgabe.
- Vollst&auml;ndige Aufgabenbeschreibung: `docs/stories/story-01-cloud-native-booking-service.md`.
