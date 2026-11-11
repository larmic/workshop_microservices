<div class="page">

<p class="kicker">Story 5 &middot; Auswertung</p>

## Gleiche St&uuml;hle, eine Wand

<div class="page-body littles">

<div class="ltable result">
<div class="ltable-head">Gez&auml;hlt bei STOPP 3</div>
<div class="ltable-head">Runde 1</div>
<div class="ltable-head">Runde 2</div>
<div class="ltable-head">Was es hei&szlig;t</div>
<div>Hotel: bei mir gesessen</div><div><code>6</code></div><div><code>3</code></div><div>Runde 1 lieferte mehr Hotels</div>
<div>Flight: abgeklatscht</div><div><code>6</code></div><div><code>11</code></div><div>in Runde 1 gesund, trotzdem ausgehungert</div>
<div>Car: beantwortet</div><div><code>4</code></div><div><code>10</code></div><div>Verf&uuml;gbarkeit vor Vollst&auml;ndigkeit</div>
<div>davon komplett mit Hotels</div><div><code>4</code></div><div><code>2</code></div><div>ehrlicher Tausch</div>
<div>davon sofort ohne Hotels</div><div><code>0</code></div><div><code>8</code></div><div>Fallback <code>hotels: []</code></div>
<div>warten hinter den St&uuml;hlen</div><div><code>4</code></div><div><code>0</code></div><div>die Queue ist weg</div>
<div>Circuit Breaker</div><div><code>CLOSED</code></div><div><code>CLOSED</code></div><div>langsam ist kein Fehler</div>
<div>St&uuml;hle bei Booking</div><div><code>2</code></div><div><code>2</code></div><div>gleich viele, nur getrennt</div>
</div>

<div class="swap swap-foot">
<div class="callout">Weniger Hotels, aber niemand wartet: Die Wand macht Hotel nicht schneller, sie h&auml;lt Flight und Car frei.</div>
</div>

</div>

</div>

Note:
- <strong>Z&auml;hler abfragen:</strong> Erst Hotel, Flight, Car f&uuml;r Runde 2, dann die gemerkten Werte aus Runde 1 daneben. Jeder darf nachrechnen.
- <strong>&bdquo;Welche Runde war besser?&ldquo;</strong> Ehrlich bleiben: Runde 1 hatte mehr Hotels, Runde 2 hat mehr beantwortet und niemand stand an.
- <strong>&bdquo;Wie viele St&uuml;hle h&auml;tte Hotel gebraucht?&ldquo;</strong> Vier: Jede Nummer kommt einer dazu, jeder bleibt vier. Hotel hat einen, deshalb bekommt jeder Vierte Hotels (16, 20, 24). Br&uuml;cke zu 18a.
- <strong>&bdquo;Steht der Breaker noch auf CLOSED?&ldquo;</strong> Ja. Ein besetzter Stuhl ist kein Fehler, der Reject erreicht den Breaker nie.
- <strong>&bdquo;Wem geh&ouml;ren die St&uuml;hle?&ldquo;</strong> Booking. Hotel wei&szlig; nichts davon. Die Wand sch&uuml;tzt Booking, nicht Hotel. F&uuml;nf Booking-Instanzen h&auml;tten f&uuml;nf Hotel-St&uuml;hle (Recap 05).
- <strong>&Uuml;bersetzung in den Code:</strong>
  - Stuhl = Slot, ein Stuhl = 10 Slots. Zwei St&uuml;hle ohne Wand = ein gemeinsamer Pool.
  - Wand = Bulkhead, Abteil = <code>max[service]</code>.
  - Hinter den St&uuml;hlen stehen = Queue. Die gibt es mit Bulkhead im Code nicht.
  - Sofort fertig ohne Hotels = <code>fallback()</code>: <code>hotels: []</code> mit <code>X-Fallback</code> und <code>X-Bulkhead-Full</code>.
- Runde 1 ist der Thread-Pool klassischer Server, die Go-Referenz hat keinen. Runde 2 ist die Referenz.
