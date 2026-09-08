<div class="page">

<p class="kicker">Quiz &middot; 5 von 5</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">POST</span> /booking/bookings<br><span class="dim">&rarr; 418 I&rsquo;m a teapot</span></div>

<div class="quiz-answer fragment">
<div class="callout">Nicht RESTful</div>
<p>Der Code existiert wirklich, seit dem Aprilscherz-RFC 2324 von 1998. Nur bedeutet er f&uuml;r Load Balancer, Monitoring und Circuit Breaker nichts. Ein Status-Code ist eine Handlungsanweisung an die Infrastruktur: wiederholen, aufgeben, cachen. Eine Teekanne gibt keine.</p>
<code class="quiz-take">witzig &ne; Vertrag</code>
</div>

</div>

</div>

Note:
- Auflockerung zum Schluss, Lacher einsammeln, dann die Pointe: Auch ein g&uuml;ltiger Code ist falsch, wenn er der Infrastruktur nichts sagt. Genau deshalb steht Regel 6 auf der Folie.
- Anekdote: 418 stammt aus RFC 2324, dem Hyper Text Coffee Pot Control Protocol, ver&ouml;ffentlicht am 1. April 1998. 2017 wollte der Vorsitzende der HTTP-Arbeitsgruppe den Code aus Node.js, Go und ASP.NET entfernen lassen. Die Community hat mit der &bdquo;Save 418&ldquo;-Kampagne dagegengehalten, der Code ist seitdem offiziell reserviert, damit ihn niemand ernsthaft belegt. Selbst der Scherz wird im Standard sauber verwaltet.
- Richtig w&auml;re je nach Fall: <code>409</code> (Konflikt), <code>422</code> (fachlich abgelehnt), <code>503</code> (Backend weg). Was es nicht sein darf: ein Code, den kein Client interpretieren kann, und erst recht kein <code>200</code> mit Fehler im Body.
- &Uuml;berleitung: &bdquo;Genug Theorie. In Story 3 baut ihr <code>POST /booking/bookings</code> selbst, mit einem Status-Code, der etwas bedeutet.&ldquo;
