# Workshop-Fragen: Design-Session REST vs. RESTful (Story 2)

Provokante Fragen für den Recap der Design-Session. Ziel: die sechs
Prinzipien nicht als Stilregeln abhaken, sondern zeigen, wo sie
Betriebsprobleme verhindern, und wo Abweichen legitim ist.

---

## Frage 1 · Der Client schickt den Storno zweimal. Was passiert?

**Frage:** Instabiles Netz, die erste Antwort geht im Timeout verloren,
die App wiederholt die Anfrage. Was passiert bei `DELETE /bookings/4711`,
was bei `POST /bookings/4711/cancellation`?

**Antwort:** Beim `DELETE` nichts Schlimmes. Der erste Aufruf storniert,
der zweite findet eine bereits stornierte Buchung und antwortet `204`
(oder `404`, je nach Modell). Der Zustand ist derselbe, die Methode ist
idempotent, der Retry ist sicher.

Beim `POST` legt der zweite Aufruf eine zweite Stornierung an, wenn der
Server nicht aufpasst. Zwei Wege:

1. **Fachliche Prüfung:** Die Buchung ist schon storniert, der Server
   antwortet `409 Conflict`. Der Client muss den Konflikt verstehen.
2. **Idempotency-Key:** Der Client schickt einen eindeutigen Schlüssel
   im Header (`Idempotency-Key: 8f3a…`). Der Server erkennt die
   Wiederholung und liefert die gespeicherte erste Antwort noch einmal.
   Stripe und viele Zahlungs-APIs arbeiten so.

**Spicy Take-away:** Idempotenz ist kein akademisches Merkmal, sondern
die Voraussetzung dafür, dass Retry überhaupt erlaubt ist. Wer POST
ohne Idempotency-Key wiederholt, bucht doppelt. Wer deshalb gar nicht
wiederholt, verliert Buchungen im Timeout. Beides ist schlechter als
fünf Zeilen Server-Code.

---

## Frage 2 · Warum ist `200 OK` mit Fehler im Body ein Betriebsproblem und nicht nur ein Stilbruch?

**Frage:** Der Client liest den Body sowieso. Wen stört es, wenn dort
`{ "status": "error" }` steht und der Status-Code trotzdem 200 ist?

**Antwort:** Alle, die den Body nicht lesen. Und das ist fast die
gesamte Infrastruktur:

| Komponente | Liest nur den Code | Folge bei 200 mit Fehler-Body |
|---|---|---|
| Load Balancer, Health-Check | ja | Kaputter Service bleibt im Pool |
| Monitoring, Alerting | ja | Fehlerrate zeigt 0 %, niemand wird geweckt |
| Retry-Logik im Client | ja | Kein Retry, obwohl er nötig wäre |
| Circuit Breaker (Story 4) | ja | Zählt Erfolge, öffnet nie |
| HTTP-Cache, CDN | ja | Fehlerantwort wird gecacht |

Richtig ist der passende 4xx- oder 5xx-Code plus Details im Body:
`404` für "Buchung existiert nicht", `409` für "schon storniert",
`422` für "Hotel lehnt die Änderung ab", `503` für "Backend nicht
erreichbar".

**Spicy Take-away:** Status-Codes sind die einzige Sprache, die alle
Beteiligten zwischen Client und Service verstehen. Wer sie umgeht,
schaltet die gesamte Betriebsintelligenz ab und wundert sich später
über grüne Dashboards vor kaputten Services.

---

## Frage 3 · `POST /flights/search` mit Filter im Body. Ist das RPC oder darf man das?

**Frage:** Eigentlich müsste eine Suche `GET /flights?from=BRE&to=LIS&…`
sein. Aber der Filter hat dreißig Felder und passt nicht in die URL.
Ist ein POST hier ein Regelbruch?

**Antwort:** Es ist eine Grauzone, und beide Varianten sind vertretbar,
wenn man die Konsequenzen kennt:

- **GET mit Query-Parametern:** cachebar, verlinkbar, idempotent per
  Definition. Grenze: URL-Länge (praktisch 2 bis 8 KB je nach Proxy)
  und Lesbarkeit.
- **POST auf eine Such-Ressource:** Der Body kann beliebig komplex sein.
  Dafür ist der Aufruf nicht cachebar und formal nicht idempotent, auch
  wenn er faktisch keine Seiteneffekte hat. Elasticsearch und viele
  Such-APIs machen es so und dokumentieren es.

Sauberer Mittelweg, wenn die Suche selbst ein Fachobjekt ist:
`POST /searches` legt eine Suche an (`201`, Location-Header), `GET
/searches/{id}` liefert die Ergebnisse. Das ist RESTful und löst das
Längenproblem, kostet aber einen zweiten Roundtrip.

**Spicy Take-away:** Nicht jede Abweichung ist ein Fehler. Ein Fehler
ist eine Abweichung, die niemand begründen und niemand dokumentiert hat.
Die Referenz-Implementierung bricht die Regeln selbst, bewusst und
klein gehalten: `POST /admin/bulkhead-reset` ist ein Knopf fürs
Dashboard, kein Fachobjekt. Genau so eine Begründung sollte jedes Team
für seine eigenen Ausnahmen liefern können.

---

## Frage 4 · Die Saga in Story 6 ruft `DELETE /bookings/{id}` an jedem Backend auf. Warum muss genau dieser Aufruf idempotent sein?

**Frage:** Die Kompensation läuft nur, wenn etwas schiefgegangen ist.
Dann läuft sie halt einmal. Wozu die Idempotenz-Forderung?

**Antwort:** Weil die Kompensation genau in dem Moment läuft, in dem das
System schon wackelt. Der Storno-Aufruf ans Hotel kann selbst im Timeout
landen. Dann weiß der Orchestrator nicht, ob storniert wurde, und muss
wiederholen. Ist `DELETE` nicht idempotent (zweiter Aufruf liefert
`500`, weil "Buchung nicht gefunden"), bleibt die Saga im Zustand
`COMPENSATING` hängen, oder sie meldet einen Fehler, obwohl fachlich
alles in Ordnung ist.

Die Regel aus der Design-Session ("DELETE zweimal ist derselbe Zustand")
ist also keine Schönheit, sondern die Voraussetzung dafür, dass
Kompensation überhaupt zuverlässig zu Ende kommt. Die Backends der
Referenz antworten deshalb auf jedes `DELETE /bookings/{id}` mit `204`,
auch beim zweiten Mal (sie speichern nichts und sind damit trivial
idempotent).

**Spicy Take-away:** Wer in Story 2 am Flipchart über Idempotenz
diskutiert hat, versteht in Story 6 sofort, warum "Kompensation muss
letztlich gelingen" nur mit idempotenten Endpoints funktioniert.
Design-Entscheidungen an der API sind Betriebsentscheidungen mit
Verzögerung.
