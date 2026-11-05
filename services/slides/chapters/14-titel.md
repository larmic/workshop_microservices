<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Circuit Breaker</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Wer ein totes Backend weiter anruft, stirbt <span class="accent">mit</span>.</p>
<p class="sub">Nach f&uuml;nf Fehlern macht der Aufrufer selbst dicht. Schnell scheitern statt im Timeout h&auml;ngen.</p>
<p class="source">Michael Nygard, <em>Release It!</em>: Der Circuit Breaker ist das Pattern, mit dem ein Aufrufer sich vor einem kranken Partner sch&uuml;tzt.</p>
</div>
<div class="split-figure">
<img src="./assets/circuitbreaker.svg" alt="Sicherungskasten mit f&uuml;nf Schutzschaltern, der zweite ist ausgel&ouml;st"/>
</div>
</div>

</div>


Note:
- <strong>Frage in die Runde:</strong> &bdquo;Wer hatte schon mal einen Threadpool voller h&auml;ngender Calls?&ldquo; Kurz Erfahrungen sammeln.
- <strong>Monolith-Check:</strong> Ein Methodenaufruf im selben Prozess h&auml;ngt nicht am Netz. &Uuml;ber das Netzwerk kann jeder Call bis zum Timeout h&auml;ngen und bindet so lange einen Thread. Das entsteht erst durch die Verteilung.
- <strong>Quelle:</strong> Das Pattern stammt aus Michael Nygards <em>Release It!</em> (2007). Die f&uuml;nf Fehler sind unsere Referenz-Einstellung, keine Norm.
- <strong>Weiter:</strong> Wie der Schalter genau arbeitet.
