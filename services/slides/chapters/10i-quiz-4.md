<div class="page">

<p class="kicker">Quiz &middot; 4 von 4</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">DELETE</span> /booking/bookings/4711<br><span class="m">POST</span> /booking/bookings/4711/cancellation</div>

<div class="quiz-answer fragment">
<div class="callout">Beides RESTful</div>
<p><code>DELETE</code> ist idempotent und braucht keinen Body, die Buchung ist danach weg. <code>POST</code> legt eine Stornierung als eigene Ressource an: mit Grund, Zeitpunkt und Geb&uuml;hr, sp&auml;ter per <code>GET</code> nachlesbar, die Buchung bleibt als Historie. Der Preis: Der zweite <code>POST</code> braucht <code>409</code> oder einen Idempotency-Key.</p>
<code class="quiz-take">stornieren &ne; l&ouml;schen</code>
</div>

</div>

</div>

Note:
- Erwartung: Die H&auml;lfte sagt &bdquo;nur DELETE ist RESTful, POST hat doch eine Aktion im Pfad&ldquo;. Pointe: <code>cancellation</code> ist ein Substantiv, die Stornierung ist ein Fachobjekt mit eigenen Daten. Nicht jedes Wort, das nach Handlung klingt, ist ein Verb.
- Entscheidungsfrage an die Teams: Muss der Storno-Grund mit? Dann braucht ihr eine Ressource, die ihn tr&auml;gt. DELETE hat daf&uuml;r keinen Platz, und die Buchung &bdquo;verschwindet&ldquo;, bis der Kundenservice fragt, wer das war.
- Nachfragen: &bdquo;Was ist beim zweiten POST?&ldquo; <code>409 Conflict</code> (schon storniert) oder Idempotency-Key im Header, mit dem der Server die Wiederholung erkennt und die erste Antwort noch einmal liefert. So arbeitet Stripe bei R&uuml;ckerstattungen (<code>POST /v1/refunds</code>).
- Dritte Variante, falls sie auf einem Flipchart steht: <code>PATCH</code> auf den Status. Vertretbar, versteckt aber die Zustandsmaschine im Body.
- Br&uuml;cke nach vorn: In Story 6 (Saga) kommt der Storno als Kompensation zur&uuml;ck. Dort ruft der Orchestrator <code>DELETE /bookings/{id}</code> an den Backends auf, und dann z&auml;hlt die Idempotenz.
