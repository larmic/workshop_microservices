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

## Teil 1: Das Rollenspiel „Ein Booking, ein Stapel" (ca. 20 Minuten)

Ihr seid die Suchanfragen. Flight ist krank, und was Booking tut, steht auf der Karte.

### Material

- Ein Kartenstapel 00–27 (A6 hoch, Karton 250 g, 7 Bogen). Oben steht „Booking tut", unten „Deine Antwort" mit Status und Dauer, Body und Headern. Graue Karten sind STOPP- bzw. ENDE-Karten
- Die Zustandskarten CLOSED, OPEN und HALF_OPEN, eine davon steht im Kartenaufsteller vor Booking, zur Schlange gedreht
- Die Fehlerkarten 1 bis 5, die 5 ist rot („5 → OPEN!"). Zustand und Fehlerkarten sind das Dashboard: Das sieht nur der Betrieb, nicht die Kundschaft
- Ein Ablagetisch, auf dem die gespielten Karten der Reihe nach landen
- Das Trainer-Skript (ein Bogen, beidseitig): Ansage, die fünf STOPPs, „Wenn es hakt", Vereinfachungen und auf Seite 2 die Kartentabelle 00–27

Gesamtliste: [`docs/orga/materialliste.md`](../orga/materialliste.md), Druckvorlagen und Trainer-Skript unter <https://larmic.github.io/workshop_microservices/material/>.

Spielbar ab vier Personen, derselbe Stapel für jede Gruppengröße.

### Rollen

1. **Booking** ist eine Person. Sie hat den Stapel, die Zustandskarten und die Fehlerkarten und arbeitet die Karten streng der Reihe nach ab: oben lesen, genau das tun, Karte dem nächsten Gast geben.
2. **Alle anderen** sind Suchanfragen auf `GET /booking/offers`. Flight, Hotel und Car spielt niemand, ihr Verhalten steht auf den Karten.
3. **Der Trainer** spielt keine Rolle. Er moderiert nur an den grauen Karten, Booking gibt sie ihm.

### So läuft es

Wer eine Karte bekommt, liest „Deine Antwort" laut vor, legt die Karte der Reihe nach auf den Ablagetisch und stellt sich **sofort hinten wieder an**. Die Karten 01–03 spielen ohne Schalter wie Story 3: Booking zählt laut bis fünf und sagt „500". Ab Karte 04 hat Booking einen Schalter und zählt nur bis drei: Jeder Fehler legt die nächste Fehlerkarte hin, bei der 5 dreht Booking die Zustandskarte auf OPEN und antwortet ab dann sofort mit „Circuit open!". An den STOPPs vergehen die 30 Sekunden im Gespräch, danach kommt HALF_OPEN und genau ein Gast ist die Probe. Ein Erfolg legt alle Fehlerkarten zurück. Zum Schluss haben sich fünf Gäste beim Datum vertippt. Der Ablauf Satz für Satz steht im Trainer-Skript.

### Wo wir reden

| STOPP | Karte | Worüber wir reden | Recap |
|---|---|---|---|
| STOPP 1 | 04 | Was habt ihr bekommen, wie lange? Hotel und Car waren gesund, was hattet ihr davon? Story 3 bricht beim ersten Fehler mit 500 ab. Ansage: Booking bekommt einen Schalter und 3 s Timeout | – |
| STOPP 2 | 13 | Am Ablagetisch: Ab welcher Nummer gab es `flights: []`, was ist ab 10 anders? Steht die Zahl der Fehler in eurer Antwort? Fallback ab dem ersten Fehler, OPEN macht ihn nur schnell und setzt `X-Circuit-Open` (09 hat Flight noch gefragt). Hotel und Car haben eigene Schalter | 02 |
| STOPP 3 | 17 | Wer war bei Flight? Was haben 15 und 16 bekommen, wie schnell? Was, wenn alle Wartenden zu Flight gegangen wären? Genau ein Probe-Slot, der Rest sofort Fallback | 03 |
| STOPP 4 | 20 | Vorhersage, nicht auflösen: Flight ist gesund, die nächsten haben sich vertippt, Datum 31.02. Gibt das Fehlerkarten? | – |
| ENDE | 27 | War Flight bei 21 bis 25 kaputt? Wer wurde bei 26 bestraft? Unser Code zählt alles ab 400. Dann: Was zeigt die Zustandskarte einer zweiten Booking-Instanz, hätte sie bei 05 bis 09 auch fünfmal 3 s gewartet, was weiß Booking nach einem Neustart? | 01, 04 |

### Bewusste Vereinfachungen

Offen ansagen:

- Bis fünf zählen steht für die 5 s Timeout in Story 3, bis drei für die 3 s in Story 4.
- Die 30 Sekunden vergehen im Gespräch. Im Code wechselt der Zustand beim nächsten Call oder beim nächsten Blick ins Dashboard.
- Die Referenz zählt auch die gescheiterte Probe (Karte 14) mit, das Dashboard zeigt dann 6/5. Im Spiel gibt es keine sechste Fehlerkarte.
- Das Datum 31.02. ist ein Gedankenspiel: Die Referenz reicht keine Suchparameter an Flight weiter.
- Wieder anstellen ist ein Client-Retry, den die Referenz nicht hat.
- Fehler in Folge gegen Fehlerquote im Fenster (Resilience4j) ist kein Spielblock, nur Gesprächsstoff: Ein wackliges Flight, das jedes zweite Mal antwortet, öffnet unseren Schalter nie, trotzdem wartet jede zweite Suche 3 s.

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
