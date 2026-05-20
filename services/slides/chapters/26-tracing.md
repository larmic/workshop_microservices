## Distributed Tracing

<p class="subtitle">Den roten Faden im Log</p>

<div class="factor-row">

<div class="factor fragment">
<h3>Trace-ID</h3>
<p>Eine ID f&uuml;r den <em>gesamten</em> Vorgang &mdash; durch alle Services hindurch.</p>
<code>32 hex, einmal pro Request</code>
<aside class="notes">Analogie: Sendungsverfolgung &uuml;ber mehrere Lieferdienste. Ohne gemeinsame Sendungsnummer wei&szlig; am Ende niemand, wo das Paket steckt. In Software: ein Buchungsvorgang ber&uuml;hrt Booking + Flight + Hotel + Car &mdash; ohne Trace-ID kann niemand sagen, was zu welcher Anfrage geh&ouml;rt. Die Trace-ID ist 16 Byte zuf&auml;llig, &uuml;blicherweise als 32 hex-Zeichen.</aside>
</div>

<div class="factor fragment">
<h3>W3C Trace Context</h3>
<p>Der <code>traceparent</code>-Header &mdash; <span class="hl">55 Zeichen</span>, sprach- und vendor-neutral.</p>
<code>00-&lt;trace&gt;-&lt;span&gt;-&lt;flags&gt;</code>
<aside class="notes">Bis ca. 2020 hatte jedes Tracing-&Ouml;kosystem sein eigenes Format (Zipkin <code>X-B3-*</code>, Jaeger <code>uber-trace-id</code>, Datadog <code>x-datadog-*</code>). W3C Trace Context vereint das &mdash; Default in OpenTelemetry und allen modernen Vendor-SDKs. Felder: <code>version</code> (immer <code>00</code>) + Trace-ID (32 hex) + Span-ID pro Hop (16 hex) + Flags (Sampling-Bit). Stolperstein: Nur-Nullen ist ung&uuml;ltig; <code>math/rand</code> mit Default-Seed reicht nicht.</aside>
</div>

<div class="factor fragment">
<h3>Propagation</h3>
<p>Server <em>&uuml;bernimmt oder erzeugt</em>, Client <em>injiziert</em> in jeden Outbound-Call.</p>
<code>Middleware + Inject</code>
<aside class="notes">Zwei Middlewares, zwei Rollen. <strong>Entry-Point</strong> (Booking): erzeugt einen Trace, falls keiner reinkommt. <strong>Downstream</strong> (Flight/Hotel/Car): &uuml;bernimmt nur einen vorhandenen Trace und legt sonst nichts an. Diese Trennung sorgt daf&uuml;r, dass Downstream-Logs ihre Trace-ID nur durch aktive Propagation erhalten &mdash; nicht zuf&auml;llig vom Service selbst.</aside>
</div>

<div class="factor fragment">
<h3>Strukturiertes Logging</h3>
<p><code>trace_id</code> als <em>Feld</em>, nicht als Teil des Format-Strings.</p>
<code>JSON-Log + jq</code>
<aside class="notes">Mit <code>log.Printf("booking %s failed", id)</code> muss man die Trace-ID manuell in jeden Format-String einsetzen &mdash; das wird nach drei Wochen vergessen. Mit strukturiertem Logger (<code>slog</code>, <code>structlog</code>, <code>logback</code> JSON, <code>pino</code>) h&auml;ngt man die ID einmal am Request-Kontext-Logger und sie taucht in <em>jeder</em> Zeile auf. Konsequenz: <code>docker compose logs | jq -c 'select(.trace_id=="abc...")'</code> liefert den vollst&auml;ndigen Vorgang sauber gefiltert.</aside>
</div>

<div class="factor fragment">
<h3>Async-Grenze</h3>
<p>Bei Events: Trace-ID als <span class="hl">Property</span> mit&uuml;bertragen &mdash; sonst zerf&auml;llt der Trace.</p>
<code>traceparent im Event-Body</code>
<aside class="notes">Sobald die Kommunikation nicht mehr synchrones HTTP ist (Story 6: Compensation-Events), reicht der HTTP-Header nicht mehr. Der Trace-Kontext muss <em>aktiv</em> als Property auf das Event mitwandern, der Konsument liest und legt ihn in seinen Worker-Kontext. Stolperstein: bei Kafka/RabbitMQ/SNS-SQS in die <em>Message-Header</em>, nicht in den Payload &mdash; sonst muss jeder Consumer das Payload-Schema kennen, nur um den Trace zu propagieren.</aside>
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
- Hook: &bdquo;In Story 5 und 6 hattet ihr eine Buchung, die durch vier Service-Logs gewandert ist. Wer von euch konnte einen einzelnen Vorgang sauber rekonstruieren?&ldquo; Antwort meistens: Timestamps zusammenpuzzeln, viel Augenma&szlig;. Dann der Effekt: <code>docker compose logs | grep &lt;trace-id&gt;</code> &mdash; alles in einem Block.
- Karten-Reihenfolge bewusst: erst das <em>Was</em> (Trace-ID), dann das <em>Wie auf der Leitung</em> (W3C, Propagation), dann das <em>Wo es sichtbar wird</em> (Logs), zuletzt die <em>Stelle, wo es bricht</em> (Async-Grenze).
- Wer was wo nutzt: <strong>OpenTelemetry</strong> ist de-facto-Standard, vendor-neutraler OTLP-Export. <strong>Jaeger</strong> klassisch lokal/CNCF, Storage skaliert nicht trivial. <strong>Grafana Tempo</strong> g&uuml;nstig (Objekt-Storage), integriert mit Loki/Mimir. SaaS-L&ouml;sungen (Datadog, Honeycomb, Dynatrace, New Relic) komfortabel, aber Pricing-Themen.
- Wichtiger Workshop-Hinweis: <strong>Wir bauen die Trace-Propagation selbst</strong> (ca. 100 Zeilen), <em>kein</em> OpenTelemetry. Der Mechanismus ist sprach- und vendor-neutral; OTel-SDKs sind pro Sprache unterschiedlich und w&uuml;rden den Lerneffekt verstecken. In Produktion nat&uuml;rlich umgekehrt.
- &Uuml;berleitung: Wir schauen uns die Mechanik konkret an &mdash; in genau der Form, wie sie im Dashboard-Spickzettel und im Story-7-Code steckt.
