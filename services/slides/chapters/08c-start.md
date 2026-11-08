<div class="page">

<p class="kicker">Start &middot; Dashboard, Schritt 0</p>

## Einklinken, dann <span class="hl">los</span>

<div class="page-body">

<p class="lead">Euer Mini-Service aus der Hausaufgabe wird Teil des Workshop-Stacks.</p>

<div class="cards cards-3 start-cards">
<div class="card">
<h3>1 &middot; Einklinken</h3>
<p>In <code>services/.env</code> auf euren Service zeigen:</p>
<p><code>CUSTOM_BOOKING_PATH=/pfad/zum/service</code></p>
<p>Stack neu starten (Ctrl+C):</p>
<p><code>make docker-up-hub</code></p>
</div>
<div class="card">
<h3>2 &middot; Testen</h3>
<p><code>curl -i localhost:8099/health</code> liefert 200</p>
<p><code>docker logs booking-custom</code> zeigt einen Kasten <strong>BEISPIEL-L&Ouml;SUNG</strong>? Dann l&auml;uft noch nicht euer Service.</p>
</div>
<div class="card">
<h3>3 &middot; Weiterarbeiten</h3>
<p>Nach jeder Code-&Auml;nderung, zweites Terminal in <code>services/</code>:</p>
<p><code>make docker-rebuild-custom</code></p>
<p>Im Dashboard: <em>Referenz</em> &harr; <em>Mein Service</em></p>
</div>
</div>

<div class="foot">
<div class="callout">Alles zum Nachlesen: http://localhost, Schritt 0 &bdquo;Start&ldquo;</div>
</div>

</div>

</div>

Note:
- Gemeinsam durchgehen, jeder klinkt jetzt seinen Service ein. Wer h&auml;ngt, setzt sich zu jemandem, bei dem es l&auml;uft.
