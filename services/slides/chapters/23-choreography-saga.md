<div class="page">

<p class="kicker">Choreography-Saga</p>

<div class="page-head">

## Die Saga wird leise

<img class="head-figure" src="./assets/choreography.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Problem aus Story 6</h3>
<p>Booking tr&auml;gt die Storno-Verantwortung allein und h&auml;ngt am kaputten Backend.</p>
<code>Booking = SPoF</code>
</div>
<div class="card">
<h3>Verantwortung verlagern</h3>
<p>Wer buchen kann, kann auch stornieren. Jedes Backend macht seinen eigenen Rollback.</p>
<code>eigene lokale Transaktion</code>
</div>
<div class="card">
<h3>Event statt Aufruf</h3>
<p>Booking publiziert ein Event und ist mit dem Befehl fertig.</p>
<code>fire-and-forget</code>
</div>
<div class="card">
<h3>Booking bleibt haftbar</h3>
<p>Saga-Status und Kundenantwort bleiben bei Booking. Ob der Rollback gelang, wei&szlig; es nicht.</p>
<code>Event raus &ne; storniert</code>
</div>
<div class="card">
<h3>Trade-off</h3>
<p>Gewinn: Entkopplung. Verlust: Das Wissen ist jetzt &uuml;berall.</p>
<code>verteilter Monolith</code>
</div>
</div>

<div class="market">
<h4>Am Markt</h4>
<div class="pills">
<span class="pill brand">Apache Kafka</span>
<span class="pill">RabbitMQ</span>
<span class="pill">NATS</span>
<span class="pill">AWS SNS + SQS</span>
<span class="pill">GCP Pub/Sub</span>
<span class="pill">Outbox-Pattern</span>
</div>
</div>

</div>

</div>

Note:
- Hook: Recap 6/01: Was, wenn der <code>DELETE</code> in 5xx l&auml;uft?
- <strong>Problem:</strong> Booking orchestriert alles allein. Ist Hotel weg, blockt Booking.
- <strong>Verlagern:</strong> Hotel kennt seine Storno-Regeln selbst.
- <strong>Event:</strong> Backend best&auml;tigt sofort, storniert asynchron.
- <strong>Haftbar:</strong> Booking wei&szlig; nicht, ob der Rollback klappte. L&ouml;sung (Reply-Events, <code>STUCK</code>) kommt im Spiel als St&ouml;rung.
- <strong>Trade-off:</strong> Saga-Wissen in jedem Service. Ohne Event-Vertr&auml;ge droht der verteilte Monolith.
- <strong>Ehrlich:</strong> <code>CompensationRequested</code> ist ein Command, Booking entscheidet, wer storniert. Echte Choreografie: <code>TripBookingFailed</code>, jedes Backend entscheidet selbst.
- &Uuml;berleitung: &bdquo;Wie sich das anf&uuml;hlt, wenn Booking nur Zettel hinlegt, spielen wir jetzt durch.&ldquo;
