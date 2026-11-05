<div class="page">

<p class="kicker">Quiz &middot; 1 von 5</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">GET</span> /booking/offers?from=BRE&amp;to=LIS&amp;date=2026-05-14</div>

<div class="quiz-answer fragment">
<div class="callout">RESTful</div>
<p><code>offers</code> ist ein Substantiv, GET ver&auml;ndert nichts, und die Auswahl steht in Query-Parametern statt in eigenen Endpoints. Genau das meint &bdquo;Sub-Ressourcen und Filter&ldquo;, Regel 5 von der Regel-Folie eben. Und genau das habt ihr in Story 1 schon gebaut.</p>
<code class="quiz-take">Filter &ne; Endpoint</code>
</div>

</div>

</div>

Note:
- <strong>Mechanik (gilt f&uuml;r alle f&uuml;nf):</strong> Daumen hoch = RESTful, runter = nicht, seitlich = kommt drauf an. Eine Stimme aus der Minderheit begr&uuml;nden lassen, dann aufl&ouml;sen.
- <strong>Kernbotschaft:</strong> Substantiv, GET ohne Seiteneffekt, Auswahl &uuml;ber Query-Parameter. Kein eigener Endpoint pro Route oder Datum.
- <strong>Bezug:</strong> Den Endpoint kennt ihr aus Story 1 (dort noch ohne Query-Parameter), Regel 5 auf der Regel-Folie.
- <strong>&Uuml;berleitung:</strong> Das war zum Aufw&auml;rmen. Beim n&auml;chsten sieht die URL genauso sauber aus, achtet auf die Antwort.
