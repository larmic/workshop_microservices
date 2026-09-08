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
- Beide Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein.
- <strong>Health-Check, meine Antwort:</strong> Unser <code>/health</code> ist eher <strong>Liveness</strong>, &bdquo;Prozess lebt&ldquo;. Eine richtige <strong>Readiness</strong> m&uuml;sste pr&uuml;fen, ob alle Abh&auml;ngigkeiten erreichbar sind: Consul (Story 3), Backends (Story 4+), DB-Pool, Message-Broker. <strong>Spicy:</strong> Wer nur Liveness baut, freut sich morgens &uuml;ber ein gr&uuml;nes Dashboard auf einem Service, der seit Stunden alle Backend-Calls in den Timeout laufen l&auml;sst. <em>Healthy ist nicht dasselbe wie Useful.</em>
- <strong>Backend nicht erreichbar, meine Antwort:</strong> Zwei v&ouml;llig verschiedene Ursachen, und beide sind in Story 1 ungel&ouml;st. <strong>1. Der Dienst ist down</strong> (Absturz, Deploy, &Uuml;berlast): Der Aufruf l&auml;uft in den Fehler oder Timeout, der Booking-Service reicht ihn ungebremst durch. Das ist ein <strong>Resilience</strong>-Thema: Timeout, Retry, Circuit Breaker, Fallback. Kommt in <strong>Stories 4 bis 6</strong>. <strong>2. Der Dienst ist umgezogen</strong>: neue URL, neuer Port, andere Instanz. Unsere <code>FLIGHT_SERVICE_URL</code> zeigt jetzt ins Leere. Jede Adress&auml;nderung erzwingt ein Re-Deployment des Booking-Service. Bei mehreren Backends mal mehreren Aufrufern wird das zur Pflege-H&ouml;lle. Genau hier setzt <strong>Story 3, Service Discovery</strong> an: Services registrieren sich selbst (Consul), gefunden wird &uuml;ber den <em>logischen Namen</em>, nicht &uuml;ber eine hartkodierte URL. <strong>Spicy:</strong> &bdquo;Steht doch in der ENV&ldquo; funktioniert, bis der dritte Service umzieht und niemand mehr wei&szlig;, welche URL noch stimmt. Ausprobiert wird im Zweifel auf Prod.
- Weitere Aspekte (Logging, Config, OpenAPI, Error-Handling) bei Bedarf: <code>docs/questions/story1.md</code>.
