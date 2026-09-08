<div class="page">

<p class="kicker">Quiz &middot; 2 von 5</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">GET</span> /booking/customers/7/bookings?status=confirmed<br><span class="dim">&rarr; 404 Not Found</span></div>

<div class="quiz-answer fragment">
<div class="callout">Nicht RESTful</div>
<p>Kunde 7 hat gerade keine best&auml;tigte Buchung. Die Sammlung ist leer, nicht weg. Richtig ist <code>200 OK</code> mit <code>[]</code>. Ein <code>404</code> sagt dem Client, dass es diesen Pfad nicht gibt, und der sucht dann an der falschen Stelle.</p>
<code class="quiz-take">leer &ne; weg</code>
</div>

</div>

</div>

Note:
- Erwartung: viele sagen &bdquo;RESTful, Pfad und Methode stimmen doch&ldquo;. Wer &bdquo;nein&ldquo; sagt, hat die zweite Zeile gelesen.
- 404 bedeutet: die Ressource existiert nicht. Die Sammlung <code>bookings</code> von Kunde 7 existiert, sie ist nur leer. Leere Sammlung = <code>200</code> mit leerem Array. 404 w&auml;re richtig, wenn Kunde 7 selbst nicht existiert.
- Wem tut das weh? Clients, die bei 404 den Pfad f&uuml;r falsch halten und Fallbacks ziehen, Monitoring mit Fehlerrate, Caches, die &bdquo;nicht vorhanden&ldquo; merken.
- Br&uuml;cke zu Regel 6: Status-Codes sind Vertrag mit der Infrastruktur, keine Kosmetik.
