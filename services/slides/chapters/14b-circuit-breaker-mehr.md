<!-- .slide: data-background-image="./assets/circuitbreaker.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Circuit Breaker kann mehr

<p class="subtitle">&hellip; als wir hier nutzen</p>

<div class="box">

- **Rate-basierte Ausl&ouml;sung** &mdash; Fehler<em>quote</em> &uuml;ber Sliding Window (z.&nbsp;B. &gt;&nbsp;50&nbsp;% bei min. 10&nbsp;Calls) statt fixer Z&auml;hler
- **Slow-Call Detection** &mdash; Aufrufe &uuml;ber X&nbsp;ms gelten als Fehler, auch wenn sie HTTP&nbsp;200 liefern
- **Mehrere Probes** in `HALF_OPEN` &mdash; konfigurierbar (3, 5 &hellip;) statt nur einer
- **Exception-Klassifizierung** &mdash; `recordExceptions` / `ignoreExceptions`: nur 5xx z&auml;hlt, 4xx nicht
- **Metriken &amp; Events** &mdash; Z&auml;hler pro Zustand exportieren, Listener f&uuml;r &Uuml;berg&auml;nge
- **Bulkhead-Integration** &mdash; CB + Thread-Pool-Isolation: max. N parallele Calls pro Backend
- **Kombiniert mit Retry** &mdash; Decorator-Kette: Retry &rarr; CB &rarr; Timeout &rarr; Fallback
- **Annotation-/Deklarativ** &mdash; `@CircuitBreaker(name="flight")`, kein Boilerplate
- ...

</div>

<p class="quote">Im Workshop nutzen wir bewusst <span class="hl">nur die count-basierte Variante</span>.</p>

Note:
- Diskussions-Anker: Rate-basiert vs. count-basiert. Bei 1 Req/s reicht count, bei 10k Req/s ist count gef&auml;hrlich (5 Fehler in 100 ms = OPEN, obwohl 99,99&nbsp;% der Requests gesund waren).
- Slow-Call Detection oft &uuml;bersehen: ein Backend, das &bdquo;nur&ldquo; langsam ist (P99 von 50&nbsp;ms auf 5&nbsp;s), ist genauso unbrauchbar wie ein 500er &mdash; aber unser count-Z&auml;hler sieht das nie.
- Annotation-Stil: Resilience4j + Spring Boot oder MicroProfile Fault Tolerance machen die State-Machine zu einer einzeiligen Annotation. Sch&ouml;n &mdash; aber dadurch wird leicht vergessen, dass es trotzdem eine echte Komponente mit Konfiguration ist.
- &Uuml;berleitung zu Service Mesh: All das oben gibt es auch auf der <em>Mesh-Ebene</em> (Envoy / Istio / Linkerd / Consul Connect) &mdash; ohne eine Zeile Anwendungscode. Trade-off: Komplexit&auml;t der Plattform vs. Komplexit&auml;t pro Service.
- Take-away: Was wir bauen, ist das Skelett. Produktive Implementierungen haben deutlich mehr Stellschrauben &mdash; und genau deswegen nimmt man in der Praxis selten den selbstgebauten CB.
