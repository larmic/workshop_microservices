<div class="page">

<p class="kicker">Story 4</p>

<div class="page-head">

## Wenn der Flug ausf&auml;llt

<span class="badge">&asymp; 60 min</span>
</div>

<p class="subtitle">Flight ist weg. Hotel und Mietwagen buchen wir <span class="hl">trotzdem</span>.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Kunde</dd>
<dt>m&ouml;chte ich</dt>
<dd>auch bei Ausfall des Flugbuchungssystems Hotel und Mietwagen buchen k&ouml;nnen,</dd>
<dt>damit</dt>
<dd>meine Reiseplanung nicht komplett blockiert wird.</dd>
</dl>
<div class="story-list">
<h4>Der Breaker</h4>
<ul>
<li>Je ein Circuit Breaker um jeden Backend-Aufruf: Flight, Hotel, Car</li>
<li>Nach <code>5</code> aufeinanderfolgenden Fehlern &ouml;ffnet er</li>
<li>Nach 30 Sekunden wechselt er in <code>HALF_OPEN</code>, ein Probe-Call geht durch</li>
<li>Aufruf-Timeout konfiguriert, maximal drei Sekunden</li>
</ul>
</div>
<div class="story-list">
<h4>Der Fallback</h4>
<ul>
<li>Bei offenem Circuit l&auml;uft ein Fallback, etwa <code>flights: []</code></li>
<li>Der Header <code>X-Circuit-Open</code> macht es dem Client sichtbar</li>
</ul>
<h4>Sichtbarkeit</h4>
<ul>
<li>Circuit-Status &uuml;ber einen Admin-Endpoint abfragbar</li>
<li>Zustandswechsel werden geloggt</li>
</ul>
</div>
</div>

</div>

</div>

