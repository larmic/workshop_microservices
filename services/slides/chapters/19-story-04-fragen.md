## Story 4 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Wozu, wenn es CB gibt?</h3>
<p>Wir haben in Story 3 schon einen Circuit Breaker. Was bringt Bulkhead, was der CB <span class="hl">nicht</span> kann?</p>
<code>slow &ne; failure</code>
<aside class="notes"><strong>Meine Antwort:</strong> Beide Patterns l&ouml;sen verschiedene Probleme. <strong>CB</strong>: &bdquo;Backend ist krank, ich versuch's gar nicht mehr&ldquo; &mdash; Reaktion auf Fehler-Rate. <strong>Bulkhead</strong>: &bdquo;Ich verbrenne maximal N Slots f&uuml;r dieses Backend, egal wie viel Last reinkommt&ldquo; &mdash; Reaktion auf Ressourcen-Druck. Killer-Szenario, das nur Bulkhead l&ouml;st: Hotel antwortet in 2&nbsp;s fehlerfrei &mdash; CB bleibt CLOSED, aber unter Last laufen alle Threads in Hotel-Calls auf, Flight und Car kommen nicht durch. Bulkhead kappt nach 10 parallelen Hotel-Calls und l&auml;sst die anderen ungest&ouml;rt.<br><strong>Spicy:</strong> &bdquo;Wir haben doch schon einen CB&ldquo; ist die Standard-Falle &mdash; CB sieht Fehler, nicht Latenz. Wer beides nicht trennt, baut sich falsche Sicherheit.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Rate Limit vs. Bulkhead</h3>
<p><strong>Rate Limit:</strong> &bdquo;Sch&uuml;tze mich vor zu vielen Anfragen von <span class="hl">au&szlig;en</span>.&ldquo;</p>
<p><strong>Bulkhead:</strong> &bdquo;Sch&uuml;tze mich davor, dass ein langsames anderes System mich <span class="hl">mitrei&szlig;t</span>.&ldquo;</p>
<code>rate &ne; concurrency</code>
<aside class="notes"><strong>Meine Antwort:</strong> Andere Dimension, andere Wirkung. <strong>Rate Limit</strong>: Requests pro Zeit (z.&nbsp;B. 100&nbsp;req/s). &Uuml;berschritten &rarr; sofort <code>429</code>. Typische Stelle: API-Gateway / Edge, pro Client oder API-Key. <strong>Bulkhead</strong>: gleichzeitig in flight (z.&nbsp;B. 10 parallele Calls). &Uuml;berschritten &rarr; sofortiger Reject im aufrufenden Service, pro Downstream. Killer-Szenario, das Rate Limit <em>nicht</em> abf&auml;ngt: 100&nbsp;req/s d&uuml;rfen rein, das Backend braucht pl&ouml;tzlich 5&nbsp;s pro Call &rarr; nach 5&nbsp;s 500 h&auml;ngende Goroutines, OOM. Rate Limit sieht nur Eingangsrate, nicht Verweildauer. Bulkhead kappt bei 10 in flight, egal wie langsam.<br><strong>Spicy:</strong> Wer nur Rate Limit hat und glaubt, gegen Slow-Downs gewappnet zu sein, verwechselt Throughput mit Concurrency: ein klassisches Architektur-Missverst&auml;ndnis. Beide Patterns sind komplement&auml;r, nicht alternativ.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Bulkhead non-blocking?</h3>
<p>Mein Server ist non-blocking, Threads sind bei I/O <span class="hl">frei</span>. Reicht da nicht Rate Limit?</p>
<code>Engpass wandert</code>
<aside class="notes"><strong>Meine Antwort:</strong> Non-blocking macht <em>Threads</em> billig, nicht <em>alle Ressourcen</em>. Der Engpass wandert nur, er verschwindet nicht. Was bei jedem in-flight Call belegt bleibt, egal ob non-blocking: (1) <strong>Connection-Pool-Slot</strong> im HTTP-Client (Go: <code>http.Transport.MaxConnsPerHost</code>, Reactor-Netty: 500 default). Ist Hotel langsam und belegt alle Pool-Slots, kommen Flight und Car nicht mehr durch, obwohl dein Server &bdquo;Threads frei&ldquo; meldet. (2) <strong>Memory</strong> f&uuml;r Request- und Response-Buffer, Context, Cancellation. 50.000 in-flight Calls = realer Heap-Druck, GC-Pause, irgendwann OOM. (3) <strong>FDs/Sockets</strong>, OS-Limit pro Prozess. Little's Law gilt unabh&auml;ngig vom Threading-Modell: <code>L = &lambda; &times; W</code>. Steigt W (Latenz), steigt L (gleichzeitig offene Calls) linear an, bei gleicher Eingangsrate. Rate Limit drosselt &lambda;, nicht L.<br>Wichtiger noch: <strong>Rate Limit kann nicht selektiv pro Downstream sheddern</strong>. Wenn Hotel h&auml;ngt, m&uuml;sstest du mit Rate Limit allein <em>alle</em> Requests drosseln, auch die, die Hotel gar nicht brauchen. Bulkhead pro Downstream sagt: &bdquo;Hotel ist dicht, aber Flight-Only-Anfragen gehen weiter durch.&ldquo;<br><strong>Spicy:</strong> &bdquo;Non-blocking spart mir Bulkhead&ldquo; ist Wunschdenken. Threads sind nur ein Engpass von vielen. In reaktiven Stacks (Spring WebFlux, Reactor) sind Bulkheads nicht zuf&auml;llig weiter empfohlen, Resilience4j hat extra einen <code>SemaphoreBulkhead</code> daf&uuml;r. Wer Bulkhead weglie&szlig;e, m&uuml;sste eine andere Form von per-Downstream-Concurrency-Limit haben, ein zentrales Rate Limit reicht nicht.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Warum keine Queue?</h3>
<p>Eine kurze Warteschlange w&auml;re doch <span class="hl">freundlicher</span> als sofort Reject?</p>
<code>Queue tarnt Last</code>
<aside class="notes"><strong>Meine Antwort:</strong> Eine begrenzte Queue (Resilience4j: <code>maxWaitDuration</code>) w&auml;re m&ouml;glich &mdash; wir haben bewusst dagegen entschieden. Zwei Gr&uuml;nde: (1) <strong>Queueing tarnt das Problem</strong>. Wenn der Bulkhead voll ist, ist das Backend bereits am Limit. Queue verschiebt das nur in die Zukunft und erh&ouml;ht End-to-End-Latency. (2) <strong>Backpressure-Signal nach oben</strong>. Sofortiger Reject sagt dem Aufrufer: &bdquo;lass mich in Ruhe, ich bin voll.&ldquo; Queue absorbiert das Signal.<br><strong>Spicy:</strong> Queue vs. kein Queue ist eine bewusste Designentscheidung. Default sollte <em>kein Queue</em> sein &mdash; eine versteckte Queue ist die Art, wie Systeme langsam und unvorhersehbar werden, ohne dass es im Monitoring auff&auml;llt.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Warum 10?</h3>
<p>Wir haben <code>maxConcurrent=10</code> hartkodiert. Wo kommt die Zahl her &mdash; <span class="hl">w&auml;re 5 sicherer, 100 freundlicher</span>?</p>
<code>L = &lambda; &times; W</code>
<aside class="notes"><strong>Meine Antwort:</strong> 10 ist eine Workshop-Default-Zahl, die fürs Demo gut funktioniert. In Produktion orientiert man die Zahl an drei Faktoren: <strong>(1) Connection-Pool</strong> des HTTP-Clients (mehr Slots als verf&uuml;gbare Connections sind sinnlos), <strong>(2) Backend-Kapazit&auml;t</strong> (Faustregel: <em>Replicas &times; maxConcurrent &le; Backend-Kapazit&auml;t</em>), <strong>(3) Little's Law</strong>: <code>concurrency = throughput &times; latency</code>. Bei 50&nbsp;ms RTT und 200&nbsp;req/s Ziel &rarr; 10. Bei 500&nbsp;ms RTT und gleichem Throughput &rarr; 100. <strong>Kontraintuitiv:</strong> langsamere Backends brauchen <em>mehr</em> Slots, nicht weniger.<br><strong>Spicy:</strong> Wer <code>maxConcurrent</code> aus dem Bauch heraus setzt (&bdquo;10 klingt gut&ldquo;), hat den Bulkhead nicht implementiert, sondern dekoriert. Das Limit geh&ouml;rt aus gemessener Backend-Kapazit&auml;t abgeleitet, nicht aus Beispiel-Code &uuml;bernommen.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Wer sch&uuml;tzt das Backend?</h3>
<p>5 Replicas &times; 10 = <span class="hl">50 parallele Calls</span> bei Hotel. Bulkhead sch&uuml;tzt wen?</p>
<code>Inbound &ne; Outbound</code>
<aside class="notes"><strong>Meine Antwort:</strong> Client-side Bulkhead sch&uuml;tzt <em>den Client</em>, nicht das Backend. Bei 5 Booking-Replicas &agrave; <code>maxConcurrent=10</code> kann Hotel bis zu 50 gleichzeitige Calls sehen &mdash; und Hotel wei&szlig; davon nichts. Komplement&auml;re Patterns auf der <em>Backend</em>-Seite: (1) <strong>Server-side Rate Limiting</strong> (Hotel lehnt nach N Calls ab &mdash; sch&uuml;tzt vor jedem Aufrufer), (2) <strong>API-Gateway / Service Mesh</strong> (zentraler Punkt &uuml;ber alle Aufrufer hinweg), (3) <strong>Backpressure</strong>: <code>429 + Retry-After</code> &mdash; Hotel kommuniziert &Uuml;berlast aktiv.<br><strong>Spicy:</strong> Bulkhead allein ist eine halbierte L&ouml;sung. Sie macht den eigenen Service stabil &mdash; aber das gesch&uuml;tzte Backend braucht erg&auml;nzende Mechanismen. Wer nur den Client sch&uuml;tzt und glaubt, das Backend sei auch gerettet, hat das Pattern falsch verstanden.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story4.md</code>.
</aside>
