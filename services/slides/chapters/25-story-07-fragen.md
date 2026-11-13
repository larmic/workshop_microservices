<div class="page">

<p class="kicker">Story 7 &middot; Recap</p>

<div class="page-head">

## Vier Fragen an euch

<span class="badge">Bonus zu 02 &darr;</span>
</div>

<div class="page-body">

<div class="numlist recap compact">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Ein Chip bleibt liegen</p>
<h3>Booking schickt das Storno per HTTP-POST, Flight ist seit der Buchung down. Wo ist das Event <span class="hl">jetzt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Event ist weg</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Das Storno kommt zweimal</p>
<h3>Ein Broker h&auml;lt das Storno fest und liefert, bis Flight best&auml;tigt. Stirbt Flight dazwischen, kommt das Event <span class="hl">noch einmal</span>.</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">eventId merken, Doppeltes ignorieren</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Booking meldet zu fr&uuml;h &bdquo;storniert&ldquo;</p>
<h3>Booking soll erst &bdquo;storniert&ldquo; sagen, wenn alle Chips zur&uuml;ck sind. Was muss zur&uuml;ckkommen, und was gilt, wenn es <span class="hl">nie</span> kommt?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Antwort zur&uuml;ck, sonst STUCK</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Brauchen wir &uuml;berhaupt einen Dirigenten?</p>
<h3>Volle Choreographie hie&szlig;e: Hotel reagiert auf Flight, Car auf Hotel. Wer merkt es, wenn Hotel das <span class="hl">nicht wei&szlig;</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Prozess braucht eine Adresse</code>
</div>
</div>

</div>

</div>

Note:
- Erst diskutieren, dann per Klick die Merks&auml;tze. Die Fragen 01 bis 03 steigen auf: vom Zuruf bis zur Antwort zur&uuml;ck.
- <strong>Nie angekommen:</strong> Das haben wir gespielt, so arbeitet die Referenz. Ein Versuch, kein Retry: <code>dispatch-failed</code> im Log, trotzdem <code>COMPENSATED</code>, das Event ist weg. Booking meldet &bdquo;storniert&ldquo;, bevor jemand gew&uuml;rfelt hat. In Story 6 hat Booking den Schmerz gesp&uuml;rt, jetzt liegt er als Chip beim Service. <strong>Spicy:</strong> Eventing macht das Problem nicht kleiner. Es macht es leiser.
- <strong>Zweimal:</strong> Nicht gespielt, ein Gedankenexperiment an den liegen gebliebenen Chips. Der Broker bringt Persistenz und Redelivery (at-least-once): Kein Storno fehlt, aber manches kommt doppelt, und Flight zahlt zweimal. Standard: <code>eventId</code> merken, in derselben Transaktion wie die Fachlogik. Unsere Backends loggen nur und deduplizieren nicht. Fan-out, Replay, Dead-Letter liefert nur ein Broker, der Webhook spart nur Infrastruktur (Bonus-Folie). <strong>Spicy:</strong> Wer Events publiziert, muss mit Duplikaten rechnen.
- <strong>Antwort zur&uuml;ck:</strong> Auch ein Gedankenexperiment. Jeder Service schickt ein Reply-Event zur&uuml;ck, zum Beispiel <code>BookingCancelled</code>. Booking wartet mit Timeout, kommt nichts, setzt es den Saga-Status <code>STUCK</code> statt <code>COMPENSATED</code>. Ehrlich, aber Booking wartet wieder, nur asynchron, also ein St&uuml;ck Orchestrierung. Das f&uuml;hrt zu 04. <strong>Spicy:</strong> Wer &bdquo;storniert&ldquo; sagt, ohne es zu wissen, l&uuml;gt h&ouml;flich.
- <strong>Dirigent:</strong> Niemand, es passiert einfach nichts. Ein Prozess mit fester Reihenfolge und Kompensation geh&ouml;rt an eine Stelle. Choreographie nur f&uuml;r einfache F&auml;lle ohne Reihenfolge (Benachrichtigung, Read-Model). Unsere Referenz ist der Mittelweg: Booking dirigiert den Hinweg, storniert wird per Zuruf. <strong>Spicy:</strong> Bei Choreographie kennt jeder jeden, nur steht es nirgends. Br&uuml;cke zu Story 8: Wer merkt es, wenn ein Event verloren geht?
- Lang: <code>docs/questions/story7.md</code>.
