<div class="page">

<p class="kicker">Service Discovery</p>

<div class="page-head">

## Was die Registry leistet

<img class="head-figure" src="./assets/service-discovery.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Logische Namen</h3>
<p>Der Code ruft <code>flight-service</code>, nicht <code>10.0.0.5:8080</code>.</p>
</div>
<div class="card">
<h3>Ein Verzeichnis</h3>
<p>Services melden sich beim Start an und beim Herunterfahren wieder ab.</p>
</div>
<div class="card">
<h3>Health Checks</h3>
<p>Wer den Check nicht besteht, wird nicht mehr ausgeliefert.</p>
</div>
<div class="card">
<h3>Dynamische Topologie</h3>
<p>Adressen &auml;ndern sich, der Code nicht. Kein <code>/etc/hosts</code>, keine URL-Liste.</p>
</div>
<div class="card">
<h3>Client- oder Server-Side</h3>
<p>Der Aufrufer l&ouml;st selbst auf, oder ein Load Balancer tut es f&uuml;r ihn.</p>
<p class="muted">Wir bauen die erste Variante.</p>
</div>
</div>

<div class="market">
<h4>Am Markt</h4>
<div class="pills">
<span class="pill brand">HashiCorp Consul</span>
<span class="pill">Kubernetes / CoreDNS</span>
<span class="pill">Netflix Eureka</span>
<span class="pill">AWS Cloud Map</span>
<span class="pill">Apache Zookeeper</span>
</div>
</div>

</div>

</div>

Note:
- Wir setzen hier Consul ein.
