## REST vs. RESTful

<p class="subtitle">Sechs Regeln, die Betriebsprobleme verhindern</p>

<div class="factor-row six">

<div class="factor fragment">
<h3><span class="numeral">1</span> Ressourcen statt Verben</h3>
<p>Pfade benennen <em>Dinge</em>, keine Aktionen. Das Verb steckt in der Methode.</p>
<code>GET /getUser?id=1423 &rarr; GET /users/1423</code>
<aside class="notes">REST ist ein Architekturstil (Fielding, 2000). RESTful hei&szlig;t im Alltag: HTTP so nutzen, wie es gemeint ist. Die meisten &bdquo;REST-APIs&ldquo; sind RPC &uuml;ber HTTP mit JSON. Erste Regel: Substantive im Pfad. Wer <code>/cancelBooking</code> schreibt, hat die Ressource noch nicht gefunden.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Methode tr&auml;gt Semantik</h3>
<p>GET liest, POST erzeugt, PUT ersetzt, PATCH &auml;ndert teilweise, DELETE entfernt.</p>
<code>POST /users/1423/delete &rarr; DELETE /users/1423</code>
<aside class="notes">Die Methode ist der Vertrag mit der Infrastruktur. Proxies, Caches und Clients treffen anhand der Methode Entscheidungen: darf ich cachen, darf ich wiederholen, darf ich vorladen? Wer alles per POST macht, nimmt der Infrastruktur jede Information.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Idempotenz</h3>
<p>GET, PUT, DELETE beliebig oft wiederholbar. POST nicht. Das macht <span class="hl">Retries sicher</span>.</p>
<code>DELETE zweimal = derselbe Zustand</code>
<aside class="notes">Idempotent hei&szlig;t: mehrfach ausf&uuml;hren ergibt denselben Zustand wie einmal. Bei instabilem Netz wei&szlig; der Client nicht, ob die erste Anfrage angekommen ist. Ein idempotenter Endpoint darf wiederholt werden, ein POST ohne Idempotency-Key bucht doppelt. Das ist der Kern der Flipchart-Aufgabe und kommt in Story 6 (Kompensation) zur&uuml;ck.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> GET ohne Seiteneffekte</h3>
<p>Prefetcher, Crawler und Caches d&uuml;rfen GET jederzeit ausf&uuml;hren.</p>
<code>GET /bookings/1423/cancel &rarr; niemals</code>
<aside class="notes">2005, Google Web Accelerator: ein Browser-Plugin lud Links vor. Bei 37signals Backpack waren &bdquo;L&ouml;schen&ldquo;-Links einfache GET-Links. Nutzern wurden Daten gel&ouml;scht, ohne dass sie geklickt hatten. Seitdem gilt: GET ver&auml;ndert nichts. Nicht als Stilfrage, sondern als Selbstschutz. Kommt im Quiz wieder.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Status-Codes statt Fehler-Body</h3>
<p>Die Infrastruktur liest <em>nur</em> den Code. <code>200</code> mit Fehler im Body ist unsichtbar.</p>
<code>404 &middot; 409 &middot; 422 &middot; 503</code>
<aside class="notes">Load Balancer, Monitoring, Retry-Logik, Circuit Breaker (Story 4) und Caches sehen nur den Status-Code. Wer <code>200 OK</code> mit <code>{ "status": "error" }</code> liefert, hat ein gr&uuml;nes Dashboard vor einem kaputten Service. Richtig: 404 existiert nicht, 409 Konflikt (schon storniert), 422 fachlich abgelehnt, 503 Backend weg. Details in den Body.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Sub-Ressourcen und Filter</h3>
<p>Zugeh&ouml;rigkeit im Pfad, Auswahl in Query-Parametern. Keine Sonder-Endpoints.</p>
<code>GET /customers/7/bookings?status=confirmed</code>
<aside class="notes">Statt <code>/getBookingsForCustomer?id=7</code> die Hierarchie im Pfad abbilden. Filter, Sortierung und Paging geh&ouml;ren in Query-Parameter. Nur Ausblick, nicht vertiefen: HATEOAS (Links in Antworten) und das Richardson Maturity Model (Fowler): Level 2, also Ressourcen plus Verben, ist das realistische Ziel f&uuml;r die meisten Teams.</aside>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>

Note:
- Zeitbox 10 bis 12 Minuten. Nicht in HATEOAS-Diskussionen abbiegen, die Teilnehmenden sollen gleich selbst denken.
- Pro Karte ein Satz und das Gegenbeispiel. Die Anekdoten (Google Web Accelerator, Stripe) f&uuml;r das Quiz aufsparen.
- &Uuml;berleitung: &bdquo;Auf der n&auml;chsten Folie einmal alles nebeneinander, dann seid ihr dran.&ldquo;
- Trainer-Hinweis: <code>docs/instructions/rest-vs-restful.md</code>, Abschnitt 2.
