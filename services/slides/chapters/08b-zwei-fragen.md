<div class="page">

<p class="kicker">Was wir in zwei Tagen beantworten</p>

## Zwei Fragen

<div class="page-body">

<div class="numlist">
<div class="numlist-row">
<div class="numlist-num">01</div>
<div>
<h3>Was m&uuml;ssen wir tun, damit es funktioniert?</h3>
<p>Health und Config, Discovery, Resilience, Konsistenz, Sichtbarkeit. Acht Stories lang.</p>
</div>
</div>
<div class="numlist-row">
<div class="numlist-num">02</div>
<div>
<h3>Was <span class="hl">kostet</span> uns das?</h3>
<p>Jedes Pattern bringt neue Fehlerf&auml;lle mit. Die sammeln wir nach jeder Story im Recap.</p>
</div>
</div>
</div>

<div class="foot">
<div class="callout">Die zweite Frage ist die unbequeme.</div>
</div>

</div>

</div>

Note:
- Was das Netzwerk uns kostet, fassen wir hier nicht im Voraus auf, die Stories machen es konkret:
  - Partielle Ausf&auml;lle &rarr; Circuit Breaker (Story 4), Bulkhead (Story 5)
  - Keine verteilten Transaktionen &rarr; Saga (Story 6/7)
  - Schnittstellen, auf die sich Infrastruktur verlassen kann &rarr; REST vs. RESTful (Story 2, Design-Session ohne Code)
  - Auffindbarkeit &rarr; Service Discovery (Story 3)
  - Nachvollziehbarkeit &rarr; Distributed Tracing (Story 8)
  - Konsistenz &rarr; Eventual Consistency, CQRS (Diskussion)
- Was on top kommt (Diskussion, kein Hands-on): Service-Schnitt (DDD, Event Storming), Auth &uuml;ber Service-Grenzen (OAuth/SAML), synchron vs. asynchron, Versionierung, Team-Schnitte (Conway's Law), Polyglot Persistence, Observability als eigene Disziplin.
- Take-away: Microservices sind ein Werkzeug, kein Ziel. Am Ende der zwei Tage habt ihr eine ehrliche Aufwandseinsch&auml;tzung und damit eine bessere Entscheidungsgrundlage f&uuml;r euer eigenes Projekt.
- &Uuml;berleitung: Bevor wir Resilienz angehen, brauchen wir das Fundament, einen Service, der &uuml;berhaupt l&auml;uft und sich beobachten l&auml;sst. Story 1.
