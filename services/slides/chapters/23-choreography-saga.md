<div class="page">

<p class="kicker">Choreography-Saga</p>

<div class="page-head">

## Choreografie: Jeder h&ouml;rt zu

<img class="head-figure" src="./assets/choreography.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Problem aus Story 6</h3>
<p>Booking storniert selbst, wartet auf jedes Ergebnis und muss jeden Storno-Fehler selbst behandeln.</p>
<code>Booking haftet f&uuml;r alles</code>
</div>
<div class="card">
<h3>Verantwortung verlagern</h3>
<p>Wer buchen kann, kann auch stornieren. Jedes Backend kompensiert selbst.</p>
<code>eigene lokale Transaktion</code>
</div>
<div class="card">
<h3>Event statt Aufruf</h3>
<p>Booking meldet das Storno nur noch an. Ausf&uuml;hren muss es das Backend.</p>
<code>wartet nur auf die Annahme</code>
</div>
<div class="card">
<h3>Booking bleibt haftbar</h3>
<p>Saga-Status und Kundenantwort bleiben bei Booking, auch wenn es nicht mehr selbst storniert.</p>
<code>Wer wei&szlig;, ob storniert ist?</code>
</div>
<div class="card">
<h3>Trade-off</h3>
<p>Gewinn: Entkopplung. Verlust: Das Wissen ist jetzt &uuml;berall.</p>
<code>verteilter Monolith</code>
</div>
<div class="card market-card">
<h3>Am Markt</h3>
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

</div>

Note:
- Hook: Recap 6/01: Was, wenn der <code>DELETE</code> in 5xx l&auml;uft?
- <strong>Problem:</strong> Booking orchestriert alles allein. Scheitert ein <code>DELETE</code>, muss Booking das selbst ausbaden (<code>COMPENSATION_FAILED</code>).
- <strong>Verlagern:</strong> Hotel kennt seine Storno-Regeln selbst.
- <strong>Event:</strong> Backend nimmt an (202), storniert danach selbst (in der Referenz: loggt nur). Booking wartet nur auf die Annahme, an der Erreichbarkeit h&auml;ngt es weiter (3-s-Timeout).
- <strong>Haftbar:</strong> Die Frage auf der Karte offen lassen, das Spiel beantwortet sie. <strong>Trade-off:</strong> Saga-Wissen in jedem Service. Ohne Event-Vertr&auml;ge droht der verteilte Monolith.
- <strong>Ehrlich:</strong> <code>CompensationRequested</code> ist ein Command, Booking entscheidet, wer storniert. Echte Choreografie: <code>TripBookingFailed</code>, jedes Backend entscheidet selbst. &Uuml;berleitung: &bdquo;Wie sich das anf&uuml;hlt, wenn Booking nur noch zuruft, spielen wir jetzt durch.&ldquo;
