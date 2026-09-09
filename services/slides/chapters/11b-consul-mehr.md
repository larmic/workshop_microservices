<div class="page">

<p class="kicker">Bewusst weggelassen</p>

## Consul kann mehr

<div class="page-body">

<div class="cards cards-3 compact">
<div class="card">
<h3>Key-Value Store</h3>
<p>Gemeinsamer Konfigurations-Speicher f&uuml;r alle Services.</p>
</div>
<div class="card">
<h3>Service Mesh</h3>
<p>Verschl&uuml;sselt den Verkehr zwischen Services automatisch.</p>
</div>
<div class="card">
<h3>DNS-Interface</h3>
<p>Aufl&ouml;sung per <code>flight-service.service.consul</code>.</p>
</div>
<div class="card">
<h3>Sidecar-Proxy</h3>
<p>Ein Proxy pro Service, ohne Zutun des Codes.</p>
</div>
<div class="card">
<h3>Intentions</h3>
<p>Hinterlegen, wer welchen Service aufrufen darf.</p>
</div>
<div class="card">
<h3>Multi-Datacenter</h3>
<p>Mehrere Rechenzentren zu einer Service-Sicht verbinden.</p>
</div>
</div>

<div class="foot">
<div class="callout">Wir nutzen bewusst nur die HTTP-Registry.</div>
</div>

</div>

</div>

Note:
- Diskussions-Anker: Wer betreibt Consul wirklich? Meist ein Platform- oder DevOps-Team, nicht das App-Team. Die Registry selbst ist ein Single Point of Failure, wenn man sie nicht clustert (3 oder 5 Server).
- Service Mesh lohnt sich, wenn ihr mTLS &uuml;berall braucht (Compliance, Zero Trust). Kostet Komplexit&auml;t (Sidecar pro Pod) und Latenz.
- Kubernetes-Welt: dort machen CoreDNS plus Service/Endpoints die HTTP-Schicht, Istio oder Linkerd den Mesh-Teil. Consul ist eher dann attraktiv, wenn ihr Kubernetes- und Nicht-Kubernetes-Workloads gemischt habt.
- KV-Store: praktisch, aber kein Ersatz f&uuml;r eine Datenbank. Eher Feature-Flags, Config-Snippets, dynamische Routing-Regeln.
- Take-away: Service Discovery ist die <em>erste</em> Funktion einer Plattform-Ebene wie Consul, nicht die einzige.
