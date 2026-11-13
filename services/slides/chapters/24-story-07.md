<div class="page">

<p class="kicker">Story 7 &middot; W&uuml;rfelspiel &bdquo;Retry oder Storno&ldquo;</p>

<div class="page-head">

## Die Saga wird leise

<span class="badge design">&asymp; 25 min <small>Spiel</small></span>
</div>

<p class="subtitle">Ziel: <span class="hl">4 Reisen</span> buchen. H&ouml;chstens 6 Reisen d&uuml;rfen starten, gescheiterte z&auml;hlen mit.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Vor dem Spiel</h4>
<ol>
<li><strong>Kunde und Booking</strong> spielen zusammen. <strong>Flight, Hotel, Car</strong> w&uuml;rfeln nur und behalten, was bei ihnen liegt.</li>
<li>Der Kunde hat 13 Chips. <strong>Eine Reise kostet 3</strong>, je einen f&uuml;r Flight, Hotel und Car.</li>
<li>Kunde und Booking machen aus, wie viele Chips der Kunde Booking als <strong>Retry-Budget</strong> gibt. Booking legt sie als eigenen Stapel zur Seite. Sie kommen nie zur&uuml;ck.</li>
<li>Kann der Kunde Booking keine 3 Chips mehr geben, ist das Spiel vorbei.</li>
</ol>
<p class="step-note"><strong><code>503</code> &rarr; Retry:</strong> vorw&auml;rts. <strong>Ausgebucht &rarr; Storno:</strong> r&uuml;ckw&auml;rts.</p>
<p class="step-note"><strong>Nach jeder Reise:</strong> Booking gibt dem Kunden zur&uuml;ck, was vom Reisegeld &uuml;brig ist.</p>
<p class="step-note"><strong>Am Ende Kassensturz:</strong> Jeder Service sollte so viele Chips haben, wie Reisen geklappt haben.</p>
</div>
<div class="steps-col">
<h4>Jede Reise</h4>
<ol>
<li>Der Kunde gibt Booking 3 Chips. Damit startet die Reise.</li>
<li>Booking fragt der Reihe nach <strong>Flight, Hotel, Car</strong> und w&uuml;rfelt einmal: Bei 4&ndash;6 buchen alle drei. Bei 1 sagt Flight nein, bei 2 Hotel, bei 3 Car. Nur wer vor dem Nein gebucht hat, bekommt einen Chip.</li>
<li>Wer nein sagt, w&uuml;rfelt den Grund: Bei 1&ndash;3 ist er nicht erreichbar (<code>503</code>). Booking zahlt ihm einen Chip aus dem Retry-Budget, er bucht doch. Wer nach ihm kommt, bucht ohne neuen Wurf, die Reise klappt.</li>
<li>Bei 4&ndash;6 ist er ausgebucht, oder das Budget ist leer: Die Reise scheitert. Booking ruft den schon Gebuchten r&uuml;ckw&auml;rts &bdquo;Storno!&ldquo; zu und sagt sofort &bdquo;Storniert&ldquo;.</li>
<li>Erst dann w&uuml;rfelt jeder Gerufene: Bei 1&ndash;4 hat er es geh&ouml;rt und gibt den Chip an den Kunden zur&uuml;ck. Bei 5&ndash;6 hat er es &uuml;berh&ouml;rt und beh&auml;lt ihn.</li>
</ol>
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
  - Bei <code>503</code> mit Budget immer retten, das ist rechnerisch richtig. Ausgebucht l&auml;sst sich nicht retten: Auch mit unbegrenztem Budget scheitert jede vierte Reise.
- <strong>Leise:</strong> &bdquo;Storniert&ldquo; kommt vor dem Wurf der Gerufenen. 5&ndash;6 hei&szlig;t: Chip bleibt liegen, niemand sagt etwas. Genau die Referenz: <code>dispatch-failed</code> im Log, Step trotzdem <code>COMPENSATED</code>.
- <strong>Kassensturz:</strong> Jeder Service sollte so viele Chips haben, wie Reisen geklappt haben. Jeder Chip dar&uuml;ber ist ein stilles Storno. Ergebnisse aller Tische an die Tafel, bei vier Tischen taucht fast immer eines auf.
- <strong>Demo:</strong> kurz, 5 Minuten, eine Saga. Car auf &bdquo;Fehler&ldquo;, Saga starten. Das Bild kennt jeder aus Story 6, neu sind zwei Dinge: der Hinweis &bdquo;Kompensation gemeldet, Ankunft ungepr&uuml;ft&ldquo; und im Log <code>phase=publishing</code> und <code>phase=dispatched</code>. Das stille Storno zeigt die Demo nicht, nur der Code (<code>compensate</code> in <code>booking/story7/handler/booking.go</code>).
- &Uuml;berleitung: Chips bis zum Recap liegen lassen. Leitfrage: Was fehlt, damit kein Chip liegen bleibt <em>und</em> Booking es wei&szlig;?
