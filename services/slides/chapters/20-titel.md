<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Saga</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Flug gebucht, Hotel sagt nein. Der Flug muss <span class="accent">wieder weg</span>.</p>
<p class="sub">Keine Transaktion &uuml;ber Service-Grenzen. Stattdessen eine Kette lokaler Schritte, jeder mit seinem Gegenst&uuml;ck.</p>
<p class="source">Hector Garcia-Molina, Kenneth Salem, <em>Sagas</em> (1987): Lange Transaktionen in Teilschritte zerlegen, jeder mit einer Kompensation.</p>
</div>
<div class="split-figure">
<img src="./assets/saga.svg" alt="Booking ruft Flight, Hotel und Car nacheinander auf, Hotel scheitert, ein Pfeil f&uuml;hrt von Hotel zur&uuml;ck zu Flight"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;Stories 4 und 5 sch&uuml;tzen <em>einen</em> Aufruf. Aber sobald ich mehrere Schritte habe und einer kippt mitten drin, brauche ich etwas anderes. Das Hotel sagt nein, der Flug ist schon gebucht. Was nun?&ldquo;
- Die Grafik: Booking links ruft Flight, Hotel und Car nacheinander auf. Flight ist gebucht, Hotel scheitert, Car wird gar nicht mehr gefragt. Der Lime-Pfeil unten ist die Kompensation: zur&uuml;ck zu Flight, Buchung weg.
- Garcia-Molina und Salem haben das Pattern 1987 f&uuml;r langlaufende Datenbank-Transaktionen beschrieben, lange vor Microservices. Die Idee: Statt eine lange Transaktion zu sperren, in Teilschritte zerlegen, jeder committet f&uuml;r sich, und f&uuml;r jeden Schritt gibt es eine Kompensation.
- Demo-Vorschau: Im Dashboard Hotel auf &bdquo;Fehler&ldquo;, dann <code>POST /booking/bookings</code>. Die Saga wandert durch <code>PENDING</code>, <code>COMPENSATING</code>, <code>FAILED</code>, und der schon gebuchte Flug wird sichtbar storniert.
