<!-- .slide: data-background-image="./assets/choreography_saga.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Choreography-Saga

<p class="subtitle">Event-Handler in Pseudo-Code</p>

<pre class="cheatsheet"><span class="cmd">// Forward bleibt wie Story 5 (Saga, sync Aufrufe).</span>
<span class="cmd">// Neu: Kompensation l&auml;uft asynchron, fire-and-forget.</span>

<span class="cmd">on bookingStepFailed(sagaId, alreadyBooked):</span>
  saga.status = COMPENSATING
  for b in reverse(alreadyBooked):
    POST /events/compensation {
      eventId, sagaId, service: b.svc, bookingId: b.id
    }
    <span class="cmd">// Backend antwortet 202 Accepted, Booking ist hier fertig.</span>
    saga.markCompensated(b.svc)
  saga.status = FAILED
  <span class="cmd">// Booking wei&szlig; NICHT, ob der fachliche Rollback geklappt hat.</span>
  <span class="cmd">// Reply, Timeout, STUCK siehe Recap-Frage 1 + 3 (Bonus).</span>

<span class="cmd">// Backend-Seite (flight / hotel / car):</span>
<span class="cmd">on POST /events/compensation:</span>
  validate(eventId, sagaId, bookingId)
  respond 202 Accepted
  async:
    rollback(bookingId)
    log "compensation done"
</pre>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 6 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- Vier Knackpunkte hervorheben:
  - <strong>Forward bleibt synchron</strong>. Nur die Kompensation l&auml;uft &uuml;ber Events. Im Workshop bewusst, um den Vergleich zu Story 5 sauber zu halten.
  - <strong>reverse(alreadyBooked)</strong>: Storno in umgekehrter Reihenfolge, gleicher Grund wie in Story 5.
  - <strong>eventId</strong> ist der Dedup-Key. Konzept zeigen, persistente Speicherung im Workshop bewusst weggelassen. In Produktion w&auml;re das Pflicht (at-least-once-Bus, Webhook-Retry).
  - <strong>Booking ist mit dem Event-Versand fertig</strong>. Saga geht direkt auf <code>FAILED</code>, der Kunde bekommt seine Antwort. Was, wenn der POST fehlschl&auml;gt oder der Rollback im Backend kaputt geht? Booking sieht nichts. Genau dieser Punkt wird in Recap-Frage 1 diskutiert.
- Reference-Code: <code>services/booking/story6/</code>. Forward-Pfad wie Story 5, Compensation-Pfad via Webhook-POST ohne Reply.
- Diskussions-Anker: Was passiert, wenn der <code>POST</code> selbst fehlschl&auml;gt? In Story 5 h&auml;tte Booking es synchron gemerkt. In Story 6 (Schmalspur) sieht Booking nichts. Reply-Events, Timeout und <code>STUCK</code>-Detection sind der n&auml;chste Schritt zur Production-Reife, im Workshop als Bonus / Diskussion.
