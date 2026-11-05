<div class="page">

<p class="kicker">Story 4 &middot; &Uuml;bung</p>

<div class="page-head">

## Wenn der Flug ausf&auml;llt

<span class="badge design">&asymp; 20 min</span>
</div>

<p class="subtitle">Ihr seid die Suchanfragen. Flight ist tot. Booking hat den <span class="hl">Schalter</span>.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> ist eine Person. Sie hat den Schalter: eine Strichliste und ein Schild mit CLOSED, OPEN, HALF_OPEN.</li>
<li><strong>Flight, Hotel, Car</strong> sind drei Leute. Hotel und Car antworten sofort. Flight ist abgest&uuml;rzt und sagt gar nichts.</li>
<li><strong>Alle anderen</strong> sind Suchanfragen auf <code>GET /booking/offers</code> und stellen sich bei Booking an.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li><strong>Runde 1, ohne Schalter:</strong> Booking fragt Flight, Flight schweigt, Booking z&auml;hlt laut bis drei: Timeout, &bdquo;Fehler&ldquo;. Die Anfrage stellt sich hinten wieder an.</li>
<li><strong>Runde 2, mit Schalter:</strong> Jeder Timeout ein Strich auf Bookings Liste. Beim f&uuml;nften Strich dreht Booking sein Schild auf OPEN.</li>
<li>Ab jetzt fragt Booking nur Hotel und Car. Karte: Hotel &#10003;, Car &#10003;, <code>flights: []</code>. Die Anfrage setzt sich.</li>
<li>Alle 30 Sekunden: HALF_OPEN, ein Probe-Call an Flight. Schweigt Flight: OPEN. Antwortet Flight: CLOSED.</li>
<li><strong>Runde 3, freies Spiel:</strong> Flight antwortet, wann es will. Booking h&auml;lt sich an seine Regeln.</li>
</ol>
</div>
</div>

</div>

</div>

Note:
- Hook: &bdquo;Story 3 hat uns geholfen, Services zu finden. Jetzt haben wir einen gefunden, und er antwortet nicht.&ldquo; Heute wird gespielt, nicht programmiert.
- Material: Strichliste, Schild mit drei Zust&auml;nden, Karten &bdquo;Hotel &#10003;, Car &#10003;, flights: []&ldquo;, eine Uhr mit Sekunden.
- Erwartung Runde 1: Die Schlange bleibt gleich lang, jede Anfrage kostet drei Sekunden, niemand setzt sich. Hotel und Car werden gar nicht gefragt, so macht es die Referenz.
- Erwartung Runde 2: Nach dem f&uuml;nften Strich leert sich die Schlange sichtbar. Hotel und Car arbeiten weiter, der Schalter gilt pro Backend.
- Du bist die Uhr: alle 30 Sekunden &bdquo;HALF_OPEN&ldquo; rufen. Flight schweigt zweimal, beim dritten Mal Flight wieder antworten lassen: CLOSED, volle Angebote.
- Runde 3: Der Schalter pendelt, ein wackliges Flight kommt nie dauerhaft auf CLOSED. Br&uuml;cke zu Recap 1 und 3: Was z&auml;hlt als Fehler, warum nur ein Probe-Call.
- Danach die Demo im Dashboard: Chaos-Schalter, Status-Grid, Header <code>X-Circuit-Open</code>.
- Vollst&auml;ndige Beschreibung und Material zum Nachbauen: <code>docs/stories/story-04-circuit-breaker.md</code>.
