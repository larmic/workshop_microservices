## Die Reise &mdash; acht Stories

<p class="subtitle">Was bleibt</p>

<div class="cols summary">
<div>

| Story | | Take-away |
|---|---|---|
| <span class="story-pill">1</span> | Fundament | Health + Config + Stateless. Der Rest baut darauf. |
| <span class="story-pill">2</span> | REST vs. RESTful | Ressourcen statt Verben. Status-Codes sind Infrastruktur. |
| <span class="story-pill">3</span> | Service Discovery | Logische Namen statt URLs. In K8s oft redundant. |
| <span class="story-pill">4</span> | Circuit Breaker | Schnell scheitern statt im Timeout h&auml;ngen. |
| <span class="story-pill">5</span> | Bulkhead | Pool pro Downstream. Async &ne; Bulkhead. |
| <span class="story-pill">6</span> | Saga | Kompensation muss letztlich gelingen. |
| <span class="story-pill">7</span> | Choreography | Eventing eliminiert nichts &mdash; es verschiebt. |
| <span class="story-pill">8</span> | Tracing | Trace-ID nur am Entry-Point. Span pro Hop. |

</div>
<div>

<div class="box">

### Wenn ihr drei Dinge mitnehmt

- **Microservices sind ein Werkzeug, kein Ziel.** Conway's Law ist der einzige zwingend gute Grund.
- **Resilience im Aufrufer, Schutz im Aufgerufenen.** Wer das vermischt, sch&uuml;tzt nichts.
- **Eventing eliminiert keine Komplexit&auml;t &mdash; es verschiebt sie.** Wer Choreography ohne durable Messaging baut, baut sich einen schlechten Broker.

</div>

</div>
</div>

Note:
- Die Tabelle ist der &bdquo;rote Faden&ldquo;: eine Zeile pro Story, eine Zeile Take-away. Wenn die Teilnehmer das Bild in zwei Wochen noch vor Augen haben, war der Workshop wirkungsvoll.
- Die drei Box-Punkte sind die Quintessenz &mdash; das, was in Architektur-Reviews z&auml;hlt.
- Provokation als Schlusspunkt: &bdquo;Welche dieser acht w&uuml;rdet ihr in eurem Projekt sofort einf&uuml;hren &mdash; welche nicht, weil ihr sie nicht braucht?&ldquo;
- Optional vorlesen: &bdquo;Im Workshop habt ihr gesehen, dass jedes Pattern eine konkrete Antwort auf ein konkretes Problem ist. Das Anti-Pattern ist nicht &bdquo;wir benutzen das falsche&ldquo; &mdash; sondern &bdquo;wir benutzen alles, weil's modern ist&ldquo;.&ldquo;
