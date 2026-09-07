# REST vs. RESTful · Workshop-Notizen

> Trainer-Notizen für Story 2 (Design-Session, ca. 55 Minuten, kein Code). Aufbau wie bei den anderen Trainer-Hinweisen: pro Abschnitt ein kleiner Hinweis, was an der Tafel oder auf der Folie passieren sollte.

## 1. Worum geht es?

**Kernbefund aus dem Kickoff (2026-05-30):** Fast jeder Service "spricht REST". Die Prinzipien dahinter sind aber oft unbekannt oder werden falsch angewendet. Typische Funde in echten Codebasen: `GET /getUser?id=1423`, `POST /users/1423/delete`, Storno per GET-Link, Fehler als `200 OK` mit `{ "error": … }` im Body. Das ist nicht nur Stil: Prefetcher, Crawler, Retries, Load Balancer und Monitoring verlassen sich auf die HTTP-Semantik. Wer sie bricht, bekommt gelöschte Daten, Doppelbuchungen und grüne Dashboards vor kaputten Services.

**Ziel der Einheit:** Die Teilnehmenden entwerfen einmal selbst eine API, bevor sie in Story 3 den ersten schreibenden Endpoint bauen. Der Lösungsraum ist bewusst offen, der Trainer moderiert Trade-offs statt eine Musterlösung zu verteilen.

> 🎯 *Einstieg:* Ein oder zwei der Kickoff-Beispiele an die Tafel schreiben und fragen: "Was ist daran falsch, und wem tut das weh?" Noch nicht auflösen, das kommt im Quiz.

---

## 2. REST vs. RESTful: sechs Prinzipien

**REST** ist ein Architekturstil (Fielding, 2000). **RESTful** heißt im Alltag: die HTTP-Semantik wird so genutzt, wie sie gemeint ist. Die meisten "REST-APIs" sind RPC über HTTP mit JSON. Für den Workshop reichen sechs Regeln:

| # | Prinzip | Gegenbeispiel | Besser |
|---|---|---|---|
| 1 | **Ressourcen statt Verben.** Pfade benennen Dinge, keine Aktionen | `GET /getUser?id=1423` | `GET /users/1423` |
| 2 | **Die HTTP-Methode trägt die Semantik.** GET liest, POST erzeugt, PUT ersetzt, PATCH ändert teilweise, DELETE entfernt | `POST /users/1423/delete` | `DELETE /users/1423` |
| 3 | **Idempotenz.** GET, PUT, DELETE dürfen beliebig oft wiederholt werden, POST nicht. Das macht Retries sicher | `POST /bookings/1423/cancel` zweimal geschickt: zweiter Storno schlägt fehl oder storniert doppelt | `DELETE /bookings/1423` zweimal: 204 dann 404 (oder 204), Zustand identisch |
| 4 | **GET hat keine Seiteneffekte.** Browser-Prefetch, Crawler und Caches dürfen GET jederzeit ausführen | `GET /bookings/1423/cancel` | `DELETE /bookings/1423` oder `POST /bookings/1423/cancellation` |
| 5 | **Status-Codes statt Fehler-Body.** Die Infrastruktur liest nur den Code | `200 OK` mit `{ "status": "error" }` | `404 Not Found`, `409 Conflict`, `422 Unprocessable Content` mit Fehlerdetails im Body |
| 6 | **Sub-Ressourcen und Filter statt Sonder-Endpoints.** Zugehörigkeit im Pfad, Auswahl in Query-Parametern | `GET /getBookingsForCustomer?id=7` | `GET /customers/7/bookings?status=confirmed` |

Nur als Ausblick, nicht im Workshop vertiefen: HATEOAS (Links in Antworten) und das Richardson Maturity Model (Level 0 bis 3, Martin Fowler). Für die Teilnehmenden ist Level 2 (Ressourcen plus HTTP-Verben) das realistische Ziel.

> 🎯 *Folie:* Die sechs Prinzipien erscheinen als Karten (`10b-rest-vs-restful.md`), die Gegenüberstellung als Tabelle (`10c-rest-vs-restful-beispiele.md`). Zeit: 10 bis 12 Minuten, nicht länger, die Teilnehmenden sollen selbst denken.

---

## 3. Ablauf der Design-Session

| Phase | Zeit | Inhalt |
|---|---|---|
| Theorie-Input | 10 bis 12 Min | Sechs Prinzipien anhand der Kickoff-Beispiele (Abschnitt 2) |
| Aufgabe stellen | 3 Min | Szenario (Abschnitt 4), Teams zu 3 bis 4 Personen, ein Flipchart pro Team |
| Teamarbeit | 20 Min | Endpoint-Tabelle: Methode, Pfad, Request-Kern, Antwort und Status-Code, idempotent ja/nein, Begründung |
| Vorstellung und Vergleich | 8 bis 10 Min | Je Team 2 Minuten, Trainer sammelt die Varianten nebeneinander (Abschnitt 5) |
| Quiz "RESTful oder nicht?" | 8 bis 10 Min | Drei Endpoints, Handzeichen, Auflösung (Abschnitt 7) |
| Brücke | 2 Min | "In Story 3 baut ihr `POST /booking/bookings`. Nehmt die Regeln mit." |

**Material:** Ein Flipchart oder eine Whiteboard-Fläche pro Team, dicke Stifte, die Endpoint-Tabelle als Vorlage (Spaltenköpfe vorzeichnen spart fünf Minuten). Remote: ein Miro- oder Excalidraw-Board pro Team.

**Teamgröße:** Drei bis vier Personen. Bei zwei Personen fehlt die Reibung, bei fünf redet einer nicht.

---

## 4. Szenario und Leitfragen

Zwei Anfragen aus dem Kundenservice:

- **Storno:** Eine Kundin möchte ihre komplette Reise (Flug, Hotel, Mietwagen) stornieren. Der Storno-Grund soll nachvollziehbar bleiben.
- **Umbuchung:** Ein Kunde möchte nur den Flug seiner bestehenden Buchung umbuchen. Hotel und Mietwagen bleiben.

Randbedingungen: Der Client hat instabile Netze und wiederholt Anfragen bei Timeout. Der Kundenservice will später nachsehen, wer wann was storniert hat.

Leitfragen (stehen auch in der Story):

1. Ist Storno ein `DELETE` auf die Buchung, eine neue Ressource oder eine Statusänderung? Wo bleibt der Grund?
2. Ist Umbuchung `PUT`, `PATCH`, Sub-Ressource oder Storno plus Neubuchung?
3. Was passiert, wenn der Client dieselbe Anfrage zweimal schickt?
4. Welcher Status-Code bei "existiert nicht", "schon storniert", "Hotel lehnt ab"?
5. Gibt es einen Endpoint, bei dem ihr die Regeln bewusst brecht? Warum?

---

## 5. Lösungsraum: typische Varianten und Trade-offs

### Storno

| Variante | Vorteil | Nachteil |
|---|---|---|
| `DELETE /bookings/{id}` | Klar, idempotent, jeder versteht es | Kein Platz für den Grund (DELETE-Body ist unüblich), Buchung "verschwindet", Historie schwierig |
| `POST /bookings/{id}/cancellation` | Storno ist eine eigene Ressource mit Grund, Zeitpunkt, Gebühr; `GET` darauf liefert die Historie; Buchung bleibt als Datensatz | POST ist nicht idempotent: zweiter Aufruf braucht `409 Conflict` oder eine Idempotency-Key-Prüfung |
| `PATCH /bookings/{id}` mit `{ "status": "cancelled", "reason": … }` | Ein Endpoint für alle Statusänderungen | Status-Maschine versteckt sich im Body, Validierung wird komplex, PATCH-Semantik (Merge vs. JSON Patch) oft unklar |

Moderationshinweis: Alle drei sind vertretbar. Die zweite ist meist die beste Antwort auf "Grund nachvollziehbar", und sie taucht in Story 6 als Kompensation wieder auf.

### Umbuchung

| Variante | Vorteil | Nachteil |
|---|---|---|
| `PUT /bookings/{id}` mit vollständigem Objekt | Idempotent, einfach | Client muss alles kennen und mitschicken; weggelassene Felder werden gelöscht |
| `PATCH /bookings/{id}` mit `{ "flightId": … }` | Nur die Änderung wird geschickt | Teilausführung: was, wenn der neue Flug nicht verfügbar ist? Idempotenz hängt von der Implementierung ab |
| `PUT /bookings/{id}/flight` mit `{ "flightId": … }` | Sub-Ressource macht klar, was sich ändert, idempotent | Mehr Endpoints; was ist mit Umbuchung von Hotel und Flug zugleich? |
| Storno plus Neubuchung | Keine neue API nötig | Zwei Vorgänge, Preis- und Verfügbarkeitsrisiko dazwischen, Kundin sieht kurz "keine Buchung" |

Moderationshinweis: Hier gibt es keine eindeutig richtige Antwort. Wichtig ist, dass die Teams die Teilausführung und die Idempotenz benennen. Das ist die Vorlage für die Saga in Story 6.

---

## 6. Häufige Fehler und wie man sie am Flipchart anspricht

- **Verben im Pfad** (`/cancelBooking`, `/doRebook`): Fragen, welche Ressource dahintersteht. Meist fällt das Substantiv sofort.
- **Alles POST:** Fragen, was bei einem Retry passiert. Wenn die Antwort "dann bucht er doppelt" ist, ist der Punkt gemacht.
- **Nur 200 und 500:** Fragen, wie der Client "Buchung existiert nicht" von "Hotel lehnt ab" unterscheidet. Und wie das Monitoring einen fachlichen Fehler von einem technischen trennt.
- **IDs als Query-Parameter** (`?id=`): Fragen, ob die Buchung eine Adresse hat. Ressourcen brauchen eine URL.
- **PUT mit Teil-Objekt:** Fragen, was mit den weggelassenen Feldern passiert. Wer "bleiben so" sagt, meint PATCH.
- **RPC-Endpoints ohne Begründung:** Nicht verbieten, sondern nach dem Grund fragen. Abweichen ist erlaubt, wenn man weiß, warum (siehe Abschnitt 9).

---

## 7. Quiz "RESTful oder nicht?"

Drei Endpoints, feste Reihenfolge, Dramaturgie: eindeutiges Nein, trügerisches Ja, trügerisches Nein. Pro Beispiel: Endpoint zeigen, Handzeichen abfragen, eine Person aus der Minderheit begründen lassen, dann auflösen (Fragment auf der Folie). Slides: `10f-quiz-1.md` bis `10h-quiz-3.md`.

### Quiz 1 · Eindeutig nicht RESTful

```
GET /booking/cancelBooking?id=4711
```

Erwartung: fast alle sagen "nicht RESTful". Auflösung: Verb im Pfad, Identität als Query-Parameter, GET mit Seiteneffekt. Anekdote: 2005 hat der Google Web Accelerator Links auf Webseiten vorgeladen, um Seiten schneller zu machen. Bei der 37signals-Anwendung Backpack waren "Löschen"-Links einfache GET-Links. Der Prefetcher hat Nutzern ihre Daten gelöscht. Seitdem ist "GET verändert nichts" Selbstschutz, keine Stilfrage. Besser: `DELETE /booking/bookings/4711` oder Quiz 2.

### Quiz 2 · Sieht falsch aus, ist RESTful

```
POST /booking/bookings/4711/cancellation
```

Erwartung: viele sagen "nein, da steht eine Aktion drin". Auflösung: "cancellation" ist ein Substantiv, also eine Ressource. Der Client legt eine Stornierung an, bekommt `201 Created` und kann sie später per `GET` nachlesen (Grund, Zeitpunkt, Gebühr). Oft besser als ein nacktes DELETE, weil die Buchung als Historie bleibt. Direkter Bezug zur Flipchart-Aufgabe, dort taucht diese Variante fast immer auf. Anker aus dem Netz: Stripe modelliert Rückerstattungen genauso als eigene Ressource (`POST /v1/refunds`).

### Quiz 3 · Sieht richtig aus, ist nicht RESTful

```
POST /booking/bookings
→ 200 OK
{ "status": "error", "message": "hotel not available" }
```

Erwartung: die meisten sagen "ja, sauber". Auflösung: Methode und Pfad stimmen, der Status-Code lügt. HTTP sagt Erfolg, der Body sagt Fehler. Load Balancer, Monitoring, Retries und der Circuit Breaker aus Story 4 sehen nur die 200 und halten den Service für gesund. Richtig: `409 Conflict` oder `422 Unprocessable Content` mit dem Fehler im Body. Beste Brücke in die Resilience-Themen: Status-Codes sind Infrastruktur, keine Kosmetik.

### Reserve (falls Zeit bleibt oder ein Team genau das gebaut hat)

- `POST /users/1423/delete` (Methode trägt keine Semantik)
- `POST /flights/search` mit komplexem Filter-Body (Grauzone: Suche als Ressource vs. `GET /flights?from=BRE&…` mit URL-Längen-Grenze; beides vertretbar, wenn der Aufrufer weiß, dass die POST-Suche keine Seiteneffekte hat und nicht gecacht wird)
- `PUT /bookings/1423` mit Teil-Objekt (Grauzone: PUT ersetzt ganz, PATCH ändert teilweise; was passiert mit weggelassenen Feldern?)
- `POST /admin/bulkhead-reset` aus der eigenen Referenz-Implementierung (bewusst RPC-artig fürs Dashboard, Beispiel für "Abweichen mit Grund", siehe Abschnitt 9)
- Ein Beispiel direkt von den Flipcharts der Teams

**Quellen für die Folie:** Martin Fowler, "Richardson Maturity Model" (zugänglichster Einstieg). Roy Fielding, "REST APIs must be hypertext-driven" (2008, strengste Lesart). Für das Quiz reicht Fowler.

---

## 8. Diskussionsfragen

Ausführliche Antworten in [`docs/questions/story2.md`](../questions/story2.md).

1. **Idempotenz bei Retries:** Der Client schickt den Storno zweimal. Was passiert bei `DELETE`, was bei `POST …/cancellation`? Wie hilft ein Idempotency-Key?
2. **Circuit Breaker auf POST:** In Story 4 wird `POST /booking/bookings` vom Circuit Breaker geschützt. Warum ist Fail-Fast bei nicht idempotenten Aufrufen heikel, und was hat das mit Status-Codes zu tun?
3. **Kompensation in Story 6:** Die Saga ruft `DELETE /bookings/{id}` an jedem Backend auf. Warum muss genau dieser Aufruf idempotent sein?
4. **Wann RPC legitim ist:** Wo würdet ihr in eurem Projekt bewusst einen Aktions-Endpoint bauen, und wie dokumentiert ihr das?

---

## 9. Bezug zum Workshop-Code: wo die Referenz bewusst abweicht

Die Referenz-Implementierung hält sich bei der Fach-API an die Regeln (`GET /booking/offers`, `POST /booking/bookings`, `GET /booking/bookings/{id}`, `DELETE /bookings/{id}` an den Backends). Sie weicht an drei Stellen bewusst ab:

| Endpoint | Wo | Warum RPC-artig |
|---|---|---|
| `POST /admin/bulkhead-reset` | `services/booking/story5..8/handler/admin.go` | Dashboard-Knopf, der Zähler zurücksetzt. Kein Fachobjekt, kein Client außer dem Dashboard |
| `POST /admin/sagas-reset` | `services/booking/story6..8/handler/admin.go` | Gleiche Begründung: Demo-Zustand leeren |
| `POST /admin/chaos` | `services/shared/chaos/chaos.go` | Steuert das Fehlverhalten der Backends für die Demo |

Das ist "Abweichen mit Grund": interne Admin-Endpoints, ein einziger Aufrufer, kein Retry-Problem, klar unter `/admin/` abgesetzt. Genau diese Begründung sollen die Teams für ihre eigenen Ausnahmen liefern können. Der Rückgriff passt in den Recap der Bulkhead-Story (Story 5), wenn der Reset-Knopf im Dashboard zum ersten Mal gedrückt wird.
