<div class="page">

<p class="kicker">Saga</p>

<div class="page-head">

## Vorw&auml;rts, und notfalls zur&uuml;ck

<img class="head-figure" src="./assets/saga.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Lokale Transaktionen</h3>
<p>Jeder Service committet nur bei sich. Kein Two-Phase-Commit &uuml;ber Service-Grenzen.</p>
<code>ACID nur lokal</code>
</div>
<div class="card">
<h3>Forward + Kompensation</h3>
<p>Pro Schritt zwei Operationen: buchen, und das fachliche Gegenst&uuml;ck dazu.</p>
<code>POST &harr; DELETE</code>
</div>
<div class="card">
<h3>Orchestrator</h3>
<p>Booking kennt Reihenfolge und Fortschritt. Die Backends bleiben dumm.</p>
<code>Booking wei&szlig; alles</code>
</div>
<div class="card">
<h3>Saga-Status</h3>
<p>Jede Saga hat einen Zustand, und der ist abfragbar.</p>
<code>PENDING &rarr; COMPLETED | FAILED</code>
</div>
<div class="card">
<h3>Eventual Consistency</h3>
<p>Zwischenzust&auml;nde sind sichtbar. Die Kompensation muss letztlich gelingen.</p>
<code>nicht atomar</code>
</div>
</div>

<div class="market">
<h4>Am Markt</h4>
<div class="pills">
<span class="pill brand">Temporal</span>
<span class="pill">Cadence</span>
<span class="pill">Camunda 8</span>
<span class="pill">Axon Framework</span>
<span class="pill">AWS Step Functions</span>
<span class="pill">Eventuate Tram</span>
</div>
</div>

</div>

</div>

Note:
- <strong>Lokale Transaktionen:</strong> Klassische DB-Transaktionen enden an der Service-Grenze. Two-Phase-Commit w&auml;re theoretisch m&ouml;glich, ist in Microservice-Stacks aber praktisch tot: fragil, langsam, sperrt Ressourcen. Stattdessen eine Kette lokaler Transaktionen, das Gesamtergebnis ist eventually consistent.
- <strong>Forward + Kompensation:</strong> Der zentrale Trick: F&uuml;r jeden Vorw&auml;rts-Schritt eine semantisch inverse Operation. Das ist nicht zwingend ein technisches Rollback, oft eine fachliche Gegenbuchung: Refund, Gutschein, Storno. Starke Annahme: Forward darf scheitern, Kompensation <em>muss</em> letztlich gelingen.
- <strong>Orchestrator:</strong> Booking kennt Reihenfolge, was schon aufgerufen wurde, was kompensiert werden muss, aktuellen Status. Hotel, Flight, Car kennen nur ihre eigene lokale Transaktion. Vorteil: Ein Bug in der Saga-Logik steckt an einer Stelle. Die Alternative, Choreography, kommt in Story 7.
- <strong>Saga-Status:</strong> <code>PENDING</code>, dann <code>COMPLETED</code>, oder bei Fehler <code>COMPENSATING</code> und <code>FAILED</code>. In Produktion muss der Status persistent sein, sonst stirbt die Saga mit dem Prozess. Im Workshop bewusst in-memory, Recap-Frage 5.
- <strong>Eventual Consistency:</strong> Wer Saga nutzt, gibt Atomicity auf. F&uuml;r kurze Momente ist der Flug gebucht und das Hotel nicht. Wer das nicht durchdenkt (UI, E-Mail, Reporting), zeigt dem Kunden inkonsistente Daten. Saga ist ein bewusster Trade-off, kein Allheilmittel.
- Am Markt: Temporal und Cadence sind heute der Standard f&uuml;r &bdquo;Saga as Code&ldquo;, die Engine &uuml;bernimmt State, Retry und Recovery. Camunda 8 in BPMN-orientierten Enterprise-Umgebungen, AWS Step Functions als Managed-Variante, Axon und Eventuate Tram als Saga-Libraries im Java-Stack. MassTransit und NServiceBus das Gegenst&uuml;ck in .NET.
- Take-away: Saga ist <em>nicht</em> Eventing. Saga kann synchron &uuml;ber HTTP laufen (Story 6) oder asynchron &uuml;ber Events (Story 7). Das eine ist das Pattern, das andere die Transportwahl.
- &Uuml;berleitung: &bdquo;Die Mechanik passt auf eine Folie.&ldquo;
