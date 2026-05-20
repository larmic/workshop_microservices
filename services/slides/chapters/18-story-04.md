## Story 4

<p class="subtitle">Isolation ist St&auml;rke <span class="time-badge">&asymp; 60 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

Wenn ein Backend-Service extrem langsam antwortet, k&ouml;nnen alle Threads / Connections des Booking-Service in diesem einen Service stecken bleiben. Aufrufe an die anderen Backends sind dann ebenfalls betroffen &mdash; obwohl sie gesund sind. Das Bulkhead-Pattern isoliert die Ressourcen pro Downstream.

#### User Story

Als <em>Betriebsteam</em> m&ouml;chte ich, <em>dass Probleme mit einem Backend-Service nicht die Aufrufe an andere Backend-Services beeintr&auml;chtigen</em>, damit <em>ein langsamer oder fehlerhafter Service nicht das gesamte System blockiert</em>.

</div>

</div>
<div>

<div class="story-card">

#### Akzeptanzkriterien

- Eigener Pool pro Backend (Flight, Hotel, Car); Gr&ouml;&szlig;e konfigurierbar
- Bei vollem Pool: <strong>sofortige</strong> Ablehnung (kein Queueing) mit <code>503</code>
- Langsamer Hotel-Service blockiert <em>nicht</em> die Aufrufe an Flight oder Car
- Aufruf-Timeout konfiguriert (max. <code>3&nbsp;s</code>)
- Metriken &uuml;ber Pool-Auslastung verf&uuml;gbar (in-flight, calls, rejected)
- Funktioniert <em>zusammen</em> mit dem Circuit Breaker aus Story 3

</div>

</div>
</div>

Note:
- Hook: &bdquo;Story 3 hat uns gegen <em>kaputte</em> Backends geh&auml;rtet. Heute geht's um die nervigere Variante: das Backend antwortet, nur eben sehr, sehr langsam. Kein Fehler &mdash; und trotzdem rei&szlig;t es alles mit.&ldquo;
- Wiedererkennung: dieselbe Karte (Kontext / User Story / Akzeptanzkriterien) im Dashboard unter Story 4 &rarr; &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story4/</code> (Go, Semaphore via <code>chan struct{}</code>).
- Drei separate Bulkheads (Flight / Hotel / Car) &mdash; gleiche Granularit&auml;t wie beim CB aus Story 3. Im Code stehen Bulkhead und CB als Decorator hintereinander: <code>bh.Execute &rarr; cb.Execute &rarr; httpCall</code>.
- Demo-Drehbuch: Dashboard &rarr; Hotel auf &bdquo;Langsam (2&nbsp;s)&ldquo; stellen, dann <code>POST /admin/burst</code> dr&uuml;cken &mdash; 20 parallele Requests. Hotel-Bulkhead zeigt rejected &asymp;&nbsp;10, Flight und Car laufen praktisch ungest&ouml;rt. <em>Recap-Hook</em>: bei genauer Beobachtung sind Flight-Rejects sauber, Hotel-Rejects unsauber &mdash; warum?
- Time-Box 60 min inkl. Demo. Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-04-bulkhead.md</code>.
