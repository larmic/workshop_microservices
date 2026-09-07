<div class="page">

<p class="kicker">Architektur &amp; Kommunikation seit 1990</p>

## Der Weg zu Microservices

<div class="page-body">

<!-- Achtung: keine Leerzeilen innerhalb des SVG, sonst zerteilt der
     Markdown-Parser den HTML-Block und rendert die <text>-Elemente als Absatz. -->
<div class="timeline">
<svg viewBox="0 0 1180 412" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Zeitlinie 1990 bis 2022: Architekturen oben (Monolithen, SOA, Microservices, Modulithen, SCS), Technologien unten (RMI, CORBA, HTTP, SOAP, REST, JSON, Docker, gRPC)">
  <!-- Verbindungslinien: oben Architekturen, unten Technologien -->
  <line class="tl-tick" x1="70"   y1="46"  x2="70"   y2="143"/>
  <line class="tl-tick" x1="408" y1="46" x2="408" y2="227"/>
  <line class="tl-tick" x1="745" y1="46" x2="745" y2="295"/>
  <line class="tl-tick" x1="1049" y1="46" x2="1049" y2="126"/>
  <line class="tl-tick" x1="1150" y1="46" x2="1150" y2="83"/>
  <line class="tl-tick" x1="70"   y1="187" x2="70"   y2="376"/>
  <line class="tl-tick" x1="408" y1="271" x2="408" y2="376"/>
  <line class="tl-tick" x1="745" y1="339" x2="745" y2="376"/>
  <line class="tl-tick" x1="880" y1="303" x2="880" y2="376"/>
  <!-- Straße, x-Achse maßstäblich: 1990 bei x=70, 2022 bei x=1150 -->
  <path class="tl-road" d="M 0 170 C 260 130, 470 330, 730 318 C 990 306, 1010 90, 1180 105"/>
  <path class="tl-dash" d="M 0 170 C 260 130, 470 330, 730 318 C 990 306, 1010 90, 1180 105"/>
  <!-- Jahres-Punkte auf der Straße -->
  <circle class="tl-dot" cx="70" cy="165" r="22"/><text class="tl-year" x="70" y="165">1990</text>
  <circle class="tl-dot" cx="408" cy="249" r="22"/><text class="tl-year" x="408" y="249">2000</text>
  <circle class="tl-dot" cx="745" cy="317" r="22"/><text class="tl-year" x="745" y="317">2010</text>
  <circle class="tl-dot" cx="880" cy="281" r="22"/><text class="tl-year" x="880" y="281">2014</text>
  <circle class="tl-dot" cx="1049" cy="148" r="22"/><text class="tl-year" x="1049" y="148">2019</text>
  <circle class="tl-dot" cx="1150" cy="105" r="22"/><text class="tl-year" x="1150" y="105">2022</text>
  <!-- Architekturen (eine Zeile oben) -->
  <text class="tl-arch" x="70"   y="36" text-anchor="middle">Monolithen</text>
  <text class="tl-arch" x="408" y="36" text-anchor="middle">SOA</text>
  <text class="tl-arch" x="745" y="36" text-anchor="middle">Microservices</text>
  <text class="tl-arch" x="1049" y="36" text-anchor="middle">Modulithen</text>
  <text class="tl-arch" x="1150" y="36" text-anchor="middle">SCS</text>
  <!-- Technologien (eine Zeile unten) -->
  <text class="tl-tech" x="70"  y="399" text-anchor="middle">RMI &#183; CORBA &#183; HTTP</text>
  <text class="tl-tech" x="408" y="399" text-anchor="middle">SOAP</text>
  <text class="tl-tech" x="745" y="399" text-anchor="middle">REST &#183; JSON</text>
  <text class="tl-tech" x="880" y="399" text-anchor="middle">Docker &#183; gRPC</text>
</svg>
</div>

</div>

</div>

Note:
- Vorher: Mainframes
- SOA: wenige, gro&szlig;e Services
- Microservices: feingranular, eine Aufgabe, ein Service
- 2014: Docker 1.0 erscheint, gRPC folgt 2015 (steht mit am 2014er-Punkt).
- SCS: Self Contained System
- REST + JSON ist der heutige Default. Was &bdquo;RESTful&ldquo; wirklich hei&szlig;t und warum das mehr als Stil ist, kommt in Story 2 (Design-Session) auf den Tisch.
- Realit&auml;tscheck: Network Latency | Deployment Hell | Debugging Horror
