<div class="page">

<p class="kicker">Architektur &amp; Kommunikation seit 2023</p>

## Und dann kam KI

<div class="page-body">

<!-- Achtung: keine Leerzeilen innerhalb des SVG, sonst zerteilt der
     Markdown-Parser den HTML-Block und rendert die <text>-Elemente als Absatz. -->
<div class="timeline">
<svg viewBox="0 0 1180 412" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Zeitlinie ab 2010, nach links geschoben: Microservices, Modulithen, SCS, dann 2025 KI; die Straße wird holprig und endet in einem Steinbruch mit beschrifteten Felsbrocken">
  <!-- Masten -->
  <line class="tl-tick" x1="125" y1="46" x2="125" y2="317"/>
  <line class="tl-tick" x1="260" y1="46" x2="260" y2="376"/>
  <line class="tl-tick" x1="429" y1="46" x2="429" y2="171"/>
  <line class="tl-tick" x1="530" y1="46" x2="530" y2="134"/>
  <line class="tl-tick" x1="800" y1="46" x2="800" y2="376"/>
  <line class="tl-tick" x1="125" y1="317" x2="125" y2="376"/>
  <!-- Straße wie auf der Folie davor, um 620 nach links geschoben; hinter 2025 führt sie
       bergab in den Steinbruch, das Ende liegt unter den Steinen -->
  <path class="tl-road" d="M -760 191 L -620 170 C -360 130, -150 330, 110 318 C 370 306, 390 120, 560 135 C 700 150, 830 170, 905 232 C 940 262, 965 282, 990 300"/>
  <path class="tl-dash" d="M -760 191 L -620 170 C -360 130, -150 330, 110 318 C 370 306, 390 120, 560 135 C 700 150, 830 170, 905 232 C 940 262, 965 282, 990 300" stroke-dasharray="34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 34 26 14 22 20 30 8 26 16 34 6 40 10 48 4 60 2 300"/>
  <!-- Fahnen -->
  <rect class="tl-flag" x="125" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="154" y="63">2010</text>
  <rect class="tl-flag" x="260" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="289" y="63">2014</text>
  <rect class="tl-flag" x="429" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="458" y="63">2019</text>
  <rect class="tl-flag" x="530" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="559" y="63">2022</text>
  <rect class="tl-flag" x="800" y="50" width="58" height="26" rx="3"/><text class="tl-year" x="829" y="63">2025</text>
  <!-- Figuren wie auf der Folie davor -->
  <g class="tl-figures" fill="#E4DDF8" stroke="none">
    <!-- 2010: viele kleine Services, dünne Fäden -->
    <g transform="translate(125 317)">
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
    <!-- 2014: Services auf dem Container -->
    <g transform="translate(260 285)">
      <rect x="-30" y="4" width="60" height="22" rx="2"/>
      <g stroke="#0A0349" stroke-width="2">
        <line x1="-20" y1="8" x2="-20" y2="22"/><line x1="-10" y1="8" x2="-10" y2="22"/><line x1="0" y1="8" x2="0" y2="22"/><line x1="10" y1="8" x2="10" y2="22"/><line x1="20" y1="8" x2="20" y2="22"/>
      </g>
      <rect x="-26" y="-14" width="15" height="15" rx="2.5"/>
      <rect x="-7"  y="-14" width="15" height="15" rx="2.5"/>
      <rect x="12"  y="-14" width="15" height="15" rx="2.5"/>
    </g>
    <!-- 2019: Modulith, ein Rahmen, innen vier Module -->
    <g transform="translate(429 171)">
      <rect x="-26" y="-26" width="52" height="52" rx="8" fill="none" stroke="#E4DDF8" stroke-width="3.5"/>
      <rect x="-18" y="-18" width="16" height="16" rx="3"/><rect x="2" y="-18" width="16" height="16" rx="3"/>
      <rect x="-18" y="2" width="16" height="16" rx="3"/><rect x="2" y="2" width="16" height="16" rx="3"/>
    </g>
    <!-- 2022: SCS, ein Rahmen mit eigener UI-Leiste oben und zwei Modulen -->
    <g transform="translate(530 134)">
      <rect x="-26" y="-26" width="52" height="52" rx="8" fill="none" stroke="#E4DDF8" stroke-width="3.5"/>
      <rect x="-18" y="-18" width="36" height="8" rx="2.5"/>
      <rect x="-18" y="-6" width="16" height="24" rx="3"/><rect x="2" y="-6" width="16" height="24" rx="3"/>
    </g>
    <!-- Risse vom Straßenrand zwischen 2010 und 2014: am Rand breit, nach innen spitz -->
    <g fill="#ffffff" stroke="none">
      <path d="M 180 361 L 190 361 L 186 349 L 190 341 L 185 333 L 183 325 L 180 334 L 183 342 L 179 350 Z"/>
      <path d="M 217 253 L 228 253 L 224 264 L 228 272 L 222 284 L 219 273 L 221 265 Z"/>
      <path d="M 238 339 L 249 339 L 245 329 L 248 322 L 243 311 L 240 322 L 241 330 Z"/>
    </g>
  </g>
  <!-- hinter der KI-Fahne wird die Straße holprig: Risse von beiden Rändern -->
  <g fill="#ffffff" stroke="none">
      <path d="M 848 148 L 857 151 L 848 162 L 848 170 L 840 175 L 838 179 L 839 170 L 843 164 L 844 156 Z"/>
      <path d="M 833 240 L 842 244 L 844 231 L 850 227 L 848 219 L 849 215 L 844 221 L 842 227 L 837 232 Z"/>
      <path d="M 906 177 L 914 182 L 905 189 L 904 195 L 897 197 L 894 199 L 896 193 L 900 189 L 901 183 Z"/>
      <path d="M 867 259 L 875 265 L 881 251 L 888 246 L 889 236 L 892 232 L 884 239 L 879 246 L 873 251 Z"/>
  </g>
  <!-- Geröll auf der Fahrbahn -->
  <g fill="#E4DDF8" stroke="none" opacity="0.7">
    <path d="M 854 186 L 850 190 L 845 190 L 842 186 L 845 182 L 850 182 Z"/>
    <path d="M 876 205 L 874 207 L 870 207 L 868 205 L 870 203 L 874 203 Z"/>
    <path d="M 899 222 L 897 226 L 890 227 L 886 222 L 890 217 L 896 218 Z"/>
    <path d="M 920 251 L 917 254 L 913 253 L 911 250 L 913 247 L 918 247 Z"/>
    <path d="M 934 268 L 932 270 L 928 270 L 926 268 L 928 266 L 932 266 Z"/>
    <path d="M 955 279 L 952 284 L 948 284 L 945 280 L 948 276 L 953 276 Z"/>
    <path d="M 883 195 L 882 197 L 879 197 L 877 195 L 879 193 L 881 193 Z"/>
  </g>
  <!-- der Steinbruch: die Grube in Lavendel, Brocken neben der Straße, Geröll, dann die beschrifteten Felsen -->
  <g fill="#EEE7FB" stroke="none">
  <path d="M 1245 308 L 1204 362 L 1160 395 L 1076 404 L 1000 389 L 927 356 L 909 274 L 935 232 L 979 199 L 1090 167 L 1154 182 L 1227 239 Z"/>
  </g>
  <g class="tl-stones" fill="#0A0349" stroke="#ffffff" stroke-width="3" stroke-linejoin="round">
    <path d="M 880 148 L 877 154 L 872 155 L 867 152 L 867 146 L 871 144 L 877 144 Z"/>
    <path d="M 950 188 L 947 194 L 937 196 L 932 192 L 932 186 L 940 181 L 946 183 Z"/>
    <path d="M 860 261 L 857 266 L 848 267 L 844 264 L 842 259 L 848 254 L 856 255 Z"/>
    <path d="M 874 277 L 873 281 L 867 282 L 862 279 L 862 274 L 867 272 L 871 272 Z"/>
    <path d="M 922 375 L 917 381 L 909 383 L 899 378 L 899 372 L 908 367 L 918 368 Z"/>
    <path d="M 1072 392 L 1065 398 L 1056 399 L 1051 396 L 1050 388 L 1057 385 L 1068 387 Z"/>
    <path d="M 1211 381 L 1206 386 L 1197 388 L 1189 383 L 1189 377 L 1199 371 L 1207 373 Z"/>
    <path d="M 1086 150 L 1080 156 L 1074 157 L 1065 153 L 1066 147 L 1074 143 L 1082 145 Z"/>
    <path d="M 1214 236 L 1211 240 L 1203 243 L 1197 238 L 1198 232 L 1204 230 L 1210 231 Z"/>
    <path d="M 960 180 L 956 185 L 947 186 L 942 183 L 942 176 L 947 173 L 957 176 Z"/>
    <path d="M 1036 301 L 1022 323 L 992 330 L 971 322 L 955 296 L 972 279 L 995 273 L 1024 280 Z"/>
    <g>
      <path d="M 1065 341 L 1057 369 L 997 382 L 947 375 L 918 357 L 909 338 L 946 316 L 1016 309 L 1051 319 Z"/>
      <text x="990" y="344" text-anchor="middle"><tspan x="990" dy="-0.35em">Deployment</tspan><tspan x="990" dy="1.25em">l&auml;uft bei mir</tspan></text>
    </g>
    <g>
      <path d="M 1185 352 L 1174 372 L 1121 388 L 1070 378 L 1049 364 L 1047 344 L 1072 323 L 1114 317 L 1163 330 Z"/>
      <text x="1112" y="350" text-anchor="middle"><tspan x="1112" dy="-0.35em">Doku</tspan><tspan x="1112" dy="1.25em">steht im Chat</tspan></text>
    </g>
    <g>
      <path d="M 1132 287 L 1117 305 L 1061 317 L 1005 310 L 971 294 L 970 276 L 1008 254 L 1054 250 L 1112 257 Z"/>
      <text x="1052" y="282" text-anchor="middle"><tspan x="1052" dy="-0.35em">Tests</tspan><tspan x="1052" dy="1.25em">macht der Agent</tspan></text>
    </g>
    <g>
      <path d="M 1224 295 L 1208 320 L 1181 325 L 1142 323 L 1121 311 L 1120 285 L 1136 275 L 1173 270 L 1214 280 Z"/>
      <text x="1170" y="296" text-anchor="middle"><tspan x="1170" dy="-0.35em">Sprache</tspan><tspan x="1170" dy="1.25em">egal</tspan></text>
    </g>
    <g>
      <path d="M 1089 222 L 1066 246 L 1028 256 L 985 254 L 937 234 L 945 216 L 985 192 L 1022 192 L 1072 202 Z"/>
      <text x="1012" y="222" text-anchor="middle"><tspan x="1012" dy="-0.35em">Architektur</tspan><tspan x="1012" dy="1.25em">sp&auml;ter</tspan></text>
    </g>
    <g>
      <path d="M 1197 228 L 1190 246 L 1150 258 L 1089 253 L 1057 237 L 1061 215 L 1087 203 L 1133 195 L 1178 206 Z"/>
      <text x="1128" y="226" text-anchor="middle"><tspan x="1128" dy="-0.35em">Patterns</tspan><tspan x="1128" dy="1.25em">irgendwie</tspan></text>
    </g>
  </g>
  <!-- Architekturen (eine Zeile oben) -->
  <text class="tl-arch" x="125" y="36" text-anchor="middle">Microservices</text>
  <text class="tl-arch" x="429" y="36" text-anchor="middle">Modulithen</text>
  <text class="tl-arch" x="530" y="36" text-anchor="middle">SCS</text>
  <text class="tl-arch" x="800" y="36" text-anchor="middle">KI</text>
  <!-- Technologien (eine Zeile unten) -->
  <text class="tl-tech" x="125" y="399" text-anchor="middle">REST &#183; JSON</text>
  <text class="tl-tech" x="260" y="399" text-anchor="middle">Docker &#183; gRPC</text>
  <text class="tl-tech" x="800" y="399" text-anchor="middle">Vibe Coding</text>
</svg>
</div>

</div>

</div>

