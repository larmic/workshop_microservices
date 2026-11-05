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
<p class="numlist-label">Nur Booking erzeugt einen Trace</p>
<h3>Flight, Hotel, Car w&uuml;rfeln <span class="hl">nie</span> selbst eine ID. Warum nicht sicherheitshalber?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Entry-Point only</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">JSON statt Format-String</p>
<h3>Warum ist <code>log.Printf</code> mit eingebauter Trace-ID nicht <span class="hl">gut genug</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">trace_id als Feld</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Die Zettel-Grenze</p>
<h3>Im Story-7-Rollenspiel lag ein Zettel auf dem Tisch. Der hat keinen HTTP-<span class="hl">Header</span>. Wie kommt der Trace mit?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">traceparent als Property</code>
</div>
<div class="numlist-row">
<div class="numlist-num">05</div>
<div>
<p class="numlist-label">Hundert Zeilen reichen</p>
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
- <strong>Entry-Point:</strong> Den Trace startet die Systemgrenze, in der Referenz <code>Middleware</code>, <code>Propagate</code> &uuml;bernimmt nur. <strong>Spicy:</strong> Wer &bdquo;sicherheitshalber&ldquo; w&uuml;rfelt, bekommt tausend Mini-Traces.
- <strong>Feld statt Text:</strong> Im Format-String wird die ID vergessen. Ein strukturierter Logger h&auml;ngt <code>trace_id</code> einmal als Feld an, danach filtert <code>jq</code>.
- <strong>Zettel-Grenze:</strong> <code>traceparent</code> wandert als Property mit, bei echten Brokern im Message-Header, nicht im Payload. <strong>Spicy:</strong> Sonst rei&szlig;t der Trace an der Bus-Grenze.
- <strong>OpenTelemetry:</strong> Propagation sind hundert Zeilen. Spans, Auto-Instrumentierung, Sampling und Export baut niemand neu. In Produktion: OTel.
- Reserve: Sampling, Korrelation, Detektiv-Runde.
- Lang: <code>docs/questions/story8.md</code>
