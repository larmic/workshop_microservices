<div class="page">

<p class="kicker">Bulkhead</p>

## Sieben Zeilen

<p class="subtitle">Deshalb tippen wir sie heute nicht ab.</p>

<div class="page-body codebody">

<div class="codeblock fit">
<div class="k">call(service):</div>
<div>  if inProgress[service] &ge; <span class="hi">max[service]</span>:</div>
<div>    rejected[service]++</div>
<div>    return fallback()   <span class="dim">// [] + X-Fallback, X-Bulkhead-Full</span></div>
<div>  inProgress[service]++</div>
<div>  try:     return service.invoke()</div>
<div>  finally: inProgress[service]--</div>
</div>

<p class="codenote">Die einzige interessante Stelle ist <code>max</code>. Woher kommt die Zahl?</p>

</div>

</div>

Note:
- <strong>Kernbotschaft:</strong> Der Code ist trivial, die Zahl darin nicht.
- <strong>Drei Knackpunkte</strong>, falls jemand nachbaut:
  - Check und Increment atomar (Mutex, CAS, Semaphore), sonst rutschen unter Last mehr Calls durch.
  - Decrement ins finally, sonst leckt der Pool und &ouml;ffnet nie wieder.
  - Ein Bulkhead pro Downstream, nie ein gemeinsamer.
- <strong>Im Spiel gleich:</strong> Platz bei Booking = Slot, drei pro Service = <code>max[service]</code>, Default-Chip von Booking = <code>fallback()</code>.
  - Stufe 2 ohne Bulkhead: Wer an der T&uuml;r steht, ist die Queue des gemeinsamen Pools. Der Code hier wartet nie.
- Referenz: <code>services/booking/story5/bulkhead/bulkhead.go</code>, knapp 100 Zeilen Go, <code>chan struct{}</code> als Semaphore. Im Dashboard unter Story 5 mit dem Burst-Knopf gegen langsame Backends treiben.
- &Uuml;berleitung: &bdquo;Die Zahl kl&auml;ren wir jetzt. Nicht am Rechner, sondern mit Chips.&ldquo;
