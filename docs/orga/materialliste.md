# Materialliste: Spielmaterial für die Rollenspiele

Packliste für „Microservices richtig gemacht“ mit 8, 12 und 16 Teilnehmenden. Sie umfasst die Rollenspiele aus Story 4 (Circuit Breaker), Story 5 (Bulkhead) und Story 7 (Choreography-Saga) sowie die Design-Session aus Story 2.

- **Druckvorlagen:** <https://larmic.github.io/workshop_microservices/material/> (Quelle: [`pages/material/index.html`](../../pages/material/index.html)). Die Seite rechnet die Mengen aus der Teilnehmerzahl aus und zeigt pro Papiersorte einen eigenen Druckstapel.
- **3D-Modelle:** [`docs/orga/3d/`](3d/) mit `card_stand.scad` für die Rollenkarten aus Story 7 (`probe_token.scad` ist nicht mehr nötig). Druckeinstellungen und STL-Export stehen in der [README](3d/README.md).
- **Regeln:** [Story 4](../stories/story-04-circuit-breaker.md), [Story 5](../stories/story-05-bulkhead.md), [Story 7](../stories/story-07-choreography-saga.md), [Story 2](../stories/story-02-api-design-session.md).

**Art:** *kaufen* heißt Büro oder Laden, *drucken* heißt über die Druckseite, *3D* heißt aus `docs/orga/3d/`, *Raum* heißt, was im Seminarraum ohnehin steht.

**Ohne Laminiergerät:** Alle Karten kommen auf **Karton 250 g** (A4, an den Linien schneiden). Die Strichliste aus Story 7 wird auf normales Papier gedruckt und in eine **klare Prospekthülle** gesteckt, dann mit Whiteboard-Marker beschrieben und mit einem Tuch weggewischt. Story 4 braucht keine Strichliste, dort zählt das Circuit-Breaker-Board mit seinen fünf Plätzen „Fehler in Folge“. In Story 5 zählen die Services an den Fingern.

## Aufteilung nach Gruppengröße

| | 8 | 12 | 16 |
|---|---|---|---|
| Story 4 | 1 Booking, 7 Requests | 1 Booking, 11 Requests | 1 Booking, 15 Requests |
| Story 5 | 1 Gruppe: Booking, Hotel, 6 Requests (Flight und Car als Schild, der Trainer zählt) | 1 Gruppe: Booking, Hotel, Flight, Car, 8 Requests | 2 Gruppen zu 8 |
| Story 7 | 1 Tisch zu 8, 4 Reisen, Car ist der Korb | 2 Tische zu 6, je 2 Reisen | 2 Tische zu 8, je 4 Reisen |
| Story 2 | 2 Teams | 3–4 Teams | 4–5 Teams |

## Story 4: Wenn der Flug ausfällt

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Kartenstapel 01–27 mit Zustandskarten | drucken | 1 | 1 | 1 | 30 Karten: 01–27 und CLOSED / OPEN / HALF_OPEN, A7 hoch, 8 pro Bogen, 4 Bogen Karton (2 Felder frei), einseitig. Nur Nummer, Wartezeit und eine Zeile, STOPP 04, 13, 17, 20, 27 dunkel. Derselbe Stapel für jede Gruppengröße |
| Circuit-Breaker-Board | drucken | 1 | 1 | 1 | A3 quer oder 2× A4 hoch zusammengeklebt: Fehler 1–5, Probe-Feld, Zustandsfeld, ohne Regelblock |
| Trainer-Skript | drucken | 1 | 1 | 1 | 1 Bogen Papier 80 g, beidseitig: Regeln und STOPP-Fragen vorn, Kartentabelle 01–27 hinten |
| Tisch für die abgelegten Karten | Raum | 1 | 1 | 1 | alle Karten auf einen Haufen |

Board und Zustandskarte liegen flach auf dem Booking-Tisch, Booking legt das Board erst bei STOPP 1 hin.

## Story 5: Der geteilte Pool

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Kartenstapel 01–27 | drucken | 1 | 1 | 2 | A7 hoch, 8 pro Bogen, 4 Bogen Karton (5 Felder frei). Nur Nummer und „Steh auf bei XX“, STOPP 04, 15, 27 dunkel. Pro Gruppe ein Stapel |
| Service-Schilder Hotel, Flight, Car | drucken | 3 | 3 | 6 | A4 hoch, Karton, 3 Bogen pro Gruppe. Hotel doppelseitig („schnell“, hinten „langsam“): Rückseite bedrucken oder zwei Blätter zusammenkleben |
| Trainer-Skript | drucken | 1 | 1 | 1 | 1 Bogen Papier 80 g, beidseitig: Ablauf, Regeln und STOPP-Fragen vorn, Kartentabelle und Durchlauf mit 10 Personen hinten |
| Stühle | Raum | 2 | 2 | 4 | zwei pro Gruppe, Sitzfläche zu Booking |
| Wand | Raum | 1 | 1 | 2 | Pinnwand, Flipchart oder ein hochkant gestellter Tisch |
| Tisch für die abgelegten Karten | Raum | 1 | 1 | 2 | alle Karten auf einen Haufen |

Kein Einkauf, kein 3D-Druck, kein A3. Ab 14 Personen spielen zwei Gruppen parallel, jede mit eigenem Satz.


## Story 7: Die Saga wird leise

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Becher | kaufen | 4 | 4 | 8 | einer pro Reise |
| Becher-Etiketten „Reise 1…“ | drucken | 4 | 4 | 8 | mit Kreppband ankleben |
| Gummibärchen | kaufen | ca. 30 + Naschreserve | ca. 30 + Naschreserve | ca. 60 + Naschreserve | etwa 7 pro Reise und Tisch, Alternative bei Allergien oder veganer Ernährung: Büroklammern |
| Schälchen für den Vorrat (Hotel, Flight) | kaufen | 2 | 4 | 4 | |
| Korb oder Schachtel | kaufen | 1 | 2 | 2 | Broker in Runde 3 |
| Zettel E1… „Hotel · Reise 1“, „Flight · Reise 1“, … | drucken | 8 | 8 | 16 | zwei pro Reise, pro Tisch beschriftet |
| Signalkarte „Runde fertig, Korb leer“ | drucken | 1 | 2 | 2 | erst dann Kassensturz |
| Buchungslisten Hotel und Flight mit 3 Knick-Spalten | drucken | 2 | 4 | 4 | Papier; Klemmbrett optional |
| Strichliste „dispatch-failed“ mit Prospekthülle | drucken + kaufen | 1 | 2 | 2 | |
| Rollenkarten Booking, Hotel, Flight, Car, Korb, Kund:innen | drucken | 6 | 12 | 12 | Car spielt den Korb unter 9 Personen am Tisch und stellt die Korb-Karte in Runde 3 dazu |
| Briefing-Karte Flight (Story 7) | drucken | 1 | 2 | 2 | geheim, für beide Tische identisch |
| Kugelschreiber und Whiteboard-Marker | kaufen | 2 + 1 | 4 + 2 | 4 + 2 | Hotel und Flight streichen durch, Booking macht Striche |
| Spickzettel Trainer | drucken | 1 | 1 | 1 | A4, Papier |

Die Leiter (Folie 24a) braucht kein eigenes Material, die Becher bleiben einfach stehen.

## Story 2: Design-Session am Flipchart

Nichts zu drucken, die Teams schreiben die Tabelle direkt aufs Flipchart.

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Flipchart-Ständer oder Whiteboard-Fläche | Büro | 2 | 4 | 5 | eine pro Team |
| Flipchart-Bogen | Büro | 2 + 2 | 4 + 2 | 5 + 2 | mit Reserve |
| Flipchart-Marker in 2 Farben | kaufen | 6 | 12 | 15 | etwa 3 pro Team |
| Kreppband | kaufen | 1 Rolle | 1 Rolle | 1 Rolle | Bogen nebeneinander aufhängen, Quiz 3 und 4 greifen darauf zurück |
| Plenum-Flipchart | Büro | 1 | 1 | 1 | Erwartungen und Schnittfrage, bleibt beide Tage hängen |

## Für alle Spiele

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Kartenaufsteller A6 quer | 3D | 6 | 12 | 12 | `card_stand.scad` für die Rollenkarten aus Story 7, Schlitz 0,8 mm hält eine gefalzte 250-g-Karte; 12 Stück passen als Raster 3 × 4 auf eine Platte |
| Karton 250 g, A4 | kaufen | 16 Bogen | 20 Bogen | 28 Bogen | plus Reserve; Story 4: 4 / 4 / 4 Bogen (Stapel mit Zustandskarten), Story 5: 7 / 7 / 14 (Karten und Schilder), Story 7: 5 / 9 / 10 |
| Papier 80 g, A4 | Büro | 6 Bogen | 9 Bogen | 9 Bogen | Trainer-Skripte Story 4 und 5 (je ein Blatt beidseitig), Spickzettel, Listen und Strichliste Story 7 |
| A3-Papier | Büro | 1 | 1 | 1 | Circuit-Breaker-Board aus Story 4, sonst 2× A4 |
| Klare Prospekthüllen A4 | kaufen | 1 + 2 | 2 + 2 | 2 + 2 | für die Strichlisten aus Story 7, plus Reserve |
| Schere oder Schneidemaschine, Lineal, Falzbein (oder Messerrücken) | Büro | 1 | 1 | 1 | 250 g vor dem Falzen rillen |
| Umschläge oder Gummibänder | Büro | ca. 6 | ca. 8 | ca. 10 | Karten pro Spiel und Tisch getrennt aufbewahren |
| Beamer bzw. Bildschirm | Büro | 1 | 1 | 1 | Folien, Dashboard-Demos, QR-Code fürs Feedback |

Beim Rollenkarten-Modus „Einzelkarte“ werden die Rollen- und Briefing-Karten beidseitig gedruckt (an der kurzen Kante wenden) statt gefalzt. Das braucht weniger Karton, der Bogen muss aber passgenau laufen.

## Zusammengefasst

**Kaufen:** Becher und Schälchen, Korb oder Schachtel (Story 7), Gummibärchen oder Büroklammern, Prospekthüllen, Whiteboard-Marker mit Tuch, Kugelschreiber, Flipchart-Marker in 2 Farben, Kreppband, Karton 250 g.

**Drucken:** alle Bogen der Druckseite, je Papiersorte ein Stapel: Karton 250 g, Papier 80 g, A3 (oder das Board als 2× A4).

**3D:** Kartenaufsteller für Story 7 (6 bzw. 12). Vorlauf einplanen, ein Raster mit 12 Aufstellern braucht 5 bis 6 h.

**Aus dem Raum:** Stühle, Tische und eine Wand (Pinnwand oder Flipchart) für Story 5, ein Tisch für den Haufen in Story 4.

## Am Vorabend prüfen

- [ ] Teilnehmerzahl bestätigt. Hat sie sich geändert, die Druckseite neu rechnen lassen und fehlende Karten nachdrucken (Story 5 ab 14 ein zweiter Satz, Story 7 ab 12 und ab 9 am Tisch).
- [ ] Story 4: Kartenstapel 01–27 lückenlos und der Reihe nach sortiert (mit der Kartentabelle auf Seite 2 des Trainer-Skripts abgleichen), 01 liegt oben, die drei Zustandskarten dabei. Circuit-Breaker-Board liegt bereit (bei 2× A4 zusammengeklebt). Tisch für den Haufen freigeräumt.
- [ ] Strichliste für Story 7 steckt in der Prospekthülle. Den Marker einmal getestet: Er schreibt und lässt sich wieder wegwischen.
- [ ] Briefing-Karte für Flight (Story 7) liegt verdeckt und getrennt von den Rollenkarten.
- [ ] Story 5: Kartenstapel 01–27 lückenlos und der Reihe nach sortiert, 01 liegt oben. Hotel-Schild zeigt vorn „schnell“, hinten „langsam“. Zwei Stühle und eine Wand (Pinnwand, Flipchart oder Tisch) sind im Raum.
- [ ] Story 7: Becher beschriftet, Schälchen gefüllt, Allergien bzw. vegane Ernährung abgefragt (sonst Büroklammern). E-Zettel pro Tisch sortiert, nur E1 bis E(2 × Reisen). Buchungslisten auf die Spalte Runde 1 geknickt.
- [ ] Story 2: Flipchart-Bogen, Marker und Kreppband liegen bereit, Wandfläche für alle Bogen nebeneinander vorhanden.
- [ ] Raum: Platz für eine stehende Schlange vor Booking (Story 4 und 5) und zwei Spieltische (Story 7 ab 12). Story 5 ab 14: zwei Spielflächen mit je zwei Stühlen und Wand.
- [ ] Dashboard: alle Services auf „Normal“, damit die Demos nach den Spielen sauber starten.
