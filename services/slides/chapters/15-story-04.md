<div class="page">

<p class="kicker">Story 4 &middot; Rollenspiel</p>

<div class="page-head">

## Wenn der Flug ausf&auml;llt

<span class="badge design">&asymp; 60 min <small>davon 20 Spiel</small></span>
</div>

<p class="subtitle">Ihr seid die Suchanfragen. Flight ist krank. Booking hat den <span class="hl">Schalter</span>.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> hat Strichliste, Schild CLOSED / OPEN / HALF_OPEN und Karten &bdquo;Hotel &#10003; &middot; Car &#10003; &middot; flights: []&ldquo;. Warten hei&szlig;t laut z&auml;hlen: eins, zwei, drei, Timeout.</li>
<li><strong>Flight</strong> folgt den geheimen Ansagen des Trainers. <strong>Hotel + Car</strong> antworten sofort.</li>
<li><strong>Alle anderen</strong> sind Suchanfragen. Wer eine Antwort hat, legt die Karte in den Korb und stellt sich <strong>sofort hinten wieder an</strong>.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li><strong>Runde 1, ohne Schalter:</strong> Flight schweigt, nach drei Sekunden sagt Booking &bdquo;500&ldquo;. <span class="fragment">&rarr; null Karten</span></li>
<li><strong>Runde 2, mit Schalter:</strong> Jeder Timeout gibt einen Strich <em>und</em> die Karte. F&uuml;nf Striche <strong>in Folge</strong>: OPEN, ab dann Karte sofort. <span class="fragment">&rarr; Karte schon beim ersten Fehler, nur langsam</span></li>
<li><strong>HALF_OPEN:</strong> 30 s nach OPEN bekommt ein Gast das PROBE-Token und fragt, alle dahinter bekommen sofort die Karte. Antwortet Flight: Liste wegwischen, CLOSED.</li>
<li><strong>Runde 3a:</strong> Flight antwortet jedes zweite Mal. Wann geht das Schild auf OPEN? <span class="fragment">&rarr; nie, aber jede zweite Suche wartet</span></li>
<li><strong>Runde 3b:</strong> F&uuml;nf rote Anfragen, Flight sagt sofort &bdquo;400&ldquo;. Gibt das einen Strich? <span class="fragment">&rarr; ja, bei uns schon</span></li>
</ol>
</div>
</div>

</div>

</div>

Note:
- <strong>Hook:</strong> &bdquo;Story 3 hat uns Services finden lassen. Jetzt haben wir einen gefunden, und er antwortet nicht.&ldquo;
- <strong>Setup:</strong> Ab 7 Personen. Runden einzeln ansagen, Fragment erst nach der Vorhersage aufdecken.
- <strong>Flight geheim briefen:</strong>
  - R1, R2 und erste Probe: schweigen. Zweite Probe: antworten.
  - R3a abwechselnd. R3b gesund, auf rote Anfragen sofort &bdquo;400&ldquo;.
  - Die Uhr bist du: 30 s nach &bdquo;OPEN&ldquo; rufst du &bdquo;HALF_OPEN&ldquo;.
- <strong>Auswerten:</strong>
  - R1 gegen R2 (Recap 02): Die Karte gibt es ab dem ersten Fehler, der Schalter spart nur das Warten.
  - HALF_OPEN (Recap 03): Einfrieren, &bdquo;Wer darf zu Flight?&ldquo; Nur einer, die anderen bekommen sofort die Karte.
  - R3a: in Folge, ein Erfolg setzt auf null.
  - R3b (Recap 01): Unser Code z&auml;hlt alles ab 400. Transfer Recap 04: &bdquo;Booking-2 hat die roten G&auml;ste nie gesehen. Welches Schild zeigt es?&ldquo;
- <strong>Vereinfacht:</strong> Neu-Anstellen ist ein Client-Retry, den die Referenz nicht hat.
- <strong>Demo:</strong> Flight auf &bdquo;H&auml;ngt&ldquo;. F&uuml;nf Calls warten je 3 s, ab dem 6. sofort 200 mit <code>flights: []</code> und <code>X-Circuit-Open</code>. Schneller: &bdquo;Fehler&ldquo;. Details: <code>docs/stories/story-04-circuit-breaker.md</code>.
