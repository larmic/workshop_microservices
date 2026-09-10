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
<code class="numlist-pill fragment" data-fragment-index="1">Grenze f&auml;llt weg</code>
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
- Alle f&uuml;nf Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein. Frage 03 hat eine Bonus-Folie darunter.
- <strong>Circuit Breaker, meine Antwort:</strong> Verschiedene Probleme. CB: &bdquo;Backend ist krank, ich versuche es gar nicht mehr&ldquo;, Reaktion auf die Fehlerrate. Bulkhead: &bdquo;Ich verbrenne maximal N Slots f&uuml;r dieses Backend, egal wie viel Last reinkommt&ldquo;, Reaktion auf Ressourcen-Druck. Das Szenario, das nur Bulkhead l&ouml;st: Hotel antwortet in 2 s fehlerfrei, der CB bleibt CLOSED, unter Last laufen alle Threads in Hotel-Calls auf. <strong>Spicy:</strong> &bdquo;Wir haben doch schon einen CB&ldquo; ist die Standardfalle. Der CB sieht Fehler, nicht Latenz.
- <strong>Rate Limit, meine Antwort:</strong> Andere Dimension. Rate Limit z&auml;hlt Requests pro Zeit, &uuml;berschritten hei&szlig;t 429, typisch am Gateway pro Client. Bulkhead z&auml;hlt gleichzeitig laufende Aufrufe pro Downstream. Das Szenario, das Rate Limit nicht abf&auml;ngt: 100 Requests pro Sekunde d&uuml;rfen rein, das Backend braucht pl&ouml;tzlich 5 s pro Call, nach 5 s h&auml;ngen 500 Aufrufe. Rate Limit sieht die Eingangsrate, nicht die Verweildauer. <strong>Spicy:</strong> Wer Rate Limit hat und sich gegen Slow-Downs gewappnet glaubt, verwechselt Durchsatz mit Nebenl&auml;ufigkeit.
- <strong>Non-blocking, meine Antwort:</strong> Non-blocking macht Threads billig, nicht alle Ressourcen. Der Engpass wandert, er verschwindet nicht. Was bleibt: Connection-Pool-Slots im HTTP-Client, Speicher f&uuml;r Puffer und Kontexte, Dateideskriptoren, und das Backend selbst. Rate Limit kann au&szlig;erdem nicht pro Downstream sheddern. Details auf der Bonus-Folie darunter.
- <strong>Warteschlange, meine Antwort:</strong> Eine begrenzte Queue (Resilience4j <code>maxWaitDuration</code>) w&auml;re m&ouml;glich, wir haben uns bewusst dagegen entschieden. Queueing tarnt das Problem, wenn der Pool voll ist, ist das Backend am Limit, die Queue verschiebt das nur und erh&ouml;ht die Latenz. Der sofortige Reject ist ein Backpressure-Signal nach oben, die Queue absorbiert es. <strong>Spicy:</strong> Versteckte Queues sind die Art, wie Systeme langsam und unvorhersehbar werden, ohne dass es im Monitoring auff&auml;llt.
- <strong>Replicas, meine Antwort:</strong> Der client-seitige Bulkhead sch&uuml;tzt den Client, nicht das Backend. Bei f&uuml;nf Replicas mit je zehn Slots kann Hotel f&uuml;nfzig gleichzeitige Calls sehen und wei&szlig; davon nichts. Komplement&auml;r auf der Backend-Seite: Server-side Rate Limiting, ein zentrales Limit im Gateway oder Mesh, Backpressure per <code>429</code> plus <code>Retry-After</code>. <strong>Spicy:</strong> Bulkhead allein ist eine halbe L&ouml;sung. Wer nur den Client sch&uuml;tzt und glaubt, das Backend sei gerettet, hat das Pattern falsch verstanden.
- Die fr&uuml;here Frage &bdquo;Warum 10?&ldquo; ist in die &Uuml;bung gewandert. Vollst&auml;ndige Antworten und Anekdoten: <code>docs/questions/story5.md</code>.
