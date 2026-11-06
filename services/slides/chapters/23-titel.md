<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Choreography-Saga</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Booking dirigiert den Hinweg. Storniert wird per <span class="accent">Zuruf</span>.</p>
<p class="sub">Die Backends stornieren selbst. Booking wartet nur noch darauf, dass jemand &bdquo;geh&ouml;rt&ldquo; sagt.</p>
<p class="source">Gewinn: Booking wartet nicht mehr auf das Storno, nur noch auf das &bdquo;geh&ouml;rt&ldquo;. Was es kostet, spielen wir gleich.</p>
</div>
<div class="split-figure">
<img src="./assets/choreography.svg" alt="Broker in der Mitte, Booking wirft ein Event ein, drei Backends holen es sich ab, eines ist gerade nicht da"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;In Story 6 wartet Booking, bis storniert ist, und muss jeden Storno-Fehler selbst ausbaden. Jetzt wartet es nur noch, bis jemand &sbquo;geh&ouml;rt&lsquo; sagt. Fachlich geh&ouml;rt die Stornierung doch zu Hotel selbst, nicht zum Aggregator.&ldquo;
- Die Grafik zeigt das Pattern mit Broker: Booking wirft ein Event ein, die Backends holen es sich ab, das Event des fehlenden Backends bleibt liegen. In unserer Referenz bleibt nichts liegen, Booking ruft direkt zu. Genau diesen Unterschied spielen wir.
- Ehrlich gesagt ist unsere Referenz ein Mittelweg: Booking dirigiert den Hinweg und entscheidet, wer ein Storno bekommt. Nur das Stornieren selbst machen die Backends. Ist Hotel weg, wartet Booking weiter bis zum Timeout (3 s), dieselbe Zeit wie in Story 6. Heute wird nicht programmiert, sondern gespielt.
