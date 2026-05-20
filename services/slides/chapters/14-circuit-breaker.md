<!-- .slide: data-background-image="./assets/circuitbreaker.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Circuit Breaker

<p class="subtitle">Wenn der Flug ausf&auml;llt</p>

<div class="factor-row">

<div class="factor fragment">
<h3>Schutzschalter</h3>
<p>Wie die Sicherung im Stromkasten &mdash; bei zu vielen Fehlern <em>raus</em>.</p>
<aside class="notes">Analogie zum Hardware-Sicherungsautomaten. Statt zu hoffen, dass der nachgelagerte Service sich erholt, kappt der Aufrufer den Stromkreis selbst. Das schont beide Seiten: das kaputte Backend bekommt Luft, der Aufrufer h&auml;ngt nicht mehr in Timeouts.</aside>
</div>

<div class="factor fragment">
<h3>Drei Zust&auml;nde</h3>
<p><code>CLOSED</code> &rarr; <code>OPEN</code> &rarr; <code>HALF_OPEN</code>. Wechsel anhand Fehlerz&auml;hler.</p>
<aside class="notes">CLOSED = alles l&auml;uft, durchwinken. OPEN = sofort Fallback, keine Anfrage geht raus. HALF_OPEN = einzelner Probe-Call testet, ob das Backend wieder lebt. Erfolg &rarr; CLOSED. Fehler &rarr; zur&uuml;ck nach OPEN.</aside>
</div>

<div class="factor fragment">
<h3>Schnell scheitern</h3>
<p>Statt 30 s Timeout sofort <span class="hl">Fallback</span>. Caller bleibt reaktiv.</p>
<aside class="notes">Cascading-Failures verhindern: ohne CB stauen sich Requests, Threads / Goroutines / Connections gehen aus, der Aufrufer wird selbst langsam und reisst seine Aufrufer mit. CB bricht die Kette &mdash; der Service bleibt schnell, auch wenn das Backend tot ist.</aside>
</div>

<div class="factor fragment">
<h3>Fallback-Strategie</h3>
<p>Leere Liste, Cache, anderer Provider &mdash; Hauptsache <span class="hl">etwas</span>.</p>
<aside class="notes">Der CB selbst entscheidet nicht, was passieren soll, wenn er feuert. Das ist Fachlogik: Bei <code>GET /booking/offers</code> reicht oft ein <code>flights:[]</code> mit Hinweis. Bei Schreib-Operationen knifflig &mdash; siehe Saga (Story 5).</aside>
</div>

<div class="factor fragment">
<h3>Outbound, nicht Inbound</h3>
<p>CB sch&uuml;tzt den <em>Aufrufer</em>, nicht den Aufgerufenen.</p>
<aside class="notes">Klassischer Denkfehler: &bdquo;Wir setzen einen CB vor unseren Service, damit er nicht &uuml;berlastet.&ldquo; Falsch &mdash; daf&uuml;r gibt's Bulkhead / Rate-Limiting. Der CB ist immer ausgehend: ich sch&uuml;tze mich davor, an einem anderen Service kaputtzugehen.</aside>
</div>

</div>

<div class="market-row">

### Am Markt

<div class="chip-row">
  <span class="chip brand">Resilience4j</span>
  <span class="chip">Polly (.NET)</span>
  <span class="chip">Spring Cloud Circuit Breaker</span>
  <span class="chip">Netflix Hystrix</span>
  <span class="chip">MicroProfile Fault Tolerance</span>
  <span class="chip">gobreaker</span>
  <span class="chip">opossum (Node)</span>
  <span class="chip">Envoy / Istio</span>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>

Note:
- Hook: &bdquo;In Story 2 haben wir die Services gefunden. Was passiert, wenn einer von ihnen kaputt ist?&ldquo; Demo-Einstieg: Flight im Dashboard auf &bdquo;Fehler&ldquo; stellen, ohne CB curlen &mdash; jeder Aufruf wartet 3 s. Mit CB &mdash; nach 5 Fehlern instant.
- Karten-Reihenfolge bewusst: erst Analogie (Schutzschalter), dann Mechanik (3 Zust&auml;nde), dann Payoff (schnell scheitern), dann das h&auml;ufig vergessene St&uuml;ck (Fallback), zuletzt die Abgrenzung (Outbound).
- Wer was wo nutzt: Resilience4j ist Standard in framework-freiem Java. Spring Cloud CB ist ein Wrapper drumherum. Hystrix bewusst dabei, ist aber seit Jahren End-of-Life &mdash; trotzdem nennen, weil viele Bestandsanwendungen es noch haben. Polly im .NET-Lager analog dominant. Envoy/Istio = CB im Service Mesh, sprach-agnostisch.
- &Uuml;berleitung: Wir schauen jetzt konkret auf die Mechanik &mdash; in genau der Form, wie sie in unserem Story-3-Code steckt.
