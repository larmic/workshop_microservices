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

