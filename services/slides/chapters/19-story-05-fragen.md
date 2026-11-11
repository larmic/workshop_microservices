<div class="page">

<p class="kicker">Story 5 &middot; Recap</p>

<div class="page-head">

## Vier Fragen an euch

<span class="badge">Bonus zu 02 &darr;</span>
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
<p class="numlist-label">Wir sind non-blocking, Threads sind bei uns billig</p>
<h3>Brauchen wir den Bulkhead dann <span class="hl">&uuml;berhaupt</span> noch?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Engpass wandert</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Warum keine Warteschlange?</p>
<h3>W&auml;re kurz warten nicht <span class="hl">freundlicher</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Queue tarnt Last</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
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
- <strong>Erst diskutieren</strong>, dann Klick f&uuml;r die Merks&auml;tze.
- <strong>01 Circuit Breaker:</strong> Der CB reagiert auf Fehler, der Bulkhead auf Ressourcen-Druck.
  - Runde 1: Hotel war fehlerfrei, der Breaker blieb CLOSED, trotzdem standen vier hinter den St&uuml;hlen.
  - Spicy: Der CB sieht Fehler, nicht Latenz.
- <strong>02 Non-blocking:</strong> Threads werden billig, Connection-Pool, Speicher und Sockets nicht. Der Engpass wandert, siehe Bonus-Folie.
- <strong>03 Warteschlange:</strong> Runde 1 war die Queue.
  - Sie w&auml;chst, solange Hotel langsam ist: Jede Nummer stand einer mehr hinter den St&uuml;hlen.
  - Warten macht die Antwort nur sp&auml;ter: Karte 14 erst nach 6 s, so lange wartet kein Client. Er l&auml;dt neu, und die Last w&auml;chst.
  - Wer wartet, h&auml;lt Verbindung und Speicher. Das Problem wandert nur vor den Pool.
  - Der sofortige Reject ist das Signal (Backpressure). Die Queue versteckt es, bis die Timeouts kommen.
  - Okay ist nur eine kurze, begrenzte Wartezeit gegen Spitzen. In Resilience4j <code>maxWaitDuration</code>, Default 0.
  - Spicy: Eine unbegrenzte Queue vor einem langsamen Downstream ist ein Timeout auf Raten.
- <strong>04 Replicas:</strong> Der Bulkhead sch&uuml;tzt den Aufrufer, nicht das Backend: Hotel sieht 50 Calls.
  - Im Spiel geh&ouml;rten die St&uuml;hle Booking, Hotel wusste nichts davon.
  - Das Backend sch&uuml;tzt sich selbst: Rate Limit, Load Shedding, <code>429</code> plus <code>Retry-After</code>.
  - Spicy: Bulkhead allein ist eine halbe L&ouml;sung.
- Lang: <code>docs/questions/story5.md</code>
