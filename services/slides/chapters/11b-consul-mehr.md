## Consul kann mehr

<p class="subtitle">&hellip; als wir hier nutzen</p>

<div class="box">

- **Key-Value Store** &mdash; verteilte Konfiguration, Feature-Flags
- **Distributed Locks / Sessions** &mdash; Leader-Election ohne extra ZooKeeper
- **Service Mesh (Consul Connect)** &mdash; mTLS, Sidecar-Proxies, Zero-Trust zwischen Services
- **Intentions** &mdash; Policy &bdquo;Service A darf Service B aufrufen&ldquo; statt Firewall-Regeln
- **DNS-Interface** &mdash; `flight-service.service.consul` statt HTTP-API
- **Watches** &mdash; Events bei Topologie-&Auml;nderungen, kein Polling
- **Multi-Datacenter Federation** &mdash; globale Service-Sicht &uuml;ber DCs hinweg
- **ACLs** &mdash; Token-basierte Zugriffskontrolle auf die Registry selbst

</div>

<p class="quote">Im Workshop nutzen wir bewusst <span class="hl">nur die HTTP-Registry</span> &mdash; der Rest ist Diskussion.</p>

Note:
- Diskussions-Anker: Wer betreibt Consul wirklich? Meist ein Platform-/DevOps-Team &mdash; nicht das App-Team. Die Registry selbst ist ein Single Point of Failure, wenn man sie nicht clustered (3 oder 5 Server).
- Service Mesh: lohnt sich, wenn ihr mTLS &uuml;berall braucht (Compliance, Zero-Trust). Kostet Komplexit&auml;t (Sidecar pro Pod) und Latenz.
- K8s-Welt: dort macht **CoreDNS + Service/Endpoints** die HTTP-Schicht; **Istio / Linkerd** machen den Mesh-Teil. Consul ist eher dann attraktiv, wenn ihr K8s und Nicht-K8s-Workloads gemischt habt.
- KV-Store: praktisch, aber nicht missbrauchen &mdash; kein Ersatz f&uuml;r eine echte Datenbank. Eher: Feature-Flags, Config-Snippets, dynamische Routing-Regeln.
- Take-away: Service Discovery ist die <em>erste</em> Funktion einer Plattform-Ebene wie Consul &mdash; nicht die einzige.
