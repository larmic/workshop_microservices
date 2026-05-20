## Story 5 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Kompensation scheitert</h3>
<p>Flug gebucht, Hotel sagt nein &mdash; und der Storno-Call l&auml;uft selbst in <span class="hl">5xx</span>. Was nun?</p>
<code>must eventually succeed</code>
<aside class="notes"><strong>Meine Antwort:</strong> Saga macht eine starke Annahme: Forward-Steps d&uuml;rfen scheitern, <em>Kompensationen m&uuml;ssen letztlich gelingen</em>. Ohne diese Annahme bricht das Konstrukt zusammen. Lehrbuch-Strategien: <strong>Idempotenz</strong> (mehrfacher DELETE &rarr; gleiches Ergebnis), <strong>Retry mit Backoff</strong>, <strong>persistenter Saga-Log</strong> (Crash darf keine offene Kompensation verlieren), <strong>Dead-Letter / Operator-Inbox</strong> (Mensch greift ein), <strong>Pivot zur fachlichen Alternative</strong> (Gutschein statt Storno), <strong>Compensation-by-design</strong> (Reservierung als Status, nicht als L&ouml;schung).<br><strong>Spicy:</strong> Eine Saga ohne Plan f&uuml;r gescheiterte Kompensation ist keine Saga, sondern eine optimistische Hoffnung. Die Frage &bdquo;was bei Misserfolg&ldquo; ist nicht optional &mdash; sie ist das eigentliche Engineering an dem Pattern.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Kein Retry &mdash; schlimm?</h3>
<p>Wir kompensieren <span class="hl">genau einmal</span>. Bei Fehlschlag <code>FAILED</code>. Verletzen wir damit das Pattern?</p>
<code>Monitoring = Retry</code>
<aside class="notes"><strong>Meine Antwort:</strong> Ja, formal verletzen wir das Pattern. Pragmatisch ist die Entscheidung trotzdem vertretbar &mdash; <em>wenn</em> bewusst getroffen. Was zwingend dazugeh&ouml;rt: (1) <strong>Saga-Status persistieren</strong> (PENDING / COMPENSATING / FAILED), sonst wei&szlig; niemand, dass eine Saga h&auml;ngt. (2) <strong>Alert auf FAILED-Sagas mit unfertiger Kompensation</strong> &mdash; das ist der Operator-Eingriff. (3) <strong>Idempotente Kompensations-Endpoints</strong>, damit ein manueller Retry gefahrlos m&ouml;glich ist.<br><strong>Spicy:</strong> Retry weglassen ist erlaubt &mdash; aber dann muss das <em>Monitoring der Retry sein</em>. Was du nicht im Code hast, musst du im Dashboard haben. Was du in keinem von beiden hast, hast du nicht.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Eventing statt sync?</h3>
<p>Booking macht synchrones <code>DELETE</code>. W&auml;re ein <span class="hl"><code>CancelBooking</code>-Event</span> nicht nat&uuml;rlicher?</p>
<code>Story 6 wartet</code>
<aside class="notes"><strong>Meine Antwort:</strong> Doch &mdash; genau das ist der Sprung von <strong>Orchestration &uuml;ber sync HTTP</strong> (Story 5) zu <strong>Event-getriebener Saga</strong> (Story 6, Choreography). Eventing verschiebt Retry-Logik in den Broker, Hotel-Ausfall wird zur Queue-Pufferung, Booking-Crash mitten drin kostet keine Recovery. <em>Aber:</em> Booking ist <em>nicht</em> fertig nach &bdquo;Event raus&ldquo;. Der Kunde will eine Antwort &mdash; Hotel publiziert <code>BookingCancelled</code> zur&uuml;ck, Booking konsumiert das Reply-Event, Saga-Status wird nachgef&uuml;hrt, Timeout-Erkennung f&uuml;r ausbleibende Replies bleibt n&ouml;tig.<br><strong>Spicy:</strong> Eventing ist die robustere Architektur &mdash; aber sie <em>verschiebt</em> die Komplexit&auml;t, sie eliminiert sie nicht. Die Verantwortung f&uuml;r die Ausf&uuml;hrung wandert zu Hotel; die Verantwortung f&uuml;r den Gesamtstatus gegen&uuml;ber dem Kunden bleibt bei Booking.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> 204 statt 404</h3>
<p><code>DELETE /bookings/{id}</code> liefert <span class="hl">204</span> auch bei unbekannter ID. Verschleiert das echte Fehler?</p>
<code>Idempotenz &gt; Ehrlichkeit</code>
<aside class="notes"><strong>Meine Antwort:</strong> Nein &mdash; und der Grund kommt direkt aus der Saga-Mechanik. (1) <strong>Idempotenz ist ein Feature</strong>: Saga-Retries d&uuml;rfen mehrfach denselben Storno absetzen &mdash; bei <code>404</code> w&uuml;rde die Retry-Logik f&auml;lschlich als &bdquo;Fehler&ldquo; werten, obwohl der Effekt l&auml;ngst eingetreten ist. (2) <strong>Ohne State kann der Service ehrlich gar nicht unterscheiden</strong> zwischen &bdquo;nie existiert&ldquo;, &bdquo;bereits storniert&ldquo; und &bdquo;gerade erst gebucht und schon wieder weg&ldquo;. (3) <strong>Saga-Praxis</strong>: Kompensations-Endpoints sind &bdquo;at-least-once safe&ldquo; &mdash; 2xx auf jeden plausiblen Aufruf.<br><strong>Spicy:</strong> Idempotenz schl&auml;gt Ehrlichkeit. Ein Kompensations-Endpoint, der &bdquo;korrekte&ldquo; 4xx liefert, zwingt jeden Aufrufer dazu, die 4xx wieder als &bdquo;eigentlich ok&ldquo; zu interpretieren &mdash; das ist die schlechtere Stelle f&uuml;r die Sonderlogik.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Wie merken?</h3>
<p>Saga h&auml;ngt zwischen Forward und Kompensation. Wer <span class="hl">schl&auml;gt Alarm</span>?</p>
<code>Status-Counter</code>
<aside class="notes"><strong>Meine Antwort:</strong> Eine Logzeile ist die Untergrenze. Eine Saga lebt von <em>Beobachtbarkeit ihres Zustands</em>: <strong>Status-Verteilung</strong> als Counter (Alarm bei <code>COMPENSATING &gt; 0</code> f&uuml;r &gt; N Min), <strong>Saga-Dauer</strong> als Histogramm (P99 &gt; erwarteter Wert), <strong>Compensation-Erfolgsrate</strong> (&lt; 99 % &rarr; Backend krank), <strong>Retry-Counter pro Schritt</strong>, <strong>DLQ-Tiefe</strong> bei Eventing, <strong>distributed Tracing</strong> mit Saga-ID in jedem Span.<br><strong>Spicy:</strong> &bdquo;Wir loggen das&ldquo; ist die Antwort von Teams, die noch keine h&auml;ngende Saga in Produktion gesehen haben. Eine h&auml;ngende Saga ist nicht laut &mdash; sie ist <em>still</em>. Sie schreibt keine Fehlermeldung mehr, weil der Step, der sie geschrieben h&auml;tte, nicht mehr l&auml;uft. Das einzige Signal: ein Counter, der zu lange auf einem Wert steht.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> Wer ist Owner?</h3>
<p>Hotel / Flight / Car wissen <span class="hl">nichts</span> von der Saga. Richtig so?</p>
<code>Wissen zentral</code>
<aside class="notes"><strong>Meine Antwort:</strong> Bei Orchestration-Saga genau richtig &mdash; und das ist ein Kernargument f&uuml;r Orchestration. Booking als Orchestrator wei&szlig; Reihenfolge, Fortschritt, was kompensiert werden muss, aktuellen Status. Hotel / Flight / Car bleiben dumm und einfach: sie kennen ihre lokalen Transaktionen, sonst nichts. Vorteile: Backends k&ouml;nnen in <em>anderen</em> Sagen mitspielen ohne Code-&Auml;nderung; ein Bug in der Saga-Logik steckt an <strong>einer</strong> Stelle, nicht in dreien; ein neues Backend (z.&nbsp;B. &bdquo;Mietboot&ldquo;) kommt ohne Anpassung der Bestehenden hinzu.<br><strong>Spicy:</strong> Orchestration konzentriert das Wissen. Choreography (Story 6) verteilt es. Beides ist legitim, aber <em>Wissen verteilen ohne Plan</em> f&uuml;hrt zu &bdquo;verteiltem Monolith&ldquo; &mdash; der schlimmsten beider Welten.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">7</span> Saga ohne DB?</h3>
<p>Wir halten den Status <span class="hl">in-memory</span>. In Produktion akzeptabel?</p>
<code>Crash-Recovery</code>
<aside class="notes"><strong>Meine Antwort:</strong> Wichtige Klarstellung: <strong>Persistenz ist kein Saga-spezifisches Thema</strong>. Jeder mehrstufige Prozess, der im RAM l&auml;uft, ist beim Crash weg. Saga macht das Problem nur <em>sichtbarer</em>, weil die Schritte externe Seiteneffekte (gebuchte Fl&uuml;ge, gestartete Zahlungen) hinterlassen. Spektrum f&uuml;r durable State: Relationale DB, Document Store, Event Log, Embedded SQLite, KV-Store (Consul KV ist im Stack), <strong>Workflow Engine</strong> (Temporal / Camunda &mdash; dreht das Modell um: Engine l&ouml;st Recovery f&uuml;r dich).<br><strong>Spicy:</strong> In Produktion ist die wichtigere Frage selten &bdquo;brauche ich eine DB?&ldquo;, sondern &bdquo;<em>schreibe ich die Orchestrator-Mechanik selbst, oder nutze ich Temporal?</em>&ldquo; &mdash; letzteres untersch&auml;tzen Teams regelm&auml;&szlig;ig und schreiben dann monatelang das, was Temporal seit Jahren in Produktion l&ouml;st.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story5.md</code>.
</aside>
