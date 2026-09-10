<div class="page">

<p class="kicker">Bonus zu Frage 03</p>

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
<h3>Das Backend selbst</h3>
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
- Kurze Diskussion, keine Vorlesung: Erst raten lassen, was knapp wird, dann per Klick aufdecken.
- <strong>Speicher:</strong> Request- und Response-Puffer, Kontext, Cancellation. 50.000 offene Aufrufe sind realer Heap-Druck, GC-Pausen, irgendwann OOM.
- <strong>Dateideskriptoren:</strong> Jede Verbindung belegt einen, eingehend wie ausgehend, plus den Connection-Pool des HTTP-Clients (Go: <code>http.Transport.MaxConnsPerHost</code>, Reactor Netty: 500 Default). Ist das OS-Limit erreicht, nimmt der Service gar nichts mehr an, auch keine Health-Checks.
- <strong>Das Backend:</strong> Ohne Grenze im Aufrufer dr&uuml;ckt jede zus&auml;tzliche Anfrage in ein System, das schon schw&auml;chelt. Little&rsquo;s Law gilt unabh&auml;ngig vom Threading-Modell: Steigt die Latenz, steigt die Zahl offener Aufrufe linear, bei gleicher Eingangsrate.
- Das Zahlenbeispiel: 150 Anfragen pro Sekunde und Instanz mal 3 s Timeout ergibt 450 offene Aufrufe. Laden Nutzer ungeduldig neu, verdreifacht sich die Last, ohne dass mehr Nutzer da sind.
- Take-away: &bdquo;Non-blocking spart mir den Bulkhead&ldquo; ist Wunschdenken. In reaktiven Stacks (Spring WebFlux, Reactor) sind Bulkheads nicht zuf&auml;llig weiter empfohlen, Resilience4j hat daf&uuml;r den <code>SemaphoreBulkhead</code>. Die Grenze setzt sonst das Betriebssystem, und das fragt nicht, welcher Downstream schuld ist.
