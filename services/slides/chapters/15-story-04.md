<div class="page">

<p class="kicker">Story 4 &middot; Rollenspiel</p>

<div class="page-head">

## Wenn der Flug ausf&auml;llt

<span class="badge design">&asymp; 60 min <small>davon 20 Spiel</small></span>
</div>

<p class="subtitle">Ihr seid die Requests. Flight ist krank, die Antwort steht auf der <span class="hl">Karte</span>.</p>

<div class="page-body">

<div class="market">
<h4>So spielen wir</h4>
<div class="cards cards-3 violet compact">
<div class="card">
<h3>Eine Person ist Booking</h3>
<p>Booking gibt den Stapel streng der Reihe nach aus und schaut aufs Circuit-Breaker-Board.</p>
</div>
<div class="card">
<h3>Alle anderen sind Requests</h3>
<p>Antwort laut vorlesen, Karte dorthin, wo das Band sagt, <strong>hinten wieder anstellen</strong>.</p>
<code>GET /booking/offers</code>
</div>
<div class="card">
<h3>Graue Karte: STOPP</h3>
<p>An den grauen Karten halten wir an und reden.</p>
</div>
</div>
</div>

</div>

</div>

Note:
- <strong>Hook:</strong> &bdquo;Story 3 hat uns Services finden lassen. Jetzt haben wir einen gefunden, und er antwortet nicht.&ldquo;
- <strong>Setup:</strong>
  - Ab 4 Personen, derselbe Stapel f&uuml;r jede Gruppengr&ouml;&szlig;e.
  - Booking bekommt Stapel 00&ndash;27, Legende und Zustandskarten. Das Board kommt erst mit STOPP 1. Der Trainer spielt keine Rolle, er moderiert nur an den STOPPs.
  - Farbe der Karte = wohin: dunkel Ablage, rot Board, gelb Probe-Feld, grau Trainer.
  - Ablauf, STOPP-Fragen und Kartentabelle: Trainer-Skript auf der Materialseite.
- <strong>Demo:</strong>
  - Flight auf &bdquo;H&auml;ngt&ldquo; (10 s, Timeout 3 s): f&uuml;nf Calls warten je 3 s, immer 200 mit <code>flights: []</code> und <code>X-Fallback: flight</code>, ab dem 6. sofort mit <code>X-Circuit-Open: flight</code>. Schneller: &bdquo;Fehler&ldquo;.
  - Flight zur&uuml;ck auf &bdquo;Normal&ldquo;: nach 30 s HALF_OPEN, das Panel zeigt es schon ohne Call. Die erfolgreiche Probe schlie&szlig;t.
  - Details: <code>docs/stories/story-04-circuit-breaker.md</code>.
