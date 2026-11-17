# Review vom 17.11.2026: Story 4, Rollenspiel „Wenn der Flug ausfällt“

Zwei offene Punkte aus dem Durchgehen des Spiels. Es wurde noch nichts geändert, dieses Dokument hält nur die Aufgaben fest.

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

**Zu entscheiden:**

- **21–25:** Verrät „Flight: 400“ die Antwort auf STOPP 4 („Kommen die aufs Board?“)? Die Frage zielt darauf, ob ein 4xx ein Ausfall ist. Steht 400 auf der Karte, wird das bei ENDE leichter zu besprechen, die Vorhersage bei STOPP 4 aber auch. Alternative: „Flight: Fehler“.
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

## Reihenfolge

Punkt 2 zuerst, weil er die Wartezeiten auf den Karten festlegt. Danach Punkt 1, dann den Stapel einmal neu drucken.
