<div class="page">

<p class="kicker">Story 2 &middot; Musterl&ouml;sung</p>

<div class="page-head">

## Eine m&ouml;gliche L&ouml;sung

<span class="hint">Nicht die einzige. Entscheidend ist die Begr&uuml;ndung.</span>
</div>

<div class="page-body sol-body">

<div class="sol-group fragment">
<h4>Storno</h4>
<div class="sol-ep">
<div class="sol-call"><span class="m">POST</span> /bookings/{id}/cancellation <span class="dim">{ reason }</span></div>
<div class="sol-codes"><span class="sc ok">201 angelegt</span><span class="sc">404 unbekannt</span><span class="sc">409 schon storniert</span></div>
<span class="sol-idem no">nicht idempotent</span>
</div>
<div class="sol-ep">
<div class="sol-call"><span class="m">GET</span> /bookings/{id}/cancellation</div>
<div class="sol-codes"><span class="sc ok">200 Grund, Zeitpunkt</span><span class="sc">404 kein Storno</span></div>
<span class="sol-idem">idempotent</span>
</div>
</div>

<div class="sol-group fragment">
<h4>Umbuchung</h4>
<div class="sol-ep">
<div class="sol-call"><span class="m">PUT</span> /bookings/{id}/flight <span class="dim">{ flightId }</span></div>
<div class="sol-codes"><span class="sc ok">200 umgebucht</span><span class="sc">404 unbekannt</span><span class="sc">409 Flug voll</span><span class="sc err">503 Backend weg</span></div>
<span class="sol-idem">idempotent</span>
</div>
</div>

<div class="sol-why fragment">
<div><strong>Storno als eigene Ressource</strong>Der Grund bleibt nachvollziehbar, die Buchung verschwindet nicht.</div>
<div><strong>Sub-Ressource mit PUT</strong>Klar, was sich &auml;ndert. Ein Retry ist harmlos.</div>
<div><strong>Fehler &uuml;ber Status-Codes</strong>Nie <code>200</code> mit Fehler-Body.</div>
</div>

<p class="sol-alt fragment">Auch vertretbar: <code>DELETE /bookings/{id}</code> oder <code>PATCH</code> mit Status, wenn die Begr&uuml;ndung tr&auml;gt.</p>

</div>

</div>

Note:
- Erst die Varianten der Teams sammeln, dann aufdecken.
- Storno kommt in Story 6 als Kompensation zur&uuml;ck.
