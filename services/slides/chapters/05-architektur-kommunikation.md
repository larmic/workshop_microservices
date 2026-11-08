<div class="page">

<p class="kicker">Architektur &amp; Kommunikation seit 1990</p>

## Der Weg zu Microservices

<div class="page-body">

<!-- Achtung: keine Leerzeilen innerhalb des SVG, sonst zerteilt der
     Markdown-Parser den HTML-Block und rendert die <text>-Elemente als Absatz. -->
<div class="timeline">
<svg viewBox="0 0 1180 412" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Zeitlinie 1990 bis 2022: Straße mit Haltepunkten, daneben Gebäude: Monolithen als Hochhaus, SOA als Bürotürme mit Verbindungsbrücke, Microservices als Dorf kleiner Häuser, Docker als Containerstapel, Modulithen als Haus mit vier Abteilen, SCS als Reihenhäuser; ein Auto fährt weiter">
  <!-- Masten: von der Fahne bis unter die Straße, bei den Technologien weiter bis nach unten -->
  <line class="tl-tick" x1="70"   y1="46" x2="70"   y2="165"/>
  <line class="tl-tick" x1="408"  y1="46" x2="408"  y2="249"/>
  <line class="tl-tick" x1="745"  y1="46" x2="745"  y2="317"/>
  <line class="tl-tick" x1="1049" y1="46" x2="1049" y2="320"/>
  <line class="tl-tick" x1="1150" y1="46" x2="1150" y2="280"/>
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
  <!-- Gebäude neben der Straße, Haltepunkte auf der Straße -->
  <g fill="#ffffff"><path d="M 800 361 L 810 361 L 806 349 L 810 341 L 805 333 L 803 325 L 800 334 L 803 342 L 799 350 Z"/><path d="M 837 253 L 848 253 L 844 264 L 848 272 L 842 284 L 839 273 L 841 265 Z"/><path d="M 858 339 L 869 339 L 865 329 L 868 322 L 863 311 L 860 322 L 861 330 Z"/></g>
  <g fill="#D3F871" stroke="#ffffff" stroke-width="4"><circle cx="70" cy="165" r="10"/><circle cx="408" cy="249" r="10"/><circle cx="745" cy="317" r="10"/><circle cx="880" cy="290" r="10"/><circle cx="1049" cy="171" r="10"/><circle cx="1150" cy="134" r="10"/></g>
  <g transform="translate(70 304)"><rect x="-6" y="-97" width="12" height="40" fill="#E4DDF8"/><rect x="-42" y="-2" width="84" height="6" rx="3" fill="#E4DDF8"/><rect x="-22" y="-80" width="44" height="78" rx="2" fill="#B9A6F2"/><rect x="-4" y="-88" width="8" height="8" fill="#B9A6F2"/><g fill="#0A0349"><rect x="-15" y="-72" width="8" height="7"/><rect x="-4" y="-72" width="8" height="7"/><rect x="7" y="-72" width="8" height="7"/><rect x="-15" y="-60" width="8" height="7"/><rect x="-4" y="-60" width="8" height="7"/><rect x="7" y="-60" width="8" height="7"/><rect x="-15" y="-48" width="8" height="7"/><rect x="-4" y="-48" width="8" height="7"/><rect x="7" y="-48" width="8" height="7"/><rect x="-15" y="-36" width="8" height="7"/><rect x="-4" y="-36" width="8" height="7"/><rect x="7" y="-36" width="8" height="7"/></g><rect x="-6" y="-16" width="12" height="14" fill="#D3F871"/></g>
  <g transform="translate(408 196)"><rect x="-6" y="0" width="12" height="16" fill="#E4DDF8"/><rect x="-62" y="-2" width="124" height="6" rx="3" fill="#E4DDF8"/><rect x="-56" y="-44" width="28" height="42" rx="2" fill="#B9A6F2"/><rect x="-14" y="-54" width="28" height="52" rx="2" fill="#7348E1"/><rect x="28" y="-44" width="28" height="42" rx="2" fill="#B9A6F2"/><rect x="-56" y="-34" width="112" height="7" fill="#D3F871"/><g fill="#0A0349"><rect x="-50" y="-22" width="7" height="7"/><rect x="-40" y="-22" width="7" height="7"/><rect x="34" y="-22" width="7" height="7"/><rect x="44" y="-22" width="7" height="7"/></g><g fill="#ffffff"><rect x="-8" y="-46" width="7" height="7"/><rect x="2" y="-46" width="7" height="7"/><rect x="-8" y="-22" width="7" height="7"/><rect x="2" y="-22" width="7" height="7"/></g></g>
  <g transform="translate(745 262)"><rect x="-6" y="0" width="12" height="16" fill="#E4DDF8"/><rect x="-60" y="-2" width="120" height="6" rx="3" fill="#E4DDF8"/><g transform="translate(-44 -2)"><path d="M-10 0 V-16 L0 -25 L10 -16 V0 Z" fill="#B9A6F2"/><rect x="-3" y="-9" width="6" height="9" fill="#0A0349"/></g><g transform="translate(-22 -2)"><path d="M-10 0 V-16 L0 -25 L10 -16 V0 Z" fill="#D3F871"/><rect x="-3" y="-9" width="6" height="9" fill="#0A0349"/></g><g transform="translate(0 -2)"><path d="M-10 0 V-16 L0 -25 L10 -16 V0 Z" fill="#7348E1"/><rect x="-3" y="-9" width="6" height="9" fill="#ffffff"/></g><g transform="translate(22 -2)"><path d="M-10 0 V-16 L0 -25 L10 -16 V0 Z" fill="#B9A6F2"/><rect x="-3" y="-9" width="6" height="9" fill="#0A0349"/></g><g transform="translate(44 -2)"><path d="M-10 0 V-16 L0 -25 L10 -16 V0 Z" fill="#D3F871"/><rect x="-3" y="-9" width="6" height="9" fill="#0A0349"/></g><g transform="translate(-11 -30)"><path d="M-10 0 V-16 L0 -25 L10 -16 V0 Z" fill="#D3F871"/><rect x="-3" y="-9" width="6" height="9" fill="#0A0349"/></g><g transform="translate(11 -30)"><path d="M-10 0 V-16 L0 -25 L10 -16 V0 Z" fill="#B9A6F2"/><rect x="-3" y="-9" width="6" height="9" fill="#0A0349"/></g></g>
  <g transform="translate(880 236)"><rect x="-6" y="0" width="12" height="14" fill="#E4DDF8"/><rect x="-42" y="-2" width="84" height="6" rx="3" fill="#E4DDF8"/><rect x="-36" y="-22" width="34" height="20" rx="2" fill="#D3F871"/><rect x="0" y="-22" width="34" height="20" rx="2" fill="#B9A6F2"/><rect x="-36" y="-44" width="34" height="20" rx="2" fill="#7348E1"/><rect x="0" y="-44" width="34" height="20" rx="2" fill="#D3F871"/><rect x="-18" y="-66" width="34" height="20" rx="2" fill="#B9A6F2"/><g stroke="#0A0349" stroke-width="1.5" opacity=".45"><path d="M-30 -19 V-5 M-24 -19 V-5 M-18 -19 V-5 M-12 -19 V-5 M6 -19 V-5 M12 -19 V-5 M18 -19 V-5 M24 -19 V-5 M6 -41 V-27 M12 -41 V-27 M18 -41 V-27 M24 -41 V-27 M-12 -63 V-49 M-6 -63 V-49 M0 -63 V-49 M6 -63 V-49"/></g><g stroke="#ffffff" stroke-width="1.5" opacity=".7"><path d="M-30 -41 V-27 M-24 -41 V-27 M-18 -41 V-27 M-12 -41 V-27"/></g></g>
  <g transform="translate(1049 300)"><rect x="-6" y="-88" width="12" height="30" fill="#E4DDF8"/><rect x="-50" y="-2" width="100" height="6" rx="3" fill="#E4DDF8"/><rect x="-38" y="-54" width="76" height="52" rx="3" fill="#E4DDF8"/><path d="M-44 -54 L0 -74 L44 -54 Z" fill="#7348E1"/><rect x="-32" y="-48" width="30" height="20" rx="2" fill="#D3F871"/><rect x="2" y="-48" width="30" height="20" rx="2" fill="#B9A6F2"/><rect x="-32" y="-24" width="30" height="20" rx="2" fill="#B9A6F2"/><rect x="2" y="-24" width="30" height="20" rx="2" fill="#D3F871"/></g>
  <g transform="translate(1158 262)"><rect x="-14" y="-90" width="12" height="36" fill="#E4DDF8"/><rect x="-52" y="-2" width="104" height="6" rx="3" fill="#E4DDF8"/><g transform="translate(-32 0)"><rect x="-14" y="-40" width="28" height="38" fill="#B9A6F2"/><path d="M-17 -40 L0 -54 L17 -40 Z" fill="#7348E1"/><rect x="-14" y="-40" width="28" height="7" fill="#D3F871"/><rect x="-4" y="-15" width="8" height="13" fill="#0A0349"/></g><g transform="translate(0 0)"><rect x="-14" y="-40" width="28" height="38" fill="#E4DDF8"/><path d="M-17 -40 L0 -54 L17 -40 Z" fill="#7348E1"/><rect x="-14" y="-40" width="28" height="7" fill="#D3F871"/><rect x="-4" y="-15" width="8" height="13" fill="#0A0349"/></g><g transform="translate(32 0)"><rect x="-14" y="-40" width="28" height="38" fill="#B9A6F2"/><path d="M-17 -40 L0 -54 L17 -40 Z" fill="#7348E1"/><rect x="-14" y="-40" width="28" height="7" fill="#D3F871"/><rect x="-4" y="-15" width="8" height="13" fill="#0A0349"/></g></g>
  <g transform="translate(1222 140) rotate(4)"><path d="M-56 -6 H-36 M-60 2 H-38" stroke="#D3F871" stroke-width="3" stroke-linecap="round"/><rect x="-22" y="-12" width="44" height="14" rx="5" fill="#D3F871"/><path d="M-12 -12 L-6 -22 H8 L14 -12 Z" fill="#D3F871"/><path d="M-7 -13 L-3 -19 H6 L10 -13 Z" fill="#0A0349"/><g fill="#ffffff" stroke="#0A0349" stroke-width="3"><circle cx="-12" cy="3" r="6"/><circle cx="12" cy="3" r="6"/></g></g>
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
- <strong>Handzeichen zum Einstieg:</strong> Wer hat einen Big Ball of Mud gepflegt? Wer hat einen ESB gesehen? Wer ist schon von Microservices zur&uuml;ck zum Modulith? Die H&auml;nde zeigen, wo ihr erz&auml;hlt und wo ihr abk&uuml;rzt.
- Vorher: Mainframes
- SOA: wenige, gro&szlig;e Services
- Microservices: feingranular, eine Aufgabe, ein Service
