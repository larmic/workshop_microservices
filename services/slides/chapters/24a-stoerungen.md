<div class="page">

<p class="kicker">Story 7 &middot; Nachbereitung</p>

<div class="page-head">

## Die Leiter

</div>

<p class="subtitle">Booking storniert nicht mehr selbst, antwortet aber weiter dem Kunden. Was fehlt, damit kein Chip liegen bleibt <span class="hl">und</span> Booking es wei&szlig;?</p>

<div class="page-body">

<div class="cards cards-4 violet">
<div class="card fragment" data-fragment-index="1">
<h3>Zuruf</h3>
<p>Booking ruft und meldet sofort &bdquo;storniert&ldquo;. Wer es nicht h&ouml;rt, beh&auml;lt den Chip.</p>
<code>Webhook &middot; unsere Referenz</code>
</div>
<div class="card fragment" data-fragment-index="2">
<h3>Korb</h3>
<p>Der Korb h&auml;lt das Storno fest und stellt zu, bis quittiert wird. Keins fehlt, manches kommt doppelt.</p>
<code>Broker &middot; at-least-once</code>
</div>
<div class="card fragment" data-fragment-index="3">
<h3>Nummer</h3>
<p>Wer sich die Storno-Nummer merkt und sie wiedererkennt, quittiert nur und zahlt nicht noch einmal.</p>
<code>eventId &middot; Idempotenz</code>
</div>
<div class="card fragment" data-fragment-index="4">
<h3>Antwort zur&uuml;ck</h3>
<p>Booking meldet &bdquo;storniert&ldquo; erst, wenn alle Chips zur&uuml;ck sind. Sonst meldet Booking nach einem Timeout <code>STUCK</code>.</p>
<code>Reply &middot; Timeout</code>
</div>
</div>

<div class="foot fragment" data-fragment-index="5">
<div class="callout">Eventing macht das Problem nicht kleiner. Es macht es leiser.</div>
</div>

</div>

</div>

Note:
- Nach dem Kassensturz die Frage stellen, dann Stufe f&uuml;r Stufe aufdecken. Gespielt wurde nur der Zuruf, er ist die Referenz.
- <strong>Zuruf</strong> (gespielt): ein Versuch, kein Retry. Booking meldet &bdquo;storniert&ldquo;, bevor jemand gew&uuml;rfelt hat (Recap 01).
- <strong>Korb</strong> (Gedankenexperiment): Persistenz und Redelivery. Kein Storno geht verloren. Kommt das &bdquo;erledigt&ldquo; nicht an, stellt der Korb noch einmal zu, und Flight zahlt doppelt (Recap 02).
- <strong>Nummer</strong> (Gedankenexperiment): <code>eventId</code> speichern, in derselben Transaktion wie das Storno. Unsere Backends loggen nur und deduplizieren nicht.
- <strong>Antwort zur&uuml;ck</strong> (Gedankenexperiment): <code>BookingCancelled</code> zur&uuml;ck an Booking, Timeout, Saga-Status <code>STUCK</code>. Ehrlich, aber Booking wartet wieder, nur asynchron.
- Merksatz: In Story 6 hat Booking den Schmerz gesp&uuml;rt. Jetzt liegt er als Chip beim Service. Wer Choreografie ernst meint, f&auml;ngt bei der Haltbarkeit des Events an. Danach Demo: Car auf &bdquo;Fehler&ldquo;, Saga starten. <code>FAILED</code>, Hotel und Flight <code>COMPENSATED</code>. Nach f&uuml;nf Klicks ist der Car-Breaker offen.
