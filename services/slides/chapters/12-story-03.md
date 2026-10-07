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

