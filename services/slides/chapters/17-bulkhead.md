<div class="page">

<p class="kicker">Bulkhead</p>

<div class="page-head">

## Schotten im Schiff

<img class="head-figure" src="./assets/bulkhead.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Pool pro Downstream</h3>
<p>Eigener Semaphore je Backend, kein gemeinsamer Topf.</p>
<code>Flight | Hotel | Car</code>
</div>
<div class="card">
<h3>Fail-Fast</h3>
<p>Pool voll, also sofort ablehnen. Keine Queue, kein Warten.</p>
<code>inProgress &ge; max &rarr; fallback()</code>
</div>
<div class="card">
<h3>Nicht Circuit Breaker</h3>
<p>Der CB reagiert auf Fehler, der Bulkhead auf Ressourcen-Druck.</p>
<code>Fehler &ne; Last</code>
</div>
<div class="card">
<h3>Sch&uuml;tzt den Aufrufer</h3>
<p>Defensives Pattern im Aufrufer, nicht im Backend.</p>
<code>Outbound, nicht Inbound</code>
</div>
<div class="card">
<h3>Die Gr&ouml;&szlig;e ist die Arbeit</h3>
<p>Der Code ist trivial. Die Zahl darin nicht.</p>
</div>
<div class="card market-card">
<h3>Am Markt</h3>
<div class="pills">
<span class="pill brand">Resilience4j Bulkhead</span>
<span class="pill">Polly</span>
<span class="pill">MicroProfile @Bulkhead</span>
<span class="pill">go-resiliency</span>
<span class="pill">Envoy / Istio</span>
</div>
</div>
</div>

</div>

</div>

Note:
- <strong>Schiff:</strong> Ein Leck flutet eine Kammer, ein h&auml;ngender Downstream nur seinen Pool.
- <strong>Pool pro Downstream:</strong> Im gemeinsamen Topf frisst Hotel alle Slots, Flight und Car verhungern.
- <strong>Fail-Fast:</strong> Ein Reject ist ein klares Backpressure-Signal, eine Queue tarnt das Problem.
  - Im Code: <code>hotels: []</code> mit <code>X-Fallback</code> und <code>X-Bulkhead-Full</code>. Der Fallback ist derselbe wie in Story 4, wieder ein eigenes Pattern.
- <strong>Nicht Circuit Breaker:</strong> Langsam, aber fehlerfrei? Da hilft nur der Bulkhead.
- <strong>Sch&uuml;tzt den Aufrufer:</strong> Booking und seine anderen Calls, nicht Hotel.
  - Serverseitige Bulkheads (Pools pro Consumer) gibt es auch, die bauen wir nicht.
- <strong>Die Gr&ouml;&szlig;e ist die Arbeit:</strong> Rechnen wir in Story 5 nach.
- Envoy/Istio: dasselbe im Sidecar (<code>max_requests</code> bzw. <code>max_connections</code>).
- &Uuml;berleitung: &bdquo;Sieben Zeilen Code.&ldquo;
