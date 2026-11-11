<div class="page">

<p class="kicker">Story 5 &middot; Auswertung</p>

## Gleiche St&uuml;hle, eine Wand

<div class="page-body littles">

<div class="ltable result">
<div class="ltable-head">Gez&auml;hlt bei STOPP 3</div>
<div class="ltable-head">Runde 1</div>
<div class="ltable-head">Runde 2</div>
<div class="ltable-head">Was es hei&szlig;t</div>
<div>beantwortet (Car)</div><div class="fragment" data-fragment-index="1"><code>4</code></div><div class="fragment" data-fragment-index="2"><code>10</code></div><div class="fragment" data-fragment-index="2">Verf&uuml;gbarkeit vor Vollst&auml;ndigkeit</div>
<div>davon mit Hotels</div><div class="fragment" data-fragment-index="1"><code>4</code></div><div class="fragment" data-fragment-index="2"><code>2</code></div><div class="fragment" data-fragment-index="2">ehrlicher Tausch</div>
<div>warten hinter den St&uuml;hlen</div><div class="fragment" data-fragment-index="1"><code>4</code></div><div class="fragment" data-fragment-index="2"><code>0</code></div><div class="fragment" data-fragment-index="2">die Queue ist weg</div>
<div>Flight abgeklatscht</div><div class="fragment" data-fragment-index="1"><code>6</code></div><div class="fragment" data-fragment-index="2"><code>11</code></div><div class="fragment" data-fragment-index="2">gesund, trotzdem ausgehungert</div>
<div>Circuit Breaker</div><div class="fragment" data-fragment-index="1"><code>CLOSED</code></div><div class="fragment" data-fragment-index="2"><code>CLOSED</code></div><div class="fragment" data-fragment-index="2">langsam ist kein Fehler</div>
</div>

<div class="swap swap-foot">
<div class="callout fragment" data-fragment-index="3">Weniger Hotels, aber niemand wartet: Die Wand macht Hotel nicht schneller, sie h&auml;lt Flight und Car frei.</div>
</div>

</div>

</div>

Note:
- <strong>Z&auml;hler abfragen:</strong> erst die gemerkten Werte aus Runde 1 (Klick), dann Runde 2 (Klick). Hotel z&auml;hlt mit, die Zahl steckt in &bdquo;davon mit Hotels&ldquo;.
- <strong>&bdquo;Welche Runde war besser?&ldquo;</strong>
  - Runde 1 hatte mehr Hotels.
  - Runde 2 hat mehr beantwortet, und niemand stand an.
- <strong>&bdquo;Wie viele St&uuml;hle h&auml;tte Hotel gebraucht?&ldquo;</strong>
  - Vier, Hotel hat einen.
  - Deshalb bekommt jeder Vierte Hotels (16, 20, 24).
  - Br&uuml;cke zu 18a.
- <strong>&bdquo;Breaker?&ldquo;</strong> CLOSED, ein besetzter Stuhl ist kein Fehler.
- <strong>&bdquo;Wem geh&ouml;ren die St&uuml;hle?&ldquo;</strong>
  - Booking, nicht Hotel.
  - F&uuml;nf Booking-Instanzen h&auml;tten f&uuml;nf Hotel-St&uuml;hle (Recap 05).
- <strong>&Uuml;bersetzung in den Code:</strong>
  - Stuhl = Slot
  - Wand = Bulkhead, <code>max[service]</code>
  - hinter den St&uuml;hlen = Queue (gibt es im Code nicht)
  - sofort ohne Hotels = <code>fallback()</code>
