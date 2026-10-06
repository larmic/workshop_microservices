<div class="page">

<p class="kicker">Story 6</p>

<div class="page-head">

## Alles oder nichts

<span class="badge">&asymp; 60 min</span>
</div>

<p class="subtitle">Drei Buchungen, ein Ergebnis: alles da, oder <span class="hl">nichts</span>.</p>

<div class="story page-body">

<div class="story-grid">
<dl class="story-user">
<dt>Als</dt>
<dd>Kunde</dd>
<dt>m&ouml;chte ich</dt>
<dd>eine Komplettbuchung aus Flug, Hotel und Mietwagen, die entweder vollst&auml;ndig gelingt oder komplett zur&uuml;ckgerollt wird,</dd>
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
<li>Jeder Service bietet <code>DELETE /bookings/{id}</code> an</li>
<li>Der Endpoint ist idempotent: <code>204</code> auch beim zweiten Aufruf</li>
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
- Hook: &bdquo;Resilience-Patterns aus Story 4 und 5 helfen <em>einem</em> Aufruf. Sobald mehrere zusammenh&auml;ngende Schritte im Spiel sind und einer kippt, brauche ich etwas anderes.&ldquo; Flight gebucht, Hotel sagt nein. Was tun mit dem Flug?
- Wiedererkennung: dieselbe Story (Kontext, User Story, Akzeptanzkriterien) im Dashboard unter Story 6, &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story6/</code> (Go, sequenzielle Saga in etwa 100 Zeilen).
- R&uuml;ckblick auf Story 2: Der Storno vom Flipchart ist jetzt die Kompensation. &bdquo;Ihr habt damals &uuml;ber DELETE zweimal diskutiert. Genau deshalb muss <code>DELETE /bookings/{id}</code> idempotent sein: Der Orchestrator wiederholt bei Timeout.&ldquo;
- Im Workshop bewusst <strong>nur ein Versuch</strong> f&uuml;r die Kompensation, kein Retry, kein persistenter Status. Das macht das Pattern sichtbar, ohne den Slot zu sprengen. Persistenz und Retry sind Diskussionspunkte im Recap.
- Demo-Drehbuch: Dashboard, Hotel auf &bdquo;Fehler&ldquo;, dann <code>POST /booking/bookings</code>. In der Saga-Karte sichtbar: Flight BOOKED, Hotel FAILED, Status COMPENSATING, Flight COMPENSATED, Status FAILED. Dann Hotel zur&uuml;ck auf normal, neue Buchung, alles gr&uuml;n.
- Time-Box 60 min inklusive Demo. Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-06-saga.md</code>.
