<div class="page">

<p class="kicker">Circuit Breaker</p>

## Die State-Machine

<p class="subtitle">In Pseudo-Code, wie im Spickzettel.</p>

<div class="page-body codebody">

<div class="codeblock">
<div class="k">call(service):</div>
<div>  if state == OPEN and now &lt; openUntil:</div>
<div>    return fallback()                    <span class="dim">// short-circuit</span></div>
<div>  if state == OPEN and now &ge; openUntil:</div>
<div>    state = HALF_OPEN                    <span class="dim">// ein Probe-Slot frei</span></div>
<div class="gap">  try service.invoke(<span class="hi">timeout = 3s</span>):</div>
<div>    success &rarr; failures = 0; state = CLOSED</div>
<div>    failure &rarr; failures++</div>
<div>              if failures &ge; <span class="hi">5</span>: state = OPEN, openUntil = now + <span class="hi">30s</span></div>
<div>              return fallback()</div>
</div>

<p class="codenote">Drei Zahlen, die den Breaker ausmachen: Timeout, Schwelle, Wartezeit. Alle drei sind Entscheidungen, keine Konstanten.</p>

</div>

</div>

Note:
- Identischer Pseudo-Code steht im Dashboard unter Story 4, &bdquo;Spickzettel&ldquo;. Wiedererkennung gewollt.
- <strong>Timeout im Aufruf selbst</strong> (3 s): Ohne ihn bringt der CB nichts, weil ein h&auml;ngender Call nie als Fehler z&auml;hlt.
- <strong>Schwelle 5</strong>: count-basiert, reicht im Workshop. In Produktion meist rate-basiert &uuml;ber ein Sliding Window, sonst &ouml;ffnen bei hohem Durchsatz f&uuml;nf Fehler in 100 ms den Breaker, obwohl 99,99 Prozent der Requests gesund waren.
- <strong>Wartezeit 30 s</strong>: Danach HALF_OPEN. In der Skizze l&auml;sst jeder Aufruf nach Ablauf die Probe durch. Im echten Code muss das atomar gegen den Probe-Storm gesch&uuml;tzt sein, Recap-Frage 4.
- Was die Referenz zus&auml;tzlich hat, hier bewusst weggelassen: Slow-Call-Detection (langsam gilt als kaputt, auch bei 200), Exception-Klassifizierung (nur 5xx z&auml;hlt, 4xx nicht), Metriken pro Zustandswechsel, Decorator-Kette Retry, CB, Timeout, Fallback. Resilience4j macht aus alldem eine Annotation.
- Referenz: <code>services/booking/story4/circuitbreaker/circuitbreaker.go</code>, etwa 80 Zeilen Go.
