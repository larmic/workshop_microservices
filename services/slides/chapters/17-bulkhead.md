## Bulkhead

<p class="subtitle">Schotten im Schiff</p>

<div class="factor-row">

<div class="factor fragment">
<h3>Schotten</h3>
<p>Wasserdichte Bereiche &mdash; l&auml;uft einer voll, bleiben die anderen <em>trocken</em>.</p>
<aside class="notes">Analogie zum Schiff: ein Leck in einem Bereich versenkt nicht das ganze Schiff. Software-&Uuml;bersetzung: ein h&auml;ngender Downstream zieht nicht alle Ressourcen des Aufrufers in seinen Sog. Klassisches Killer-Szenario, das das Pattern motiviert: Hotel antwortet langsam in 2&nbsp;s &mdash; kein Fehler, der CB schl&auml;gt nicht aus &mdash; aber alle Threads / Connections / Goroutines des Booking-Service h&auml;ngen in Hotel-Calls. Flight und Car bekommen nichts mehr durch.</aside>
</div>

<div class="factor fragment">
<h3>Pool pro Downstream</h3>
<p>Eigener Semaphore / Thread-Pool je Service &mdash; kein gemeinsamer Topf.</p>
<code>Flight | Hotel | Car</code>
<aside class="notes">Ein <em>separater</em> Pool pro Backend. Ein gemeinsamer Pool f&uuml;r alle Aufrufe w&uuml;rde das Pattern ad absurdum f&uuml;hren &mdash; ein hungriger Hotel-Call w&uuml;rde alle Slots belegen und Flight / Car aushungern. Im Workshop-Code: <code>flightBulkhead</code>, <code>hotelBulkhead</code>, <code>carBulkhead</code> &mdash; je <code>maxConcurrent=10</code>.</aside>
</div>

<div class="factor fragment">
<h3>Fail-Fast</h3>
<p>Pool voll &rarr; sofort <span class="hl">ablehnen</span>. Keine Queue, kein Warten.</p>
<code>inFlight &ge; max &rarr; 503</code>
<aside class="notes">Drei Strategien bei Limit-&Uuml;berschreitung: Fail-Fast, Bounded Queue, Wait + Timeout. Default in Microservices: <strong>Fail-Fast</strong>. Begr&uuml;ndung: Queueing tarnt das Problem (Backend ist eh am Limit, Queue verschiebt nur), und ein sofortiger Reject ist ein klares Backpressure-Signal nach oben. Im Code: atomarer Check &amp; Increment unter Mutex, sonst rutschen mehr Calls durch als erlaubt.</aside>
</div>

<div class="factor fragment">
<h3>Sch&uuml;tzt den Aufrufer</h3>
<p>Defensives Pattern <em>im Aufrufer</em>. Nicht das Backend.</p>
<code>Outbound, nicht Inbound</code>
<aside class="notes">Wichtigste p&auml;dagogische Klarstellung: Bulkhead ist <strong>client-seitig</strong>. Gesch&uuml;tzt wird (1) der Booking-Service selbst, (2) die anderen Downstream-Calls im selben Service, (3) die eigenen Aufrufer. <em>Nicht</em> gesch&uuml;tzt: Hotel, Flight, Car. Wer das Backend vor &Uuml;berlast sch&uuml;tzen will, braucht <strong>Inbound-Patterns</strong> (Rate Limiting, Backpressure, Load Shedding) &mdash; <em>im</em> Backend, nicht im Booking-Service.</aside>
</div>

<div class="factor fragment">
<h3>Nicht Circuit Breaker</h3>
<p>CB reagiert auf <strong>Fehler</strong>, Bulkhead auf <strong>Ressourcen-Druck</strong>. Komplement&auml;r.</p>
<code>Fehler &ne; Last</code>
<aside class="notes">H&auml;ufige Verwechslung. CB = &bdquo;Backend ist krank, ich versuche es eine Weile nicht.&ldquo; Reaktion auf Fehler-Rate. Bulkhead = &bdquo;Ich verbrenne maximal N gleichzeitige Slots f&uuml;r dieses Backend, egal wie viel Last reinkommt.&ldquo; Reaktion auf Ressourcen-Druck, nicht auf Fehler. Klassisches Szenario, das <em>nur</em> Bulkhead l&ouml;st: Backend antwortet langsam, aber fehlerfrei &mdash; CB sieht keinen Grund, Bulkhead kappt nach N parallelen Calls.</aside>
</div>

</div>

<div class="market-row">

### Am Markt

<div class="chip-row">
  <span class="chip brand">Resilience4j Bulkhead</span>
  <span class="chip">Polly (.NET)</span>
  <span class="chip">MicroProfile @Bulkhead</span>
  <span class="chip">Netflix Hystrix</span>
  <span class="chip">Spring Cloud Circuit Breaker</span>
  <span class="chip">go-resiliency</span>
  <span class="chip">Envoy / Istio</span>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>

Note:
- Hook: &bdquo;Story 3 hat uns gegen <em>kaputte</em> Backends geh&auml;rtet. Aber was, wenn ein Backend gar nicht kaputt ist &mdash; nur langsam? Der CB bleibt CLOSED, weil 200 zur&uuml;ckkommt &mdash; und der Booking-Service ger&auml;t trotzdem ins Stocken.&ldquo;
- Karten-Reihenfolge bewusst: erst Analogie (Schotten), dann Mechanik (Pool pro Downstream, Fail-Fast), dann zwei wichtige Abgrenzungen (Outbound nicht Inbound, Bulkhead nicht CB).
- Demo-Tipp: Im Dashboard auf Story 4 wechseln, Backend auf &bdquo;Langsam&ldquo; stellen, dann den <code>POST /admin/burst</code>-Button dr&uuml;cken &mdash; 20 parallele Requests, Rejects werden in der Bulkhead-Karte sichtbar.
- Wer was wo nutzt: Resilience4j ist Standard im framework-freien Java (Semaphore- und ThreadPool-Variante). Hystrix Pionier mit ThreadPool-Isolation, heute End-of-Life. Polly im .NET-Lager. Im Service Mesh: Envoy/Istio macht das &uuml;ber <code>circuit_breakers.max_pending_requests</code> &mdash; sprach-agnostisch, ohne Anwendungscode.
- &Uuml;berleitung: Wir schauen uns die Mechanik konkret an &mdash; in genau der Form, wie sie im Dashboard-Spickzettel und im Story-4-Code steckt.
