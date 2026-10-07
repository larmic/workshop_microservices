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

