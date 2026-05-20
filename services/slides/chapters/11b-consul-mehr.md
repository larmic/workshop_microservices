<!-- .slide: data-background-image="./assets/service_discovery.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Consul kann mehr

<p class="subtitle">&hellip; als wir hier nutzen</p>

<div class="box">

- **Key-Value Store** &mdash; gemeinsamer Konfigurations-Speicher f&uuml;r alle Services
- **Service Mesh** &mdash; verschl&uuml;sselt Service-zu-Service-Verkehr automatisch
- **DNS-Interface** &mdash; Services per `flight-service.service.consul` finden
- **Sidecar-Proxy** &mdash; automatischer Proxy zwischen Service und Client
- **Intentions (Firewallersatz)** &mdash; hinterlegen, wer welchen Service aufrufen darf
- **Multi-Datacenter** &mdash; mehrere Rechenzentren zu einer Service-Sicht verbinden
- ...

</div>

<p class="quote">Im Workshop nutzen wir bewusst <span class="hl">nur die HTTP-Registry</span>.</p>

Note:
- Diskussions-Anker: Wer betreibt Consul wirklich? Meist ein Platform-/DevOps-Team &mdash; nicht das App-Team. Die Registry selbst ist ein Single Point of Failure, wenn man sie nicht clustered (3 oder 5 Server).
- Service Mesh: lohnt sich, wenn ihr mTLS &uuml;berall braucht (Compliance, Zero-Trust). Kostet Komplexit&auml;t (Sidecar pro Pod) und Latenz.
- K8s-Welt: dort macht **CoreDNS + Service/Endpoints** die HTTP-Schicht; **Istio / Linkerd** machen den Mesh-Teil. Consul ist eher dann attraktiv, wenn ihr K8s und Nicht-K8s-Workloads gemischt habt.
- KV-Store: praktisch, aber nicht missbrauchen &mdash; kein Ersatz f&uuml;r eine echte Datenbank. Eher: Feature-Flags, Config-Snippets, dynamische Routing-Regeln.
- Take-away: Service Discovery ist die <em>erste</em> Funktion einer Plattform-Ebene wie Consul &mdash; nicht die einzige.
