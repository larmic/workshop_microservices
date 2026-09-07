## Choreography-Saga

<p class="subtitle">Event-Handler in Pseudo-Code</p>

<div class="cols">
<div>

<pre class="cheatsheet"><span class="cmd">STATE:</span>
  saga   = { id, status: PENDING, steps: [] }
  booked = []                                    // f&uuml;r Rollback

<span class="cmd">try:</span>
  flight = POST flight/bookings { ... }
  booked.push({ svc: "flight", id: flight.id })
  saga.steps.push({ svc: "flight", status: BOOKED })

  hotel = POST hotel/bookings { ... }
  booked.push({ svc: "hotel",  id: hotel.id })
  saga.steps.push({ svc: "hotel",  status: BOOKED })

  car = POST car/bookings { ... }
  saga.steps.push({ svc: "car",    status: BOOKED })

  saga.status = COMPLETED
</pre>

</div>
<div>

<pre class="cheatsheet"><span class="cmd">catch error at step X:</span>
  saga.status   = COMPENSATING
  saga.failedAt = X
  for b in reverse(booked):
    POST b.svc/events/compensation             <span class="cmd">// statt DELETE: Event</span>
      { eventId, sagaId, bookingId: b.id }
    mark step COMPENSATED                      <span class="cmd">// dispatched, kein Reply</span>
  saga.status = FAILED
  <span class="cmd">// Booking wei&szlig; NICHT, ob der fachliche Rollback geklappt hat.</span>
  <span class="cmd">// Reply, Timeout, STUCK siehe Recap-Frage 1 + 3 (Bonus).</span>

<span class="cmd">// Backend (flight / hotel / car):</span>
<span class="cmd">on POST /events/compensation:</span>
  validate(eventId, sagaId, bookingId)
  respond 202 Accepted
  async: rollback(bookingId)
</pre>

</div>
</div>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 7 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- <strong>Diff zu Story 6: genau zwei Zeilen.</strong> <code>POST .../events/compensation</code> statt <code>DELETE</code>, plus Kommentar <code>// dispatched, kein Reply</code>. Der Rest (STATE, try, catch) ist Wort f&uuml;r Wort identisch. Visuell darauf zeigen.
- Knackpunkte hervorheben:
  - <strong>Forward bleibt synchron</strong>. Der <code>try</code>-Block ist 1:1 Story 6. Nur die Kompensation l&auml;uft &uuml;ber Events. Im Workshop bewusst, um den Vergleich sauber zu halten.
  - <strong>reverse(booked)</strong>: identisch zu Story 6, gleiche fachliche Reihenfolge-Annahme. In der Schmalspur wartet Booking zwar nicht, aber das Event soll trotzdem in der erwarteten Storno-Reihenfolge rausgehen.
  - <strong>eventId</strong> ist der Dedup-Key. Konzept zeigen, persistente Speicherung im Workshop bewusst weggelassen. In Produktion w&auml;re das Pflicht (at-least-once-Bus, Webhook-Retry).
  - <strong><code>mark step COMPENSATED</code> hei&szlig;t hier nur: Event ist raus.</strong> Nicht: Backend hat fachlich storniert. Saga geht direkt auf <code>FAILED</code>, der Kunde bekommt seine Antwort. Was, wenn der POST fehlschl&auml;gt oder der Rollback im Backend kaputt geht? Booking sieht nichts. Genau dieser Punkt wird in Recap-Frage 1 diskutiert.
  - <strong>Backend antwortet sofort 202 Accepted</strong> und macht den Rollback in einer Goroutine. Aus Sicht des Senders fire-and-forget.
- Reference-Code: <code>services/booking/story7/</code>. Forward-Pfad wie Story 6, Compensation-Pfad via Webhook-POST ohne Reply.
- Diskussions-Anker: Was passiert, wenn der <code>POST</code> selbst fehlschl&auml;gt? In Story 6 h&auml;tte Booking es synchron gemerkt. In Story 7 (Schmalspur) sieht Booking nichts. Reply-Events, Timeout und <code>STUCK</code>-Detection sind der n&auml;chste Schritt zur Production-Reife, im Workshop als Bonus / Diskussion.
