<div class="page">

<p class="kicker">Saga</p>

## Der Orchestrator

<p class="subtitle">In Pseudo-Code, wie im Spickzettel.</p>

<div class="page-body codebody">

<div class="codeblock">
<div class="k">book(request):</div>
<div>  saga = { status: <span class="hi">PENDING</span>, steps: [] }</div>
<div class="gap">  try:</div>
<div>    for svc in [flight, hotel, car]:</div>
<div>      id = POST svc/bookings                 <span class="dim">// buchen</span></div>
<div>      saga.steps.push({ svc, id })</div>
<div>    saga.status = <span class="hi">COMPLETED</span></div>
<div class="gap">  catch error at step X:</div>
<div>    saga.status = <span class="hi">COMPENSATING</span></div>
<div>    for step in <span class="hi">reverse</span>(saga.steps):</div>
<div>      DELETE step.svc/bookings/{step.id}     <span class="dim">// Kompensation, ein Versuch</span></div>
<div>    saga.status = <span class="hi">FAILED</span></div>
</div>

<p class="codenote">R&uuml;ckw&auml;rts, genau einmal, ohne Netz. Was das kostet, kl&auml;ren wir im Recap.</p>

</div>

</div>

Note:
- <strong>Zeigen:</strong>
  - <code>steps</code>: Nur Gebuchtes wird kompensiert. Car wird nach dem Hotel-Fehler nie angefasst.
  - <code>reverse</code>: Die Kompensation l&auml;uft r&uuml;ckw&auml;rts.
  - Ein Versuch: L&auml;uft der <code>DELETE</code> in einen 5xx-Fehler, steht die Saga auf <code>FAILED</code> (Recap 01).
- <strong>Achtung, Timeout und trotzdem gebucht:</strong> Hotel auf &bdquo;Timeout&ldquo;. Nach 3 s ist die Saga <code>FAILED</code>, nach 10 s bucht Hotel trotzdem. Abhilfe: eine Buchungs-ID von Booking oder ein Idempotency-Key.
- <strong><code>DELETE</code> idempotent:</strong> Recap 02.
- <strong>Im echten Code zus&auml;tzlich:</strong> der Pre-Check (Breaker OPEN, keine Saga) und die Kompensation im eigenen Context.
- <strong>Bewusst weggelassen:</strong> Retry mit Backoff, persistenter Saga-Log, Workflow-Engine, fachlicher Pivot. Mehr in <code>docs/questions/story6.md</code>.
- <strong>Referenz:</strong> <code>services/booking/story6/handler/booking.go</code>, <code>CreateBookingHandler</code> und <code>compensate</code>.
- <strong>&Uuml;berleitung:</strong> &bdquo;Bewusst die einfachste Variante: synchron, in-memory, ein Versuch. Jetzt baut ihr sie.&ldquo;
