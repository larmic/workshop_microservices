## Bulkhead kann mehr

<p class="subtitle">&hellip; als wir hier nutzen</p>

<div class="box">

- **ThreadPool-Bulkhead** &mdash; eigener Pool inkl. echtem Timeout (Hystrix-Stil) statt nur Semaphore-Z&auml;hler
- **Wait + Timeout** &mdash; kurz warten statt sofort Fail-Fast (<code>maxWaitDuration</code>) f&uuml;r kleine Bursts
- **Adaptive Bulkhead** &mdash; Limit folgt automatisch Last und Latenz (Netflix concurrency-limits, AIMD)
- **Pool-Sizing via Little's Law** &mdash; <code>concurrency = throughput &times; latency</code>: langsames Backend braucht <em>mehr</em> Slots, nicht weniger
- **Per-Replica vs. global** &mdash; bei N Replicas sehen Backends N &times; Pool-Gr&ouml;&szlig;e; erg&auml;nzendes Limit im API-Gateway / Service Mesh
- **Backpressure-Signal** &mdash; <code>429 Too Many Requests</code> + <code>Retry-After</code> statt nur 503
- **Metriken &amp; Events** &mdash; <code>bulkhead.available.concurrent.calls</code>, Reject-Counter, Listener f&uuml;r &bdquo;Pool voll&ldquo;
- ...

</div>

<p class="quote">Im Workshop nutzen wir bewusst <span class="hl">Semaphore mit Fail-Fast</span>.</p>

Note:
- Zentrale Abgrenzung gegen Circuit Breaker &mdash; oft verwechselt, deshalb auf der Folie zweimal erw&auml;hnt:
  <pre>CB:       reagiert auf Fehler-Rate    &rarr; &bdquo;Backend ist krank&ldquo;
Bulkhead: reagiert auf Ressourcen-Druck &rarr; &bdquo;ich verbrenne mich nicht selbst&ldquo;</pre>
  Beide d&uuml;rfen gleichzeitig feuern &mdash; sie reagieren auf unterschiedliche Symptome.
- ThreadPool vs. Semaphore: ThreadPool isoliert echt (eigener Stack, eigene Threads, hartes Timeout per Thread-Interrupt). Semaphore ist leichtgewichtig (nur Z&auml;hler), kann aber Aufrufe nicht selbst abbrechen &mdash; Timeout muss <em>im Aufruf</em> selbst sitzen. Hystrix-Default war ThreadPool, Resilience4j-Default ist Semaphore.
- Little's Law am Beispiel: bei 50&nbsp;req/s und 500&nbsp;ms Latenz braucht man 25 Slots. Wer aus dem Bauch heraus &bdquo;10 klingt gut&ldquo; einstellt, hat den Bulkhead nicht implementiert, sondern dekoriert.
- Adaptive Bulkhead: Netflix concurrency-limits, AIMD &auml;hnlich TCP-Congestion-Control. Cool, aber operativ aufwendig &mdash; in den meisten F&auml;llen reicht eine gemessen feste Zahl.
- Service Mesh: Envoy regelt <code>max_pending_requests</code> und <code>max_connections</code> sprach-agnostisch im Sidecar &mdash; gleicher Effekt ohne Anwendungscode.
- Take-away: Das Skelett bleibt simpel. Die spannenden Hebel sind <em>Sizing</em> und <em>Wechselwirkung</em> mit anderen Patterns &mdash; nicht der Algorithmus selbst.
