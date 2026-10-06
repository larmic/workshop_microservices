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
<p class="numlist-label">Der Zettel kommt nie an</p>
<h3>Booking schickt per HTTP-POST, Hotel ist down. Wo ist das Event <span class="hl">jetzt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Event ist weg</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Der Zettel kommt zweimal</p>
<h3>At-least-once ist Standard. Was, wenn das Backend echten State <span class="hl">&auml;ndert</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">eventId als Dedup-Key</code>
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
<p class="numlist-label">Wo ist das Wissen hin?</p>
<h3>Jeder reagiert auf Events der anderen. Wann wird das zum <span class="hl">verteilten Monolithen</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Orchestration &ne; Choreography</code>
</div>
</div>

</div>

</div>

Note:
- Alle vier Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein. Frage 03 hat eine Bonus-Folie darunter.
- <strong>Zettel kommt nie an, meine Antwort:</strong> In der Webhook-Referenz passiert genau nichts. Der Fehler wird geloggt, der Step trotzdem als <code>COMPENSATED</code> markiert, die Saga geht auf <code>FAILED</code>, der Kunde bekommt seine Antwort. Booking sagt &bdquo;Event ist raus&ldquo;, in Wahrheit wurde es nie empfangen. Lehrbuch-Antworten: Retry mit Backoff, Outbox-Pattern (Event in derselben Transaktion wie der Saga-State, ein Worker publiziert), Dead-Letter-Queue, ein Broker mit at-least-once. <strong>Spicy:</strong> Eventing eliminiert das &bdquo;Backend kurz weg&ldquo;-Problem nicht, es verschiebt es. In Story 6 hat Booking den Schmerz gesp&uuml;rt. In Story 7 sieht Booking gar nichts.
- <strong>Zettel kommt zweimal, meine Antwort:</strong> Unsere Referenz ist idempotent durch Zustandslosigkeit, das Backend loggt nur. In Produktion ist das die Ausnahme. Doppelte Zustellung passiert, wenn die Antwort auf den POST verloren geht und der Sender wiederholt, wenn der Broker nach einem Consumer-Crash erneut liefert, wenn der Outbox-Worker zwischen Senden und Markieren abst&uuml;rzt. Ohne Dedup: R&uuml;ckerstattung zweimal, Reservierung doppelt freigegeben. Standard: <code>eventId</code> in einer kleinen Tabelle merken, in derselben Transaktion wie die Fachlogik. <strong>Spicy:</strong> Wer Events publiziert, muss davon ausgehen, dass sie mehrfach ankommen. Wer das ignoriert, hat einen Bug, der erst unter Last zuschl&auml;gt.
- <strong>Webhooks reichen doch, meine Antwort:</strong> F&uuml;r genau die Dinge, die der Webhook strukturell nicht kann: Persistenz, Redelivery, Fan-out, Backpressure, Dead-Letter, Consumer-Crash-Recovery, Ordering, Replay, Entkopplung in der Zeit. Webhook gewinnt nur beim Infrastruktur-Aufwand, und genau darauf fallen Teams herein. Details auf der Bonus-Folie.
- <strong>Wo ist das Wissen hin, meine Antwort:</strong> In Story 6 lag das Saga-Wissen bei Booking. In Story 7 reagiert jeder Service auf Events anderer, ein Bug kann &uuml;berall sitzen. Zum verteilten Monolithen wird es, wenn niemand mehr wei&szlig;, wer auf welches Event wie reagiert, wenn Event-Vertr&auml;ge nicht versioniert sind, wenn eine &Auml;nderung drei andere Services mitzieht. Gegenmittel: Event-Katalog, Schema-Registry, klare Ownership pro Event-Typ. <strong>Spicy:</strong> Beides ist legitim. Wissen verteilen ohne Plan ergibt die schlimmste beider Welten: lose gekoppelt zur Laufzeit, fest gekoppelt zur Build-Zeit.
- Die fr&uuml;here f&uuml;nfte Frage (Exactly-once) ist in die Bonus-Folie gewandert. Vollst&auml;ndige Antworten: <code>docs/questions/story7.md</code>.
