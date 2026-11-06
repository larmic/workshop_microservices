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

<p class="codenote">Das <code>Interval</code> entscheidet, wie lange nach einem Absturz Verkehr ins Leere l&auml;uft. Der zweite Wert r&auml;umt nur auf.</p>

</div>

</div>

Note:
- <strong>Drei Calls:</strong> Anmelden beim Start, Abmelden beim Shutdown, Aufl&ouml;sen bei jedem Aufruf. Reines HTTP, kein SDK n&ouml;tig.
- <strong>passing=true:</strong> Consul liefert nur Instanzen mit gr&uuml;nem Health-Check.
- <strong>Die zwei Zahlen:</strong> Das Interval entscheidet, wie lange nach einem Absturz Verkehr ins Leere l&auml;uft, hier bis zu 10 s. DeregisterCriticalServiceAfter r&auml;umt nach 1 min nur den Katalog auf.
- <strong>Graceful Shutdown:</strong> Bei SIGTERM (<code>docker stop</code>) melden sich unsere Backends selbst ab, der Eintrag ist sofort weg. Nur bei <code>kill</code> oder Absturz wartet ihr auf den Check.
- <strong>Weiter:</strong> Consul kann noch viel mehr. Das lassen wir bewusst weg.
