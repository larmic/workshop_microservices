## Story 1 &mdash; Recap

<p class="subtitle">Zwei Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Health-Check</h3>
<p>Unser <code>/health</code> gibt 200 zur&uuml;ck. Hei&szlig;t das, der Service ist <span class="hl">gesund</span>?</p>
<code>healthy &ne; useful</code>
<aside class="notes"><strong>Meine Antwort:</strong> Unser <code>/health</code> ist eher <strong>Liveness</strong> &mdash; &bdquo;Prozess lebt&ldquo;. Eine richtige <strong>Readiness</strong> m&uuml;sste pr&uuml;fen, ob alle Abh&auml;ngigkeiten erreichbar sind: Consul (Story 2), Backends (Story 3+), DB-Pool, Message-Broker.<br><strong>Spicy:</strong> Wer nur Liveness baut, freut sich morgens &uuml;ber ein gr&uuml;nes Dashboard auf einem Service, der seit Stunden alle Backend-Calls in den Timeout laufen l&auml;sst. <em>Healthy ist nicht dasselbe wie Useful.</em></aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Backend nicht erreichbar?</h3>
<p>Flight, Hotel oder Car antworten nicht. <span class="hl">Wann</span> passiert das &mdash; und wer f&auml;ngt es auf?</p>
<code>down &ne; umgezogen</code>
<aside class="notes"><strong>Meine Antwort:</strong> Zwei v&ouml;llig verschiedene Ursachen &mdash; und beide sind in Story 1 ungel&ouml;st:<br><strong>1. Der Dienst ist down</strong> (Absturz, Deploy, &Uuml;berlast). Der Aufruf l&auml;uft in den Fehler oder Timeout, der Booking-Service reicht ihn ungebremst durch. Das ist ein <strong>Resilience</strong>-Thema &mdash; Timeout, Retry, Circuit Breaker, Fallback. Kommt in <strong>Stories 3&ndash;5</strong>.<br><strong>2. Der Dienst ist umgezogen</strong> &mdash; neue URL, neuer Port, andere Instanz. Unsere <code>FLIGHT_SERVICE_URL</code> zeigt jetzt ins Leere. Jede Adress&auml;nderung erzwingt ein Re-Deployment des Booking-Service. Bei mehreren Backends &times; mehreren Aufrufern wird das zur Pflege-H&ouml;lle. Genau hier setzt <strong>Story 2 &mdash; Service Discovery</strong> an: Services registrieren sich selbst (Consul), gefunden wird &uuml;ber den <em>logischen Namen</em>, nicht &uuml;ber eine hartkodierte URL.<br><strong>Spicy:</strong> &bdquo;Steht doch in der ENV&ldquo; funktioniert, bis der dritte Service umzieht und niemand mehr wei&szlig;, welche URL noch stimmt &mdash; ausprobiert wird im Zweifel auf Prod.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Weitere Aspekte (Logging, Config, OpenAPI, Error-Handling) bei Bedarf: <code>docs/questions/story1.md</code>.
</aside>
