<div class="page">

<p class="kicker">Story 5 &middot; Recap</p>

<div class="page-head">

## F&uuml;nf Fragen an euch

<span class="badge">Bonus zu 03 &darr;</span>
</div>

<div class="page-body">

<div class="numlist recap compact dense">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Wir haben doch einen Circuit Breaker</p>
<h3>Was kann der Bulkhead, was der CB <span class="hl">nicht</span> kann?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">langsam &ne; kaputt</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Ein Rate Limit erlaubt hundert Aufrufe pro Sekunde</p>
<h3>Und wenn alle hundert <span class="hl">gleichzeitig</span> h&auml;ngen?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Rate &ne; Nebenl&auml;ufigkeit</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Wir sind non-blocking, Threads sind bei uns billig</p>
<h3>Brauchen wir den Bulkhead dann <span class="hl">&uuml;berhaupt</span> noch?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Engpass wandert</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Warum keine Warteschlange?</p>
<h3>W&auml;re kurz warten nicht <span class="hl">freundlicher</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Queue tarnt Last</code>
</div>
<div class="numlist-row">
<div class="numlist-num">05</div>
<div>
<p class="numlist-label">F&uuml;nf Replicas mal zehn Slots sind f&uuml;nfzig Aufrufe</p>
<h3><span class="hl">Wen</span> sch&uuml;tzt der Bulkhead dann?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Outbound &ne; Inbound</code>
</div>
</div>

</div>

</div>

Note:
- Erst diskutieren, dann Klick f&uuml;r die Merks&auml;tze.
- <strong>Circuit Breaker:</strong> CB reagiert auf Fehler, Bulkhead auf Ressourcen-Druck. Runde 1: Hotel fehlerfrei, die Zustandskarte bleibt CLOSED, trotzdem wartet alles. <strong>Spicy:</strong> Der CB sieht Fehler, nicht Latenz.
- <strong>Rate Limit:</strong> Z&auml;hlt Requests pro Zeit, der Bulkhead gleichzeitige Aufrufe. Probelauf und Runde 1: gleicher Takt, nur Hotel langsamer. 100/s bei pl&ouml;tzlich 5 s pro Call hei&szlig;t 500 h&auml;ngende Aufrufe. <strong>Spicy:</strong> Durchsatz ist nicht Nebenl&auml;ufigkeit.
- <strong>Non-blocking:</strong> Threads werden billig, Connection-Pool, Speicher und Sockets nicht. Der Engpass wandert, siehe Bonus-Folie.
- <strong>Warteschlange:</strong> Runde 1 war die Queue: Jede Suche kam sp&auml;ter. Der sofortige Reject ist das Backpressure-Signal, die Queue schluckt es.
- <strong>Replicas:</strong> Der Bulkhead im Aufrufer sch&uuml;tzt den Aufrufer, nicht das Backend: Hotel sieht 50 Calls. Das Backend sch&uuml;tzt sich selbst, mit Rate Limit, Load Shedding, <code>429</code> plus <code>Retry-After</code>. <strong>Spicy:</strong> Bulkhead allein ist eine halbe L&ouml;sung.
- Lang: <code>docs/questions/story5.md</code>
