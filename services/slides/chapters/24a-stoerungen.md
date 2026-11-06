<div class="page">

<p class="kicker">Story 7 &middot; Nachbereitung</p>

<div class="page-head">

## Die Leiter

</div>

<p class="subtitle">Booking storniert nicht mehr selbst, antwortet aber weiter der Kund:in. Was fehlt, damit jeder Becher stimmt <span class="hl">und</span> Booking es wei&szlig;?</p>

<div class="page-body">

<div class="cards cards-4 violet">
<div class="card fragment" data-fragment-index="1">
<h3>Zuruf</h3>
<p>Booking wartet nur auf &bdquo;geh&ouml;rt&ldquo;. Ist Flight weg, ist das Storno weg.</p>
<code>Webhook &middot; unsere Referenz</code>
</div>
<div class="card fragment" data-fragment-index="2">
<h3>Korb</h3>
<p>Der Korb h&auml;lt den Zettel fest und liest vor, bis quittiert wird. Keiner fehlt, manches kommt doppelt.</p>
<code>Broker &middot; at-least-once</code>
</div>
<div class="card fragment" data-fragment-index="3">
<h3>Nummer</h3>
<p>Wer sich die Zettelnummer merkt und sie wiedererkennt, quittiert nur, ohne Gummib&auml;rchen.</p>
<code>eventId &middot; Idempotenz</code>
</div>
<div class="card fragment" data-fragment-index="4">
<h3>Antwortzettel</h3>
<p>Booking meldet &bdquo;storniert&ldquo; erst, wenn beide Antworten da sind. Sonst z&auml;hlt Booking bis zehn und meldet <code>STUCK</code>.</p>
<code>Reply &middot; Timeout</code>
</div>
</div>

<div class="foot fragment" data-fragment-index="5">
<div class="callout">Eventing macht das Problem nicht kleiner. Es macht es leiser.</div>
</div>

</div>

</div>

Note:
- Nach dem dritten Kassensturz die Frage stellen, dann Stufe f&uuml;r Stufe aufdecken. Nur der Zuruf ist Referenz.
- <strong>Zuruf</strong> (Runde 2): ein Versuch, kein Retry. Booking merkt den Strich und meldet trotzdem &bdquo;storniert&ldquo; (Recap 01).
- <strong>Korb</strong> (Runde 3): Persistenz und Redelivery. Der Korb liefert nach, auch doppelt, denn Flight hat storniert, aber nicht quittiert (Recap 02).
- <strong>Nummer</strong> (nicht gespielt): Im Spiel merkt sich Flight die Nummer nicht. Echt: <code>eventId</code> speichern, in derselben Transaktion wie das Storno. Unsere Backends loggen nur und deduplizieren nicht.
- <strong>Antwortzettel</strong> (nicht gespielt): <code>BookingCancelled</code> zur&uuml;ck an Booking, Timeout, Saga-Status <code>STUCK</code>. Ehrlich, aber Booking wartet wieder, nur asynchron.
- Merksatz: In Story 6 hat Booking den Schmerz gesp&uuml;rt. Jetzt steht er nur noch auf der Strichliste. Wer Choreografie ernst meint, f&auml;ngt bei der Haltbarkeit des Events an. Danach Demo: Car auf &bdquo;Fehler&ldquo;, Saga starten. <code>FAILED</code>, Hotel und Flight <code>COMPENSATED</code>. Nach f&uuml;nf Klicks ist der Car-Breaker offen.
