<div class="page">

<p class="kicker">Bonus zu Frage 03</p>

## Was ein Broker kann

<p class="subtitle">&hellip; und ein Webhook nicht.</p>

<div class="page-body bonus">

<div class="swap">
<div class="cards cards-3 ghost fragment fade-out" data-fragment-index="1">
<div class="card">?</div>
<div class="card">?</div>
<div class="card">?</div>
</div>
<div class="cards cards-3 violet compact fragment" data-fragment-index="1">
<div class="card">
<h3>Behalten</h3>
<p>Das Event liegt auf Platte und &uuml;berlebt den Absturz von Sender und Empf&auml;nger. Wer sp&auml;ter kommt, bekommt es trotzdem.</p>
<code>Persistenz &middot; Redelivery &middot; Replay</code>
</div>
<div class="card">
<h3>Verteilen</h3>
<p>Beliebig viele Empf&auml;nger pro Thema, dynamisch dazukommen lassen. Ein langsamer Empf&auml;nger bremst den Sender nicht.</p>
<code>Fan-out &middot; Backpressure &middot; Ordering</code>
</div>
<div class="card">
<h3>Scheitern lassen</h3>
<p>Was nach N Versuchen nicht zustellbar ist, landet sichtbar in einer Inbox statt im Nichts.</p>
<code>Dead-Letter-Queue</code>
</div>
</div>
</div>

<p class="bonus-text fragment" data-fragment-index="1">Exactly-once gibt es trotzdem nicht. Die praktische N&auml;herung: at-least-once beim Sender, Idempotenz beim Empf&auml;nger.</p>

<div class="foot fragment" data-fragment-index="1">
<div class="callout">Wer Eventing ohne Broker will, baut sich in zw&ouml;lf Monaten einen schlechten.</div>
</div>

</div>

</div>

Note:
- Kurze Diskussion: Erst raten lassen, was der Broker leistet, dann per Klick aufdecken.
- <strong>Behalten:</strong> Persistenz, Redelivery mit at-least-once, Replay vom alten Offset, etwa um ein Read-Model neu aufzubauen. Dazu das Outbox-Pattern als Br&uuml;cke: Saga-State und Event in derselben Datenbank-Transaktion, ein Worker publiziert mit Retry. Sonst gibt es State ohne Event oder Event ohne State.
- <strong>Verteilen:</strong> Fan-out per Topic, Backpressure durch die Queue, Ordering pro Partition. Ein Webhook kennt genau einen Empf&auml;nger und blockiert, wenn der langsam ist.
- <strong>Scheitern lassen:</strong> Dead-Letter-Queue als Operator-Inbox. Beim Webhook verschwindet ein nicht zustellbares Event im Log, wenn &uuml;berhaupt.
- <strong>Exactly-once:</strong> Streng genommen nicht erreichbar ohne globale Koordination. Kafka bietet eine Exactly-once-Semantik an, sie ist teuer (transaktionale Producer, viel Koordination) und deckt nur Producer bis Broker ab, nicht den Empf&auml;nger. Was man wirklich will: at-least-once Delivery plus idempotente Verarbeitung. Nicht hipper, aber machbar.
- Das Antipattern: &bdquo;Wir wollen Eventing, aber keinen Broker betreiben.&ldquo; Sechs Monate sp&auml;ter: Outbox-Tabelle, Retry-Worker, Dedup-Logik, Dead-Letter-Inbox, Replay-Skript, alles selbst gebaut. Zw&ouml;lf Monate sp&auml;ter: ein schlechter Broker.
- Br&uuml;cke zu Story 8: Bei Events wandert der HTTP-Header nicht mit. Die Trace-ID muss als Property ins Event, sonst rei&szlig;t der rote Faden genau an der Stelle, wo es leise wird.
