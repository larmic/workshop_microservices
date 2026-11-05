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
- <strong>Nur &Uuml;berblick:</strong> Nicht jede Karte erkl&auml;ren. F&uuml;r Story 3 z&auml;hlen nur Registry und Health-Check von der Folie davor.
- <strong>Optional:</strong> Der Key-Value Store steht als Zusatzaufgabe in der Story, f&uuml;r Schnelle.
- <strong>Ausblick:</strong> Ein Mesh verlagert Discovery und mTLS in den Sidecar-Proxy, Intentions regeln die Aufrufrechte. Operativ teuer. DNS ist bequem, aber TTL-Caching verz&ouml;gert den Failover.
- <strong>&Uuml;berleitung:</strong> &bdquo;Genug Theorie, jetzt l&ouml;st euer Booking-Service die Backends &uuml;ber Namen auf.&ldquo;
