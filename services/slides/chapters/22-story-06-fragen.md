<div class="page">

<p class="kicker">Story 6 &middot; Recap</p>

## F&uuml;nf Fragen an euch

<div class="page-body">

<div class="numlist recap compact dense">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<p class="numlist-label">Die Kompensation scheitert</p>
<h3>Flug gebucht, Hotel sagt nein, und der Storno l&auml;uft selbst in <span class="hl">5xx</span>. Was nun?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">must eventually succeed</code>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<p class="numlist-label">Wir kompensieren genau einmal</p>
<h3>Beim Fehlschlag einfach <code>FAILED</code>. Verletzen wir damit das <span class="hl">Pattern</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Monitoring = Retry</code>
</div>
<div class="numlist-row">
<div class="numlist-num">03</div>
<div>
<p class="numlist-label">204 statt 404</p>
<h3><code>DELETE</code> auf eine unbekannte ID liefert <code>204</code>. Verschleiert das <span class="hl">echte Fehler</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Idempotenz &gt; Ehrlichkeit</code>
</div>
<div class="numlist-row">
<div class="numlist-num">04</div>
<div>
<p class="numlist-label">Die Saga h&auml;ngt</p>
<h3>Zwischen Forward und Kompensation bleibt sie stehen. Wer <span class="hl">schl&auml;gt Alarm</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Status-Counter</code>
</div>
<div class="numlist-row">
<div class="numlist-num">05</div>
<div>
<p class="numlist-label">Der Status lebt im RAM</p>
<h3>Booking st&uuml;rzt mitten in <code>COMPENSATING</code> ab. Was ist jetzt mit dem <span class="hl">Flug</span>?</h3>
</div>
<code class="numlist-pill fragment" data-fragment-index="1">Crash-Recovery</code>
</div>
</div>

</div>

</div>

Note:
- Alle f&uuml;nf Fragen stehen sofort da, erst diskutieren lassen. Ein Klick blendet rechts die Merks&auml;tze ein.
- <strong>Kompensation scheitert, meine Antwort:</strong> Saga macht eine starke Annahme: Forward darf scheitern, Kompensation <em>muss letztlich gelingen</em>. Ohne diese Annahme bricht das Konstrukt zusammen. Strategien: Idempotenz (mehrfacher DELETE, gleiches Ergebnis), Retry mit Backoff, persistenter Saga-Log (Crash darf keine offene Kompensation verlieren), Dead-Letter oder Operator-Inbox (Mensch greift ein), Pivot zur fachlichen Alternative (Gutschein statt Storno), Compensation-by-design (Reservierung als Status, nicht als L&ouml;schung). <strong>Spicy:</strong> Eine Saga ohne Plan f&uuml;r gescheiterte Kompensation ist keine Saga, sondern eine optimistische Hoffnung. Die Frage &bdquo;was bei Misserfolg&ldquo; ist das eigentliche Engineering an dem Pattern.
- <strong>Genau einmal, meine Antwort:</strong> Ja, formal verletzen wir das Pattern. Pragmatisch vertretbar, <em>wenn</em> bewusst entschieden. Dann geh&ouml;rt zwingend dazu: Saga-Status persistieren, sonst wei&szlig; niemand, dass eine Saga h&auml;ngt. Alert auf <code>FAILED</code> mit unfertiger Kompensation, das ist der Operator-Eingriff. Idempotente Kompensations-Endpoints, damit ein manueller Retry gefahrlos ist. <strong>Spicy:</strong> Retry weglassen ist erlaubt, aber dann muss das Monitoring der Retry sein. Was du nicht im Code hast, musst du im Dashboard haben. Was du in keinem von beiden hast, hast du nicht.
- <strong>204 statt 404, meine Antwort:</strong> Nein, und der Grund kommt aus der Saga-Mechanik. Saga-Retries d&uuml;rfen denselben Storno mehrfach absetzen; bei <code>404</code> w&uuml;rde der Aufrufer f&auml;lschlich einen Fehler sehen, obwohl der Effekt l&auml;ngst eingetreten ist. Ohne State kann das Backend ohnehin nicht zwischen &bdquo;nie existiert&ldquo; und &bdquo;bereits storniert&ldquo; unterscheiden. Kompensations-Endpoints sind at-least-once safe: 2xx auf jeden plausiblen Aufruf, 4xx nur bei strukturell falschen Anfragen. <strong>Spicy:</strong> Ein Endpoint, der &bdquo;korrekte&ldquo; 4xx liefert, zwingt jeden Aufrufer, die 4xx wieder als &bdquo;eigentlich ok&ldquo; zu interpretieren. Das ist die schlechtere Stelle f&uuml;r die Sonderlogik.
- <strong>Saga h&auml;ngt, meine Antwort:</strong> Eine Logzeile ist die Untergrenze. Eine Saga lebt von der Beobachtbarkeit ihres Zustands: Status-Verteilung als Counter (Alarm bei <code>COMPENSATING &gt; 0</code> l&auml;nger als N Minuten), Saga-Dauer als Histogramm, Compensation-Erfolgsrate, Retry-Counter pro Schritt, bei Eventing die DLQ-Tiefe, Tracing mit Saga-ID in jedem Span. <strong>Spicy:</strong> &bdquo;Wir loggen das&ldquo; ist die Antwort von Teams, die noch keine h&auml;ngende Saga in Produktion gesehen haben. Eine h&auml;ngende Saga ist nicht laut, sie ist still. Das einzige Signal ist ein Counter, der zu lange auf einem Wert steht.
- <strong>Status im RAM, meine Antwort:</strong> Der Flug bleibt gebucht, und niemand wei&szlig; es mehr. Wichtige Klarstellung: Persistenz ist kein Saga-spezifisches Thema, jeder mehrstufige Prozess im RAM ist beim Crash weg. Saga macht es nur sichtbarer, weil die Schritte externe Spuren hinterlassen. Spektrum: relationale DB, Document Store, Event Log, SQLite, Consul KV (ist seit Story 3 im Stack), Workflow Engine. <strong>Spicy:</strong> In Produktion ist die wichtigere Frage selten &bdquo;brauche ich eine DB?&ldquo;, sondern &bdquo;schreibe ich die Orchestrator-Mechanik selbst, oder nehme ich Temporal?&ldquo; Letzteres untersch&auml;tzen Teams regelm&auml;&szlig;ig und bauen monatelang nach, was Temporal seit Jahren in Produktion l&ouml;st.
- <strong>Reserve, Eventing statt sync:</strong> W&auml;re ein <code>CancelBooking</code>-Event nicht nat&uuml;rlicher als der synchrone <code>DELETE</code>? Doch, genau das ist der Sprung zu Story 7. Aber: Booking ist nach &bdquo;Event raus&ldquo; nicht fertig, der Kunde will eine Antwort. Eventing verschiebt die Komplexit&auml;t, es eliminiert sie nicht.
- <strong>Reserve, Owner:</strong> Hotel, Flight, Car wissen nichts von der Saga. Richtig so? Bei Orchestration ja: Die Backends k&ouml;nnen in anderen Sagen mitspielen, ein Bug steckt an einer Stelle, ein neues Backend kommt ohne Anpassung der Bestehenden hinzu. Choreography verteilt das Wissen, Wissen verteilen ohne Plan gibt den verteilten Monolithen.
- Vollst&auml;ndige Antworten und Anekdoten: <code>docs/questions/story6.md</code>.
