# Abschließendes Review: „Microservices richtig gemacht“

> **Pfad-Konvention:** Folien-Dateien ohne Pfad liegen in `services/slides/chapters/`. „Dashboard“ meint `services/dashboard/static/index.html`. Die übrigen Pfade sind relativ zum Repo-Root angegeben.

---

## 1. Gesamturteil

**Roter Faden (gut, an den Rändern brüchig):** Die These „Monolithen mit Netzwerkproblemen“ trägt den Workshop. Die drei Kapitel steigern sich (Fundament, ein Aufruf, mehrere Schritte), die Stories greifen echt aufeinander zurück (2 zu 3 und 6, 6 zu 7), und 16a ist ein starker Cliffhanger. Gebrochen wird der Faden an drei Stellen: Die Agenda verspricht CQRS, BFF und Zero-Downtime und liefert sie nicht. Die Schlussbotschaft zum Bulkhead widerspricht Story 5. Story 8 setzt Code aus dem Rollenspiel Story 7 voraus.

**Übungsmix (gut verteilt, Mechanik muss nachgeschärft werden):** 4 × Hands-on, 1 × Design und 3 × Rollenspiel plus Quiz wechseln sich ab Tag 1 nachmittags sauber ab. Drei Probleme bleiben: Tag 1 beginnt mit rund 165 Minuten Zuhören. Für ein Architekt:innen-Publikum gibt es zu wenig Design, der Service-Schnitt fehlt. Mehrere Spielmechaniken funktionieren so nicht (Story 4 Runde 3, Story 5 Pool, Story 7 Tisch-Metapher).

**Themen (Kern stark, drei Lücken):** Idempotenz und der Resilience-Block sind überdurchschnittlich gut. Es fehlen Service-Schnitt, Timeouts/Retries und eine Begründung für Datenhoheit. Dazu kommt das CQRS-Versprechen.

**Texte (klar, aber Widersprüche an Merkstellen):** Sprachlich gut und weitgehend deckungsgleich zwischen Folie, Dashboard und Docs. Die Widersprüche sitzen ausgerechnet dort, wo etwas hängen bleiben soll: 29a, die Pill „Grenze fällt weg“, der Fallback in Story 4. Dazu kommt ein Begriffsmix (Rollback/Kompensation/Storno, Übung/Rollenspiel).

**Notes (zweigeteilt):** Spiel-Folien und die Recaps zu Story 3 und 4 sind vorbildlich. Die Recaps zu Story 5 bis 8 sind dagegen Textwände mit 380 bis 560 Wörtern. An wichtigen Stellen fehlen Notes ganz (Story 3, 14a, 13a, Quiz, 05a, Einstieg).

**Fragen (provokant und gut, Transfer fehlt):** Die Labels als Behauptung zum Widerlegen funktionieren. Es fehlt aber jede „Und bei euch?“-Frage. Die Fragen zu Story 4 stammen noch aus der Coding-Variante, Story 6 bis 8 haben zu viele Fragen, die sich überschneiden.

**Spiele (Story 7 ist das stärkste):** Billiges Material, unabhängig vom Tech-Stack, mit eingebauten Vorhersagen. Die Mechaniken müssen nachgeschärft werden, damit sie das lehren, was Referenz und Recap danach sagen.

**Dashboard/Design (Highlight mit Schwachstelle):** Das Dashboard ist laut Feedback die Spaßmaschine. Optisch wirkt es aber wie ein anderes Produkt als die Folien: Systemschrift, Grau, Emoji, abweichende Titel. Die Story-Box ist der zentrale Schwachpunkt.

**Spaßfaktor (Mittelteil sehr gut, Ränder schwach):** Tag 1 vormittags und Tag 2 nachmittags sind die Energielöcher. Dort liegen auch die größten Hebel.

---

## 2. Top 10 Maßnahmen mit größtem Hebel

Status: `[ ]` offen, `[~]` in Arbeit, `[x]` erledigt. Erledigte Punkte bekommen den Commit-Hash dahinter.

- [x] **1. Schlussbotschaft 29a zum Bulkhead korrigieren.** `29a-drei-dinge.md:20-21` ordnet den Bulkhead dem Aufgerufenen zu und widerspricht damit 17-bulkhead, Recap 05 und der eigenen Note. Sieben Reviewer haben das unabhängig gefunden. Fünf Minuten Aufwand (Formulierung siehe 4.0). **Umgesetzt** mit eigener Formulierung („sitzen vor dem Call … Keins ersetzt das andere“), Überschrift beibehalten. Mitkorrigiert: Notes 14:58, 17:61, 19:65/67, `docs/instructions/circuitbreaker.md:236`, `docs/questions/story5.md` Frage 2. Dashboard nicht betroffen. (fd162c5)
- [x] **2. CQRS, BFF und Deployment entscheiden.** `03-agenda.md:36`, `25-story-07-fragen.md:58` und `30-was-fehlt.md:28-31` widersprechen sich, dazu sind 70 Minuten ohne Folie verplant. Entweder drei Diskussionsfolien bzw. eine Architektur-Kata bauen oder Agenda, Note und themen.md bereinigen. **Entschieden: Variante B, ehrlich streichen, ohne neue Folie.** Agenda 03 und Note 25 bereinigt, Blöcke 7, 9 und 10 in `themen.md` aus dem Ablauf genommen und als Reserve markiert, Verweise in `docs/questions/story3`, `story6`, `story7` angepasst. 30-was-fehlt bleibt unverändert. Die 70 freien Minuten bleiben bewusst Puffer (f12cfe0). Nebenbefund (offen): `docs/stories/story-07-choreography-saga.md:125` nennt ein Feature-Flag für den Mischbetrieb sync/async, das es in `services/booking/story7` nicht gibt.
- [ ] **3. Story-4-Rollenspiel an die Referenz angleichen.** In `15-story-04.md:28-32` fehlen der Fallback bei jedem Fehler, das Wegwischen der Strichliste bei Erfolg, die richtige Erwartung für Runde 3 sowie Nachschub und Rundenuhr. Sonst widerspricht das Spiel Recap 02 und der Demo.
- [ ] **4. Story 8 von Story 7 entkoppeln.** Das Event-Kriterium in `27-story-08.md:30` und `docs/stories/story-08-tracing.md:47-50` als optional markieren und den Hook in `26-titel.md:21` umschreiben.
- [ ] **5. Recap-Notes zu Story 5 bis 8 auf höchstens 150 Wörter kürzen** (`19`, `22`, `25`, `28-*-fragen.md`), Karten-Notes auf höchstens 100. Vorlage sind 13 und 16. Die Tiefe wandert nach `docs/questions`.
- [ ] **6. Fehlende Notes ergänzen:** `12-story-03.md`, `14a`, `13a`, `05a`, `01` bis `04`, Quiz `10f` bis `10j` sowie Frage 2 in `10-story-01-fragen.md`.
- [ ] **7. Mechanik von Story 5 und Story 7 reparieren.** Story 5 bekommt Request-Karten statt Personen (`18-story-05.md:20-32`). Story 7 wird zur Leiter Zuruf → Korb → Antwortzettel (`24-story-07.md:31`, `24a-stoerungen.md`).
- [ ] **8. Tag 1 vormittags aktivieren.** Live-Chaos-Opener nach 04, Design-Übung „Schneidet die Reiseplattform“ statt eines Teils von Block 5, Zeitbudget realistisch planen (`docs/themen.md:15-28`: es fehlen 45 Minuten, und es gibt keine Kaffeepause).
- [ ] **9. Story-Box im Dashboard umbauen und Story-Art und Zeit sichtbar machen** (Dashboard `868-958`, `STORY_META 2982-2999`), dazu ein einheitliches Vokabular Hands-on / Design-Session / Rollenspiel (siehe Abschnitt 5).
- [ ] **10. Fachliche Angriffsflächen schließen:** Readiness prüft keine Backends (`10-story-01-fragen.md:33`), Story 7 ist keine echte Choreografie (`23-titel.md:9`), Consul-Deregistrierung bei `docker stop` (`13-story-03-fragen.md:30`), Probe-Slot im CB-Pseudocode (`14a:15-17`), parent-id im traceparent (`26-tracing.md`-Note), Pill „Idempotenz > Ehrlichkeit“ (`22:32`).

---

## 3. Übungsmix-Tabelle

| Story | Typ | Dauer (belegt) | Bewertung |
|---|---|---|---|
| 1 Erst lauffähig, dann smart | Hands-on | 90 min | Solide. Es fehlt ein prüfbares „Fertig, wenn“, und das Doc weicht von der Folie ab. Davor liegen rund 165 Minuten Zuhören. |
| 2 Erst denken, dann tippen | Design-Session | ca. 55 min (Badge zeigt 20) | Gutes Ergebnisformat (Spalte „Idempotent“), das zweimal wieder aufgegriffen wird. Notes sind zu dünn, die Team-Vorstellung fehlt auf den Folien, die Umbuchung kommt nie wieder vor. |
| Quiz 10f bis 10j | Quiz | (in Story 2) | Starke Dramaturgie (Teapot, „Beides RESTful“). Keine Notes, keine Mechanik, Quiz 2 ist mehrdeutig. |
| 3 Services finden sich selbst | Hands-on | 60 min | Zwei Aufgaben (Resolver plus POST) in 60 Minuten. Keine Notes, kein Hook. Der Recap fragt die Brücke zu Story 2 nicht ab. |
| 4 Circuit Breaker | Rollenspiel | ca. 60 min (Badge zeigt 20) | Gute Idee. Die Mechanik weicht von der Referenz ab (Fallback, Zähler-Reset) und kommt ohne Nachschub zum Stillstand. |
| 5 Bulkhead | Rollenspiel | ca. 60 min | Starke Vorhersage plus Little's Law. Bei normaler Gruppengröße wird der Pool nie leer, die Regeln stehen nur in der Note. |
| 6 Alles oder nichts | Hands-on | 60 min | Gute Saga-Story. `/admin/sagas` ist nicht spezifiziert, Gegebenes und Aufgabe sind vermischt, und zu Retry und Idempotenz gibt es widersprüchliche Notes. |
| 7 Choreography | Rollenspiel | ca. 60 min | Das stärkste Spiel (drei Störungen). Die Metapher Tisch/Papier/Postfach ist widersprüchlich, es gibt zu viele Zuschauende, und „jede zweite Buchung“ macht das Ergebnis zufällig. |
| 8 Tracing | Hands-on | 60 min | Setzt Events aus Story 7 voraus. Der grep-Moment verpufft allein im Terminal. Log-Filter siehe Kritiker-Hinweise. |
| CQRS / BFF / Deployment | Vortrag (geplant) | 30 + 20 + 20 min | Es gibt keine Folien. Liegt direkt nach dem Mittagessen. |
| Abschluss | Plenum | 30 min | Sieben Fragen auf drei Folien, keine aktive Phase. |

**Summe laut Befund:** etwa 270 min Code, 180 min Rollenspiel, 55 min Design. Für Architekt:innen ist der Design-Anteil zu klein.

---

## 4. Befunde je Dimension

### 4.0 Übergreifende Befunde (mehrere Reviewer, zusammengeführt)

**[hoch] Bulkhead „nach innen“ in der Schlussbotschaft**
- **Ort:** `29a-drei-dinge.md:20-21`, Note `:40`; Gegenstellen `17-bulkhead.md:31-33`, `:61`, `19-story-05-fragen.md:53`, `:67`, `14-circuit-breaker.md:58`
- **Problem:** „Circuit Breaker und Timeout gehören nach außen, Bulkhead und Rate Limit nach innen.“ Folie 17 sagt „Outbound, nicht Inbound“, Recap 05 „Outbound ≠ Inbound“, die Referenz hat den Bulkhead im Aufrufer. Damit kippt der letzte Merksatz, was das Becher-Spiel gezeigt hat. Fachlich kommt hinzu: Serverseitige Bulkheads gibt es auch (Nygard, Azure), wir bauen sie nur nicht.
- **Vorschlag:**
  ```html
  <h3>Resilience gehört in den <span class="hl">Aufrufer</span>.</h3>
  <p>Timeout, Circuit Breaker und Bulkhead schützen den, der aufruft.
  Das Backend schützt sich selbst, mit Rate Limit und Load Shedding.
  Wer das vermischt, schützt nichts.</p>
  ```
  Note `:40` ergänzen: „Die Becher standen bei Booking, nicht bei Hotel.“ Note auf 17: „Serverseitige Bulkheads, etwa Pools pro Consumer, gibt es auch. Die bauen wir nicht.“ `14-circuit-breaker.md:58`: „Das ist Rate Limiting oder Load Shedding im Backend, kein Circuit Breaker.“

**[hoch] CQRS, BFF und Zero-Downtime: angekündigt, nicht geliefert, am Ende „nicht angeschaut“**
- **Ort:** `03-agenda.md:36`, `25-story-07-fragen.md:58`, `30-was-fehlt.md:28-31`, `services/slides/index.html:111-129`, `docs/themen.md:38-41`, `173-260`
- **Problem:** 70 Minuten sind in themen.md geplant, es gibt keine Folie. Note 25 baut eine „Brücke zu CQRS am Nachmittag“, die ins Leere führt. 30-was-fehlt nimmt das Agenda-Versprechen ausdrücklich zurück. Der Block liegt außerdem direkt nach dem Mittagessen, und am Ende fehlt eine aktive Phase.
- **Variante A (liefern):** Nach 28 bzw. 25b einen Stack mit drei kurzen Diskussionsfolien einfügen. (a) „CQRS: Lesen ist nicht Schreiben“, Pill „CQRS ≠ Event Sourcing“ bzw. Statement „Wer liest, muss nicht warten, bis geschrieben ist.“ (b) „BFF: ein Backend pro Client?“, Pointe „Booking ist schon eins“. (c) „Deployen ohne Downtime“: N & N+1, Expand/Contract, Feature Toggle, Rolling/Canary. Je drei Note-Bullets aus themen.md. Spaßigere Alternative: eine **Architektur-Kata** (25 min) mit drei Karten: „Meine Reisen soll auch bei Hotel-Ausfall funktionieren“, „Smart-Watch-App braucht nur Gate-Info“, „Flight-Team will Freitag 17 Uhr deployen“, danach drei Auflösungsfolien. Oder CQRS als Design-Session im Format von Story 2: „Entwerft GET /customers/{id}/trips: Woher kommen die Daten, was zeigt ihr bei Hotel down, wie alt darf eine Zeile sein?“ Story 8 direkt nach die Mittagspause legen. Auf 30-was-fehlt dann nur noch „Event Sourcing“: „Events als Quelle der Wahrheit. CQRS hatten wir, Event Sourcing ist der nächste Schritt.“
- **Variante B (streichen):** `03-agenda.md:36` → „Distributed Tracing, ein Rückblick und was wir bewusst nicht gemacht haben.“ In `25:58` den CQRS-Satz ersetzen durch „Brücke zu Story 8: Wer merkt es, wenn ein Event verloren geht?“ Die themen.md-Blöcke 7, 9, 10 als optional markieren. Fällt nur BFF weg, reicht in der Agenda „CQRS und Deployment ohne Downtime“.
- **Fachlich zusätzlich:** Note 25 rückt CQRS Richtung Eventing. Besser: „Ein Read-Model, das nur zuhört, ist ein guter Ort für Events. Nötig ist das für CQRS aber nicht.“

**[hoch] Story 8 verlangt Events aus Story 7, die niemand programmiert hat**
- **Ort:** `27-story-08.md:30`, `:36`, `:52`; `26-titel.md:21`; `28-story-08-fragen.md:38`; `docs/stories/story-08-tracing.md:47-55`, `124-133`; Dashboard `1924-1928`
- **Problem:** Story 7 ist ein Rollenspiel. Der eigene Service hat weder Events noch CB oder Bulkhead. Bei freiem Tech-Stack entsteht so die Frage „Müssen wir jetzt noch Events bauen?“
- **Vorschlag:** `27:30` → „Auch die DELETE-Kompensation aus Story 6 trägt ihn.“ Unter „Optional“: „Wer mag: die Kompensation als Event wie in der Story-7-Referenz, traceparent als Property.“ `27:36` → „Jede Logzeile eures Booking-Service trägt die trace_id. Flight, Hotel und Car übernehmen sie schon.“ Hook `26-titel:21` → „In Story 6 ist eure Buchung durch vier Service-Logs gewandert. Wer konnte einen einzelnen Vorgang sauber rekonstruieren?“ Recap 04 (`28:38`) → „Der Zettel auf dem Tisch hat keinen HTTP-Header. Wie kommt der Trace mit?“ Im Doc-Bezug `124-126` „werden im Trace sichtbar“ ersetzen durch „in der Referenz sichtbar“.

**[hoch] Story 4: Das Spiel lehrt beim Fallback das Gegenteil von Recap 02 und Referenz**
- **Ort:** `15-story-04.md:28-30`, `:44`; Dashboard `1615-1617`; `docs/stories/story-04-circuit-breaker.md:38-40`; `16-story-04-fragen.md:22-24`, `:50`; `services/booking/story4/handler/booking.go:221-224`
- **Problem:** Im Spiel gibt es `flights: []` erst ab OPEN. Recap 02 löst auf mit „Ja, schon beim ersten Fehler“, und die Referenz macht `markFallback` bei jedem Fehler. Die Note „Hotel und Car werden gar nicht gefragt, so macht es die Referenz“ gilt nur für Story 3 und widerspricht außerdem der User Story („Hotel und Mietwagen auch dann sehen“).
- **Vorschlag:** `15:29`: „Runde 2, mit Schalter: Jeder Timeout gibt einen Strich und trotzdem die Karte Hotel ✓, Car ✓, flights: [], nach drei Sekunden. Beim fünften Strich: Schild auf OPEN.“ `15:30`: „Ab jetzt fragt Booking Flight gar nicht mehr, die Karte kommt sofort.“ Note `15:44`: „so macht es die Referenz aus Story 3“. Note `16:50`: „Ihr habt es gespielt: Die Karte gab es schon vor OPEN, der Schalter hat nur das Warten gespart.“ Dashboard und Doc genauso anpassen.

**[hoch] Recap-Notes Story 5 bis 8 als Textwand**
- **Ort:** `22-story-06-fragen.md:56-65` (561 Wörter), `28-story-08-fragen.md` (495), `25-story-07-fragen.md:53-59` (478), `19-story-05-fragen.md` (383). Zum Vergleich 13 (123) und 16 (107)
- **Vorschlag:** Pro Frage ein fettes Stichwort und höchstens zwei Sätze, Ziel höchstens 150 Wörter, letzte Zeile „Lang: docs/questions/storyN.md“. Vor dem Kürzen mit docs/questions abgleichen. Beispiele:
  - 22/01: „**Kompensation scheitert:** Forward darf scheitern, Kompensation muss gelingen: Idempotenz, Retry mit Backoff, persistenter Log, Operator-Inbox. **Spicy:** Ohne Plan dafür ist es keine Saga, sondern Hoffnung.“
  - 25/Dirigent: „**Dirigent:** Ein Prozess mit fester Reihenfolge und Kompensation gehört an eine Stelle. Choreographie nur für einfache Fälle ohne Reihenfolge (Benachrichtigung, Read-Model). **Spicy:** Bei Choreographie kennt jeder jeden, nur steht es nirgends.“
  - Reserve-Fragen in 22:63-64 und 28 in eine Zeile: „Reserve: Eventing statt sync, Owner, s. docs/questions/story6.md“.

**[mittel/hoch] Story 7 heißt Choreografie, ist aber asynchrone Orchestrierung**
- **Ort:** `23-titel.md:9-10`; `23-choreography-saga.md:27`, `:62`; `24-story-07.md:31`; `25-story-07-fragen.md:58`; `29-zusammenfassung.md:35-36`; `services/booking/story7/handler/booking.go:259-290`
- **Problem:** `CompensationRequested` ist ein Command an einen bekannten Empfänger. Nach Richardson ist das keine Choreografie. Die Titelfolie sagt „Niemand dirigiert“, Note 25 sagt „Booking bleibt Dirigent“. Außerdem wartet Booking beim Webhook weiter synchron auf das 202.
- **Vorschlag:** `23-titel:10` → „Wir gehen heute den halben Weg: Booking dirigiert die Buchung, das Stornieren organisiert jedes Backend selbst.“ Optional `23-titel:9` → „Booking ruft nicht mehr an. Es legt einen Zettel hin und geht.“ `24:31` → „Runde 2, Kompensation per Event:“. Note auf 23-choreography: „CompensationRequested ist ein Command: Booking entscheidet, wer storniert. Echte Choreografie wäre ein Event wie TripBookingFailed {sagaId}, und jedes Backend entscheidet selbst.“ Note zu 29: „halbe Choreografie, die Kompensation ist entkoppelt“. Optional Variante B (TripBookingFailed) als Bonus im Spiel.

**[mittel] Pill „Grenze fällt weg“ sagt das Gegenteil der Bonus-Folie**
- **Ort:** `19-story-05-fragen.md:37`
- **Vorschlag:** „Engpass wandert, verschwindet nicht“ oder kurz „Engpass wandert“.

**[mittel] Derselbe Hook drei- bis viermal hintereinander**
- **Ort:** `19a:16`, `20-titel:21`, `21:51`; `23-titel:21`, `:23`, `23-choreography-saga:59`, `24:43`; `26-titel:21`, `27:52`; „Heute wird gespielt“ in `17-bulkhead:62`, `18:43`, `23-titel:23`
- **Vorschlag:** Einen Hook pro Block, immer auf der Titelfolie. `19a:16` wird zur Kapitelbrücke ohne Zitat. `21:51` → „Ihr kennt die Mechanik, jetzt baut ihr sie. 60 min inkl. Demo.“ `23-choreography:59` → „Recap 6/01: Was, wenn DELETE in 5xx läuft?“ `24:43` → „Zweimal dieselbe Saga, mit Dirigent und ohne. Achtet darauf, wer wartet.“ In `27:52` nur die Demo-Anweisung behalten. „Heute wird gespielt“ nur auf 18 und 24 lassen.

**[mittel] Uneinheitliche Story-Art und Zeit-Badges**
- **Ort:** Kicker `15:3`, `18:3`, `18a:3`, `24:3`, `24a:3` („Übung“); Dashboard `1370-1376`, `1600`, `1658`, `1824`, `2990-2996`; `docs/stories/story-05-bulkhead.md:3`; Badges `10d:9`, `15:9`
- **Problem:** „Übung“ bzw. „Rollenspiel“ wechseln je nach Ort. Die Badges zeigen bei Code-Stories den ganzen Slot, bei Spielen und Design nur die Phase (20 statt 55/60 min). Wer nach den Folien plant, verschätzt sich um bis zu 40 Minuten, und Spiele wirken wie Pausenfüller.
- **Vorschlag:** Feste Wörter Hands-on (1, 3, 6, 8), Design-Session (2), Rollenspiel (4, 5, 7), Quiz. Kicker „Story 4 · Rollenspiel“, Phasenfolien „Story 5 · Rollenspiel, Nachrechnen“ und „Story 7 · Rollenspiel, Störungen“. Regel: Die erste Folie zeigt den Slot (`15:9` → „≈ 60 min <small>davon 20 Spiel</small>“, `10d:9` → „≈ 55 min <small>davon 20 Flipchart</small>“), Phasenfolien zeigen „Teil 2 · ≈ 20 min“.

**[mittel] Die Kostenfrage aus 08b wird nicht eingelöst**
- **Ort:** `08b-zwei-fragen.md:20-21`; `29-zusammenfassung.md:9-56`
- **Vorschlag:** Note 08b: „Flipchart Preisliste aufhängen, nach jedem Recap ein Zettel (z. B. Consul = neuer SPOF, CB-State nur im RAM, Event kommt doppelt).“ In jeder Recap-Note: „- Kosten-Wand: ein Post-it pro Pill.“ Auf 29 pro Karte eine gedämpfte Zeile „Preis: …“ (CB „Zustand pro Instanz“, Saga „Zwischenzustände sind sichtbar“, Choreography „niemand hat den Überblick“). Callout: „Jedes Pattern löst ein Problem und kauft ein neues.“ Optional ein Tag „auch im Modulith“ (Health/Config, REST, CB, Bulkhead) bzw. „nur verteilt“ (Discovery, Saga, Choreography, Tracing). Pointe: „Die Hälfte der Patterns gibt es nur, weil wir verteilt haben.“

---

### 4.1 Roter Faden

**Stärken:** Die Titel-These kommt auf 08 zurück. Die Kapitelübergänge sind gut begründet, die Rückgriffe zwischen den Stories echt. 16a ist ein Cliffhanger, die Domäne zieht sich konsequent durch. Die Recap-Folien sind einheitlich gebaut. „Acht Stories, acht Sätze“ ist ein guter Abbinder.

**Befunde** (zusätzlich zu 4.0):
- **[mittel] Pattern vor dem Problem erklärt.** `index.html:72-77`, `87-90`, `103-110`; `14-circuit-breaker.md:14-37`; `23-choreography-saga.md:30-33`; `24-story-07.md:32`, `:47`. Feedback D1 („Wie würdet ihr lösen? vor der Lösung“) wird so nicht eingelöst. **Vorschlag:** Story 4 Runde 1 direkt nach 14-titel spielen (15 im Stack vor 14-circuit-breaker), Runden 2 und 3 als Fragment erst nach 14a. Bei Story 7 die Karte „Booking bleibt haftbar“ aus 23 nach 24a oder in den Recap verschieben. Story 5 bleibt.
- **[mittel] Tag 2 hat keinen Einstieg.** Zwischen 16a und 17-titel (`index.html:78-83`). **Vorschlag:** Eine Folie „Tag 2“: links die Tag-2-Spalte der Agenda, rechts „Gestern offen: Flight-Breaker OPEN, die Kundin klickt Buchen. Was würdet ihr tun?“. Zwei Minuten Murmelrunde, die Antworten am Flipchart parken und in Story 6 aufgreifen. Note 17-titel: „Gestern war Flight tot, und der Schalter hat geholfen. Heute ist Hotel nur langsam, und der Schalter merkt nichts.“
- **[mittel] Zwei Gliederungen ohne Abbildung aufeinander.** `03-agenda.md:28`, `08a:9`, `13a:5-9`, `19a:9`. **Vorschlag:** Auf jeder Kapitel-Folie eine Story-Leiste 1 bis 8 im Stil des Steppers, in der Agenda das Kapitel als Kicker. Untertitel 08a: „Ein Service, der läuft, sauber spricht und gefunden wird.“ Agenda Tag 2 Vormittag: „Langsam, und dann halb gebucht“.
- **[mittel] Zeitbudget Tag 1 Vormittag.** `docs/themen.md:15-28`; `03-agenda.md:16`; `docs/orga/einladung-email.md:12`. Geplant sind 255 Minuten in einem Fenster von 210, eine Kaffeepause gibt es nirgends. **Vorschlag:** Geschichte auf 25 Minuten, Thesen-Diskussion auf 20. Story 1 bewusst über die Mittagspause legen (Start etwa 11:15, „Story 1 (Ende)“ im Nachmittags-Slot). Pro Halbtag 15 Minuten Kaffeepause eintragen.
- **[mittel] Domäne und Systemlandschaft werden nie eingeführt.** Zwischen 08a und 09. **Vorschlag:** Folie „Unser Spielfeld“ (Booking in der Mitte, Flight/Hotel/Car, Consul, Traefik, Dashboard). Note: „Das Dashboard ist euer Cockpit: Story lesen, Spickzettel, Chaos-Schalter. Jetzt kurz prüfen: Lädt localhost bei allen?“ Optional eine Kundin mit Namen, einheitlich in 10d, 16a, 24 und 24a.
- **[mittel] Die Klammer zum Anfang fehlt.** `05a` (keine Note), `08:32`, `29:52-56`, `themen.md:49`. **Vorschlag:** Auf 29 die Karten nach „Langsam“ (Bulkhead), „Lügt“ (Health, Status-Codes) und „Kommt nie an“ (CB, Saga, Choreography, Tracing) gruppieren, Discovery zum Fundament. Note 05a: „Mit KI wird Code billig, Architekturentscheidungen nicht. Genau die üben wir.“ Note 29 als ersten Punkt: „Erwartungs-Flipchart vom Morgen holen und gemeinsam abhaken.“
- **[mittel] Die Brücke von Recap 1 zu Discovery reißt.** `10-story-01-fragen.md:32-34`, `11-titel`, `12-story-03`, `13a`. **Vorschlag:** Note 10: „Backend nicht erreichbar: Die URL steht fest in der Env-Variable. Tot oder umgezogen? Parken wir bis Story 3.“ Note 11-titel: „Erinnert ihr euch an Frage 2 von heute Vormittag? Vielleicht ist Flight nicht tot, sondern hat nur eine neue Adresse.“ Note 13a: „Discovery sagt, wo ein Service ist. Ob er antwortet, sagt sie nicht.“
- **[niedrig] Reihenfolge am Einstieg.** `index.html:40-48`, `06a:9`. Optional umsortieren: 04, 05/05a, 08/08b, 08a, 06/06a, 09.
- **[niedrig] Schwaches Ende.** `index.html:122-129`, `30-was-fehlt`, `30a:20`. **Vorschlag:** 30-was-fehlt (Kicker „Für den nächsten Workshop“) vor 29, dann 29, 29a, 30a. 02-feedback am Ende ein zweites Mal einbinden und die Note 30a:20 anpassen.
- **[niedrig] Reste vor dem Umbau.** `16:32` Pill → „einer probiert, alle warten“. Note `16:49` → „Die Referenz zählt alles ab 400, in Produktion wäre das ein Bug.“ `09:58` → „Resilienz kommt in Kapitel 2 (Stories 4 und 5), Fehler über mehrere Schritte ab Story 6.“
- **[niedrig] Zusammenfassung Karte 01.** `29:12` → „Health, Config aus der Umgebung, Logs auf stdout. Der Rest baut darauf.“ Auf 06a ist Stateless ausgegraut.
- **[niedrig] Vorstellung der Story-2-Ergebnisse fehlt.** Note 10d: „Danach je Team 2 Min am Flipchart: ein Endpoint, auf den ihr stolz seid, und eine bewusste Abweichung. Frage in die Runde: Was passiert bei euch, wenn der Storno zweimal kommt? Dann das Quiz.“
- **[niedrig] Gleiche Aufstellung Story 4/5 nicht genutzt.** Note `18:43`: „Gleiche Aufstellung wie gestern bei Story 4. Neu: Hotel ist nicht tot, sondern langsam, und Booking hat keinen Schalter, sondern einen Becher.“

### 4.2 Übungsmix und Spiele

**Stärken:** Gute Abwechslung, nie zwei Code-Blöcke direkt hintereinander, jeder Halbtag endet ohne Laptop. Billiges Material, Vorhersagen vor jeder Auflösung, Demo nach jedem Spiel. Story 7 mit den Störungen, die Quiz-Dramaturgie und das Ergebnisformat von Story 2 funktionieren.

**Befunde:**
- **[hoch] Story 4 Runde 3: Die Erwartung ist falsch.** `docs/stories/story-04-circuit-breaker.md:42`, `:48`; `15-story-04.md:31-32`, `:47`; Dashboard `1621-1622`; `circuitbreaker.go:183`. Bei Erfolg wird der Zähler auf 0 gesetzt, ein wackliges Flight öffnet den Schalter also nie. **Vorschlag:** Regel ergänzen: „Antwortet Flight, wischt Booking die Strichliste weg.“ Erwartung: „Flight antwortet jedes zweite Mal: Der Schalter geht nie auf, jede zweite Anfrage wartet trotzdem drei Sekunden. Fünf Fehler in Folge sieht ein wackliges Backend nie.“ Optional eine Runde 3b mit Fehlerquote im Fenster („von den letzten zehn Strichen fünf rot, also OPEN“).
- **[hoch] Story 4 kommt zum Stillstand.** `15-story-04.md:28-31`, `:44-46`; Doc `:38-41`. **Vorschlag:** Auf die Folie: „Wer sitzt, steht nach zehn Sekunden wieder auf und stellt sich neu an.“, „Jede Runde dauert 90 Sekunden.“, „Eine Person zählt die beantworteten Anfragen.“, „Nach 30 Sekunden dreht Booking beim nächsten Gast auf HALF_OPEN.“ Optional den Schalter als eigene Person vor Flight stellen.
- **[hoch] Story 5: Der Pool wird nie leer.** `18-story-05.md:20-24`, `:32`, `:45-46`; Doc `:32-48`. **Vorschlag:** Request-Karten (F/H/C gemischt, ein Drittel H). Bei Hotel bleiben Karte und Chip liegen, die Person stellt sich sofort mit der nächsten Karte an. Alternativ Chips = Request-Personen minus eins. Auf die Folie: „Runde 2: Flight 4, Car 3, Hotel 3 Chips“, „6. Kein Chip frei? „503“, zurück ans Ende der Schlange.“ Dazu ein Hinweis für kleine Gruppen. `18:30` → „Jede:r schreibt vorher auf einer Karteikarte auf: …“ (es gibt keine Teams).
- **[niedrig] Story 5: 20 s ohne Timeout.** `18:24`, `:43`. Folie: „Hotel ist langsam, nicht kaputt: zwanzig Sekunden, kein Fehler, kein Strich.“ Note: „Das CLOSED-Schild aus Story 4 neben den Becher stellen. Es steht die ganze Runde auf CLOSED, der Breaker sieht Fehler, keine Latenz.“
- **[niedrig] Story 5: Runde 2 widerspricht Little's Law.** Note 18a: „Rückgriff: In Runde 2 hatte Hotel nur drei Chips. War das richtig? Und was, wenn Hotel gar keine sechzig Slots verkraftet?“ Die Selbst-Verteilen-Runde mit L = λ × W als optionale Kür einplanen.
- **[hoch] Story 7: Die Metapher ist widersprüchlich.** `24-story-07.md:21`, `:31`; `24a-stoerungen.md:35-37`, `:47-48`; Doc `:56-62`. Der Tisch persistiert wie ein Broker, die Referenz ist aber ein Webhook. **Vorschlag:** Eine steigende Leiter. Runde 2: „Booking ruft Flight den Storno quer durch den Raum zu.“ Störung 1: Flight ist draußen, der Zuruf verhallt. Danach den Korb mit Quittung einführen. Störung 2: Die Quittung geht verloren, der Korb liefert doppelt, die Zettel bekommen Nummern. Störung 3: der Antwortzettel. Merksatz: „Zuruf vergisst. Das Postfach erinnert sich, manchmal doppelt. Wer wissen will, ob es geklappt hat, braucht einen Antwortzettel.“
- **[mittel] Story 7: Zuschauende und Zufall.** `24:20-24`; Doc `:35-39`; `25:18-19`. **Vorschlag:** „Hotel lehnt in diesem Spiel jede Buchung ab.“ Zwei bis drei Kundinnen parallel. Beobachtende als „Kundenservice-Hotline“, die nur antworten darf, was sie am Tisch und bei Booking sieht. Ab 12 Personen zwei Tische. Recap 01: „Booking schickt den Storno per HTTP-POST, Flight ist gerade down. Wo ist das Event jetzt?“
- **[mittel] Quiz ohne Notes und Mechanik, Quiz 2 mehrdeutig.** `10f` bis `10j`; `10g:9`, `:13`; Plan `plan-d8-rest-vs-restful.md:67`, `153-154`, `185-186`. **Vorschlag:** Note pro Quiz: „- Daumen hoch = RESTful, runter = nicht, seitlich = kommt drauf an. - Eine Stimme aus der Minderheit begründen lassen, dann klicken.“ `10g:9`: `<span class="dim">Kunde 7 existiert</span>`, Note: „Gibt es Kunde 7 nicht, ist 404 richtig, genau diese Unterscheidung ist der Punkt.“ Optional ein Punktestand der Flipchart-Teams.
- **[mittel] Quiz 4 „stornieren ≠ löschen“ und DELETE in Story 6 werden nicht verbunden.** Note `21:54`: „Rückblick auf Story 2: Kompensation ist kein Kunden-Storno. Die Reservierung wird technisch zurückgenommen, bevor die Kundin sie bestätigt bekommt, deshalb passt DELETE. Für den Storno einer fertigen Reise bleibt POST …/cancellation richtig.“ Optional auf der Folie: „zurücknehmen, nicht stornieren“. (Zum Retry-Teil siehe 4.4.)
- **[hoch] Tag 1: fast drei Stunden Zuhören, zu wenig Design.** `themen.md:15-22`; `30-was-fehlt.md:44`; Backlog D2. **Vorschlag:** Block 5 zur Design-Session „Schneidet die Reiseplattform“ umbauen (siehe 4.3). Optional eine Aufstellung im Raum Monolith ↔ Microservices in Block 1 und noch einmal bei 29a.
- **[mittel] Abschluss ohne Spiel.** **Vorschlag:** 15 Minuten Postmortem-Spiel. Der Trainer liest eine kombinierte Störung vor, die Teams legen die Pattern-Karten aus 29 dazu und nennen die Kosten, mit der Kosten-Wand als Spickzettel.
- **[niedrig] Code-Stories: kein Pairing, kein Aufholen.** Note-Zeilen: „- Paare nach Sprache, gemischte Erfahrung.“, „- Pflicht ist das erste Kriterium, der Rest ist Kür.“ In 21: „- Kein POST aus Story 3? Mit dem Nachbar-Paar zusammengehen oder die Saga als Pseudo-Code entwerfen und gegen die Referenz demonstrieren.“
- **[niedrig] Code-Stories ohne Spielmoment.** Optional fünf Minuten Chaos-Runde: Story 3 eine Flight-Replica stoppen, Story 6 Hotel auf Fehler („Wessen Saga kompensiert sauber?“), Story 8 als „Trace-Detektiv“. Nach Story 6 optional „Chaos-Karten“: Wer nach 60 s ein sauberes FAILED mit kompensiertem Flug zeigt, bekommt einen Punkt.
- **[niedrig] Notes Story 2 zu dünn.** `10d:48-49`. Ergänzen: „- Hook: „Bei uns hat ein Browser-Prefetcher über einen GET-Link Buchungen storniert. Heute sorgen wir dafür, dass das nie wieder passiert.“ - 3 min Aufgabe, 20 min Flipchart, je Team 2 min Vorstellung. - Erwartbar für Storno: DELETE, POST …/cancellation, PATCH status. Umbuchung: PUT, PATCH, Storno plus Neubuchung. - Nach 10 min ohne Idempotenz-Spalte: „Was passiert beim zweiten Klick?““ Weiter: „Jetzt ihr gegen fünf Endpoints, Handzeichen: RESTful oder nicht?“

### 4.3 Themen

**Stärken:** Das „Ob“ wird ehrlich gestellt (04, 08b, 29a), Idempotenz zieht sich vorbildlich durch, der Resilience-Block ist tief (Little's Law, non-blocking). Messaging-Basics sind gut dosiert, Reduktionen werden offen benannt, die Story-Kette ist problemgetrieben.

**Befunde:**
- **[hoch] Service-Schnitt/DDD nur als Ausblick-Card.** `30-was-fehlt.md:12-15`, `:44`; `08b:37`; Feedback D2. **Vorschlag:** Design-Übung „Wo gehört das hin?“ (20 bis 25 min) mit 6 bis 8 Karten: Storno-Gebühr, Kundenprofil, Gesamtpreis, Treuepunkte, Zahlung, Zimmerverfügbarkeit, Reiseunterlagen-PDF, Umbuchung. Die Teams legen sie auf Flight/Hotel/Car/Booking oder einen neuen Service. Pointe: Gesamtpreis und Storno fallen nirgends sauber hin. Danach mit dem Ist vergleichen: „Wir schneiden nach Anbietern. Ist das gut? Wo landet Zahlung?“ Platz: nach 08/08b, die Zeit aus Geschichte (45 → 30) und Block 5 (30 → 20) nehmen. Note 30-was-fehlt:44 ersetzen: „Den Schnitt haben wir euch vorgegeben. Im echten Projekt ist das die erste und teuerste Entscheidung.“
- **[mittel] Timeouts und Retries fehlen.** `15-story-04.md:28-32`; `16-story-04-fragen.md`; `14a:24`; `20a:35`. **Vorschlag:** Runde 1 schärfen: „Abgewiesene Anfragen stellen sich sofort wieder an, und Booking selbst fragt Flight dreimal hintereinander.“ Erwartung: „Die Schlange wächst statt gleich zu bleiben. Das ist der Retry-Sturm, den der Breaker in Runde 2 stoppt.“ Recap-Frage: „Booking retried dreimal, das Gateway auch dreimal. Wie viele Calls bekommt Flight?“ (9, bei drei Schichten 27). Callout: „Retry ohne Backoff ist ein DDoS auf das eigene Backend.“
- **[mittel] Datenhoheit falsch etikettiert.** `30-was-fehlt.md:24-27`; `20-saga.md:16-18`; Backlog B8/F2. **Vorschlag:** Card auf 20-saga: „Jeder Service besitzt seine Daten. Keine geteilte DB, also kein gemeinsames COMMIT.“ Note: „Geteilte DB = geteiltes Schema = gemeinsames Deployment = verteilter Monolith. Im Workshop in-memory, bewusst reduziert.“ Card auf 30 umbenennen in „Daten zwischen Services: Wer besitzt welche Daten? Replikation per Event, Read-Models, Polyglot Persistence.“
- **[niedrig] Contract Testing fehlt ganz.** `30:20-23` → „Verträge & Versionen: Breaking Changes finden, bevor sie deployt sind. Consumer-Driven Contracts, OpenAPI-Diff, N & N+1 parallel.“
- **[niedrig] SAML.** `08b:38`, `themen.md:119` → „OIDC am Rand, Token weiterreichen oder tauschen, mTLS zwischen Services“.
- **[niedrig] Observability.** Callout auf 26-tracing oder 27: „Metriken sagen, dass etwas passiert. Traces sagen, wo. Logs sagen, was.“
- **[niedrig] Outbox als Produkt-Pill.** Die Pill aus `23-choreography-saga:50` entfernen und den Halbsatz in der Note `:65` streichen. Auf 25b im Code-Tag „Behalten“ Outbox ergänzen.
- **[niedrig] Sync vs. async.** Callout auf 25: „Sync oder async ist eine Transportwahl. Orchestration oder Choreography ist eine Verantwortungswahl.“
- **[niedrig] KI-Folie ohne Anschluss.** Note 05a: „- Provokation: Ein Agent schreibt euch einen Circuit Breaker in 30 Sekunden. Ob ihr einen braucht, wo er hingehört und welche drei Zahlen er bekommt, entscheidet er nicht. Darum geht es die nächsten zwei Tage.“ Optional in 29a: „Welchen der drei Sätze hätte euch ein Agent gesagt?“
- **[niedrig] Kulturwandel nur in den Notes.** Zwischen 29a und 30 eine numlist-Folie: „01 Wer von euch hätte heute Nacht Bereitschaft für Booking?“, „02 Wer betreibt Consul, Broker und Tracing-Backend?“, „03 Monolith, Modulith, Microservices: Was wäre euer Trigger?“. In der Note 30a ein Lesetipp: Team Topologies.
- **[mittel] Monolith-Check pro Pattern (B7 offen).** Auf 14-, 17-, 20-, 23- und 26-titel: „Monolith-Check (Daumen hoch/seitwärts/runter): Gibt es das Problem im Monolithen auch?“. CB seitwärts, Bulkhead seitwärts, Saga runter, Tracing seitwärts.

### 4.4 Fachliche Haltbarkeit

**Stärken:** Die CB-Zahlen, der Bulkhead-Semaphore und die Saga-Zustände passen zur Referenz. Die Consul-API ist korrekt, Little's Law ist korrekt gerechnet. Die Zuschreibungen stimmen (12-Factor 2011, Sagas 1987, RFC 2324, W3C 2020). Die Trennung von CB und Bulkhead ist sauber, der Webhook-Schwachpunkt in Story 7 wird ehrlich benannt.

**Befunde** (29a und Choreografie siehe 4.0):
- **[hoch] Readiness prüft Backends (Kaskaden-Anti-Pattern).** `10-story-01-fragen.md:33`. **Vorschlag:** „Readiness: Kann diese Instanz jetzt bedienen? Eigene Ressourcen wie Warm-up, eigene DB und Pool, nicht die Downstream-Services. Wer Flight in die Readiness nimmt, schaltet bei einem Flight-Ausfall den ganzen Booking-Service ab. Gegen kranke Backends hilft der Circuit Breaker, Story 4.“
- **[mittel] Der CB-Pseudocode erlaubt einen Probe-Storm.** `14a:15-21`; `14:22`. **Vorschlag:** „if state == HALF_OPEN and not probeSlot.tryAcquire(): return fallback()   // nur einer probt“. Kommentar „// genau ein Versuch“ und vor Z. 19 „failure in HALF_OPEN → state = OPEN, openUntil = now + 30s“. Karte 14: „CLOSED → OPEN → HALF_OPEN → CLOSED | OPEN“.
- **[mittel] Consul-Deregistrierung.** `13-story-03-fragen.md:30-32`; `11a:25`; Dashboard `1515`. Bei SIGTERM melden sich die Backends ab. **Vorschlag:** Frage 03: „docker stop oder docker kill: Wann verschwindet der Eintrag?“, Pill „stop: sofort · kill: bis zum nächsten Check“. 11a: „Das Interval entscheidet, wie lange nach einem Absturz Verkehr ins Leere läuft. DeregisterCriticalServiceAfter räumt nur den Katalog auf.“ Dashboard: „Läuft eine zweite Instanz, übernimmt sie, sobald der Health-Check der ersten fehlschlägt.“
- **[mittel] „Langsamere Backends brauchen mehr Slots“.** `18a:30`, `:41`. Callout: „Die Pool-Größe folgt aus normaler Latenz × Durchsatz, gedeckelt durch die Backend-Kapazität. Wird das Backend langsamer, ist der volle Pool das Signal, kein Grund zum Vergrößern.“ Note: „Wer auf die Timeout-Latenz dimensioniert (60), hat keinen Bulkhead mehr.“
- **[mittel] „Downstreams erzeugen niemals einen Trace“.** `26b`-Note, `28:21-23`, `27`-Note. Frage 28/02: „Ein Downstream bekommt einen traceparent und würfelt trotzdem eine neue ID. Was passiert?“, Pill „übernehmen, nie überschreiben“. Note: „Ohne Header darf jeder Service einen Trace starten, das macht OTel auch. Wir lassen es bei Flight, Hotel und Car bewusst weg, damit der Kontrast sichtbar bleibt.“
- **[mittel] parent-id falsch beschrieben.** `26-tracing.md`-Note → „Pro Hop wird eine frische Span-ID gewürfelt und als parent-id in den Header geschrieben. Der Empfänger sieht sie als seinen Parent. Die Trace-ID bleibt.“ und „parent-id (die Span-ID dieses Aufrufs)“.
- **[niedrig] Details zum traceparent.** 26b: „Feste Längen prüfen, Nur-Nullen und Version ff verwerfen.“ (Achtung: Die Referenz `tracing.go` lehnt alles außer 00 ab, das ist eine Code-Änderung.) `26:32` „Version 00: 55 Zeichen“. 26-titel „traceparent (plus optional tracestate)“. 28: „Wer Fehler-Traces immer will, braucht Tail-based.“ Note: „crypto/rand, weil IDs nicht vorhersagbar sein sollen.“ Das Dashboard nutzt traceparent im Response, standardisiert wäre traceresponse (Draft).
- **[mittel] Saga: „Timeout, aber gebucht“.** `20a:16-17`. Note: „Achtung: booked enthält nur Bestätigtes. Timeout beim Forward heißt Ergebnis unbekannt. Abhilfe: von Booking vergebene Buchungs-ID oder Idempotency-Key.“ Reserve-Frage 22: „Hotel antwortet nicht, hat aber gebucht. Wer storniert?“
- **[mittel] Idempotenz und Retry in Story 6.** `21:37`, `:54` (gegen `:55`, `20a:36`); `22:32`, `:60`. Note 21:54: „Heute wiederholt niemand. Sobald ein Retry dazukommt oder jemand manuell nachstößt, muss DELETE idempotent sein.“ Pill 22/03: „204 ist Pragmatik, nicht Idempotenz“. Note 22: „Idempotent sind beide Varianten. 204 erspart dem Orchestrator nur die Sonderlogik 404 = ok.“ Beginn Note 22:60: „Am Flipchart in Story 2 war 404 noch vertretbar. Jetzt ruft ein Orchestrator auf, der wiederholt …“ In `story2.md:17` und `rest-vs-restful.md:23`: „Welche Antwort ein wiederholender Aufrufer braucht, klärt Story 6.“
- **[niedrig] Atomarität versprochen.** `21:12`, `:21` und Dashboard `1735`, `1742`. User Story: „… die am Ende entweder vollständig steht oder vollständig storniert ist“. Untertitel: „Drei Buchungen, am Ende ein Ergebnis: alles da oder alles storniert.“ Dashboard: „… ersetzt die fehlende Transaktion durch Kompensationsschritte. Atomar ist das nicht.“
- **[niedrig] REST vs. RESTful nie aufgelöst.** Note 10a: „Auflösung: REST = Architekturstil (Fielding 2000, inkl. Hypermedia). RESTful im Alltag = HTTP-Semantik richtig nutzen, Richardson Level 2. Mehr verlangen wir nicht.“
- **[niedrig] Status-Codes.** `10b:56` → „Caches: speichern auch 404 oder 301, und leider auch den 200 mit Fehler drin“. Quiz 5 → „Die Infrastruktur liest ihn wie 400: Client-Fehler, nicht wiederholen.“ Note 10d → „Unerwartete 5xx muss man nicht einzeln dokumentieren, 503 mit Retry-After schon.“
- **[niedrig] Bulkhead „503“ und „Slots“.** `17a:15` → „return fallback()   // Referenz: flights: [] + X-Bulkhead-Full“. Dashboard `1662` → „hängen immer mehr Aufrufe gleichzeitig an Hotel: Threads, Connections, Speicher. Irgendwann fehlt das für Flight und Car.“
- **[niedrig] Infrastruktur-Defaults.** `17:63` → „max_requests bzw. max_connections“. `19b:46` → „Go: MaxConnsPerHost ist per Default unbegrenzt.“ `25b:36` → „Exactly-once-Zustellung gibt es nicht. At-least-once plus idempotente Verarbeitung ergibt exactly-once-Wirkung.“ `25b:51` → „Kafka-EOS gilt innerhalb von Kafka, nicht für den DELETE beim Hotel.“
- **[niedrig] Pauschale Schlusssätze.** 30a: „nach Melvin Conway, 1968“. Note: „… und wer das weiß, schneidet die Teams bewusst passend (Inverse Conway Maneuver).“ 04-foot: „Monolith first (Martin Fowler): ein Standpunkt, kein Gesetz.“ 29a-Note: „höhere Velocity“ streichen. Optional `29a:14` „der stärkste Grund“.
- **[niedrig] Zeitleiste.** `05:33`, `:97`: SCS auf etwa 2015 setzen. Auf `05a:3`, `28-29` Kicker und Fahne angleichen.
- **[niedrig] Nygard-„Zitate“.** `14-titel:11`, `17-titel:11` → „nach Michael Nygard, Release It! (2007)“.

### 4.5 Verständlichkeit

**Stärken:** Gleicher Wortlaut von Folie, Dashboard und Docs bei den Spielen. Quiz-Merksätze wie „leer ≠ weg“, Titelfolien mit festem Aufbau, Little's Law mit drei Zeilen. Die Pseudo-Code-Folien entsprechen den Spickzetteln, die Analogien sind gut gewählt.

**Befunde** (zusätzlich zu 4.0):
- **[hoch] Story 6: `/admin/sagas` nicht spezifiziert, Gegebenes und Aufgabe vermischt.** `21:34-41`; Dashboard `1740-1743`; `story-06-saga.md:21-22`; `docs/custom-setup.md`. **Vorschlag:** Überschrift „Gegeben: Flight, Hotel, Car bieten DELETE /bookings/{id}, idempotent: 204 auch beim zweiten Mal“. `21:41`: „GET /admin/sagas liefert eure Sagas mit Status PENDING … FAILED, Format im Spickzettel. Dann zeigt das Dashboard sie an.“ Beispiel-JSON im Spickzettel und Admin-Endpoints in custom-setup.md.
- **[mittel] Vier Begriffe für „zurücknehmen“.** `21:21`; `23:22`, `:32`; `story-06-saga.md:13`, `:20`; Dashboard `1735-1742`. „(Kompensation in umgekehrter Reihenfolge)“; `23:22` „Jedes Backend storniert selbst.“; `23:32` „Ob das Storno gelang, weiß es nicht.“
- **[niedrig] Denglisch.** Agenda `03:29`: „Die Saga räumt auf, wenn Flug, Hotel und Mietwagen nicht zusammenpassen, erst zentral gesteuert, dann per Event.“ Dashboard `1874`: „… wird die Kompensation als Event an die Backends verschickt, ohne auf Antwort zu warten.“ Optional eine Regel in CLAUDE.md: Service-Namen englisch, Fachbegriffe deutsch, Pattern-Namen englisch.
- **[mittel] Story-Titel weichen ab** (siehe 5 bzw. Dashboard). Zusätzlich `23-choreography-saga:7` → „Wer storniert, wenn keiner dirigiert?“ oder `24-story-07:7` → „Der Tisch ist der Broker“, damit „Die Saga wird leise“ nicht doppelt vorkommt.
- **[niedrig] Kryptische Pills.** `16:32` „genau ein Probe-Call“; `22:16` „Kompensation muss am Ende klappen“; `22:24` „kein Retry → dann Alarm“; `22:32` „zweimal stornieren = einmal“; `22:40` „hängende Saga ist still → Counter“; `25:37` „Broker hält fest, Webhook vergisst“; `28:16` „IDs zum Finden, Spans zum Messen“ bzw. „Logs mit Trace-ID: 80 % Nutzen“; `10:24` „umgezogen → Story 3, kaputt → Story 4“.
- **[niedrig] Kein „Fertig, wenn“.** „Fertig, wenn im Dashboard unter Story N mit Quelle „Custom“ alle Buttons grün sind.“ Story 6: „Hotel auf „Fehler“ ergibt eine Saga FAILED mit Flight COMPENSATED.“ Bei Story 3 Body und Antwort in die Liste aufnehmen. Das Doc zu Story 1 an die Folie angleichen.
- **[niedrig] 04-praxis.** Kicker: „Aus der Praxis: TN Skills, von Microservices zurück zum Modulith“ („TN“ wird sonst als Teilnehmende gelesen). „Fast 3 Teams“ → „2–3 Teams, mehrere Services“.
- **[niedrig] Story-1-User-Story.** `09:23`, Dashboard `1399` → „… damit wir Probleme früh sehen und dasselbe Image ohne Umbau in jeder Umgebung läuft.“
- **[niedrig] Zeitangabe Story 2.** Doc `:7` → „ca. 55 Minuten: Input 10, Flipchart 20, Vorstellung 10, Quiz 15“.

### 4.6 Speaker-Notes (Übersicht in Abschnitt 6)

**Stärken:** Die Spiel-Notes haben Hook, Material und Erwartung pro Runde. Überleitungen in Zitatform (17a, 20a, 26b), die Recaps 13 und 16 als Muster, 16a als Muster für eine Diskussionsfolie. 18a nennt die erwartete Fehleinschätzung, die Kapitel-Trenner 08a und 19a sind gut, die Klickabläufe werden genannt.

**Vorschlag für ein einheitliches Schema:** `Kern / Frage / Zeit·Klick / Weiter / Mehr`. Recap höchstens 150 Wörter, Karten höchstens 100, ein Hook pro Block auf der Titelfolie. Den Dashboard-Hinweis nur in 09 belassen. Backticks und `<code>` nicht mischen (09). Die interne TODO-Zeile in `28:65` entfernen und `docs/questions/story8.md` anlegen.

### 4.7 Fragen nach den Folien

**Stärken:** Labels als Behauptung zum Widerlegen, Story 7 nutzt die Sprache des Rollenspiels, 16a als offener Abschluss. Idempotenz wird quer verbunden, das Muster „meine Antwort“ plus „Spicy“, Vorhersagen in den Bonus-Folien, docs/questions als Nachschlagewerk.

**Befunde** (zusätzlich zu 4.0):
- **[mittel] Kein Transfer ins eigene Projekt.** Alle Recap-Folien. **Vorschlag:** Pro Recap eine letzte, abgesetzte Zeile „Und bei euch?“ ohne Pill. Story 1: „Was prüft der Health-Check eures wichtigsten Service wirklich?“ Story 4: „Welcher Aufruf in eurem System hat keinen Timeout?“ Story 7: „Wo verschickt ihr Events ohne Quittung?“ Dafür eine Fachfrage streichen.
- **[mittel] Story 4: Die Fragen 01, 03 und 04 kamen im Spiel nicht vor.** In Runde 3 eine Karte „Flug-ID 999“ einbauen, Flight antwortet sofort „404, gibt es nicht“, und Booking entscheidet, ob es einen Strich macht. Pill 03 → „einer probiert, alle warten“.
- **[niedrig] docs/questions/story4.md veraltet.** In `16:52` „Lang: docs/questions/story4.md“ ergänzen. Eine neue Frage „Kommt der Fallback erst bei OPEN?“ anlegen. Frage 1 neu einleiten: „Booking hat im Spiel bis fünf gezählt. Warum nicht bis eins oder bis fünfzig?“
- **[mittel] story5.md widerspricht dem Recap.** `docs/questions/story5.md:54-56` → „Pro Booking-Instanz sieht jedes Backend höchstens 10 gleichzeitige Aufrufe. Das begrenzt, was diese eine Instanz anrichten kann. Geschützt ist damit Booking, nicht das Backend, siehe Frage 8.“ Die Fragen 2 und 3 mit „(nur nach der Dashboard-Demo)“ kennzeichnen.
- **[mittel] Recap Story 7 wiederholt 24a.** Die Fragen 01 und 02 zusammenlegen: „Unsere Referenz verschickt per HTTP-POST. Welche der drei Störungen trifft sie, und welche nicht?“ (Pill „Blatt Papier, kein Postfach“). Die Verweise in `24a:48-49` anpassen. Falls 02 bleibt (`25:27`, bisher keine Frage): „Flight stirbt nach dem Storno, aber vor der Quittung. Der Broker liefert noch einmal. Was passiert, wenn auf dem Zettel <span class="hl">Rückerstattung</span> steht?“
- **[mittel] Story 1 Frage 02 ohne Antwort.** Note: „**Backend weg:** Zwei Fälle. Umgezogen (neue IP, Skalierung): Die feste URL zeigt ins Leere, das löst Story 3. Wirklich kaputt: Jede Anfrage hängt im Timeout, das löst Story 4. Heute fängt es niemand auf, mit Absicht.“ Zusätzlich in docs/questions/story1.md aufnehmen.
- **[niedrig] Story 3: Folie widerspricht docs.** `13:49-50` gegen `story3.md:30-35`. Auf eine Linie bringen: „Kurz weg: letzte Auflösung mit kurzer TTL weiterverwenden. Länger weg: ehrlicher Fehler statt Raten.“ Die Kubernetes-Frage in story3.md ergänzen.
- **[mittel] Zu viele, sich überschneidende Fragen.** Story 6: 01 und 02 zusammenlegen („Der Storno läuft selbst in 5xx, und wir versuchen es nur einmal. Was nun?“), 04 und 05 zusammenlegen („Booking stürzt mitten in COMPENSATING ab. Wer merkt es, und wann?“). Story 8: 01 und 05 zusammenlegen („Kein Span-Baum, kein OTel: Was fehlt uns?“), 03 wird Reserve.
- **[mittel] DELETE: `200` oder `204`?** *(vom Autor ergänzt, nicht aus dem Review-Lauf)* Die Frage gehört in den RESTful-Abschnitt und in die Fragen, ist aber bisher nicht ausdrücklich gestellt. `10b-rest-vs-restful.md:16` zeigt `DELETE … → 204 No Content` als gesetzt, und `10i-quiz-4.md` behandelt nur DELETE gegen POST `/cancellation`, nicht den Status-Code. **Vorschlag:** Die Frage „DELETE antwortet mit `200` oder `204`?“ aufnehmen (`200` mit Body, z. B. die gelöschte Ressource oder ein Storno-Beleg, gegen `204` ohne Body). Dazu gehört, was beim zweiten DELETE zurückkommt (`204` oder `404`, siehe Story 6 in 4.4). **Offen für die Umsetzung:** Form (eigenes Quiz, Frage im RESTful-Abschnitt oder Recap-Frage) und Ort (`10b`, Quiz `10f` bis `10j`, `22-story-06-fragen.md`, `docs/questions/story2.md`) werden erst dort entschieden.
- **[niedrig] Story 8 Frage 03 setzt Go voraus.** `28:30` → „Die Trace-ID steht im Log-Text. Warum reicht das nicht?“
- **[niedrig] Alle Pills mit einem Klick.** Den `data-fragment-index` pro Zeile hochzählen. Note: „Pro Frage ein Klick, der Merksatz erscheint erst nach der Diskussion.“
- **[mittel] Sieben Recaps im gleichen Format.** Methoden-Zeile in der Note: Recap 3 „Wette in Sekunden“, Recap 5 „Think-Pair-Share“, Recap 6 „eine Frage pro Team, 3 Min“, Recap 8 „Daumen-Abstimmung vor dem Klick“. In 19 und 22 nur drei Fragen live.

### 4.8 Dashboard und Design

**Stärken:** Kernpalette synchron, „gestrichelt = kein Code“ auf beiden Seiten. Story 2 ist schlank, die Spieltexte decken sich. Notes verweisen gezielt aufs Dashboard, Spielkarten sind standardmäßig offen, der Stepper hat Pfeiltasten und localStorage. Der Umschalter Referenz/Custom und die Demo zu Story 8 zum Anfassen funktionieren.

**Befunde** (Story-Box siehe Abschnitt 5):
- **[hoch] Story-Art und Zeit fehlen.** Dashboard `1364-1378` (nur Tooltips), `844-845` (gestrichelt auf solid zurückgesetzt), `2982-2999`. **Vorschlag:** `STORY_META` um `kind`, `kindLabel` und `minutes` erweitern und in `showStory()` Kicker plus Zeit-Pille rendern. Stepper-Knoten mit `<small>Code|Design|Spiel</small>`. Den gestrichelten Rand in #7348E1 zeichnen und beim aktiven Knoten erhalten (Z. 845 streichen). Drei Varianten (gefüllt/gestrichelt/gepunktet) nur zusammen mit `.badge.game` in `theme.css:323`.
- **[mittel] Drei Titelsysteme.** Dashboard `2983-2998`, Panel-h2 `1453`, `1556`, `1628`, `1688`, `1785`, `1858`, `1978`. **Vorschlag:** Folientitel als kanonischen Titel in `STORY_META` und im Doc-H1 übernehmen, `topic` als eigenes Feld. Panel-h2 einheitlich „Demo mit der Referenz · Circuit Breaker“, bei 1, 3, 6 und 8 „Ausprobieren · Service Discovery“.
- **[mittel] Beamer.** Badge und h4 0.65 bis 0.7rem in #888 (≈ 3,2:1), Body 0.85rem, die Chaos-Schalter liegen ganz unten (`2029-2032`). **Vorschlag:** #888 → #5B4E93, Body 0.95rem. Eine kompakte Chaos-Zeile über dem Demo-Panel für Story 4, 5 und 6 („Flight: normal | langsam | Fehler“). Spielkarten beim ersten Demo-Klick zuklappen. Header-h1 `1283` → „Microservices richtig gemacht“. Einen Beamer-Modus nur bauen, wenn 125 % Zoom das Layout bricht.
- **[mittel] Visuelle Sprache.** `11-13`, `879`, Emoji `1288`, `1303`, `1311-1327`, `2153`, `3159-3160`. **Vorschlag:** Poppins 400/600/800 plus OFL nach `static/fonts/`, `body{font-family:"Poppins",…;color:#0A0349}`, Story-Titel in 800. Story-Box in #EEE7FB. Emoji in Header und Werkzeugmenü durch ein Inline-SVG-Set ersetzen (z. B. Lucide). ✈ 🏨 🚗 im Status-Grid bleiben.
- **[niedrig] Lime-Abweichung.** `1033`, `1046-1049`, `1074`: #CCF465 → #D3F871, Schatten `rgba(211,248,113,0.55)`. Danach die Skill `sync-palette` prüfen lassen.
- **[mittel] Story-7-Legende verwässert die Pointe.** `1868-1874`, `1883`, `2767`, `490-493`. **Vorschlag:** Unterzeile: „Vorwärts weiter orchestriert: flight → hotel → car. Scheitert ein Schritt, legt Booking ein Storno-Event ab und antwortet sofort, ohne auf die Backends zu warten.“ Legende „✉ STORNO ANGEFORDERT“ mit eigener Klasse `.saga-step.requested` (gestrichelt #7348E1, weiß). Darunter: „Ob Flight wirklich storniert hat, sieht Booking nicht. Genau das war Störung 3.“ `2767`: „— Kompensation angefordert, ob sie geklappt hat, weiß Booking nicht“. Demo inszenieren: Hotel langsam plus Fehler, POST auslösen, Flight auf 0 skalieren. Note 24a: „Das Dashboard lügt wie das Blatt Papier.“
- **[mittel] Story-8-Filter im Custom-Modus.** `2812-2814`. **Vorschlag:** `const svc = BookingMode[8] === "custom" ? "booking-custom" : "booking-ref-story8";` Im Custom-Modus zusätzlich eine grep-Variante anzeigen. Die Funktionsnamen in `callTracing8` bzw. `resetSagas7` umbenennen.
- **[mittel] Dreifach gepflegte Inhalte driften.** Story 3 `1518-1524`: „<code>POST /booking/bookings</code> orchestriert die drei Backends und antwortet mit <code>201 Created</code>“, dazu „Aus Story 2 mitgenommen: Ressource bookings statt Verb, Fehler als Status-Code statt 200 mit Fehler-Body.“ und „Optional: Consul KV“. Story 6 `1737` „halbe Reise“. Dazu ein kleiner Konsistenz-Check im Stil von sync-palette (Titel, Zeit, Kriterienzahl).
- **[niedrig] User Story „Als System“.** `1830`, Doc `:16` → „Als <em>Betriebsteam</em> möchte ich, <em>dass eine fehlgeschlagene Buchung sofort beantwortet wird, auch wenn ein Backend beim Stornieren hängt</em>, damit <em>ein langsamer Flight-Service nicht jede Buchung blockiert</em>.“
- **[niedrig] Stepper ohne Fortschritt.** In `showStory()`: `document.querySelectorAll('.stepper-connector').forEach((c,i)=>c.classList.toggle('done', i < n-1));`, dazu `.visited{background:#EEE7FB;border-color:#7348E1}`.
- **[niedrig] Sprachgemisch in der UI.** Umschalter „Referenz | Mein Service“. `1640`: „… pro Backend ein eigener Circuit Breaker. Nach fünf Fehlern öffnet er und liefert sofort einen Fallback.“ `1711`: „… Abgewiesene Anfragen pro Pool werden sichtbar.“
- **[niedrig] CB-Demo ohne Burst.** `1627-1648`. Ein Button „10× feuern“ analog zu `bh-burst-btn`. Note `15:48` als Drehbuch: „1. Flight auf Fehler. 2. 10× feuern: fünf 500er, dann X-Circuit-Open. 3. Grid zeigt OPEN, Countdown laut mitzählen. 4. Flight auf Normal, nach HALF_OPEN wieder CLOSED.“

### 4.9 Spaßfaktor

**Stärken:** Das Dashboard als Spaßmaschine (Chaos-Knöpfe, Countdown, Burst), deutliche Reaktion auf das Kickoff-Feedback, Vorhersagen mit Aha-Effekt. Pointen wie Teapot, Vibe Coding und „dekoriert“, die TN-Skills-Geschichte, Cliffhanger und Hero-Folien. Tag 1 endet und Tag 2 beginnt mit einem Spiel.

**Befunde** (zusätzlich zu den obigen):
- **[mittel] Tag-1-Einstieg.** **Vorschlag:** (1) Live-Chaos-Opener in 5 Minuten nach 04: Flight auf Fehler/Langsam, Offers klicken, „Bis morgen Nachmittag baut ihr, dass das nicht weh tut.“ (2) Den Geschichtsblock mit einer Aufstellung entlang einer Zeitachse beginnen („Wo steht euer Projekt?“), das ersetzt die Vorstellungsrunde. (3) Den Setup-Check auf 10 Minuten in den Start von Story 1 legen. (4) Note 05 mit Erzähl-Ankern pro Epoche (1990 RMI/CORBA, 2000 SOA/ESB als Engpass, 2010 REST/JSON, 2014 Docker, 2019 Modulith, SCS) und „Wo steht ihr auf der Straße?“.
- **[niedrig] Kriegsgeschichten fehlen.** Note 13a optional: CrowdStrike Juli 2024 (Check-in weltweit down, Bordkarten von Hand). Note 19a optional: Southwest Dezember 2022. Zahlen vorher prüfen. Google Web Accelerator/Backpack auf 10h (siehe Quiz).
- **[niedrig] Story 8 ohne Vorher-Erlebnis.** Jedes Team bekommt einen Ausdruck mit gemischten Story-7-Logs von drei parallelen Buchungen: „Welche Buchung wurde kompensiert und warum?“ Drei Minuten auf der Stoppuhr. Nach Story 8 dieselbe Frage per grep in 10 Sekunden. Die Logs einmalig aus der Referenz erzeugen und ins Repo legen.
- **[niedrig] 14a-Untertitel.** `14a:7` → „Die Regeln, nach denen Booking gleich spielt.“ oder „In Pseudo-Code. Gleich spielen wir ihn.“ Den Spickzettel gibt es nicht mehr.

### 4.10 Weitere Hinweise (Kritiker-Durchgang, **nicht gegengeprüft**, vor Umsetzung verifizieren)

- **[hoch] Kein Auffangplan bei Rückstand in Code-Stories.** Die Stories sind kumulativ (`docs/custom-setup.md:7`), die Kotlin-Lösung deckt nur Story 1 ab. **Vorschlag:** Checkpoints `services/booking/custom-checkpoints/story3|story6|story8` oder Git-Tags. Note: „Wer hängt: CUSTOM_BOOKING_PATH auf custom-checkpoints/storyN-1 setzen, make docker-up-hub, dann von dort weiter.“ Auf der Story-6-Folie ein story-foot: „Startpunkt: euer Stand aus Story 3 oder Checkpoint story3“.
- **[hoch] Das Timeout-Szenario ist im Dashboard nicht zeigbar.** „Langsam“ = 2000 ms liegt unter dem 3000-ms-Timeout (`docker-compose.yml:15`, `chaos.go:125-130`), Story 3 hat 5000 ms. **Vorschlag:** Chaos-Modus „Hängt“ (z. B. 10 s) oder Flight-Latenz 4000. Demo mit `curl -w "%{time_total}\n"`. Spielregel bzw. Story-3-Timeout angleichen, `docs/instructions/circuitbreaker.md:13` korrigieren.
- **[mittel] Setup-Check ohne Folie, Custom-Service nie vorab im Stack getestet.** `vorbereitung.md:15`, `52-56`. **Vorschlag:** Schritt 2b in vorbereitung.md (CUSTOM_BOOKING_PATH setzen, `curl -i http://localhost:8099/health` liefert 200, Custom-Toggle aktiv). Folie „Startklar?“ mit vier Häkchen. Note: „Wer kein grün hat, setzt sich zu jemandem mit grün.“
- **[mittel] Trainer-Hinweise zum Quiz passen nicht zu den Folien.** `rest-vs-restful.md:30`, `108-152`: falsche Nummerierung, eine nicht existierende Folie 10c, „sechs Prinzipien“ statt zehn Regeln. Abschnitt 7 an 10f bis 10j angleichen. Prüfen, ob Quiz 2 gegen „POST /booking/bookings → 200 OK { "status": "error" }“ getauscht wird, Merksatz „Status-Codes sind Infrastruktur“.
- **[mittel] Trainer-Hinweise zu CB und Consul veraltet.** `circuitbreaker.md:9`, `157-190` (Coding-Variante, 30 s, „Echte Atomicity“). `service-discovery-consul.md:79`, `140`, `269` sind widersprüchlich. Auf das Rollenspiel umstellen, 3 s, „Saga: kein Alles-oder-nichts, sondern ein definierter Weg zurück“, „Kritisch nach spätestens 10 s, aus der Registry nach 1 Minute“.
- **[niedrig] Keine Trainer-Hinweise für Story 1, 6 und 7.** `docs/instructions/saga.md` und `twelve-factor-health.md` anlegen, danach die Recap-Notes kürzen.
- **[mittel] Log-Filter in Story 8 bricht.** Das `docker compose logs`-Präfix lässt jq scheitern, ohne die `-f`-Liste fehlen Services. **Vorschlag:** `make trace ID=…` mit `--no-log-prefix` und `jq -cR "fromjson? | select(...)"`, oder `COMPOSE_FILE` in `.env.example`.
- **[mittel] Story 3 enthält zwei Aufgaben, der Recap fragt POST nicht ab.** „Teil A · Resolver (30 min)“ und „Teil B · Buchen (30 min)“. Recap-Frage: „Hotel sagt Nein. Was antwortet euer POST: 201, 200 oder 5xx, und was ist mit dem schon gebuchten Flug?“
- **[niedrig] Die Umbuchung aus Story 2 verpufft.** Reserve-Frage in 22 oder 25: „Umbuchung ist eine Saga: neuen Flug buchen, alten stornieren. Welche Reihenfolge, und was passiert, wenn das Storno scheitert?“
- **[mittel] KI-Nutzung ungeregelt, keine Zusatzaufgaben.** Block „Schon fertig?“ auf jeder Hands-on-Folie. Regel auf der Story-1-Folie: „KI ist erlaubt. Bevor ihr Custom auf grün stellt, erklärt ihr eurem Nachbar-Team in zwei Sätzen, was euer Code bei einem Ausfall tut.“ Das schließt auch an 05a an.
- **[mittel] Kontrast der Quiz-Merksätze.** `theme.css:625-630` (#8880AE ≈ 3,6:1). `.quiz-take` → #5B4E93, mindestens 20px. Ebenso `.muted` und `.compare-note`. `.dim` Opacity 0.28 → 0.45. Am Beamer im Raum prüfen.
- **[niedrig] Packliste und Druckvorlagen fehlen.** `docs/orga/materialliste.md` (Mengen für 12 Teilnehmende) und Rollenkarten zum Drucken (vorne der Name, hinten zwei bis drei Regeln).
- **[niedrig] Story-6-Doc nennt Webhook-Callbacks und Choreography.** `story-06-saga.md:27-28`, `:34` → „Ihr baut die Orchestration: Booking ruft nacheinander auf und kompensiert selbst. Die Alternative ohne Dirigent spielen wir in Story 7.“
- **[niedrig] Feedback „anonym“, aber der Name bleibt gespeichert.** `pages/feedback/index.html:225-376`. Namen nicht persistieren oder nur per Checkbox. Unterzeile: „Anonym, außer du trägst unten einen Namen ein“. Dropdown „Bezieht sich auf“. Für die Discord-Webhooks im ausgelieferten HTML ein Relay mit Rate-Limit vorschalten oder die Webhooks nach jedem Durchlauf rotieren.
- **[niedrig] Titelfolie ohne offiziellen Namen.** `01-intro.md:7-8`: Kicker „Microservices richtig gemacht“ ergänzen.
- **[niedrig] Die Ausschreibung verspricht Service-Grenzen.** `docs/orga/intranet.md:13` → „Was passiert, wenn Services ausfallen, langsam werden oder sich widersprechen: Discovery, Circuit Breaker, Bulkhead, Saga, Tracing, hands-on und im Rollenspiel.“ Oder die Schnitt-Übung wirklich einplanen.

---

## 5. Story-Box im Dashboard: Redesign

**Ist-Zustand** (Dashboard `868-884`, `916-924`, `940-958`, `1392ff`): Bis zu drei `<details>` stehen in einer Flex-Zeile (`flex: 1 1 320px`, `align-items: flex-start`). Eine geöffnete „Story lesen“ wird zu einer etwa 350px schmalen Textsäule, daneben bleiben zwei graue Balken und leere Fläche. Die Summary heißt überall gleich, das Badge (0.65rem, #888 auf #f5f6f8, „User Story + Akzeptanzkriterien“) sagt nichts aus. Kicker, Story-Art, Zeit-Pille und Unterzeile fehlen. Die User Story steht als Fließtext statt als dl-Karte wie auf der Folie (`theme.css:461-485`). Bei den Spielen (`1599`, `1657`, `1823`) ist sie offen, rund 25 bis 30 Zeilen Textwand, und die Demo rutscht unter den Fold. Optisch wirkt das wie ein FAQ-Akkordeon.

**Variante C, minimaler Fix (zuerst, sofort sichtbar):**
- `.story-helpers details[open]{flex-basis:100%;order:-1}`, die geöffnete Box also in voller Breite.
- User Story als `<dl class="story-user">` (Als / möchte ich / damit) auf Lavendel #EEE7FB, wie auf der Folie.
- Das Summary-Badge durch eine Art-/Zeit-Pille ersetzen: „Hands-on · 90 min“ bzw. gestrichelt „Rollenspiel · 60 min“.
- Body 0.95rem, #888 → #5B4E93.

**Variante A, Folien-Spiegel (empfohlen als zweiter Schritt):**
- Story-Kopf aus `STORY_META` gerendert: Kicker „Story 4 · Rollenspiel“, Folientitel in Poppins 800, Zeit-Pille, Unterzeile.
- Darunter ein dreispaltiges Story-Grid statt Akkordeon-Leiste.
- Gleiche Titel wie auf den Folien (siehe 4.8).

**Variante B, Tabs mit typabhängigem Start-Tab und Phasenanzeige:** Für Live-Workshops nicht nötig und mit zusätzlicher Bedienlogik verbunden. Eher nicht.

**Speziell für die Rollenspiele (4, 5, 7):**
- `open` entfernen. Darüber eine kompakte Spielkarte „So spielt ihr“: Rollen als Chips (Booking · Flight · Hotel · Car · Requests), Runden als drei nummerierte Pills, eine Leitfrage. Langtext eingeklappt bzw. als Link zum Story-Doc.
- Optional ein Spielwerkzeug: Story 4 ein 30-s-Countdown mit Signal „HALF_OPEN!“ plus 90-s-Rundenuhr (übernimmt „Du bist die Uhr“), Story 5 ein 20-s-Hotel-Timer und ein 503-Zähler.
- Die Minuten aus den h4 (`1607`, `1831`, `1845`) streichen, nur noch Kopf-Pille.

**Begleitend:** Die Chaos-Zeile direkt über dem Demo-Panel (4.8) und die Spielkarten beim ersten Demo-Klick zuklappen. Reihenfolge: 1) Art und Zeit vereinheitlichen und in `STORY_META` aufnehmen, 2) Variante C, dann A mit Poppins und Lavendel, 3) Beamer-Tauglichkeit, 4) kleine Text- und Bugfixes.

---

## 6. Speaker-Notes: Status je Folie

| Folie | Status | Vorschlag (knapp) |
|---|---|---|
| 01-intro | fehlt | „Vorstellungsrunde: Name, Rolle, eine Erwartung (Flipchart).“ |
| 02-feedback | fehlt | „Nur anteasern, Formular am Ende.“ Am Schluss erneut einbinden. |
| 03-agenda | fehlt | Blockzeiten aus themen.md als Zeile. Kaffeepausen ergänzen. |
| 04-praxis | fehlt | „Erzählen: kein Aushelfen über Teams, Pipeline-Overhead. Frage: Wer hat Microservices wegen der Teams, wer, weil es modern war? Kommt am Schluss zurück.“ |
| 05 Geschichte | zu dünn | Ein Bullet pro Station plus „Wo steht ihr auf der Straße?“ |
| 05a KI | fehlt | „Kern: KI macht Code billig, nicht Architektur. Frage: Welcher Brocken liegt in eurem Projekt? Weiter: 2011 hat das jemand aufgeschrieben.“ |
| 06a 12-Faktor | Fokus falsch | III, VII und XI erklären (statt VIII/XII), um die geht es beim Klick. |
| 08a, 19a | gut | 19a:16 Hook streichen (Doppelung). |
| 08b | ergänzen | Anweisung zur Kosten-Wand. |
| 09 Story 1 | Format | Backticks und `<code>` vereinheitlichen. `:58` Satz zu Resilienz aktualisieren. Dashboard-Hinweis nur hier lassen. |
| 10 Recap 1 | lückenhaft, fachlich falsch | Frage 2 beantworten. Readiness ohne Backends. |
| 10a | ergänzen | Auflösung REST/RESTful (Fielding, Richardson Level 2). |
| 10b | Lücken | Fragment 2 bekommt data-notes, dazu ein Slide-Note für den Ausgangszustand. |
| 10d Story 2 | zu dünn | Hook Prefetcher, Timing 3/20/2 min, erwartbare Varianten, Überleitung zum Quiz. |
| 10f bis 10j Quiz | fehlen | Mechanik (Daumen, Stimme aus der Minderheit), auf 10g Kunde-7-Hinweis, auf 10h Backpack/Google Web Accelerator, auf 10i DELETE vs. POST /cancellation. |
| 11-titel | fehlt | „Erinnert ihr euch an Frage 2? Vielleicht ist Flight nicht tot, sondern hat nur eine neue Adresse.“ |
| 11a Consul | fehlt | passing=true heißt raus nach etwa einem Intervall (10 s). Deregister nach 1 m ist Aufräumen. Graceful Shutdown. |
| 11b | ok ohne | Optional in die Note von 11a eindampfen. |
| 12-story-03 | fehlt | Hook, erster schreibender Endpoint (201), Referenz, Demo (vorher prüfen), Time-Box 60 min. |
| 13 Recap 3 | Muster, fachlich korrigieren | `docker stop` vs. `kill`. Linie zu story3.md. |
| 13a Kapitel 2 | fehlt | „Kapitel 1 ließ den Service laufen, ab jetzt: was tun, wenn die anderen nicht laufen. Kaputt → CB, langsam → Bulkhead. Zwei Spiele, kein Code.“ |
| 14-titel | fehlt | Hook aus 15:42 hierher. „Wer hatte schon mal einen Threadpool voller hängender Calls?“ Monolith-Check. |
| 14a CB-Code | fehlt | Drei Zahlen, HALF_OPEN genau ein Probe-Call, Referenzpfad, Weiter: „Das tippen wir nicht ab, wir spielen es.“ |
| 15 Story 4 | anpassen | Hook kürzen. „so macht es die Referenz aus Story 3“. Demo-Drehbuch. |
| 16 Recap 4 | Muster | CompareAndSwap in die Note, „Lang: docs/questions/story4.md“. |
| 16a | Muster | Behalten. |
| 17-titel | fehlt | Hook aus 18:43 hierher. „500 sofort oder 20 s?“ Brücke zu Tag 2. |
| 17 Bulkhead | zu lang (Karten-Prosa) | Fünf Bullets (Pool, Fail-Fast, Nicht CB, Markt, Weiter). Serverseitige Bulkheads erwähnen. |
| 17a | gut | Behalten. |
| 18 Story 5 | anpassen | „Gleiche Aufstellung wie Story 4. Neu: …“ CLOSED-Schild neben den Becher. |
| 18a | gut | Rückgriff auf die 3 Chips und Satz zum Timeout-Dimensionieren ergänzen. |
| 19 Recap 5 | zu lang (383) | Höchstens 150 Wörter, Methoden-Zeile. |
| 20-titel | aufgebläht | „Die Grafik“ → „Auf den Lime-Pfeil zeigen: die Kompensation.“ |
| 20 Saga | zu lang | Karten-Notes auf höchstens 100, Datenhoheit ergänzen. |
| 20a | ergänzen | Fall „Timeout, aber gebucht“. |
| 21 Story 6 | widersprüchlich | Retry-Satz `:54` korrigieren. Hook kürzen. Pairing/Aufholen. |
| 22 Recap 6 | zu lang (561) | Kürzen, Fragen zusammenlegen, Pill zu 204/Idempotenz korrigieren. |
| 23-titel | aufgebläht | Zeigeanweisung statt Grafikbeschreibung, Hook-Doppelung entfernen. |
| 23-choreography | zu lang, doppelter Hook | Hook → Rückbezug Recap 6/01. Note Command vs. Event. |
| 24 Story 7 | kürzen | „Zweimal dieselbe Saga, mit Dirigent und ohne. Achtet darauf, wer wartet.“ |
| 24a | ergänzen | „Das Dashboard lügt wie das Blatt Papier.“ |
| 25 Recap 7 | zu lang (478) | Kürzen, CQRS-Satz klären. |
| 26-titel | anpassen | Hook auf Story 6, Grafik-Bullet streichen. |
| 26 Tracing | fachlich falsch, lang | parent-id korrigieren, crypto/rand-Begründung. |
| 26b | fachlich | Entry-Point als Workshop-Entscheidung, Versionsbehandlung. |
| 27 Story 8 | anpassen | Hook ohne Story-7-Events, Demo-Drehbuch, Filter-Befehl. |
| 28 Recap 8 | zu lang (495), TODO | Kürzen, story8.md anlegen, TODO-Zeile entfernen. |
| 29 | ergänzen | Erwartungs-Flipchart abhaken, nur eine Frage. |
| 29a | Widerspruch | Bulkhead-Satz anpassen, keine Frage, „höhere Velocity“ streichen. |
| 30-was-fehlt | Backlog-Notiz | `:44` durch den Trainer-Satz zum Schnitt ersetzen. |
| 30a | anpassen | Conway-Frage plus „Was nehmt ihr mit?“, dann Feedback-Folie (ohne Rücksprung). |
| CQRS/BFF/Deployment | fehlen ganz | Je drei Bullets oder Blöcke streichen (siehe 4.0). |

---

## 7. Themen

**Sollte rein:**
- **Service-Schnitt** als 20- bis 30-Minuten-Design-Übung „Wo gehört das hin?“ / „Schneidet die Reiseplattform“. Das ist der Top-Wunsch aus dem Feedback und wird so auch in der Ausschreibung versprochen.
- **Timeouts und Retries** (Retry-Sturm, Backoff), eingebaut in Story 4 Runde 1 plus eine Recap-Frage (9 bzw. 27 Calls).
- **Datenhoheit** als Card bzw. Note auf 20-saga („keine geteilte DB, also kein gemeinsames COMMIT“).
- **CQRS** (und ggf. Deployment ohne Downtime), wenn es in der Agenda bleibt: als Design-Session oder Kata mit Folien.

**Als Ausblick (Card oder Satz genügt):**
- Verträge & Versionen (Consumer-Driven Contracts, OpenAPI-Diff, N & N+1)
- Security: OIDC am Rand, Token weiterreichen oder tauschen, mTLS
- Observability: „Metriken sagen, dass. Traces sagen, wo. Logs sagen, was.“
- Sync vs. async als Transportwahl gegenüber Orchestration vs. Choreography als Verantwortungswahl
- Kulturwandel und Plattform als sichtbare Diskussionsfolie, Team Topologies als Lesetipp
- Event Sourcing (auf 30-was-fehlt)
- „Auch im Modulith / nur verteilt“ als Tag auf der Zusammenfassung

**Bewusst weglassen oder kürzen:**
- BFF (laut themen.md ist Booking schon eins), allenfalls als Pointe
- Consul-Feature-Liste 11b (in die Note von 11a eindampfen)
- Byte-Details auf der Tracing-Konzeptfolie
- KI-Folie 05a, falls sie nicht per Note und 29a mit dem Rest verbunden wird
- Outbox als Produkt-Pill auf 23 (nach 25b verschieben)

---

## 8. Quick Wins (je unter 15 Minuten)

- [ ] `29a-drei-dinge.md:20-21` Bulkhead-Satz ersetzen, Note `:40` ergänzen
- [ ] `19-story-05-fragen.md:37` Pill → „Engpass wandert, verschwindet nicht“
- [ ] `10-story-01-fragen.md:33` Readiness-Note korrigieren und Antwort auf Frage 2 ergänzen
- [ ] `21-story-06.md:54` Retry-Widerspruch auflösen
- [ ] `22-story-06-fragen.md:32` Pill → „204 ist Pragmatik, nicht Idempotenz“
- [ ] `26-tracing.md`-Note: parent-id richtig beschreiben
- [ ] `14a-circuit-breaker-code.md:7` Untertitel → „Die Regeln, nach denen Booking gleich spielt.“ und eine Zeile für den Probe-Slot
- [ ] `16-story-04-fragen.md:32` Pill → „einer probiert, alle warten“, Note `:49` umformulieren
- [ ] `09-story-01.md:58` Resilienz-Satz aktualisieren
- [ ] `27-story-08.md:30` und `story-08-tracing.md:47-50` Event-Kriterium als optional markieren, Hook `26-titel:21` umschreiben
- [ ] `25-story-07-fragen.md:58` CQRS-Brücke klären (Variante A oder B)
- [ ] `03-agenda.md:36` an die Entscheidung zu CQRS/BFF/Deployment anpassen
- [ ] Kicker 15, 18, 18a, 24, 24a → „Rollenspiel“, dazu Dashboard `1372`, `2990`, `2992`, `story-05-bulkhead.md:3`
- [ ] Zeit-Badges `15:9`, `10d:9` auf den Slot umstellen („≈ 60 min <small>davon 20 Spiel</small>“)
- [ ] `18-story-05.md:32` Becher benennen („Flight 4, Car 3, Hotel 3 Chips“) und die 503-Regel auf die Folie
- [ ] `15-story-04.md:29-31` Fallback-Satz, Regel „Strichliste wegwischen“, Rundenuhr ergänzen
- [ ] `24-story-07.md`: „Hotel lehnt in diesem Spiel jede Buchung ab.“, Titel → „Der Tisch ist der Broker“
- [ ] Notes anlegen: `12-story-03`, `13a`, `05a`, `14-titel`, `17-titel`, `10f` bis `10j`
- [ ] `10g-quiz-2.md:9` „Kunde 7 existiert“ ergänzen
- [ ] `13-story-03-fragen.md:30-32` `docker stop` vs. `kill`
- [ ] `29-zusammenfassung.md:12` Karte 01 → „Health, Config aus der Umgebung, Logs auf stdout.“
- [ ] `14-titel.md:11`, `17-titel.md:11` → „nach Michael Nygard, Release It! (2007)“
- [ ] `30a-conway.md` „1968“, `04-praxis` Kicker und Fußzeile präzisieren
- [ ] `28-story-08-fragen.md:65` TODO-Zeile entfernen, Frage 03 ohne `log.Printf`
- [ ] `20-titel`, `23-titel`, `26-titel` „Die Grafik“-Bullets kürzen oder streichen
- [ ] Dashboard: #CCF465 → #D3F871 (`1033`, `1046-1049`, `1074`)
- [ ] Dashboard: Story-7-Legende `1883` → „STORNO ANGEFORDERT“, Z. 2767 anpassen
- [ ] Dashboard: Story-8-Filter mit `svc` je nach Modus (`2812-2814`)
- [ ] Dashboard: Header-h1 → „Microservices richtig gemacht“, Umschalter „Referenz | Mein Service“
- [ ] Dashboard: `.stepper-connector.done` in `showStory()` setzen
- [ ] Dashboard: Story 3 Kriterien um „201 Created“ und Brücke zu Story 2 ergänzen
- [ ] Dashboard: `.story-helpers details[open]{flex-basis:100%;order:-1}` und #888 → #5B4E93
- [ ] `docs/questions/story5.md:54-56` „schützt Booking, nicht das Backend“
- [ ] `08b-zwei-fragen.md:38`, `themen.md:119` SAML ersetzen
- [ ] `story-06-saga.md:27-28`, `:34` Webhook-/Choreography-Hinweis ersetzen (Kritiker-Hinweis, vorher prüfen)