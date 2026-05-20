## Saga kann mehr

<p class="subtitle">&hellip; als wir hier nutzen</p>

<div class="box">

- **Idempotente Kompensation** &mdash; <code>DELETE /bookings/{id}</code> liefert 2xx auch bei unbekannter ID. Idempotenz schl&auml;gt Ehrlichkeit.
- **Retry mit Backoff** &mdash; Kompensation bei transienten Fehlern wiederholen, ohne Doppel-Storno.
- **Persistenter Saga-Log** &mdash; Crash-Recovery aus durable Store (Postgres, SQLite, Consul KV, &hellip;).
- **Workflow Engines** &mdash; Temporal, Cadence, Camunda: Saga als Code, Engine macht State + Retry + Recovery + Timeout.
- **Pivot zur fachlichen Alternative** &mdash; Gutschein statt Storno, wenn der Flug schon abgehoben ist.
- **Choreography statt Orchestration** &mdash; Wissen verteilt &uuml;ber Events (Br&uuml;cke zu Story 6).
- **Observability** &mdash; Saga-Status-Counter, Compensation-Erfolgsrate, DLQ-Tiefe, distributed Tracing mit <code>Saga-ID</code>.
- ...

</div>

<p class="quote">Im Workshop bewusst die einfachste Variante &mdash; <span class="hl">synchron, in-memory, ein Versuch</span>.</p>

Note:
- Drei zentrale Take-aways, die in der Folie nur angerissen sind, aber im Gespr&auml;ch betont werden sollten:
  <pre>Saga ohne Plan f&uuml;r gescheiterte Kompensation = optimistische Hoffnung.
Persistenz ist generelle Orchestrator-Robustheit, nicht Saga-spezifisch.
Workflow Engines drehen das Modell um: Engine l&ouml;st Recovery f&uuml;r dich.</pre>
- Idempotenz: <code>DELETE</code> auf unbekannte ID &rarr; <code>204</code>, nicht <code>404</code>. Saga-Retries d&uuml;rfen mehrfach denselben Storno absetzen &mdash; jeder Aufruf nach dem ersten w&uuml;rde bei <code>404</code> f&auml;lschlich als Fehler interpretiert.
- Workflow Engines: Statt selbst State, Retry und Recovery zu coden, schreibt man die Gesch&auml;ftslogik als Workflow. Temporal/Cadence persistieren jeden Schritt, retryen, starten nach Crash neu, garantieren exactly-once-Semantik. In Produktion ist die wichtigere Frage selten &bdquo;brauche ich eine DB?&ldquo;, sondern &bdquo;<em>schreibe ich die Orchestrator-Mechanik selbst oder nutze ich Temporal?</em>&ldquo;.
- Pivot: Refund ist eine fachliche Gegenbuchung, kein technisches Delete. Saga muss damit umgehen k&ouml;nnen, dass nicht jede Kompensation symmetrisch zum Forward-Step ist.
- Observability: Eine h&auml;ngende Saga ist nicht laut, sie ist <em>still</em> &mdash; das einzige sichtbare Signal ist ein <strong>Counter, der zu lange auf einem Wert steht</strong>. &bdquo;Wir loggen das&ldquo; ist die Antwort von Teams, die noch keine h&auml;ngende Saga in Produktion gesehen haben.
- Br&uuml;cke zu Story 6 (Choreography): Forward bleibt &auml;hnlich, Kompensation l&auml;uft asynchron &uuml;ber Events. Booking ist nach &bdquo;Event raus&ldquo; trotzdem nicht fertig &mdash; es braucht Reply-Events plus Timeout-Erkennung. Eventing eliminiert keine Komplexit&auml;t, es verschiebt sie.
