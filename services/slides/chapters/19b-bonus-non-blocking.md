<div class="page">

<p class="kicker">Bonus zu Frage 02</p>

## Wo ist die Grenze bei non-blocking?

<p class="subtitle">Threads sind billig. Was wird stattdessen knapp?</p>

<div class="page-body bonus">

<div class="swap">
<div class="cards cards-3 ghost fragment fade-out" data-fragment-index="1">
<div class="card">?</div>
<div class="card">?</div>
<div class="card">?</div>
</div>
<div class="cards cards-3 violet compact fragment" data-fragment-index="1">
<div class="card">
<h3>Speicher</h3>
<p>Jede offene Anfrage h&auml;lt ihre Daten, Puffer und Verbindung, linear mit der Gleichzeitigkeit.</p>
</div>
<div class="card">
<h3>Dateideskriptoren</h3>
<p>Jede Verbindung ist einer, eingehend wie ausgehend. Ist die Tabelle voll, nimmt der Service <strong>gar nichts</strong> mehr an.</p>
</div>
<div class="card">
<h3>Der Downstream selbst</h3>
<p>Wer unbegrenzt nachschiebt, dr&uuml;ckt immer mehr gleichzeitige Anfragen in ein System, das schon schw&auml;chelt.</p>
</div>
</div>
</div>

<p class="bonus-text fragment" data-fragment-index="1">150 Anfragen pro Sekunde und Instanz, Hotel bei 3 Sekunden: 450 gleichzeitig offene Aufrufe. Dann laden die Nutzer neu, und aus 450 werden 1350.</p>

<div class="foot fragment" data-fragment-index="1">
<div class="callout">Die Grenze verschwindet nicht. Sie wird nur vom Betriebssystem gesetzt statt von euch.</div>
</div>

</div>

</div>

Note:
- <strong>Erst raten lassen</strong>, was knapp wird, dann Klick.
- <strong>Speicher:</strong> Puffer, Kontext, Cancellation. 50.000 offene Aufrufe sind Heap-Druck und GC-Pausen, irgendwann OOM.
- <strong>Dateideskriptoren:</strong>
  - Jede Verbindung belegt einen, dazu der Connection-Pool des HTTP-Clients. In Go ist <code>MaxConnsPerHost</code> per Default unbegrenzt.
  - Ist das Limit erreicht, kommen auch keine Health-Checks mehr durch.
- <strong>Downstream:</strong> Little&rsquo;s Law gilt unabh&auml;ngig vom Threading-Modell. Steigt die Latenz, steigen die offenen Aufrufe.
- <strong>Zahlenbeispiel:</strong> 150/s &times; 3 s = 450 offene Aufrufe. Ungeduldiges Neuladen verdreifacht die Last ohne einen neuen Nutzer.
- <strong>Take-away:</strong>
  - &bdquo;Non-blocking spart mir den Bulkhead&ldquo; ist Wunschdenken.
  - Reaktive Stacks empfehlen ihn weiter, Resilience4j hat daf&uuml;r den <code>SemaphoreBulkhead</code>.
  - Sonst setzt das Betriebssystem die Grenze, und das fragt nicht, welcher Downstream schuld ist.
