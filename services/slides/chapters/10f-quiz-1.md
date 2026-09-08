<div class="page">

<p class="kicker">Quiz &middot; 1 von 4</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">GET</span> /booking/offers?from=BRE&amp;to=LIS&amp;date=2026-05-14</div>

<div class="quiz-answer fragment">
<div class="callout">RESTful</div>
<p><code>offers</code> ist ein Substantiv, GET ver&auml;ndert nichts, und die Auswahl steht in Query-Parametern statt in eigenen Endpoints. Genau das meint &bdquo;Sub-Ressourcen und Filter&ldquo;, Regel 6 von der Regel-Folie eben. Und genau das habt ihr in Story 1 schon gebaut.</p>
<code class="quiz-take">Filter &ne; Endpoint</code>
</div>

</div>

</div>

Note:
- Ablauf f&uuml;r alle drei: Endpoint zeigen, Handzeichen &bdquo;RESTful?&ldquo; abfragen, eine Person aus der Minderheit begr&uuml;nden lassen, dann Fragment aufl&ouml;sen.
- Erwartung: einige sagen &bdquo;nicht RESTful&ldquo;, weil Query-Parameter nach RPC aussehen. Pointe: Filter auf eine Sammlung sind genau der Zweck von Query-Parametern (Regel 6, Sub-Ressourcen und Filter). Eigene Endpoints pro Filterkombination w&auml;ren das Anti-Pattern.
- Bezug: <code>GET /booking/offers</code> ist der Aggregations-Endpoint aus Story 1. Die Filter sind hier hinzugedacht, das Prinzip ist dasselbe.
- &Uuml;berleitung: &bdquo;Das war der Einstieg. Jetzt der Status-Code.&ldquo;
