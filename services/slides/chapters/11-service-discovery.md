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
- Hook: &bdquo;Im Monolithen kennen sich Module &uuml;ber Funktionsaufrufe. Im verteilten System kennen sich Services &uuml;ber &hellip; was eigentlich?&ldquo; Statische URLs (Story 1) funktionieren genau so lange, bis ihr skaliert, deployt oder eine Instanz ausf&auml;llt.
- <strong>Logische Namen:</strong> Hinter dem Namen stehen 1 bis n Instanzen. Das, was DNS f&uuml;r Hosts macht, nur dynamischer und mit Health-Wissen.
- <strong>Verzeichnis:</strong> Beim Start registriert sich der Service (Name, Adresse, Port, Tags), beim sauberen Stop deregistriert er sich. Wer wissen will, wo der Service l&auml;uft, fragt die Registry, nicht die Wiki-Seite.
- <strong>Health Checks:</strong> Die Registry pollt oder bekommt Heartbeats und nimmt kranke Instanzen aus dem Pool. Der Aufrufer bekommt nur lebende Adressen. Failover ohne manuelles Eingreifen.
- <strong>Dynamische Topologie:</strong> Skalieren = neue Instanz, Crash = neue IP, Rolling Deploy = alte raus, neue rein. Abgel&ouml;st werden <code>/etc/hosts</code>-Pflege, hartkodierte <code>BOOKING_URL</code>-Variablen und Excel-Listen im Ops-Wiki.
- <strong>Client- oder Server-Side:</strong> Client-Side hat jeder Service einen Resolver und w&auml;hlt selbst (bauen wir). Server-Side fragt ein Load Balancer oder Gateway die Registry (Traefik, AWS ALB mit Cloud Map). Client-Side spart einen Hop, Server-Side ist sprachunabh&auml;ngig.
- Am Markt: Eureka stark in der Spring-Welt, etcd unter der Haube von Kubernetes, Consul polyglott. In Kubernetes reicht meist Service plus CoreDNS, eine extra Registry ist dort selten.
- &Uuml;berleitung: Wir schauen jetzt konkret auf Consul, weil wir es im Workshop benutzen.
