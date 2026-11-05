<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Story 4 &middot; und jetzt schreiben</p>

<div class="page-body">

<p class="statement">Flight-Breaker ist OPEN.<br>Die Kundin klickt auf <span class="accent">Buchen</span>.</p>

<p class="source">Hotel und Mietwagen buchen wir trotzdem? Ist das gut? Was w&auml;ren die Alternativen?</p>

</div>

</div>

Note:
- Nur Fragen stellen, keine L&ouml;sung zeigen. Lesen und Schreiben sind zwei verschiedene F&auml;lle.
- <strong>Nein, sofort:</strong> klare Meldung statt h&auml;ngender Seite. Fail fast ist auch ein Fallback.
- <strong>Asynchron:</strong> Hotel und Auto buchen, Flug vormerken und nachbuchen, sobald der Breaker schlie&szlig;t. Braucht Zustand und Benachrichtigung.
- <strong>Kompensation:</strong> Hotel und Auto buchen, Flug scheitert, beides stornieren. Das ist die Saga, Story 6.
- <strong>Alternativer Anbieter:</strong> zweiter Flug-Provider als Fallback. Teuer, selten.
