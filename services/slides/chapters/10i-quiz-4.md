<div class="page">

<p class="kicker">Quiz &middot; 4 von 5</p>

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

