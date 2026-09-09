<div class="page">

<p class="kicker">Service Discovery &middot; Consul</p>

## Drei HTTP-Aufrufe

<p class="subtitle">Mehr braucht eine Registry nicht.</p>

<div class="page-body codebody">

<div class="codeblock">
<div class="k">Anmelden, beim Start</div>
<div>  PUT consul:8500/v1/agent/service/register</div>
<div>  { Name: "flight-service", ID: "flight-service-{hostname}",</div>
<div>    Address: "{hostname}", Port: 8080,</div>
<div>    Check: { HTTP: "http://{host}:8080/health",</div>
<div>             <span class="hi">Interval: "10s"</span>, <span class="hi">DeregisterCriticalServiceAfter: "1m"</span> } }</div>
<div class="gap k">Abmelden, beim Shutdown</div>
<div>  PUT consul:8500/v1/agent/service/deregister/{serviceID}</div>
<div class="gap k">Aufl&ouml;sen, bei jedem Aufruf</div>
<div>  GET consul:8500/v1/health/service/{name}?passing=true</div>
<div class="dim">  &rarr; [ { Service: { Address, Port } }, &hellip; ]</div>
</div>

<p class="codenote">Diese beiden Zahlen entscheiden, wie lange nach einem Ausfall noch Verkehr ins Leere l&auml;uft.</p>

</div>

</div>

Note:
- Drei Aufrufe: Lifecycle der Instanz (anmelden, abmelden) und das, was der Aufrufer bei jedem Request tut (aufl&ouml;sen). Danach normaler HTTP-Call an <code>http://{Address}:{Port}</code>.
- Die zwei markierten Zahlen: <code>Interval: 10s</code> ist der Takt des Health-Checks, <code>DeregisterCriticalServiceAfter: 1m</code> die Frist, nach der Consul eine tote Instanz selbst entfernt, falls der Shutdown-Hook nicht durchkommt (kill -9, OOM). Zwischen Ausfall und Entfernen l&auml;uft Verkehr ins Leere. Kleinere Werte = schnellere Heilung, mehr Last auf Consul.
- <code>?passing=true</code> ist der Knackpunkt: ohne den Filter kommen auch ungesunde Instanzen zur&uuml;ck. Der Health-Check ist nur wertvoll, wenn der Client ihn respektiert.
- Identischer Pseudo-Code inklusive <code>resolve(name)</code> und der Schleife &uuml;ber Flight, Hotel, Car steht im Dashboard unter Story 3, &bdquo;Spickzettel&ldquo;. Wiedererkennung gewollt.
- Im Referenz-Image (<code>booking/story3/</code>) ist genau das implementiert: <code>services/shared/consul/register.go</code> und <code>resolver.go</code>.
- Trade-off: vor jedem Call aufl&ouml;sen ist stets aktuell, kostet aber Last auf Consul. Alternative: Cache mit Watch oder TTL.
