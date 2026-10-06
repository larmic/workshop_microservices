<div class="page">

<p class="kicker">Circuit Breaker</p>

<div class="page-head">

## Der Schutzschalter

<img class="head-figure" src="./assets/circuitbreaker.svg" alt=""/>
</div>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Schutzschalter</h3>
<p>Wie die Sicherung im Stromkasten: bei zu vielen Fehlern raus.</p>
</div>
<div class="card">
<h3>Drei Zust&auml;nde</h3>
<p>Wechsel anhand Fehlerz&auml;hler. Ein Probe-Call pr&uuml;ft, ob das Backend wieder lebt.</p>
<code>CLOSED &rarr; OPEN &rarr; HALF_OPEN</code>
</div>
<div class="card">
<h3>Schnell scheitern</h3>
<p>Statt 3 s Timeout sofort Fallback. Der Aufrufer bleibt reaktiv.</p>
</div>
<div class="card">
<h3>Fallback-Strategie</h3>
<p>Leere Liste, Cache, anderer Provider. Hauptsache etwas.</p>
<code>flights: []</code>
</div>
<div class="card">
<h3>Outbound, nicht Inbound</h3>
<p>Sch&uuml;tzt den Aufrufer, nicht den Aufgerufenen.</p>
</div>
</div>

<div class="market">
<h4>Am Markt</h4>
<div class="pills">
<span class="pill brand">Resilience4j</span>
<span class="pill">Polly</span>
<span class="pill">Spring Cloud Circuit Breaker</span>
<span class="pill">gobreaker</span>
<span class="pill">Envoy / Istio</span>
</div>
</div>

</div>

</div>

Note:
- <strong>Schutzschalter:</strong> Analogie zum Sicherungsautomaten. Statt zu hoffen, dass das Backend sich erholt, kappt der Aufrufer den Stromkreis selbst. Das schont beide Seiten: Das kaputte Backend bekommt Luft, der Aufrufer h&auml;ngt nicht mehr in Timeouts.
- <strong>Drei Zust&auml;nde:</strong> CLOSED heisst alles l&auml;uft, durchwinken. OPEN heisst sofort Fallback, keine Anfrage geht raus. HALF_OPEN heisst ein einzelner Probe-Call testet, ob das Backend wieder lebt. Erfolg f&uuml;hrt zu CLOSED, Fehler zur&uuml;ck nach OPEN.
- <strong>Schnell scheitern:</strong> Ohne CB stauen sich Requests, Threads und Connections gehen aus, der Aufrufer wird selbst langsam und reisst seine Aufrufer mit. Der CB bricht die Kette.
- <strong>Fallback:</strong> Der CB entscheidet nicht, was passiert, wenn er feuert. Das ist Fachlogik. Bei <code>GET /booking/offers</code> reicht oft <code>flights: []</code> mit Hinweis. Bei Schreib-Operationen knifflig, siehe Saga in Story 6.
- <strong>Outbound:</strong> Klassischer Denkfehler: &bdquo;Wir setzen einen CB vor unseren Service, damit er nicht &uuml;berlastet.&ldquo; Daf&uuml;r gibt es Bulkhead und Rate Limiting. Der CB ist immer ausgehend.
- Am Markt: Resilience4j ist Standard im Java-Umfeld, Spring Cloud Circuit Breaker ein Wrapper darum. Polly im .NET-Lager, gobreaker f&uuml;r Go. Envoy und Istio machen dasselbe im Service Mesh, ohne Anwendungscode. Hystrix nur erw&auml;hnen, wenn jemand fragt: End-of-Life, aber in Bestandsanwendungen noch verbreitet.
- &Uuml;berleitung: &bdquo;Die Mechanik passt auf eine Folie.&ldquo;
