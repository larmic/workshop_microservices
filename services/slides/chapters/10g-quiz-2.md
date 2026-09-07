## Quiz 2/3

<p class="subtitle">RESTful oder nicht?</p>

<div class="box">

### <code>POST /booking/bookings/4711/cancellation</code>

</div>

<div class="box fragment">

<strong>RESTful.</strong> <em>cancellation</em> ist ein Substantiv, also eine Ressource. Der Client legt eine Stornierung an: <code>201 Created</code>, sp&auml;ter per <code>GET</code> nachlesbar (Grund, Zeitpunkt, Geb&uuml;hr).

Oft besser als ein nacktes <code>DELETE</code>, weil die Buchung als Historie bleibt.

</div>

Note:
- Erwartung: viele sagen &bdquo;nein, da steht doch eine Aktion drin&ldquo;. Genau das ist die Pointe: nicht jedes Wort, das nach Handlung klingt, ist ein Verb. Stornierung ist ein Fachobjekt.
- Nachfragen: &bdquo;Was ist beim zweiten POST?&ldquo; Antwort: <code>409 Conflict</code> (schon storniert) oder Idempotency-Key im Header. POST ist nicht idempotent, der Server muss das abfangen.
- Anker aus dem Netz: Stripe modelliert R&uuml;ckerstattungen genauso als eigene Ressource (<code>POST /v1/refunds</code>), mit Idempotency-Key.
- Bezug zur Flipchart-Aufgabe: diese Variante haben meist ein oder zwei Teams gebaut. Aufgreifen und loben, sie l&ouml;st die Anforderung &bdquo;Grund nachvollziehbar&ldquo; am sauberesten.
- Br&uuml;cke nach vorn: In Story 6 (Saga) kommt der Storno als Kompensation zur&uuml;ck. Dort ruft der Orchestrator <code>DELETE /bookings/{id}</code> an den Backends auf, und dann z&auml;hlt die Idempotenz.
