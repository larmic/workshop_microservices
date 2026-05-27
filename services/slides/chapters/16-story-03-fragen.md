## Story 3 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Was z&auml;hlt als Fehler?</h3>
<p>5xx ja, Timeout ja. Aber: ein <code>404</code> &mdash; ist das wirklich ein <span class="hl">Backend-Problem</span>?</p>
<code>4xx &ne; Backend kaputt</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein. Eine 404 hei&szlig;t: <em>der Aufrufer</em> hat eine nicht-existente ID benutzt &mdash; das Backend ist gesund. Grobe Klassifizierung: <strong>5xx, Timeout, Connection-Refused &rarr; CB-relevant</strong>; <strong>4xx (au&szlig;er 429) &rarr; nicht relevant</strong>; <strong>429 &rarr; diskutabel</strong> (Backend &uuml;berlastet). Unser Workshop-Code ist bewusst grob: alles mit <code>status &gt;= 400</code> z&auml;hlt &mdash; in Produktion w&auml;re das ein Bug. Reale Libraries (Resilience4j) machen das &uuml;ber <code>recordExceptions</code> / <code>ignoreExceptions</code>.<br><strong>Spicy:</strong> &bdquo;Fehler&ldquo; ist nicht bin&auml;r. Wer den CB nur an <code>error != nil</code> h&auml;ngt, baut ihm ein zu sensibles Fr&uuml;hwarnsystem.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> 5 Fehler = 5 Anfragen</h3>
<p>Erster Fehler bei Flight &mdash; h&auml;mmert der Service in <span class="hl">derselben</span> Anfrage noch 4&times; nach, um die 5 vollzukriegen?</p>
<code>1 Execute pro Service pro Request</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein. Ein <code>GET /booking/offers</code> ruft <code>cb.Execute</code> je Service <em>genau einmal</em> auf, kein Retry-Loop. Bei Flight-Fehler: Z&auml;hler +1, Fallback, fertig. Die 5 Fehler entstehen &uuml;ber <strong>5 eingehende Anfragen</strong> (5 Kunden / 5 Requests) &mdash; der Z&auml;hler ist instanz-weit geteilter Zustand &uuml;ber Requests hinweg (daher der Mutex). Nach dem 5. Fehler &rarr; OPEN, danach kommen alle weiteren Anfragen sofort als Fallback zur&uuml;ck, ohne Flight noch anzufassen. <strong>Wichtig:</strong> Der CB macht selbst <em>keine</em> Retries &mdash; Retry ist ein eigenes Pattern, das man bewusst und vorsichtig kombiniert.<br><strong>Spicy:</strong> Retry und Circuit Breaker naiv zusammenwerfen &mdash; dann z&auml;hlt ein Retry-Sturm den Breaker k&uuml;nstlich schnell hoch, oder die Retries halten das tote Backend unter Dauerfeuer. Erst CB, dann (sparsam) Retry.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Fallback &ne; Circuit Breaker</h3>
<p>Erster Fehler bei Flight, Z&auml;hler steht auf <code>1</code>, CB noch <code>CLOSED</code>. Kommt <span class="hl">jetzt schon</span> <code>flights: []</code> &mdash; oder erst bei OPEN?</p>
<code>err != nil &rarr; fallback</code>
<aside class="notes"><strong>Meine Antwort:</strong> Schon beim ersten Fehler. Die leere Liste ist die Reaktion auf <em>jeden</em> fehlgeschlagenen Call, nicht auf den offenen CB &mdash; im Code: <code>if err != nil { markFallback(); return [] }</code>, unabh&auml;ngig vom State. Der CB &auml;ndert nur das <em>Wie</em>: <strong>CLOSED</strong> &rarr; Call geht raus, kostet bis 3&nbsp;s Timeout, <em>dann</em> leere Liste (Header <code>X-Fallback</code>). <strong>OPEN</strong> (ab dem 5. Fehler) &rarr; sofort, ganz ohne Call, zus&auml;tzlich <code>X-Circuit-Open</code>. Der CB &bdquo;produziert&ldquo; die leere Liste also nicht &mdash; er sorgt nur daf&uuml;r, dass sie <em>sofort statt nach Timeout</em> kommt und das tote Backend nicht weiter beschossen wird.<br><strong>Spicy:</strong> Wer Fallback und CB gleichsetzt, sucht den Fehler an der falschen Stelle. Den Fallback baut die Fachlogik (was zeige ich bei Ausfall?), der CB entscheidet nur, ob der Aufruf &uuml;berhaupt rausgeht. Nebenbei ein UX-Antipattern: <code>flights: []</code> sieht f&uuml;r den User aus wie &bdquo;keine Fl&uuml;ge&ldquo;, nicht wie &bdquo;System kaputt&ldquo; &mdash; der Header <code>X-Circuit-Open</code> macht es transparent, die UI muss ihn aber auch nutzen.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Probe-Storm</h3>
<p>Im <code>HALF_OPEN</code> lassen wir genau <span class="hl">einen</span> Call durch. Warum nicht alle gleichzeitig?</p>
<code>CompareAndSwap(false, true)</code>
<aside class="notes"><strong>Meine Antwort:</strong> Weil alle wartenden Aufrufer das gerade erholende Backend sofort wieder umlegen w&uuml;rden. Ohne Probe-Lock laufen bei <code>1000&nbsp;req/s</code> nach 30&nbsp;s genau 1000 Calls gleichzeitig los &mdash; klassischer <em>Probe-Storm</em>. Mit atomarem Slot (<code>CompareAndSwap</code>) kommt genau ein Call durch, alle anderen werden als <code>in_flight</code> abgek&uuml;rzt. Erfolg &rarr; CLOSED. Fehler &rarr; weitere 30 s OPEN.<br><strong>Spicy:</strong> Genau hier trennt sich Library-Qualit&auml;t von schnell selbstgeschriebenem Code. Eine naive Implementierung (&bdquo;Wartezeit abgelaufen? Dann CLOSED&ldquo;) tr&auml;gt zum Cascading-Failure bei, statt ihn zu verhindern.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Fl&uuml;chtiger, lokaler Zustand</h3>
<p>Der State lebt im RAM <em>einer</em> Instanz. Zwei Replicas: Flight-CB von <span class="hl">A</span> ist OPEN &mdash; soll <span class="hl">B</span> auch dichtmachen? Und was gilt nach <code>docker restart</code>?</p>
<code>kein shared state &middot; Restart &rarr; CLOSED</code>
<aside class="notes"><strong>Meine Antwort:</strong> Der Zustand ist instanz-lokaler RAM (im Code: <code>sync.Mutex</code> + Felder im Prozess, kein Redis/Consul-KV) &mdash; mit zwei Konsequenzen. <strong>(1) Nicht geteilt:</strong> Der CB misst, was <em>diese</em> Instanz beobachtet. A kann ein Netz-/Routing-Problem zu Flight haben oder gerade eine kranke Flight-Replica erwischen, w&auml;hrend B Flight problemlos erreicht. Geteilter Zustand w&uuml;rde B grundlos mitblockieren und genau die Isolation zerst&ouml;ren, f&uuml;r die man den CB baut. So machen es auch Resilience4j, Polly und das Service Mesh (Sidecar pro Pod). <strong>(2) Fl&uuml;chtig:</strong> Nach <code>docker restart</code> ist alles wieder CLOSED &mdash; auch wenn Hotel tot ist. In den ersten Sekunden verbrennt der Service erneut 5 Calls &agrave; 3&nbsp;s, bevor der CB greift; bei Rolling Deploy sind alle Replicas gleichzeitig in dieser Lernphase. L&ouml;sungsstufen: hinnehmen &rarr; Initial Probe beim Start &rarr; Shared State (Redis) &rarr; Service Mesh.<br><strong>Spicy:</strong> &bdquo;Globaler CB-Zustand&ldquo; klingt nach Kontrolle, opfert aber die instanz-lokale Isolation. Und lokale Statistik ist per Definition ephemer: Wer das vergisst, hat nach jedem Deploy einen Latency-Spike im Monitoring &mdash; und niemand kann erkl&auml;ren, woher.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Granularit&auml;t</h3>
<p>Ein CB f&uuml;r alles? Pro Service? Pro <span class="hl">Endpoint</span>?</p>
<code>Flight / Hotel / Car</code>
<aside class="notes"><strong>Meine Antwort:</strong> Drei Stufen mit unterschiedlichen Trade-offs. <strong>Global</strong> (ein CB f&uuml;r alle Backends): trivial, aber ein kranker Service kappt alle. <strong>Pro Service</strong> (unsere Wahl): isoliert Ausf&auml;lle, aber innerhalb des Service kein Schutz. <strong>Pro Endpoint</strong>: sehr feingranular, aber viele CBs zu pflegen. Pro-Service ist die &uuml;bliche Default-Wahl &mdash; ein &bdquo;kranker&ldquo; Service betrifft meist alle Endpoints (Connection-Pool tot, Container down). Pro-Endpoint lohnt sich, wenn ein Service sehr unterschiedliche Workloads hat (schnell <code>search</code> vs. langsam <code>book</code>).<br><strong>Spicy:</strong> Granularit&auml;t ist eine Designentscheidung, keine Pattern-Eigenschaft. Wer einen langsamen <code>book</code>-Endpoint hat und einen schnellen <code>search</code>, fasst die zusammen &mdash; und kappt <code>search</code>, weil <code>book</code> unter Last steht.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story3.md</code>.
</aside>
