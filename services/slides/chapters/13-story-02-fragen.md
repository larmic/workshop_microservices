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
<h3><span class="numeral">3</span> Client-Side LB vs. Service Mesh</h3>
<p>Unser Resolver w&auml;hlt zuf&auml;llig. Warum baut die halbe Branche stattdessen <span class="hl">Envoy / Istio / Linkerd</span>?</p>
<code>N Sprachen = N Implementierungen</code>
<aside class="notes"><strong>Meine Antwort:</strong> Client-Side LB hat vier unangenehme Eigenschaften: (1) Logik wird pro Sprach-Stack neu gebaut, (2) kein zentrales Tuning (jede &Auml;nderung = Rebuild aller Services), (3) Beobachtbarkeit ist verteilt, (4) mTLS muss jeder Client selbst k&ouml;nnen. Service Mesh packt das in den Sidecar-Proxy, sprachunabh&auml;ngig.<br><strong>Spicy:</strong> Client-Side LB ist richtig f&uuml;r homogene Sprachlandschaft + wenige Services. Bei vielen Teams / Sprachen / Services ist Mesh g&uuml;nstiger &mdash; initiale Kosten (Ops, Komplexit&auml;t) aber hoch. Bewusste Architektur-Entscheidung, kein No-Brainer.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Was, wenn der Health-Check l&uuml;gt?</h3>
<p>Consul nimmt 200 als Beweis. Was, wenn der Service kaputt ist, aber `/health` <span class="hl">trotzdem 200</span> sagt?</p>
<code>Consul = Indexer, kein Orakel</code>
<aside class="notes"><strong>Meine Antwort:</strong> Dann routet Consul munter Traffic auf eine kaputte Instanz. Drei typische Failure-Modes: <strong>Zombie-Service</strong> (HTTP lebt, Worker-Pool tot), <strong>Backend-Abh&auml;ngigkeit weg</strong> (DB nicht erreichbar, aber `/health` antwortet trotzdem), <strong>langsame Degradation</strong> (P99 von 50 ms auf 5 s &mdash; technisch gesund, faktisch unbrauchbar). L&ouml;sung: der <em>richtige</em> Check (TCP &rarr; HTTP &rarr; Probe-Logik &rarr; Synthetic).<br><strong>Spicy:</strong> Service Discovery ist <strong>so gut wie der schlechteste Health-Check, der ihr darin pflegt</strong>. Ein Flight-Service, der l&uuml;gt, ist schlimmer als gar keine Registry &mdash; weil ihr euch in Sicherheit wiegt.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Canary &amp; Versionierung</h3>
<p>Wie deploye ich `flight-service` v2 mit <span class="hl">5 %</span> Traffic, ohne dass alle sofort drauf gehen?</p>
<code>5 % v2 / 95 % v1?</code>
<aside class="notes"><strong>Meine Antwort:</strong> Mit unserem Stand: <strong>gar nicht.</strong> Zufallsauswahl unter allen gesunden Instanzen &mdash; v2 bekommt sofort den gleichen Anteil wie v1. Drei L&ouml;sungspfade: (1) <strong>Tags / Metadaten</strong> in Consul (Resolver liest Tag, gewichtet &mdash; muss in jedem Client gebaut werden), (2) <strong>separate logische Namen</strong> (sauberer, aber Versionierung leakt nach oben), (3) <strong>Service Mesh / API-Gateway</strong> (Traffic-Splitting zentral konfiguriert).<br><strong>Spicy:</strong> Service Discovery l&ouml;st &bdquo;wo l&auml;uft das?&ldquo;, nicht &bdquo;welche Version will ich?&ldquo;. Sobald Canary / Blue-Green / A/B-Tests im Spiel sind, kommt eine zweite Schicht ins Spiel.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Was passiert beim STOP?</h3>
<p>Beim Start melden sich die Services an. Beim `docker stop` &mdash; verschwindet der Eintrag <span class="hl">sofort</span>?</p>
<code>10&ndash;30 s Traffic ins Leere</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein. Ohne Graceful Shutdown bleibt der Eintrag, bis der Health-Check nach mehreren Misses ausschl&auml;gt &mdash; typisch 10&ndash;30 s, in denen Traffic auf eine tote Instanz l&auml;uft (Connection-Refused). L&ouml;sungen: <strong>Graceful Shutdown</strong> mit aktivem `deregister`, schnellerer Health-Check-Intervall (kostet Last), <strong>Out-of-Service-Mode</strong> (keine neuen Requests, laufende beenden, dann Stop).<br><strong>Spicy:</strong> Service Discovery ist immer <strong>eventually consistent</strong>. Es gibt <strong>immer</strong> ein Zeitfenster, in dem Aufrufer auf tote Endpoints sto&szlig;en &mdash; genau einer der Gr&uuml;nde, warum wir in Story 3 / 4 Resilience-Patterns brauchen.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story2.md</code>.
</aside>
