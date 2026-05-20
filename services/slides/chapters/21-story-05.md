## Story 5

<p class="subtitle">Alles oder nichts &mdash; aber richtig <span class="time-badge">&asymp; 60 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

Eine Reisebuchung umfasst Flug, Hotel und Mietwagen. Wenn die Hotelbuchung fehlschl&auml;gt, <em>nachdem</em> der Flug bereits gebucht wurde, muss der Flug storniert werden. Klassische Datenbank-Transaktionen funktionieren nicht &uuml;ber Service-Grenzen hinweg.

#### User Story

Als <em>Kunde</em> m&ouml;chte ich <em>eine Komplettbuchung (Flug + Hotel + Mietwagen) durchf&uuml;hren, die entweder vollst&auml;ndig erfolgreich ist oder komplett zur&uuml;ckgerollt wird</em>, damit <em>ich nicht mit einer unvollst&auml;ndigen Buchung dastehe</em>.

</div>

</div>
<div>

<div class="story-card">

#### Akzeptanzkriterien

- Buchungsanfrage f&uuml;r Flug, Hotel und Mietwagen
- Services werden nacheinander aufgerufen (<em>Orchestration-Saga</em>)
- Bei Fehler in einem Schritt werden alle vorherigen Schritte <strong>kompensiert</strong> (Rollback)
- Jeder Service bietet einen <code>DELETE /bookings/{id}</code> als Kompensation an
- Saga-Status abfragbar (<code>PENDING</code>, <code>COMPLETED</code>, <code>COMPENSATING</code>, <code>FAILED</code>)
- Kompensations-Endpoints sind <strong>idempotent</strong>

</div>

</div>
</div>

Note:
- Hook: &bdquo;Resilience-Patterns aus Story 3 und 4 helfen <em>einem</em> Aufruf. Aber sobald ich mehrere zusammenh&auml;ngende Schritte habe (Flug + Hotel + Auto) und einer kippt, brauche ich etwas anderes.&ldquo; Klassisches Beispiel: Flight gebucht, Hotel sagt nein. Was tun mit dem Flug?
- Wiedererkennung: dieselbe Karte (Kontext / User Story / Akzeptanzkriterien) im Dashboard unter Story 5 &rarr; &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story5/</code> (Go, sequenzielle Saga in ca. 100 Zeilen).
- Im Workshop bewusst <strong>nur ein Versuch</strong> f&uuml;r die Kompensation, kein Retry, kein persistenter Status &mdash; das macht das Pattern sichtbar, ohne den 60-Min-Slot zu sprengen. Persistenz + Retry sind explizite Diskussionspunkte im Recap (Fragen 1, 2, 7).
- Demo-Drehbuch: Dashboard &rarr; Hotel auf &bdquo;Fehler&ldquo;, dann <code>POST /booking/bookings</code> &mdash; in der Saga-Karte ist sichtbar: Flight BOOKED &rarr; Hotel FAILED &rarr; status COMPENSATING &rarr; Flight COMPENSATED &rarr; status FAILED. Anschlie&szlig;end Hotel zur&uuml;ck auf normal, neue Buchung &mdash; alles gr&uuml;n.
- Time-Box 60 min inkl. Demo. Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-05-saga.md</code>.
