<div class="page">

<p class="kicker">Story 6 &middot; Recap</p>

## F&uuml;nf Fragen an euch

<div class="page-body">

<div class="numlist recap compact dense">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Die Kompensation scheitert</p>
<h3>Flug gebucht, Hotel sagt nein, und der Storno l&auml;uft selbst in <span class="hl">5xx</span>. Was nun?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">must eventually succeed</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Wir kompensieren genau einmal</p>
<h3>Beim Fehlschlag einfach <code>FAILED</code>. Verletzen wir damit das <span class="hl">Pattern</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Monitoring = Retry</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">204 statt 404</p>
<h3><code>DELETE</code> auf eine unbekannte ID liefert <code>204</code>. Verschleiert das <span class="hl">echte Fehler</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">204 ist Pragmatik</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Die Saga h&auml;ngt</p>
<h3>Zwischen Forward und Kompensation bleibt sie stehen. Wer <span class="hl">schl&auml;gt Alarm</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Status-Counter</code>
</div>
<div class="numlist-row">
<div class="numlist-num">05</div>
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
- <strong>Kompensation scheitert:</strong> Forward darf scheitern, Kompensation muss gelingen: Idempotenz, Retry mit Backoff, persistenter Log, Operator-Inbox. <strong>Spicy:</strong> Ohne Plan daf&uuml;r ist es keine Saga, sondern Hoffnung.
- <strong>Genau einmal:</strong> Formal ja, vertretbar nur als bewusste Entscheidung mit persistiertem Status und Alert auf <code>FAILED</code>. <strong>Spicy:</strong> Ohne Retry im Code muss das Monitoring der Retry sein.
- <strong>204 statt 404:</strong> Nein. Idempotent sind beide Varianten, 204 erspart dem Orchestrator nur die Sonderlogik &bdquo;404 = ok&ldquo;, sobald er wiederholt.
- <strong>Saga h&auml;ngt:</strong> Logzeilen reichen nicht, n&ouml;tig ist ein Status-Counter mit Alarm bei <code>COMPENSATING &gt; 0</code> &uuml;ber N Minuten. <strong>Spicy:</strong> H&auml;ngende Sagas sind still.
- <strong>Status im RAM:</strong> Der Flug bleibt gebucht, und niemand wei&szlig; es mehr. Persistenz braucht jeder mehrstufige Prozess, Saga macht es nur sichtbar. <strong>Spicy:</strong> Die echte Frage ist selbst bauen oder Temporal.
- Reserve: Eventing statt sync, Owner.
- Lang: <code>docs/questions/story6.md</code>.
