<!-- .slide: data-background-image="./assets/tracing.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Distributed Tracing

<p class="subtitle">Den roten Faden im Log</p>

<div class="factor-row">

<div class="factor fragment">
<h3>Trace</h3>
<p>Der <em>gesamte</em> Gesch&auml;ftsvorgang &mdash; eine ID f&uuml;r alle Services entlang einer Anfrage.</p>
<code>Trace-ID: 16 Byte, 32 hex</code>
<aside class="notes">Analogie: Sendungsverfolgung &uuml;ber mehrere Lieferdienste. Ohne gemeinsame Sendungsnummer wei&szlig; am Ende niemand, wo das Paket steckt. Software-Pendant: ein Buchungsvorgang ber&uuml;hrt Booking + Flight + Hotel + Car &mdash; ohne Trace-ID kann niemand sagen, was zu welcher Anfrage geh&ouml;rt. Die Trace-ID bleibt &uuml;ber <em>alle</em> Hops gleich; das ist genau der Punkt.</aside>
</div>

<div class="factor fragment">
<h3>Span</h3>
<p>Ein einzelner Arbeitsabschnitt im Trace &mdash; ein HTTP-Call, ein DB-Query, eine Funktionsausf&uuml;hrung.</p>
<code>Span-ID: 8 Byte, 16 hex</code>
<aside class="notes">Jede Span hat eine eigene Span-ID und eine <strong>Parent-Span-ID</strong>, die auf die Span verweist, in deren Kontext sie entstanden ist. Daraus entsteht ein Baum: oben die eingehende Anfrage am Booking-Service, darunter pro Outbound-Call eine eigene Span. Jaeger / Tempo / Datadog visualisieren das als &bdquo;Wasserfall&ldquo; mit Dauer pro Span. Im Workshop bauen wir bewusst nur Trace-IDs, keinen vollwertigen Span-Baum &mdash; daf&uuml;r reichen 60 Min nicht. Die Span-IDs sind im Workshop-Code da, dienen aber nur dazu, das <code>traceparent</code>-Format g&uuml;ltig zu halten.</aside>
</div>

<div class="factor fragment">
<h3>Hop</h3>
<p>Ein Service-zu-Service-&Uuml;bergang. Pro Hop <span class="hl">neue Span-ID</span>, <em>gleiche</em> Trace-ID.</p>
<code>Booking &rarr; Flight = 1 Hop</code>
<aside class="notes">Der Hop ist der Mechanismus, der den Span-Baum aufspannt. Booking ruft Flight &rarr; pro Hop wird eine frische Span-ID generiert, die alte wandert als Parent-Span-ID in den Header. Trace-ID bleibt konstant &mdash; sie ist der gemeinsame Identifier des gesamten Vorgangs. So entstehen pro Buchung typischerweise 4&ndash;6 Spans (eingehende Booking-Anfrage + Outbound-Calls zu Flight / Hotel / Car + ggf. Kompensation).</aside>
</div>

<div class="factor fragment">
<h3>traceparent</h3>
<p>Der W3C-Header. Sprach- und vendor-neutral, fest <span class="hl">55 Zeichen</span>.</p>
<code>00-&lt;trace 32&gt;-&lt;span 16&gt;-&lt;flags 2&gt;</code>
<aside class="notes">Bis ca. 2020 hatte jedes Tracing-&Ouml;kosystem sein eigenes Header-Format (Zipkin <code>X-B3-*</code>, Jaeger <code>uber-trace-id</code>, Datadog <code>x-datadog-*</code>). W3C Trace Context vereint das, Default in OpenTelemetry und allen modernen SDKs. Vier Felder mit festen L&auml;ngen: <strong>version</strong> (2 hex, derzeit immer <code>00</code>), <strong>trace-id</strong> (32 hex, &uuml;berall gleich), <strong>parent-id</strong> = Span-ID (16 hex, pro Hop neu), <strong>flags</strong> (2 hex, Bitmaske). Stolperstein: Nur-Nullen-Trace oder Nur-Nullen-Span sind ung&uuml;ltig &mdash; <code>math/rand</code> mit Default-Seed reicht nicht, immer <code>crypto/rand</code>.</aside>
</div>

<div class="factor fragment">
<h3>Sampling-Flag</h3>
<p>Bit&nbsp;0 in <code>flags</code>: &bdquo;soll dieser Trace <em>exportiert</em> werden?&ldquo; Sonst zu teuer.</p>
<code>flags = 01 &harr; sampled</code>
<aside class="notes">In Produktion sind Traces teuer (Storage, Netzwerk, Backend-Lizenzen) &mdash; niemand exportiert 100&nbsp;%. Das Flag steuert die Entscheidung pro Trace, propagiert sich &uuml;ber alle Hops: wenn der Entry-Point &bdquo;sample&ldquo; sagt, sehen alle nachgelagerten Services dieses Bit und exportieren ihre Spans mit. Sonst werden Spans lokal erzeugt, aber nicht ans Backend gesendet. Strategien (Head-based / Tail-based / Adaptive) kommen im Recap.</aside>
</div>

</div>

<div class="market-row">

### Am Markt

<div class="chip-row">
  <span class="chip brand">OpenTelemetry</span>
  <span class="chip">Jaeger</span>
  <span class="chip">Grafana Tempo</span>
  <span class="chip">Zipkin</span>
  <span class="chip">Datadog APM</span>
  <span class="chip">Honeycomb</span>
  <span class="chip">Dynatrace</span>
  <span class="chip">AWS X-Ray</span>
  <span class="chip">Micrometer Tracing</span>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>

Note:
- Hook: &bdquo;In Story 5 und 6 hattet ihr eine Buchung, die durch vier Service-Logs gewandert ist. Wer konnte einen einzelnen Vorgang sauber rekonstruieren?&ldquo; Antwort meistens: Timestamps zusammenpuzzeln, viel Augenma&szlig;. Dann der Effekt sp&auml;ter in der Demo: <code>docker compose logs | grep &lt;trace-id&gt;</code> &mdash; alles in einem Block.
- Karten-Reihenfolge bewusst: erst die <em>Begriffe</em>, dann das <em>Wire-Format</em>, zuletzt die <em>Cost-Frage</em>. Mechanik (Propagation, Logging, Async-Grenze) kommt auf der n&auml;chsten Folie als Pseudo-Code.
- Trace vs. Span explizit trennen: <em>Trace</em> ist der &bdquo;rote Faden&ldquo;, <em>Span</em> ist ein einzelnes Wegst&uuml;ck, <em>Hop</em> ist der &Uuml;bergang. Wir bauen im Workshop nur Trace-Korrelation &mdash; den Span-Baum baut OpenTelemetry, daf&uuml;r ist die n&auml;chste Folie der &Uuml;bergang.
- Wer was wo nutzt: <strong>OpenTelemetry</strong> ist de-facto-Standard, vendor-neutraler OTLP-Export. <strong>Jaeger</strong> klassisch lokal/CNCF, Storage skaliert nicht trivial. <strong>Grafana Tempo</strong> g&uuml;nstig (Objekt-Storage), integriert mit Loki/Mimir. SaaS-L&ouml;sungen (Datadog, Honeycomb, Dynatrace, New Relic) komfortabel, aber Pricing-Themen.
- Wichtiger Workshop-Hinweis: <strong>Wir bauen die Trace-Propagation selbst</strong> (ca. 100 Zeilen), <em>kein</em> OpenTelemetry. Der Mechanismus ist sprach- und vendor-neutral; OTel-SDKs sind pro Sprache unterschiedlich und w&uuml;rden den Lerneffekt verstecken. In Produktion nat&uuml;rlich umgekehrt.
- &Uuml;berleitung: Jetzt schauen wir konkret auf den Pseudo-Code, der diese Konzepte umsetzt &mdash; Middleware, Inject, Logging, Async-Grenze.
