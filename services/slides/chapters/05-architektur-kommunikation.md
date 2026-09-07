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
  <line class="tl-tick" x1="401"  y1="46"  x2="401"  y2="225"/>
  <line class="tl-tick" x1="732"  y1="46"  x2="732"  y2="296"/>
  <line class="tl-tick" x1="1031" y1="46"  x2="1031" y2="141"/>
  <line class="tl-tick" x1="1148" y1="46"  x2="1130" y2="86"/>
  <line class="tl-tick" x1="70"   y1="187" x2="70"   y2="376"/>
  <line class="tl-tick" x1="401"  y1="269" x2="401"  y2="376"/>
  <line class="tl-tick" x1="732"  y1="340" x2="732"  y2="376"/>
  <line class="tl-tick" x1="865"  y1="310" x2="865"  y2="376"/>
  <!-- Straße, x-Achse maßstäblich: 1990 bei x=70, 2022 bei x=1130 -->
  <path class="tl-road" d="M 0 170 C 260 130, 470 330, 730 318 C 990 306, 1010 90, 1180 105"/>
  <path class="tl-dash" d="M 0 170 C 260 130, 470 330, 730 318 C 990 306, 1010 90, 1180 105"/>
  <!-- Jahres-Punkte auf der Straße -->
  <circle class="tl-dot" cx="70"   cy="165" r="22"/><text class="tl-year" x="70"   y="165">1990</text>
  <circle class="tl-dot" cx="401"  cy="247" r="22"/><text class="tl-year" x="401"  y="247">2000</text>
  <circle class="tl-dot" cx="732"  cy="318" r="22"/><text class="tl-year" x="732"  y="318">2010</text>
  <circle class="tl-dot" cx="865"  cy="288" r="22"/><text class="tl-year" x="865"  y="288">2014</text>
  <circle class="tl-dot" cx="1031" cy="163" r="22"/><text class="tl-year" x="1031" y="163">2019</text>
  <circle class="tl-dot" cx="1130" cy="108" r="22"/><text class="tl-year" x="1130" y="108">2022</text>
  <!-- Architekturen (eine Zeile oben) -->
  <text class="tl-arch" x="70"   y="36" text-anchor="middle">Monolithen</text>
  <text class="tl-arch" x="401"  y="36" text-anchor="middle">SOA</text>
  <text class="tl-arch" x="732"  y="36" text-anchor="middle">Microservices</text>
  <text class="tl-arch" x="1031" y="36" text-anchor="middle">Modulithen</text>
  <text class="tl-arch" x="1170" y="36" text-anchor="end">SCS</text>
  <!-- Technologien (eine Zeile unten) -->
  <text class="tl-tech" x="70"  y="399" text-anchor="middle">RMI &#183; CORBA &#183; HTTP</text>
  <text class="tl-tech" x="401" y="399" text-anchor="middle">SOAP</text>
  <text class="tl-tech" x="732" y="399" text-anchor="middle">REST &#183; JSON</text>
  <text class="tl-tech" x="865" y="399" text-anchor="middle">Docker &#183; gRPC</text>
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
