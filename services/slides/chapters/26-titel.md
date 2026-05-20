<!-- .slide: data-background-image="./assets/tracing.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.40" data-background-repeat="no-repeat" -->

## Distributed Tracing

<p class="subtitle">Den roten Faden im Log</p>

Note:
- Hook: &bdquo;In Story 5 und 6 hattet ihr eine Buchung, die durch vier Service-Logs gewandert ist. Wer von euch konnte einen einzelnen Vorgang sauber rekonstruieren? Eben &mdash; Timestamps und Augenma&szlig;.&ldquo;
- Demo-Vorschau: Dieselbe Buchung mit Trace-ID &mdash; <code>docker compose logs | grep &lt;trace-id&gt;</code> zeigt den ganzen Vorgang in einem Block, &uuml;ber alle vier Services hinweg.
- &Uuml;bergang zur Karten-Slide: &bdquo;Erst die Begriffe &mdash; Trace, Span, Hop, traceparent &mdash; dann das konkrete Beispiel.&ldquo;
