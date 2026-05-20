## Saga

<p class="subtitle">Orchestrator in Pseudo-Code</p>

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

<span class="cmd">catch error at step X:</span>
  saga.status   = COMPENSATING
  saga.failedAt = X
  for b in reverse(booked):
    DELETE b.svc/bookings/{b.id}              // Kompensation
    mark step COMPENSATED
  saga.status = FAILED
</pre>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 5 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- Drei Knackpunkte hervorheben:
  - <strong>reverse(booked)</strong> &mdash; Kompensation in umgekehrter Reihenfolge. Sonst kompensiert man Schritte, die nie ausgef&uuml;hrt wurden, oder verletzt fachliche Reihenfolge-Annahmen (Auto kann nur storniert werden, wenn Hotel noch existiert).
  - <strong>booked als Snapshot</strong> statt aus <code>saga.steps</code> ableiten &mdash; entkoppelt &bdquo;was wurde gebucht&ldquo; vom Saga-Logging. Falls eine Status-Update-Operation fehlschl&auml;gt, ist der Rollback-Pfad immer noch korrekt.
  - <strong>DELETE muss idempotent sein</strong> &mdash; bei Retry darf der Service nicht erschrecken (siehe Recap-Frage 4).
- Reference-Code: <code>services/booking/story5/saga/</code> &mdash; identische Logik in ca. 100 Zeilen Go.
- Diskussions-Anker: Was, wenn der <code>DELETE</code> selbst in 5xx l&auml;uft? Im Workshop genau <em>einen</em> Versuch, ansonsten <code>FAILED</code>. Produktion braucht Retry mit Backoff (siehe Recap-Frage 1+2). Brille von dem Code aufsetzen: in den Schritten <code>saga.steps.push</code> liegt der State; ohne Persistenz ist alles weg, wenn Booking abst&uuml;rzt &mdash; siehe Recap-Frage 7.
