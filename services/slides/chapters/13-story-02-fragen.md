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
<h3><span class="numeral">2</span> Twelve-Factor-Bruch?</h3>
<p>Der Service kennt jetzt seine eigene IP und meldet sich aktiv an. Ist das eine <span class="hl">Regression</span>?</p>
<code>App kennt Plattform</code>
<aside class="notes"><strong>Meine Antwort:</strong> Halb. Zwei Modelle: <strong>Self-Registration im Code</strong> (Workshop, funktioniert &uuml;berall, aber Service kennt seine Umgebung) vs. <strong>Sidecar / Plattform</strong> (K8s Service, Consul-Connect &mdash; umgebungsneutral, setzt Plattform voraus).<br><strong>Spicy:</strong> Self-Registration im Code ist die schnellere L&ouml;sung, aber sie verteilt Plattform-Wissen in jeden Service. In Ordnung f&uuml;r 5 Services, H&ouml;lle bei 50. Branchen-Trend geht klar zu &bdquo;Plattform &uuml;bernimmt das&ldquo;.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Client-Side LB <em>vs.</em> Service Mesh?</h3>
<p>Falsche Dichotomie &mdash; Mesh <span class="hl">macht</span> Client-Side LB. Brauchen wir beides?</p>
<code>Discovery &rarr; LB-Sidecar &rarr; Mesh</code>
<aside class="notes"><strong>Meine Antwort:</strong> Zwei Achsen trennen: <strong>wer entscheidet</strong> (Client-Side vs. Server-Side) und <strong>wo l&auml;uft der Code</strong> (Library / Sidecar / Plattform). Client-Side LB gibt es als Library (Workshop), als reinen Envoy-Sidecar (ohne Mesh) und implizit in Plattformen (`kube-proxy`). Das Sprach-Argument gegen Library ist mit einem reinen LB-Sidecar bereits erledigt &mdash; daf&uuml;r braucht's noch kein Mesh.<br><strong>Was Mesh wirklich differenziert:</strong> L7-Resilience (Retry, CB, Timeout, Outlier Detection), mTLS by default, zentrale Control Plane (YAML-&Auml;nderung, alle Sidecars ziehen nach), Traffic-Splitting (Canary), einheitliche Observability pro Hop.<br><strong>Brauchen wir beides?</strong> Keine Oder-Frage &mdash; <strong>Mesh enth&auml;lt Service Discovery</strong> (in K8s der API-Server, mit Consul Connect Consul selbst, sonst xDS). Drei Stufen: <strong>1.</strong> nur Discovery (Workshop, App-Team), <strong>2.</strong> Discovery + LB-Sidecar (oft &uuml;bersehener Mittelweg), <strong>3.</strong> Mesh (Plattform-Team, Resilience + mTLS + Routing + Observability).<br><strong>Spicy:</strong> Echte Frage ist nicht &bdquo;Mesh ja/nein&ldquo;, sondern: <em>brauchen wir mTLS, Traffic-Splitting und L7-Resilience &uuml;berall, oder reicht Stufe 1 oder 2?</em> Wer direkt zu Istio greift, hat sechs Monate sp&auml;ter Sidecar-Wildwuchs ohne Team, das ihn operiert. Linkerd ist schlanker, Consul Connect f&uuml;r Nicht-K8s.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Was, wenn der Health-Check l&uuml;gt?</h3>
<p>Consul nimmt 200 als Beweis. Was, wenn der Service kaputt ist, aber `/health` <span class="hl">trotzdem 200</span> sagt?</p>
<code>Consul = Indexer, kein Orakel</code>
<aside class="notes"><strong>Meine Antwort:</strong> Dann routet Consul munter Traffic auf eine kaputte Instanz. Drei typische Failure-Modes: <strong>Zombie-Service</strong> (HTTP lebt, Worker-Pool tot), <strong>Backend-Abh&auml;ngigkeit weg</strong> (DB nicht erreichbar, aber `/health` antwortet trotzdem), <strong>langsame Degradation</strong> (P99 von 50 ms auf 5 s &mdash; technisch gesund, faktisch unbrauchbar). L&ouml;sung: der <em>richtige</em> Check (TCP &rarr; HTTP &rarr; Probe-Logik &rarr; Synthetic).<br><strong>Spicy:</strong> Service Discovery ist <strong>so gut wie der schlechteste Health-Check, der ihr darin pflegt</strong>. Ein Flight-Service, der l&uuml;gt, ist schlimmer als gar keine Registry &mdash; weil ihr euch in Sicherheit wiegt.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Consul im K8s-Cluster?</h3>
<p>Kubernetes bringt seine eigene Service Discovery mit. Ist Consul daneben <span class="hl">noch sinnvoll</span>?</p>
<code>kube-dns vs. Consul</code>
<aside class="notes"><strong>Meine Antwort:</strong> Im reinen K8s-Setup: <strong>nein.</strong> Der API-Server ist die Registry, `kube-dns` / CoreDNS l&ouml;st <code>flight-service.default.svc.cluster.local</code> auf, `kube-proxy` macht die Lastverteilung &mdash; alles ohne Zusatzkomponente, ohne Self-Registration im Code. Ein zweites Consul daneben w&auml;re reines Duplikat plus Sync-Problem.<br><strong>Wann doch?</strong> (1) <strong>Hybrid-/Multi-Cluster</strong>: Services laufen teils in K8s, teils auf VMs / Bare Metal / anderem Cluster &mdash; eine cluster-&uuml;bergreifende Registry. (2) <strong>Consul Connect als Service Mesh</strong> (mTLS, L7-Routing, Intentions) &mdash; dann ist Discovery nur Beiwerk, das Mesh ist der Grund. (3) <strong>Legacy-Brownfield</strong> mit existierender Consul-Infrastruktur, bei der K8s schrittweise dazukommt.<br><strong>Spicy:</strong> &bdquo;Wir nehmen Consul, weil wir das immer so machen&ldquo; ist in einem reinen K8s-Setup eine teure Gewohnheit. Die Frage ist nicht &bdquo;Consul ja/nein?&ldquo;, sondern &bdquo;welches Problem habe ich, das K8s nicht schon l&ouml;st?&ldquo; &mdash; meist hei&szlig;t die Antwort dann Mesh, nicht Discovery.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Was passiert beim STOP?</h3>
<p>Beim Start melden sich die Services an. Beim `docker stop` &mdash; verschwindet der Eintrag <span class="hl">sofort</span>?</p>
<code>10&ndash;30 s Traffic ins Leere</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein. Ohne Graceful Shutdown bleibt der Eintrag, bis der Health-Check nach mehreren Misses ausschl&auml;gt &mdash; typisch 10&ndash;30 s, in denen Traffic auf eine tote Instanz l&auml;uft (Connection-Refused). L&ouml;sungen: <strong>Graceful Shutdown</strong> mit aktivem `deregister`, schnellerer Health-Check-Intervall (kostet Last), <strong>Out-of-Service-Mode</strong> (keine neuen Requests, laufende beenden, dann Stop).<br><strong>Spicy:</strong> Service Discovery ist immer <strong>eventually consistent</strong>. Es gibt <strong>immer</strong> ein Zeitfenster, in dem Aufrufer auf tote Endpoints sto&szlig;en &mdash; genau einer der Gr&uuml;nde, warum wir in Story 3 / 4 Resilience-Patterns brauchen.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">7</span> Gefunden &ne; gesund</h3>
<p>Consul sagt: Service ist da. Er antwortet trotzdem mit <span class="hl">500</span>, im Schneckentempo oder gar nicht. Wer f&auml;ngt das ab?</p>
<code>Flight tot &rarr; Booking tot?</code>
<aside class="notes"><strong>&Uuml;berleitung zu Story 3:</strong> Discovery beantwortet nur <em>wo</em> ein Service ist, nicht <em>ob er funktioniert</em>. Selbst mit perfektem Health-Check bleibt das eventually-consistent-Fenster aus Frage 6 &ndash; und der Fall, dass ein registrierter Service einfach kaputt antwortet. Genau hier setzt Story 3 an: Circuit Breaker, Timeout, Fallback, damit ein ausgefallener Flight-Service nicht den Booking-Service mitrei&szlig;t (Teilbuchung Hotel + Mietwagen ohne Flug).<br><strong>Knackig:</strong> &bdquo;Story 2 hat den Service <em>gefunden</em>. Story 3 fragt: und wenn das, was wir gefunden haben, nicht antwortet?&ldquo;</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story2.md</code>.
</aside>
