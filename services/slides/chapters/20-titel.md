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
- <strong>Hook:</strong> &bdquo;Flug gebucht, Hotel sagt nein. Was nun?&ldquo;
- <strong>Grafik:</strong> Hotel scheitert, Car wird gar nicht mehr gefragt. Der lime Pfeil storniert den Flug.
- <strong>Herkunft:</strong> Garcia-Molina und Salem, 1987, f&uuml;r lange Datenbank-Transaktionen. Jeder Schritt committet f&uuml;r sich und hat eine Kompensation.
- <strong>Spicy:</strong> Kompensation ist kein Rollback. Der Flug war kurz wirklich gebucht, und jemand h&auml;tte ihn sehen k&ouml;nnen.
