## Story 2 &mdash; Recap

<p class="subtitle">Was bleibt von der Design-Session</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Storno: DELETE oder Ressource?</h3>
<p><code>DELETE /bookings/{id}</code> ist idempotent und klar. Wo bleibt der <span class="hl">Grund</span>?</p>
<code>POST /bookings/{id}/cancellation &rarr; 201</code>
<aside class="notes"><strong>Meine Antwort:</strong> Beide sind vertretbar. DELETE ist die einfachste idempotente Variante, aber der Storno-Grund hat keinen Platz und die Buchung &bdquo;verschwindet&ldquo;. Die Stornierung als eigene Ressource tr&auml;gt Grund, Zeitpunkt und Geb&uuml;hr, bleibt per GET nachlesbar und l&auml;sst die Buchung als Historie stehen. Preis: POST ist nicht idempotent, der zweite Aufruf braucht 409 oder einen Idempotency-Key. PATCH auf den Status ist die dritte Variante, versteckt aber die Zustandsmaschine im Body.<br><strong>Spicy:</strong> &bdquo;Wir l&ouml;schen einfach&ldquo; funktioniert, bis der Kundenservice fragt, warum die Buchung weg ist und wer das war.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Umbuchung: PUT, PATCH oder Storno + Neu?</h3>
<p>Nur der Flug &auml;ndert sich. Was passiert mit den <span class="hl">weggelassenen Feldern</span>?</p>
<code>PUT /bookings/{id}/flight</code>
<aside class="notes"><strong>Meine Antwort:</strong> PUT ersetzt die ganze Ressource, weggelassene Felder werden gel&ouml;scht. Wer nur den Flug &auml;ndern will, meint PATCH oder eine Sub-Ressource (<code>PUT /bookings/{id}/flight</code>). Die Sub-Ressource ist idempotent und macht sichtbar, was sich &auml;ndert. Storno plus Neubuchung braucht keine neue API, hat aber zwei Vorg&auml;nge mit Preis- und Verf&uuml;gbarkeitsrisiko dazwischen. Entscheidend ist, dass die Teams die Teilausf&uuml;hrung benennen: der neue Flug ist nicht verf&uuml;gbar, was dann?<br><strong>Spicy:</strong> Genau diese Frage ist die Saga in Story 6. Wer sie hier am Flipchart hatte, erkennt sie dort wieder.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Retry ohne Reue?</h3>
<p>Der Client schickt den Storno zweimal. <span class="hl">Was passiert</span>?</p>
<code>Idempotency-Key: 8f3a&hellip;</code>
<aside class="notes"><strong>Meine Antwort:</strong> Bei DELETE nichts Schlimmes, der Zustand ist derselbe. Bei POST legt der Server ohne Schutz eine zweite Stornierung an. Zwei Wege: fachliche Pr&uuml;fung (schon storniert, 409) oder ein Idempotency-Key im Header, mit dem der Server die Wiederholung erkennt und die erste Antwort noch einmal liefert. So arbeiten Stripe und die meisten Zahlungs-APIs.<br><strong>Spicy:</strong> Idempotenz ist die Voraussetzung daf&uuml;r, dass Retry &uuml;berhaupt erlaubt ist. Wer POST ohne Key wiederholt, bucht doppelt. Wer deshalb nicht wiederholt, verliert Buchungen im Timeout. Beides ist schlechter als f&uuml;nf Zeilen Server-Code.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Wann ist RPC legitim?</h3>
<p>Unsere Referenz hat <code>POST /admin/bulkhead-reset</code>. <span class="hl">Regelbruch</span>?</p>
<code>Abweichen mit Grund</code>
<aside class="notes"><strong>Meine Antwort:</strong> Ja, bewusst. Ein Knopf f&uuml;rs Dashboard, kein Fachobjekt, ein einziger Aufrufer, kein Retry-Problem, klar unter <code>/admin/</code> abgesetzt. Genau diese Begr&uuml;ndung sollen die Teams f&uuml;r ihre eigenen Ausnahmen liefern k&ouml;nnen. Grauzonen wie <code>POST /flights/search</code> (Filter zu lang f&uuml;r die URL) sind ebenso vertretbar, wenn dokumentiert ist, dass der Aufruf keine Seiteneffekte hat.<br><strong>Spicy:</strong> Ein Fehler ist nicht die Abweichung, sondern die Abweichung, die niemand begr&uuml;nden kann. Kommt im Recap der Bulkhead-Story (Story 5) wieder, wenn der Reset-Knopf zum ersten Mal gedr&uuml;ckt wird.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Br&uuml;cke: &bdquo;In Story 3 baut ihr <code>POST /booking/bookings</code>. Nehmt die Regeln mit.&ldquo; Vollst&auml;ndige Antworten: <code>docs/questions/story2.md</code>.
</aside>
