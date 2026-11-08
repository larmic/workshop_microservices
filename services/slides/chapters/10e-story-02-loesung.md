<div class="page">

<p class="kicker">Story 2 &middot; Musterl&ouml;sung</p>

<div class="page-head">

## Eine m&ouml;gliche L&ouml;sung

<span class="hint">Nicht die einzige. Entscheidend ist die Begr&uuml;ndung.</span>
</div>

<div class="page-body sol-body">

<div class="sol-grid">
<h4>Storno</h4>
<p class="sol-reason">Ein Storno hat eigene Daten: Grund und Zeitpunkt. Als eigene Ressource kann der Kundenservice sie sp&auml;ter nachlesen, die Buchung bleibt erhalten.</p>
<div class="sol-call"><span class="m">POST</span> /bookings/{id}/cancellation <span class="dim">{ reason }</span></div>
<div class="sol-codes"><span class="sc ok">201 angelegt</span><span class="sc">404 unbekannt</span><span class="sc">409 schon storniert</span></div>
<span class="sol-idem no">nicht idempotent</span>
<div class="sol-call"><span class="m">GET</span> /bookings/{id}/cancellation</div>
<div class="sol-codes"><span class="sc ok">200 Grund, Zeitpunkt</span><span class="sc">404 kein Storno</span></div>
<span class="sol-idem">idempotent</span>
<h4>Umbuchung</h4>
<p class="sol-reason">Nur der Flug &auml;ndert sich, Hotel und Mietwagen bleiben. PUT auf die Sub-Ressource ersetzt genau das, ein Retry ist harmlos.</p>
<div class="sol-call"><span class="m">PUT</span> /bookings/{id}/flight <span class="dim">{ flightId }</span></div>
<div class="sol-codes"><span class="sc ok">200 umgebucht</span><span class="sc">404 unbekannt</span><span class="sc">409 Flug voll</span><span class="sc err">503 Backend weg</span></div>
<span class="sol-idem">idempotent</span>
</div>

</div>

</div>

Note:
- Erst die Varianten der Teams sammeln, dann aufdecken.
- Storno kommt in Story 6 als Kompensation zur&uuml;ck.
