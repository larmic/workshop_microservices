## Choreography-Saga

<p class="subtitle">Die Saga wird leise</p>

<div class="factor-row">

<div class="factor fragment">
<h3>Problem aus Story 5</h3>
<p>Booking tr&auml;gt allein die Storno-Verantwortung &mdash; und h&auml;ngt synchron am kaputten Backend.</p>
<code>Booking = SPoF</code>
<aside class="notes">Story 5 hatte Booking als Orchestrator: <code>DELETE</code> gegen jeden Backend-Service, eigene Retry-Schleife, eigener Saga-State, eigene Antwort zum Kunden. Wenn Hotel kurz weg ist, blockt Booking. Wenn der <code>DELETE</code>-Call in 5xx l&auml;uft, hat Booking die Wahl: weiterh&auml;ngen, eskalieren, Pivot. Alle Last liegt bei Booking. Es ist Aggregator <em>und</em> Stornier-Stelle in einem.</aside>
</div>

<div class="factor fragment">
<h3>Verantwortung verlagern</h3>
<p>Wer buchen kann, kann auch stornieren. Backends &uuml;bernehmen ihren eigenen Rollback.</p>
<code>eigene lokale TX</code>
<aside class="notes">Fachlich ist die Stornierung Aufgabe des jeweiligen Backends &mdash; nicht des Orchestrators. Hotel kennt seinen Buchungsbestand am besten, kennt seine Stornierungs-Sonderregeln, kennt seine eigenen Downstream-Abh&auml;ngigkeiten. Diese Logik bei Booking zu zentralisieren ist eine k&uuml;nstliche Schicht zwischen Hotel und sich selbst.</aside>
</div>

<div class="factor fragment">
<h3>Event statt RPC</h3>
<p>Booking publiziert <code>CompensationRequested</code> und ist mit dem <em>Befehl</em> fertig.</p>
<code>fire-and-forget?</code>
<aside class="notes">Statt synchron auf jede DELETE-Antwort zu warten, schickt Booking einen <em>Hinweis</em> raus: &bdquo;Diese Buchung soll storniert werden.&ldquo; Backend zieht das Event aus dem Bus (in unserer Workshop-Variante: per Webhook), macht seinen Rollback und antwortet asynchron. Booking ist von der Backend-Verf&uuml;gbarkeit entkoppelt.</aside>
</div>

<div class="factor fragment">
<h3>Booking bleibt verantwortlich</h3>
<p>Reply-Events + Timeout &mdash; der <span class="hl">Gesamtstatus</span> gegen&uuml;ber dem Kunden bleibt bei Booking.</p>
<code>nicht fire-and-forget</code>
<aside class="notes">Der h&auml;ufigste Denkfehler: &bdquo;Event raus &rarr; Booking ist fertig.&ldquo; Falsch. Der Kunde will am Ende eine Antwort &mdash; &bdquo;Reise gebucht? storniert? steckt fest?&ldquo;. Backends m&uuml;ssen mit <code>BookingCancelled</code> / <code>CancellationFailed</code> antworten. Booking konsumiert die Replies und f&uuml;hrt den Saga-Status nach. Plus Timeout-Erkennung: bleibt ein Reply nach X Sekunden aus &rarr; <code>STUCK</code>, Operator eingreifen.</aside>
</div>

<div class="factor fragment">
<h3>Trade-off</h3>
<p>Gewinn: Entkopplung. Verlust: Wissen ist jetzt <span class="hl">&uuml;berall</span>.</p>
<code>verteilter Monolith</code>
<aside class="notes">Orchestration konzentriert das Saga-Wissen an einer Stelle (Booking). Choreography verteilt es: Booking publiziert, Backends reagieren auf Events anderer. Ein Bug in der Saga-Logik kann jetzt in jedem beteiligten Service sitzen. Ohne klare Event-Vertr&auml;ge und Dokumentation hat man <em>verteilten Monolithen</em> &mdash; die schlimmste beider Welten. Choreography ist eleganter, aber operativ aufwendiger.</aside>
</div>

</div>

<div class="market-row">

### Am Markt (Broker)

<div class="chip-row">
  <span class="chip brand">Apache Kafka</span>
  <span class="chip">RabbitMQ</span>
  <span class="chip">NATS</span>
  <span class="chip">AWS SNS + SQS</span>
  <span class="chip">GCP Pub/Sub</span>
  <span class="chip">Solace</span>
  <span class="chip">Eventuate Tram</span>
  <span class="chip">Outbox-Pattern</span>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>

Note:
- Hook direkt aus Story 5: &bdquo;Erinnert ihr euch an die Saga-Frage 1 &mdash; was, wenn der <code>DELETE</code> selbst in 5xx l&auml;uft? Story 5 lie&szlig; es scheitern. Story 6 schiebt das Problem ins Backend &mdash; und reisst damit ein neues auf, das wir gleich besprechen.&ldquo;
- Wichtiger Framing-Punkt: Choreography ist <em>kein neues Pattern</em>, sondern eine andere Verantwortungsverteilung f&uuml;r dieselbe Saga. Forward bleibt sync (in unserer Variante), Kompensation wandert auf Events.
- Karten-Reihenfolge bewusst: Problem (Story 5) &rarr; L&ouml;sungsidee (Verantwortung verlagern, Event statt RPC) &rarr; ehrliches Caveat (Booking bleibt verantwortlich) &rarr; Trade-off (verteiltes Wissen).
- Wer was wo nutzt: Kafka als Log-orientierter Broker (Persistenz, Replay), RabbitMQ als klassischer Queue-Broker, NATS als leichtgewichtige Alternative, AWS SNS+SQS / GCP Pub/Sub als managed. Outbox-Pattern ist kein Broker, sondern die Br&uuml;cke vom App-State zum Broker &mdash; Pflicht-Lekt&uuml;re.
- Im Workshop nutzen wir <strong>keinen</strong> echten Broker, sondern HTTP-Webhooks. Das ist eine bewusste Vereinfachung &mdash; die n&auml;chste Folie zeigt, wo unsere Variante an ihre strukturellen Grenzen st&ouml;sst.
