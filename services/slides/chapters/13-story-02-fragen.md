## Story 2 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Neuer SPOF?</h3>
<p>Jeder Aufruf fragt zuerst Consul. Was, wenn Consul <span class="hl">kippt</span>?</p>
<code>Consul down &rarr; Booking down</code>
<aside class="notes"><strong>Meine Antwort:</strong> Ja, <strong>wenn man es naiv implementiert</strong> &mdash; und genau so haben wir es gebaut. In Produktion baut man mehrstufig: lokaler Cache der letzten Aufl&ouml;sung, Consul-Agent pro Node (puffert), TTL gegen alte Eintr&auml;ge.<br><strong>Spicy:</strong> Service Discovery l&ouml;st das Problem statischer URLs &mdash; und schafft sich dadurch eine neue zentrale Komponente, die hochverf&uuml;gbar sein muss. Wer Consul / Eureka / etcd &bdquo;mal eben&ldquo; einf&uuml;hrt, ohne &uuml;ber deren Resilienz nachzudenken, hat das Problem nur verschoben.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Consul im K8s-Cluster?</h3>
<p>Kubernetes bringt seine eigene Service Discovery mit. Ist Consul daneben <span class="hl">noch sinnvoll</span>?</p>
<code>kube-dns vs. Consul</code>
<aside class="notes"><strong>Meine Antwort:</strong> Im reinen K8s-Setup: <strong>nein.</strong> Der API-Server ist die Registry, `kube-dns` / CoreDNS l&ouml;st <code>flight-service.default.svc.cluster.local</code> auf, `kube-proxy` macht die Lastverteilung &mdash; alles ohne Zusatzkomponente, ohne Self-Registration im Code. Ein zweites Consul daneben w&auml;re reines Duplikat plus Sync-Problem.<br><strong>Wann doch?</strong> (1) <strong>Hybrid-/Multi-Cluster</strong>: Services laufen teils in K8s, teils auf VMs / Bare Metal / anderem Cluster &mdash; eine cluster-&uuml;bergreifende Registry. (2) <strong>Consul Connect als Service Mesh</strong> (mTLS, L7-Routing, Intentions) &mdash; dann ist Discovery nur Beiwerk, das Mesh ist der Grund. (3) <strong>Legacy-Brownfield</strong> mit existierender Consul-Infrastruktur, bei der K8s schrittweise dazukommt.<br><strong>Spicy:</strong> &bdquo;Wir nehmen Consul, weil wir das immer so machen&ldquo; ist in einem reinen K8s-Setup eine teure Gewohnheit. Die Frage ist nicht &bdquo;Consul ja/nein?&ldquo;, sondern &bdquo;welches Problem habe ich, das K8s nicht schon l&ouml;st?&ldquo; &mdash; meist hei&szlig;t die Antwort dann Mesh, nicht Discovery.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Was passiert beim STOP?</h3>
<p>Beim Start melden sich die Services an. Beim `docker stop` &mdash; verschwindet der Eintrag <span class="hl">sofort</span>?</p>
<code>10&ndash;30 s Traffic ins Leere</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein. Ohne Graceful Shutdown bleibt der Eintrag, bis der Health-Check nach mehreren Misses ausschl&auml;gt &mdash; typisch 10&ndash;30 s, in denen Traffic auf eine tote Instanz l&auml;uft (Connection-Refused). L&ouml;sungen: <strong>Graceful Shutdown</strong> mit aktivem `deregister`, schnellerer Health-Check-Intervall (kostet Last), <strong>Out-of-Service-Mode</strong> (keine neuen Requests, laufende beenden, dann Stop).<br><strong>Spicy:</strong> Service Discovery ist immer <strong>eventually consistent</strong>. Es gibt <strong>immer</strong> ein Zeitfenster, in dem Aufrufer auf tote Endpoints sto&szlig;en &mdash; genau einer der Gr&uuml;nde, warum wir in Story 3 / 4 Resilience-Patterns brauchen.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Registriert &ne; gesund</h3>
<p>Consul nimmt 200 als Beweis. Doch der Check kann <span class="hl">l&uuml;gen</span> &ndash; und selbst ein ehrlicher ist im Moment des Calls l&auml;ngst veraltet. Wer sch&uuml;tzt den Aufrufer?</p>
<code>gepr&uuml;fte URL &ne; laufender Call</code>
<aside class="notes"><strong>Zwei Wege, dasselbe Problem.</strong> (1) <strong>Der Check l&uuml;gt</strong> (inhaltlich falsch): <code>/health</code> sagt 200, der Service ist aber faktisch tot &ndash; Zombie (HTTP lebt, Worker-Pool tot), Backend-Abh&auml;ngigkeit weg (DB nicht erreichbar), oder langsame Degradation (P99 von 50 ms auf 5 s). Gegenmittel: der <em>richtige</em> Check (TCP &rarr; HTTP &rarr; Probe-Logik &rarr; Synthetic). (2) <strong>Der Check ist veraltet</strong> (zeitlich falsch): selbst ein ehrlicher Check ist Sekunden alt; der Service kann genau zwischen zwei Checks ausfallen oder unter Last wegbrechen (siehe STOP-Fenster, Frage 3). Bessere Checks verkleinern Weg 1, Weg 2 bleibt <strong>prinzipiell</strong>. Beide enden gleich: <strong>registriert garantiert nicht gesund im Moment des Calls.</strong><br><strong>&Uuml;berleitung zu Story 3:</strong> Auf Discovery-Ebene ist das nicht l&ouml;sbar &ndash; der <strong>Aufrufer selbst</strong> muss sich sch&uuml;tzen: Circuit Breaker, Timeout, Fallback (Teilbuchung Hotel + Mietwagen ohne Flug).<br><strong>Knackig:</strong> &bdquo;Discovery sagt dir, wo der Service <em>war</em>. Nicht, ob er noch da ist, wenn du anrufst.&ldquo;</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story2.md</code>.
</aside>
