<div class="page">

<p class="kicker">REST vs. RESTful</p>

<div class="page-head">

## Sechs Regeln

<span class="hint hint-late">Hypermedia lassen wir bewusst weg.</span>
</div>

<p class="subtitle">&hellip; die Betriebsprobleme verhindern.</p>

<div class="page-body rules-body">

<div class="quiz-code wide"><span class="ex ex-m">DELETE</span> <span class="ex ex-p">/customers/7/bookings/4711</span> &nbsp;&nbsp;<span class="ex ex-a">&rarr; 204 No Content</span></div>

<div class="rules">
<div class="rules-group fragment" data-fragment-index="1" data-notes="Frage in die Runde: Wir loggen jeden GET-Aufruf und schreiben ihn in eine Datenbank. Ist das schon ein Verstoß gegen „GET ohne Seiteneffekte“?
Antwort: Nein. Die Regel meint die Ressource, nicht den Server. GET darf Logs, Zähler und Caches füllen, solange der Client danach denselben Zustand sieht. Verboten ist, was die Ressource ändert oder fachlich etwas auslöst.">
<h4>Die Methode sagt, was passiert</h4>
<div class="rule">
<h5><span class="num">1</span> Tr&auml;gt die Semantik</h5>
<p>GET liest, POST erzeugt, PUT ersetzt, PATCH &auml;ndert, DELETE entfernt.</p>
<code>POST /users/1/delete &rarr; DELETE /users/1</code>
</div>
<div class="rule">
<h5><span class="num">2</span> Idempotenz</h5>
<p>GET, PUT und DELETE beliebig oft wiederholbar. POST nicht.</p>
<code>DELETE zweimal = derselbe Zustand</code>
</div>
<div class="rule">
<h5><span class="num">3</span> GET ohne Seiteneffekte</h5>
<p>Prefetcher, Crawler und Caches f&uuml;hren GET ungefragt aus.</p>
<code>GET /bookings/1/cancel &rarr; niemals</code>
</div>
</div>
<div class="rules-group fragment" data-fragment-index="2" data-notes="- <strong>Ressourcen statt Verben:</strong> Die Aktion steckt schon in der Methode. Wer Verben in den Pfad schreibt, baut RPC &uuml;ber HTTP.
- <strong>Pfad vs. Query:</strong> Der Pfad sagt, wem etwas geh&ouml;rt (Kunde 7), die Query, welche davon ich will (best&auml;tigte).
- <strong>Warum das z&auml;hlt:</strong> Einheitliche Pfade machen Routing, Logs und Berechtigungen pro Ressource einfach.">
<h4>Der Pfad sagt, woran</h4>
<div class="rule">
<h5><span class="num">4</span> Ressourcen statt Verben</h5>
<p>Pfade benennen Dinge, keine Aktionen.</p>
<code>GET /getUser?id=1 &rarr; GET /users/1</code>
</div>
<div class="rule">
<h5><span class="num">5</span> Sub-Ressourcen und Filter</h5>
<p>Zugeh&ouml;rigkeit in den Pfad, Auswahl in die Query.</p>
<code>/customers/7/bookings?status=confirmed</code>
</div>
</div>
<div class="rules-group fragment" data-fragment-index="3" data-notes="Nur der Status-Code zählt, zum Beispiel für:
- Load Balancer: nimmt eine Instanz bei 5xx aus der Rotation
- Retry-Logik: wiederholt bei 503
- Circuit Breaker (Story 4): zählt Fehler
- Monitoring: alarmiert auf Fehlerraten
- Caches: speichern auch 404 oder 301, und leider auch den 200 mit Fehler drin
Ein 200 mit Fehler im Body sieht keiner davon.
Überleitung: Jetzt wendet ihr die Regeln selbst an, Storno und Umbuchung am Flipchart.">
<h4>Die Antwort sagt, wie es ausging</h4>
<div class="rule">
<h5><span class="num">6</span> Status-Codes statt Fehler-Body</h5>
<p>Die Infrastruktur liest nur den Code. Ein 200 mit Fehler im Body ist unsichtbar.</p>
<code>404 &middot; 409 &middot; 422 &middot; 503</code>
</div>
</div>
</div>

</div>

</div>

Note:
- <strong>Kernbotschaft:</strong> Sechs Regeln, jede verhindert ein Betriebsproblem. Das ist kein Stilthema.
- <strong>Das Beispiel oben:</strong> Ein Aufruf beantwortet drei Fragen: Die Methode sagt, was passiert, der Pfad, woran, der Status-Code, wie es ausging. Genau so bauen sich die drei Gruppen auf.
- <strong>Hypermedia:</strong> REST nach Fielding verlangt auch Links in den Antworten. Unser Ziel ist Richardson Level 2: Ressourcen plus HTTP-Methoden.
- <strong>Mechanik:</strong> Gruppe f&uuml;r Gruppe aufdecken, insgesamt 10 bis 12 Minuten.
