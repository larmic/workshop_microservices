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
<code>inProgress &ge; max &rarr; 503</code>
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
</div>

<div class="market">
<h4>Am Markt</h4>
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

Note:
- Analogie zum Schiff: Ein Leck in einem Bereich versenkt nicht das ganze Schiff. In Software: Ein h&auml;ngender Downstream zieht nicht alle Ressourcen des Aufrufers in seinen Sog.
- <strong>Pool pro Downstream:</strong> Ein gemeinsamer Pool f&uuml;r alle Aufrufe w&uuml;rde das Pattern ad absurdum f&uuml;hren, ein hungriger Hotel-Call w&uuml;rde alle Slots belegen und Flight und Car aushungern. In der Referenz: <code>flightBulkhead</code>, <code>hotelBulkhead</code>, <code>carBulkhead</code>.
- <strong>Fail-Fast:</strong> Drei Strategien bei &Uuml;berschreitung: Fail-Fast, Bounded Queue, Wait plus Timeout. Default in Microservices ist Fail-Fast: Queueing tarnt das Problem, der sofortige Reject ist ein klares Backpressure-Signal nach oben. Im Code m&uuml;ssen Check und Increment atomar sein.
- <strong>Nicht Circuit Breaker:</strong> CB hei&szlig;t &bdquo;Backend ist krank, ich versuche es eine Weile nicht&ldquo;, Reaktion auf die Fehlerrate. Bulkhead hei&szlig;t &bdquo;ich verbrenne maximal N Slots f&uuml;r dieses Backend&ldquo;, Reaktion auf Ressourcen-Druck. Das Szenario, das nur Bulkhead l&ouml;st: Backend antwortet langsam, aber fehlerfrei.
- <strong>Sch&uuml;tzt den Aufrufer:</strong> Bulkhead ist client-seitig. Gesch&uuml;tzt werden der Booking-Service, die anderen Downstream-Calls im selben Service und die eigenen Aufrufer. Nicht gesch&uuml;tzt: Hotel, Flight, Car. Daf&uuml;r braucht das Backend Inbound-Patterns wie Rate Limiting, Backpressure, Load Shedding.
- <strong>Die Gr&ouml;&szlig;e ist die Arbeit:</strong> Deshalb programmieren wir heute nicht, sondern spielen und rechnen. N&auml;chste Folie zeigt, wie wenig Code das ist.
- Am Markt: Resilience4j ist Standard im Java-Umfeld (Semaphore- und ThreadPool-Variante), Polly im .NET-Lager, MicroProfile per Annotation, go-resiliency f&uuml;r Go. Envoy und Istio regeln dasselbe im Sidecar &uuml;ber <code>max_pending_requests</code>, sprachunabh&auml;ngig.
