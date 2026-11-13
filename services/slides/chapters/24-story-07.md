<div class="page">

<p class="kicker">Story 7 &middot; W&uuml;rfelspiel</p>

<div class="page-head">

## Die Saga wird leise

<span class="badge design">&asymp; 60 min <small>davon 25 Spiel</small></span>
</div>

<p class="subtitle">Ziel: <span class="hl">4 Reisen</span> buchen. H&ouml;chstens 6 Reisen d&uuml;rfen starten, gescheiterte z&auml;hlen mit.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Vor dem Spiel</h4>
<ol>
<li><strong>Booking und Kunde</strong> spielen als Team. <strong>Flight, Hotel, Car</strong> w&uuml;rfeln nur und behalten, was bei ihnen liegt.</li>
<li>Der Kunde hat 13 Chips. <strong>Eine Reise kostet 3</strong>: je einen f&uuml;r Flight, Hotel und Car.</li>
<li>Vor der ersten Reise entscheidet das Team: Wie viele Chips bekommt Booking als <strong>Retry-Budget</strong>? Sie bezahlen nur Retries und kommen nie zur&uuml;ck.</li>
<li>Hat der Kunde vor einer Reise nur 0, 1 oder 2 Chips, ist das Team raus, auch mit Budget. Mit 3 geht es weiter.</li>
</ol>
</div>
<div class="steps-col">
<h4>Jede Reise</h4>
<ol>
<li>Der Kunde legt 3 Chips vor Booking. Damit startet die Reise.</li>
<li>Booking w&uuml;rfelt: 4&ndash;6 alle buchen und nehmen je einen Chip. 1 Flight, 2 Hotel, 3 Car sagt nein, wer davor dran war, hat schon gebucht.</li>
<li>Wer nein sagt, w&uuml;rfelt den Grund: 1&ndash;3 nicht erreichbar (<code>503</code>). Mit Budget zahlt Booking einen Budget-Chip, die Reise klappt ohne neuen Wurf.</li>
<li>4&ndash;6 ausgebucht oder kein Budget: Die Reise scheitert. Booking ruft den schon Gebuchten r&uuml;ckw&auml;rts &bdquo;Storno!&ldquo; zu und sagt sofort &bdquo;Storniert&ldquo;.</li>
<li>Erst dann w&uuml;rfelt jeder Gerufene: 1&ndash;4 geh&ouml;rt, Chip zur&uuml;ck an den Kunden. 5&ndash;6 &uuml;berh&ouml;rt, er beh&auml;lt ihn. Was danach noch vor Booking liegt, nimmt der Kunde zur&uuml;ck.</li>
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
