# Review vom 17.11.2026: Story 4, Rollenspiel „Wenn der Flug ausfällt“

Drei offene Punkte aus dem Durchgehen des Spiels. Es wurde noch nichts geändert, dieses Dokument hält nur die Aufgaben fest.

> **Pfad-Konvention:** Folien ohne Pfad liegen in `services/slides/chapters/`. „Dashboard“ meint `services/dashboard/static/index.html`, „Druckseite“ meint `pages/material/index.html`, „Story-Doc“ meint `docs/stories/story-04-circuit-breaker.md`.

**Schon erledigt (Commit `73a3751`):** „Booking zieht die oberste Karte, sagt ihre Nummer und gibt sie dem Nächsten“ statt „Booking ruft Nummern“, Kartentext 01–03 „500 · Flight hängt, Hotel + Car weg“. Punkt 1 baut darauf auf und ersetzt diesen Kartentext.

**Aufwand:** klein = unter 30 Minuten, mittel = ein bis drei Stunden.

---

## 1. Karten klarer: erst Flight, dann der Rest

**Befund:** Auf den Karten steht nur, was die Kundin bekommt („`flights: []`, Hotel + Car da“, „`flights: []` + `X-Circuit-Open`“). Was mit Flight los ist, muss man sich aus Untertitel und Farbe erschließen. Beim Durchlesen war unklar, welcher Service eigentlich kaputt ist und was „500, keine Hotels, keine Autos“ bedeutet.

**Ziel:** Jede Karte liest sich in zwei Schritten. Erste Zeile: Was ist mit Flight? Zweite Zeile: Was bekommt die Kundin vom Rest?

**Vorschlag** (Wartezeiten schon mit Punkt 2):

| Karte | Wartezeit | Zeile 1: Flight | Zeile 2: der Rest | heute |
|---|---|---|---|---|
| 01–03 | 3 s | Flight: Timeout | 500, auch Hotel + Car weg | 500 · Flight hängt, Hotel + Car weg |
| 05–09 | 3 s | Flight: Timeout | Hotel + Car da, `flights: []` | `flights: []` · Hotel + Car da |
| 10–12, 15, 16 | sofort | Flight: nicht gefragt | Hotel + Car da, `X-Circuit-Open` | `flights: []` · + X-Circuit-Open |
| 14 (PROBE) | 3 s | Flight: Timeout! | Hotel + Car da | Timeout! · `flights: []` |
| 18 (PROBE) | sofort | Flight antwortet! | alles da | Flight antwortet! · alles da |
| 19 | sofort | Flight antwortet | alles da | alles da · Flight, Hotel, Car |
| 21–25 (31.02.) | sofort | Flight: 400 | Hotel + Car da, `flights: []` | `flights: []` · Hotel + Car da |
| 26 (14.03.) | sofort | Flight: nicht gefragt | Hotel + Car da, `X-Circuit-Open` | `flights: []` · + X-Circuit-Open |

**Geklärt:** „Flight: 400“ auf 21–25 verrät die Vorhersage bei STOPP 4 nicht, denn die Gruppe tippt, bevor Karte 21 gezogen wird. Bei ENDE hilft es, weil alle sehen, dass Flight gesund war.

**Zu entscheiden:**

- **Welche Zeile ist groß?** Heute ist die Antwort an die Kundin groß. Mit dem Vorschlag wäre Flight groß, die Kundin klein. Für die Pointe von STOPP 1 („Hotel und Car waren gesund, was hattet ihr davon?“) muss Zeile 2 trotzdem gut lesbar bleiben.

**Betroffene Stellen:**

- Druckseite: `S4_TXT` und `S4_DECK` (Kartentexte), Kartentabelle im Trainer-Skript (Seite 2), Durchlauf mit 10 Personen
- Story-Doc: Abschnitt „Die Karten“ (Beschreibung und Tabelle), „Durchlauf mit 10 Personen“
- Folie 15 „Wenn der Flug ausfällt“ und Dashboard Story 4: Untertitel „die Antwort steht auf der Karte“ prüfen, ggf. „was mit Flight los ist, steht auf der Karte“
- Notes zu Folie 15: STOPP 2 und ENDE beziehen sich auf die Kartentexte

**Folge:** Der gedruckte Stapel (4 Karton-Bogen) muss neu gedruckt werden.

**Aufwand:** mittel (Texte, Kartenlayout prüfen, neu drucken)

---

## 2. Timeout in beiden Phasen 3 s

**Befund:** Karten 01–03 warten 5 s (Timeout von Story 3), ab Karte 05 3 s (Timeout von Story 4). Bei STOPP 1 ändern sich damit zwei Dinge gleichzeitig: Der Breaker kommt dazu und der Timeout sinkt. Das legt nahe, der Breaker hätte das Warten verkürzt; das tut er aber erst ab OPEN. Außerdem verwechselt man die 5 s leicht mit den 5 Plätzen auf dem Board (Schwelle `FailureThreshold = 5`), die nichts damit zu tun haben. Das Teilnehmer-Review vom 13.11. merkt die „5 s“ ebenfalls an.

**Ziel:** Ein Timeout von 3 s in Story 3 und Story 4. Bei STOPP 1 kommt nur der Breaker dazu.

**Betroffene Stellen:**

- Code: `services/booking/story3/main.go:23` (`Timeout: 5000` → `3000`), `services/booking/story3/api/openapi.yaml:116` (Beispiel im `/info`)
- Druckseite: Karten 01–03 (`"5 s"` → `"3 s"`, Zeilen 651–653), Durchlauf (Zeile 681), Kommentar Zeile 640, Spielregel „Bei 5 s oder 3 s laut zählen“ (Zeile 452), Vereinfachung „5 s Timeout in Story 3, 3 s in Story 4“ (Zeile 474)
- Folie 15, Notes: Regel „Bei 5 s oder 3 s“ (Zeile 40); STOPP 1 „Die 5 s sind der Timeout der Referenz …“ und „Booking bekommt einen Circuit Breaker und 3 s Timeout“ → nur noch „einen Circuit Breaker“
- Story-Doc: Material (Zeile 26), Regel 2 (Zeile 47), Kartentabelle (Zeile 63), Durchlauf (Zeile 84), Vereinfachungen (Zeile 120), STOPP-1-Zeile in „Wo wir reden“
- `docs/instructions/circuitbreaker.md:15`: „Ohne CB (Story 3) wartet `/booking/offers` 5 s“

**Zu entscheiden:**

- **Story 1:** `services/booking/story1/main.go:26` hat ebenfalls 5 s. Mitziehen, damit Story 1 bis 4 denselben Timeout haben, oder bewusst lassen? Im Spiel kommt Story 1 nicht vor.

**Aufwand:** klein (ein Wert im Code, Texte); Gates danach komplett, weil Go und OpenAPI betroffen sind.

---

## 3. Status-Codes der Downstream-Services auf den Karten

**Befund:** Die Karten nennen keine Status-Codes von Flight, Hotel und Car. Gerade die Pointen hängen aber daran: Bei 01–03 bricht Booking nach dem Flight-Timeout ab und fragt Hotel und Car gar nicht erst. Bei 21–25 antwortet Flight mit 400, ist also gesund, und der Breaker zählt es trotzdem. Ohne Status-Code auf der Karte muss der Trainer das erzählen.

**Ziel:** Auf jeder Karte steht irgendwo, was Flight, Hotel und Car geantwortet haben. Ob als eigene Zeile, kleine Leiste oder Teil von Zeile 1 aus Punkt 1, ist offen.

**Vorschlag** (eine kleine Leiste pro Karte, ergänzt Punkt 1):

| Karte | Flight | Hotel | Car | Booking an die Kundin |
|---|---|---|---|---|
| 01–03 | Timeout | nicht gefragt | nicht gefragt | 500 |
| 05–09 | Timeout | 200 | 200 | 200, `flights: []` |
| 10–12, 15, 16 | nicht gefragt | 200 | 200 | 200, `X-Circuit-Open` |
| 14 (PROBE) | Timeout | 200 | 200 | 200, `flights: []` |
| 18 (PROBE), 19 | 200 | 200 | 200 | 200, alles da |
| 21–25 (31.02.) | 400 | 200 | 200 | 200, `flights: []` |
| 26 (14.03.) | nicht gefragt | 200 | 200 | 200, `X-Circuit-Open` |

Damit wird auch „Das war immer HTTP 200 mit `X-Fallback`“ (Ansage bei STOPP 2) sichtbar, statt nur gesagt. Die Vereinfachung „HTTP 200 und `X-Fallback` stehen nicht auf den Karten“ im Story-Doc und im Trainer-Skript entfällt dann.

**Vor dem Umsetzen prüfen:**

- Fragt die Referenz von Story 3 Hotel und Car wirklich nicht mehr, wenn Flight hängt (Reihenfolge der Aufrufe in `services/booking/story3/handler/booking.go`)? Davon hängt „nicht gefragt“ bei 01–03 ab.
- Passt die Leiste auf eine A7-Karte, ohne dass Nummer, Wartezeit und Vorlesezeile kleiner werden?

**Betroffene Stellen:** wie Punkt 1, dazu die Vereinfachungen im Story-Doc und im Trainer-Skript (Druckseite) und die STOPP-2-Ansage in den Notes zu Folie 15.

**Aufwand:** mittel, am besten zusammen mit Punkt 1.

---

## Reihenfolge

Punkt 2 zuerst, weil er die Wartezeiten auf den Karten festlegt. Danach Punkt 1 und 3 zusammen, weil beide das Kartenlayout ändern, dann den Stapel einmal neu drucken.
