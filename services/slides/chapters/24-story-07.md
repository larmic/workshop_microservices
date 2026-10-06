<div class="page">

<p class="kicker">Story 7 &middot; &Uuml;bung</p>

<div class="page-head">

## Die Saga wird leise

<span class="badge">&asymp; 25 min</span>
</div>

<p class="subtitle">Ihr seid die Services. Der Broker ist ein Tisch mit Zetteln.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Aufbau</h4>
<ol>
<li>Vier Leute sind Booking, Flight, Hotel und Car. Eine Person ist der Kunde.</li>
<li>Ein Tisch in der Mitte ist der Broker. Events sind Zettel, die dort abgelegt werden.</li>
<li>Jedes Backend hat einen Stapel Buchungskarten, die es ausstellen und wieder einziehen kann.</li>
<li>Der Kunde bucht eine Reise: Flug, Hotel, Mietwagen. Hotel lehnt jede zweite Buchung ab.</li>
<li>Alle anderen beobachten und notieren, wer wann was wusste.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li><strong>Runde 1, Orchestration:</strong> Booking ruft die Backends nacheinander an. Scheitert Hotel, ruft Booking Flight an und l&auml;sst stornieren.</li>
<li><strong>Runde 2, Choreography:</strong> Booking legt nur einen Zettel <code>CompensationRequested</code> auf den Tisch und antwortet dem Kunden. Flight holt ihn sich und zieht die Karte ein.</li>
<li>Vergleichen: Wer hat gewartet? Wer wusste am Ende, ob storniert wurde?</li>
<li>Weiter nach unten: drei St&ouml;rungen.</li>
</ol>
</div>
</div>

</div>

</div>

Note:
- Hook: &bdquo;Story 6 habt ihr programmiert. Heute spielen wir dieselbe Saga zweimal durch, einmal mit Dirigent, einmal ohne. Achtet darauf, wer wartet und wer am Ende Bescheid wei&szlig;.&ldquo;
- Material: ein Tisch oder ein Flipchart als Broker, Zettel in zwei Farben (Events und Buchungskarten), f&uuml;nf Rollenkarten (Kunde, Booking, Flight, Hotel, Car), eine Karteikarte pro Beobachter.
- Runde 1 nachspielen wie Story 6: Booking spricht Flight an (Karte), Hotel an (lehnt ab), spricht Flight erneut an (Storno), antwortet dem Kunden. Booking hat die ganze Zeit gewartet und wei&szlig; am Ende alles.
- Runde 2: Booking spricht Flight an (Karte), Hotel an (lehnt ab), legt einen Zettel auf den Tisch, antwortet dem Kunden sofort. Flight schaut irgendwann auf den Tisch, holt den Zettel, zieht die Karte ein. Booking hat nicht gewartet und wei&szlig; am Ende nicht, ob storniert wurde.
- Erwartung: Runde 2 f&uuml;hlt sich schneller und lockerer an. Die Beobachter merken, dass niemand mehr den &Uuml;berblick hat. Das ist der Trade-off von der Karten-Folie, jetzt am eigenen Leib.
- Wiedererkennung: Die &Uuml;bung steht im Dashboard unter Story 7, &bdquo;Story lesen&ldquo;. Das Saga-Panel darunter zeigt die Referenz-Implementierung mit Webhook-Events, wenn ihr das Echte sehen wollt.
- Weiter nach unten: drei St&ouml;rungen, die den Unterschied zwischen Webhook und Broker zeigen.
- Vollst&auml;ndige Beschreibung: <code>docs/stories/story-07-choreography-saga.md</code>.
