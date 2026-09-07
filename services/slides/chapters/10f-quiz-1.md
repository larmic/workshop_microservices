## Quiz 1/3

<p class="subtitle">RESTful oder nicht?</p>

<div class="box">

### <code>GET /booking/cancelBooking?id=4711</code>

</div>

<div class="box fragment">

<strong>Nicht RESTful.</strong> Verb im Pfad, Identit&auml;t als Query-Parameter, und ein GET mit Seiteneffekt.

Besser: <code>DELETE /booking/bookings/4711</code> oder <code>POST /booking/bookings/4711/cancellation</code>

</div>

Note:
- Ablauf: Endpoint zeigen, Handzeichen &bdquo;RESTful?&ldquo; abfragen, eine Person aus der Minderheit begr&uuml;nden lassen, dann Fragment aufl&ouml;sen.
- Erwartung: fast alle sagen &bdquo;nicht RESTful&ldquo;. Das ist der Aufw&auml;rmer, damit die Regeln sitzen.
- Anekdote: 2005 hat der Google Web Accelerator Links auf Webseiten vorgeladen, um Seiten schneller zu machen. Bei der 37signals-Anwendung Backpack waren &bdquo;L&ouml;schen&ldquo;-Links einfache GET-Links. Der Prefetcher hat Nutzern ihre Daten gel&ouml;scht, ohne dass jemand geklickt hatte. Seitdem ist &bdquo;GET ver&auml;ndert nichts&ldquo; Selbstschutz, keine Stilfrage.
- Drei Fehler in einer Zeile benennen: Verb (<code>cancelBooking</code>), Identit&auml;t als Query statt Pfad, Seiteneffekt per GET.
- &Uuml;berleitung: &bdquo;Das war leicht. Jetzt wird es gemeiner.&ldquo;
