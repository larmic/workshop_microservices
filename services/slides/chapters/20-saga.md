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
<h3>Buchen + Kompensation</h3>
<p>Pro Schritt zwei Operationen: buchen, und das fachliche Gegenst&uuml;ck dazu.</p>
<code>POST &harr; DELETE</code>
</div>
<div class="card">
<h3>Orchestrator</h3>
<p>Booking kennt Reihenfolge und Fortschritt. Die andere Variante, Choreografie, kommt in Story 7.</p>
<code>Variante 1 von 2</code>
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
<div class="card market-card">
<h3>Am Markt</h3>
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

</div>

Note:
- <strong>Lokal:</strong> 2PC ist praktisch tot, fragil, langsam und sperrend.
  - Geteilte DB hei&szlig;t geteiltes Schema und gemeinsames Deployment, also verteilter Monolith.
  - Im Workshop in-memory, bewusst reduziert.
- <strong>Kompensation:</strong> fachliche Gegenbuchung (Storno, Gutschein), kein technisches Rollback.
  - Buchen darf scheitern, die Kompensation <em>muss</em> gelingen.
  - Scheitert ein Schritt: r&uuml;ckw&auml;rts kompensieren (Backward Recovery) oder bei technischem Fehler wie <code>503</code> wiederholen (Forward Recovery). Unsere Referenz kann nur r&uuml;ckw&auml;rts, Story 7 spielt beides.
- <strong>Orchestrator:</strong> eine von zwei Arten, die Saga zu koordinieren (Chris Richardson). Die Saga-Logik an einer Stelle, die Downstreams bleiben dumm. Choreografie, ohne Dirigent, kommt in Story 7.
- <strong>Status:</strong> bei Fehler <code>COMPENSATING</code>, dann <code>FAILED</code>. Produktiv persistent (Recap 04).
- <strong>Eventual Consistency:</strong> Flug gebucht, Hotel noch nicht. UI und E-Mail m&uuml;ssen das aushalten.
- <strong>Markt:</strong> f&uuml;r .NET zus&auml;tzlich MassTransit und NServiceBus.
- <strong>Spicy:</strong> Saga ist nicht Eventing. Ob sync oder async, ist eine Frage des Transports.
