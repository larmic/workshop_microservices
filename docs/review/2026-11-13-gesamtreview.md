# Gesamtreview vom 13.11.2026: „Microservices richtig gemacht“

Geprüft wurden Folien, Dashboard, Story- und Fragen-Docs, Materialliste und Druckseite, Vorbereitung, Feedback, Landing-Page, Einladungstexte und das Easter Egg. Es wurde nichts geändert, dieses Dokument hält nur Befunde fest.

> **Pfad-Konvention:** Folien ohne Pfad liegen in `services/slides/chapters/`. „Dashboard“ meint `services/dashboard/static/index.html`, „Egg“ meint `services/dashboard/static/voxel.css`.

**Wie geprüft**

- Vier parallele Lese-Reviews (Folien, Dashboard gegen Folien/Docs/Code, Drumherum, Design und Egg).
- Alle 62 Folien im Browser gerendert, mit allen aufgedeckten Klick-Elementen: keine läuft über. Am knappsten sind 20a „Der Orchestrator“, 24 „Die Saga wird leise“ und 26b „Vier Stellen im Code“.
- Dashboard (Stories 0, 4 bis 8) mit und ohne Egg gerendert, alle vier Pages bei 1280, 390 und 320 px.
- Die Befunde unter „Unbedingt“ sowie W1, W3, W5, W9 und W29 sind an den Dateien nachgeprüft. Der Rest stammt aus den Lese-Reviews mit Fundstelle.
- **Nicht geprüft:** Inhalt des QR-Codes (`services/slides/assets/qr-code.png`), echtes Backend, Safari und Firefox, Schriften unter macOS.

**Aufwand:** klein = unter 30 Minuten, mittel = ein bis drei Stunden, hoch = mehr als ein halber Tag.

---

## 1. Kurzurteil

| Frage | Urteil |
|---|---|
| Passen die Inhalte zueinander? | **Weitgehend ja.** Zahlen und Regeln stimmen zwischen Folien, Dashboard, Docs und Go-Code überein (Schwellen, Timeouts, Chaos-Modi, Story-7-Regeln und Quoten). Abweichungen sitzen in der Story-5-Demo, in `docs/themen.md` und in der Materialliste. |
| Passen die Designs? | **Teilweise.** Palette und Footer sind über alle Oberflächen erkennbar gleich. Dashboard und Pages nutzen aber Systemschrift statt Poppins, eigene Ampelfarben und Emojis und wirken dadurch wie ein anderes Produkt als das Deck. Die Vorbereitungsseite ist auf dem Handy kaputt. |
| Passt das Easter Egg noch? | **Ja, bis auf vier Stellen.** `voxel.css` deckt das neue Spiel-Markup fast vollständig ab. Unlesbar ist nur die Überschrift des Demo-Panels in Story 4, 5 und 7. Texte und Aktivierung passen, Reste des alten Story-7-Spiels gibt es nicht. |
| Ist der rote Faden erkennbar? | **Ja.** Die Kette Recap 1 → Discovery → Circuit Breaker → Bulkhead → Cliffhanger 16a → Saga → Zuruf → Tracing ist in Hooks und Überleitungen sauber verbunden. Brüchig ist er am ersten Vormittag (kaum Notes, keine Überleitungen), bei der Kostenfrage aus 08b (nie eingelöst) und zwischen Agenda (vier Halbtage) und Kapiteln (drei). |

**Was gut ist und so bleiben sollte**

- Alle „Recap NN“-, „Frage NN“- und Bonus-Verweise in den Notes zeigen auf das richtige Ziel, die Zähl-Überschriften stimmen.
- Der Umbau von Story 7 ist überall angekommen: keine Becher, Gummibärchen, Körbe oder Leiter mehr in Folien, Dashboard, Druckseite und Materialliste (eine Ausnahme: U6).
- Kicker und Badges der Story-Folien sind einheitlich, Titel und Zeiten stimmen in Folie, Dashboard und `themen.md` für alle acht Stories überein.
- Befehle, Ports und Namen der Vorbereitung passen zu Makefile und Compose. Vorbereitungsseite und `docs/vorbereitung.md` sind eine Quelle.
- Einladungstexte versprechen nichts mehr, was der Workshop nicht liefert.

---

## 2. Unbedingt

Fehler, die im Workshop auffallen oder etwas Falsches lehren.

> **Stand 13.11.:** U1 bis U7 sind erledigt. U3 als Textänderung: Dashboard und Story-Doc sagen jetzt, dass das Story-5-Panel den Breaker nicht zeigt. Ein eigenes Breaker-Raster bleibt offen.

| Nr. | Bereich | Befund | Wo | Vorschlag | Aufwand |
|---|---|---|---|---|---|
| U1 | Vorbereitung | **Seite auf dem Handy kaputt.** Die Inhaltsspalte hat bei 390 und 320 px Breite 0 px, der Text läuft buchstabenweise am Rand. Die Datei hat keine einzige `@media`-Regel. Teilnehmer öffnen die Seite vorab oft am Handy. | `pages/vorbereitung/index.html:136` | Breakpoint aus `pages/material/index.html:291-294` übernehmen (unter 900 px einspaltig, Seitenleiste statisch). | klein |
| U2 | Folien | **Story-5-Note nennt 6 s „genau der Timeout“.** Der Timeout ist überall 3 s, eine Zeile darüber steht es richtig. | `18-story-05.md:55` | „Bei 26, also 6 s nach dem Ankommen. Das Doppelte des Timeouts, der Client wäre längst weg.“ | klein |
| U3 | Dashboard | **Story-5-Demo behauptet etwas, das das Panel nicht zeigt.** „Der Hotel-Breaker bleibt CLOSED“ ist die Kernaussage (langsam ist nicht kaputt), das Story-5-Panel hat aber kein Breaker-Raster. Das einzige steht in Story 4 und fragt einen anderen Service ab. | Dashboard `:2057`, `docs/stories/story-05-bulkhead.md:191` | Entweder ein Breaker-Raster für `booking-ref-story5` ins Panel (Endpoint existiert) oder den Halbsatz auf „im Log“ ändern. | klein (Text) / mittel (Raster) |
| U4 | Dashboard | **„hotel rej: etwa zehn“ stimmt nur beim ersten Burst.** Der Wert summiert über die letzten zehn Aufrufe im Browser, der Reset-Knopf leert nur die Server-Zähler. Beim zweiten Burst stehen dort rund 20. | Dashboard `:2057`, `:3135-3221`, `18-story-05.md:66` | Reset leert auch `bulkheadBursts`, oder der Wert zeigt nur den letzten Burst. | klein |
| U5 | Docs | **`themen.md` zählt den Story-5-Recap falsch:** „fünf Fragen“, Folie und Story-Doc haben vier. | `docs/themen.md:154` | „vier Fragen“. | klein |
| U6 | Docs | **Altlast „Schalter und Becher aus dem Rollenspiel“.** Becher gibt es nirgends mehr. | `docs/stories/story-08-tracing.md:155-157` | „das Board aus Story 4 und die Stühle aus Story 5“. | klein |
| U7 | Material | **Karton-Menge passt nicht zur eigenen Empfehlung.** Die Liste nennt 11 / 11 / 18 Bogen, empfiehlt aber, das Hotel-Schild zweimal einseitig zu drucken. Die Druckseite erzeugt dann 12 / 12 / 20. Die Abweichung entstand heute beim Umstellen auf 200 g. | `docs/orga/materialliste.md:11`, `:74` | „12 / 12 / 20 (mit Karton-Duplex 11 / 11 / 18)“, Story 5 auf „8 / 8 / 16“. | klein |

---

## 3. Wünschenswert

Spürbare Verbesserungen für Verständnis, Konsistenz und Auftritt.

> **Stand 13.11.:** Erledigt sind W1 bis W7, W10 bis W13, W15 bis W17, W19 bis W23, W25, W26, W29 bis W34 und W36 bis W38.
>
> - **W3** als Note-Zeile „Preis: …“ in den Recaps zu Story 3 bis 8, Folie 29 ist unverändert.
> - **W13:** Spielname im Kicker und Kassensturz-Zeile, kein Tischblatt.
> - **W31** ist nicht am Bild geprüft, weil die Chaos-Knöpfe ohne Backend nicht erscheinen.
>
> **Teilweise erledigt:**
> - **W14:** Die Reihenfolge der Blöcke 4 und 5 und die zusätzlichen Recap-Fragen in `themen.md` sind unverändert.
> - **W18:** Nur die Legende („im Raum“) ist angepasst, Oberbegriff und Strichart sind offen.
> - **W35:** Kreppband ist optional, das Plenum-Flipchart steht noch in der Liste.
>
> **Bewusst nicht umgesetzt:** W8 (Überleitungen am ersten Vormittag) und W9 (Zeitplan Tag 1 und Story 3).
>
> **Nachgezogen:** W27 (`services/CLAUDE.md` und Kopfkommentar von `theme.css`).
>
> **Nachgezogen:** W24 und W28.
> - **W24:** Health- und Startup-Symbole auf hellem Grund nutzen die Status-Farben. Die drei kräftigen Punkte im dunklen Kopf bleiben und sind dokumentiert. Das Gate `TestPaletteNoForeignColours` meldet jede Farbe außerhalb der Palette in Dashboard und `pages/`.
> - **W28:** Die Druckseite nutzt Poppins aus `pages/material/fonts/`. Die Trainer-Skripte zu Story 4 und 5 sind etwas enger gesetzt, alle Bogen laufen für 7 bis 20 Teilnehmer ohne Überlauf.

### Folien und roter Faden

| Nr. | Befund | Wo | Vorschlag | Aufwand |
|---|---|---|---|---|
| W1 | **Titelfolie Choreography passt nicht zum Würfelspiel.** Sie sagt, Booking warte „nur noch darauf, dass jemand ‚gehört‘ sagt“. Im Spiel sagt Booking sofort „Storniert“, erst danach würfeln die Gerufenen. | `23-titel.md:10-11`, `23-choreography-saga.md:28` | „Booking ruft zu und sagt ‚storniert‘, ob es jemand gehört hat oder nicht.“ Karte 3: „ruft zu, fragt nicht nach“. | klein |
| W2 | **Einstiegsfrage im Story-8-Recap läuft ins Leere.** „Wo hätte euch Tracing in Stories 4 bis 7 geholfen?“ Drei davon sind Spiele ohne Logs. | `28-story-08-fragen.md:57`, `docs/themen.md:224` | „… in Story 3 und 6? Und in den Demos zu 4, 5 und 7?“ | klein |
| W3 | **Kostenfrage aus 08b wird nie eingelöst.** Die Folie verspricht: „Jedes Pattern bringt neue Fehlerfälle mit. Die sammeln wir nach jeder Story im Recap.“ Kein Recap sammelt etwas. | `08b-zwei-fragen.md:21` | Pro Recap-Note eine Zeile „Preis: …“ und auf Folie 29 eine gedämpfte Preiszeile je Karte. Oder den Satz auf 08b abschwächen. | mittel |
| W4 | **Fallback „200 mit leerer Liste“ gegen Regel 6.** Story 2 lehrt „ein 200 mit Fehler im Body ist unsichtbar“, zwei Stunden später liefert Story 4 genau das. Keine Note fängt die Frage ab. | `10b:64`, `10d:50`, `14a:38`, `15:49` | Ein Satz in der Note zu Recap 4, Frage 02: Teilantwort ist kein Fehler, deshalb der Header `X-Fallback`. | klein |
| W5 | **Brücke Story 2 → 6 zeigt auf die falsche Variante.** „Der Storno vom Flipchart ist jetzt die Kompensation“, die Musterlösung war aber `POST …/cancellation` (nicht idempotent), Story 6 nutzt `DELETE` (idempotent). | `21-story-06.md:53`, `10e:14-16` | Note: „In Story 2 war das Storno eine Ressource mit Grund. Die Kompensation ist die technische Variante aus Quiz 4: DELETE, idempotent.“ | klein |
| W6 | **„Forward“ hat zwei Bedeutungen.** Mal Buchungsschritt (20a, 22, 27), mal Forward Recovery (20, 24). Die Karte auf Folie 20 heißt seit heute „Buchen“, die übrigen Stellen nicht. | `20a:16`, `22:30`, `:50`, `27:29` | Dort „buchen“ bzw. „Buchung“ schreiben, „Forward/Backward“ dem Begriffspaar vorbehalten. | klein |
| W7 | **Agenda und Kapitel sind zwei Gliederungen.** „Kapitel 1 von 3“ kommt unangekündigt, die Agenda zeigt vier Halbtage. Der Untertitel von 08a deckt nur Story 1. | `03-agenda.md`, `08a:5-9` | Kapitelmarken in der Agenda (K1 bis K3) oder ein Note-Satz, 08a-Untertitel erweitern. | klein bis mittel |
| W8 | **Erster Vormittag ohne Überleitungen.** Rund 135 Minuten Programm (05 bis 08b) mit sehr kurzen oder fehlenden Notes. Teils bewusst gekürzt, aber die Überleitungen fehlen als einzige im Deck. | `05`, `06`, `08`, `08b` | Je Folie eine Zeile „Weiter: …“, auf 08b ein Satz zum Eisberg. | mittel |
| W9 | **Zeitbudget Tag 1 geht nicht auf.** `themen.md` rechnet 255 Minuten Vormittag, bei 9:00 Beginn und Mittag um 12:30 sind es 210. Der Tag liegt bei 490 Minuten ohne Kaffeepause. Story 3 verplant 30 + 30 Minuten Bauen bei 60 Minuten Slot. | `docs/themen.md:28`, `03-agenda.md:16`, `12-story-03.md:59-60` | Überhang benennen („Story 1 endet nach dem Mittag“) oder Story 4 auf Tag 2 ziehen. Story 3 auf 25 + 25. | mittel |
| W10 | **Eisberg (08b) und „Was wir nicht angeschaut haben“ (30) greifen nicht ineinander.** Unterschiedliche Themenlisten, die Klammer Anfang/Ende wird verschenkt. Agenda und `themen.md` nennen Folie 30 zudem anders als ihr Titel. | `08b:26`, `30-was-fehlt.md`, `03:36`, `themen.md:41` | Note 30: „Der Eisberg vom ersten Vormittag, hier die größten Brocken.“ Titel angleichen. | klein |
| W11 | **Cliffhanger 16a wird in Kapitel 3 nicht aufgelöst.** „Die Kundin klickt auf Buchen“ bleibt offen, 19a und 20 greifen es nicht auf. | `16a:9-11`, `19a`, `20-titel` | Ein Satz in Note 19a oder 20-titel. | klein |
| W12 | **Schreibweisen, die auffallen.** Choreografie / Choreographie / Choreography, der / das Storno, Replicas / Instanzen, Backend / Downstream (auf Folie 17 sogar in einer Karte), Auto / Mietwagen. | u. a. `25:43`, `22:14`, `16:38`, `17:16-17`, `03:29` | Je eine Schreibweise festlegen: Pattern-Name englisch, Fließtext „Choreografie“, „das Storno“, „Booking-Instanzen“, „Backend“, „Mietwagen“. | klein |
| W13 | **Spielfolie Story 7 ist dichter als die von Story 4 und 5** (neun Regeln in ganzen Sätzen gegen zwei bis drei Zeilen). Das war die bewusste Entscheidung für Klarheit. Der Spielname „Retry oder Storno“ und der Kassensturz stehen aber nicht auf der Folie. | `24-story-07.md:20-35` | Spielname in den Kicker, Kassensturz als letzte Zeile. Optional die Würfelregeln auf ein Tischblatt auslagern. | klein / mittel |
| W14 | **`themen.md` an weiteren Stellen veraltet.** „Story 2 ist die Ausnahme“ (4, 5 und 7 sind auch ohne Code), Reihenfolge der Blöcke 4 und 5 vertauscht, Recap-Fragen zu Story 1, 5 und 6, die es auf keiner Folie gibt. | `docs/themen.md:20-21`, `:76`, `:131`, `:135`, `:156-161` | Durchgehen und angleichen. | klein |

### Dashboard

| Nr. | Befund | Wo | Vorschlag | Aufwand |
|---|---|---|---|---|
| W15 | **Labels im Story-5-Demo-Text passen nicht zur Oberfläche.** Der Text sagt „failed bleibt 0“, angezeigt wird „fail“. Der Leertext nennt Knöpfe, die anders heißen. | Dashboard `:2057`, `:3189`, `:3204` | Texte an die Beschriftungen angleichen. | klein |
| W16 | **„Downstream-Services“ im Dashboard gegen „Backends“ auf Folien und in Docs**, auch in den User Storys von Story 3, 5 und 7, die wortgleich sein sollen. | Dashboard `:1845`, `:2106`, `:2296` | In den Story-Texten „Backend(s)“. | klein |
| W17 | **Folie 20a sagt „wie im Spickzettel“, nutzt aber andere Namen** (`booked` statt `steps`). Die Referenz heißt wie der Spickzettel. | `20a-saga-code.md:7-17`, Dashboard `:2139-2159` | Folie auf `steps` umstellen. | klein |
| W18 | **Drei Wörter und zwei Stricharten für die Spiel-Stories.** Pille „Spiel“, Kicker „Rollenspiel“ bzw. „Würfelspiel“, Folien-Badge gestrichelt violett, Dashboard gepunktet dunkel. Die Legende sagt „am Tisch“, Story 4 und 5 spielen im Raum. | Dashboard `:904`, `:1675`, `KIND_LABEL`, Folien 15, 18, 24 | Ein Oberbegriff und eine Strichart festlegen. | klein |
| W19 | **Story 4 „Was passieren sollte“: Antworten ohne Fragen.** „Ohne Breaker: nichts.“ und „31.02.: Ja, bei uns schon.“ sind ohne die STOPP-Fragen nicht lesbar. | Dashboard `:1963-1966` | Je Punkt die Frage voranstellen. | klein |
| W20 | **Story 7: „Nichts zu drucken“, die Materialseite druckt aber einen Trainer-Spickzettel.** Story 4 und 5 verlinken die Materialseite, Story 7 nicht. | Dashboard `:2297`, `materialliste.md:56`, `pages/material/index.html:339` | „Für die Tische nichts zu drucken, Trainer-Spickzettel: Materialseite“ mit Link `#story7`. | klein |
| W21 | **Titel der Story-Docs 1 bis 3 weichen von Folie und Dashboard ab** (z. B. „Der erste cloud-native Booking-Service“ gegen „Erst lauffähig, dann smart“). | `docs/stories/story-01…03*.md:1` | Titel angleichen. | klein |
| W22 | **Zehn Links im Browser-Standardblau.** Für Story-Texte, Start-Karte und Hilfe-Link gibt es keine Link-Regel. | Dashboard (`.story-info-body`, `.start-card`, `.hw-help`) | Links auf `#7348E1`. | klein |

### Design über alle Oberflächen

| Nr. | Befund | Wo | Vorschlag | Aufwand |
|---|---|---|---|---|
| W23 | **Zwei Lime-Töne.** `#CCF465` im Dashboard-Footer, am Start-Knoten und in allen Page-Footern, `#D3F871` überall sonst. Die Palette kennt nur den zweiten. | Dashboard `:881`, `:1282`, `:1295`, `:1320`, alle Pages | Überall `#D3F871`. | klein |
| W24 | **Fremde Farben und ein Gate, das sie nicht findet.** Das Dashboard nutzt eine zweite Ampel (`#22c55e`, `#f59e0b`, `#ef4444` …), `#b8860b` für „langsam“, eigene Service-Farben. Das Palette-Gate prüft nur, ob die Kernfarben vorkommen, nicht ob fremde dazukommen, und deckt `pages/` nicht ab. | Dashboard `:1357-1474`, `services/archtest/palette_test.go:52-74` | Health-Punkte auf die Statusfarben umstellen, Service-Farben dokumentieren, Gate optional um eine Positivliste erweitern. | mittel |
| W25 | **Kontrast unter 4,5:1.** `#aaa` und `#888` als kleiner Text auf Weiß (Feedback, Landing, Dashboard), Weiß auf `#b8860b` beim aktiven „Langsam“-Knopf. | `pages/feedback/index.html:101`, `:122`, Dashboard `:256`, `:654` | `#5B4E93` für kleinen Text, „Langsam“ auf `#8a6d00`. | klein |
| W26 | **Kein sichtbarer Tastaturfokus** in Dashboard und Pages (nur im Egg). | Dashboard, alle Pages | Globale `:focus-visible`-Regel. | klein |
| W27 | **Dokumentation der Palette driftet.** `services/CLAUDE.md` nennt Folientext `#000000` (tatsächlich Dunkelblau), „ausschließlich Poppins“ (Sacramento ist eingebunden) und kennt `#4d27a8` nicht. | `services/CLAUDE.md:25`, `:44`, `theme.css:6`, `:29-34` | Tabelle und Schrift-Absatz nachziehen. | klein |
| W28 | **Druckseite hängt an der Systemschrift.** Unter Linux meldet der eigene Überlauf-Detektor einen Überlauf im Trainer-Skript Story 4, Bogen 2. Auf dem Mac vermutlich in Ordnung, aber nicht deterministisch. | `pages/material/index.html:275`, `:1110-1114` | Poppins selbst hosten oder dem Bogen Luft geben. Vor dem Druck einmal am eigenen Rechner prüfen. | mittel |

### Easter Egg

| Nr. | Befund | Wo | Vorschlag | Aufwand |
|---|---|---|---|---|
| W29 | **Überschrift „Downstream-Services“ im Demo-Panel unlesbar** (Kontrast etwa 1,6:1) in Story 4, 5 und 7. Dort hängt das Panel in einem `details.story-info` und erbt dessen Summary-Stil. Wird „Unbedingt“, wenn das Egg im Workshop gezeigt wird. | `voxel.css:356`, `:614`, Dashboard `:3725-3731` | Selektor auf `details.story-info > summary` einschränken. | klein |
| W30 | **`.play-rule` und `.play-lead` bleiben im Basis-Look** (fliederfarbener Kasten und violette Zeile im Holz-Panel), Story 4, 5 und 7. | Dashboard `:1007`, `:1019`, keine Regel in `voxel.css` | Zwei Regeln in `voxel.css` ergänzen. | klein |
| W31 | **Demo-Panel im Egg zu eng.** Zu viel Innenabstand, die Chaos-Knöpfe „Normal“ und „Langsam“ laufen leicht über. | `voxel.css:607`, `:650` | Eigene Regel für `.ref-slot .backend-panel`. | klein |

### Material, Vorbereitung, Orga

| Nr. | Befund | Wo | Vorschlag | Aufwand |
|---|---|---|---|---|
| W32 | **Papier 80 g zu niedrig angesetzt, „Bogen“ meint Seiten.** Die Liste nennt 2 Bogen, es fehlen der Story-7-Spickzettel und das Board als 2× A4 (real 3, ohne A3 5). Die Druckseite zählt bei Duplex Seiten, nennt sie aber Bogen. | `materialliste.md:75`, `pages/material/index.html:336`, `:1035-1040` | Liste korrigieren, auf der Druckseite „Seiten“ schreiben. | klein / mittel |
| W33 | **Druckseite schlägt A3 als Standard vor**, obwohl A3 laut Liste die Ausnahme ist. | `pages/material/index.html:325-326` | Reihenfolge tauschen, 2× A4 als Standard. | klein |
| W34 | **Legende „Art“ veraltet.** „3D“ kommt in keiner Zeile mehr vor, „Büro“ steht zehnmal in der Tabelle, aber nicht in der Legende. | `materialliste.md:9` | Legende auf kaufen, Büro, drucken, Raum. | klein |
| W35 | **Kreppband als Pflichtkauf, Plenum-Flipchart ohne Gegenstück.** Kreppband braucht es nur, wenn Bogen an die Wand sollen. Das Plenum-Flipchart („Erwartungen und Schnittfrage“) kommt im Ablauf nicht vor. | `materialliste.md:67-68`, `:83`, `themen.md:48` | Kreppband optional, Plenum-Flipchart streichen oder im Ablauf verankern. | klein |
| W36 | **Toter Link auf der Vorbereitungsseite.** Verweise aus dem Troubleshooting auf `vorbereitung.md` springen zu `#schritt-0`, diese ID gibt es nicht. | `pages/vorbereitung/index.html:353` | Fallback auf eine vorhandene ID. | klein |
| W37 | **Story-7-Doc nennt 20 und 25 Minuten Spiel.** Kopf, Folie, Dashboard und Spickzettel sagen 25. | `docs/stories/story-07-choreography-saga.md:22` | „(ca. 25 Minuten inklusive Kassensturz)“. | klein |
| W38 | **Einladungs-Mail kündigt Punkte an, die nicht folgen.** „Erledigt bitte vorab die folgenden Punkte:“, danach steht nur die URL. Der Intranet-Text nennt die drei Schritte. | `docs/orga/einladung-email.md:27-30` | Die drei Schritte übernehmen. | klein |

---

## 4. Nice-to-have

Feinschliff ohne Risiko, wenn Zeit ist.

> **Stand 13.11.:** Erledigt sind N1 bis N4, N6 bis N10, N12 bis N17 und N19.
>
> - **N1:** Leerer Stapel aufgelöst, Story-4-Recap mit Badge. Fünf Fragen in Story 8, Titelfolien ohne Zitat und Story 7 ohne Code-Folie bleiben bewusst.
> - **N2:** Nur die Notes zum Story-7-Recap sind gekürzt und nummeriert. Die Spielskripte zu Story 4 und 5 bleiben lang.
> - **N5:** Zuordnung der Karten vom Anfang, Feedback-Folie am Schluss und Note zu Timeouts und Retry sind drin. Die Domäne führt die neue Folie „Unsere Reiseplattform“ ein, direkt vor „Einklinken, dann los“ (damit 64 Folien). Die Transferfrage pro Recap bleibt bewusst weg.
> - **N9:** Überschriften in Poppins, Fließtext bleibt Systemschrift.
> - **N16:** Der eine Link auf `larmic.de` nutzt jetzt wie alle anderen `larmic.github.io`. Welche Adresse die kanonische sein soll, ist nicht entschieden.
>
> - **N11:** Entschieden: Die Emojis bleiben, wo sie sind.
> - **N18:** Regel: „ihr“ für alles, was im Workshop gelesen wird, „du“ für die Vorbereitung allein zu Hause, durchgehend „Teilnehmende“. Umgesetzt in Docs und READMEs. Nicht angefasst: das verallgemeinernde „du“ in den Hintergrundtexten (`docs/questions/`, `docs/instructions/`), das Feedback-Formular, die Intranet-Vorlage sowie `.claude/` und `docs/feedback/`.
>
> Damit ist aus der Nice-to-have-Liste nichts mehr offen.

| Nr. | Bereich | Befund | Wo | Aufwand |
|---|---|---|---|---|
| N1 | Folien | Strukturabweichungen zwischen den Story-Blöcken: Story 7 als einzige ab Story 3 ohne Code-Folie, Story 8 mit fünf Fragen, Titelfolien 23 und 26 ohne Zitat, 16 ohne „↓“-Badge auf 16a, leerer Stapel um Folie 18 in `index.html:91-93`. | Deck | klein |
| N2 | Folien | Lange Notes: 15 (561 Wörter) und 18 (496) sind Spielskripte mit doppelter Pflege zur Materialseite. 25, 24, 29a, 17 und 25b liegen über 150 Wörtern. Note-Label „Nie angekommen“ passt nicht zum Folien-Label „Ein Chip bleibt liegen“. | `15`, `18`, `25:55` | klein bis mittel |
| N3 | Folien | Intro-Grafik kompensiert nur Flight, nach der Saga-Logik müsste auch Hotel zurück. | `01-intro.md:25-31` | klein |
| N4 | Folien | Kleinigkeiten: Hook „Flug gebucht, Hotel sagt nein“ dreimal fast wortgleich (20-titel, 21), „Regel-Folie eben“ nach 20 Minuten Flipchart (10f), Schritt-Status nur in Notes. | `20-titel:9`, `21:51`, `10f:13` | klein |
| N5 | Folien | Offen aus dem Review vom 05.11.: Domäne wird nie auf einer Folie eingeführt, die drei Karten von Folie 08 kehren auf 29 nicht wieder, schwaches Ende (Rücksprung zur Feedback-Folie), keine „Und bei euch?“-Frage auf den Recaps, Timeouts und Retry als eigenes Thema fehlen. | Deck | mittel bis hoch |
| N6 | Dashboard | Story 6: Reihenfolge der Akzeptanzkriterien weicht von Folie und Doc ab, ein Kriterium fehlt. | Dashboard `:2124-2128` | klein |
| N7 | Dashboard | Reste: `resetSagas6()` für Story 7, Kommentar „Rollenspiele“, Story-7-Doc spricht von „Pflicht-Variante“ ohne Kriterienliste, Story-4-Doc verspricht Chaos „pro Replica“. | Dashboard `:2262`, `story-07…md:145-161`, `story-04…md:190` | klein |
| N8 | Dashboard | Kleinkram: „durch sind“ gegen „geklappt haben“ beim Kassensturz, Panel-Überschriften uneinheitlich, Stepper-Tooltips nur bei fünf von neun Schritten. | Dashboard, `24:52` | klein |
| N9 | Design | Poppins nur auf den Folien, Dashboard und Pages mit Systemschrift. Entweder als bewusst dokumentieren oder Poppins für Überschriften einbinden. | Dashboard, Pages | mittel |
| N10 | Design | Inline-Code und Hervorhebung uneinheitlich (Folie: Lime-Marker, Dashboard: violettes fettes `em`, das wie ein Link aussieht). | Dashboard `:995` | klein |
| N11 | Design | Emojis nur außerhalb der Folien (Landing, Feedback, Dashboard-Kopf, Service-Karten). | Pages, Dashboard | mittel |
| N12 | Design | Druckseite auf dem Handy: Tabelle und Auswahlfelder ragen über. Feedback: Stern-Hinweis klebt am Namensfeld. | `pages/material`, `pages/feedback` | klein |
| N13 | Egg | Tote Regeln `.bh-btn`, `.bh-actions` in beiden Dateien, Saga-Liste zeigt im Egg eine Zeile weniger, Sprechblase behält ihre Position und kann bei schmalem Fenster Scrollen erzeugen. | `voxel.css:102`, `:123`, `:514`, Dashboard `:3811-3829` | klein |
| N14 | Material | Umschläge-Menge stammt aus der Karten-Zeit von Story 7, Randangabe 3 mm gegen 4 mm, toter Code (Falz, Duplex kurze Kante, `cards70`). | `materialliste.md:78`, `pages/material/index.html:338`, `:343`, `:927` | klein |
| N15 | Orga | Titel „für die Rollenspiele“ in Materialliste, Landing und 3D-README, gemeint sind auch Würfelspiel und Design-Session. | `materialliste.md:1`, `pages/landing/index.html:153` | klein |
| N16 | Orga | Zwei Domains für dieselbe Site: `larmic.github.io/…` leitet auf `larmic.de/…` weiter. Beides funktioniert, die Folie zeigt nicht die kanonische Adresse. | Dashboard `:1685`, `02-feedback.md` | klein |
| N17 | Orga | READMEs, Workflow-Name und `CLAUDE.md` kennen die Materialseite nicht. | `README.md:33-49`, `services/README.md:5`, `deploy-pages.yml:1` | klein |
| N18 | Orga | Anrede und Gendern uneinheitlich: du / ihr gemischt (`custom-setup.md` in einer Datei), Teilnehmer / Teilnehmende / Teilnehmer:innen über die Docs verteilt. Die Notes-Regel „ohne Gendern“ gilt nur für die Notes. | Docs, Pages | mittel |
| N19 | Vorbereitung | Troubleshooting verweist auf „Windows-Hinweise“ in `vorbereitung.md`, die es so nicht gibt. Folie 08c zeigt den Schritt „im Team auf einen Service einigen“ nicht. | `docs/troubleshooting.md:196`, `:238`, `08c-start.md:15-20` | klein |

---

## 5. Empfohlene Reihenfolge

1. **U1 bis U7** in einem Rutsch: alle klein, zusammen etwa eine Stunde. U3 zunächst als Textänderung.
2. **Story-5-Demo** einmal am laufenden System durchspielen (U3, U4, W15), sie hat die meisten Abweichungen.
3. **Kleine Folien-Korrekturen** W1, W2, W4, W5, W6, W11, W12: alle klein, zusammen ein bis zwei Stunden.
4. **`themen.md` und Materialliste durchgehen** (W9, W14, W32 bis W35): dort sitzen die meisten veralteten Stellen.
5. **Egg** W29 bis W31, falls es gezeigt wird.
6. Danach nach Wahl: W3 (Kostenfrage), W8 (Überleitungen am ersten Vormittag) und W24 (Farben und Gate) sind die drei Punkte mit dem größten Hebel und mittlerem Aufwand.
