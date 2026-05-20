## Story 6 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Event-POST schl&auml;gt fehl</h3>
<p>Booking publiziert per HTTP-POST. Hotel ist <span class="hl">down</span> &mdash; Connection-Refused. Was nun?</p>
<code>Event ist weg</code>
<aside class="notes"><strong>Meine Antwort:</strong> Aktuell genau <em>nichts</em>. Der Fehler wird geloggt, der Step trotzdem als <code>COMPENSATED</code> markiert, die Saga geht auf <code>FAILED</code>, der Kunde bekommt seine Antwort. Bewusst die fragilste m&ouml;gliche Variante &mdash; sie zeigt das Problem von Eventing-ohne-Broker unverstellt: <em>Booking sagt &bdquo;Event ist raus&ldquo;, Realit&auml;t: Event wurde nie empfangen.</em> Lehrbuch-Antworten: <strong>Retry mit Backoff</strong> (transiente H&auml;nger), <strong>Outbox-Pattern</strong> (Event in derselben DB-TX wie der Saga-State, Worker publiziert), <strong>Dead-Letter-Queue</strong> (Operator entscheidet manuell), <strong>at-least-once-Bus</strong> (Broker garantiert Zustellung).<br><strong>Spicy:</strong> Eventing eliminiert das &bdquo;Backend kurz weg&ldquo;-Problem nicht &mdash; es verschiebt es. In Story 5 hat Booking den Schmerz gesp&uuml;rt. In Story 6 sieht Booking gar nichts. Wer Choreography ernst meint, f&auml;ngt nicht beim Event-Versand an, sondern bei der <em>Durability des Events</em>.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Idempotenz bei Events</h3>
<p>At-least-once ist Standard. Unser Backend tut nichts au&szlig;er <span class="hl">loggen</span> &mdash; was, wenn es echten State &auml;ndert?</p>
<code>eventId als Dedup-Key</code>
<aside class="notes"><strong>Meine Antwort:</strong> Unsere Workshop-Variante ist <em>idempotent durch Zustandslosigkeit</em> &mdash; das Backend loggt nur. In Produktion ist das die Ausnahme. Doppelte Zustellung passiert bei (1) POST &rarr; 200 verloren &rarr; Sender retryt &rarr; Backend bekommt das Event zweimal, (2) Broker liefert nach Consumer-Crash erneut, (3) Outbox-Worker crasht zwischen Send und Markieren. Ohne Dedup-Logik: <em>Refund zweimal aufs Konto, Lager-Reservierung doppelt freigegeben, Inkasso-Stop mehrfach</em>. Standard-Mechanismus: <code>eventId</code> in einer kleinen Tabelle merken, Wiederholungen ignorieren &mdash; in derselben Transaktion wie die Business-Logik.<br><strong>Spicy:</strong> Idempotenz ist nicht &bdquo;nice to have&ldquo; &mdash; sie ist das, was einen Event-Handler von einem Random-Number-Generator unter Last unterscheidet. Wer Events publiziert, muss davon ausgehen, dass sie mehrfach ankommen. Wer das ignoriert, hat einen Bug, der erst unter Last zuschl&auml;gt.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Brauchen wir einen Broker?</h3>
<p>Webhooks sind einfach und kosten keine Infra. Wozu Kafka / RabbitMQ / NATS &uuml;berhaupt?</p>
<code>durable messaging</code>
<aside class="notes"><strong>Meine Antwort:</strong> F&uuml;r genau die Dinge, die unsere Variante <em>strukturell nicht leisten kann</em>: <strong>Persistenz</strong>, <strong>Redelivery</strong>, <strong>Fan-out</strong>, <strong>Backpressure</strong>, <strong>Dead-Letter</strong>, <strong>Consumer-Crash-Recovery</strong>, <strong>Ordering</strong>, <strong>Replay</strong>, <strong>Entkopplung in der Zeit</strong>. Webhook gewinnt nur in einer Disziplin: Infrastruktur-Aufwand &mdash; und genau auf diesen Punkt fallen Teams herein.<br><strong>Spicy:</strong> Choreography ohne Broker ist eine p&auml;dagogische &Uuml;bung. Die Architekturidee ist richtig &mdash; aber sie braucht durable Messaging, sonst wird sie zur Verbesserung der schlechten Art: &bdquo;Es f&uuml;hlt sich lockerer gekoppelt an, ist aber stiller im Fehlerfall.&ldquo; Wer Broker scheut, baut sich einen &mdash; meist schlechter als das fertige Werkzeug.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Wo wandert das Wissen?</h3>
<p>Saga-Logik ist jetzt verteilt &mdash; wann wird Choreography zum <span class="hl">verteilten Monolith</span>?</p>
<code>Bug-Lokalisierung</code>
<aside class="notes"><strong>Meine Antwort:</strong> In Story 5 lag das Saga-Wissen an einer Stelle (Booking). In Story 6 reagiert <em>jeder</em> Service auf Events anderer &mdash; ein Bug in der Saga-Logik kann jetzt &uuml;berall sitzen. Choreography wird zum <em>verteilten Monolith</em>, wenn (1) niemand mehr wei&szlig;, wer auf welches Event wie reagiert, (2) Event-Vertr&auml;ge nicht versioniert / dokumentiert sind, (3) eine &Auml;nderung in einem Service drei andere mitziehen muss. Gegenma&szlig;nahmen: zentraler Event-Katalog, Schema-Registry, klare Ownership pro Event-Typ, Visualisierung der Event-Fl&uuml;sse.<br><strong>Spicy:</strong> Orchestration konzentriert das Wissen, Choreography verteilt es. Beides ist legitim &mdash; aber <em>Wissen verteilen ohne Plan</em> f&uuml;hrt zum verteilten Monolith. Das ist die schlimmste beider Welten: lose gekoppelt zur Laufzeit, fest gekoppelt zur Build-Zeit.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Exactly-once?</h3>
<p>Kann man &uuml;berhaupt garantieren, dass ein Event genau einmal verarbeitet wird?</p>
<code>at-least-once + idempotent</code>
<aside class="notes"><strong>Meine Antwort:</strong> Streng genommen nein &mdash; nicht ohne globale Koordination, die in verteilten Systemen praktisch unbezahlbar ist. Die <em>praktische N&auml;herung</em> ist: <strong>at-least-once beim Sender</strong> (Broker garantiert &bdquo;mindestens einmal&ldquo;) + <strong>Idempotenz beim Empf&auml;nger</strong> (Dedup-Key, mehrfache Zustellung wird ignoriert). Effektives Ergebnis: jedes Event genau einmal wirksam, auch wenn es technisch mehrfach geschickt wird. Kafka und einige andere Broker bieten zwar Exactly-once-Semantik an &mdash; sie ist aber teuer (Transaktional-IDs, Idempotenz-Producer, viel Koordination) und l&ouml;st nur einen Teil des Problems (Producer &rarr; Broker), nicht den Empf&auml;nger-seitigen Teil.<br><strong>Spicy:</strong> &bdquo;Exactly-once&ldquo; ist Marketing. Was du wirklich willst, ist <em>at-least-once delivery + idempotent processing</em>. Das ist nicht hipper, aber tats&auml;chlich machbar.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story6.md</code>.
</aside>
