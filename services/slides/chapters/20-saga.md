<!-- .slide: data-background-image="./assets/saga.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Saga

<p class="subtitle">Alles oder nichts &mdash; aber richtig</p>

<div class="factor-row">

<div class="factor fragment">
<h3>Lokale Transaktionen</h3>
<p>Jeder Service hat seine eigene Transaktion &mdash; <strong>kein</strong> 2PC &uuml;ber Service-Grenzen.</p>
<code>ACID nur lokal</code>
<aside class="notes">Klassische DB-Transaktionen funktionieren nicht &uuml;ber Service-Grenzen. Zwei-Phasen-Commit ist theoretisch m&ouml;glich, praktisch in modernen Microservice-Stacks tot &mdash; zu fragil, zu langsam, sperrt Ressourcen. Stattdessen: eine Kette lokaler Transaktionen, jeder Service committed lokal, das Gesamtergebnis ist eventually consistent.</aside>
</div>

<div class="factor fragment">
<h3>Forward + Kompensation</h3>
<p>Pro Schritt zwei Operationen: <code>book</code> und <code>cancel</code>.</p>
<code>POST &harr; DELETE</code>
<aside class="notes">Der zentrale Saga-Trick: f&uuml;r jeden Vorw&auml;rts-Schritt eine semantisch inverse Operation. Wichtig: das ist nicht zwingend ein technisches Rollback &mdash; oft eine fachliche Gegenbuchung (Refund, Gutschein, Storno). Saga macht eine starke Annahme: Forward-Steps d&uuml;rfen scheitern, Kompensationen <em>m&uuml;ssen letztlich gelingen</em>. Ohne diese Annahme bricht das Konstrukt zusammen.</aside>
</div>

<div class="factor fragment">
<h3>Orchestrator</h3>
<p>Zentraler Koordinator steuert den Ablauf. Wissen liegt an <em>einer</em> Stelle.</p>
<code>Booking wei&szlig; alles</code>
<aside class="notes">Booking kennt: Reihenfolge der Schritte, was wurde schon aufgerufen, was muss kompensiert werden, aktueller Saga-Status. Hotel / Flight / Car bleiben <em>dumm und einfach</em> &mdash; sie kennen nur ihre eigenen lokalen Transaktionen. Vorteil: ein Bug in der Saga-Logik steckt an <em>einer</em> Stelle. Alternative ist Choreography (Story 6) &mdash; verteiltes Saga-Wissen &uuml;ber Events.</aside>
</div>

<div class="factor fragment">
<h3>Saga-Status</h3>
<p><code>PENDING</code> &rarr; <code>COMPLETED</code> oder <code>COMPENSATING</code> &rarr; <code>FAILED</code>.</p>
<code>idempotent by design</code>
<aside class="notes">Der Status ist die Lebensader der Saga. Pflicht: <strong>idempotente Endpoints</strong> &mdash; mehrfacher <code>DELETE</code> derselben Buchung liefert dasselbe Ergebnis, ohne Fehler. In Produktion muss der Status persistent sein, sonst stirbt die Saga mit dem Prozess. Im Workshop bewusst in-memory, weil 60-Min-Slot &mdash; daf&uuml;r aber explizit besprochen (siehe Recap-Frage 7).</aside>
</div>

<div class="factor fragment">
<h3>Eventual Consistency</h3>
<p>Zwischenzust&auml;nde sind <em>sichtbar</em>. Kompensation muss <span class="hl">letztlich</span> gelingen.</p>
<code>nicht atomar</code>
<aside class="notes">Wer Saga benutzt, gibt Atomicity auf. F&uuml;r kurze Momente ist Flug gebucht, Hotel nicht. Wer das nicht durchdenkt (UI / E-Mail / Reporting), zeigt dem Kunden inkonsistente Daten. Akzeptiere das &mdash; oder w&auml;hle ein anderes Pattern (zentrale DB, eine Datenbank pro Booking-Replica, &hellip;). Saga ist kein magisches Allheilmittel, sondern ein bewusstes Trade-off zwischen Konsistenz und Verteilung.</aside>
</div>

</div>

<div class="market-row">

### Am Markt

<div class="chip-row">
  <span class="chip brand">Temporal</span>
  <span class="chip">Cadence</span>
  <span class="chip">Camunda 8 (Zeebe)</span>
  <span class="chip">Axon Framework</span>
  <span class="chip">AWS Step Functions</span>
  <span class="chip">Eventuate Tram Saga</span>
  <span class="chip">MassTransit Saga</span>
  <span class="chip">NServiceBus Saga</span>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>

Note:
- Hook: &bdquo;In Story 3 haben wir bei POST Fail-Fast gemacht &mdash; wenn ein CB OPEN ist, ganzen Buchungsversuch abbrechen. Aber: was, wenn der Flug schon gebucht ist und <em>dann</em> kippt das Hotel?&ldquo; Demo: Flight normal, Hotel auf &bdquo;Fehler&ldquo;, dann <code>POST /booking/bookings</code> &mdash; im Dashboard ist sichtbar, wie Flight gebucht und dann kompensiert wird.
- Karten-Reihenfolge bewusst: erst die Abgrenzung gegen ACID (Lokale Transaktionen), dann die Mechanik (Forward + Kompensation, Orchestrator, Status), zuletzt das ehrliche Trade-off (Eventual Consistency).
- Wer was wo nutzt: <strong>Temporal / Cadence</strong> sind heute Branchen-Standard f&uuml;r &bdquo;Saga as Code&ldquo; &mdash; Engine k&uuml;mmert sich um State, Retry, Recovery. <strong>Camunda 8</strong> stark in BPMN-orientierten Enterprise-Umgebungen. <strong>AWS Step Functions</strong> als managed Variante. <strong>Eventuate / MassTransit / NServiceBus</strong> sind Saga-Libraries im jeweiligen .NET-/Java-Stack.
- Wichtigster Take-away f&uuml;r die Brain Bridge: Saga ist <em>nicht</em> Eventing. Saga kann sync HTTP sein (Story 5) oder asynchron &uuml;ber Events laufen (Choreography, Story 6). Das eine ist das Pattern, das andere die Transport-Wahl.
- &Uuml;berleitung: Wir schauen uns die Mechanik konkret an &mdash; in genau der Form, wie sie im Dashboard-Spickzettel und im Story-5-Code steckt.
