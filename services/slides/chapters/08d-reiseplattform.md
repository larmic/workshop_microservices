<div class="page">

<p class="kicker">Die Beispiel-Dom&auml;ne</p>

## Unsere <span class="hl">Reiseplattform</span>

<div class="page-body">

<p class="lead">Eine Reise sind drei Buchungen: Flug, Hotel, Mietwagen. Booking h&auml;lt sie zusammen.</p>

<!-- Achtung: keine Leerzeilen innerhalb des SVG, sonst zerteilt der
     Markdown-Parser den HTML-Block und rendert die <text>-Elemente als Absatz. -->
<div class="platform">
<svg viewBox="0 0 1000 350" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Der Kunde bucht eine Reise bei Booking. Booking ruft Flight, Hotel und Car auf.">
  <line class="pf-line" x1="500" y1="50" x2="500" y2="96"/>
  <path class="pf-line" d="M 500 186 V 222 M 170 256 V 222 H 830 V 256"/>
  <line class="pf-line" x1="500" y1="222" x2="500" y2="256"/>
  <rect class="pf-pill" x="400" y="4" width="200" height="46" rx="23"/>
  <text class="pf-pill-text" x="500" y="28">Kunde</text>
  <text class="pf-call" x="516" y="74">POST /bookings</text>
  <rect class="pf-main" x="330" y="96" width="340" height="90" rx="12"/>
  <text class="pf-main-name" x="500" y="131">Booking</text>
  <text class="pf-main-sub" x="500" y="162">euer Service</text>
  <rect class="pf-box" x="40" y="256" width="260" height="88" rx="12"/>
  <text class="pf-name" x="170" y="290">Flight</text>
  <text class="pf-sub" x="170" y="320">Flug buchen</text>
  <rect class="pf-box" x="370" y="256" width="260" height="88" rx="12"/>
  <text class="pf-name" x="500" y="290">Hotel</text>
  <text class="pf-sub" x="500" y="320">Zimmer buchen</text>
  <rect class="pf-box" x="700" y="256" width="260" height="88" rx="12"/>
  <text class="pf-name" x="830" y="290">Car</text>
  <text class="pf-sub" x="830" y="320">Mietwagen buchen</text>
</svg>
</div>

</div>

</div>

Note:
- Einmal das ganze Bild, bevor jeder seinen Service einklinkt: Der Kunde bucht eine Reise, Booking ruft daf&uuml;r Flight, Hotel und Car auf.
- Flight, Hotel und Car laufen fertig im Stack. Die Teams bauen nur Booking, und zwar zwei Tage lang denselben Service, Story f&uuml;r Story erweitert.
- Die drei Services sind absichtlich simpel. Spannend wird, was zwischen ihnen passiert: einer ist langsam, einer f&auml;llt aus, einer sagt nein.
