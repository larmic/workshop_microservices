<div class="page">

<p class="kicker">Architektur &amp; Kommunikation seit 1990</p>

## Der Weg zu Microservices

<div class="page-body">

<!-- Achtung: keine Leerzeilen innerhalb des SVG, sonst zerteilt der
     Markdown-Parser den HTML-Block und rendert die <text>-Elemente als Absatz. -->
<div class="timeline">
<svg viewBox="0 0 1180 412" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Zeitlinie 1990 bis 2022: Architekturen oben (Monolithen, SOA, Microservices, Modulithen, SCS), Technologien unten (RMI, CORBA, HTTP, SOAP, REST, JSON, Docker, gRPC)">
  <!-- Masten: von der Fahne bis unter die Straße, bei den Technologien weiter bis nach unten -->
  <line class="tl-tick" x1="70"   y1="46" x2="70"   y2="165"/>
  <line class="tl-tick" x1="408"  y1="46" x2="408"  y2="249"/>
  <line class="tl-tick" x1="745"  y1="46" x2="745"  y2="317"/>
  <line class="tl-tick" x1="1049" y1="46" x2="1049" y2="171"/>
  <line class="tl-tick" x1="1150" y1="46" x2="1150" y2="134"/>
  <line class="tl-tick" x1="70"   y1="165" x2="70"   y2="376"/>
  <line class="tl-tick" x1="408"  y1="249" x2="408"  y2="376"/>
  <line class="tl-tick" x1="745"  y1="317" x2="745"  y2="376"/>
  <line class="tl-tick" x1="880"  y1="46" x2="880"  y2="376"/>
  <!-- Straße, x-Achse maßstäblich: 1990 bei x=70, 2022 bei x=1150. Läuft über
       beide Folienränder hinaus (overflow: visible), die Enden sind abgeschnitten. -->
  <path class="tl-road" d="M -140 191 L 0 170 C 260 130, 470 330, 730 318 C 990 306, 1010 120, 1180 135 L 1320 147"/>
  <path class="tl-dash" d="M -140 191 L 0 170 C 260 130, 470 330, 730 318 C 990 306, 1010 120, 1180 135 L 1320 147"/>
  <!-- Fahnen mit der Jahreszahl am Mast, alle in einer Reihe unter den Architekturen -->
  <rect class="tl-flag" x="70" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="99" y="63">1990</text>
  <rect class="tl-flag" x="408" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="437" y="63">2000</text>
  <rect class="tl-flag" x="745" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="774" y="63">2010</text>
  <rect class="tl-flag" x="880" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="909" y="63">2014</text>
  <rect class="tl-flag" x="1049" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="1078" y="63">2019</text>
  <rect class="tl-flag" x="1150" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="1179" y="63">2022</text>
  <!-- Figuren auf der Straße, auf dem Jahrespunkt, aufrecht und in genau einer
       Farbe; Struktur entsteht durch Lücken, in denen die Straße durchscheint.
       Lokal um (0,0) gezeichnet, Position in transform. -->
  <g class="tl-figures" fill="#E4DDF8" stroke="none">
    <!-- 1990: Big Ball of Mud -->
    <g transform="translate(70 165)">
      <path d="M -30 4 C -32 -22, -12 -34, 6 -30 C 24 -38, 40 -22, 34 -4 C 44 12, 30 34, 10 32 C -8 42, -30 30, -30 4 Z"/>
      <path d="M -20 2 C -10 -20, 20 -16, 10 2 S -16 18, 0 22 S 26 10, 14 -6 S -14 -10, -10 10 S 16 24, 24 8" fill="none" stroke="#0A0349" stroke-width="3" stroke-linecap="round"/>
    </g>
    <!-- 2000: SOA, drei große Blöcke am Bus -->
    <g transform="translate(408 249)">
      <rect x="-30" y="12" width="60" height="6" rx="3"/>
      <rect x="-23" y="4" width="4" height="8"/><rect x="-2" y="4" width="4" height="8"/><rect x="19" y="4" width="4" height="8"/>
      <rect x="-30" y="-16" width="18" height="20" rx="3"/>
      <rect x="-9"  y="-16" width="18" height="20" rx="3"/>
      <rect x="12"  y="-16" width="18" height="20" rx="3"/>
    </g>
    <!-- 2010: viele kleine Services, dünne Fäden -->
    <g transform="translate(745 317)">
      <g fill="none" stroke="#E4DDF8" stroke-width="2">
        <line x1="-22" y1="-18" x2="4" y2="-22"/><line x1="4" y1="-22" x2="26" y2="-10"/>
        <line x1="4" y1="-22" x2="6" y2="2"/><line x1="-18" y1="8" x2="-10" y2="24"/>
        <line x1="-10" y1="24" x2="14" y2="22"/><line x1="26" y1="-10" x2="14" y2="22"/>
      </g>
      <rect x="-28" y="-24" width="12" height="12" rx="2.5"/><rect x="-2" y="-28" width="12" height="12" rx="2.5"/>
      <rect x="20" y="-16" width="12" height="12" rx="2.5"/><rect x="-24" y="2" width="12" height="12" rx="2.5"/>
      <rect x="0" y="-4" width="12" height="12" rx="2.5"/><rect x="-16" y="18" width="12" height="12" rx="2.5"/>
      <rect x="8" y="16" width="12" height="12" rx="2.5"/>
    </g>
    <!-- Risse vom Straßenrand zwischen 2010 und 2014: am Rand breit, nach innen spitz -->
    <g fill="#ffffff" stroke="none">
      <path d="M 800 361 L 810 361 L 806 349 L 810 341 L 805 333 L 803 325 L 800 334 L 803 342 L 799 350 Z"/>
      <path d="M 837 253 L 848 253 L 844 264 L 848 272 L 842 284 L 839 273 L 841 265 Z"/>
      <path d="M 858 339 L 869 339 L 865 329 L 868 322 L 863 311 L 860 322 L 861 330 Z"/>
    </g>
    <!-- 2014: Services auf dem Container -->
    <g transform="translate(880 285)">
      <rect x="-30" y="4" width="60" height="22" rx="2"/>
      <g stroke="#0A0349" stroke-width="2">
        <line x1="-20" y1="8" x2="-20" y2="22"/><line x1="-10" y1="8" x2="-10" y2="22"/><line x1="0" y1="8" x2="0" y2="22"/><line x1="10" y1="8" x2="10" y2="22"/><line x1="20" y1="8" x2="20" y2="22"/>
      </g>
      <rect x="-26" y="-14" width="15" height="15" rx="2.5"/>
      <rect x="-7"  y="-14" width="15" height="15" rx="2.5"/>
      <rect x="12"  y="-14" width="15" height="15" rx="2.5"/>
    </g>
    <!-- 2019: Modulith, ein Rahmen, innen vier Module -->
    <g transform="translate(1049 171)">
      <rect x="-26" y="-26" width="52" height="52" rx="8" fill="none" stroke="#E4DDF8" stroke-width="3.5"/>
      <rect x="-18" y="-18" width="16" height="16" rx="3"/><rect x="2" y="-18" width="16" height="16" rx="3"/>
      <rect x="-18" y="2" width="16" height="16" rx="3"/><rect x="2" y="2" width="16" height="16" rx="3"/>
    </g>
    <!-- 2022: SCS, ein Rahmen mit eigener UI-Leiste oben und zwei Modulen -->
    <g transform="translate(1150 134)">
      <rect x="-26" y="-26" width="52" height="52" rx="8" fill="none" stroke="#E4DDF8" stroke-width="3.5"/>
      <rect x="-18" y="-18" width="36" height="8" rx="2.5"/>
      <rect x="-18" y="-6" width="16" height="24" rx="3"/><rect x="2" y="-6" width="16" height="24" rx="3"/>
    </g>
  </g>
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
