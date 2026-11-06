<div class="page">

<p class="kicker">Quiz &middot; 5 von 5</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">POST</span> /booking/bookings<br><span class="dim">&rarr; 418 I&rsquo;m a teapot</span></div>

<div class="quiz-answer fragment">
<div class="callout">Nicht RESTful</div>
<p>Der Code existiert wirklich, seit dem Aprilscherz-RFC 2324 von 1998. Load Balancer, Monitoring und Circuit Breaker kennen ihn nicht und lesen ihn wie 400: Client-Fehler, nicht wiederholen. Was eigentlich schiefging, erf&auml;hrt niemand.</p>
<code class="quiz-take">witzig &ne; Vertrag</code>
</div>

</div>

</div>

Note:
- <strong>Mechanik:</strong> Daumen wie gehabt. Den Lacher mitnehmen, dann aufl&ouml;sen.
- <strong>Kernbotschaft:</strong> Status-Codes sind Infrastruktur, kein Freitext. Wer einen Code w&auml;hlt, w&auml;hlt, was Load Balancer, Retry und Alerting damit tun.
- <strong>Richtig je nach Fall:</strong> <code>409</code> bei Konflikt, <code>422</code> bei fachlich ung&uuml;ltiger Anfrage, <code>503</code>, wenn ein Backend weg ist und ein Retry Sinn ergibt.
- <strong>&Uuml;berleitung:</strong> In Story 3 baut ihr <code>POST /booking/bookings</code> selbst, mit <code>201</code> und Fehler-Codes, die etwas bedeuten. Vorher die Frage: Wie findet der Booking-Service Flight, Hotel und Car?
