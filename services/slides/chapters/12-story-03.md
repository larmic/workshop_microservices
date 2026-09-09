<div class="page">

<p class="kicker">Story 3</p>

<div class="page-head">

## Services finden sich selbst

<span class="badge">&asymp; 60 min</span>
</div>

<p class="subtitle">Consul kennt die Adressen. Der Code kennt nur <span class="hl">Namen</span>.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Betriebsteam</dd>
<dt>m&ouml;chte ich</dt>
<dd>Backend-Services &uuml;ber logische Namen finden, statt URLs per Hand zu pflegen,</dd>
<dt>damit</dt>
<dd>Skalierung und Ausf&auml;lle ohne n&auml;chtliche Config-Edits funktionieren.</dd>
</dl>
<div class="story-list">
<h4>Der Service kann</h4>
<ul>
<li>Ein Resolver fragt <code>GET /v1/health/service/{name}?passing=true</code> bei Consul</li>
<li><code>GET /booking/offers</code> l&ouml;st Flight, Hotel und Car &uuml;ber den Resolver auf, keine statischen URLs mehr</li>
<li><code>POST /booking/bookings</code> orchestriert die drei Backends und antwortet mit <code>201 Created</code></li>
</ul>
</div>
<div class="story-list">
<h4>Namen statt Adressen</h4>
<ul>
<li><code>flight-service</code>, <code>hotel-service</code>, <code>car-service</code></li>
<li>Die Backends registrieren sich beim Start selbst bei Consul</li>
</ul>
<h4>Optional</h4>
<ul>
<li>Client-Side Load Balancing: bei mehreren Instanzen zuf&auml;llig eine w&auml;hlen</li>
<li>Consul Key-Value Store f&uuml;r zus&auml;tzliche Konfiguration</li>
</ul>
</div>
</div>

<div class="story-foot">
<strong>Aus Story 2 mitgenommen:</strong> Ressource <code>bookings</code> statt Verb, <code>POST</code> antwortet mit <code>201 Created</code>, Fehler kommen als Status-Code, nicht als <code>200</code> mit Fehler-Body.
</div>

</div>

</div>

Note:
- Hook: &bdquo;Story 1 hatte URLs in ENV-Variablen. Was passiert, wenn Flight umzieht? Re-Deploy. Was passiert, wenn ihr Flight skaliert? Ein Backend bekommt allen Traffic.&ldquo;
- Wiedererkennung: dieselbe Story (Kontext / User Story / Akzeptanzkriterien) im Dashboard unter Story 3, &bdquo;Story lesen&ldquo;.
- Br&uuml;cke zu Story 2: <code>POST /booking/bookings</code> ist der erste schreibende Endpoint. Die Regeln vom Flipchart gelten, siehe Fu&szlig;zeile. Kurz fragen, wer sie schon eingebaut hat.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story3/</code> (Go).
- Self-Registration im Code: Flight, Hotel und Car melden sich beim Start aktiv an. Trade-off zu Sidecar oder plattformbasierter Registrierung diskutieren wir im Recap.
- Client-Side Load Balancing: Der Resolver bekommt eine Liste und w&auml;hlt eine Instanz zuf&auml;llig. Keine Load-Balancer-Magie n&ouml;tig.
- Time-Box 60 min. Dashboard <code>http://localhost</code> zeigt den Story-3-Modus inklusive Spickzettel.
- Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-03-service-discovery.md</code>.
