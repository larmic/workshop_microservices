<div class="page">

<p class="kicker">Story 3 &middot; Recap</p>

## Vier Fragen an euch

<div class="page-body">

<div class="numlist recap compact">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Neuer Single Point of Failure?</p>
<h3>Jeder Aufruf fragt zuerst Consul. Was, wenn Consul <span class="hl">kippt</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Consul down &rarr; Booking down</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Consul im Kubernetes-Cluster?</p>
<h3>Kubernetes hat eigene Discovery. Braucht es Consul <span class="hl">daneben</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">CoreDNS vs. Consul</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Was passiert beim Stop?</p>
<h3>Beim <code>docker stop</code>: Verschwindet der Eintrag <span class="hl">sofort</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">10 &ndash; 30 s ins Leere</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Registriert &ne; gesund</p>
<h3>Der Check sagt 200. Ist der Service im Moment des Calls noch <span class="hl">da</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">war &ne; ist</code>
</div>
</div>

</div>

</div>

Note:
- Alle vier Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein.
- <strong>SPOF, meine Antwort:</strong> Ja, wenn man es naiv implementiert, und genau so haben wir es gebaut. In Produktion baut man mehrstufig: lokaler Cache der letzten Aufl&ouml;sung, Consul-Agent pro Node (puffert), TTL gegen alte Eintr&auml;ge. <strong>Spicy:</strong> Service Discovery l&ouml;st das Problem statischer URLs und schafft sich dadurch eine neue zentrale Komponente, die hochverf&uuml;gbar sein muss. Wer Consul, Eureka oder etcd &bdquo;mal eben&ldquo; einf&uuml;hrt, ohne &uuml;ber deren Resilienz nachzudenken, hat das Problem nur verschoben.
- <strong>Kubernetes, meine Antwort:</strong> Im reinen Kubernetes-Setup nein. Der API-Server ist die Registry, CoreDNS l&ouml;st <code>flight-service.default.svc.cluster.local</code> auf, kube-proxy verteilt die Last, alles ohne Zusatzkomponente und ohne Self-Registration im Code. Ein zweites Consul daneben w&auml;re Duplikat plus Sync-Problem. <strong>Wann doch?</strong> Hybrid- oder Multi-Cluster (Services teils in Kubernetes, teils auf VMs), Consul Connect als Service Mesh (dann ist Discovery Beiwerk), Brownfield mit bestehender Consul-Infrastruktur. <strong>Spicy:</strong> Die Frage ist nicht &bdquo;Consul ja oder nein&ldquo;, sondern &bdquo;welches Problem habe ich, das Kubernetes nicht schon l&ouml;st?&ldquo; Meist hei&szlig;t die Antwort Mesh, nicht Discovery.
- <strong>Stop, meine Antwort:</strong> Nein. Ohne Graceful Shutdown bleibt der Eintrag, bis der Health-Check nach mehreren Fehlversuchen ausschl&auml;gt, typisch 10 bis 30 Sekunden, in denen Traffic auf eine tote Instanz l&auml;uft (Connection refused). L&ouml;sungen: Graceful Shutdown mit aktivem Deregister, k&uuml;rzeres Check-Intervall (kostet Last), Out-of-Service-Modus (keine neuen Requests, laufende beenden, dann Stop). <strong>Spicy:</strong> Service Discovery ist immer eventually consistent. Es gibt immer ein Zeitfenster mit toten Endpoints, einer der Gr&uuml;nde f&uuml;r die Resilience-Patterns in Story 4 und 5.
- <strong>Registriert ungleich gesund, meine Antwort:</strong> Zwei Wege, dasselbe Problem. Der Check l&uuml;gt inhaltlich: <code>/health</code> sagt 200, der Service ist faktisch tot (Zombie, Worker-Pool tot, DB nicht erreichbar, langsame Degradation). Gegenmittel ist der richtige Check: TCP, HTTP, Probe-Logik, Synthetic. Oder der Check ist zeitlich veraltet: selbst ein ehrlicher Check ist Sekunden alt, der Service kann genau zwischen zwei Checks ausfallen. Bessere Checks verkleinern Weg eins, Weg zwei bleibt prinzipiell. <strong>&Uuml;berleitung zu Story 4:</strong> Auf Discovery-Ebene nicht l&ouml;sbar, der Aufrufer muss sich selbst sch&uuml;tzen: Circuit Breaker, Timeout, Fallback (Teilbuchung Hotel plus Mietwagen ohne Flug). <strong>Knackig:</strong> &bdquo;Discovery sagt dir, wo der Service <em>war</em>. Nicht, ob er noch da ist, wenn du anrufst.&ldquo;
- Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story3.md</code>.
