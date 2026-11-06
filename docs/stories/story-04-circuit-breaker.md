# Story 4: Wenn der Flug ausfällt

**Thema:** Circuit Breaker, als Rollenspiel mit Schalter und Strichliste statt Code
**Zeitrahmen:** ca. 60 Minuten (20 Minuten Spiel, 10 Minuten Demo mit der Referenz, Recap)

## Kontext

Der FlightService ist zeitweise nicht erreichbar (Netzwerkprobleme, Überlast, Deployment). Ohne Schutz wartet der Booking-Service bei jeder Anfrage drei Sekunden auf den Timeout, bindet dabei Threads und Verbindungen und reißt seine eigenen Aufrufer mit. Ein Circuit Breaker zählt die Fehler, schaltet nach fünf Fehlern in Folge ab, antwortet sofort mit einem Fallback und probiert nach einer Pause vorsichtig wieder zu.

Der Breaker ist Workshop-Didaktik in achtzig Zeilen, in Produktion nimmt man die Library. Deshalb tippen wir ihn nicht ab, sondern spielen ihn.

## User Story

Als **Kundin**
möchte ich **Hotel und Mietwagen auch dann sehen und vergleichen können, wenn die Flugsuche ausfällt**,
damit **ich nicht vor einer hängenden Seite sitze**.

---

## Teil 1: Das Rollenspiel (ca. 20 Minuten)

Ihr seid die Suchanfragen. Booking hat den Schalter.

### Material

- Eine Strichliste (in Prospekthülle, mit Whiteboard-Marker beschreibbar, damit Booking sie wegwischen kann) und ein Schild mit drei Seiten: CLOSED, OPEN, HALF_OPEN
- Karten „Hotel ✓ · Car ✓ · flights: []" für die Fallback-Antworten
- Blaue Karten „Hotel ✓ · Car ✓ · Flight ✓" für erfolgreiche Antworten (geglückte Probe, Runde 3a und 3b)
- Fünf rote Karten „Datum 31.02." für Runde 3b
- Ein Korb, in den die Gäste ihre Antwortkarte zurücklegen
- Ein PROBE-Token (3D-gedruckt oder ein beliebiger Gegenstand), es gibt nur eins
- Eine Uhr mit Sekundenanzeige (der Trainer ist die Uhr)

Gesamtliste und Druckvorlagen: [`docs/orga/materialliste.md`](../orga/materialliste.md), druckfertig unter <https://larmic.github.io/workshop_microservices/material/>.

Spielbar ab sieben Personen: Booking, Flight, Hotel + Car und mindestens vier Suchanfragen.

### Rollen

1. **Booking** ist eine Person. Sie hat den Schalter: Strichliste, Schild und die Karten. Wartet Booking auf Flight, zählt es laut: eins, zwei, drei, Timeout.
2. **Flight** folgt den geheimen Ansagen des Trainers. **Hotel und Car** antworten sofort, bei kleiner Gruppe übernimmt eine Person beides.
3. **Alle anderen** sind Suchanfragen auf `GET /booking/offers`. Wer eine Antwort hat, legt die Karte in den Korb und stellt sich **sofort hinten wieder an**. So bleibt immer jemand in der Schlange, auch für die Probe, und Booking nimmt sich aus dem Korb Nachschub, die Karten gehen nie aus.

### Ablauf

1. **Runde 1, ohne Schalter (wie Story 3):** Flight schweigt. Nach drei Sekunden sagt Booking „500" und fragt Hotel und Car gar nicht erst.
2. **Runde 2, mit Schalter:** Bei jedem Timeout macht Booking einen Strich, fragt Hotel und Car und gibt die Karte „Hotel ✓ · Car ✓ · flights: []". Nach fünf Strichen **in Folge** ruft Booking „OPEN" und dreht das Schild. Ab dann fragt Booking Flight nicht mehr: Karte sofort, dazu „Circuit open!".
3. **HALF_OPEN:** 30 Sekunden nach „OPEN" ruft der Trainer „HALF_OPEN". Der nächste Gast ist die Probe: Booking gibt ihm das PROBE-Token und fragt Flight. Nur wer das Token hält, darf zu Flight. Wer dahinter steht, bekommt sofort die Karte und „Circuit open!". Bleibt Flight still: Strich, OPEN, die 30 Sekunden laufen neu. Antwortet Flight: blaue Karte, Liste wegwischen, CLOSED, volle Angebote. Das Token geht in beiden Fällen zurück an Booking.
4. **Runde 3a, wackliges Flight:** Flight antwortet jedes zweite Mal. Jede Antwort gibt die blaue Karte und wischt die Liste weg, jeder Timeout gibt Strich und Karte.
5. **Runde 3b, rote Gäste:** Flight ist gesund. Die nächsten fünf Anfragen bekommen eine rote Karte „Datum 31.02." (kaputtes Datum), Flight sagt darauf sofort „400". Ein Gedankenspiel: Die echte Referenz reicht keine Suchparameter an Flight weiter.

Geheimes Briefing für Flight: Runde 1 und 2 schweigen, erste Probe schweigen, zweite antworten. Runde 3a abwechselnd. Runde 3b gesund, auf rote Karten sofort „400".

### Was passieren sollte

- Runde 1: Null Karten. Jede Anfrage kostet drei Sekunden und endet mit einem Fehler.
- Runde 2: Die Karte gibt es schon ab dem ersten Fehler, nur langsam. Der Schalter macht sie schnell, er erfindet den Fallback nicht (Recap 02). Hotel und Car arbeiten weiter, der Schalter gilt pro Backend.
- HALF_OPEN: Nur einer probiert, die anderen warten nicht, sondern bekommen sofort die Karte (Recap 03).
- Runde 3a: Das Schild geht nie auf OPEN, weil ein Erfolg den Zähler auf null setzt. Trotzdem wartet jede zweite Suche drei Sekunden. Brücke zur Fehlerquote im Fenster (Resilience4j, siehe unten).
- Runde 3b: Ja, das gibt Striche. Unser Code zählt alles ab 400 als Fehler und macht wegen Tippfehlern der Kundschaft dicht (Recap 01).
- Transfer zu Recap 04: Neben Booking steht Booking-2 mit eigener Liste. Es hat die roten Gäste nie gesehen und zeigt CLOSED.

Bewusste Vereinfachungen: Booking zählt bis drei wie der Story-4-Timeout (Story 3 hatte 5 Sekunden). Das Neu-Anstellen in Runde 1 ist ein Client-Retry, den die Referenz nicht hat. HALF_OPEN ruft der Trainer nach Uhr aus, im Code passiert der Wechsel erst beim nächsten Call.

**Variante ab 16 Personen:** Ein echtes Booking-2 mit eigener Schlange und eigener Strichliste. Die roten Gäste gehen nur zu Booking-1.

---

## Teil 2: Die Demo mit der Referenz (ca. 10 Minuten)

Im Dashboard unter Story 4: Flight auf „Hängt" stellen (antwortet erst nach 10 s, also nach dem 3-s-Timeout von Booking), dann `GET /booking/offers` auslösen. Die ersten fünf Calls warten je 3 Sekunden, genau wie im Spiel. Schneller geht es mit „Fehler" (sofort ein 500, ohne Warten). Der Client bekommt immer 200 mit `flights: []` und `X-Fallback: flight`. Der fünfte Call öffnet den Breaker, ab dem sechsten steht zusätzlich `X-Circuit-Open: flight` in der Antwort. Flight zurück auf „Normal": Nach 30 Sekunden geht der Breaker beim nächsten Call auf HALF_OPEN, die erfolgreiche Probe schließt ihn. Drehbuch im Detail weiter unten.

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
- [ ] Nach 30 Sekunden geht der Circuit auf HALF_OPEN, ein erfolgreicher Probe-Call schließt ihn wieder
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
  Achtung: Resilience4j zählt hier eine Fehlerquote im Fenster, unsere Referenz zählt Fehler in Folge. Das wacklige Flight aus Runde 3a würde Resilience4j öffnen, unsere Referenz nicht.

### Bonus (optional)

- Visualisiere die Circuit-States in einem Dashboard.

## Reference-Implementierung

Unter `services/booking/story4/` liegt eine vollständige Go-Reference, die alle Akzeptanzkriterien samt Bonus abdeckt:

- **Selbstgebauter Circuit Breaker** (keine Library) — `services/booking/story4/circuitbreaker/circuitbreaker.go`. Drei Zustände, lazy Übergang von OPEN nach HALF_OPEN beim nächsten Call, atomarer Probe-Slot gegen Probe-Storms in HALF_OPEN.
- **Drei separate CBs** für Flight, Hotel und Car. Jeder Backend-Aufruf läuft durch seinen eigenen Breaker; bei jedem Fehler wird der Teilbereich leer zurückgegeben (`flights: []`). Der POST macht fail fast: Ist ein Breaker OPEN, antwortet ein Pre-Check sofort mit 503, ohne etwas zu buchen.
- **Detailliertes Logging**: jede Zustandsänderung, jede SHORT-CIRCUIT-Entscheidung mit Restzeit, jeder Fehlerzähler-Inkrement und jeder Probe-Versuch wird mit Begründung im Log ausgegeben — so ist im Workshop nachvollziehbar, *warum* der Breaker gerade so handelt.
- **Dashboard-Visualisierung** (`services/dashboard/`): pro CB eine Karte mit Badge (CLOSED/OPEN/HALF_OPEN), Countdown bis HALF_OPEN, Failure-Counter, Anzahl Short-Circuits. Aktualisierung im 1-Sekunden-Takt.
- **Chaos-Steuerung**: Backend-Services (Flight, Hotel, Car) bieten `/admin/chaos` mit den Modi `normal`, `slow` (2 s, unter dem Timeout), `fail` (sofort 500) und `hang` (10 s, über jedem Timeout im Stack). Das Dashboard erlaubt das Schalten pro Service oder pro Replica, sodass HALF_OPEN-Recovery live demonstriert werden kann.

### Workshop-Drehbuch

1. Stack starten: `docker compose -f services/docker-compose.yml -f services/docker-compose.infra.yml -f services/docker-compose.reference.yml up -d`
2. Dashboard öffnen (`http://localhost/dashboard`). Alle drei CBs zeigen CLOSED.
3. Flight im Dashboard auf "Fehler" stellen. Sechsmal `curl -i localhost:8088/booking/offers` schicken — jede Antwort ist 200 mit `flights: []` und `X-Fallback: flight`, der fünfte Call öffnet den Flight-CB, ab dem sechsten kommt `X-Circuit-Open: flight` dazu.
4. Logs des `booking-ref-story4` Containers ansehen (`docker compose logs -f booking-ref-story4`) — jede Entscheidung des Breakers ist sichtbar.
5. Flight zurück auf "Normal". Nach 30 Sekunden geht der CB beim nächsten Call automatisch in HALF_OPEN, der Probe-Call schließt den Breaker wieder.
