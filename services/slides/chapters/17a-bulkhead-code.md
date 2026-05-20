## Bulkhead

<p class="subtitle">Semaphore in Pseudo-Code</p>

<pre class="cheatsheet"><span class="cmd">STATE:</span>
  max      = { flight: 5, hotel: 5, car: 5 }   // pro Service
  inFlight = { flight: 0, hotel: 0, car: 0 }
  rejected = { flight: 0, hotel: 0, car: 0 }

<span class="cmd">call(service):</span>
  if inFlight[service] &gt;= max[service]:
    rejected[service]++
    return 503                              // sofort ablehnen
  inFlight[service]++
  try:
    return service.invoke()
  finally:
    inFlight[service]--

<span class="cmd">// Aggregator (Booking) ruft alle drei parallel auf:</span>
//   ein langsamer Hotel f&uuml;llt nur seinen Pool,
//   Flight + Car liefern weiter normal.
</pre>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 4 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- Drei Knackpunkte hervorheben:
  - <strong>Check &amp; Increment m&uuml;ssen atomar</strong> sein (Mutex / Compare-and-Set / Semaphore-Primitive). Sonst rutschen unter Last mehr Calls durch als erlaubt &mdash; die Isolation ist dahin.
  - <strong>release() im finally</strong> &mdash; vergessen hei&szlig;t Slot-Lecks, der Pool f&uuml;llt sich &uuml;ber die Zeit, der Bulkhead &ouml;ffnet nie wieder.
  - <strong>Ein Bulkhead pro Downstream</strong>. Ein gemeinsamer Pool h&auml;tte das Pattern ad absurdum gef&uuml;hrt.
- Reference-Code: <code>services/booking/story4/bulkhead/bulkhead.go</code> &mdash; ca. 60 Zeilen Go mit <code>chan struct{}</code> als Semaphore.
- Diskussions-Anker: Sollte ein Bulkhead-Reject als CB-Failure z&auml;hlen? (Antwort im Recap: nein &mdash; CB und Bulkhead sind komplement&auml;r und unabh&auml;ngig.)
