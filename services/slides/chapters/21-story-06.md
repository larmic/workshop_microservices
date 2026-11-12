<div class="page">

<p class="kicker">Story 6 &middot; Hands-on</p>

<div class="page-head">

## Alles oder nichts

<span class="badge">&asymp; 60 min</span>
</div>

<p class="subtitle">Drei Buchungen, am Ende ein Ergebnis: alles da oder nichts mehr <span class="hl">gebucht</span>.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Kunde</dd>
<dt>m&ouml;chte ich</dt>
<dd>eine Komplettbuchung aus Flug, Hotel und Mietwagen, die am Ende entweder vollst&auml;ndig steht oder nichts mehr gebucht ist,</dd>
<dt>damit</dt>
<dd>ich nicht mit einer halben Reise dastehe.</dd>
</dl>
<div class="story-list">
<h4>Die Saga</h4>
<ul>
<li>Eine Anfrage bucht Flug, Hotel und Mietwagen</li>
<li>Booking ruft die Services nacheinander auf (Orchestration)</li>
<li>Scheitert ein Schritt, werden alle vorherigen kompensiert, in umgekehrter Reihenfolge</li>
</ul>
</div>
<div class="story-list">
<h4>Die Kompensation</h4>
<ul>
<li>Jeder Service bietet <code style="white-space:nowrap">DELETE /bookings/{id}</code> an</li>
<li>Der Endpoint ist idempotent: zweiter Aufruf, gleicher Zustand, wieder <code>204</code></li>
</ul>
<h4>Sichtbarkeit</h4>
<ul>
<li>Saga-Status abfragbar: <code>PENDING</code>, <code>COMPLETED</code>, <code>COMPENSATING</code>, <code>FAILED</code></li>
</ul>
</div>
</div>

</div>

</div>

Note:
- <strong>Hook:</strong> &bdquo;Flight gebucht, Hotel sagt nein. Jetzt baut ihr den R&uuml;ckweg.&ldquo;
- <strong>Wiedererkennung:</strong> dieselbe Story im Dashboard unter Story 6, &bdquo;Aufgabe&ldquo;. Sprache und Framework wieder frei.
- <strong>R&uuml;ckblick Story 2:</strong> Der Storno vom Flipchart ist jetzt die Kompensation.
  - Sobald ein Retry dazukommt oder jemand nachst&ouml;&szlig;t, muss <code>DELETE</code> idempotent sein.
- <strong>Nur ein Versuch</strong> f&uuml;r die Kompensation, kein Retry, kein persistenter Status. Das macht das Pattern sichtbar, Persistenz und Retry kommen im Recap.
- <strong>Demo:</strong> Hotel auf &bdquo;Fehler&ldquo;, dann <code>POST /booking/bookings</code>.
  - Saga-Karte: Flight COMPENSATED, Hotel FAILED, Car nicht aufgerufen, Saga FAILED. <code>COMPENSATING</code> sieht man nicht, die Kompensation dauert Millisekunden.
  - Nicht zu oft klicken: Nach 5 Fehlern &ouml;ffnet der Breaker f&uuml;r Hotel, dann startet die Saga gar nicht erst.
  - Hotel zur&uuml;ck auf normal, neue Buchung, alles gr&uuml;n.
- <strong>Referenz:</strong> <code>services/booking/story6/handler/booking.go</code>. Volle Aufgabe: <code>docs/stories/story-06-saga.md</code>.
