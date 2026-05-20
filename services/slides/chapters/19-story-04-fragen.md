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
<h3><span class="numeral">2</span> Reject-Kaskade</h3>
<p>20er-Burst bei langsamen Backends &mdash; Flight rejected <span class="hl">10</span>, Hotel <span class="hl">7</span>, Car <span class="hl">1</span>. Warum sinkend?</p>
<code>Traffic-Shaper</code>
<aside class="notes"><strong>Meine Antwort:</strong> Sequenzielle Aufruf-Reihenfolge plus Timing-Jitter. Flight sieht den Burst in voller Wucht (20 zeitgleich &rarr; 10 rejected). Die abgewiesenen 10 marschieren sofort zu Hotel, die erfolgreichen 10 kommen 50&ndash;200&nbsp;ms sp&auml;ter nach &mdash; Hotel sieht den Burst zeitlich entzerrt, Reject-Rate sinkt. Car sieht nur noch die geringe Restwelle.<br><strong>Spicy:</strong> Eine Bulkhead-Kette ist gleichzeitig ein <em>Traffic-Shaper</em>. Jede Stufe sch&uuml;tzt nicht nur sich selbst, sondern gl&auml;ttet auch die Last f&uuml;r alle nachgelagerten &mdash; ein kostenloser Nebeneffekt, der in Produktion oft mehr Wirkung zeigt als das Pattern selbst. Wer nur den letzten Hop instrumentiert und sich freut, dass dort fast nichts abgelehnt wird, hat den Punkt verpasst.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Reject &ne; Abbruch</h3>
<p>Flight-Bulkhead voll &rarr; Hotel und Car laufen <span class="hl">trotzdem weiter</span>. Pattern oder Code-Entscheidung?</p>
<code>Mechanismus &ne; Strategie</code>
<aside class="notes"><strong>Meine Antwort:</strong> Code-Entscheidung. <strong>Bulkhead-Pattern</strong>: &bdquo;Pool voll &rarr; sofort ablehnen.&ldquo; Was der Aufrufer mit dem Reject macht, ist nicht Teil des Patterns. <strong>Unser Aggregator</strong>: f&auml;hrt nach jedem Fehler (CB-Open, Bulkhead-Full, echter Backend-Fehler) mit leerem Teilergebnis fort &mdash; <em>Fail-Soft</em> beim GET. Im POST-Pfad (<code>/booking/bookings</code>) machen wir dagegen <em>Fail-Fast</em>: ein Reject = ganze Buchung 503, weil eine halbe Buchung schlimmer ist als gar keine.<br><strong>Spicy:</strong> Resilience-Patterns liefern einen <em>Mechanismus</em>. Die <em>Strategie</em> dahinter (Teilergebnis vs. Komplettabbruch vs. Retry mit anderem Backend) ist Fachlogik &mdash; und h&auml;ngt davon ab, ob die Operation idempotent ist und ob Teilergebnisse sinnvoll sind.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> CB unbeeinflusst</h3>
<p>Ein Bulkhead-Reject z&auml;hlt <span class="hl">nicht</span> als CB-Failure. Warum nicht?</p>
<code>bh &rarr; cb &rarr; call</code>
<aside class="notes"><strong>Meine Antwort:</strong> Weil ein Bulkhead-Reject sagt &bdquo;ich sch&uuml;tze mich selbst&ldquo;, nicht &bdquo;das Backend ist krank&ldquo;. Im Code-Aufbau steht der Bulkhead <em>vor</em> dem CB: <code>bh.Execute(...cb.Execute(...httpCall))</code>. Ist der Pool voll, gibt der Bulkhead sofort <code>ErrBulkheadFull</code> zur&uuml;ck &mdash; der CB wird gar nicht erst aufgerufen, seine Z&auml;hler bleiben sauber. W&uuml;rden wir Rejects als Failures werten, w&uuml;rde der CB f&auml;lschlich auf OPEN gehen, obwohl das Backend gesund antwortet &mdash; selbst eingebuddelte T&uuml;r.<br><strong>Spicy:</strong> CB und Bulkhead sind komplement&auml;r, nicht alternativ. Sie d&uuml;rfen gleichzeitig feuern &mdash; das ist gut so, weil sie auf verschiedene Symptome reagieren.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Warum keine Queue?</h3>
<p>Eine kurze Warteschlange w&auml;re doch <span class="hl">freundlicher</span> als sofort Reject?</p>
<code>Queue tarnt Last</code>
<aside class="notes"><strong>Meine Antwort:</strong> Eine begrenzte Queue (Resilience4j: <code>maxWaitDuration</code>) w&auml;re m&ouml;glich &mdash; wir haben bewusst dagegen entschieden. Zwei Gr&uuml;nde: (1) <strong>Queueing tarnt das Problem</strong>. Wenn der Bulkhead voll ist, ist das Backend bereits am Limit. Queue verschiebt das nur in die Zukunft und erh&ouml;ht End-to-End-Latency. (2) <strong>Backpressure-Signal nach oben</strong>. Sofortiger Reject sagt dem Aufrufer: &bdquo;lass mich in Ruhe, ich bin voll.&ldquo; Queue absorbiert das Signal.<br><strong>Spicy:</strong> Queue vs. kein Queue ist eine bewusste Designentscheidung. Default sollte <em>kein Queue</em> sein &mdash; eine versteckte Queue ist die Art, wie Systeme langsam und unvorhersehbar werden, ohne dass es im Monitoring auff&auml;llt.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Warum 10?</h3>
<p>Wir haben <code>maxConcurrent=10</code> hartkodiert. Wo kommt die Zahl her &mdash; <span class="hl">w&auml;re 5 sicherer, 100 freundlicher</span>?</p>
<code>L = &lambda; &times; W</code>
<aside class="notes"><strong>Meine Antwort:</strong> 10 ist eine Workshop-Default-Zahl, die fürs Demo gut funktioniert. In Produktion orientiert man die Zahl an drei Faktoren: <strong>(1) Connection-Pool</strong> des HTTP-Clients (mehr Slots als verf&uuml;gbare Connections sind sinnlos), <strong>(2) Backend-Kapazit&auml;t</strong> (Faustregel: <em>Replicas &times; maxConcurrent &le; Backend-Kapazit&auml;t</em>), <strong>(3) Little's Law</strong>: <code>concurrency = throughput &times; latency</code>. Bei 50&nbsp;ms RTT und 200&nbsp;req/s Ziel &rarr; 10. Bei 500&nbsp;ms RTT und gleichem Throughput &rarr; 100. <strong>Kontraintuitiv:</strong> langsamere Backends brauchen <em>mehr</em> Slots, nicht weniger.<br><strong>Spicy:</strong> Wer <code>maxConcurrent</code> aus dem Bauch heraus setzt (&bdquo;10 klingt gut&ldquo;), hat den Bulkhead nicht implementiert, sondern dekoriert. Das Limit geh&ouml;rt aus gemessener Backend-Kapazit&auml;t abgeleitet, nicht aus Beispiel-Code &uuml;bernommen.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">7</span> Wer sch&uuml;tzt das Backend?</h3>
<p>5 Replicas &times; 10 = <span class="hl">50 parallele Calls</span> bei Hotel. Bulkhead sch&uuml;tzt wen?</p>
<code>Inbound &ne; Outbound</code>
<aside class="notes"><strong>Meine Antwort:</strong> Client-side Bulkhead sch&uuml;tzt <em>den Client</em>, nicht das Backend. Bei 5 Booking-Replicas &agrave; <code>maxConcurrent=10</code> kann Hotel bis zu 50 gleichzeitige Calls sehen &mdash; und Hotel wei&szlig; davon nichts. Komplement&auml;re Patterns auf der <em>Backend</em>-Seite: (1) <strong>Server-side Rate Limiting</strong> (Hotel lehnt nach N Calls ab &mdash; sch&uuml;tzt vor jedem Aufrufer), (2) <strong>API-Gateway / Service Mesh</strong> (zentraler Punkt &uuml;ber alle Aufrufer hinweg), (3) <strong>Backpressure</strong>: <code>429 + Retry-After</code> &mdash; Hotel kommuniziert &Uuml;berlast aktiv.<br><strong>Spicy:</strong> Bulkhead allein ist eine halbierte L&ouml;sung. Sie macht den eigenen Service stabil &mdash; aber das gesch&uuml;tzte Backend braucht erg&auml;nzende Mechanismen. Wer nur den Client sch&uuml;tzt und glaubt, das Backend sei auch gerettet, hat das Pattern falsch verstanden.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story4.md</code>.
</aside>
