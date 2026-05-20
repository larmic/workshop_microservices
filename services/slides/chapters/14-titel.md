<!-- .slide: data-background-image="./assets/circuitbreaker.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.40" data-background-repeat="no-repeat" -->

## Circuit Breaker

<p class="subtitle">Wenn der Flug ausf&auml;llt</p>

Note:
- Hook: &bdquo;In Story 2 haben wir gelernt, Services zu <em>finden</em>. Heute kl&auml;ren wir, was passiert, wenn wir einen gefunden haben &mdash; und er antwortet nicht.&ldquo;
- Demo-Vorschau: Im Dashboard Flight auf &bdquo;Fehler&ldquo; stellen, dann ein paar Requests gegen <code>/booking/offers</code> &mdash; jeder h&auml;ngt 3 s im Timeout. Mit Circuit Breaker in Story 3: nach den ersten f&uuml;nf Fehlern instant Fallback.
- &Uuml;bergang zur Karten-Slide: &bdquo;F&uuml;nf Bausteine, die den Kollaps verhindern.&ldquo;
