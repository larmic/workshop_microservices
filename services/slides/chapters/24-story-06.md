## Story 6

<p class="subtitle">Die Saga wird leise <span class="time-badge">&asymp; 60 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

In Story 5 tr&auml;gt Booking die volle Verantwortung f&uuml;r die Kompensation: synchroner <code>DELETE</code> gegen jedes Backend, eigener Retry, eigener State. F&auml;llt ein Backend aus, blockiert Booking. Fachlich ist die Stornierung aber Aufgabe des jeweiligen Backends.

Booking publiziert deshalb ein Event &mdash; die Backends abonnieren es und kompensieren selbst. Wir wechseln von <strong>Orchestration</strong> zu <strong>Choreography</strong>: derselbe fachliche Ablauf, andere Verantwortungsverteilung.

#### User Story

Als <em>System</em> m&ouml;chte ich <em>bei einer fehlgeschlagenen Buchung die Kompensation asynchron an die zust&auml;ndigen Backends delegieren</em>, damit <em>Booking nicht f&uuml;r deren Verf&uuml;gbarkeit haften muss und die Backends ihre eigene Stornierungslogik kapseln k&ouml;nnen</em>.

</div>

</div>
<div>

<div class="story-card">

#### Akzeptanzkriterien

- Booking publiziert bei Fehlschlag ein <code>CompensationRequested</code>-Event (statt <code>DELETE</code>)
- Backends nehmen das Event entgegen, antworten sofort <code>202 Accepted</code> und kompensieren asynchron
- Booking setzt den Step nach Event-Dispatch auf <code>COMPENSATED</code>, die Saga auf <code>FAILED</code>
- <strong>Idempotenz</strong>: <code>eventId</code> wird mitgesendet (persistente Speicherung ist Bonus)
- <em>Bonus</em>: Reply-Events, Timeout-Erkennung, <code>STUCK</code>-Status (siehe Recap)

</div>

</div>
</div>

Note:
- Hook: &bdquo;Wir verbessern Story 5 an genau einer Stelle &mdash; der Kompensation. Forward bleibt synchron, das ist Absicht. So sehen wir sauber, was sich verschiebt, wenn man Events einf&uuml;hrt.&ldquo;
- Wiedererkennung: dieselbe Karte (Kontext / User Story / Akzeptanzkriterien) im Dashboard unter Story 6 &rarr; &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story6/</code>.
- <strong>Wir nutzen bewusst keinen Broker</strong> &mdash; HTTP-Webhook-POST statt Kafka / RabbitMQ. Das ist die einfachste m&ouml;gliche Variante, um das Pattern sichtbar zu machen, und gleichzeitig die fragilste m&ouml;gliche Variante. Im Recap (Frage 1 und 3) wird das ausf&uuml;hrlich besprochen &mdash; in Produktion w&uuml;rde man hier einen echten Broker einf&uuml;hren.
- Demo-Drehbuch: Dashboard &rarr; Hotel auf &bdquo;Fehler&ldquo;, dann <code>POST /booking/bookings</code>. Sichtbar: Booking publiziert <code>CompensationRequested</code> f&uuml;r Flight, Flight antwortet sofort <code>202 Accepted</code> und macht den Storno asynchron (Log-Eintrag). Booking-Status: <code>COMPENSATING</code> &rarr; <code>FAILED</code>, noch bevor der Rollback im Backend abgeschlossen ist. <em>Anker f&uuml;r Recap-Frage 1</em>: der Rollback im Backend ist im Log sichtbar, aber nicht mehr in der Booking-Response. Was, wenn der Rollback fehlschl&auml;gt? Booking wei&szlig; es nicht.
- Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-06-choreography-saga.md</code>.
