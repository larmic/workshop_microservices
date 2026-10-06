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
<p class="numlist-label">Ohne Broker: Der Zettel kommt nie an</p>
<h3>Booking schickt das Event per HTTP-POST, Hotel ist down. Wo ist das Event <span class="hl">jetzt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Event ist weg</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Mit Broker: Der Zettel kommt zweimal</p>
<h3>Der Broker liefert, bis das Backend best&auml;tigt. Stirbt es dazwischen, kommt das Event <span class="hl">noch einmal</span>.</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Nummer merken, Doppeltes ignorieren</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Webhooks reichen doch</p>
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
- Alle vier Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein. Frage 03 hat eine Bonus-Folie darunter.
- <strong>Zettel kommt nie an, meine Antwort:</strong> In der Webhook-Referenz passiert genau nichts. Der Fehler wird geloggt, der Step trotzdem als <code>COMPENSATED</code> markiert, die Saga geht auf <code>FAILED</code>, der Kunde bekommt seine Antwort. Booking sagt &bdquo;Event ist raus&ldquo;, in Wahrheit wurde es nie empfangen. Lehrbuch-Antworten: Retry mit Backoff, Outbox-Pattern (Event in derselben Transaktion wie der Saga-State, ein Worker publiziert), Dead-Letter-Queue, ein Broker mit at-least-once. <strong>Spicy:</strong> Eventing eliminiert das &bdquo;Backend kurz weg&ldquo;-Problem nicht, es verschiebt es. In Story 6 hat Booking den Schmerz gesp&uuml;rt. In Story 7 sieht Booking gar nichts.
- <strong>Zettel kommt zweimal, meine Antwort:</strong> Das l&ouml;st der Broker nicht, er verursacht es. At-least-once hei&szlig;t: liefern, bis der Empf&auml;nger best&auml;tigt. Geht die Best&auml;tigung verloren, kommt das Event noch einmal. Unsere Referenz ist idempotent durch Zustandslosigkeit, das Backend loggt nur. In Produktion ist das die Ausnahme. Doppelte Zustellung passiert, wenn die Antwort auf den POST verloren geht und der Sender wiederholt, wenn der Broker nach einem Consumer-Crash erneut liefert, wenn der Outbox-Worker zwischen Senden und Markieren abst&uuml;rzt. Ohne Dedup: R&uuml;ckerstattung zweimal, Reservierung doppelt freigegeben. Standard: <code>eventId</code> in einer kleinen Tabelle merken, in derselben Transaktion wie die Fachlogik. <strong>Spicy:</strong> Wer Events publiziert, muss davon ausgehen, dass sie mehrfach ankommen. Wer das ignoriert, hat einen Bug, der erst unter Last zuschl&auml;gt.
- <strong>Webhooks reichen doch, meine Antwort:</strong> F&uuml;r genau die Dinge, die der Webhook strukturell nicht kann: Persistenz, Redelivery, Fan-out, Backpressure, Dead-Letter, Consumer-Crash-Recovery, Ordering, Replay, Entkopplung in der Zeit. Webhook gewinnt nur beim Infrastruktur-Aufwand, und genau darauf fallen Teams herein. Details auf der Bonus-Folie.
- <strong>Dirigent, meine Antwort:</strong> F&uuml;r unsere Buchung ist Orchestrierung der bessere Weg. Drei Schritte in fester Reihenfolge mit Kompensation sind ein Prozess, und den will man an einer Stelle lesen, testen und im Log verfolgen. Volle Choreographie verteilt die Reihenfolge auf die K&ouml;pfe der Services: Jeder muss wissen, auf wen er h&ouml;rt. Fehlt ein Abonnement, bricht die Kette, und niemand bekommt einen Fehler, es passiert einfach nichts. Das ist der verteilte Monolith: lose gekoppelt zur Laufzeit, fest gekoppelt im Kopf. Was wir in Runde 2 gespielt haben, ist deshalb der sinnvolle Mittelweg: Booking bleibt Dirigent f&uuml;r den Hinweg, die Kompensation macht jedes Backend selbst. <strong>Faustregel nach Chris Richardson</strong> (<em>Microservices Patterns</em>, 2018, Kapitel 4): Choreographie f&uuml;r einfache Sagas mit wenigen Beteiligten und ohne Reihenfolge, etwa Benachrichtigungen oder Read-Models. Orchestrierung, sobald es einen Ablauf mit mehreren Schritten oder Verzweigungen gibt, weil sonst die Abh&auml;ngigkeiten zyklisch werden und niemand mehr den Prozess sieht. <strong>Spicy:</strong> Choreographie klingt eleganter, weil kein Service den anderen kennt. In Wahrheit kennt jeder jeden, nur steht es nirgends. Br&uuml;cke zu CQRS am Nachmittag: Dort sind Events am richtigen Platz, weil das Read-Model nur zuh&ouml;rt und niemanden braucht.
- Die fr&uuml;here f&uuml;nfte Frage (Exactly-once) ist in die Bonus-Folie gewandert. Vollst&auml;ndige Antworten: <code>docs/questions/story7.md</code>.
