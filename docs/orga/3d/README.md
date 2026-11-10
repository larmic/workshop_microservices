# 3D-Druck: Requisiten für die Rollenspiele

Drei parametrische OpenSCAD-Modelle für die Rollenspiele aus Story 4 und 5
(das PROBE-Token ist seit dem Kartenstapel in Story 4 optional).
Im Repo liegen nur die `.scad`-Quellen; die STL-Dateien entstehen per Export
(siehe unten) und werden nicht eingecheckt.

| Datei | Teil | Maße (B × T × H) | Material ca. | Druckzeit ca. |
|---|---|---|---|---|
| `bulkhead_tray.scad` | Bulkhead-Tablett | 161,5 × 109,5 × 12,4 mm | 70 g | 1,5–2 h |
| | 2 Trennwände (liegend) | je 103,2 × 34,7 × 3 mm, eingesteckt 25 mm über dem Rand | je 12 g | je 15 min |
| `probe_token.scad` | PROBE-Token (nicht mehr nötig, optional) | 60 × 70 × 5 mm | 15 g | 25 min |
| `card_stand.scad` | Aufsteller für A6-Karten | 70 × 40 × 18 mm | je 17 g | je 25 min |

Material und Zeit sind grobe Schätzungen ohne Slicer (Volumen aus dem STL,
umgerechnet mit den Einstellungen unten); der Slicer liefert die genauen Werte.
Alle Teile passen bequem auf die Bambu Lab X1C (256 × 256 mm). Das Tablett
mit beiden Wänden passt als `part="all"` auf eine Platte (161,5 × 193 mm),
15 Aufsteller als Raster 3 × 5 ebenfalls (222 × 224 mm, ca. 250 g, 5–6 h).

## Druckeinstellungen

- Material: PLA, Düse 0,4 mm
- Schichthöhe: 0,2 mm (Token und Wände mit Schrift gern 0,16 mm)
- Wände: 3 Perimeter, oben/unten je 5 Schichten
- Infill: 15 % Gyroid
- Stützstrukturen: keine nötig. Tablett und Aufsteller liegen flach auf der
  Unterseite, Wände und Token liegen flach, Schrift ist oben vertieft oder
  erhaben bzw. auf der Unterseite vertieft
- Mehrfarbig ist optional: Beim Token kann man auf Höhe 4 mm (Beginn der
  erhabenen Schrift) einen Filamentwechsel einplanen, beim Tablett genügt eine
  Farbe. Das AMS wird nicht gebraucht.
- Passung: Die Schlitze im Tablett sind 0,3 mm breiter und länger als die
  Wände (`wall_fit`). Sitzen die Wände zu stramm, `wall_fit = 0.4` setzen und
  nur das Tablett neu drucken.

## STL exportieren

Mit OpenSCAD (getestet mit 2021.01) direkt oder über Docker:

```sh
openscad -D 'part="all"'  -o bulkhead_plate.stl bulkhead_tray.scad   # Tablett + 2 Wände
openscad -D 'part="tray"' -o bulkhead_tray.stl  bulkhead_tray.scad
openscad -D 'part="wall"' -o bulkhead_wall.stl  bulkhead_tray.scad   # 2× drucken
openscad -o probe_token.stl probe_token.scad
openscad -D grid_x=3 -D grid_y=4 -o card_stands.stl card_stand.scad  # 12 Stück
```

`part="assembled"` zeigt das Tablett mit eingesteckten Wänden, nur zur
Ansicht. Alle Parameter stehen oben in den Dateien, z. B. `chip_d`, `chip_h`,
`stack`, `clearance`, `wall_t`, `labels`, `columns`, `rows` beim Tablett oder
`card_t`, `slot_extra` und `tilt` beim Aufsteller. Die Schrift nutzt
„DejaVu Sans Bold"; ist sie lokal nicht installiert, `font` z. B. auf
`"Liberation Sans:style=Bold"` setzen.

## Einsatz im Spiel

**Bulkhead-Tablett (Story 5).** Das Tablett ersetzt die Becher und steht bei
Booking. Drei Spalten FLIGHT, HOTEL und CAR mit je zwei Mulden für 40-mm-Pokerchips
(jede Mulde fasst bis zu drei Chips, die Griffmulde in der Mitte der Spalte
hilft beim Herausnehmen). Im Probelauf und in Runde 1 bleiben die Wände
draußen: Die sechs Chips liegen in einem gemeinsamen Pool, jeder Aufruf nimmt
irgendeinen freien Chip, und sobald Hotel alle sechs hält, wartet auch Flight.
Vor Runde 2 steckt Booking die beiden Wände sichtbar in die Schlitze, jetzt
hat jeder Service genau seine zwei Chips. Ist die Hotel-Spalte leer, gibt es
sofort `hotels: []`, Flight und Car greifen weiter in ihre eigene Spalte.
Das Einstecken der Wände ist der Moment, in dem die Gruppe das Pattern sieht:
Die Ressourcen werden nicht mehr, sie werden aufgeteilt.

**PROBE-Token (Story 4, nicht mehr nötig).** Seit dem Umbau auf den
Kartenstapel steht die Probe auf den Karten 14 und 18: Der Request mit dieser
Karte legt sie aufs Probe-Feld des Circuit-Breaker-Boards und wartet dort. Ein Token braucht es dafür
nicht mehr. `probe_token.scad` bleibt optional im Repo, wer eins hat, kann es
dem Probe-Request zusätzlich in die Hand geben.

**Kartenaufsteller (Story 5 und 7).** Die Aufsteller halten A6-Karten im
Querformat (148 × 105 mm) leicht nach hinten geneigt (10°), damit die ganze
Schlange sie lesen kann: Rollenschilder für Story 5 und 7 sowie Hinweiskarten
wie die Timeout-Regel in Story 5. Die Zustandskarten aus Story 4 brauchen
keinen Aufsteller mehr, sie liegen im Zustandsfeld des Circuit-Breaker-Boards.
Der Schlitz ist auf 250-g-Karton ausgelegt (`card_t = 0.3`,
Schlitz 0,8 mm). Für laminierte Karten `card_t` auf die gemessene Dicke setzen.
