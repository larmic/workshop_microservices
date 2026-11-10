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
<div>    Check: { HTTP: "http://{hostname}:8080/health",</div>
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
- Reines HTTP, kein SDK n&ouml;tig. Ihr baut in Story 3 nur das Aufl&ouml;sen; Anmelden und Abmelden machen Flight, Hotel und Car schon.
- <strong>passing=true:</strong> nur Instanzen mit gr&uuml;nem Health-Check.
- <strong>Die zwei Zahlen:</strong>
  - Interval: Nach einem Absturz l&auml;uft bis zu 10 s Verkehr ins Leere.
  - DeregisterCriticalServiceAfter r&auml;umt nach 1 min nur den Katalog auf.
- <strong>docker stop:</strong> Unsere Backends melden sich selbst ab, der Eintrag ist sofort weg. Nur bei <code>kill</code> oder Absturz greift der Check.
- <strong>&Uuml;berleitung:</strong> Consul kann noch viel mehr, das lassen wir bewusst weg.
