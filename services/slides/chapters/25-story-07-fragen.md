<div class="page">

<p class="kicker">Story 7 &middot; Recap</p>

<div class="page-head">

## Vier Fragen an euch

<span class="badge">Bonus zu 03 &darr;</span>
</div>

<div class="page-body">

<div class="numlist recap compact">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Runde 2: Ein Gummib&auml;rchen fehlt</p>
<h3>Booking schickt das Storno per HTTP-POST, Flight ist seit der Buchung down. Wo ist das Event <span class="hl">jetzt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Event ist weg</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Runde 3: Ein Becher hat drei</p>
<h3>Der Broker liefert, bis das Backend best&auml;tigt. Stirbt es dazwischen, kommt das Event <span class="hl">noch einmal</span>.</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Nummer merken, Doppeltes ignorieren</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Runde 2 gegen Runde 3: Webhooks reichen doch</p>
<h3>Kein Broker, keine Infrastruktur. Wozu Kafka, RabbitMQ oder NATS <span class="hl">&uuml;berhaupt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">durable messaging</code>
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
- <strong>Nie angekommen:</strong> Die Referenz loggt <code>dispatch-failed</code>, markiert trotzdem <code>COMPENSATED</code>, das Event ist weg. Abhilfe: Retry, Outbox, Dead-Letter, Broker mit at-least-once. <strong>Spicy:</strong> Eventing verschiebt das Problem, Booking sieht nichts mehr.
- <strong>Zweimal:</strong> Verursacht der Broker: At-least-once liefert, bis der Empf&auml;nger best&auml;tigt. Standard: <code>eventId</code> merken, in derselben Transaktion wie die Fachlogik. <strong>Spicy:</strong> Wer Events publiziert, muss mit Duplikaten rechnen.
- <strong>Webhooks reichen:</strong> Persistenz, Redelivery, Fan-out, Backpressure, Replay und Entkopplung in der Zeit liefert nur ein Broker. Webhook spart nur Infrastruktur (Bonus-Folie).
- <strong>Dirigent:</strong> Niemand, es passiert einfach nichts. Ein Prozess mit fester Reihenfolge und Kompensation geh&ouml;rt an eine Stelle. Choreographie nur f&uuml;r einfache F&auml;lle ohne Reihenfolge (Benachrichtigung, Read-Model). Unsere Referenz ist der Mittelweg: Booking dirigiert den Hinweg, storniert wird per Zuruf. <strong>Spicy:</strong> Bei Choreographie kennt jeder jeden, nur steht es nirgends. Br&uuml;cke zu Story 8: Wer merkt es, wenn ein Event verloren geht?
- Lang: <code>docs/questions/story7.md</code>.
