<!-- .slide: data-background-image="./assets/bulkhead.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.40" data-background-repeat="no-repeat" -->

## Bulkhead

<p class="subtitle">Schotten im Schiff</p>

Note:
- Hook: &bdquo;Story 3 sch&uuml;tzt gegen <em>kaputte</em> Backends. Aber was, wenn das Backend gar nicht kaputt ist &mdash; nur langsam? Hotel antwortet in 2&nbsp;s, fehlerfrei. Der Circuit Breaker bleibt CLOSED. Trotzdem h&auml;ngen unter Last alle Threads in Hotel-Calls fest, Flight und Car bekommen nichts mehr durch.&ldquo;
- Provokation: &bdquo;Async-Aufrufe l&ouml;sen das nicht &mdash; eine Goroutine, die an HTTP h&auml;ngt, ist genauso teuer wie ein blockierter Thread.&ldquo;
- &Uuml;bergang zur Karten-Slide: &bdquo;F&uuml;nf Bausteine, die jedem Backend nur seinen eigenen Schoss geben.&ldquo;
