# Story 4: Wenn der Flug ausfällt

**Thema:** Circuit Breaker, als Rollenspiel mit Kartenstapel statt Code
**Zeitrahmen:** ca. 60 Minuten (20 Minuten Spiel, 10 Minuten Demo mit der Referenz, Recap)

## Kontext

Der FlightService ist zeitweise nicht erreichbar (Netzwerkprobleme, Überlast, Deployment). Ohne Schutz wartet der Booking-Service bei jeder Anfrage drei Sekunden auf den Timeout, bindet dabei Threads und Verbindungen und reißt seine eigenen Aufrufer mit. Ein Circuit Breaker zählt die Fehler, schaltet nach fünf Fehlern in Folge ab, antwortet sofort mit einem Fallback und probiert nach einer Pause vorsichtig wieder zu.

Der Breaker ist Workshop-Didaktik in achtzig Zeilen, in Produktion nimmt man die Library. Deshalb tippen wir ihn nicht ab, sondern spielen ihn.

## User Story

Als **Kundin**
möchte ich **Hotel und Mietwagen auch dann sehen und vergleichen können, wenn die Flugsuche ausfällt**,
damit **ich nicht vor einer hängenden Seite sitze**.

---

## Teil 1: Das Rollenspiel „Ein Booking, ein Board" (ca. 20 Minuten)

Ihr seid die Requests. Flight ist krank, und was Booking tut, steht auf der Karte.

### Material

- Ein Kartenstapel 00–27 (A7 hoch, Karton 250 g, 8 pro A4-Bogen). Oben sieht Booking Flight, Hotel und Car mit grünem Haken, rotem Kreuz oder grau für „nicht fragen". Unten steht „Request · deine Antwort": Wartezeit, „Antwort HTTP …", Body und Header
- Die Farbe einer Karte sagt, wohin sie danach geht: dunkel auf die Ablage, rot aufs Board, gelb aufs Probe-Feld, grau zum Trainer
- Das Circuit-Breaker-Board (A3 quer oder 2 × A4): fünf Plätze „Fehler in Folge", Platz 5 heißt OPEN, dazu das Probe-Feld und das Zustandsfeld. Das Board ist das Dashboard: Das sieht nur der Betrieb, nicht die Kundschaft
- Die Zustandskarten CLOSED, OPEN und HALF_OPEN (A7 hoch) für das Zustandsfeld und eine Legendenkarte
- Ein Ablagetisch, auf dem die Karten der Nummer nach landen
- Das Trainer-Skript (ein Bogen, beidseitig): Ansage, Aufbau, die fünf STOPPs, „Wenn es hakt", Vereinfachungen und auf Seite 2 die Kartentabelle 00–27

Gesamtliste: [`docs/orga/materialliste.md`](../orga/materialliste.md), Druckvorlagen und Trainer-Skript unter <https://larmic.github.io/workshop_microservices/material/>.

Spielbar ab vier Personen, derselbe Stapel für jede Gruppengröße.

### Rollen

1. **Booking** ist eine Person. Sie hat den Stapel und gibt die Karten streng der Reihe nach aus, die nächste erst, wenn die vorige liegt. Nur bei der gelben Probe geht es sofort weiter. Booking legt nichts ab, schaut aufs Board, ruft „OPEN!" und tauscht die Zustandskarte.
2. **Alle anderen** sind Requests auf `GET /booking/offers`. Flight, Hotel und Car spielt niemand, ihr Verhalten steht auf den Karten.
3. **Der Trainer** spielt keine Rolle. Er moderiert nur an den grauen Karten, Booking gibt sie ihm.

### So läuft es

Wer eine Karte bekommt, wartet die Zeit darauf ab, liest die Antwort laut vor, legt die Karte dorthin, wo das Band unten sagt, und stellt sich **sofort hinten wieder an**.

- **01–03** spielen ohne Circuit Breaker wie Story 3: 5 Sekunden warten, HTTP 500, keine Hotels, keine Autos. Das Board liegt noch nicht auf dem Tisch.
- **STOPP 1 (04)** ist das Update: Booking bekommt einen Circuit Breaker und 3 s Timeout, legt das Board hin und CLOSED ins Zustandsfeld.
- **05–09:** Booking fragt Flight, Flight schweigt. Der Request wartet 3 s und bekommt den Fallback. Die Karten sind rot und kommen auf Platz 1 bis 5. Liegt Platz 5, ruft Booking „OPEN!".
- **10–12:** OPEN, Flight wird nicht gefragt, die Antwort kommt sofort mit `X-Circuit-Open`.
- **STOPP 2 und 3:** Die 30 Sekunden vergehen im Gespräch, danach HALF_OPEN. Die gelben Karten 14 und 18 sind die Probe und liegen auf dem Probe-Feld, solange sie läuft. 14 scheitert (wieder OPEN, die 30 s laufen von vorn), 18 gelingt: Booking räumt das Board ab, CLOSED.
- **21–25:** Flight ist gesund, aber das Datum 31.02. bringt einen 400. Die Karten kommen trotzdem aufs Board, bei 25 wieder „OPEN!".

Der Ablauf Satz für Satz steht im Trainer-Skript.

### Wo wir reden

| STOPP | Karte | Worüber wir reden | Recap |
|---|---|---|---|
| STOPP 1 | 04 | Was habt ihr bekommen, wie lange? Hotel und Car waren gesund, was hattet ihr davon? Story 3 bricht beim ersten Fehler mit 500 ab. Ansage: Booking bekommt einen Circuit Breaker und 3 s Timeout | – |
| STOPP 2 | 13 | Am Board: Welche Karten liegen hier, ab welcher Nummer gab es `flights: []`, was ist ab 10 anders? Was davon war der Breaker, was der Fallback? Fallback ab dem ersten Fehler, auch bei CLOSED. OPEN fragt Flight nicht und setzt `X-Circuit-Open` (09 hat Flight noch gefragt). Hotel und Car haben eigene Breaker | 02 |
| STOPP 3 | 17 | Wer lag auf dem Probe-Feld? Was haben 15 und 16 bekommen, wie schnell? Was, wenn alle zu Flight gegangen wären? Genau ein Probe-Slot, der Rest sofort Fallback | 03 |
| STOPP 4 | 20 | Vorhersage, nicht auflösen: Flight ist gesund, die nächsten haben sich vertippt, Datum 31.02. Kommen die aufs Board? | – |
| ENDE | 27 | War Flight bei 21 bis 25 kaputt? Wer wurde bei 26 bestraft? Unser Code zählt alles ab 400. Dann: Eine zweite Booking-Instanz hat ein eigenes Board, was liegt darauf? Was weiß Booking nach einem Neustart? | 01, 04 |

### Bewusste Vereinfachungen

Offen ansagen:

- Der Request wartet selbst: 5 s Timeout in Story 3, 3 s in Story 4. Booking gibt die nächste Karte erst aus, wenn die vorige liegt. Im Code laufen Requests parallel.
- Die 30 Sekunden vergehen im Gespräch. Im Code wechselt der Zustand beim nächsten Call oder beim nächsten Blick ins Dashboard.
- Die Referenz zählt auch die gescheiterte Probe (Karte 14) mit, das Dashboard zeigt dann 6/5. Im Spiel geht 14 vom Probe-Feld zur Ablage, das Board ist schon voll.
- Das Datum 31.02. ist ein Gedankenspiel: Die Referenz reicht keine Suchparameter an Flight weiter.
- Wieder anstellen ist ein Client-Retry, den die Referenz nicht hat.
- Fehler in Folge gegen Fehlerquote im Fenster (Resilience4j, [Frage 8](../questions/story4.md)) ist kein Spielblock, nur Gesprächsstoff: Ein wackliges Flight, das jedes zweite Mal antwortet, öffnet unseren Breaker nie, trotzdem wartet jede zweite Suche 3 s.

---

## Teil 2: Die Demo mit der Referenz (ca. 10 Minuten)

Im Dashboard unter Story 4: Flight auf „Hängt" stellen (antwortet erst nach 10 s, also nach dem 3-s-Timeout von Booking), dann `GET /booking/offers` auslösen. Die ersten fünf Calls warten je 3 Sekunden, genau wie die Karten 05–09. Schneller geht es mit „Fehler" (sofort ein 500, ohne Warten). Der Client bekommt immer 200 mit `flights: []` und `X-Fallback: flight`. Der fünfte Call öffnet den Breaker, ab dem sechsten steht zusätzlich `X-Circuit-Open: flight` in der Antwort. Flight zurück auf „Normal": Nach 30 Sekunden zeigt das Panel HALF_OPEN schon ohne Call, die erfolgreiche Probe schließt den Breaker. Drehbuch im Detail weiter unten.

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
  Achtung: Resilience4j zählt hier eine Fehlerquote im Fenster, unsere Referenz zählt Fehler in Folge. Ein wackliges Flight, das jedes zweite Mal antwortet, würde Resilience4j öffnen, unsere Referenz nicht.

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
