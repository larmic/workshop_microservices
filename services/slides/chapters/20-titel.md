<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Saga</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Flug gebucht, Hotel sagt nein. Der Flug muss <span class="accent">wieder weg</span>.</p>
<p class="sub">Keine Transaktion &uuml;ber Service-Grenzen. Stattdessen eine Kette lokaler Schritte, jeder mit seinem Gegenst&uuml;ck.</p>
<figure class="foot-quote">
<blockquote>Entweder laufen alle Teiltransaktionen durch, oder die schon ausgef&uuml;hrten werden kompensiert.</blockquote>
<figcaption>Hector Garcia-Molina, Kenneth Salem, <em>Sagas</em>, Paper auf der ACM SIGMOD 1987 (sinngem&auml;&szlig;)</figcaption>
</figure>
</div>
<div class="split-figure">
<img src="./assets/saga.svg" alt="Booking ruft Flight, Hotel und Car nacheinander auf, Hotel scheitert, ein Pfeil f&uuml;hrt von Hotel zur&uuml;ck zu Flight"/>
</div>
</div>

</div>

Note:
- <strong>Grafik:</strong> Hotel scheitert, Car wird gar nicht mehr gefragt. Der lime Pfeil storniert den Flug.
- <strong>Herkunft:</strong> Kein Buch, sondern das Paper, das den Begriff Saga gepr&auml;gt hat. Zwei Datenbank-Forscher aus Princeton, 1987.
  - Problem damals: lange Transaktionen halten Sperren zu lange. L&ouml;sung: in Teilschritte zerlegen, jeder committet f&uuml;r sich und hat eine Kompensation.
  - In die Microservice-Welt gebracht hat es Chris Richardson (<em>Microservices Patterns</em>, 2018). Das Muster ist also &auml;lter als Microservices.
- <strong>Spicy:</strong> Kompensation ist kein Rollback. Der Flug war kurz wirklich gebucht, und jemand h&auml;tte ihn sehen k&ouml;nnen.
