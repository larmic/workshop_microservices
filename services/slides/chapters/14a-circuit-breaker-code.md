<div class="page">

<p class="kicker">Circuit Breaker</p>

## Die State-Machine

<p class="subtitle">Die Regeln, nach denen Booking gleich spielt.</p>

<div class="page-body codebody">

<div class="codeblock">
<div class="k">call(service):</div>
<div>  if state == OPEN and now &lt; openUntil:</div>
<div>    return fallback()                    <span class="dim">// short-circuit</span></div>
<div>  if state == OPEN and now &ge; openUntil:</div>
<div>    state = HALF_OPEN</div>
<div>  if state == HALF_OPEN and not probeInFlight.CAS(false, true):</div>
<div>    return fallback()                    <span class="dim">// nur einer probt</span></div>
<div class="gap">  try service.invoke(<span class="hi">timeout = 3s</span>):</div>
<div>    success &rarr; failures = 0; state = CLOSED</div>
<div>    failure &rarr; failures++</div>
<div>               if failures &ge; <span class="hi">5</span>: state = OPEN, openUntil = now + <span class="hi">30s</span></div>
<div>               return fallback()</div>
</div>

<p class="codenote">Drei Zahlen, die den Breaker ausmachen: Timeout, Schwelle, Wartezeit. Alle drei sind Entscheidungen, keine Konstanten.</p>

</div>

</div>

Note:
- <strong>Kernbotschaft:</strong> Das tippen wir nicht ab, wir spielen es gleich.
  - Pl&auml;tze auf dem Board = <code>failures</code>, Zustandskarte = <code>state</code>, Probe-Feld = <code>probeInFlight</code>, <code>flights: []</code> auf der Antwort = <code>fallback()</code>.
- <strong>Drei Zahlen:</strong> 3 s, 5 in Folge, 30 s. Jede h&auml;ngt am Backend und an eurem Budget.
- <strong>Reset bei Erfolg:</strong> Jeder Erfolg setzt auf null. Darum &bdquo;in Folge&ldquo;: Ein wackliges Flight, das jedes zweite Mal antwortet, &ouml;ffnet den Schalter nie.
- <strong>Probe-Slot:</strong> In HALF_OPEN darf genau einer durch (CAS), alle anderen bekommen den Fallback. Scheitert die Probe: sofort wieder OPEN, die 30 s laufen von vorn.
- Referenz: <code>services/booking/story4/circuitbreaker/circuitbreaker.go</code>. Weiter: Booking bestimmen, Stapel &uuml;bergeben.
