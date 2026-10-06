<div class="page">

<p class="kicker">Story 7 &middot; &Uuml;bung</p>

<div class="page-head">

## Drei St&ouml;rungen

<span class="badge">&asymp; 20 min</span>
</div>

<p class="subtitle">Runde 2 noch einmal, jedes Mal mit einem Haken.</p>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Flight ist nicht da</h3>
<p>Flight verl&auml;sst den Raum, bevor der Zettel auf dem Tisch liegt. Booking hat dem Kunden schon geantwortet.</p>
<code>Wer merkt es? Wann?</code>
</div>
<div class="card">
<h3>Der Zettel kommt zweimal</h3>
<p>Booking ist sich nicht sicher, ob der Zettel angekommen ist, und legt einen zweiten hin. Flight findet beide.</p>
<code>Wie oft wird storniert?</code>
</div>
<div class="card">
<h3>Niemand sagt Bescheid</h3>
<p>Flight hat storniert. Der Kunde fragt Booking: &bdquo;Ist mein Flug jetzt weg?&ldquo;</p>
<code>Was antwortet Booking?</code>
</div>
</div>

<div class="swap swap-foot">
<p class="ask fragment fade-out" data-fragment-index="1">Vor jeder St&ouml;rung festlegen: Was passiert, wenn der Tisch ein <strong>Blatt Papier</strong> ist? Was, wenn er ein <strong>Postfach mit Quittung</strong> ist?</p>
<div class="fragment" data-fragment-index="1">
<p class="ask">Das Postfach h&auml;lt den Zettel fest und liefert ihn nach, auch doppelt. Deshalb tr&auml;gt jeder Zettel eine <strong>Nummer</strong>: Wer sie schon kennt, legt ihn weg. Und wer eine Antwort braucht, wartet auf einen <strong>Antwortzettel</strong>.</p>
<div class="callout">Eventing macht das Problem nicht kleiner. Es macht es leiser.</div>
</div>
</div>

</div>

</div>

Note:
- Jede St&ouml;rung einmal durchspielen, vorher die Beobachter tippen lassen. Die Frage &bdquo;Blatt Papier oder Postfach mit Quittung&ldquo; ist der Unterschied zwischen unserem Webhook und einem echten Broker, ohne das Wort Broker zu benutzen.
- <strong>Flight ist nicht da:</strong> Beim Blatt Papier liegt der Zettel ewig, niemand merkt es, der Flug bleibt gebucht. Booking hat l&auml;ngst &bdquo;storniert&ldquo; gemeldet. Beim Postfach bleibt der Zettel sicher liegen, Flight holt ihn beim Zur&uuml;ckkommen ab. Das ist Persistenz und Redelivery. Beim Webhook-POST unserer Referenz ist der Zettel einfach weg (Connection refused), Recap-Frage 1.
- <strong>Der Zettel kommt zweimal:</strong> Flight storniert zweimal. Bei einer Stornierung ist das harmlos, bei einer R&uuml;ckerstattung nicht. Gegenmittel: Jeder Zettel tr&auml;gt eine Nummer (<code>eventId</code>), Flight merkt sich bearbeitete Nummern. At-least-once plus Idempotenz, Recap-Frage 2.
- <strong>Niemand sagt Bescheid:</strong> Booking kann nur sagen: &bdquo;Ich habe einen Zettel hingelegt.&ldquo; Ob Flight ihn verarbeitet hat, wei&szlig; Booking nicht. L&ouml;sung: Flight legt einen Antwortzettel <code>BookingCancelled</code> hin, Booking wartet darauf mit Timeout. Das ist das Reply-Pattern aus dem Bonus der Story, und der Grund, warum die Saga einen <code>STUCK</code>-Status braucht.
- Klick: die drei Antworten in Spielsprache (Postfach = Persistenz und Redelivery, Nummer = <code>eventId</code> als Dedup-Key, Antwortzettel = Reply-Event) und der Merksatz. In Story 6 hat Booking den Schmerz gesp&uuml;rt und konnte reagieren. In Story 7 sieht Booking nichts. Wer Choreography ernst meint, f&auml;ngt nicht beim Event-Versand an, sondern bei der Haltbarkeit des Events.
