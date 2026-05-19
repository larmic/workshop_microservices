## Service Discovery

<p class="subtitle">Services finden sich &uuml;ber Namen &mdash; nicht &uuml;ber URLs</p>

<div class="cols">
<div>

### Was ist das?

- Services finden einander &uuml;ber <span class="hl">logische Namen</span> statt fester IPs/URLs
- **Registry** als Single Source of Truth: Services melden sich an und ab
- **Health Checks** blenden ungesunde Instanzen automatisch aus
- Topologie ist **dynamisch**: Skalierung, Failover, Rolling Deploys ohne Config-Edit
- Ersetzt manuell gepflegte URL-Listen / `/etc/hosts` / hardcoded Env-Vars
- Zwei Spielarten: **Client-Side** (Resolver fragt Registry) vs. **Server-Side** (Load Balancer fragt Registry)

</div>
<div>

### Am Markt

<div class="chip-row">
  <span class="chip brand">HashiCorp Consul</span>
  <span class="chip">Netflix Eureka</span>
  <span class="chip">Apache Zookeeper</span>
  <span class="chip">etcd</span>
  <span class="chip">Kubernetes (Service + CoreDNS)</span>
  <span class="chip">AWS Cloud Map</span>
  <span class="chip">Spring Cloud Discovery</span>
</div>

</div>
</div>

Note:
- Hook: &bdquo;Im Monolithen kennen sich Module &uuml;ber Funktionsaufrufe. Im verteilten System kennen sich Services &uuml;ber &hellip; was eigentlich?&ldquo; Statische URLs (Story 1) funktionieren genau so lange, bis ihr skaliert, deployt oder eine Instanz ausf&auml;llt.
- Client-Side vs. Server-Side: Wir bauen im Workshop Client-Side &mdash; jeder Service hat einen Resolver und w&auml;hlt selbst. Server-Side hei&szlig;t Load Balancer / API-Gateway fragt die Registry (Traefik kann das, AWS ALB mit Cloud Map, etc.).
- Wer was wo nutzt: Eureka stark in Spring-Welt, etcd unter der Haube von Kubernetes, Consul polyglot. In K8s selten extra Service Discovery &mdash; Service + CoreDNS reicht meist.
- &Uuml;berleitung: Wir schauen jetzt konkret auf Consul, weil wir es im Workshop benutzen.
