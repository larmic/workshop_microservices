## Story 3

<p class="subtitle">Wenn der Flug ausf&auml;llt <span class="time-badge">&asymp; 60 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

Der Flight-Service ist zeitweise nicht erreichbar (Netzwerkprobleme, &Uuml;berlast, Deployment). Der Booking-Service soll in diesem Fall nicht ebenfalls ausfallen, sondern <em>graceful degradieren</em> und dem Nutzer eine sinnvolle Alternative bieten.

#### User Story

Als <em>Kunde</em> m&ouml;chte ich <em>auch bei Ausfall des Flugbuchungssystems eine Teilbuchung (Hotel + Mietwagen) durchf&uuml;hren k&ouml;nnen</em>, damit <em>meine Reiseplanung nicht komplett blockiert wird</em>.

</div>

</div>
<div>

<div class="story-card">

#### Akzeptanzkriterien

- Je ein Circuit Breaker um <em>jeden</em> Backend-Aufruf (Flight, Hotel, Car); nach <code>5</code> aufeinanderfolgenden Fehlern &ouml;ffnet er
- Bei offenem Circuit l&auml;uft ein Fallback (z.&nbsp;B. <code>flights: []</code> mit Header <code>X-Circuit-Open</code>)
- Nach <code>30&nbsp;s</code> wechselt der Breaker in <code>HALF_OPEN</code> und l&auml;sst einen Probe-Call durch
- Aufruf-Timeout konfiguriert (max. <code>3&nbsp;s</code>)
- Aktueller Circuit-Status &uuml;ber einen Admin-Endpoint abfragbar
- Zustands&auml;nderungen werden geloggt

</div>

</div>
</div>

Note:
- Hook: &bdquo;Story 2 hat uns geholfen, Services zu <em>finden</em>. Heute kl&auml;ren wir, was passiert, wenn wir einen gefunden haben &mdash; und er antwortet nicht.&ldquo;
- Wiedererkennung: dieselbe Karte (Kontext / User Story / Akzeptanzkriterien) im Dashboard unter Story 3 &rarr; &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter <code>services/booking/story3/</code> (Go, selbstgebauter CB &ndash; ca. 80 Zeilen).
- Drei separate CBs (Flight / Hotel / Car) statt einem globalen &mdash; isoliert Ausf&auml;lle pro Backend. Granularit&auml;t kommt im Recap (Frage 7) zur&uuml;ck.
- Demo-Tipp: Dashboard &ouml;ffnet die Chaos-Schalter pro Service / pro Replica. Flight auf &bdquo;Fehler&ldquo; &rarr; nach 5 Calls geht der CB OPEN &rarr; Antwort enth&auml;lt sofort <code>flights: []</code> und Header <code>X-Circuit-Open: flight</code>. Zur&uuml;ck auf &bdquo;Normal&ldquo; &rarr; nach 30 s schliesst der Breaker &uuml;ber HALF_OPEN.
- Resilience-Libraries: Resilience4j (Java), Polly (.NET), Spring Cloud Circuit Breaker, MicroProfile Fault Tolerance, gobreaker. Selber bauen ist Workshop-Didaktik &mdash; in Produktion <em>nimmt</em> man die Library.
- Time-Box 60 min inkl. Demo. Dashboard <code>http://localhost</code> zeigt den Story-3-Modus inkl. Spickzettel mit Pseudo-Code.
- Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-03-circuit-breaker.md</code>.
