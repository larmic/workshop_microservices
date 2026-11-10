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
<p>Der Fehlerz&auml;hler schaltet CLOSED auf OPEN. Ein Probe-Call pr&uuml;ft, ob das Backend wieder lebt.</p>
<code>HALF_OPEN &rarr; CLOSED | OPEN</code>
</div>
<div class="card">
<h3>Schnell scheitern</h3>
<p>Bei OPEN gar kein Call. Statt 3 s Timeout sofort eine Antwort, der Aufrufer bleibt reaktiv.</p>
</div>
<div class="card own-pattern">
<span class="tag">eigenes Pattern</span>
<h3>Fallback</h3>
<p>Was der Kunde bekommt: leere Liste, Cache, anderer Provider. Der Breaker liefert das nicht, ihr schreibt es selbst.</p>
<code>flights: []</code>
</div>
<div class="card">
<h3>Outbound, nicht Inbound</h3>
<p>F&uuml;r Calls, die ihr macht. Nicht f&uuml;r Calls, die ihr bekommt.</p>
</div>
<div class="card market-card">
<h3>Am Markt</h3>
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

</div>

Note:
- <strong>Outbound:</strong> Der Schalter geh&ouml;rt vor die Calls, die euer Service macht.
- <strong>Sch&uuml;tzt beide Seiten:</strong>
  - Aufrufer: keine Threads mehr in Timeouts.
  - Backend: keine Anfragen mehr, kann sich erholen.
- <strong>Warum beim Aufrufer:</strong> Ein totes Backend kann sich nicht selbst abschalten. Nur der Aufrufer sieht die Fehler.
- <strong>Fallback ist nicht der Breaker:</strong> Der Breaker entscheidet nur, ob gerufen wird. Was zurückkommt, macht der Fallback, auch schon bei CLOSED (Karten 05–09). Ohne Fallback: bei OPEN sofort Fehler, &bdquo;fail fast&ldquo;. Resilience4j: <code>CallNotPermittedException</code> plus eigene <code>fallbackMethod</code>.
- <strong>Denkfehler:</strong> &bdquo;CB vor unseren Service, damit er nicht &uuml;berlastet.&ldquo; Das ist Rate Limiting oder Load Shedding. Auch Bulkhead sitzt im Aufrufer (Story 5).
