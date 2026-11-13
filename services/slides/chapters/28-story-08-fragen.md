<div class="page">

<p class="kicker">Story 8 &middot; Recap</p>

## F&uuml;nf Fragen an euch

<div class="page-body">

<div class="numlist recap compact dense">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Trace ohne Span-Baum</p>
<h3>Wir haben Trace-IDs propagiert, aber keinen Baum gebaut. Was <span class="hl">fehlt</span> uns?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">80 % f&uuml;r 20 %</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Eingehende ID &uuml;berschreiben</p>
<h3>Ein Downstream bekommt einen <code>traceparent</code> und w&uuml;rfelt trotzdem eine <span class="hl">neue</span> ID. Was passiert?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">&uuml;bernehmen, nie &uuml;berschreiben</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">Feld statt Flie&szlig;text</p>
<h3>Die Trace-ID steht im Logtext. Warum ist das nicht <span class="hl">gut genug</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">trace_id als Feld</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Die Event-Grenze</p>
<h3>Im Story-7-Spiel wurde das Storno nur zugerufen. Ein Zuruf hat keinen HTTP-<span class="hl">Header</span>. Wie kommt der Trace mit?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">traceparent als Property</code>
</div>
<div class="numlist-row">
<div class="numlist-num">05</div>
<div>
<p class="numlist-label">120 Zeilen reichen</p>
<h3>Wozu dann <span class="hl">OpenTelemetry</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">OTel = SDK + Export</code>
</div>
</div>

</div>

</div>

Note:
- Erst diskutieren, dann Daumen-Abstimmung, dann der Klick auf die Merks&auml;tze. Einstieg: &bdquo;Wo h&auml;tte euch Tracing in Stories 4 bis 7 schon geholfen?&ldquo;
- <strong>Span-Baum:</strong> Eine Span ist ein Wegst&uuml;ck mit eigener ID und Parent-ID, erst das ergibt den Wasserfall in Jaeger. Wir haben nur die Trace-ID. <strong>Spicy:</strong> Logs mit Trace-ID: 80 Prozent Nutzen f&uuml;r 20 Prozent Aufwand.
- <strong>&Uuml;bernehmen:</strong> Ohne Header darf jeder starten, wie in OTel. Die Backends lassen es weg (<code>Propagate</code>). <strong>Spicy:</strong> Wer &uuml;berschreibt, zerrei&szlig;t den Trace.
- <strong>Feld statt Text:</strong> Im Flie&szlig;text wird die ID mal vergessen, mal anders geschrieben. Ein strukturierter Logger h&auml;ngt <code>trace_id</code> einmal als Feld an, danach filtert <code>jq</code>.
- <strong>Event-Grenze:</strong> <code>traceparent</code> wandert als Property mit, bei echten Brokern im Message-Header, nicht im Payload. <strong>Spicy:</strong> Sonst rei&szlig;t der Trace an der Bus-Grenze.
- <strong>OpenTelemetry:</strong> Propagation sind gut 120 Zeilen. Spans, Auto-Instrumentierung, Sampling und Export baut niemand neu. In Produktion: OTel.
- Reserve: Sampling, Korrelation, Detektiv-Runde.
- Lang: <code>docs/questions/story8.md</code>
