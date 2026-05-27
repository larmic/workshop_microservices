## Story 3 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Warum gerade 5 Fehler?</h3>
<p>Akzeptanzkriterium sagt &bdquo;nach <span class="hl">5</span> Fehlern&ldquo;. Wo kommt die Zahl her &mdash; w&auml;re 3 sicherer, 20 stabiler?</p>
<code>failures &gt;= 5 &rarr; OPEN</code>
<aside class="notes"><strong>Meine Antwort:</strong> Niemand wei&szlig; es ohne Last- und Backend-Profil. <strong>Zu klein</strong> (1&ndash;2): ein Netz-Glitch oder eine GC-Pause &ouml;ffnet den CB &mdash; <em>false positives</em>. <strong>Zu gro&szlig;</strong> (50): bei 100&nbsp;req/s und totem Backend brennen 50 Fehler-Calls in &lt;&nbsp;1&nbsp;s durch, jeder mit 3&nbsp;s Timeout &mdash; der Booking-Service h&auml;ngt an 50 Threads, bevor der CB &uuml;berhaupt reagiert (Br&uuml;cke zu Story 4 &mdash; Bulkhead). <strong>Alternative:</strong> Resilience4j arbeitet mit <em>Failure-Rate</em> &uuml;ber Sliding Window (z.&nbsp;B. &gt;&nbsp;50&nbsp;% der letzten 100 Calls) &mdash; robuster, aber komplexer und braucht Last.<br><strong>Spicy:</strong> Die 5 ist eine Designentscheidung, kein Naturgesetz. In Produktion gehört der Wert ans aktuelle Last- und Failure-Profil getuned, nicht ins erste Code-Beispiel kopiert.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Probe-Storm</h3>
<p>Im <code>HALF_OPEN</code> lassen wir genau <span class="hl">einen</span> Call durch. Warum nicht alle gleichzeitig?</p>
<code>CompareAndSwap(false, true)</code>
<aside class="notes"><strong>Meine Antwort:</strong> Weil alle wartenden Aufrufer das gerade erholende Backend sofort wieder umlegen w&uuml;rden. Ohne Probe-Lock laufen bei <code>1000&nbsp;req/s</code> nach 30&nbsp;s genau 1000 Calls gleichzeitig los &mdash; klassischer <em>Probe-Storm</em>. Mit atomarem Slot (<code>CompareAndSwap</code>) kommt genau ein Call durch, alle anderen werden als <code>in_flight</code> abgek&uuml;rzt. Erfolg &rarr; CLOSED. Fehler &rarr; weitere 30 s OPEN.<br><strong>Spicy:</strong> Genau hier trennt sich Library-Qualit&auml;t von schnell selbstgeschriebenem Code. Eine naive Implementierung (&bdquo;Wartezeit abgelaufen? Dann CLOSED&ldquo;) tr&auml;gt zum Cascading-Failure bei, statt ihn zu verhindern.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Was z&auml;hlt als Fehler?</h3>
<p>5xx ja, Timeout ja. Aber: ein <code>404</code> &mdash; ist das wirklich ein <span class="hl">Backend-Problem</span>?</p>
<code>4xx &ne; Backend kaputt</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein. Eine 404 hei&szlig;t: <em>der Aufrufer</em> hat eine nicht-existente ID benutzt &mdash; das Backend ist gesund. Grobe Klassifizierung: <strong>5xx, Timeout, Connection-Refused &rarr; CB-relevant</strong>; <strong>4xx (au&szlig;er 429) &rarr; nicht relevant</strong>; <strong>429 &rarr; diskutabel</strong> (Backend &uuml;berlastet). Unser Workshop-Code ist bewusst grob: alles mit <code>status &gt;= 400</code> z&auml;hlt &mdash; in Produktion w&auml;re das ein Bug. Reale Libraries (Resilience4j) machen das &uuml;ber <code>recordExceptions</code> / <code>ignoreExceptions</code>.<br><strong>Spicy:</strong> &bdquo;Fehler&ldquo; ist nicht bin&auml;r. Wer den CB nur an <code>error != nil</code> h&auml;ngt, baut ihm ein zu sensibles Fr&uuml;hwarnsystem.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> L&uuml;gt der Fallback?</h3>
<p>Bei OPEN liefern wir <code>flights: []</code>. Der User sieht <span class="hl">&bdquo;keine Fl&uuml;ge&ldquo;</span> &mdash; t&auml;uschen wir ihn?</p>
<code>X-Circuit-Open: flight</code>
<aside class="notes"><strong>Meine Antwort:</strong> Streng genommen: ja, das ist ein UX-Antipattern. Der User kann nicht zwischen &bdquo;keine Fl&uuml;ge auf dieser Strecke&ldquo; und &bdquo;unser System ist gerade kaputt&ldquo; unterscheiden &mdash; und entscheidet im Zweifel falsch (&bdquo;dann fahre ich halt nur Hotel + Auto&ldquo;). Ehrlichere Varianten: (1) <strong>Teilantwort mit Hinweis-Feld</strong> (<code>unavailable: ["flight"]</code>), (2) <strong>aussagekr&auml;ftige Header</strong> (<code>X-Circuit-Open: flight</code> &mdash; haben wir bereits, UI nutzt's aber nicht), (3) <strong>Stale-while-revalidate</strong> &mdash; letztes erfolgreiches Ergebnis mit Hinweis &bdquo;Daten von vor 2 Min&ldquo;.<br><strong>Spicy:</strong> Resilience-Patterns sind <em>nicht ehrlich von selbst</em>. Sie liefern eine Antwort, die wie Erfolg aussieht &mdash; wer das nicht durch UX und Header transparent macht, baut gut funktionierende Systeme, in denen der Mensch falsche Entscheidungen trifft.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> 30 s warten?</h3>
<p>Backend in 1 s wieder gesund &rarr; 29 s blockiert. Backend 5 min tot &rarr; alle 30 s drauf<span class="hl">h&auml;mmern</span>.</p>
<code>5s &rarr; 10s &rarr; 20s &rarr; 40s &hellip;</code>
<aside class="notes"><strong>Meine Antwort:</strong> 30 s ist ein Kompromiss f&uuml;r den typischen Fall (Container restartet, ist in 10&ndash;20 s wieder oben). Zwei smartere Strategien: (1) <strong>Exponential Backoff</strong>: 5&nbsp;s &rarr; 10&nbsp;s &rarr; 20&nbsp;s &rarr; 40&nbsp;s &rarr; 80&nbsp;s, gedeckelt bei 5 Min. Sobald ein Probe erfolgreich ist, Reset. (2) <strong>Adaptive Wartezeit</strong> am letzten erfolgreichen RTT orientiert &mdash; schnelles Backend darf nach 1 s wieder probiert werden.<br><strong>Spicy:</strong> Eine fixe Wartezeit ist die schlechteste der guten Optionen. Sie reicht f&uuml;r die Demo. Resilience4j, Polly &amp; Co. liefern Backoff out of the box &mdash; sich darauf <em>nicht</em> zu verlassen ist ein Symptom f&uuml;r &bdquo;CB selbst gebaut&ldquo;.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> CB-State nach Restart</h3>
<p>Der Zustand lebt im RAM. Nach <code>docker restart</code> ist alles wieder <code>CLOSED</code> &mdash; auch wenn Hotel <span class="hl">tot</span> ist.</p>
<code>CLOSED &rarr; 5&times; fail &rarr; OPEN</code>
<aside class="notes"><strong>Meine Antwort:</strong> Ja, in der Tendenz gef&auml;hrlich. In den ersten Sekunden nach Restart ist der Service blind und verbraucht 5 Aufrufe &agrave; 3 s Wartezeit, bevor der CB greift &mdash; bei mehreren Replicas multipliziert sich das, bei Rolling Deploy sind alle gleichzeitig in der Lernphase. L&ouml;sungsstufen: (1) hinnehmen, (2) <strong>Initial Probe</strong> beim Service-Start gegen jedes Backend, (3) <strong>Shared State</strong> in Redis o.&Auml;. (Aufwand hoch), (4) <strong>Service Mesh</strong> &mdash; Sidecar h&auml;lt den State, App-Restart hat keinen Einfluss.<br><strong>Spicy:</strong> Lokale CB-Statistik ist per Definition ephemerer Zustand. Wer das vergisst, hat nach jedem Deploy einen Latency-Spike im Monitoring &mdash; und niemand kann erkl&auml;ren, woher.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">7</span> Granularit&auml;t</h3>
<p>Ein CB f&uuml;r alles? Pro Service? Pro <span class="hl">Endpoint</span>?</p>
<code>Flight / Hotel / Car</code>
<aside class="notes"><strong>Meine Antwort:</strong> Drei Stufen mit unterschiedlichen Trade-offs. <strong>Global</strong> (ein CB f&uuml;r alle Backends): trivial, aber ein kranker Service kappt alle. <strong>Pro Service</strong> (unsere Wahl): isoliert Ausf&auml;lle, aber innerhalb des Service kein Schutz. <strong>Pro Endpoint</strong>: sehr feingranular, aber viele CBs zu pflegen. Pro-Service ist die &uuml;bliche Default-Wahl &mdash; ein &bdquo;kranker&ldquo; Service betrifft meist alle Endpoints (Connection-Pool tot, Container down). Pro-Endpoint lohnt sich, wenn ein Service sehr unterschiedliche Workloads hat (schnell <code>search</code> vs. langsam <code>book</code>).<br><strong>Spicy:</strong> Granularit&auml;t ist eine Designentscheidung, keine Pattern-Eigenschaft. Wer einen langsamen <code>book</code>-Endpoint hat und einen schnellen <code>search</code>, fasst die zusammen &mdash; und kappt <code>search</code>, weil <code>book</code> unter Last steht.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">8</span> Fallback &ne; Circuit Breaker</h3>
<p>Erster Fehler bei Flight, Z&auml;hler steht auf <code>1</code>, CB noch <code>CLOSED</code>. Kommt <span class="hl">jetzt schon</span> <code>flights: []</code> &mdash; oder erst bei OPEN?</p>
<code>err != nil &rarr; fallback</code>
<aside class="notes"><strong>Meine Antwort:</strong> Schon beim ersten Fehler. Die leere Liste ist die Reaktion auf <em>jeden</em> fehlgeschlagenen Call, nicht auf den offenen CB &mdash; im Code: <code>if err != nil { markFallback(); return [] }</code>, unabh&auml;ngig vom State. Der CB &auml;ndert nur das <em>Wie</em>: <strong>CLOSED</strong> &rarr; Call geht raus, kostet bis 3&nbsp;s Timeout, <em>dann</em> leere Liste (Header <code>X-Fallback</code>). <strong>OPEN</strong> (ab dem 5. Fehler) &rarr; sofort, ganz ohne Call, zus&auml;tzlich <code>X-Circuit-Open</code>. Der CB &bdquo;produziert&ldquo; die leere Liste also nicht &mdash; er sorgt nur daf&uuml;r, dass sie <em>sofort statt nach Timeout</em> kommt und das tote Backend nicht weiter beschossen wird.<br><strong>Spicy:</strong> Wer Fallback und CB gleichsetzt, sucht den Fehler an der falschen Stelle. Den Fallback baut die Fachlogik (was zeige ich bei Ausfall?), der CB entscheidet nur, ob der Aufruf &uuml;berhaupt rausgeht.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">9</span> CB &uuml;ber mehrere Instanzen?</h3>
<p>Zwei Booking-Replicas. Flight-CB von <span class="hl">A</span> ist OPEN &mdash; soll <span class="hl">B</span> auch dichtmachen?</p>
<code>state lebt im RAM <em>der</em> Instanz</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein &mdash; lokal pro Instanz ist richtig (im Code: <code>sync.Mutex</code> + Felder im Prozess, kein Redis/Consul-KV). Der CB misst, was <em>diese</em> Instanz beobachtet. A kann ein Netz-/Routing-Problem zu Flight haben oder gerade eine kranke Flight-Replica erwischen, w&auml;hrend B Flight problemlos erreicht. Ein geteilter Zustand w&uuml;rde B grundlos mitblockieren &mdash; und damit genau die Isolation zerst&ouml;ren, f&uuml;r die man den CB baut. So machen es auch die Libraries (Resilience4j, Polly) und das Service Mesh (Sidecar pro Pod). <strong>Trade-off:</strong> bei komplett totem Backend muss jede Instanz die 5 Fehler selbst lernen &mdash; bei vielen Replicas eine Lastspitze in der Lernphase (Br&uuml;cke zu Frage 6, CB-State nach Restart).<br><strong>Spicy:</strong> &bdquo;Globaler CB-Zustand&ldquo; klingt nach mehr Kontrolle, opfert aber genau die instanz-lokale Fehler-Isolation, die der Sinn der Sache ist.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">10</span> 5 Fehler = 5 Anfragen</h3>
<p>Erster Fehler bei Flight &mdash; h&auml;mmert der Service in <span class="hl">derselben</span> Anfrage noch 4&times; nach, um die 5 vollzukriegen?</p>
<code>1 Execute pro Service pro Request</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein. Ein <code>GET /booking/offers</code> ruft <code>cb.Execute</code> je Service <em>genau einmal</em> auf, kein Retry-Loop. Bei Flight-Fehler: Z&auml;hler +1, Fallback, fertig. Die 5 Fehler entstehen &uuml;ber <strong>5 eingehende Anfragen</strong> (5 Kunden / 5 Requests) &mdash; der Z&auml;hler ist instanz-weit geteilter Zustand &uuml;ber Requests hinweg (daher der Mutex). Nach dem 5. Fehler &rarr; OPEN, danach kommen alle weiteren Anfragen sofort als Fallback zur&uuml;ck, ohne Flight noch anzufassen. <strong>Wichtig:</strong> Der CB macht selbst <em>keine</em> Retries &mdash; Retry ist ein eigenes Pattern, das man bewusst und vorsichtig kombiniert.<br><strong>Spicy:</strong> Retry und Circuit Breaker naiv zusammenwerfen &mdash; dann z&auml;hlt ein Retry-Sturm den Breaker k&uuml;nstlich schnell hoch, oder die Retries halten das tote Backend unter Dauerfeuer. Erst CB, dann (sparsam) Retry.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story3.md</code>.
</aside>
