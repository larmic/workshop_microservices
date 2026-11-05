<div class="page">

<p class="kicker">Story 4 &middot; Rollenspiel</p>

<div class="page-head">

## Wenn der Flug ausf&auml;llt

<span class="badge design">&asymp; 20 min</span>
</div>

<p class="subtitle">Ihr seid die Suchanfragen. Flight ist krank. Booking hat den <span class="hl">Schalter</span>.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> hat Strichliste, Schild CLOSED / OPEN / HALF_OPEN und Karten <code>flights: []</code>. Warten hei&szlig;t laut z&auml;hlen: eins, zwei, drei, Timeout.</li>
<li><strong>Flight</strong> folgt den geheimen Ansagen des Trainers. <strong>Hotel + Car</strong> antworten sofort.</li>
<li><strong>Alle anderen</strong> sind Suchanfragen. Wer eine Antwort hat, stellt sich <strong>sofort hinten wieder an</strong>.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li><strong>Runde 1, ohne Schalter:</strong> Flight schweigt, nach drei Sekunden sagt Booking &bdquo;500&ldquo;. <span class="fragment">&rarr; null Karten</span></li>
<li><strong>Runde 2, mit Schalter:</strong> Jeder Timeout gibt einen Strich <em>und</em> die Karte. F&uuml;nf Striche <strong>in Folge</strong>: OPEN, ab dann Karte sofort. <span class="fragment">&rarr; Karte schon beim ersten Fehler, nur langsam</span></li>
<li><strong>HALF_OPEN:</strong> Ein Gast ist die Probe, alle dahinter bekommen sofort die Karte. Antwortet Flight: Liste wegwischen, CLOSED.</li>
<li><strong>Runde 3a:</strong> Flight antwortet jedes zweite Mal. Wann geht das Schild auf OPEN? <span class="fragment">&rarr; nie, aber jede zweite Suche wartet</span></li>
<li><strong>Runde 3b:</strong> F&uuml;nf rote Anfragen, Flight sagt sofort &bdquo;400&ldquo;. Gibt das einen Strich? <span class="fragment">&rarr; ja, bei uns schon</span></li>
</ol>
</div>
</div>

</div>

</div>

Note:
- Hook: &bdquo;Story 3 hat uns geholfen, Services zu finden. Jetzt haben wir einen gefunden, und er antwortet nicht.&ldquo; Ab 7 Personen spielbar. Runden einzeln ansagen, Fragment erst nach der Vorhersage aufdecken.
- Flight geheim briefen: R1 und R2 schweigen, erste Probe schweigen, zweite antworten. R3a abwechselnd. R3b gesund, auf rote Anfragen sofort &bdquo;400&ldquo;. Du bist die Uhr: 30 s nach &bdquo;OPEN&ldquo; rufst du &bdquo;HALF_OPEN&ldquo;.
- R1 gegen R2 (Recap 02): Die Karte gibt es ab dem ersten Fehler, der Schalter spart nur das Warten. Hotel und Car werden in R2 immer gefragt. Vereinfachung: Booking z&auml;hlt bis drei (Story-4-Timeout), Neu-Anstellen ist ein Client-Retry, den die Referenz nicht hat.
- HALF_OPEN (Recap 03): Einmal einfrieren, &bdquo;Wer darf zu Flight?&ldquo; Nur einer. Die anderen warten nicht, sie bekommen sofort die Karte.
- R3a: Gez&auml;hlt wird in Folge, ein Erfolg setzt auf null. R3b (Recap 01) ist ein Gedankenspiel: Unser Code z&auml;hlt alles ab 400. Transfer Recap 04: &bdquo;Booking-2 hat die roten G&auml;ste nie gesehen. Welches Schild zeigt es?&ldquo;
- Demo: Flight auf &bdquo;Fehler&ldquo; ist sofort ein 500, das Warten sieht man dort nicht. Client sieht immer 200 mit <code>flights: []</code>, <code>X-Circuit-Open</code> erst ab dem 6. Call. Details: <code>docs/stories/story-04-circuit-breaker.md</code>.
