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

