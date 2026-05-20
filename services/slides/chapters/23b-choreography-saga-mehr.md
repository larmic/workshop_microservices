## Eventing braucht einen Broker

<p class="subtitle">&hellip; was unsere Webhook-Variante nicht kann</p>

<div class="box">

- **Persistenz** &mdash; Event auf Disk, &uuml;berlebt Sender- und Empf&auml;nger-Crash
- **Redelivery / at-least-once** &mdash; Broker garantiert Zustellung; HTTP-POST tut das nicht
- **Dead-Letter-Queue** &mdash; nicht zustellbare Events landen sichtbar in einer Operator-Inbox
- **Outbox-Pattern** &mdash; Event in derselben DB-Transaktion wie der Saga-State; ein Worker publiziert
- **Fan-out per Topic** &mdash; beliebig viele Subscriber, dynamisch dazukommen lassen
- **Backpressure** &mdash; Queue puffert, langsamer Empf&auml;nger blockiert nicht den Sender
- **Replay** &mdash; Consumer springt an alten Offset, baut Read-Model neu auf
- **Ordering** &mdash; Garantie pro Partition / Topic-Sequence
- ...

</div>

<p class="quote">Choreography ohne durable Messaging ist eine <span class="hl">p&auml;dagogische &Uuml;bung</span> &mdash; kein Production-Pattern.</p>

Note:
- Zentrales Take-away in der Note: <em>&bdquo;Eventing eliminiert das Backend-kurz-weg-Problem nicht &mdash; es macht es stiller.&ldquo;</em> In Story 5 hat Booking den Schmerz gef&uuml;hlt und konnte reagieren. In Story 6 sieht Booking gar nichts. Solange wir nichts gegen verlorene Events haben, ist das das Gegenteil von resilient.
- Das Antipattern explizit benennen: <pre>&bdquo;Wir wollen Eventing, aber keinen Broker betreiben.&ldquo;
&rarr; 6 Monate sp&auml;ter: Outbox-Tabelle, Retry-Worker, Dedup-Logik, DLQ-Inbox, Replay-Skript, alles selbst gebaut.
&rarr; 12 Monate sp&auml;ter: das ist ein schlechter Broker.</pre>
- Outbox-Pattern erkl&auml;ren: Saga-State und Event m&uuml;ssen <em>transaktional</em> zusammen geschrieben werden, sonst gibt's Inkonsistenzen (State gespeichert, Event verloren oder umgekehrt). L&ouml;sung: Event in eine <code>outbox</code>-Tabelle in derselben DB-TX, separater Worker zieht raus und publiziert &mdash; mit Retry und at-least-once-Garantie.
- Br&uuml;cke zu Recap-Frage 3: Dort steht die Vergleichstabelle Workshop-Variante vs. echter Broker im Detail.
- Take-away: Die Architekturidee (Verantwortung verteilen, Backends selbstst&auml;ndig kompensieren) ist richtig. Aber die Idee braucht durable Messaging, sonst wird sie zur Verbesserung der schlechten Art: &bdquo;Es f&uuml;hlt sich lockerer gekoppelt an, ist aber stiller im Fehlerfall.&ldquo;
