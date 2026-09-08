<div class="page">

<p class="kicker">REST vs. RESTful</p>

<div class="page-head">

## Sechs Regeln

<span class="hint fragment" data-fragment-index="3">Hypermedia lassen wir bewusst weg.</span>
</div>

<p class="subtitle">&hellip; die Betriebsprobleme verhindern.</p>

<div class="page-body rules-body">

<div class="quiz-code wide"><span class="m">DELETE</span> /customers/7/bookings/4711 &nbsp;&nbsp;<span class="dim">&rarr; 204 No Content</span></div>

<div class="rules">
<div class="rules-group fragment" data-fragment-index="1">
<h4>Die Methode sagt, was passiert</h4>
<div class="rule">
<h5><span class="num">1</span> Tr&auml;gt die Semantik</h5>
<p>GET liest, POST erzeugt, PUT ersetzt, PATCH &auml;ndert, DELETE entfernt.</p>
<code>POST /users/1/delete &rarr; DELETE /users/1</code>
</div>
<div class="rule">
<h5><span class="num">2</span> Idempotenz</h5>
<p>GET, PUT und DELETE beliebig oft wiederholbar. POST nicht.</p>
<code>DELETE zweimal = derselbe Zustand</code>
</div>
<div class="rule">
<h5><span class="num">3</span> GET ohne Seiteneffekte</h5>
<p>Prefetcher, Crawler und Caches f&uuml;hren GET ungefragt aus.</p>
<code>GET /bookings/1/cancel &rarr; niemals</code>
</div>
</div>
<div class="rules-group fragment" data-fragment-index="2">
<h4>Der Pfad sagt, woran</h4>
<div class="rule">
<h5><span class="num">4</span> Ressourcen statt Verben</h5>
<p>Pfade benennen Dinge, keine Aktionen.</p>
<code>GET /getUser?id=1 &rarr; GET /users/1</code>
</div>
<div class="rule">
<h5><span class="num">5</span> Sub-Ressourcen und Filter</h5>
<p>Zugeh&ouml;rigkeit in den Pfad, Auswahl in die Query.</p>
<code>/customers/7/bookings?status=confirmed</code>
</div>
</div>
<div class="rules-group fragment" data-fragment-index="3">
<h4>Die Antwort sagt, wie es ausging</h4>
<div class="rule">
<h5><span class="num">6</span> Status-Codes statt Fehler-Body</h5>
<p>Die Infrastruktur liest nur den Code. Ein 200 mit Fehler im Body ist unsichtbar.</p>
<code>404 &middot; 409 &middot; 422 &middot; 503</code>
</div>
</div>
</div>

</div>

</div>

Note:
- Zeitbox 10 bis 12 Minuten. Die Beispielzeile oben ist der rote Faden: eine Anfrage, an der alle sechs Regeln sichtbar werden. Drei Klicks: Methode, Pfad, Antwort.
- <strong>Methode.</strong> Regel 1: Die Methode ist der Vertrag mit der Infrastruktur. Proxies, Caches und Clients entscheiden anhand der Methode: darf ich cachen, wiederholen, vorladen? Wer alles per POST macht, nimmt der Infrastruktur jede Information. Regel 2: Idempotent hei&szlig;t, mehrfach ausf&uuml;hren ergibt denselben Zustand wie einmal. Bei instabilem Netz wei&szlig; der Client nicht, ob die erste Anfrage ankam. Ein idempotenter Endpoint darf wiederholt werden, ein POST ohne Idempotency-Key bucht doppelt. Kern der Flipchart-Aufgabe, kommt in Story 6 (Kompensation) zur&uuml;ck. Regel 3: 2005, Google Web Accelerator, 37signals Backpack: L&ouml;schen-Links als GET, der Prefetcher l&ouml;schte Daten. Seitdem gilt: GET ver&auml;ndert nichts. Anekdote f&uuml;rs Quiz aufsparen.
- <strong>Pfad.</strong> Regel 4: REST ist ein Architekturstil (Fielding, 2000), RESTful hei&szlig;t im Alltag: HTTP so nutzen, wie es gemeint ist. Die meisten &bdquo;REST-APIs&ldquo; sind RPC &uuml;ber HTTP mit JSON. Wer <code>/cancelBooking</code> schreibt, hat die Ressource noch nicht gefunden. Regel 5: Hierarchie im Pfad statt <code>/getBookingsForCustomer?id=7</code>; Filter, Sortierung und Paging in Query-Parameter.
- <strong>Antwort.</strong> Regel 6: Load Balancer, Monitoring, Retry-Logik, Circuit Breaker (Story 4) und Caches sehen nur den Status-Code. <code>200 OK</code> mit <code>{ "status": "error" }</code> ist ein gr&uuml;nes Dashboard vor einem kaputten Service. 404 existiert nicht, 409 Konflikt, 422 fachlich abgelehnt, 503 Backend weg. Details in den Body.
- Hypermedia (HATEOAS) und das Richardson Maturity Model nur als Ausblick nennen: Level 2, Ressourcen plus Verben, ist das realistische Ziel f&uuml;r die meisten Teams. Nicht abbiegen, die Teilnehmenden sollen gleich selbst denken.
- &Uuml;berleitung zur Story-Folie: &bdquo;Jetzt ihr. Storno und Umbuchung, am Flipchart, in Teams.&ldquo;
- Trainer-Hinweis: <code>docs/instructions/rest-vs-restful.md</code>, Abschnitt 2.
