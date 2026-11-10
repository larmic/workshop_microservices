<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Circuit Breaker</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Wer ein totes Backend weiter anruft, stirbt <span class="accent">mit</span>.</p>
<p class="sub">Nach f&uuml;nf Fehlern in Folge macht der Aufrufer dicht. Schnell scheitern statt im Timeout h&auml;ngen.</p>
<figure class="foot-quote">
<blockquote>Der Circuit Breaker ist das Pattern, mit dem ein Aufrufer sich vor einem kranken Partner sch&uuml;tzt.</blockquote>
<figcaption>Michael Nygard, <em>Release It!</em>, 2007 (sinngem&auml;&szlig;)</figcaption>
</figure>
</div>
<div class="split-figure">
<img src="./assets/circuitbreaker.svg" alt="Sicherungskasten mit drei Schutzschaltern f&uuml;r Flight, Hotel und Car, der von Flight ist ausgel&ouml;st"/>
</div>
</div>

</div>


Note:
- <strong>Frage in die Runde:</strong> &bdquo;Wer hatte schon mal einen Threadpool voller h&auml;ngender Calls?&ldquo;
- <strong>Monolith-Check:</strong> Ein Methodenaufruf h&auml;ngt nicht am Netz. &Uuml;ber das Netz bindet jeder Call bis zum Timeout einen Thread. Das entsteht erst durch die Verteilung.
- <strong>Pro Backend ein Schalter:</strong> Flight ist raus, Hotel und Car laufen weiter.
- Die f&uuml;nf Fehler sind unsere Referenz-Einstellung, keine Norm.
- <strong>Weiter:</strong> Wie der Schalter genau arbeitet.
