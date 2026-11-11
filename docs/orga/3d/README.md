# 3D-Druck: Requisiten für die Rollenspiele

Zwei parametrische OpenSCAD-Modelle: der Kartenaufsteller für die
Rollenkarten aus Story 7 und das PROBE-Token aus Story 4, das seit dem
Kartenstapel optional ist. Story 5 braucht keinen 3D-Druck, gespielt wird mit
zwei Stühlen und einer Wand aus dem Raum.
Im Repo liegen nur die `.scad`-Quellen; die STL-Dateien entstehen per Export
(siehe unten) und werden nicht eingecheckt.

| Datei | Teil | Maße (B × T × H) | Material ca. | Druckzeit ca. |
|---|---|---|---|---|
| `card_stand.scad` | Aufsteller für A6-Karten | 70 × 40 × 18 mm | je 17 g | je 25 min |
| `probe_token.scad` | PROBE-Token (nicht mehr nötig, optional) | 60 × 70 × 5 mm | 15 g | 25 min |

Material und Zeit sind grobe Schätzungen ohne Slicer (Volumen aus dem STL,
umgerechnet mit den Einstellungen unten); der Slicer liefert die genauen Werte.
Alle Teile passen bequem auf die Bambu Lab X1C (256 × 256 mm),
15 Aufsteller als Raster 3 × 5 auf eine Platte (222 × 224 mm, ca. 250 g, 5–6 h).

## Druckeinstellungen

- Material: PLA, Düse 0,4 mm
- Schichthöhe: 0,2 mm (Token mit Schrift gern 0,16 mm)
- Wände: 3 Perimeter, oben/unten je 5 Schichten
- Infill: 15 % Gyroid
- Stützstrukturen: keine nötig. Aufsteller und Token liegen flach auf der
  Unterseite, die Schrift des Tokens ist erhaben bzw. auf der Unterseite
  vertieft
- Mehrfarbig ist optional: Beim Token kann man auf Höhe 4 mm (Beginn der
  erhabenen Schrift) einen Filamentwechsel einplanen. Das AMS wird nicht
  gebraucht.

## STL exportieren

Mit OpenSCAD (getestet mit 2021.01) direkt oder über Docker:

```sh
openscad -D grid_x=3 -D grid_y=4 -o card_stands.stl card_stand.scad  # 12 Stück
openscad -o probe_token.stl probe_token.scad
```

Alle Parameter stehen oben in den Dateien, z. B. `card_t`, `slot_extra` und
`tilt` beim Aufsteller. Die Schrift des Tokens nutzt „DejaVu Sans Bold"; ist
sie lokal nicht installiert, `font` z. B. auf `"Liberation Sans:style=Bold"`
setzen.

## Einsatz im Spiel

**Kartenaufsteller (Story 7).** Die Aufsteller halten A6-Karten im
Querformat (148 × 105 mm) leicht nach hinten geneigt (10°), damit der ganze
Tisch die Rollenschilder lesen kann. Die Zustandskarten aus Story 4 brauchen
keinen Aufsteller, sie liegen im Zustandsfeld des Circuit-Breaker-Boards.
Der Schlitz ist auf 250-g-Karton ausgelegt (`card_t = 0.3`,
Schlitz 0,8 mm). Für laminierte Karten `card_t` auf die gemessene Dicke setzen.

**PROBE-Token (Story 4, nicht mehr nötig).** Seit dem Umbau auf den
Kartenstapel steht die Probe auf den Karten 14 und 18: Der Request mit dieser
Karte legt sie aufs Probe-Feld des Circuit-Breaker-Boards und wartet dort. Ein Token braucht es dafür
nicht mehr. `probe_token.scad` bleibt optional im Repo, wer eins hat, kann es
dem Probe-Request zusätzlich in die Hand geben.
