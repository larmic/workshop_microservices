# 3D-Druck: Requisiten für die Rollenspiele

Ein parametrisches OpenSCAD-Modell: das PROBE-Token aus Story 4, das seit dem
Kartenstapel optional ist. Story 5 und Story 7 brauchen keinen 3D-Druck,
gespielt wird mit Stühlen und einer Wand aus dem Raum bzw. mit Würfeln und
Pokerchips.
Im Repo liegen nur die `.scad`-Quellen; die STL-Dateien entstehen per Export
(siehe unten) und werden nicht eingecheckt.

| Datei | Teil | Maße (B × T × H) | Material ca. | Druckzeit ca. |
|---|---|---|---|---|
| `probe_token.scad` | PROBE-Token (nicht mehr nötig, optional) | 60 × 70 × 5 mm | 15 g | 25 min |

Material und Zeit sind grobe Schätzungen ohne Slicer (Volumen aus dem STL,
umgerechnet mit den Einstellungen unten); der Slicer liefert die genauen Werte.
Das Token passt bequem auf die Bambu Lab X1C (256 × 256 mm).

## Druckeinstellungen

- Material: PLA, Düse 0,4 mm
- Schichthöhe: 0,2 mm (Token mit Schrift gern 0,16 mm)
- Wände: 3 Perimeter, oben/unten je 5 Schichten
- Infill: 15 % Gyroid
- Stützstrukturen: keine nötig. Das Token liegt flach auf der
  Unterseite, die Schrift des Tokens ist erhaben bzw. auf der Unterseite
  vertieft
- Mehrfarbig ist optional: Beim Token kann man auf Höhe 4 mm (Beginn der
  erhabenen Schrift) einen Filamentwechsel einplanen. Das AMS wird nicht
  gebraucht.

## STL exportieren

Mit OpenSCAD (getestet mit 2021.01) direkt oder über Docker:

```sh
openscad -o probe_token.stl probe_token.scad
```

Alle Parameter stehen oben in der Datei. Die Schrift des Tokens nutzt „DejaVu Sans Bold"; ist
sie lokal nicht installiert, `font` z. B. auf `"Liberation Sans:style=Bold"`
setzen.

## Einsatz im Spiel

**PROBE-Token (Story 4, nicht mehr nötig).** Seit dem Umbau auf den
Kartenstapel steht die Probe auf den Karten 14 und 18: Der Request mit dieser
Karte legt sie aufs Probe-Feld des Circuit-Breaker-Boards und wartet dort. Ein Token braucht es dafür
nicht mehr. `probe_token.scad` bleibt optional im Repo, wer eins hat, kann es
dem Probe-Request zusätzlich in die Hand geben.
