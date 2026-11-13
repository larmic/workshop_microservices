<div class="page">

<p class="kicker">Story 6 &middot; Recap</p>

## Vier Fragen an euch

<div class="page-body">

<div class="numlist recap compact">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Die Kompensation scheitert</p>
<h3>Das Storno l&auml;uft selbst in <span class="hl">5xx</span>, und wir versuchen es nur einmal. Was nun?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Storno muss irgendwann klappen</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">204 statt 404</p>
<h3><code>DELETE</code> auf eine unbekannte ID liefert <code>204</code>. Verschleiert das <span class="hl">echte Fehler</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">204 ist Pragmatik</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Die Saga h&auml;ngt</p>
<h3>Zwischen Buchen und Kompensation bleibt sie stehen. Wer <span class="hl">schl&auml;gt Alarm</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Status-Counter</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Der Status lebt im RAM</p>
<h3>Booking st&uuml;rzt mitten in <code>COMPENSATING</code> ab. Was ist jetzt mit dem <span class="hl">Flug</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Crash-Recovery</code>
</div>
</div>

</div>

</div>

Note:
- Erst diskutieren, dann per Klick die Merks&auml;tze.
- <strong>Kompensation scheitert:</strong> F&uuml;r einen gescheiterten Buchungsschritt gibt es die Kompensation. F&uuml;r eine gescheiterte Kompensation gibt es nichts mehr, also muss sie gelingen, nicht beim ersten Versuch, aber irgendwann.
  - Daf&uuml;r: idempotentes <code>DELETE</code>, Retry mit Backoff, persistenter Status, am Ende ein Mensch.
  - Alert auf <code>COMPENSATION_FAILED</code> (Dashboard: &bdquo;KOMP. FEHLER&ldquo;), nicht auf <code>FAILED</code>: Das ist auch jede sauber kompensierte Saga.
  - <strong>Spicy:</strong> Ohne Retry im Code muss das Monitoring der Retry sein.
- <strong>204 statt 404:</strong> Nein. Idempotent sind beide Varianten, 204 erspart dem Orchestrator nur die Sonderlogik &bdquo;404 = ok&ldquo;, sobald er wiederholt.
- <strong>Saga h&auml;ngt:</strong> N&ouml;tig ist ein Status-Counter mit Alarm bei <code>COMPENSATING &gt; 0</code> &uuml;ber N Minuten. <strong>Spicy:</strong> H&auml;ngende Sagas sind still.
- <strong>Status im RAM:</strong> Der Flug bleibt gebucht, und niemand wei&szlig; es mehr. Saga macht fehlende Persistenz nur sichtbar. <strong>Spicy:</strong> Die echte Frage ist selbst bauen oder Temporal.
- Reserve: Eventing statt sync, Owner.
- <strong>Preis:</strong> Zwischenzust&auml;nde sind sichtbar, und die Kompensation muss irgendwann klappen, sonst bleibt eine halbe Reise.
- Lang: <code>docs/questions/story6.md</code>.
