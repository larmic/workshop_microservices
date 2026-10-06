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
- Hook direkt aus Story 6: &bdquo;Erinnert ihr euch an die Saga-Frage, was passiert, wenn der <code>DELETE</code> selbst in 5xx l&auml;uft? Story 6 lie&szlig; es scheitern. Story 7 schiebt das Problem ins Backend, und rei&szlig;t damit ein neues auf.&ldquo;
- <strong>Problem aus Story 6:</strong> Booking war Orchestrator: <code>DELETE</code> gegen jedes Backend, eigene Retry-Schleife, eigener Saga-State, eigene Antwort zum Kunden. Wenn Hotel kurz weg ist, blockt Booking. Es ist Aggregator und Stornierstelle in einem.
- <strong>Verantwortung verlagern:</strong> Hotel kennt seinen Buchungsbestand, seine Storno-Sonderregeln, seine eigenen Abh&auml;ngigkeiten. Diese Logik bei Booking zu zentralisieren ist eine k&uuml;nstliche Schicht zwischen Hotel und sich selbst.
- <strong>Event statt Aufruf:</strong> Statt synchron auf jede Antwort zu warten, schickt Booking einen Hinweis: &bdquo;Diese Buchung soll storniert werden.&ldquo; Das Backend nimmt ihn entgegen, best&auml;tigt sofort und macht den Rollback asynchron.
- <strong>Booking bleibt verantwortlich:</strong> Der h&auml;ufigste Denkfehler: &bdquo;Event raus, Booking ist fertig.&ldquo; Der Kunde bekommt seine Antwort, aber Booking wei&szlig; nicht, ob der Rollback im Backend geklappt hat. Reply-Events, Timeout-Erkennung und ein <code>STUCK</code>-Status w&auml;ren der n&auml;chste Schritt. Kommt im Spiel als St&ouml;rung dran.
- <strong>Trade-off:</strong> Orchestration konzentriert das Saga-Wissen an einer Stelle. Choreography verteilt es: Ein Bug in der Saga-Logik kann jetzt in jedem Service sitzen. Ohne klare Event-Vertr&auml;ge entsteht der verteilte Monolith, die schlimmste beider Welten.
- Am Markt: Kafka als Log-orientierter Broker (Persistenz, Replay), RabbitMQ als klassischer Queue-Broker, NATS als leichtgewichtige Alternative, SNS plus SQS und Pub/Sub als Managed-Varianten. Das Outbox-Pattern ist kein Broker, sondern die Br&uuml;cke vom Datenbank-State zum Broker, Pflichtlekt&uuml;re.
- &Uuml;berleitung: &bdquo;Wie sich das anf&uuml;hlt, wenn niemand mehr dirigiert, spielen wir jetzt durch.&ldquo;
