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
- Derselbe Ablauf steht ausf&uuml;hrlicher im Dashboard unter Story 6, &bdquo;Spickzettel&ldquo;. Wiedererkennung gewollt, hier auf das Skelett gek&uuml;rzt.
- Drei Knackpunkte: <strong>reverse(booked)</strong>, die Kompensation l&auml;uft in umgekehrter Reihenfolge, sonst verletzt man fachliche Reihenfolge-Annahmen. <strong>booked als eigene Liste</strong>, nur was wirklich gebucht wurde, wird kompensiert; Car wird nach dem Hotel-Fehler gar nicht angefasst. <strong>DELETE muss idempotent sein</strong>, bei einem Retry darf das Backend nicht erschrecken, Recap-Frage 3.
- Was hier bewusst fehlt und in Produktion dazugeh&ouml;rt: <strong>Retry mit Backoff</strong> f&uuml;r die Kompensation (transiente Fehler aussitzen, ohne Doppel-Storno). <strong>Persistenter Saga-Log</strong> f&uuml;r Crash-Recovery (Postgres, SQLite, Consul KV). <strong>Workflow Engines</strong> wie Temporal oder Camunda drehen das Modell um: Die Engine persistiert jeden Schritt, retryt, startet nach Crash neu. <strong>Pivot zur fachlichen Alternative</strong>: Gutschein statt Storno, wenn der Flug schon abgehoben ist. <strong>Observability</strong>: Status-Counter, Compensation-Erfolgsrate, Saga-ID in jedem Span.
- Diskussions-Anker: Was, wenn der <code>DELETE</code> selbst in 5xx l&auml;uft? Im Workshop genau ein Versuch, dann <code>FAILED</code>. Das ist die erste Recap-Frage.
- Referenz: <code>services/booking/story6/saga/</code>, dieselbe Logik in etwa 100 Zeilen Go.
- &Uuml;berleitung: &bdquo;Im Workshop bewusst die einfachste Variante: synchron, in-memory, ein Versuch. Jetzt baut ihr sie.&ldquo;
