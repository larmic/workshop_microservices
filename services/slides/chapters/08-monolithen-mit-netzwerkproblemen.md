## Microservices

<p class="subtitle">&hellip; sind auch nur Monolithen mit <span class="hl">Netzwerkproblemen</span></p>

<div class="box">

### Zwei Fragen, die wir im Workshop beantworten

- Was m&uuml;ssen wir f&uuml;r eine <span class="hl">funktionierende Architektur</span> tun?
- Was <span class="hl">kostet</span> uns das?

</div>

Note:
- These provokant in den Raum stellen: stimmt das? Wo trifft sie zu, wo nicht?
- Diskussions-Anker: Welche Probleme h&auml;tte man im Monolithen auch &mdash; und welche entstehen erst durch das Netzwerk?
- Was das Netzwerk uns kostet, fassen wir hier nicht im Voraus auf &mdash; die Stories machen es konkret:
  - Partielle Ausf&auml;lle &rarr; Circuit Breaker (Story 3), Bulkhead (Story 4)
  - Keine verteilten Transaktionen &rarr; Saga (Story 5/6)
  - Auffindbarkeit &rarr; Service Discovery (Story 2)
  - Nachvollziehbarkeit &rarr; Distributed Tracing (Story 7)
  - Konsistenz &rarr; Eventual Consistency, CQRS (Diskussion)
- Was on top kommt (Diskussion, kein Hands-on): Service-Schnitt (DDD, Event Storming), Auth &uuml;ber Service-Grenzen (OAuth/SAML), synchron vs. asynchron, Versionierung, Team-Schnitte (Conway's Law), Polyglot Persistence, Observability als eigene Disziplin.
- Take-away: Microservices sind ein Werkzeug, kein Ziel. Am Ende der zwei Tage habt ihr eine ehrliche Aufwandseinsch&auml;tzung &mdash; und damit eine bessere Entscheidungsgrundlage f&uuml;r euer eigenes Projekt.
- &Uuml;berleitung: Bevor wir Resilienz angehen, brauchen wir das Fundament &mdash; einen Service, der &uuml;berhaupt l&auml;uft und sich beobachten l&auml;sst. Story 1.
