<div class="page">

<p class="kicker">Quiz &middot; 3 von 4</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">GET</span> /booking/cancelBooking?id=4711</div>

<div class="quiz-answer fragment">
<div class="callout">Nicht RESTful</div>
<p>Verb im Pfad, Identit&auml;t als Query-Parameter, und ein GET, das etwas ver&auml;ndert: drei Regeln auf einmal. Wie es stattdessen aussieht, habt ihr gerade selbst entworfen.</p>
<code class="quiz-take">Aktion &ne; Pfad</code>
</div>

</div>

</div>

Note:
- Erwartung: fast alle sagen &bdquo;nicht RESTful&ldquo;. Drei Fehler in einer Zeile benennen lassen: Verb (<code>cancelBooking</code>), Identit&auml;t als Query statt Pfad, Seiteneffekt per GET.
- Anekdote: 2005 hat der Google Web Accelerator Links auf Webseiten vorgeladen, um Seiten schneller zu machen. Bei der 37signals-Anwendung Backpack waren &bdquo;L&ouml;schen&ldquo;-Links einfache GET-Links. Der Prefetcher hat Nutzern ihre Daten gel&ouml;scht, ohne dass jemand geklickt hatte. Seitdem ist &bdquo;GET ver&auml;ndert nichts&ldquo; Selbstschutz, keine Stilfrage.
- Auf die Flipcharts zeigen: <code>DELETE /booking/bookings/4711</code> oder <code>POST /booking/bookings/4711/cancellation</code>, je nachdem, ob der Grund mit soll.
- Br&uuml;cke nach vorn: In Story 3 baut ihr <code>POST /booking/bookings</code>. Nehmt die Regeln mit. In Story 6 (Saga) kommt der Storno als Kompensation zur&uuml;ck, dann z&auml;hlt die Idempotenz.
- &Uuml;berleitung: &bdquo;Und wie sieht es richtig aus? Zwei Varianten, beide von euren Flipcharts.&ldquo;
- Reserve, falls Zeit bleibt: <code>POST /flights/search</code> (Grauzone), <code>PUT</code> mit Teil-Objekt, <code>POST /admin/bulkhead-reset</code> aus unserer Referenz (bewusst RPC-artig). Alle in <code>docs/instructions/rest-vs-restful.md</code>, Abschnitt 7. Vollst&auml;ndige Antworten zur Design-Session: <code>docs/questions/story2.md</code>.
