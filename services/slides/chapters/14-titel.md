<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Circuit Breaker</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Wer ein totes Backend weiter anruft, stirbt <span class="accent">mit</span>.</p>
<p class="sub">Nach f&uuml;nf Fehlern macht der Aufrufer selbst dicht. Schnell scheitern statt im Timeout h&auml;ngen.</p>
<p class="source">Michael Nygard, <em>Release It!</em>: Der Circuit Breaker ist das Pattern, mit dem ein Aufrufer sich vor einem kranken Partner sch&uuml;tzt.</p>
</div>
<div class="split-figure">
<img src="./assets/circuitbreaker.svg" alt="Sicherungskasten mit f&uuml;nf Schutzschaltern, der zweite ist ausgel&ouml;st"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;In Story 3 haben wir gelernt, Services zu <em>finden</em>. Heute kl&auml;ren wir, was passiert, wenn wir einen gefunden haben, und er antwortet nicht.&ldquo;
- Die Grafik: ein Sicherungskasten, vier Automaten oben (CLOSED), der zweite ist gefallen (OPEN). Genau das macht der Aufrufer f&uuml;r jedes Backend einzeln.
- Nygard, Release It! (2007): Der Circuit Breaker ist das bekannteste seiner Stabilit&auml;tspatterns. Grundgedanke: Ein Aufrufer, der einen kaputten Partner weiter anruft, verbrennt seine eigenen Threads und reisst seine eigenen Aufrufer mit.
- Demo-Vorschau: Im Dashboard Flight auf &bdquo;Fehler&ldquo; stellen, ein paar Requests gegen <code>/booking/offers</code>, jeder h&auml;ngt 3 s im Timeout. Mit Circuit Breaker in Story 4: nach den ersten f&uuml;nf Fehlern sofort Fallback.
