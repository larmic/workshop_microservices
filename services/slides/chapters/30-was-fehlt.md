<div class="page">

<p class="kicker">Zwei Tage reichen nicht f&uuml;r alles</p>

## Was wir nicht angeschaut haben

<p class="subtitle">Jedes Thema verdient einen eigenen Tag.</p>

<div class="page-body">

<div class="cards cards-3 violet compact">
<div class="card">
<h3>Service-Schnitt</h3>
<p>Wo h&ouml;rt ein Service auf? DDD, Event Storming, Bounded Contexts.</p>
</div>
<div class="card">
<h3>Authn / Authz</h3>
<p>OAuth, OIDC, Identit&auml;t &uuml;ber Service-Grenzen weiterreichen.</p>
</div>
<div class="card">
<h3>Versionierung</h3>
<p>API-Vertr&auml;ge, Breaking Changes, zwei Versionen parallel betreiben.</p>
</div>
<div class="card">
<h3>Synchron vs. Asynchron</h3>
<p>Wann Request/Response, wann Event? Zeitliche Kopplung, Queues, Backpressure. Story 7 hat es nur gestreift.</p>
</div>
<div class="card">
<h3>Daten &uuml;ber Service-Grenzen</h3>
<p>Eine Datenbank pro Service, CQRS, Event Sourcing.</p>
</div>
<div class="card">
<h3>Observability als Disziplin</h3>
<p>RED-Metriken, SLOs, Alerts. Tracing war erst der Anfang.</p>
</div>
</div>

</div>

</div>

Note:
- Das ist der Eisberg vom ersten Vormittag (Folie &bdquo;Zwei Fragen&ldquo;), hier die gr&ouml;&szlig;ten Brocken unter Wasser.
- <strong>Synchron vs. Asynchron:</strong> Timeout, Breaker und Bulkhead sind die Rechnung f&uuml;r synchrone Aufrufe. Wer asynchron entkoppelt, zahlt mit Eventual Consistency und einem Broker, siehe Story 7.
- Jeder Punkt verdient seinen eigenen Tag. Wir hatten zwei. Br&uuml;cke zu Folgeterminen und weiterf&uuml;hrenden Workshops.
- Service-Schnitt kam im Feedback als Wunsch nach Aufgaben ohne Coding. Kandidat f&uuml;r eine 30-Minuten-&Uuml;bung am ersten Tag beim n&auml;chsten Durchlauf.
- Auch nicht als eigenes Thema: Timeouts und Retry mit Backoff. Timeout kam in Story 4 vor, Retry als Forward Recovery in Story 7.
- Repo-Link und Kontakt f&uuml;r Follow-up: <code>github.com/larmic/workshop_microservices</code>, Stories, Code, Fragen als GitHub Issue.
