<div class="page">

<p class="kicker">Story 7 &middot; W&uuml;rfelspiel</p>

<div class="page-head">

## Die Saga wird leise

<span class="badge design">&asymp; 60 min <small>davon 25 Spiel</small></span>
</div>

<p class="subtitle">Booking und Kunde wollen <span class="hl">4 Reisen</span> buchen und haben daf&uuml;r 6 Versuche.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Vor dem Spiel</h4>
<ol>
<li><strong>Booking und Kunde</strong> spielen zusammen. <strong>Flight, Hotel, Car</strong> w&uuml;rfeln nur und behalten, was bei ihnen liegen bleibt.</li>
<li>Der Kunde hat 13 Chips. Jede gebuchte Leistung kostet ihn einen Chip.</li>
<li>Vorher gibt der Kunde Booking ein paar Chips als <strong>Retry-Budget</strong>. Die kommen nicht zur&uuml;ck.</li>
<li>Hat der Kunde vor einer Reise weniger als 3 Chips, ist Schluss.</li>
</ol>
</div>
<div class="steps-col">
<h4>Jede Reise</h4>
<ol>
<li>Booking w&uuml;rfelt: Bei 4&ndash;6 buchen alle. Bei 1&ndash;3 sagt ein Service nein (1 Flight, 2 Hotel, 3 Car), wer davor dran war, hat schon gebucht.</li>
<li>Wer nein sagt, w&uuml;rfelt den Grund: 1&ndash;3 nicht erreichbar (<code>503</code>), Booking zahlt einen Budget-Chip und der Service bucht doch. 4&ndash;6 ausgebucht, die Reise scheitert.</li>
<li>Scheitert sie, ruft Booking den schon gebuchten Services &bdquo;Storno!&ldquo; zu und sagt sofort &bdquo;Storniert&ldquo;.</li>
<li>Erst dann w&uuml;rfelt jeder Gerufene, ob er es geh&ouml;rt hat: 1&ndash;4 Chip zur&uuml;ck an den Kunden, 5&ndash;6 nicht geh&ouml;rt, der Chip bleibt liegen.</li>
</ol>
<p class="step-note"><strong><code>503</code> &rarr; Retry:</strong> vorw&auml;rts. <strong>Ausgebucht &rarr; Storno:</strong> r&uuml;ckw&auml;rts.</p>
</div>
</div>

</div>

</div>

Note:
- <strong>Hook:</strong> &bdquo;Story 6 habt ihr gebaut. Jetzt w&uuml;rfeln wir sie, und Booking ruft das Storno nur noch zu.&ldquo;
- <strong>Tische zu f&uuml;nft</strong>, Rollen bleiben das ganze Spiel. Ein sechster f&uuml;hrt den Zettel. Pro Tisch ein W&uuml;rfel, 13 Chips. Chips laut abz&auml;hlen lassen, die Balance h&auml;ngt daran.
- <strong>Erste Entscheidung:</strong> Wie viel Retry-Budget? Jeder Chip fehlt f&uuml;rs Buchen. Simuliert: 0 Chips 32&nbsp;%, 1 Chip 54&nbsp;%, 2 Chips 44&nbsp;%, 3 Chips 17&nbsp;% gewonnen. Resilienz kostet Kapazit&auml;t.
- <strong>Vorw&auml;rts oder r&uuml;ckw&auml;rts:</strong> Scheitert ein Schritt, kann Booking ihn wiederholen (Forward Recovery) oder das Gebuchte kompensieren (Backward Recovery).
  - Vorw&auml;rts nur bei technischen Fehlern wie <code>503</code>, ausgebucht bleibt ausgebucht.
  - Unsere Referenz kennt nur r&uuml;ckw&auml;rts, der Circuit Breaker wiederholt nicht.
- <strong>Leise:</strong> &bdquo;Storniert&ldquo; kommt vor dem Wurf der Gerufenen. 5&ndash;6 hei&szlig;t: Chip bleibt liegen, niemand sagt etwas. Genau die Referenz: <code>dispatch-failed</code> im Log, Step trotzdem <code>COMPENSATED</code>.
- <strong>Kassensturz:</strong> Jeder Service sollte so viele Chips haben, wie Reisen durch sind. Jeder Chip dar&uuml;ber ist ein stilles Storno. Ergebnisse aller Tische an die Tafel, bei vier Tischen taucht fast immer eines auf.
- &Uuml;berleitung: Chips liegen lassen, die Leiter sortiert sie. Die Demo steht in der Note dort.
