# Story 2: Design-Session: REST vs. RESTful

> Unser Booking-Service liefert jetzt Angebote. Als Nächstes sollen Kundinnen und Kunden buchen können, und wo gebucht wird, wird auch storniert und umgebucht. Bevor jemand losprogrammiert, hätten wir gern eine API, die man nicht in drei Monaten schon wieder umbauen muss.
>
> Wir haben im Haus genug Schnittstellen, die sich "REST" nennen und dann `GET /getUser?id=1423` machen. Beim letzten Mal hat ein Prefetcher im Browser über einen solchen GET-Link Buchungen storniert. Das soll uns nicht noch einmal passieren. Erst denken, dann tippen: Einigt euch am Flipchart darauf, wie Storno und Umbuchung als Ressourcen aussehen, und begründet, warum.

**Zeitrahmen:** ca. 55 Minuten (Theorie-Input, Teamarbeit am Flipchart, Vorstellung, Quiz)

## Thema

- **REST vs. RESTful** → HTTP benutzen ist nicht dasselbe wie HTTP richtig benutzen
- **Ressourcen statt Verben** → Substantive im Pfad, die Semantik trägt die HTTP-Methode
- **Idempotenz und Status-Codes** → Retries werden sicher, Infrastruktur kann Fehler sehen

---

## Ziel

Eine am Flipchart entworfene API für Storno und Umbuchung, die:
- Ressourcen und Sub-Ressourcen statt Aktionen im Pfad verwendet
- die HTTP-Methoden so einsetzt, dass Idempotenz und Seiteneffekte klar sind
- Fehler über Status-Codes meldet, nicht über einen Fehler-Body mit `200 OK`
- jede Abweichung von den Regeln bewusst begründet

**Diese Story hat keine Referenz-Implementierung, keinen Port und kein Docker-Image.** Gearbeitet wird am Flipchart in Teams zu drei bis vier Personen. Die Regeln aus dieser Session wendet ihr direkt in [Story 3](story-03-service-discovery.md) an, wenn `POST /booking/bookings` entsteht.

---

## Aufgaben

### 1. Szenario verstehen

Zwei Anfragen aus dem Kundenservice:

- Eine Kundin möchte ihre komplette Reise (Flug, Hotel, Mietwagen) stornieren. Der Storno-Grund soll nachvollziehbar bleiben.
- Ein Kunde möchte nur den Flug seiner bestehenden Buchung umbuchen. Hotel und Mietwagen bleiben.

Randbedingungen: Der Client (Web und App) hat instabile Netze und wiederholt Anfragen bei Timeout. Der Kundenservice möchte später nachsehen können, wer wann was storniert hat.

### 2. API am Flipchart entwerfen

Ergebnis ist eine Endpoint-Tabelle mit einer Zeile pro Endpoint:

| Methode | Pfad | Request (Kern) | Antwort und Status-Code | Idempotent? | Begründung |
|---|---|---|---|---|---|
| … | … | … | … | ja / nein | … |

Leitfragen für die Teamarbeit:

- Ist Storno ein `DELETE` auf die Buchung, eine neue Ressource (`…/cancellation`) oder eine Statusänderung per `PATCH`? Was passiert mit dem Storno-Grund?
- Ist Umbuchung ein `PUT` auf die Buchung, ein `PATCH`, eine Sub-Ressource (`…/flight`) oder Storno plus Neubuchung?
- Was passiert, wenn der Client dieselbe Anfrage zweimal schickt, weil die erste Antwort im Timeout verloren ging?
- Welcher Status-Code kommt zurück, wenn die Buchung nicht existiert, schon storniert ist oder das Hotel die Änderung ablehnt?
- Gibt es einen Endpoint, bei dem ihr die Regeln bewusst brecht? Warum?

### 3. Vorstellen und vergleichen

- Jedes Team stellt seine Tabelle in zwei Minuten vor
- Der Trainer sammelt die Varianten nebeneinander und moderiert die Trade-offs
- Anschließend Quiz im Plenum: drei Endpoints, Handzeichen, "RESTful oder nicht?"

### 4. Brücke zu Story 3

- Nehmt eure Regeln mit: In Story 3 baut ihr `POST /booking/bookings` selbst
- In Story 6 (Saga) kommt der Storno als Kompensation zurück. Dann zahlt sich die Diskussion ein zweites Mal aus
