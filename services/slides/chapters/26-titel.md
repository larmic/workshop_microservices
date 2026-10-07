<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Distributed Tracing</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Vier Logs, eine Buchung. Wer findet den <span class="accent">roten Faden</span>?</p>
<p class="sub">Eine Trace-ID pro Vorgang, in jeder Logzeile, &uuml;ber jede Service-Grenze. Dann reicht ein <code>grep</code>.</p>
<p class="source">W3C Trace Context (2020): ein Header, 55 Zeichen, vendor-neutral. Vorher hatte jedes Tool sein eigenes Format.</p>
</div>
<div class="split-figure">
<img src="./assets/tracing.svg" alt="Vier Service-Logs nebeneinander, eine Linie zieht sich durch alle vier"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;In Story 6 und 7 hattet ihr eine Buchung, die durch vier Service-Logs gewandert ist. Wer von euch konnte einen einzelnen Vorgang sauber rekonstruieren? Eben: Timestamps und Augenma&szlig;.&ldquo;
- Die Grafik: vier Logs, Booking, Flight, Hotel, Car. Die Lime-Linie ist der Faden: in jedem Log genau die Zeile, die zu dieser einen Buchung geh&ouml;rt.
- W3C Trace Context: Bis etwa 2020 hatte jedes Tracing-&Ouml;kosystem sein eigenes Header-Format (Zipkin <code>X-B3-*</code>, Jaeger <code>uber-trace-id</code>, Datadog <code>x-datadog-*</code>). Der Standard vereint das, heute Default in OpenTelemetry und allen modernen SDKs.
- Demo-Vorschau: Dieselbe Buchung mit Trace-ID, <code>docker compose logs | grep &lt;trace-id&gt;</code> zeigt den ganzen Vorgang in einem Block, &uuml;ber alle vier Services hinweg.
