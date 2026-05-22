<!-- .slide: data-background-image="./assets/zusammenfassung.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Die Reise &mdash; sieben Stories

<p class="subtitle">Was bleibt</p>

| Story | | Take-away |
|---|---|---|
| **1** | Fundament | Health + Config + Stateless. Der Rest baut darauf. |
| **2** | Service Discovery | Logische Namen statt URLs. In K8s oft redundant. |
| **3** | Circuit Breaker | Schnell scheitern statt im Timeout h&auml;ngen. |
| **4** | Bulkhead | Pool pro Downstream. Async &ne; Bulkhead. |
| **5** | Saga | Kompensation muss letztlich gelingen. |
| **6** | Choreography | Eventing eliminiert nichts &mdash; es verschiebt. |
| **7** | Tracing | Trace-ID nur am Entry-Point. Span pro Hop. |

<div class="box">

### Wenn ihr drei Dinge mitnehmt

- **Microservices sind ein Werkzeug, kein Ziel.** Conway's Law ist der einzige zwingend gute Grund.
- **Resilience im Aufrufer, Schutz im Aufgerufenen.** Wer das vermischt, sch&uuml;tzt nichts.
- **Eventing eliminiert keine Komplexit&auml;t &mdash; es verschiebt sie.** Wer Choreography ohne durable Messaging baut, baut sich einen schlechten Broker.

</div>

Note:
- Die Tabelle ist der &bdquo;rote Faden&ldquo;: eine Zeile pro Story, eine Zeile Take-away. Wenn die Teilnehmer das Bild in zwei Wochen noch vor Augen haben, war der Workshop wirkungsvoll.
- Die drei Box-Punkte sind die Quintessenz &mdash; das, was in Architektur-Reviews z&auml;hlt.
- Provokation als Schlusspunkt: &bdquo;Welche dieser sieben w&uuml;rdet ihr in eurem Projekt sofort einf&uuml;hren &mdash; welche nicht, weil ihr sie nicht braucht?&ldquo;
- Optional vorlesen: &bdquo;Im Workshop habt ihr gesehen, dass jedes Pattern eine konkrete Antwort auf ein konkretes Problem ist. Das Anti-Pattern ist nicht &bdquo;wir benutzen das falsche&ldquo; &mdash; sondern &bdquo;wir benutzen alles, weil's modern ist&ldquo;.&ldquo;
