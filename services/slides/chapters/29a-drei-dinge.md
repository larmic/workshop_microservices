<div class="page">

<p class="kicker">Wenn ihr nur drei Dinge mitnehmt</p>

## Drei S&auml;tze

<div class="page-body">

<div class="numlist">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<h3>Microservices sind ein <span class="hl">Werkzeug</span>, kein Ziel.</h3>
<p>Unabh&auml;ngige Teams sind der einzige zwingend gute Grund. Alles andere geht auch im Modulith.</p>
</div>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<h3>Resilience im Aufrufer, Schutz im <span class="hl">Aufgerufenen</span>.</h3>
<p>Circuit Breaker und Timeout geh&ouml;ren nach au&szlig;en, Bulkhead und Rate Limit nach innen. Wer das vermischt, sch&uuml;tzt nichts.</p>
</div>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<h3>Eventing eliminiert keine Komplexit&auml;t. Es <span class="hl">verschiebt</span> sie.</h3>
<p>Wer Choreography ohne durable Messaging baut, baut sich einen schlechten Broker.</p>
</div>
</div>
</div>

</div>

</div>

Note:
- Die drei S&auml;tze sind die Quintessenz, das, was in Architektur-Reviews z&auml;hlt. Sie spiegeln die zwei Fragen vom Anfang: Was m&uuml;ssen wir tun, und was kostet es?
- <strong>Werkzeug:</strong> R&uuml;ckgriff auf &bdquo;Brauchen wir Microservices?&ldquo; vom Anfang. TN Skills ist vom Microservice-Setup zur&uuml;ck zum Modulith ger&uuml;ckt: gleiches Produkt, 5 statt 14 Leute, h&ouml;here Velocity. Das ist kein Versagen, das ist eine gesunde Architektur-Entscheidung.
- <strong>Aufrufer und Aufgerufener:</strong> Der h&auml;ufigste Denkfehler aus Story 4 und 5. Der Circuit Breaker ist immer ausgehend, der Bulkhead sch&uuml;tzt den eigenen Pool.
- <strong>Verschieben:</strong> Story 7 im Rollenspiel: Der Tisch als Blatt Papier verliert Zettel, das Postfach mit Quittung liefert doppelt. Beides muss jemand behandeln.
- Diskussion Kulturwandel aus <code>docs/themen.md</code>, falls Zeit ist: Microservices ohne DevOps? Spoiler: nein. &bdquo;You build it, you run it&ldquo;, was hei&szlig;t das konkret f&uuml;r Bereitschaft? Muss man Cloud Native sein? Cloud Native ist ein Toolkit, kein Zwang, aber on-prem ohne diese Toolchain ist deutlich teurer.
- Grenzfrage: Monolith, Modulith, Microservices. Wer f&uuml;hlt sich heute auf welcher Stufe zu Hause? Was w&auml;re der Trigger, eine Stufe weiterzugehen?
