# Materialliste: Spielmaterial für die Spiele

Packliste für „Microservices richtig gemacht“ mit 8, 12 und 16 Teilnehmenden. Sie umfasst die Rollenspiele aus Story 4 (Circuit Breaker) und Story 5 (Bulkhead), das Würfelspiel aus Story 7 (Choreography-Saga) sowie die Design-Session aus Story 2.

- **Druckvorlagen:** <https://larmic.github.io/workshop_microservices/material/> (Quelle: [`pages/material/index.html`](../../pages/material/index.html)). Die Seite rechnet die Mengen aus der Teilnehmerzahl aus und zeigt pro Papiersorte einen eigenen Druckstapel.
- **3D-Modelle:** werden nicht mehr gebraucht. [`docs/orga/3d/`](3d/) enthält nur noch das optionale PROBE-Token aus Story 4.
- **Regeln:** [Story 4](../stories/story-04-circuit-breaker.md), [Story 5](../stories/story-05-bulkhead.md), [Story 7](../stories/story-07-choreography-saga.md), [Story 2](../stories/story-02-api-design-session.md).

**Art:** *kaufen* heißt im Laden besorgen, *Büro* heißt, was im Büro ohnehin vorhanden ist, *drucken* heißt über die Druckseite, *Raum* heißt, was im Seminarraum ohnehin steht.

**Ohne Laminiergerät:** Alle Karten kommen auf **Karton 200–250 g** (A4, an den Linien schneiden). Viele Bürodrucker schaffen nur bis 200 g (z. B. Xerox C315), das reicht. Beidseitig drucken sie Karton meist nicht automatisch: Das Hotel-Schild aus Story 5 dann zweimal einseitig drucken und Rücken an Rücken kleben. Story 4 braucht keine Strichliste, dort zählt das Circuit-Breaker-Board mit seinen fünf Plätzen „Fehler in Folge“. In Story 5 zählt niemand, die Chips in der Hand sind das Ergebnis.

## Aufteilung nach Gruppengröße

| | 8 | 12 | 16 |
|---|---|---|---|
| Story 4 | 1 Booking, 7 Requests | 1 Booking, 11 Requests | 1 Booking, 15 Requests |
| Story 5 | Booking, Flight, Hotel, Car, 4 Requests: 3 parallel, 1 pro Service | Booking, Flight, Hotel, Car, 8 Requests: 6 parallel, 2 pro Service | Booking, Flight, Hotel, Car, 12 Requests: 9 parallel, 3 pro Service |
| Story 7 | 2 Tische zu 4 | 2 Tische zu 6 | 3 Tische (5, 5, 6) |
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
| Pokerchips in 4 Farben | kaufen | 80 | 80 | 80 | je 20 pro Farbe: Flight, Hotel, Car und Default. Dieselben Chips wie in Story 7. Nach jeder Stufe gehen alle Chips zurück an die Services, sonst reichen 20 Flight-Chips nicht |
| Service-Schilder Hotel, Flight, Car | drucken | 3 | 3 | 3 | A4 hoch, Karton, 4 Bogen pro Gruppe (3 mit Karton-Duplex). Hotel doppelseitig („schnell“, hinten „langsam“): Rückseite bedrucken oder zwei Blätter zusammenkleben |
| Trainer-Skript | drucken | 1 | 1 | 1 | 1 Bogen Papier 80 g, beidseitig: Aufbau, Regeln und STOPP-Fragen vorn, Durchlauf mit 14 Personen und Ergebnis hinten |
| Kreppband (optional) | Büro | 1 Rolle | 1 Rolle | 1 Rolle | zwei Felder auf dem Boden: „bei Booking“ und „wartet bei Hotel“ |

Kein Kartenstapel, keine Stühle, kein 3D-Druck, kein A3. Es braucht einen Request mehr als parallele Plätze. Ab 18 Personen spielen zwei Gruppen parallel, jede mit eigenen Schildern, die Chips werden geteilt.


## Story 7: Die Saga wird leise

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Würfel | kaufen | 2 | 2 | 3 | einer pro Tisch, Ersatz: Würfel-App |
| Pokerchips | kaufen | 26 | 26 | 39 | 13 pro Tisch, Farbe egal; vorher abzählen, die Balance hängt daran |
| Zettel und Stift für Booking | Büro | 2 | 2 | 3 | Striche für „Storniert“ und gelungene Reisen |
| Flipchart oder Whiteboard | Raum | 1 | 1 | 1 | Ergebnisse aller Tische: Budget, gewonnen, stille Stornos |
| Spickzettel Trainer | drucken | 1 | 1 | 1 | Papier 80 g, einseitig |

Für die Tische ist nichts zu drucken. Die Chips bleiben nach dem Kassensturz bis zum Recap liegen.

## Story 2: Design-Session am Flipchart

Nichts zu drucken, die Teams schreiben die Tabelle direkt aufs Flipchart.

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Flipchart-Ständer oder Whiteboard-Fläche | Büro | 2 | 4 | 5 | eine pro Team |
| Flipchart-Bogen | Büro | 2 + 2 | 4 + 2 | 5 + 2 | mit Reserve |
| Flipchart-Marker in 2 Farben | kaufen | 6 | 12 | 15 | etwa 3 pro Team |
| Kreppband (optional) | Büro | 1 Rolle | 1 Rolle | 1 Rolle | nur wenn es weniger Flipchart-Ständer als Teams gibt und Bogen an die Wand sollen. Quiz 3 und 4 greifen auf die Bogen zurück |
| Plenum-Flipchart | Büro | 1 | 1 | 1 | Erwartungen und Schnittfrage, bleibt beide Tage hängen |

## Für alle Spiele

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Karton 200–250 g, A4, laser-geeignet | kaufen | 8 Bogen | 8 Bogen | 8 Bogen | plus Reserve; Story 4: 4 / 4 / 4 Bogen (Stapel mit Zustandskarten), Story 5: 4 / 4 / 4 (nur die Schilder, das Hotel-Schild als zwei einseitige Bogen). Mit Karton-Duplex einer weniger: 7 / 7 / 7 |
| Papier 80 g, A4 | Büro | 3 Bogen | 3 Bogen | 3 Bogen | Trainer-Skripte Story 4 und 5 (je ein Blatt beidseitig), Spickzettel Story 7. Ohne A3-Drucker 5 Bogen: das Board als 2× A4 |
| A3-Papier (optional) | Büro | 1 | 1 | 1 | nur mit A3-Drucker, sonst das Circuit-Breaker-Board aus Story 4 als 2× A4 |
| Schere oder Schneidemaschine, Lineal | Büro | 1 | 1 | 1 | Karten nur schneiden, nichts falzen |
| Umschläge oder Gummibänder | Büro | 1 | 1 | 1 | für den Kartenstapel aus Story 4 |
| Beamer bzw. Bildschirm | Büro | 1 | 1 | 1 | Folien, Dashboard-Demos, QR-Code fürs Feedback |

## Zusammengefasst

**Kaufen:** Würfel (Story 7) und Pokerchips in 4 Farben zu je 20 (Story 5 und 7, oft schon vorhanden), Flipchart-Marker in 2 Farben, Karton 200–250 g (eine Packung mit 50 Bogen reicht).

**Drucken:** alle Bogen der Druckseite, je Papiersorte ein Stapel: Karton 200–250 g, Papier 80 g, A3 nur mit A3-Drucker (sonst das Board als 2× A4).

**Aus dem Raum:** Tische für Story 7, ein Tisch für den Haufen in Story 4.

## Am Vorabend prüfen

- [ ] Teilnehmerzahl bestätigt. Hat sie sich geändert, die Druckseite neu rechnen lassen und fehlende Karten nachdrucken (Story 5 ab 18 ein zweiter Satz Schilder).
- [ ] Story 4: Kartenstapel 01–27 lückenlos und der Reihe nach sortiert (mit der Kartentabelle auf Seite 2 des Trainer-Skripts abgleichen), 01 liegt oben, die drei Zustandskarten dabei. Circuit-Breaker-Board liegt bereit (bei 2× A4 zusammengeklebt). Tisch für den Haufen freigeräumt.
- [ ] Story 5: je 20 Chips in vier Farben abgezählt (Flight, Hotel, Car, Default), Farben den Services zugeordnet. Hotel-Schild zeigt vorn „schnell“, hinten „langsam“. Aus der Teilnehmerzahl die Zahl der parallelen Plätze festgelegt (Trainer-Skript, Seite 1).
- [ ] Story 7: pro Tisch ein Würfel und genau 13 Chips abgezählt, Zettel und Stift für Booking.
- [ ] Story 2: Flipchart-Bogen und Marker liegen bereit, pro Team ein Ständer (sonst Kreppband und Wandfläche für alle Bogen nebeneinander).
- [ ] Raum: Platz für eine stehende Schlange vor Booking (Story 4 und 5) und ein Spieltisch pro Story-7-Tisch. Story 5: Flight, Hotel, Car nebeneinander, Platz für Wartende vor Hotel; ab 18 zwei Spielflächen.
- [ ] Dashboard: alle Services auf „Normal“, damit die Demos nach den Spielen sauber starten.
