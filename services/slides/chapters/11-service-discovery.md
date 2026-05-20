## Service Discovery

<p class="subtitle">Services finden sich &uuml;ber Namen &mdash; nicht &uuml;ber URLs</p>

<div class="factor-row">

<div class="factor fragment">
<h3>Logische Namen</h3>
<p>Code ruft <code>flights-service</code> &mdash; nicht <code>10.0.0.5:8080</code>.</p>
<aside class="notes">Statt fester IPs/URLs gibt es einen logischen Namen, hinter dem 1..n Instanzen stehen. Genau das, was DNS f&uuml;r Hosts macht &mdash; nur dynamischer und mit Health-Wissen.</aside>
</div>

<div class="factor fragment">
<h3>Registry</h3>
<p>Single Source of Truth: Services melden sich an und ab.</p>
<aside class="notes">Beim Start registriert sich der Service (Name, Adresse, Port, Tags). Beim sauberen Stop deregistriert er sich. Wer wissen will, wo der Service l&auml;uft, fragt die Registry &mdash; nicht die Wiki-Seite.</aside>
</div>

<div class="factor fragment">
<h3>Health Checks</h3>
<p>Ungesunde Instanzen fliegen automatisch raus.</p>
<aside class="notes">Registry pollt (oder bekommt Heartbeats) und entfernt kranke Instanzen aus dem Pool. Caller bekommt nur &bdquo;lebende&ldquo; Adressen zur&uuml;ck. Vorteil: Failover ohne manuelles Eingreifen.</aside>
</div>

<div class="factor fragment">
<h3>Dynamische Topologie</h3>
<p>Adressen &auml;ndern sich, Code nicht &mdash; kein <code>/etc/hosts</code>, keine URL-Liste.</p>
<aside class="notes">In Container-Welten ist die Topologie kein statisches Bild mehr: Skalieren = neue Instanz, Crash = neue IP, Rolling Deploy = alte raus / neue rein. Service Discovery folgt automatisch. Was sie abl&ouml;st: <code>/etc/hosts</code>-Pflege, hardcoded <code>BOOKING_URL=...</code> Env-Vars, Excel-Listen im Ops-Wiki.</aside>
</div>

<div class="factor fragment">
<h3>Client- vs. Server-Side</h3>
<p>Resolver beim Caller &mdash; oder Load Balancer fragt die Registry.</p>
<aside class="notes">Client-Side: jeder Service hat einen Resolver und w&auml;hlt selbst (was wir im Workshop bauen). Server-Side: ein Load Balancer / API-Gateway davor fragt die Registry (Traefik, AWS ALB + Cloud Map, &hellip;). Trade-off: Client-Side spart einen Hop, Server-Side ist sprach-agnostisch.</aside>
</div>

</div>

<div class="market-row">

### Am Markt

<div class="chip-row">
  <span class="chip brand">HashiCorp Consul</span>
  <span class="chip">Netflix Eureka</span>
  <span class="chip">Apache Zookeeper</span>
  <span class="chip">AWS Cloud Map</span>
  <span class="chip">Spring Cloud Discovery</span>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>

Note:
- Hook: &bdquo;Im Monolithen kennen sich Module &uuml;ber Funktionsaufrufe. Im verteilten System kennen sich Services &uuml;ber &hellip; was eigentlich?&ldquo; Statische URLs (Story 1) funktionieren genau so lange, bis ihr skaliert, deployt oder eine Instanz ausf&auml;llt.
- Karten-Reihenfolge bewusst: erst das Konzept (Name), dann das Werkzeug (Registry + Health), dann der Payoff (dynamische Topologie), zuletzt die Geschmacks-Frage (Client vs. Server).
- Wer was wo nutzt: Eureka stark in Spring-Welt, etcd unter der Haube von Kubernetes, Consul polyglot. In K8s selten extra Service Discovery &mdash; Service + CoreDNS reicht meist.
- &Uuml;berleitung: Wir schauen jetzt konkret auf Consul, weil wir es im Workshop benutzen.
