<div class="page">

<p class="kicker">Story 3 &middot; Recap</p>

## Vier Fragen an euch

<div class="page-body">

<div class="numlist recap compact">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Neuer Single Point of Failure?</p>
<h3>Jeder Aufruf fragt zuerst Consul. Was, wenn Consul <span class="hl">kippt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Consul down &rarr; Booking down</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Consul im Kubernetes-Cluster?</p>
<h3>Kubernetes hat eigene Discovery. Braucht es Consul <span class="hl">daneben</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">CoreDNS vs. Consul</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Stop oder Kill?</p>
<h3><code>docker stop</code> oder <code>docker kill</code>: <span class="hl">Wann</span> verschwindet der Eintrag?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">stop: sofort &middot; kill: &le; 10 s</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Registriert &ne; gesund</p>
<h3>Der Check sagt 200. Ist der Service im Moment des Calls noch <span class="hl">da</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">war &ne; ist</code>
</div>
</div>

</div>

</div>

Note:
- <strong>SPOF:</strong> Consul ist ein HA-Cluster und muss erreichbar sein. Ist er weg: Fehler, nicht raten. Der lokale Consul-Agent puffert, nicht der Service. Clients halten die Aufl&ouml;sung h&ouml;chstens Sekunden, das ist fl&uuml;chtig, kein State.
- <strong>Kubernetes:</strong> Nein, API-Server plus CoreDNS sind die Registry. Consul nur bei Hybrid, Multi-Cluster oder Brownfield. Die Frage ist meist Mesh, nicht Discovery.
- <strong>Stop vs. Kill:</strong> <code>docker stop</code> schickt SIGTERM, unsere Backends melden sich selbst ab, der Eintrag ist sofort weg. Bei <code>kill</code> oder Absturz l&auml;uft Traffic bis zum n&auml;chsten fehlgeschlagenen Check ins Leere, hier bis zu 10 s. Discovery ist immer eventually consistent, deshalb Story 4 und 5.
- <strong>Registriert &ne; gesund:</strong> Der Check kann l&uuml;gen oder veraltet sein. Discovery sagt, wo der Service <em>war</em>, nicht ob er noch da ist. Der Aufrufer muss sich selbst sch&uuml;tzen, Story 4.
- Lang: <code>docs/questions/story3.md</code>.
