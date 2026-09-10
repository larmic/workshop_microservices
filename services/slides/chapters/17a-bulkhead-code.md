<div class="page">

<p class="kicker">Bulkhead</p>

## Sieben Zeilen

<p class="subtitle">Deshalb tippen wir sie heute nicht ab.</p>

<div class="page-body codebody">

<div class="codeblock fit">
<div class="k">call(service):</div>
<div>  if inProgress[service] &ge; max[service]:</div>
<div>    rejected[service]++</div>
<div>    return 503</div>
<div>  inProgress[service]++</div>
<div>  try:     return service.invoke()</div>
<div>  finally: inProgress[service]--</div>
</div>

<p class="codenote">Die einzige interessante Stelle ist <code>max</code>. Woher kommt die Zahl?</p>

</div>

</div>

Note:
- Drei Knackpunkte, falls jemand nachbaut: Check und Increment m&uuml;ssen atomar sein (Mutex, Compare-and-Set, Semaphore), sonst rutschen unter Last mehr Calls durch als erlaubt. Das Decrement geh&ouml;rt ins finally, sonst leckt der Pool und &ouml;ffnet nie wieder. Ein Bulkhead pro Downstream, nie ein gemeinsamer.
- Referenz: <code>services/booking/story5/bulkhead/bulkhead.go</code>, etwa 60 Zeilen Go mit <code>chan struct{}</code> als Semaphore. Im Dashboard unter Story 5 l&auml;sst sich die Referenz gegen langsame Backends treiben (Burst-Knopf).
- &Uuml;berleitung: &bdquo;Die Zahl kl&auml;ren wir jetzt. Nicht am Rechner, sondern mit Bechern und Chips.&ldquo;
