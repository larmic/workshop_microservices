# Story 4: Wenn der Flug ausfällt

**Thema:** Circuit Breaker, als Rollenspiel mit Schalter und Strichliste statt Code
**Zeitrahmen:** ca. 60 Minuten (20 Minuten Spiel, 10 Minuten Demo mit der Referenz, Recap)

## Kontext

Der FlightService ist zeitweise nicht erreichbar (Netzwerkprobleme, Überlast, Deployment). Ohne Schutz wartet der Booking-Service bei jeder Anfrage drei Sekunden auf den Timeout, bindet dabei Threads und Verbindungen und reißt seine eigenen Aufrufer mit. Ein Circuit Breaker zählt die Fehler, schaltet nach fünf ab, antwortet sofort mit einem Fallback und probiert nach einer Pause vorsichtig wieder zu.

Der Breaker ist Workshop-Didaktik in achtzig Zeilen, in Produktion nimmt man die Library. Deshalb tippen wir ihn nicht ab, sondern spielen ihn.

## User Story

Als **Kundin**
möchte ich **Hotel und Mietwagen auch dann sehen und vergleichen können, wenn die Flugsuche ausfällt**,
damit **ich nicht vor einer hängenden Seite sitze**.

---

## Teil 1: Das Rollenspiel (ca. 20 Minuten)

Ihr seid die Suchanfragen. Booking hat den Schalter.

### Material

- Eine Strichliste und ein Schild mit drei Seiten: CLOSED, OPEN, HALF_OPEN
- Karten „Hotel ✓, Car ✓, flights: []" für die Fallback-Antworten
- Eine Uhr mit Sekundenanzeige

### Rollen

1. **Booking** ist eine Person. Sie hat den Schalter: eine Strichliste und ein Schild mit CLOSED, OPEN, HALF_OPEN.
2. **Flight, Hotel, Car** sind drei Leute. Hotel und Car antworten sofort. Flight ist abgestürzt und sagt gar nichts. Bei kleiner Gruppe reicht eine Person für Hotel und Car.
3. **Alle anderen** sind Suchanfragen auf `GET /booking/offers` und stellen sich bei Booking an.

### Ablauf

1. **Runde 1, ohne Schalter:** Booking fragt Flight, Flight schweigt, Booking zählt laut bis drei: Timeout. Booking sagt „Fehler", die Anfrage stellt sich hinten wieder an. Hotel und Car werden gar nicht erst gefragt, so macht es die Referenz aus Story 3.
2. **Runde 2, mit Schalter:** Jeder Timeout ein Strich auf Bookings Liste. Beim fünften Strich dreht Booking sein Schild auf OPEN.
3. Ab jetzt fragt Booking Flight gar nicht mehr, nur Hotel und Car, und gibt die Karte aus: Hotel ✓, Car ✓, flights: []. Die Anfrage setzt sich.
4. Alle 30 Sekunden dreht Booking auf HALF_OPEN und fragt Flight genau einmal (der Trainer ist die Uhr). Schweigt Flight: zurück auf OPEN. Beim dritten Mal antwortet Flight wieder: CLOSED, volle Angebote.
5. **Runde 3, freies Spiel:** Flight entscheidet selbst, wann es antwortet und wann nicht, auch mal nur jede zweite Anfrage. Booking hält sich nur an seine Regeln: zählen, nach fünf Fehlern OPEN, alle 30 Sekunden ein Probe-Call.

### Was passieren sollte

- Runde 1: Die Schlange bleibt gleich lang, jede Anfrage kostet drei Sekunden, niemand setzt sich.
- Runde 2: Nach dem fünften Strich leert sich die Schlange sichtbar. Hotel und Car arbeiten weiter, der Schalter gilt pro Backend.
- Runde 3: Der Schalter pendelt. Ein wackliges Flight kommt nie dauerhaft auf CLOSED, weil der Zähler nach jedem Erfolg bei null beginnt und nach fünf Fehlern wieder öffnet. Das ist der Moment für die Recap-Fragen 1 und 3: Was zählt als Fehler, und warum nur ein Probe-Call.

---

## Teil 2: Die Demo mit der Referenz (ca. 10 Minuten)

Im Dashboard unter Story 4: Flight auf „Fehler" stellen, fünfmal `GET /booking/offers` auslösen. Der Flight-Breaker geht auf OPEN, die Antwort enthält sofort `flights: []` und den Header `X-Circuit-Open: flight`. Flight zurück auf „Normal", nach 30 Sekunden schließt der Breaker über HALF_OPEN. Drehbuch im Detail weiter unten.

---

## Und die Buchung?

Lesen und Schreiben sind zwei Fälle. Für die Angebote ist der Fallback sinnvoll, bei der Buchung ist eine Teilbuchung eine Fachentscheidung, die der Breaker nicht trifft. Nach dem Recap stellen wir die Frage offen: Flight-Breaker ist OPEN, die Kundin klickt auf Buchen. Hotel und Mietwagen trotzdem buchen? Mögliche Antworten: sofort und ehrlich Nein, asynchron nachbuchen, Kompensation per Saga (Story 6), alternativer Anbieter.

---

## Zum Nachbauen zu Hause

Wer den Breaker selbst implementieren will, findet hier die ursprüngliche Coding-Aufgabe. Sprache und Framework sind frei.

### Akzeptanzkriterien

- [ ] Je ein Circuit Breaker ist um jeden Backend-Aufruf (Flight, Hotel, Car) implementiert, jeweils mit eigenem State
- [ ] Nach 5 aufeinanderfolgenden Fehlern öffnet sich der Circuit
- [ ] Bei offenem Circuit wird ein Fallback ausgeführt (z.B. `flights: []`), und der Header `X-Circuit-Open` nennt dem Client das betroffene Backend
- [ ] Der Circuit schließt sich nach 30 Sekunden wieder (Half-Open State)
- [ ] Der aktuelle Circuit-Status ist über einen Endpoint abfragbar
- [ ] Timeouts sind konfiguriert (max. 3 Sekunden Wartezeit)
- [ ] Metriken über Circuit-Zustandswechsel werden geloggt

### Technische Hinweise

- **Empfohlene Libraries:**
  - Resilience4j (empfohlen, framework-agnostisch)
  - Spring Cloud Circuit Breaker
  - MicroProfile Fault Tolerance (Quarkus)
- **Circuit Breaker Zustände:**
  - CLOSED: Normale Funktion, Fehler werden gezählt
  - OPEN: Calls werden sofort mit Fallback beantwortet
  - HALF_OPEN: Testweise einzelne Calls durchlassen
- **Konfigurationsbeispiel (Resilience4j):**
  ```yaml
  resilience4j:
    circuitbreaker:
      instances:
        flightService:
          failureRateThreshold: 50
          waitDurationInOpenState: 30000
          slidingWindowSize: 10
  ```

### Bonus (optional)

- Visualisiere die Circuit-States in einem Dashboard.

## Reference-Implementierung

Unter `services/booking/story4/` liegt eine vollständige Go-Reference, die alle Akzeptanzkriterien samt Bonus abdeckt:

- **Selbstgebauter Circuit Breaker** (keine Library) — `services/booking/story4/circuitbreaker/circuitbreaker.go`. Drei Zustände, lazy Übergang von OPEN nach HALF_OPEN beim nächsten Call, atomarer Probe-Slot gegen Probe-Storms in HALF_OPEN.
- **Drei separate CBs** für Flight, Hotel und Car. Jeder Backend-Aufruf läuft durch seinen eigenen Breaker; bei Fallback wird der Teilbereich leer zurückgegeben (`flights: []` bzw. `flight: null`).
- **Detailliertes Logging**: jede Zustandsänderung, jede SHORT-CIRCUIT-Entscheidung mit Restzeit, jeder Fehlerzähler-Inkrement und jeder Probe-Versuch wird mit Begründung im Log ausgegeben — so ist im Workshop nachvollziehbar, *warum* der Breaker gerade so handelt.
- **Dashboard-Visualisierung** (`services/dashboard/`): pro CB eine Karte mit Badge (CLOSED/OPEN/HALF_OPEN), Countdown bis HALF_OPEN, Failure-Counter, Anzahl Fallbacks. Aktualisierung im 1-Sekunden-Takt.
- **Chaos-Steuerung**: Backend-Services (Flight, Hotel, Car) bieten `/admin/chaos` mit den Modi `normal`, `slow`, `fail`. Das Dashboard erlaubt das Schalten pro Service oder pro Replica, sodass HALF_OPEN-Recovery live demonstriert werden kann.

### Workshop-Drehbuch

1. Stack starten: `docker compose -f services/docker-compose.yml -f services/docker-compose.infra.yml -f services/docker-compose.reference.yml up -d`
2. Dashboard öffnen (`http://localhost/dashboard`). Alle drei CBs zeigen CLOSED.
3. Flight im Dashboard auf "Fehler" stellen. Fünfmal `curl localhost:8088/booking/offers` schicken — der Flight-CB wechselt nach OPEN, Antwort enthält `flights: []` und Header `X-Circuit-Open: flight`.
4. Logs des `booking-ref-story4` Containers ansehen (`docker compose logs -f booking-ref-story4`) — jede Entscheidung des Breakers ist sichtbar.
5. Flight zurück auf "Normal". Nach 30 Sekunden geht der CB beim nächsten Call automatisch in HALF_OPEN, der Probe-Call schließt den Breaker wieder.
