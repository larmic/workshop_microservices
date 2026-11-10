# Materialliste: Spielmaterial für die Rollenspiele

Packliste für „Microservices richtig gemacht“ mit 8, 12 und 16 Teilnehmenden. Sie umfasst die Rollenspiele aus Story 4 (Circuit Breaker), Story 5 (Bulkhead) und Story 7 (Choreography-Saga) sowie die Design-Session aus Story 2.

- **Druckvorlagen:** <https://larmic.github.io/workshop_microservices/material/> (Quelle: [`pages/material/index.html`](../../pages/material/index.html)). Die Seite rechnet die Mengen aus der Teilnehmerzahl aus und zeigt pro Papiersorte einen eigenen Druckstapel.
- **3D-Modelle:** [`docs/orga/3d/`](3d/) mit `bulkhead_tray.scad` und `card_stand.scad` (`probe_token.scad` ist nicht mehr nötig). Druckeinstellungen und STL-Export stehen in der [README](3d/README.md).
- **Regeln:** [Story 4](../stories/story-04-circuit-breaker.md), [Story 5](../stories/story-05-bulkhead.md), [Story 7](../stories/story-07-choreography-saga.md), [Story 2](../stories/story-02-api-design-session.md).

**Art:** *kaufen* heißt Büro oder Laden, *drucken* heißt über die Druckseite, *3D* heißt aus `docs/orga/3d/`.

**Ohne Laminiergerät:** Alle Karten kommen auf **Karton 250 g** (A4, an den Linien schneiden). Die Strichliste aus Story 7 wird auf normales Papier gedruckt und in eine **klare Prospekthülle** gesteckt, dann mit Whiteboard-Marker beschrieben und mit einem Tuch weggewischt. Story 4 und 5 brauchen keine Strichliste mehr, dort zählen die Fehlerkarten 1 bis 5.

## Aufteilung nach Gruppengröße

| | 8 | 12 | 16 |
|---|---|---|---|
| Story 4 | 1 Booking, 7 Suchanfragen | 1 Booking, 11 Suchanfragen | 1 Booking, 15 Suchanfragen |
| Story 5 | 1 Tisch zu 8 (Flight + Car in einer Person) | 1 Tisch zu 12 (Flight, Car getrennt, zwei bei Booking) | 2 Tische zu 8 |
| Story 7 | 1 Tisch zu 8, 4 Reisen, Car ist der Korb | 2 Tische zu 6, je 2 Reisen | 2 Tische zu 8, je 4 Reisen |
| Story 2 | 2 Teams | 3–4 Teams | 4–5 Teams |

## Story 4: Wenn der Flug ausfällt

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Kartenstapel 00–27 mit Legende und Zustandskarten | drucken | 1 | 1 | 1 | 28 Karten plus Legende und CLOSED / OPEN / HALF_OPEN, A7 hoch, 8 pro Bogen, 4 Bogen Karton, einseitig. Derselbe Stapel für jede Gruppengröße |
| Circuit-Breaker-Board | drucken | 1 | 1 | 2 | A3 quer oder 2× A4 hoch zusammengeklebt: Fehler 1–5, Probe-Feld, Zustandsfeld. Eins pro Story-5-Tisch; jeder weitere Tisch braucht dazu 1 Bogen Karton mit Zustandskarten |
| Trainer-Skript | drucken | 1 | 1 | 1 | 1 Bogen Papier 80 g, beidseitig: STOPPs vorn, Kartentabelle 00–27 hinten |
| Ablagetisch | Büro | 1 | 1 | 1 | die Karten landen dort der Reihe nach, die Nummer auf der Karte reicht |

Board und Zustandskarte liegen flach auf dem Booking-Tisch, Booking legt das Board erst bei STOPP 1 hin.

## Story 5: Der geteilte Pool

| Was | Art | 8 | 12 | 16 | Hinweis |
|---|---|---|---|---|---|
| Bulkhead-Tablett mit 2 Trennwänden | 3D | 1 | 1 | 2 | ohne Wände ein Pool (Probelauf, Runde 1), mit Wänden 2 Chips je Service (Runde 2). Ersatz: 1 großer und 3 kleine Becher je Tisch |
| Pokerchips 40 mm, eine Farbe | kaufen | 6 + 4 | 6 + 4 | 12 + 4 | 6 pro Tisch plus Reserve |
| Tablett-Streifen „Flight \| Hotel \| Car“ und „ein Pool“ | drucken | 1 Bogen | 1 Bogen | 1 Bogen | Spaltenabstand 52 mm wie im Tablett. Auf dem Bogen sind auch die Ablage-Schilder und Becher-Etiketten |
| Suchkarten #01–#20 | drucken | 20 | 20 | 40 | pro Tisch ein Satz |
| Orange Karten „hotels: []“ | drucken | 20 | 20 | 40 | gebraucht werden 14 pro Tisch, der Rest ist Reserve |
| Uhrmatte mit 8 Feldern | drucken | 1 | 1 | 2 | A3 quer oder 2× A4 hoch zusammengeklebt |
| Rollenkarten | drucken | 5 | 7 | 10 | Takt, Booking, Hotel, Flight + Car bzw. Flight und Car, Suchanfragen; ab 12 am Tisch dazu „Booking an den Stapeln“ |
| Circuit-Breaker-Board und Zustandskarten | aus Story 4 | 1 | 1 | 2 | eins pro Tisch, nach Story 4 aufheben |
| Spickzettel Trainer mit Takt-für-Takt-Verlauf | drucken | 1 | 1 | 1 | A4, Papier |

Folie 18a (Little's Law) braucht kein Material.

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
| Kartenaufsteller A6 quer | 3D | 6 | 12 | 12 | `card_stand.scad` für die Rollenkarten aus Story 5 und 7, Schlitz 0,8 mm hält eine gefalzte 250-g-Karte; 12 Stück passen als Raster 3 × 4 auf eine Platte |
| Karton 250 g, A4 | kaufen | 18 Bogen | 23 Bogen | 32 Bogen | plus Reserve; Story 4 allein: 4 / 4 / 5 Bogen (Stapel mit Legende und Zustandskarten, ab 14 ein Bogen Zustandskarten mehr) |
| Papier 80 g, A4 | Büro | 6 Bogen | 9 Bogen | 9 Bogen | Strichliste Story 7, Listen, Spickzettel; dazu das Trainer-Skript Story 4 als ein Blatt beidseitig |
| A3-Papier | Büro | 2 | 2 | 4 | Circuit-Breaker-Board und Uhrmatte je Tisch, sonst je 2× A4 |
| Klare Prospekthüllen A4 | kaufen | 1 + 2 | 2 + 2 | 2 + 2 | für die Strichlisten aus Story 7, plus Reserve |
| Schere oder Schneidemaschine, Lineal, Falzbein (oder Messerrücken) | Büro | 1 | 1 | 1 | 250 g vor dem Falzen rillen |
| Umschläge oder Gummibänder | Büro | ca. 6 | ca. 8 | ca. 10 | Karten pro Spiel und Tisch getrennt aufbewahren |
| Beamer bzw. Bildschirm | Büro | 1 | 1 | 1 | Folien, Dashboard-Demos, QR-Code fürs Feedback |

Beim Rollenkarten-Modus „Einzelkarte“ werden die Rollen- und Briefing-Karten beidseitig gedruckt (an der kurzen Kante wenden) statt gefalzt. Das braucht weniger Karton, der Bogen muss aber passgenau laufen.

## Zusammengefasst

**Kaufen:** Pokerchips 40 mm, Becher und Schälchen, Korb oder Schachtel (Story 7), Gummibärchen oder Büroklammern, Prospekthüllen, Whiteboard-Marker mit Tuch, Kugelschreiber, Flipchart-Marker in 2 Farben, Kreppband, Karton 250 g.

**Drucken:** alle Bogen der Druckseite, je Papiersorte ein Stapel: Karton 250 g, Papier 80 g, A3 (oder Board und Uhrmatte als 2× A4).

**3D:** Bulkhead-Tablett mit 2 Trennwänden (1, ab 14 Personen 2), Kartenaufsteller (6 bzw. 12). Vorlauf einplanen, das Tablett braucht etwa 2 h, ein Raster mit 12 Aufstellern 5 bis 6 h.

## Am Vorabend prüfen

- [ ] Teilnehmerzahl bestätigt. Hat sie sich geändert, die Druckseite neu rechnen lassen und fehlende Karten nachdrucken (Story 5 ab 14, Story 7 ab 12 und ab 9 am Tisch).
- [ ] Story 4: Kartenstapel 00–27 lückenlos und der Reihe nach sortiert (mit der Kartentabelle auf Seite 2 des Trainer-Skripts abgleichen), 00 liegt oben, Legende und Zustandskarten dabei. Circuit-Breaker-Board liegt bereit (bei 2× A4 zusammengeklebt). Ablagetisch freigeräumt.
- [ ] Strichliste für Story 7 steckt in der Prospekthülle. Den Marker einmal getestet: Er schreibt und lässt sich wieder wegwischen.
- [ ] Briefing-Karte für Flight (Story 7) liegt verdeckt und getrennt von den Rollenkarten.
- [ ] Story 5: Die Wände sitzen im Tablett, sechs Chips plus Reserve passen in die Mulden, der Streifen liegt dabei. Suchkarten #01–#20 sind sortiert, die orangen Karten liegen getrennt. Uhrmatte liegt bereit (bei 2× A4 zusammengeklebt).
- [ ] Story 5 läuft am zweiten Tag: Circuit-Breaker-Board und Zustandskarten aus Story 4 über Nacht aufheben.
- [ ] Story 7: Becher beschriftet, Schälchen gefüllt, Allergien bzw. vegane Ernährung abgefragt (sonst Büroklammern). E-Zettel pro Tisch sortiert, nur E1 bis E(2 × Reisen). Buchungslisten auf die Spalte Runde 1 geknickt.
- [ ] Story 2: Flipchart-Bogen, Marker und Kreppband liegen bereit, Wandfläche für alle Bogen nebeneinander vorhanden.
- [ ] Raum: Platz für eine stehende Schlange vor dem Booking-Tisch (Story 4 und 5) und zwei Spieltische (Story 7 ab 12, Story 5 ab 14).
- [ ] Dashboard: alle Services auf „Normal“, damit die Demos nach den Spielen sauber starten.
