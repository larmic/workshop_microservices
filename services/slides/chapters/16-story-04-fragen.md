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
- Alle vier Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein.
- <strong>Fehler, meine Antwort:</strong> Nein. Eine 404 heisst, der Aufrufer hat eine nicht existente ID benutzt, das Backend ist gesund. Grob: 5xx, Timeout, Connection refused sind CB-relevant. 4xx ausser 429 nicht. 429 ist diskutabel, das Backend ist &uuml;berlastet. Unser Workshop-Code z&auml;hlt alles ab 400, in Produktion w&auml;re das ein Bug. Resilience4j macht das &uuml;ber <code>recordExceptions</code> und <code>ignoreExceptions</code>. <strong>Spicy:</strong> Wer den CB an <code>error != nil</code> h&auml;ngt, baut ein zu sensibles Fr&uuml;hwarnsystem.
- <strong>Nachfrage zu Frage 2, falls jemand &bdquo;f&uuml;nf Fehler in einer Anfrage&ldquo; vermutet:</strong> Nein. Ein <code>GET /booking/offers</code> ruft den Breaker je Service genau einmal auf, kein Retry-Loop. Die f&uuml;nf Fehler entstehen &uuml;ber f&uuml;nf eingehende Anfragen, der Z&auml;hler ist instanzweiter Zustand &uuml;ber Requests hinweg, daher der Mutex. Der CB macht selbst keine Retries, Retry ist ein eigenes Pattern. <strong>Spicy:</strong> Retry und CB naiv zusammenwerfen, dann z&auml;hlt ein Retry-Sturm den Breaker k&uuml;nstlich hoch oder h&auml;lt das tote Backend unter Dauerfeuer. Erst CB, dann sparsam Retry.
- <strong>Fallback, meine Antwort:</strong> Schon beim ersten Fehler. Die leere Liste ist die Reaktion auf jeden fehlgeschlagenen Call, unabh&auml;ngig vom State. Der CB &auml;ndert nur das Wie: CLOSED heisst Call geht raus, kostet bis 3 s, dann leere Liste mit Header <code>X-Fallback</code>. OPEN heisst sofort, ohne Call, mit Header <code>X-Circuit-Open</code>. <strong>Spicy:</strong> <code>flights: []</code> sieht f&uuml;r den User aus wie &bdquo;keine Fl&uuml;ge&ldquo;, nicht wie &bdquo;System kaputt&ldquo;. Der Header macht es transparent, die UI muss ihn auch nutzen.
- <strong>Probe-Storm, meine Antwort:</strong> Weil alle wartenden Aufrufer das gerade erholende Backend sofort wieder umlegen w&uuml;rden. Ohne Probe-Lock laufen bei 1000 Requests pro Sekunde nach 30 s genau 1000 Calls gleichzeitig los. Mit atomarem Slot (<code>CompareAndSwap</code>) kommt genau einer durch, alle anderen bekommen den Fallback. <strong>Spicy:</strong> Genau hier trennt sich Library-Qualit&auml;t von selbstgeschriebenem Code.
- <strong>Replicas, meine Antwort:</strong> Der Zustand ist instanzlokal, und das ist richtig so. A kann ein Routing-Problem zu Flight haben oder eine kranke Flight-Replica erwischen, w&auml;hrend B Flight problemlos erreicht. Geteilter Zustand w&uuml;rde B grundlos mitblockieren. So machen es Resilience4j, Polly und das Service Mesh. Zweite Konsequenz: Nach <code>docker restart</code> ist alles wieder CLOSED, auch wenn Hotel tot ist. Die ersten f&uuml;nf Calls kosten erneut je 3 s. L&ouml;sungsstufen: hinnehmen, Initial Probe beim Start, Shared State, Service Mesh. <strong>Spicy:</strong> Wer das vergisst, hat nach jedem Deploy einen Latenz-Spike im Monitoring, den niemand erkl&auml;ren kann.
- <strong>Reserve, Granularit&auml;t (ein Breaker f&uuml;r alles, pro Service, pro Endpoint?):</strong> Global ist trivial, aber ein kranker Service kappt alle. Pro Service (unsere Wahl) isoliert Ausf&auml;lle. Pro Endpoint ist feingranular, aber viel zu pflegen. Pro Service ist der &uuml;bliche Default, ein kranker Service betrifft meist alle Endpoints. Pro Endpoint lohnt sich bei sehr unterschiedlichen Workloads, etwa schnelles <code>search</code> gegen langsames <code>book</code>. <strong>Spicy:</strong> Wer beide zusammenfasst, kappt <code>search</code>, weil <code>book</code> unter Last steht.
- Vollst&auml;ndige Antworten und Anekdoten: <code>docs/questions/story4.md</code>.
