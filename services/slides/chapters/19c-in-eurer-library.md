<div class="page">

<p class="kicker">Kapitel 2 &middot; R&uuml;ckblick</p>

## Breaker und Bulkhead in eurer <span class="hl">Library</span>

<div class="page-body">

<p class="lead">In Produktion schreibt ihr das nicht selbst. Ihr stellt dieselben Schrauben ein.</p>

<table class="knobs">
<thead>
<tr><th>Im Spiel</th><th>Stellschraube</th><th>Beispiel Resilience4j</th></tr>
</thead>
<tbody>
<tr class="knob-group"><td colspan="3">Circuit Breaker &middot; Story 4</td></tr>
<tr><td>Laut z&auml;hlen, bis der Timeout kommt</td><td>Timeout pro Aufruf</td><td><code>timeoutDuration: 3s</code></td></tr>
<tr><td>F&uuml;nf rote Karten auf dem Board</td><td>Ab wann &ouml;ffnet der Breaker?</td><td><code>failureRateThreshold: 50</code></td></tr>
<tr><td>OPEN, 30 s Pause</td><td>Wie lange bleibt er offen?</td><td><code>waitDurationInOpenState: 30s</code></td></tr>
<tr><td>Eine Karte auf dem Probe-Feld</td><td>Wie viele Proben?</td><td><code>permittedNumberOfCallsInHalfOpenState: 1</code></td></tr>
<tr class="knob-group"><td colspan="3">Bulkhead &middot; Story 5</td></tr>
<tr><td>Ein Stuhl hinter der Wand</td><td>Wie viele gleichzeitig?</td><td><code>maxConcurrentCalls: 10</code></td></tr>
<tr><td>Hinter den St&uuml;hlen stehen</td><td>Wie lange auf einen Platz warten?</td><td><code>maxWaitDuration: 0</code></td></tr>
</tbody>
</table>

<p class="step-note">Die Namen wechseln, die Schrauben bleiben: Polly (.NET), opossum (Node), gobreaker (Go), Envoy / Istio (Mesh).</p>

</div>

</div>

Note:
- Zweck: Abschluss von Kapitel 2, beide Spiele (Breaker aus Story 4, Bulkhead aus Story 5) als Stellschrauben einer Library. Der Handgriff f&uuml;r Montag. Niemand muss Java k&ouml;nnen, die Spalte rechts ist nur ein Beispiel. Die mittlere Spalte ist die, die jede Library hat.
- Frage in den Raum: &bdquo;Welche Library nutzt ihr? Wie hei&szlig;en die sechs Schrauben dort?&ldquo;
- Stolperstein: Resilience4j z&auml;hlt eine Fehlerquote im Fenster (<code>slidingWindowSize</code>), unsere Referenz und das Spiel z&auml;hlen Fehler in Folge. Ein wackliges Flight, das jedes zweite Mal antwortet, w&uuml;rde Resilience4j &ouml;ffnen, die Referenz nicht.
- <code>maxWaitDuration: 0</code> hei&szlig;t: nicht anstellen, sofort Fallback. Das war Runde 2.
- Die Zahlen kommen aus der Messung, nicht aus dem Bauch: Timeout aus der gemessenen Latenz, Slots aus Aufrufen pro Sekunde mal Dauer.
