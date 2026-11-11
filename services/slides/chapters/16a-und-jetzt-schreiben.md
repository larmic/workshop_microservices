<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Story 4 &middot; und jetzt schreiben</p>

<div class="page-body">

<p class="statement">Flight-Breaker ist OPEN.<br>Die Kundin klickt auf <span class="accent">Buchen</span>.</p>

<p class="source">Hotel und Mietwagen buchen wir trotzdem? Ist das gut? Was w&auml;ren die Alternativen?</p>

</div>

</div>

Note:
- Nur fragen, keine L&ouml;sung zeigen. Lesen und Schreiben sind zwei F&auml;lle.
  - Suche: <code>flights: []</code> ist ein guter Fallback.
  - Buchung: &bdquo;ohne Flug&ldquo; w&auml;re eine falsche Buchung.
- <strong>Unsere Referenz:</strong> <code>POST /booking/bookings</code> pr&uuml;ft vorher alle Breaker. Einer OPEN: sofort 503, nichts gebucht.
  - Scheitert Hotel nach dem Flug, bleibt der Flug gebucht, kein Rollback.
- <strong>Alternativen:</strong>
  - Nein, sofort: klare Meldung statt h&auml;ngender Seite. Fail fast ist auch ein Fallback.
  - Asynchron: Hotel und Auto buchen, Flug nachbuchen, sobald der Breaker schlie&szlig;t. Braucht Zustand und Benachrichtigung.
  - Kompensation: buchen, bei Fehler stornieren. Das ist die Saga, Story 6.
  - Anderer Anbieter: zweiter Flug-Provider. Teuer, selten.
