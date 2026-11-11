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
<p>H&auml;ngt Hotel, bindet es h&ouml;chstens <code class="nowrap">max</code> Calls von Booking. Der Rest bleibt frei f&uuml;r Flight und Car.</p>
<code>Outbound, nicht Inbound</code>
</div>
<div class="card">
<h3>Die Gr&ouml;&szlig;e ist die Arbeit</h3>
<p>Zu klein: unn&ouml;tige Rejects. Zu gro&szlig;: sch&uuml;tzt nicht. Die Zahl h&auml;ngt an der Latenz.</p>
<code>Slots = Aufrufe/s &times; Dauer</code>
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
- <strong>Sch&uuml;tzt den Aufrufer:</strong> Booking hat begrenzte Threads und Connections. H&auml;ngt Hotel, wartet bald alles in Booking auf Hotel, und Booking kann auch Flight und Car nicht mehr rufen. Booking f&auml;llt um, obwohl nur Hotel krank ist.
  - Der Bulkhead begrenzt, wie viel von Booking an Hotel h&auml;ngen darf (Referenz: 10). Gesch&uuml;tzt ist Booking, nicht Hotel.
  - Hotel selbst sch&uuml;tzt sich mit Rate Limiting oder Load Shedding. Serverseitige Bulkheads (Pools pro Consumer) gibt es auch, die bauen wir nicht.
- <strong>Die Gr&ouml;&szlig;e ist die Arbeit:</strong> Zu klein lehnt ein gesundes Hotel ab, zu gro&szlig; sch&uuml;tzt nicht.
  - Little&rsquo;s Law: Slots = Aufrufe pro Sekunde &times; Dauer. 20/s &times; 0,5 s = 10. Bei 3 s Timeout w&auml;ren es 60.
  - Dimensioniert wird auf die normale Latenz. Story 5 spielt das mit zwei St&uuml;hlen und einer Wand: Hotel br&auml;uchte vier, hat einen.
- Envoy/Istio: dasselbe im Sidecar (<code>max_requests</code> bzw. <code>max_connections</code>).
- &Uuml;berleitung: &bdquo;Sieben Zeilen Code.&ldquo;
