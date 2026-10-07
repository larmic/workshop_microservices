<div class="page">

<p class="kicker">Circuit Breaker</p>

<div class="page-head">

## Der Schutzschalter

<img class="head-figure" src="./assets/circuitbreaker.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Schutzschalter</h3>
<p>Wie die Sicherung im Stromkasten: bei zu vielen Fehlern raus.</p>
</div>
<div class="card">
<h3>Drei Zust&auml;nde</h3>
<p>Wechsel anhand Fehlerz&auml;hler. Ein Probe-Call pr&uuml;ft, ob das Backend wieder lebt.</p>
<code>CLOSED &rarr; OPEN &rarr; HALF_OPEN</code>
</div>
<div class="card">
<h3>Schnell scheitern</h3>
<p>Statt 3 s Timeout sofort Fallback. Der Aufrufer bleibt reaktiv.</p>
</div>
<div class="card">
<h3>Fallback-Strategie</h3>
<p>Leere Liste, Cache, anderer Provider. Hauptsache etwas.</p>
<code>flights: []</code>
</div>
<div class="card">
<h3>Outbound, nicht Inbound</h3>
<p>F&uuml;r Calls, die ihr macht. Nicht f&uuml;r Calls, die ihr bekommt.</p>
</div>
</div>

<div class="market">
<h4>Am Markt</h4>
<div class="pills">
<span class="pill">Resilience4j</span>
<span class="pill">Polly</span>
<span class="pill">Spring Cloud Circuit Breaker</span>
<span class="pill">gobreaker</span>
<span class="pill">Envoy / Istio</span>
</div>
</div>

</div>

</div>

Note:
- <strong>Outbound hei&szlig;t:</strong> Der Schalter geh&ouml;rt vor die Calls, die euer Service macht. Nicht vor die, die er bekommt.
- <strong>Er sch&uuml;tzt trotzdem beide Seiten:</strong> Der Aufrufer bindet keine Threads mehr in Timeouts. Das Backend bekommt keine Anfragen mehr und kann sich erholen.
- <strong>Warum beim Aufrufer:</strong> Ein totes oder &uuml;berlastetes Backend kann sich nicht selbst abschalten. Nur der Aufrufer sieht Timeouts und Fehler, nur er kann aufh&ouml;ren.
- <strong>Der Denkfehler:</strong> &bdquo;Wir setzen einen CB vor unseren Service, damit er nicht &uuml;berlastet.&ldquo; Das ist Rate Limiting oder Bulkhead (Story 5), kein Circuit Breaker.
