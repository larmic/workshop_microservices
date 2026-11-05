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

Note:
- <strong>Mechanik:</strong> Daumen wie gehabt, eine Stimme aus der Minderheit begr&uuml;nden lassen, dann aufl&ouml;sen.
- <strong>Erwartung:</strong> Viele sagen &bdquo;nur DELETE, im POST steckt eine Aktion&ldquo;. Aber <code>cancellation</code> ist ein Substantiv, also eine Ressource.
- <strong>Entscheidungsfrage:</strong> Muss der Storno-Grund mit? Dann braucht es eine Ressource, die ihn tr&auml;gt. Direkter Bezug zu eurem Flipchart.
- <strong>Ausblick Story 6:</strong> Die Kompensation ruft <code>DELETE /bookings/{id}</code> auf. Das ist kein Kunden-Storno: Die Reservierung wird zur&uuml;ckgenommen, bevor die Kundin sie best&auml;tigt bekommt. Dann z&auml;hlt die Idempotenz. F&uuml;r den Storno einer fertigen Reise bleibt <code>POST &hellip;/cancellation</code> richtig.
- <strong>Weiter:</strong> Letztes Quiz, ein Status-Code.
