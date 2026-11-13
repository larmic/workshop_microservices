<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Choreography-Saga</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Booking dirigiert den Hinweg. Storniert wird per <span class="accent">Zuruf</span>.</p>
<p class="sub">Die Backends stornieren selbst. Booking ruft zu und meldet &bdquo;storniert&ldquo;, ob es jemand geh&ouml;rt hat oder nicht.</p>
<p class="source">Gewinn: Booking wartet nicht mehr auf das Storno. Was es kostet, w&uuml;rfeln wir gleich.</p>
</div>
<div class="split-figure">
<img src="./assets/choreography.svg" alt="Broker in der Mitte, Booking wirft ein Event ein, drei Backends holen es sich ab, eines ist gerade nicht da"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;In Story 6 wartet Booking, bis storniert ist, und muss jeden Storno-Fehler selbst ausbaden. Jetzt ruft es das Storno nur noch zu. Fachlich geh&ouml;rt die Stornierung doch zu Hotel selbst, nicht zum Aggregator.&ldquo;
- Grafik: Pattern mit Broker, das Event des fehlenden Backends bleibt liegen. Unsere Referenz ruft direkt zu, genau diesen Unterschied spielen wir.
- Ehrlich gesagt ist unsere Referenz ein Mittelweg: Booking dirigiert den Hinweg und entscheidet, wer ein Storno bekommt. Nur das Stornieren selbst machen die Backends. H&auml;ngt Hotel, wartet Booking weiter bis zum Timeout (3 s), dieselbe Zeit wie in Story 6. Heute wird nicht programmiert, sondern gew&uuml;rfelt.
