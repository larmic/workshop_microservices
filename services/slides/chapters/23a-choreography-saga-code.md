## Choreography-Saga

<p class="subtitle">Event-Handler in Pseudo-Code</p>

<pre class="cheatsheet"><span class="cmd">// Forward bleibt wie Story 5 (Saga, sync Aufrufe).</span>
<span class="cmd">// Neu: Kompensation l&auml;uft asynchron &uuml;ber Events.</span>

<span class="cmd">on bookingStepFailed(sagaId, alreadyBooked):</span>
  saga.status = COMPENSATING
  for b in reverse(alreadyBooked):
    publish "CompensationRequested" {
      eventId, sagaId, service: b.svc, bookingId: b.id
    }
  schedule timeout(sagaId, expectedReplies = len(alreadyBooked))

<span class="cmd">on event "BookingCancelled" {sagaId, service}:</span>
  saga.markCompensated(service)
  if saga.allCompensated(): saga.status = FAILED

<span class="cmd">on event "CancellationFailed" {sagaId, service}:</span>
  saga.markCompensationFailed(service)
  saga.status = COMPENSATION_INCOMPLETE

<span class="cmd">on timeout(sagaId):</span>
  if saga still has open replies:
    saga.status = STUCK
    log.alert("operator eingreifen", sagaId)
</pre>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 6 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- Vier Knackpunkte hervorheben:
  - <strong>Forward bleibt synchron</strong>. Nur die Kompensation l&auml;uft &uuml;ber Events. Im Workshop bewusst, um den Vergleich zu Story 5 sauber zu halten.
  - <strong>reverse(alreadyBooked)</strong> &mdash; Storno in umgekehrter Reihenfolge, gleicher Grund wie in Story 5.
  - <strong>eventId</strong> ist der Dedup-Key. Backends m&uuml;ssen sich gegen doppelte Zustellung absichern &mdash; bei at-least-once-Bus die Realit&auml;t, bei Webhook-Retry ebenfalls.
  - <strong>timeout()</strong> ist das Sicherheitsnetz f&uuml;r ausbleibende Replies. Ohne dieses Netz h&auml;ngt die Saga still im <code>COMPENSATING</code> &mdash; niemand merkt etwas.
- Reference-Code: <code>services/booking/story6/</code> &mdash; Forward-Pfad wie Story 5, Compensation-Pfad via Webhook-Publish + Reply-Konsum.
- Diskussions-Anker: Was passiert, wenn der <code>publish</code>-POST selbst fehlschl&auml;gt? In Story 5 h&auml;tte Booking gemerkt, dass etwas nicht stimmt. In Story 6 sieht Booking nichts &mdash; siehe Recap-Frage 1.
