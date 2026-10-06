<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Choreography-Saga</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Niemand dirigiert. Jeder h&ouml;rt zu und <span class="accent">handelt selbst</span>.</p>
<p class="sub">Booking wirft ein Event ein und ist fertig. Die Backends stornieren sich selbst.</p>
<p class="source">Gewinn: Entkopplung. Preis: Niemand hat mehr den &Uuml;berblick.</p>
</div>
<div class="split-figure">
<img src="./assets/choreography.svg" alt="Broker in der Mitte, Booking wirft ein Event ein, drei Backends holen es sich ab, eines ist gerade nicht da"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;In Story 6 trug Booking allein die Verantwortung f&uuml;r die Stornierung. Wenn Hotel kurz weg ist, blockt der Storno, und Booking auch. Aber fachlich geh&ouml;rt die Stornierung doch zu Hotel selbst, nicht zum Aggregator.&ldquo;
- Die Grafik: Booking links wirft ein Event in den Broker, die Backends rechts holen es sich ab. Das dritte Backend ist gerade nicht da, sein Event bleibt im Broker liegen. Genau diese Situation spielen wir gleich nach.
- Choreography ist kein neues Pattern, sondern eine andere Verantwortungsverteilung f&uuml;r dieselbe Saga. Heute wird nicht programmiert, sondern gespielt: Jeder ist ein Service.
