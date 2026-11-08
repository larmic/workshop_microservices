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
<strong>Schon vorbereitet:</strong> <code>docker-compose.custom.yml</code> setzt <code>FLIGHT_SERVICE_URL</code>, <code>HOTEL_SERVICE_URL</code> und <code>CAR_SERVICE_URL</code>, ihr lest sie nur aus. Einklinken und Testen: Dashboard, Schritt 0 &bdquo;Start&ldquo;.
</div>

</div>

</div>

Note:
- Wiedererkennung: dieselbe Story (Kontext / User Story / Akzeptanzkriterien) findet ihr im Dashboard unter &bdquo;Aufgabe&ldquo;.
- Sprache und Framework sind frei (Go, Java, Quarkus, Node, &hellip;). Die Referenz unter `services/booking/story1/` ist nur ein Go-Beispiel.
- Kein Error-Handling in Story 1, das ist Absicht. Wir bauen das Skelett; Resilienz kommt in Story 4 und 5, Konsistenz in 6.
- Das Grundger&uuml;st (Mini-Service, Dockerfile, Port 8080) war Hausaufgabe. Einklinken und Testen steht im Dashboard unter Schritt 0 &bdquo;Start&ldquo;, dort landen alle beim ersten &Ouml;ffnen. Achtung: Ohne eigenen Pfad l&auml;uft die Kotlin-Beispiell&ouml;sung, dann ist auch alles gr&uuml;n. Erst <code>docker logs booking-custom</code> zeigt, wessen Service l&auml;uft.
- Ab Story 3 erweitern wir denselben Service, das Setup bleibt unver&auml;ndert. Time-Box 90 min, Spickzettel im Dashboard zugeklappt unter der Aufgabe.
- Vollst&auml;ndige Aufgabenbeschreibung: `docs/stories/story-01-cloud-native-booking-service.md`.
