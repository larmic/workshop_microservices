# Teilnehmer-Review „Microservices richtig gemacht“

Stand: 13.11.2026, Deck mit 64 Folien.

**Wer hier schreibt:** ein durchschnittlicher deutscher Entwickler, etwa sechs
Jahre Projekterfahrung (Java/Spring oder Node, kein Go), angehender Architekt,
in einem Projekt ein wenig mit Microservices gearbeitet. Ich kenne nur die
Einladung, den Inhalt nicht.

**Wie das Review entstanden ist:** Der Workshop wurde virtuell durchgespielt,
in acht Abschnitten (Einladung und Vorbereitung, Tag 1 bis Story 2, Story 3 bis
5, Story 6 und 7, Story 8 und Abschluss, Lösbarkeit der Code-Stories, Design und
Easter Egg, Mischung und Zeit). Die Folien wurden als „das sehe ich“, die Notes
als „das höre ich einmal“ gelesen. Jeder Kritikpunkt wurde danach von einer
zweiten Instanz gegen die Dateien geprüft; die Schwere unten ist die nach der
Prüfung. Frühere Reviews wurden nicht gelesen.

**Was nicht geprüft ist:** Der Workshop wurde nicht mit echten Menschen
gespielt. Die Code-Stories 1, 3 und 6 wurden in Node wirklich gebaut und gegen
die echten Backends gestartet, Story 8 nur geplant; ein Java-Team wurde nicht
nachgestellt. Die Zeiten innerhalb der Blöcke sind Schätzungen. Windows wurde
nicht geprüft.

## Kurzfazit

Ich würde den Workshop empfehlen: Kolleginnen und Kollegen, die schon einmal
REST zwischen Services gebaut haben und jetzt verstehen wollen, was im
Fehlerfall passiert. Wer wissen will, wie man Services schneidet, wird nicht
bedient, obwohl die Einladung genau damit lockt.

Das Material ist ungewöhnlich sorgfältig. Jedes Pattern kommt erst, nachdem ich
das Problem erlebt habe. Ein Service wächst über zwei Tage, das Dashboard zeigt
mir sofort, ob es klappt, und der Schluss ist rund.

Drei Dinge hätten mich im Raum wirklich Zeit oder Nerven gekostet:

1. **Der Log-Filter in Story 8 läuft nicht wie gedruckt.** Der Befehl ist das
   Erfolgserlebnis der Story und bricht mit `invalid compose project` ab.
2. **Story 3 passt nicht in 60 Minuten.** Die Notes rechnen allein 30 + 30
   Minuten Arbeit, dazu kommen vier Theoriefolien, Demo und Recap, am Ende eines
   überbuchten ersten Tages.
3. **In Story 6 bleibt das Saga-Panel für „Mein Service“ leer**, weil es zwei
   Endpunkte braucht, die nirgends im Teilnehmermaterial stehen. Die Meldung
   „nicht erreichbar“ schickt mich auf die falsche Fährte.

## Antworten auf die Fragen

| Frage | Antwort |
|---|---|
| Nachvollziehbar? | Überwiegend ja. Lücken bei den Spielregeln von Story 5 und 7 im Detail und beim Dashboard-Vertrag in Story 6. |
| Lösbar? | Ja, der Kern jeder Code-Story. Story 1 hat Luft, Story 3, 6 und 8 sind knapp. |
| Anspruch angemessen? | Passend. Story 1 eher leicht, Story 4 als „Request“ im Spiel unterfordernd, Story 5 bei der Dimensionierung zu schnell. |
| Passt das Design? | Ja, Folien, Dashboard und Pages wirken wie ein Produkt. Ausreißer: die Spielfolie von Story 7 ist eine Textwand. |
| Mischung Folien, Dashboard, Stories, Spiele? | Gut verteilt, aber Tag 1 ist überbucht und Tag 2 hat fast drei Stunden Luft. Ich programmiere rund drei Stunden. |
| Fehlt mir etwas für meine Arbeit? | Service-Schnitt, ein Handgriff für Circuit Breaker und Bulkhead im eigenen Stack, Weiterlese-Hinweise. |
| Empfehlen? | Ja, mit dem Hinweis, dass es um Betrieb und Kommunikation geht, nicht um den Schnitt. |
| Easter Egg gefunden? | Am ersten Tag vermutlich nicht. Der einzige sichtbare Hinweis ist eine blasse Spitzhacke im Dashboard-Footer; deutlicher ist die Meldung in der Browser-Konsole. |

## Durchgespielt

### Vor dem Workshop

Die Einladung holt mich ab: „verteilte Monolithen“ und „unklare Grenzen“ kenne
ich. Die Teilnahmebedingungen sind konkret, ich erkenne mich wieder. Was ich
nicht erfahre: welche Themen drankommen, dass im Team gearbeitet wird und dass
es Spiele gibt. Meinem Teamleiter könnte ich nicht in einem Satz sagen, was ich
danach kann.

Die Vorbereitungsseite ist die beste, die ich für einen Workshop gesehen habe:
drei Schritte, Checklisten, Hilfe direkt am Schritt, und am Ende sagt mir der
Knopf „Hausaufgabe prüfen“: „Vor dem Workshop ist sonst nichts zu tun.“ Nach gut
einer Stunde bin ich fertig. Gestolpert bin ich an Schritt 2: „zieht fertige
Images“ stimmt nicht ganz, es wird zusätzlich ein Kotlin-Service mit Gradle
gebaut, und „läuft ohne Fehler“ kann ich an durchrauschenden Logs nicht ablesen.

### Tag 1, Vormittag

Der Einstieg ist ehrlich („Brauchen wir Microservices?“ mit einem echten
Projekt). Die neue Folie „Unsere Reiseplattform“ sagt mir in einem Bild, was ich
zwei Tage lang baue. Bis zur ersten Zeile Code vergehen laut Plan aber 165
Minuten, davon 45 für die Zeitlinie.

Story 1 ist klar und für mich eher leicht: Allein wäre ich in 30 bis 45 Minuten
durch, angesetzt sind 90. Eine Aufgabe für schnelle Teams gibt es nicht.
Gleichzeitig passt der Vormittag nicht vor die Mittagspause, die auf der
Agenda-Folie steht: Bei Beginn 9:00 endet er rechnerisch um 13:15.

Story 2 (REST am Flipchart) und das Quiz sind fair und eine gute Abwechslung
nach dem Essen.

### Story 3 bis 5

Story 3 ist genau mein Niveau: neues Werkzeug, bekannte Technik, und die Folie
„Drei HTTP-Aufrufe“ nimmt mir die Angst vor Consul. Resolver und Angebote
schaffe ich, die Buchung per POST wird gehetzt oder bleibt liegen. Auf die baut
Story 6 am nächsten Morgen auf.

Das Circuit-Breaker-Spiel verstehe ich beim ersten Hören. Lautes Zählen gegen
„sofort“ sitzt. Als Request lese ich allerdings nur vor.

Beim Bulkhead-Spiel verliere ich den Faden an zwei Stellen: Das Abklatschen von
Flight und Car höre ich einmal und sehe es nirgends, und die Umrechnung zwischen
Nummern, Stühlen und Slots rauscht vorbei. Die wichtigste Aussage („Die Größe
ist die Arbeit“) bleibt bei mir deshalb vage.

Nach dem ganzen Resilience-Kapitel habe ich keine Zeile dazu geschrieben. Die
Begründung („in Produktion nimmt man die Library“) überzeugt mich, aber dann
will ich die Library einmal in meinem Stack konfiguriert sehen.

### Story 6 und 7

Die Saga ist gut hergeleitet, der Pseudo-Code auf „Der Orchestrator“ lässt sich
fast eins zu eins umsetzen. Netto bleiben mir geschätzt gut 30 Minuten zum
Programmieren. Dann schalte ich auf „Mein Service“ und das Panel meldet „Booking
Story 6 nicht erreichbar“, obwohl mein Service läuft.

Das Würfelspiel macht Spaß und ist schnell verstanden. Beim ersten 503 bleibt
aber offen, ob danach neu gewürfelt wird und wohin der übrige Chip geht. Die
Demo danach zeigt dasselbe Bild wie in Story 6; das stille Storno, um das es
geht, sehe ich nur als Go-Code.

### Story 8 und Abschluss

Die Folie „Vier Stellen im Code“ und der Spickzettel ergeben einen klaren
Bauplan. Ich kopiere den Filter-Befehl aus dem Panel und bekomme `service
"flight" depends on undefined service "consul": invalid compose project`. Fünf
Minuten weg. Danach stehen Buchung und Storno über vier Services untereinander
in einem Block, und das ist der stärkste Moment des zweiten Tages.

„Acht Stories, acht Sätze“ und „Drei Sätze“ schließen rund ab, das Conway-Zitat
ist ein würdiger Schluss. Bei „Was wir nicht angeschaut haben“ stutze ich:
Service-Schnitt steht dort an erster Stelle. Das Feedback-Formular ist in einer
Minute ausgefüllt.

## Befunde

Schwere nach der Gegenprüfung. Aufwand ist meine Schätzung für die Behebung.

### Hoch

| # | Befund | Beleg | Aufwand |
|---|---|---|---|
| T1 | **Story 8: Filter-Befehl läuft nicht wie gedruckt.** `docker compose logs …` ohne `-f` lädt nur `docker-compose.yml` und bricht ab. `COMPOSE_FILE` steht nur in `docs/troubleshooting.md:29`. Nachgestellt, mit und ohne Export. Betroffen sind Story, Dashboard-Panel, die Folie „Den roten Faden im Log“ und eine Note, dazu derselbe Befehl in Story 4. | `docs/stories/story-08-tracing.md:70-71`, `services/dashboard/static/index.html:3439`, `services/slides/chapters/27-story-08.md:42`, `docs/stories/story-04-circuit-breaker.md:197` | klein (z. B. `COMPOSE_FILE` in `.env.example` oder ein `make logs`) |
| T2 | **Story 3: 60 Minuten reichen nicht.** Die Notes verplanen 30 + 30 Minuten reine Arbeit, der Block enthält zusätzlich vier Theoriefolien, Demo und vier Recap-Fragen. Tag 1 hat keinen Puffer. | `services/slides/chapters/12-story-03.md:9` und `:59-60`, `docs/themen.md:25` | mittel (Zeitplan) |

### Mittel

| # | Befund | Beleg | Aufwand |
|---|---|---|---|
| T3 | **Story 6: Saga-Panel für „Mein Service“ braucht undokumentierte Endpunkte.** Das Dashboard ruft `GET /admin/sagas` und `POST /admin/sagas-reset` am eigenen Service. Beide Pfade und das erwartete JSON stehen nur in der Referenz-Spec. Fehlen sie, heißt es „nicht erreichbar“. Auch der Pfad der Statusabfrage und der Statuscode im Fehlerfall sind nicht genannt. | `services/dashboard/handler/saga.go:13` und `:29`, `services/dashboard/static/index.html:3545`, `docs/stories/story-06-saga.md:22` | klein (Spickzettel ergänzen, Meldung präzisieren) |
| T4 | **Tag 1 ist überbucht, Tag 2 hat fast drei Stunden Luft.** Tag 1 plant 430 Minuten ohne Kaffeepause, Tag 2 nur 255. Die Agenda-Folie nennt Mittag um 12:30, der Vormittag endet laut Plan erst um 13:15. | `docs/themen.md:15-30`, `services/slides/chapters/03-agenda.md:16` | mittel (Block nach Tag 2 schieben) |
| T5 | **Bei Story 3 und 6 ist unklar, was in „≈ 60 min“ steckt.** Die Spiele schlüsseln die Zeit auf, Story 8 sagt „inkl. Demo“, Story 3 und 6 nicht. Netto bleiben in Story 6 geschätzt gut 30 Minuten. | `services/slides/chapters/21-story-06.md:9`, `docs/stories/story-06-saga.md:4` | klein |
| T6 | **Story 1 ist für 90 Minuten zu dünn, eine Zusatzaufgabe fehlt.** Die Referenz hat 115 Zeilen. Kürzen würde nebenbei T4 entschärfen. | `docs/stories/story-01-cloud-native-booking-service.md:7` | klein |
| T7 | **Kein Aufholweg von Story 3 nach Story 6.** Story 8 hat einen Abschnitt „Story 6 noch nicht fertig?“, Story 6 hat nichts Vergleichbares. Die mitgelieferte Beispiel-Lösung deckt nur Story 1 ab. | `docs/stories/story-06-saga.md`, zum Vergleich `docs/stories/story-08-tracing.md:79` | klein |
| T8 | **Einladung nennt weder Inhalte noch Format.** Kein Pattern, keine Teams, keine Spiele. | `docs/orga/intranet.md:14` | klein |
| T9 | **Die Einladung wirbt mit Service-Schnitt, der Workshop behandelt ihn nicht.** „Zu kleinteilige Services, unklare Grenzen“ ist der Aufhänger; am Ende steht das Thema als erste Karte auf „Was wir nicht angeschaut haben“. | `docs/orga/intranet.md:12`, `services/slides/chapters/30-was-fehlt.md:13-14` | klein (Einladung) oder hoch (Übung dazu) |
| T10 | **Story 4 und 5: kein Handgriff für den eigenen Stack.** Ein Resilience4j-Beispiel steht nur unter „Zum Nachbauen zu Hause“ und zählt anders als das Spiel; für Bulkhead gibt es kein Konfigurationsbeispiel. Die Luft an Tag 2 böte Platz dafür. | `docs/stories/story-04-circuit-breaker.md:142` und `:176`, `services/slides/chapters/13a-kapitel-2-resilience.md:19` | mittel |
| T11 | **Story 5: Abklatsch-Regeln stehen nicht auf der Folie.** Dass ein Request Flight, Hotel und Car nacheinander ruft, wird vor dem Spiel nirgends gesagt. | `services/slides/chapters/18-story-05.md:27` gegenüber `:44-47` | klein |
| T12 | **Story 5: drei Maßstäbe ohne sichtbare Umrechnung.** Nummern, Stühle, Slots; die Tabelle steht nur in der Story-Datei. | `services/slides/chapters/17-bulkhead.md:38`, `docs/stories/story-05-bulkhead.md:205-210` | klein |
| T13 | **Story 5 Demo: Flight-Rejects sind möglich**, das Dashboard verspricht aber „Flight und Car laufen weiter“. Aus Texten und Code abgeleitet, nicht live gesehen. | `services/dashboard/static/index.html:2038`, `docs/stories/story-05-bulkhead.md:191` | klein bis mittel |
| T14 | **Story 7: Der Retry-Fall ist auf der Folie nicht eindeutig.** Wird nach dem 503 neu gewürfelt? Wohin geht der übrige Reisegeld-Chip? Die Antwort steht nur im Trainer-Spickzettel und in der Story-Datei. Verteilt ein Tisch falsch, stimmt der Kassensturz nicht. | `services/slides/chapters/24-story-07.md:33`, `docs/stories/story-07-choreography-saga.md:67` | klein |
| T15 | **Story 7: Die Demo zeigt das stille Storno nicht.** Das Material sagt es selbst; im Dashboard sehe ich dasselbe Bild wie in Story 6. | `docs/stories/story-07-choreography-saga.md:114` | hoch (Demo umbauen) oder klein (Demo kürzen) |
| T16 | **Teamarbeit beim Programmieren ist nicht geregelt.** Teamgröße steht nur bei Story 2; wer tippt, ob gewechselt wird, sagt niemand. | `services/slides/chapters/08c-start.md:43` | klein |
| T17 | **Ein Service pro Team, aber das Dockerfile muss dann überall bauen.** Kopiert es ein lokal gebautes Artefakt (bei Java üblich), bauen die anderen ohne Toolchain nicht. Ein Hinweis auf einen mehrstufigen Build fehlt. | `docs/vorbereitung.md:28` und `:36` | klein |
| T18 | **„Zieht fertige Images“ verschweigt den lokalen Gradle-Build.** Hinter einem TLS-aufbrechenden Proxy scheitert er auf jedem Betriebssystem; die passende Hilfe heißt „Nur Windows/WSL2“. | `docs/vorbereitung.md:72`, `services/Makefile:62-64` | klein |
| T19 | **Windows-Voraussetzungen stehen nicht bei den Voraussetzungen.** WSL2 taucht erst im Kasten von Schritt 2 auf. | `docs/vorbereitung.md:11-16` und `:63` | klein |

### Klein

Einzeln geprüft, alle bestätigt, alle mit wenig Aufwand zu beheben.

**Vorbereitung und Landing**
- „Der Stack läuft ohne Fehler“ ist nicht prüfbar formuliert; wie man den Stack beendet, steht nur in der Hilfe.
- Port 8099 antwortet nach dem Start, aber mit der Beispiel-Lösung, nicht mit meinem Service.
- Die Einladung sagt „anspruchsvoll“ und „Vorwissen brauchst du nicht“.
- Die Kachel „Spielmaterial“ auf der Landing-Seite ist Trainer-Material und nicht so gekennzeichnet.
- Die Spring-Falle `/actuator/health` wird erst im Fehlerfall erklärt.

**Tag 1**
- Die Zeitlinie zeigt SCS, CORBA und RMI ohne Erklärung.
- Quiz 1 sagt, die Filter seien in Story 1 schon gebaut.
- Die Musterlösung von Story 2 lässt Fragen zu Idempotenz und Status-Codes offen; die Pfad-Präfixe wechseln zwischen Regel-Folie, Musterlösung und Quiz.
- Bei „Logs auf stdout“ und OpenAPI ist unklar, wann das Kriterium erfüllt ist.
- Die Backend-Pfade stehen nur im zugeklappten Spickzettel.
- „Die 12 Faktoren“: zwölf Kacheln, die Code-Chips sind von hinten nicht lesbar.

**Story 3 bis 5**
- `CONSUL_URL` wird in Story 3 nicht erwähnt; der Request-Body der Backend-Buchung steht nirgends im Story-Text.
- Die Folie „Die State-Machine“ nennt CAS ohne Erklärung.
- Die „5 s“ auf den Karten 01 bis 03 haben für mich keinen Bezug zu meinem Service.
- Das Badge „Und jetzt schreiben“ ist missverständlich.
- Story 5 hat einen falschen Recap-Verweis; „drei von elf“ und „komplett mit Hotels: 2“ passen nicht zusammen.

**Story 6 und 7**
- Das Akzeptanzkriterium zu DELETE liest sich, als müsste ich es in den Backends bauen.
- Wie ich in Story 6 gezielt einen Fehler auslöse, steht nur in der Trainer-Note.
- Forward und Backward Recovery stehen auf keiner Folie.
- Recap 02 (doppeltes Storno) wirkt wie ein Widerspruch zum idempotenten DELETE aus Story 6.
- Die Titelgrafik von Story 7 zeigt einen Broker, gebaut ist ein Zuruf per POST.
- Die Spielfolie von Story 7 ist eine Textwand, die von Story 4 und 5 sind fast leer.

**Story 8 und Abschluss**
- Die technischen Hinweise der Story sprechen Go (`slog`, „die Library“ mit `Middleware` und `Propagate`) und verweisen auf Trainer-Notizen.
- Der Merksatz „nur am Entry-Point erzeugen“ steht in der Zusammenfassung ohne die Einschränkung aus der Story.
- „Was wir nicht angeschaut haben“ nennt keine Links oder Bücher; Testing, Kubernetes und Messaging mit echtem Broker fehlen dort ganz.
- Das Feedback-Formular zeigt die Leitfragen der Folie nicht.
- Im Modus „Mein Service“ bleiben Chip und Button-Beschriftung auf der Referenz stehen.

**Design und Easter Egg**
- Das Landing-Banner hat einen anderen Stil als die Titelfolie.
- Die Agenda zeigt weder Start- noch Endzeit; Kaffeepausen sind nicht geplant.
- Der Hinweis auf das Easter Egg ist praktisch unsichtbar.
- Im Voxel-Theme sind Ziffern und Buchstaben in der Pixelschrift schwer zu unterscheiden, und der Inhalt beginnt weit unten.

## Mischung und Zeit

Blocksummen aus `docs/themen.md`; die Aufteilung innerhalb der Blöcke ist
geschätzt.

| Arbeitsform | Minuten | Anteil |
|---|---|---|
| Folien und Vortrag | 160 | 23 % |
| Dashboard und Demo | 75 | 11 % |
| Eigenes Programmieren | 178 | 26 % |
| Spiel und Design-Session | 90 | 13 % |
| Recap und Diskussion | 182 | 27 % |
| Summe | 685 | Tag 1: 430, Tag 2: 255 |

Nach dem ersten Vormittag höre ich nie länger als 10 bis 15 Minuten am Stück
zu, das ist gut. Die Arbeitsformen wechseln sinnvoll (Flipchart nach dem Mittag,
Spiel am Ende von Tag 1). Drei Spiele plus Design-Session sind von der Menge her
richtig; Story 4 und 5 liegen direkt hintereinander.

Rund drei Stunden eigenes Programmieren in zwei Tagen sind für „hands-on“ am
unteren Rand. Die Luft an Tag 2 ist der naheliegende Ort für mehr.

## Offene Fragen aus dem Raum

- Welche Themen kommen dran, und arbeite ich allein oder im Team?
- Wann fängt der Tag an, wann ist wirklich Mittag?
- Was mache ich mit meinem eigenen Service, wenn das Team den eines Kollegen nimmt?
- Meint „≈ 60 min“ die reine Arbeitszeit?
- Story 3: Nehme ich `consul:8500` fest in den Code? Cacht man die Abfrage in echt nicht?
- Story 4: Wenn der Fallback 200 liefert, widerspricht das nicht „Fehler kommen als Status-Code“?
- Story 6: Was muss mein Service liefern, damit das Panel etwas zeigt? Welcher Statuscode bei kompensierter Saga?
- Story 7: Warum ist ein doppeltes Storno ein Problem, wenn DELETE idempotent ist?
- Story 7: Booking bekommt beim Zuruf doch einen Fehler zurück. Warum setzt es trotzdem COMPENSATED?
- Story 8: Wie viel davon baue ich im echten Projekt selbst, wenn Spring mit Micrometer oder OTel das automatisch macht?
- Bekomme ich Folien und Merksätze zum Mitnehmen, und wo lese ich weiter?
- Das Dashboard meldet das gefundene Easter Egg an einen Webhook. Erfahre ich das irgendwo außer im Quelltext?

## Was mir gefallen hat

- Problem vor Pattern, in jeder Story. Ich erlebe erst den Schmerz, dann die Lösung.
- Ein Service über zwei Tage: Am Ende habe ich etwas Zusammenhängendes gebaut.
- Die Vorbereitung mit eindeutigem Fertig-Signal und einer Fehlertabelle, die genau meine typischen Fehler trifft.
- Die Spickzettel im Dashboard sind sprachneutraler Pseudo-Code und an den kniffligen Stellen genau richtig (Form der Consul-Antwort, Saga in drei Schritten, Tracing in fünf).
- Story 1 und 3 liefen mit Story-Text und Spickzettel beim ersten Versuch gegen die echten Backends.
- Die Backends sind gut vorbereitet: idempotentes DELETE, der reservierte Fehlerfall C9, fertige Übernahme des `traceparent`.
- Das Circuit-Breaker-Spiel mit lautem Zählen und das Würfelspiel „Retry oder Storno“.
- Die Folien „Unsere Reiseplattform“, „Drei HTTP-Aufrufe“, „Der Orchestrator“ und „Vier Stellen im Code“.
- „Acht Stories, acht Sätze“ und „Drei Sätze“: Das kann ich am Montag wörtlich in ein Architektur-Review mitnehmen.
- Die Ehrlichkeit: ein Auftakt mit „Brauchen wir das überhaupt?“, eine Folie über das, was fehlt, und Windows-Hilfen, die als ungeprüft markiert sind.
- Das Easter Egg ist liebevoll gemacht, stört nicht und lässt sich mit demselben Code wieder abschalten.
