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
- <strong>Health-Check:</strong> Unser <code>/health</code> ist eher <strong>Liveness</strong>, der Prozess lebt. <strong>Readiness</strong> fragt: Kann diese Instanz jetzt bedienen? Eigene Ressourcen wie Warm-up, eigene DB und Pool, nicht die Downstream-Services. Wer Flight in die Readiness nimmt, schaltet bei einem Flight-Ausfall den ganzen Booking-Service ab. Healthy ist nicht dasselbe wie Useful.
- <strong>Backend weg:</strong> Zwei F&auml;lle. <strong>Umgezogen</strong> (neue IP, Skalierung): Die feste URL aus der Env-Variable zeigt ins Leere, das l&ouml;st Story 3. <strong>Wirklich kaputt:</strong> Jede Anfrage h&auml;ngt im Timeout, das l&ouml;st Story 4. Heute f&auml;ngt es niemand auf, mit Absicht.
- Weitere Aspekte (Logging, Config, OpenAPI, Error-Handling) bei Bedarf: <code>docs/questions/story1.md</code>.
- <strong>&Uuml;berleitung:</strong> Unser Service antwortet per HTTP. Aber benutzen wir HTTP auch richtig?
