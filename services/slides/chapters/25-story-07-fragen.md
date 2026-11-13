<div class="page">

<p class="kicker">Story 7 &middot; Recap</p>

<div class="page-head">

## Vier Fragen an euch

<span class="badge">Bonus zu 02 &darr;</span>
</div>

<div class="page-body">

<div class="numlist recap compact">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Ein Chip bleibt liegen</p>
<h3>Booking schickt das Storno per HTTP-POST, Flight ist seit der Buchung down. Wo ist das Event <span class="hl">jetzt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Event ist weg</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Das Storno kommt zweimal</p>
<h3>Ein Broker h&auml;lt das Storno fest und liefert, bis Flight best&auml;tigt. Stirbt Flight dazwischen, kommt das Event <span class="hl">noch einmal</span>.</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">eventId merken, Doppeltes ignorieren</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Booking meldet zu fr&uuml;h &bdquo;storniert&ldquo;</p>
<h3>Booking soll erst &bdquo;storniert&ldquo; sagen, wenn alle Chips zur&uuml;ck sind. Was muss zur&uuml;ckkommen, und was gilt, wenn es <span class="hl">nie</span> kommt?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Antwort zur&uuml;ck, sonst STUCK</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Brauchen wir &uuml;berhaupt einen Dirigenten?</p>
<h3>Volle Choreografie hie&szlig;e: Hotel reagiert auf Flight, Car auf Hotel. Wer merkt es, wenn Hotel das <span class="hl">nicht wei&szlig;</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Prozess braucht eine Adresse</code>
</div>
</div>

</div>

</div>

Note:
- Erst diskutieren, dann per Klick die Merks&auml;tze. 01 bis 03 steigen auf: vom Zuruf bis zur Antwort zur&uuml;ck.
- <strong>01 Chip bleibt liegen:</strong> Das haben wir gespielt, so arbeitet die Referenz.
  - Ein Versuch, kein Retry: <code>dispatch-failed</code> im Log, trotzdem <code>COMPENSATED</code>.
  - <strong>Spicy:</strong> Eventing macht das Problem nicht kleiner. Es macht es leiser.
- <strong>02 Storno kommt zweimal:</strong> Gedankenexperiment an den liegen gebliebenen Chips.
  - Ein Broker liefert, bis quittiert ist (at-least-once): Nichts fehlt, manches kommt doppelt.
  - Abhilfe: <code>eventId</code> merken, in derselben Transaktion wie die Fachlogik. Unsere Backends tun das nicht.
  - Einwand aus Story 6: &bdquo;Das Storno war doch idempotent, zweimal 204.&ldquo; Stimmt, solange Stornieren nur &bdquo;weg&ldquo; hei&szlig;t. H&auml;ngt mehr daran (Gutschrift, Mail, Kontingent plus eins), z&auml;hlt jedes Doppelte.
  - Was ein Broker sonst kann: Bonus-Folie.
- <strong>03 Antwort zur&uuml;ck:</strong> Auch ein Gedankenexperiment.
  - Reply-Event (<code>BookingCancelled</code>) plus Timeout, sonst <code>STUCK</code> statt <code>COMPENSATED</code>.
  - Booking wartet wieder, nur asynchron. Das f&uuml;hrt zu 04.
- <strong>04 Dirigent:</strong> Niemand merkt es, es passiert einfach nichts.
  - Ein fester Ablauf mit Kompensation geh&ouml;rt an eine Stelle. Choreografie nur ohne Reihenfolge (Benachrichtigung, Read-Model).
  - Unsere Referenz ist der Mittelweg: Booking dirigiert den Hinweg, storniert wird per Zuruf.
  - <strong>Spicy:</strong> Bei Choreografie kennt jeder jeden, nur steht es nirgends.
- <strong>Preis:</strong> Fehler werden still, und mit einem Broker kommen Duplikate dazu.
- Br&uuml;cke zu Story 8: Wer merkt es, wenn ein Event verloren geht?
- Lang: <code>docs/questions/story7.md</code>.
