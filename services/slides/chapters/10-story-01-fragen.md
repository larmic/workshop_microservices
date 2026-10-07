<div class="page">

<p class="kicker">Story 1 &middot; Recap</p>

## Zwei Fragen an euch

<div class="page-body">

<div class="numlist recap">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Health-Check</p>
<h3>Unser <code>/health</code> gibt 200 zur&uuml;ck. Hei&szlig;t das, der Service ist <span class="hl">gesund</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">healthy &ne; useful</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Backend nicht erreichbar</p>
<h3>Flight, Hotel oder Car antworten nicht. <span class="hl">Wann</span> passiert das, und wer f&auml;ngt es auf?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">down &ne; umgezogen</code>
</div>
</div>

</div>

</div>

Note:
- <strong>Health-Check:</strong> Unser <code>/health</code> ist eher <strong>Liveness</strong>, der Prozess lebt. <strong>Readiness</strong> ist mehr: pr&uuml;ft, ob die Abh&auml;ngigkeiten erreichbar sind (Consul, Backends, DB, Broker). Healthy ist nicht dasselbe wie Useful.
- Weitere Aspekte (Logging, Config, OpenAPI, Error-Handling) bei Bedarf: <code>docs/questions/story1.md</code>.
