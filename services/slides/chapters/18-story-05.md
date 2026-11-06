<div class="page">

<p class="kicker">Story 5 &middot; Rollenspiel</p>

<div class="page-head">

## Der geteilte Pool

<span class="badge design">&asymp; 25 min</span>
</div>

<p class="subtitle">Gleiche Schlange wie Story 4. Neu: Takt, Chips und ein <span class="hl">langsames</span> Hotel.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Takt</strong> z&auml;hlt laut 1 bis 8 und klatscht. Pro Klatscher gibt die vorderste Suchanfrage eine Karte ab und stellt sich <strong>sofort hinten wieder an</strong>.</li>
<li><strong>Booking</strong> ruft Flight, Hotel und Car, jeder Aufruf braucht einen Chip. Kein Chip: Karte wartet.</li>
<li><strong>Hotel</strong> legt Karte und Chip aufs Uhrfeld, nach acht Takten zur&uuml;ck. <strong>Flight + Car</strong> sind sofort fertig.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li><strong>Probelauf:</strong> Hotel ist schnell, drei Takte zum Einspielen.</li>
<li><strong>Runde 1, ein Becher, sechs Chips:</strong> Wie viele Karten warten nach 20 Takten? <span class="fragment">&rarr; vier, ab Suche 7 wartet jede schon auf Flight</span></li>
<li><strong>Runde 2, drei Becher, je zwei Chips:</strong> Hotel-Becher leer, sofort <code>hotels: []</code>. Wie viele bekommen Hotels? <span class="fragment">&rarr; sechs, die anderen 14 sofort</span></li>
<li><strong>Schild:</strong> Noch CLOSED? <span class="fragment">&rarr; ja, kein Timeout, Rejects erreichen den Breaker nie</span></li>
</ol>
</div>
</div>

</div>

</div>

Note:
- Hook: &bdquo;Story 4 hatte ein kaputtes Backend. Heute antwortet Hotel, nur sehr langsam.&ldquo; Ab 7 Personen. Erst Handzeichen, dann Fragment. Material, Detailregeln (Reihenfolge im Takt, wer den freien Chip bekommt) und Spickzettel im Doc und im Dashboard. Nach jeder Runde Stapel und Uhr zur&uuml;cksetzen.
- Gez&auml;hlt wird in Takten, nicht in Sekunden. Ein freier Chip kann im selben Takt mehrere Aufrufe bedienen.
- R1 (Recap 01): Der gemeinsame Becher ist der Thread-Pool klassischer Server, die Go-Referenz hat keinen. Flight ist gesund und wartet trotzdem. Ehrlich: R1 hat mehr komplette Antworten (10 gegen 4), weil Hotel alle sechs Chips nehmen durfte.
- R2 ist die Referenz: Pool je Backend, keine Queue, Client bekommt 200 mit <code>hotels: []</code> und <code>X-Bulkhead-Full</code>.
- Br&uuml;cke zu 18a: Im Probelauf lag nie mehr als ein Chip bei Hotel, langsam sind es acht. Der leere Becher ist das Signal, kein Grund zum Vergr&ouml;&szlig;ern.
- Demo: Alle Services auf Normal, dann Hotel auf &bdquo;Langsam&ldquo;, Burst. &bdquo;hotel rej&ldquo; zeigt etwa zehn, &bdquo;failed&ldquo; bleibt 0. Danach Hotel zur&uuml;ck. Details: <code>docs/stories/story-05-bulkhead.md</code>.
