## Quiz 3/3

<p class="subtitle">RESTful oder nicht?</p>

<div class="box">

### <code>POST /booking/bookings</code>

<pre><code>&rarr; 200 OK
{ "status": "error", "message": "hotel not available" }</code></pre>

</div>

<div class="box fragment">

<strong>Nicht RESTful.</strong> Methode und Pfad stimmen, aber der Status-Code l&uuml;gt: HTTP sagt Erfolg, der Body sagt Fehler.

Richtig: <code>409 Conflict</code> oder <code>422 Unprocessable Content</code>, Details im Body. Status-Codes sind Infrastruktur, keine Kosmetik.

</div>

Note:
- Erwartung: die meisten sagen &bdquo;ja, sieht sauber aus&ldquo;. Wer &bdquo;nein&ldquo; sagt, hat die zweite Zeile gelesen.
- Wem tut das weh? Allen, die den Body nicht lesen: Load Balancer und Health-Checks (Service bleibt im Pool), Monitoring (Fehlerrate 0 %), Retry-Logik (kein Retry), Caches (Fehler wird gecacht) und der Circuit Breaker aus Story 4, der Erfolge z&auml;hlt und nie &ouml;ffnet.
- Welcher Code? 404 Buchung existiert nicht, 409 Konflikt (schon storniert, Doppelbuchung), 422 fachlich abgelehnt (Hotel nicht verf&uuml;gbar), 503 Backend nicht erreichbar. Diskutieren, ob 409 oder 422 hier besser passt, beides ist vertretbar.
- Das ist die beste Br&uuml;cke in Kapitel 2: Resilience-Patterns funktionieren nur, wenn die Fehler sichtbar sind.
- Reserve, falls Zeit bleibt: <code>POST /flights/search</code> (Grauzone), <code>PUT</code> mit Teil-Objekt, <code>POST /admin/bulkhead-reset</code> aus unserer eigenen Referenz (bewusst RPC-artig, ein Aufrufer, unter <code>/admin/</code>). Alle in <code>docs/instructions/rest-vs-restful.md</code>, Abschnitt 7.
- Quelle f&uuml;r Nachleser: Martin Fowler, &bdquo;Richardson Maturity Model&ldquo;.
