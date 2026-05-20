## Circuit Breaker

<p class="subtitle">Die State-Machine in Pseudo-Code</p>

<pre class="cheatsheet"><span class="cmd">STATE:</span>
  state     = CLOSED   // CLOSED | OPEN | HALF_OPEN
  failures  = 0
  openUntil = 0

<span class="cmd">call(service):</span>
  if state == OPEN and now &lt; openUntil:
    return fallback()                     // short-circuit
  if state == OPEN and now &gt;= openUntil:
    state = HALF_OPEN                     // probe slot frei

  try service.invoke(timeout = 3s):
    success &rarr; failures = 0
              state    = CLOSED
              return result
    failure &rarr; failures++
              if failures &gt;= 5:
                state     = OPEN
                openUntil = now + 30s
              return fallback()
</pre>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 3 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- Drei Knackpunkte hervorheben:
  - <strong>Timeout im Aufruf selbst</strong> (3 s) &mdash; ohne den bringt der CB nichts, weil er nie als &bdquo;failure&ldquo; getriggert wird.
  - <strong>failures &gt;= 5</strong> &mdash; count-basiert. Reicht im Workshop, in Produktion oft rate-basiert &uuml;ber Sliding Window.
  - <strong>HALF_OPEN ohne Probe-Schutz</strong> &mdash; in der Skizze l&auml;sst <em>jeder</em> Aufruf nach Ablauf die Probe durch. Im echten Code muss das atomar gegen Probe-Storm gesch&uuml;tzt sein.
- Real Code: <code>services/booking/story3/circuitbreaker/circuitbreaker.go</code> &mdash; ca. 80 Zeilen Go.
- Diskussions-Anker: Was z&auml;hlt als <em>failure</em>? HTTP 5xx ja, Timeout ja, Connection refused ja &mdash; aber HTTP 4xx? (Antwort: nein, der Aufrufer hat Mist gemacht, Backend ist gesund.)
