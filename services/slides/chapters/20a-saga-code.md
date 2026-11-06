<div class="page">

<p class="kicker">Saga</p>

## Der Orchestrator

<p class="subtitle">In Pseudo-Code, wie im Spickzettel.</p>

<div class="page-body codebody">

<div class="codeblock">
<div class="k">book(request):</div>
<div>  saga = { status: <span class="hi">PENDING</span>, booked: [] }</div>
<div class="gap">  try:</div>
<div>    for svc in [flight, hotel, car]:</div>
<div>      id = POST svc/bookings                 <span class="dim">// Forward</span></div>
<div>      saga.booked.push({ svc, id })</div>
<div>    saga.status = <span class="hi">COMPLETED</span></div>
<div class="gap">  catch error at step X:</div>
<div>    saga.status = <span class="hi">COMPENSATING</span></div>
<div>    for b in <span class="hi">reverse</span>(saga.booked):</div>
<div>      DELETE b.svc/bookings/{b.id}           <span class="dim">// Kompensation, ein Versuch</span></div>
<div>    saga.status = <span class="hi">FAILED</span></div>
</div>

<p class="codenote">R&uuml;ckw&auml;rts, genau einmal, ohne Netz. Was das kostet, kl&auml;ren wir im Recap.</p>

</div>

</div>

Note:
- Bewusst wie der Spickzettel im Dashboard (Story 6), nur aufs Skelett gek&uuml;rzt.
- Zeigen: <strong>reverse(booked)</strong>, Kompensation r&uuml;ckw&auml;rts.
- Zeigen: <strong>booked</strong> als eigene Liste, nur Gebuchtes wird kompensiert. Car wird nach dem Hotel-Fehler gar nicht angefasst.
- Achtung <strong>Timeout, aber gebucht</strong>: Ergebnis unbekannt, niemand storniert. Live: Hotel auf &bdquo;H&auml;ngt&ldquo;, nach 3 s ist die Saga <code>FAILED</code>, nach 10 s loggt Hotel trotzdem <code>hotel booking confirmed</code>. Abhilfe: Buchungs-ID von Booking oder Idempotency-Key.
- <strong>DELETE idempotent</strong>, Recap-Frage 3.
- Bewusst weggelassen: Retry mit Backoff, persistenter Saga-Log, Workflow Engine, fachlicher Pivot. Details in <code>docs/questions/story6.md</code>.
- Anker: <code>DELETE</code> l&auml;uft in 5xx? Ein Versuch, dann <code>FAILED</code>. Recap-Frage 1.
- Referenz: <code>services/booking/story6/saga/</code>, ca. 100 Zeilen Go.
- &Uuml;berleitung: &bdquo;Bewusst die einfachste Variante: synchron, in-memory, ein Versuch. Jetzt baut ihr sie.&ldquo;
