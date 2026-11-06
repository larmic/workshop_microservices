<div class="page">

<p class="kicker">Microservices sind auch nur Monolithen</p>

## &hellip; mit <span class="hl">Netzwerkproblemen</span>

<div class="page-body">

<p class="lead">&rarr; und ein Methodenaufruf wird zum Netzwerkaufruf</p>

<div class="cards cards-3 cards-display">
<div class="card">
<h3>Langsam</h3>
<p>blockiert alles, was hinter ihm wartet</p>
</div>
<div class="card">
<h3>L&uuml;gt</h3>
<p>200 zur&uuml;ck, obwohl nichts mehr geht</p>
</div>
<div class="card">
<h3>Kommt nie an</h3>
<p>niemand wei&szlig;, ob etwas passiert ist</p>
</div>
</div>

</div>

</div>

Note:
- These in den Raum stellen, dann <strong>live beweisen</strong> (Dashboard, Story-1-Abschnitt, 5 bis 7 Minuten):
  1. <code>GET /booking/offers</code> klicken: alles da, schnell.
  2. Flight auf <strong>Langsam</strong>, nochmal klicken: rund 2 Sekunden warten. Hotel und Auto warten mit.
  3. Wette: &bdquo;Flight auf Fehler. Bekommen wir wenigstens Hotel und Auto?&ldquo;
  4. Flight auf <strong>Fehler</strong>: 500, alles weg. Der Health-Check von Booking bleibt gr&uuml;n.
  5. Optional: Flight auf <strong>H&auml;ngt</strong>. F&uuml;nf Sekunden nichts, dann 500. Das ist &bdquo;Kommt nie an&ldquo;.
  6. Flight zur&uuml;ck auf <strong>Normal</strong>.
- Zu jeder Karte Daumen hoch oder runter: Gibt es das im Monolithen auch?
- Kern: Ein Methodenaufruf ist schnell, ehrlich und kommt an. Ein Netzwerkaufruf nicht. Alles danach ist der Umgang mit diesen drei F&auml;llen.
- L&auml;uft der Stack nicht: Demo weglassen, nur die Daumen-Runde.
