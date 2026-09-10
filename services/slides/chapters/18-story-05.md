<div class="page">

<p class="kicker">Story 5 &middot; &Uuml;bung</p>

<div class="page-head">

## Der geteilte Pool

<span class="badge">&asymp; 25 min</span>
</div>

<p class="subtitle">Ihr seid die Requests. Der Becher ist der Pool.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Aufbau</h4>
<ol>
<li>Zehn Chips in einem Becher, der Pool des Booking-Service.</li>
<li>Drei Leute sind Flight, Hotel und Car.</li>
<li>Alle anderen sind Requests und stellen sich an.</li>
<li>Wer bedient werden will, braucht einen Chip. Zur&uuml;ck gibt es ihn erst, wenn das Backend fertig ist.</li>
<li>Flight und Car antworten sofort. Hotel braucht zwanzig Sekunden.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li>Jedes Team schreibt vorher auf: Was passiert, wenn ein Drittel der Requests zu Hotel will?</li>
<li><strong>Runde 1:</strong> ein Becher f&uuml;r alle drei Backends.</li>
<li><strong>Runde 2:</strong> drei Becher mit vier, drei und drei Chips.</li>
<li>Karten aufdecken, vergleichen.</li>
</ol>
</div>
</div>

</div>

</div>

Note:
- Hook: &bdquo;Story 4 hat uns gegen kaputte Backends geh&auml;rtet. Heute die nervigere Variante: Das Backend antwortet, nur sehr, sehr langsam. Kein Fehler, und trotzdem rei&szlig;t es alles mit.&ldquo; Heute wird nicht programmiert, sondern gespielt.
- Material: ein gro&szlig;er Becher, drei kleine Becher, zehn Chips (M&uuml;nzen, Pokerchips, Zettel), Karteikarten f&uuml;r die Vorhersagen, eine Uhr mit Sekunden.
- Rollen: Flight und Car geben den Chip sofort zur&uuml;ck. Hotel z&auml;hlt laut bis zwanzig und h&auml;lt den Chip solange. Requests, die keinen Chip bekommen, gehen mit &bdquo;503&ldquo; zur&uuml;ck ans Ende der Schlange.
- Erwartung Runde 1: Nach wenigen Sekunden liegen alle zehn Chips bei Hotel. Flight- und Car-Requests scheitern, obwohl beide sofort antworten k&ouml;nnten. Das ist der Kern des Patterns, sichtbar in einer Minute.
- Erwartung Runde 2: Hotel-Requests scheitern weiter (drei Chips reichen nicht), Flight und Car laufen ungest&ouml;rt. Genau das haben die Karten vorher meist nicht vorhergesagt.
- Wiedererkennung: Die &Uuml;bung steht im Dashboard unter Story 5, &bdquo;Story lesen&ldquo;. Das Bulkhead-Panel darunter zeigt die Referenz-Implementierung, wenn ihr das Echte sehen wollt: Hotel auf langsam stellen, Burst dr&uuml;cken.
- Weiter nach unten: Warum ausgerechnet zehn?
- Vollst&auml;ndige Beschreibung: <code>docs/stories/story-05-bulkhead.md</code>.
