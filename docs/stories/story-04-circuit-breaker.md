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

Ihr seid die Requests. Flight ist krank, und was ihr zurückbekommt, steht auf der Karte.

### Material

- Ein Kartenstapel mit 30 Karten: 01–27 und die drei Zustandskarten CLOSED, OPEN, HALF_OPEN (A7 hoch, Karton 200–250 g, 8 pro A4-Bogen, 4 Bögen, 2 Felder frei). Auf jeder Karte stehen nur die große Nummer, die Wartezeit („5 s", „3 s" oder „sofort"), eine Zeile zum Vorlesen und eine kleine Zeile darunter, bei 14, 18 und 21–26 dazu ein Etikett (PROBE, Datum). Unten rechts klein „Story 4", damit sich die Stapel von Story 4 und 5 nicht mischen
- Die Farbe sagt, was passiert ist: weiß mit violettem Rahmen heißt Haufen, rot heißt Flight hat einen Fehler geliefert, gelb ist die Probe, dunkel ist STOPP. Die STOPP-Karten 04, 13, 17, 20 und 27 tragen nur „STOPP n" bzw. „ENDE"
- Das Circuit-Breaker-Board (A3 quer oder 2 × A4): fünf Plätze „Fehler in Folge", Platz 5 heißt OPEN, dazu das Probe-Feld („nur eine") und das Zustandsfeld. Kein Regelblock, die Regeln sagt der Trainer. Das Board ist das Dashboard: Das sieht nur der Betrieb, nicht die Kundschaft
- Das Trainer-Skript (ein Blatt Papier, beidseitig): vorn Aufbau, Ansage mit den vier Regeln und den Ergänzungen bei STOPP 1 und 2, die fünf STOPPs, „Wenn es hakt" und Vereinfachungen, hinten die Kartentabelle 01–27
- Aus dem Raum: ein Tisch für alle abgelegten Karten (ein Haufen)

Gesamtliste: [`docs/orga/materialliste.md`](../orga/materialliste.md), Druckvorlagen und Trainer-Skript unter <https://larmic.github.io/workshop_microservices/material/>.

Spielbar ab vier Personen (Booking und drei Requests, damit während der Probe zwei für 15 und 16 da sind), derselbe Stapel für jede Gruppengröße.

### Rollen

1. **Der Trainer** spielt keine Rolle. Er erklärt die Regeln, nimmt die STOPP-Karten, leitet das Gespräch und sagt „30 Sekunden sind um".
2. **Booking** ist eine Person. Sie hält den Stapel 01–27 und die drei Zustandskarten, ruft die Nummer und gibt die Karte dem Nächsten. Booking schaut aufs Board, ruft „OPEN!", tauscht die Zustandskarte und räumt ab. Sonst liest Booking nichts vor.
3. **Alle anderen** sind Requests auf `GET /booking/offers` und stehen in einer Schlange vor Booking. Flight, Hotel und Car spielt niemand, ihr Verhalten steht auf den Karten.

### Die Regeln

Der Trainer sagt sie am Anfang einmal (stehen auch auf Folie 15 und im Dashboard). Bei jeder Karte, immer in dieser Reihenfolge:

1. **Booking ruft die Nummer** und gibt die Karte dem Nächsten in der Schlange. Die nächste erst, wenn die vorige liegt.
2. **Warten, was auf der Karte steht.** Bei „5 s" oder „3 s" laut zählen: „einundzwanzig, zweiundzwanzig, dreiundzwanzig". Bei „sofort" nicht.
3. **Die Antwort laut vorlesen.**
4. **Karte ablegen, hinten anstellen.** Weiße Karten kommen auf den Haufen.

Dunkle Karte: Booking ruft „STOPP" und gibt sie dem Trainer. Wer noch zählt, zählt fertig.

**Bei STOPP 1 dazu** (das Board kommt, CLOSED): Rot heißt, Flight hat einen Fehler geliefert. Rote Karten kommen aufs Board, auf den nächsten freien Platz. Liegt Platz 5, ruft Booking „OPEN!" und legt OPEN ins Zustandsfeld.

**Bei STOPP 2 dazu** (30 s sind um, HALF_OPEN): Gelb ist die Probe. Die Karte kommt zuerst aufs Probe-Feld, dann zählen und vorlesen, danach auf den Haufen. Booking hört zu: „Timeout!" heißt OPEN, die 30 s laufen von vorn. „Flight antwortet!" heißt Board abräumen, CLOSED.

Weil die gelbe Karte zuerst aufs Probe-Feld kommt, liegt sie schon: Booking gibt 15 und 16 sofort aus, ohne Ausnahme von Regel 1. Die Karten 21–25 brauchen keine neue Regel, rot heißt Board.

### Die Karten

| Karte | Farbe | Wartezeit | Vorlesen | Danach |
|---|---|---|---|---|
| 01–03 | weiß | 5 s | „500, keine Hotels, keine Autos" | Haufen. Noch kein Breaker, das Board liegt noch nicht |
| 04 | dunkel | – | STOPP 1 | Trainer. Booking legt das Board hin, CLOSED |
| 05–09 | rot | 3 s | „`flights: []`, Hotel + Car da" | Board, Platz 1–5. Bei 09 ruft Booking „OPEN!" |
| 10–12 | weiß | sofort | „`flights: []` + `X-Circuit-Open`" | Haufen |
| 13 | dunkel | – | STOPP 2 | Trainer, danach HALF_OPEN |
| 14 | gelb, PROBE | 3 s | „Timeout!", `flights: []` | Probe-Feld, dann Haufen. Booking: OPEN, 30 s von vorn |
| 15, 16 | weiß | sofort | „`flights: []` + `X-Circuit-Open`" | Haufen (Probe-Feld belegt) |
| 17 | dunkel | – | STOPP 3 | Trainer, danach HALF_OPEN |
| 18 | gelb, PROBE | sofort | „Flight antwortet!", alles da | Probe-Feld, dann Haufen. Booking räumt 05–09 auf den Haufen, CLOSED |
| 19 | weiß | sofort | „alles da", Flight, Hotel, Car | Haufen |
| 20 | dunkel | – | STOPP 4 | Trainer |
| 21–25 | rot, 31.02. | sofort | „`flights: []`, Hotel + Car da" | Board, Platz 1–5. Bei 25 ruft Booking „OPEN!" |
| 26 | weiß, 14.03. | sofort | „`flights: []` + `X-Circuit-Open`" | Haufen. Gültiges Datum, trotzdem abgewiesen |
| 27 | dunkel | – | ENDE | Trainer. Board bleibt liegen (21–25, OPEN) |

### Durchlauf mit 10 Personen

Booking und die Requests P1 bis P9.

| Karte | an | Was passiert | Board | Zustand danach |
|---|---|---|---|---|
| 01–03 | P1, P2, P3 | 5 s zählen, „500, keine Hotels, keine Autos", Haufen | – | kein Breaker |
| 04 | Trainer | STOPP 1. Booking legt das Board hin, CLOSED | 0/5 | CLOSED |
| 05 | P4 | 3 s, „`flights: []`, Hotel + Car da", rot: Platz 1 | 1/5 | CLOSED |
| 06–08 | P5, P6, P7 | wie 05, Plätze 2–4 | 4/5 | CLOSED |
| 09 | P8 | wie 05, Platz 5. Booking: „OPEN!", OPEN hinlegen | 5/5 | OPEN |
| 10–12 | P9, P1, P2 | sofort, „`flights: []` + `X-Circuit-Open`", Haufen | 5/5 | OPEN |
| 13 | Trainer | STOPP 2. Gespräch, dann „30 Sekunden sind um": HALF_OPEN | 5/5 | HALF_OPEN |
| 14 | P3 | gelb: Probe-Feld, 3 s zählen, „Timeout!". Booking: OPEN, 30 s von vorn. Karte auf den Haufen | 5/5 | OPEN |
| 15, 16 | P4, P5 | während P3 zählt: sofort, `X-Circuit-Open`, Haufen (Probe-Feld belegt) | 5/5 | OPEN |
| 17 | Trainer | STOPP 3. Gespräch, dann „30 Sekunden sind um, Flight ist zurück": HALF_OPEN | 5/5 | HALF_OPEN |
| 18 | P6 | gelb: Probe-Feld, „Flight antwortet!". Booking räumt 05–09 auf den Haufen, CLOSED | 0/5 | CLOSED |
| 19 | P7 | sofort, „alles da" | 0/5 | CLOSED |
| 20 | Trainer | STOPP 4. Vorhersage zum 31.02., nicht auflösen | 0/5 | CLOSED |
| 21–24 | P8, P9, P1, P2 | 31.02., sofort, rot: Plätze 1–4 | 4/5 | CLOSED |
| 25 | P4 | 31.02., Platz 5. Booking: „OPEN!" | 5/5 | OPEN |
| 26 | P5 | 14.03., sofort, `X-Circuit-Open`: gültig und trotzdem abgewiesen | 5/5 | OPEN |
| 27 | Trainer | ENDE. Board bleibt liegen (21–25, OPEN) | 5/5 | OPEN |

Booking gibt der Reihe nach aus, deshalb wandert die Schlange gleichmäßig. Nur bei der Probe 14 ändert sich die Reihenfolge: P4 und P5 sind sofort zurück, P3 erst nach dem Zählen. Ob P3 vor oder nach P5 wieder ansteht, hängt vom Tempo ab. An Karten und Antworten ändert das nichts. Bei 8 Personen ist jeder bis zu viermal dran, bei 12 genau zweimal.

### Wo wir reden

| STOPP | Karte | Trainer fragt | Erwartete Antwort | Recap |
|---|---|---|---|---|
| STOPP 1 | 04 | Was habt ihr bekommen, wie lange? Hotel und Car waren gesund, was hattet ihr davon? Dann: Booking bekommt einen Circuit Breaker und 3 s Timeout. Regel für Rot und Platz 5 ansagen | Story 3 bricht beim ersten Fehler mit 500 ab, Hotel und Car werden gar nicht gefragt. Das Board ist das Dashboard, das sieht der Betrieb | – |
| STOPP 2 | 13 | Am Board: Welche Karten liegen hier, ab welcher Nummer gab es `flights: []`, was ist ab 10 anders? Was davon war der Breaker, was der Fallback? Einmal sagen: „Das war immer HTTP 200 mit `X-Fallback`." Dann „30 Sekunden sind um", Regel für Gelb ansagen | Fallback ab dem ersten Fehler, auch bei CLOSED. Der Breaker spart nur das Warten: OPEN fragt Flight nicht, darum `X-Circuit-Open` (09 hat Flight noch gefragt). Hotel und Car haben eigene Breaker | 02 |
| STOPP 3 | 17 | Wer lag auf dem Probe-Feld? Was haben 15 und 16 bekommen, wie schnell? Was, wenn alle zu Flight gegangen wären? Dann: „30 Sekunden sind um, Flight ist zurück." | Genau ein Probe-Slot (CompareAndSwap), der Rest sofort Fallback. Nach der gescheiterten Probe laufen die 30 s von vorn | 03 |
| STOPP 4 | 20 | Vorhersage, nicht auflösen: Flight ist gesund, die nächsten haben sich vertippt, Datum 31.02. Kommen die aufs Board? | Handzeichen, weiterspielen. Erst bei ENDE auflösen | – |
| ENDE | 27 | War Flight bei 21 bis 25 kaputt? Wer wurde bei 26 bestraft? Dann: Eine zweite Booking-Instanz hat ein eigenes Board, was liegt darauf? Was weiß Booking nach einem Neustart? | 4xx ist kein Ausfall, unser Code zählt alles ab 400. Der Zustand lebt im RAM einer Instanz. Geteilter Zustand hilft, kostet aber Kopplung | 01, 04 |

Die Fragen stehen nicht auf den STOPP-Karten, sondern im Trainer-Skript und in den Notes zu Folie 15.

### Bewusste Vereinfachungen

Offen ansagen:

- Der Request wartet selbst: 5 s Timeout in Story 3, 3 s in Story 4. Booking gibt die nächste Karte erst aus, wenn die vorige liegt. Im Code laufen Requests parallel.
- Die 30 Sekunden vergehen im Gespräch. Im Code wechselt der Zustand beim nächsten Call oder beim nächsten Blick ins Dashboard.
- Die Referenz zählt auch die gescheiterte Probe (Karte 14) mit, das Dashboard zeigt dann 6/5. Im Spiel geht 14 vom Probe-Feld auf den Haufen, das Board ist schon voll.
- HTTP 200 und `X-Fallback` stehen nicht auf den Karten. Der Trainer sagt sie einmal bei STOPP 2, die Demo zeigt sie.
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
- **Chaos-Steuerung**: Backend-Services (Flight, Hotel, Car) bieten `/admin/chaos` mit den Modi `normal`, `slow` (2 s, unter dem Timeout), `fail` (sofort 500) und `hang` (10 s, über jedem Timeout im Stack). Das Dashboard erlaubt das Schalten pro Service, sodass HALF_OPEN-Recovery live demonstriert werden kann.

### Workshop-Drehbuch

1. Stack starten: `docker compose -f services/docker-compose.yml -f services/docker-compose.infra.yml -f services/docker-compose.reference.yml up -d`
2. Dashboard öffnen (`http://localhost/dashboard`). Alle drei CBs zeigen CLOSED.
3. Flight im Dashboard auf "Fehler" stellen. Sechsmal `curl -i localhost:8088/booking/offers` schicken — jede Antwort ist 200 mit `flights: []` und `X-Fallback: flight`, der fünfte Call öffnet den Flight-CB, ab dem sechsten kommt `X-Circuit-Open: flight` dazu.
4. Logs des `booking-ref-story4` Containers ansehen (`docker compose logs -f booking-ref-story4`) — jede Entscheidung des Breakers ist sichtbar.
5. Flight zurück auf "Normal". Nach 30 Sekunden geht der CB beim nächsten Call automatisch in HALF_OPEN, der Probe-Call schließt den Breaker wieder.
