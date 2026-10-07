<div class="page">

<p class="kicker">Story 4 &middot; Recap</p>

## Vier Fragen an euch

<div class="page-body">

<div class="numlist recap compact">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Was z&auml;hlt als Fehler?</p>
<h3>5xx ja, Timeout ja. Aber ein <code>404</code>, ist das wirklich ein <span class="hl">Backend-Problem</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">4xx &ne; kaputt</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Fallback &ne; Circuit Breaker</p>
<h3>Erster Fehler, Z&auml;hler auf 1, Breaker noch CLOSED. Kommt <span class="hl">jetzt schon</span> <code>flights: []</code>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Fehler &rarr; Fallback, OPEN &rarr; sofort</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Probe-Storm</p>
<h3>In <code>HALF_OPEN</code> lassen wir genau einen Call durch. Warum nicht <span class="hl">alle</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">CompareAndSwap</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Der Zustand lebt im RAM einer Instanz</p>
<h3>Zwei Replicas: Flight-Breaker von A ist OPEN. Soll B auch <span class="hl">dichtmachen</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">kein shared state</code>
</div>
</div>

</div>

</div>

Note:
- <strong>Fehler:</strong> Nein. 404 hei&szlig;t falsche ID, das Backend ist gesund. CB-relevant sind 5xx, Timeout, Connection refused. 4xx nicht, 429 diskutabel. Unser Code z&auml;hlt alles ab 400, in Produktion ein Bug.
- <strong>Fallback:</strong> Ja, schon beim ersten Fehler. Die leere Liste kommt bei jedem gescheiterten Call. Der CB &auml;ndert nur das Wie: CLOSED kostet bis 3 s, OPEN antwortet sofort mit <code>X-Circuit-Open</code>.
- <strong>Probe-Storm:</strong> Sonst legen alle Wartenden das erholende Backend sofort wieder um. Ein atomarer Slot l&auml;sst genau einen durch.
- <strong>Replicas:</strong> Nein, Zustand bleibt instanzlokal. B erreicht Flight vielleicht problemlos. Folge: Nach Neustart ist alles CLOSED, die ersten f&uuml;nf Calls kosten wieder 3 s.
- Reserve: Retry und CB nicht naiv kombinieren. Granularit&auml;t pro Service ist der Default. Lang: <code>docs/questions/story4.md</code>.
