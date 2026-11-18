# Konzept: Einklinken in „Start (0)“ vereinfachen

Stand: 18.11.2026, während eines laufenden Workshops.

**Worum es geht:** Der Schritt „Start (0) – eigenen Service einklinken“ läuft
im Workshop holprig. Dieses Dokument wertet die Beobachtungen des Trainers aus,
vergleicht drei Lösungswege und empfiehlt eine Sofortmaßnahme für den laufenden
Workshop sowie einen Umbau danach.

**Was nicht passiert ist:** Es wurde nichts umgesetzt. Alle Aussagen sind gegen
die Dateien und mit `docker compose --dry-run` bzw. `config` (Docker 29.7.2,
Compose 5.7.1) geprüft. Windows und ältere Compose-Versionen wurden nicht
nachgestellt.

## Kurzfazit

1. **Die `image:`-Zeile in `docker-compose.custom.yml` ist der eigentliche
   Fehler.** Sie hat für den Ablauf keinen Zweck, versteckt aber den
   eingeklinkten Service im Hausaufgaben-Check. Außerdem lässt sie ein nacktes
   `docker compose pull/up` die Beispiel-Lösung vom Hub ziehen. Sie kann
   ersatzlos raus, und das wirkt ohne neues Dashboard-Image.
2. **Der Port 8099 ist nur das Symptom.** Das eigentliche Missverständnis
   lautet: „Mein Service läuft irgendwo“. Dashboard und Gateway erreichen ihn
   aber nur als Container `booking-custom` im Compose-Netz. Die Vorbereitung
   trainiert mit `docker run -p 8099:8080` genau den falschen Reflex.
3. **Die `.env` kann für Teilnehmer zur Datei werden, die man nie anfasst.**
   Der Pfad gehört dann direkt in `docker-compose.custom.yml`, wie der Trainer
   vorschlägt. Ganz weg fällt die `.env` aber nicht, Story 8 braucht
   `COMPOSE_FILE`.
4. **Einklinken als Hausaufgabe:** technisch möglich und sinnvoll als
   Generalprobe. Es hebt aber die Entscheidung D5 auf und spart im Workshop nur
   teilweise Zeit, weil die meisten auf den Team-Service umbiegen.

## Anlass

Beobachtung im Workshop:

- Die Anpassung der `.env` (`CUSTOM_BOOKING_PATH`) wirkt umständlich und ist
  vielleicht unnötig.
- Teilnehmer starten ihren Service „irgendwo“ (per `docker run -p 8099:8080`,
  aus der IDE). Der Stack belegt 8099 aber schon mit der Beispiel-Lösung, und
  der Start scheitert.
- Wer die `image:`-Zeile in `docker-compose.custom.yml` nicht auskommentiert,
  sieht sein neues Image beim Test „Hausaufgabe prüfen“ im Dashboard nicht.
  Das betrifft nur einige.

## Ist-Zustand

### Ablauf aus Sicht eines Teilnehmers

1. **Hausaufgabe:** eigenes Projekt, `docker build -t my-booking-service .`,
   `docker run --rm -p 8099:8080 my-booking-service`, `curl localhost:8099/health`
   (`docs/vorbereitung.md:43-50`). Der Text sagt: „Unter `localhost:8099` ist
   dein Service auch später im Workshop erreichbar“ (`docs/vorbereitung.md:53`).
   Danach soll der Container gestoppt werden, weil der Stack 8099 belegt
   (`docs/vorbereitung.md:55`).
2. **Vorbereitung:** Repo klonen, `cp .env.example .env`, `make docker-up-hub`.
   „Die .env funktioniert so, wie sie ist“ (`docs/vorbereitung.md:70-74`). Der
   Stack startet die Kotlin-Beispiel-Lösung als `booking-custom` auf 8099.
3. **Dashboard:** „Hausaufgabe prüfen“ sucht `my-booking-service`
   (`services/dashboard/handler/homework.go:148-152, 208-213`).
4. **Workshop Start (0):** Das Team einigt sich auf einen Service, die anderen
   klonen ihn (`services/slides/chapters/08c-start.md`, Notes). In
   `services/.env` wird `CUSTOM_BOOKING_PATH` absolut eingetragen (Windows:
   Forward-Slashes, LF). Danach Ctrl+C, `make docker-up-hub`,
   `curl localhost:8099/health` und in `docker logs booking-custom` prüfen,
   dass kein BEISPIEL-LÖSUNG-Kasten mehr erscheint
   (`services/dashboard/static/index.html:1694-1720`, `docs/custom-setup.md:12-40`).
5. **Story 1 ff.:** Code ändern, `make docker-rebuild-custom`.

### Mechanik

| Baustein | Stand | Beleg |
|---|---|---|
| Build-Kontext | `context: ${CUSTOM_BOOKING_PATH}`, ohne Default | `services/docker-compose.custom.yml:5` |
| Image-Name | `image: larmic/workshop-microservices-booking:custom` | `services/docker-compose.custom.yml:6` |
| Host-Port | `8099:8080`, `container_name: booking-custom` | `services/docker-compose.custom.yml:3, 7-8` |
| `.env.example` | `CUSTOM_BOOKING_PATH=./booking/custom`, dazu `COMPOSE_PATH_SEPARATOR`/`COMPOSE_FILE` (alle vier Dateien) | `services/.env.example:1-5` |
| Make-Targets | Alle setzen `-f` selbst und laden `custom.yml` immer mit. `docker-up-hub` zieht ohne `custom.yml`, baut `booking-custom` danach lokal. | `services/Makefile:55-67` |
| Dashboard-Zugriff | fest `http://booking-custom:8080` im Compose-Netz | `services/docker-compose.infra.yml:64`, `services/traefik/dynamic.yml:248-251` |
| Live-Status Schritt 0 | prüft nur `/health`, erkennt die Beispiel-Lösung nicht | `services/dashboard/main.go:102-119`, `index.html:2736-2757` |
| Hausaufgaben-Filter | Wirft ein Image weg, sobald *ein* Tag `workshop-microservices-` enthält oder in einer `image:`-Zeile einer Compose-Datei steht. Liest dafür auch die live gemountete `custom.yml`. | `homework.go:514-528`, `homework.go:564-584`, `docker-compose.infra.yml:63, 70` |
| Compose-Container | Container mit Label `com.docker.compose.project` zählen nie als Hausaufgabe | `homework.go:285-290` |
| Hub-Image `:custom` | wird von der CI gebaut und gepusht, aber von keinem Make-Target gezogen | `.github/workflows/docker.yml:214-232` |
| Compose-Gate | setzt `CUSTOM_BOOKING_PATH` selbst | `scripts/gates.sh:95-96` |

### Herkunft

- Erster Entwurf (`d2412e2`, `docker-compose.workshop.yml`): kein Host-Port,
  eigener Target `docker-up-workshop`, also nicht im Standard-Stack. Die
  Begründung für die `.env` lautete damals: „Das Repo bleibt unverändert.“
- Port 8099 und das Mitstarten im Standard-Stack kamen mit dem Umbau auf
  `custom` (`98aac8b`, `4c5ae7c`), ohne dokumentierte Begründung.
- Das Präfix `larmic/` kam mit dem Hub-Publishing der Kotlin-Beispiel-Lösung
  (`923aac3`).
- `COMPOSE_FILE` in der `.env` wurde am 13.11. nachgerüstet (`e4fa808`,
  Finding T1), damit die gedruckten `docker compose logs`-Befehle aus Story 8
  laufen.
- Entscheidung D5: Das Einklinken passiert „bewusst erst gemeinsam im
  Workshop“ (`docs/feedback/2026-05-30/backlog.md:194-198`). Offenes Finding
  im Abschlussreview: „Custom-Service nie vorab im Stack getestet“
  (`docs/review/2026-11-05-abschlussreview.md:280`).

## Bewertung der Trainer-Hypothesen

### 1. „Die .env-Anpassung ist vielleicht gar nicht nötig.“ Teilweise bestätigt

- **Für `make` überflüssig ist `COMPOSE_FILE`:** Alle Targets setzen `-f`
  selbst. Die Variable ist nur Komfort für nackte `docker compose`-Befehle. Die
  braucht aber Story 8 (`docs/stories/story-08-tracing.md:70-74`), und
  `docs/troubleshooting.md:22-31` verweist darauf.
- **Pflicht ist heute `CUSTOM_BOOKING_PATH`:** Ohne Variable warnt Compose nur
  und baut aus `services/` mit dem generischen Go-Dockerfile. Weil dort kein
  `main`-Paket liegt, scheitert der Build mit einer irreführenden
  Go-Fehlermeldung (`services/Dockerfile:8, 13`).
- **Kopien veralten:** Die lokale `.env` des Trainers enthält nur
  `CUSTOM_BOOKING_PATH`, kein `COMPOSE_FILE`. Nacktes `docker compose config`
  bricht dort mit `service "hotel" depends on undefined service "consul"` ab.
  Neue Zeilen aus `.env.example` erreichen bestehende `.env`-Dateien nie.

Fazit: Für den *Pfad* braucht man die `.env` nicht, wenn er in `custom.yml`
steht. Eine `.env` mit `COMPOSE_FILE` bleibt aber nützlich, solange Story 8
nackte Compose-Befehle druckt.

### 2. „custom.yml startet den Service mit, der Port 8099 blockiert den eigenen Start.“ Bestätigt

- Jeder Stack-Start (`docker-up`, `docker-up-hub`, `docker-rebuild-custom`,
  nacktes `docker compose up` über `COMPOSE_FILE`) startet `booking-custom`
  mit Host-Port 8099. Bis zum Einklinken läuft dort die Beispiel-Lösung.
- Die Vorbereitung trainiert `docker run -p 8099:8080` und verspricht, der
  Service sei „auch später im Workshop“ unter 8099 erreichbar
  (`docs/vorbereitung.md:47, 53`). Das lädt zum Konflikt ein.
- Der Konflikt besteht in beiden Richtungen. Läuft der Smoke-Test-Container
  noch, bricht der Stack-Start ab (`docs/troubleshooting.md:140-142`). Diese
  Richtung ist dokumentiert. Startet jemand nach dem Stack seinen Container
  auf 8099, scheitert dieser. Das steht nirgends.
- Wer aus der IDE auf Port 8080 startet (Standard bei Spring Boot oder Ktor),
  kollidiert zusätzlich mit Traefik (`services/docker-compose.infra.yml:4-6`).
- **Wichtiger als der Port:** Auch ohne Konflikt wäre ein Service aus der IDE
  oder per `docker run` im Bridge-Netz für Dashboard und Gateway unsichtbar.
  Beide sprechen `http://booking-custom:8080` im Compose-Netz an, und der
  Service selbst braucht `flight:8080` usw. Port 8099 ist nur ein bequemer
  Direktzugang und hat im Ablauf keine Funktion.

### 3. „Ohne .env, Ansage ‚passt custom.yml an‘. Dann muss auch kein Image zu finden sein.“ Teilweise bestätigt

- **Image: richtig.** Kein Make-Target zieht `larmic/…:custom`. Compose baut
  immer aus `build.context`. Ohne `image:`-Zeile heißt das Image
  `services-booking-custom`, `up` baut es, `pull` überspringt es („No image to
  be pulled“, per `--dry-run` geprüft).
- **Heute gibt es sogar einen Fallstrick durch die Zeile:** Ein nacktes
  `docker compose pull` (mit `COMPOSE_FILE` aus der `.env`) zieht die
  Kotlin-Beispiel-Lösung vom Hub unter genau diesen Tag. Ein `docker compose up`
  ohne `--build` zieht sie, wenn lokal kein Image existiert. In beiden Fällen
  läuft das Beispiel, obwohl der Pfad stimmt.
- **custom.yml statt .env: tragfähig, mit Kosten.**
  - Die Datei ist versioniert, die Arbeitskopie wird dirty. `git pull`
    scheitert aber nur, wenn ein Upstream-Commit genau diese Datei ändert.
    Zuletzt geändert wurde sie am 22.05. (`26a7338`). Ausgerechnet der Umbau
    selbst ändert sie, wer vorher schon lokal editiert hat, bekommt also einmal
    einen Konflikt. Bei `pull.rebase=true` ohne Autostash bricht `git pull`
    bei jeder lokalen Änderung ab.
  - Bestehende `.env`-Dateien mit `CUSTOM_BOOKING_PATH` würden bei festem Pfad
    still ignoriert, und dann läuft wieder das Beispiel. Eine Brücke
    `${CUSTOM_BOOKING_PATH:-./booking/custom}` verhindert das.
  - Wer in custom.yml `image:` umbenennt statt löscht, versteckt sein Image
    erneut (siehe Hypothese 5). Die Ansage muss also lauten: „Zeile
    **löschen**“, oder die Zeile ist schon im Repo weg.
  - Im Dashboard-Container wird custom.yml beim Scaling mit ausgewertet
    (`services/dashboard/handler/scaling.go:108-111`). Ein dort nicht
    existierender Kontext stört `up --scale flight` nicht (geprüft per
    `--dry-run`). Erst ein `up` auf `booking-custom` aus dem Dashboard heraus
    würde scheitern, und das gibt es heute nicht.

### 4. „Das Einklinken könnte Teil der Hausaufgabe werden.“ Teilweise bestätigt

- **Dafür:** Docker-, Pfad- und Windows-Probleme tauchen zu Hause auf statt
  bei 15 Leuten gleichzeitig im Raum. Das schließt das offene Finding aus dem
  Abschlussreview (Schritt 2b). Der Smoke-Test per `docker run -p 8099` und
  damit der Portkonflikt entfallen.
- **Dagegen:**
  - Entscheidung D5 („nicht vorab, kein Spoiler“) müsste bewusst aufgehoben
    werden. Die Anleitung verrät allerdings nur Mechanik, keine Story-Inhalte.
  - Team-Setup (T16): Im Workshop biegen drei von vier Leuten ohnehin auf
    einen fremden Service um. Zu Hause wird also etwas geprobt, das im Raum
    anders abläuft. Der Handgriff ist dann aber bekannt.
  - Der heutige Hausaufgaben-Check ignoriert Compose-Container
    (`homework.go:285-290`) und sieht deshalb genau den eingeklinkten Service
    nicht. Ohne Umbau des Checks sieht der Teilnehmer nach dem Einklinken
    rot oder eine leere Liste.
  - Das Erkennen der Beispiel-Lösung über das Vendor-Label „Team Neusta“ ist
    zu grob. Dasselbe Label tragen `services/Dockerfile:22` und
    `services/dashboard/Dockerfile:27`, und wer eins davon als Vorlage kopiert,
    gilt als Beispiel. Nötig wäre ein eigenes Label (z. B.
    `workshop.example=true`) nur in `services/booking/custom/Dockerfile`.

### 5. „Ohne Auskommentieren der image:-Zeile taucht das neue Image im Hausaufgaben-Test nicht auf.“ Bestätigt

Ursache im Code: `filterHomeworkImages` verwirft ein Image, sobald *einer*
seiner Tags `workshop-microservices-` enthält oder in einer `image:`-Zeile der
gemounteten Compose-Dateien steht (`homework.go:514-528, 564-584`). Der
Container `booking-custom` zählt wegen seines Compose-Labels ohnehin nicht.

Warum „bei einigen“, je nach Weg:

| Wer | Was passiert | Sichtbar? |
|---|---|---|
| baut nach dem Einklinken nur noch per `make` | Das neue Image trägt nur `larmic/workshop-microservices-booking:custom`. | nein; nur das alte `my-booking-service` (veralteter Stand) |
| hat `image:` auf den eigenen Namen umgeschrieben, z. B. `my-booking-service` | `composeImageRepos` liest die Zeile aus der live gemounteten custom.yml und behandelt das Repo als Infrastruktur. Mit `build:` überschreibt Compose zusätzlich den Hausaufgaben-Tag. | nein |
| hat erst mit, dann ohne `image:`-Zeile gebaut | Bei vollem Cache ist die Image-ID gleich, das Image trägt beide Tags, der Filter greift pro Image. | nein, bis `docker rmi larmic/workshop-microservices-booking:custom` |
| hat zusätzlich `docker build -t my-booking-service .` ausgeführt | Compose setzt das Label `com.docker.compose.image.builder`, die Image-ID unterscheidet sich (Compose 5.7.1 geprüft). | ja |
| hat die `image:`-Zeile von Anfang an gelöscht | Das Image heißt `services-booking-custom`, kein Filter greift. | ja, in der Auswahlliste (nicht automatisch) |

Bei älteren Compose-Versionen ohne Builder-Label könnte auch ein reiner
`docker build` dieselbe ID bekommen und mit versteckt werden. Das ist nicht
verifiziert.

## Varianten im Vergleich

| | **A: Minimal** | **T: Eine Zeile, ein Ort** | **R: Robust per make** |
|---|---|---|---|
| Kern | `image:`-Zeile raus, `.env` bleibt, Ansage „nur im Stack“ | Pfad direkt in custom.yml (relativ, Repos nebeneinander), keine `.env`, Einklinken als Hausaufgabe, Check am Container | `make einklinken PFAD=…` schreibt die `.env`, eigener Image-Name `custom-local` + `pull_policy: build`, Ampel in Schritt 0, `make diagnose` |
| Hypothese 1/3 (keine .env) | nein | ja (bis auf `COMPOSE_FILE`) | nein, `.env` wird aber nie von Hand bearbeitet |
| Hypothese 2 (Port) | nur Ansage | Smoke-Test entfällt | Smoke-Test auf 9090 |
| Hypothese 5 (Filter) | behoben | behoben | durch Filter pro Tag und weil niemand custom.yml anfasst |
| Im laufenden Workshop ausrollbar | ja, custom.yml ist live gemountet | nur Teile | nur Teile; Dashboard braucht CI-Build + Pull |
| Git-Konflikte | keine | selten (nur bei Änderung von custom.yml im Upstream) | keine |
| Aufwand | sehr klein | groß (Dashboard, Doku, `compose.yaml`) | am größten (Make-Shell-Logik, Handler, Doku) |
| Schwächen | Port bleibt Ansage; `services-booking-custom` kann ein fremder Team-Service sein und wird im Picker vorausgewählt | Fester Pfad bricht bestehende `.env` still; Nebeneinander-Regel kollidiert unter Windows mit „Repo ins WSL-Dateisystem“; `compose.yaml` neben `docker-compose.yml` erzeugt bei jedem Befehl die Warnung „Found multiple config files“ | Die „Ampel für den Trainer“ ist ein Trugschluss, weil jedes Dashboard lokal läuft; Backslash-Umsetzung in make wirkungslos (die Shell frisst sie); `pull_policy: build` baut bei jedem `up` neu |
| Bewertung der Kritik | 6/10 | 6/10 (zwei schwere Mängel, beide behebbar) | 5/10 |

## Empfehlung

### Sofortmaßnahme für den laufenden Workshop

Ziel: ohne neues Dashboard-Image, ohne CI-Abhängigkeit, mit möglichst wenigen
Handgriffen pro Teilnehmer.

1. **`services/docker-compose.custom.yml`: `image:`-Zeile ersatzlos löschen.**
   Kein `pull_policy: build`, denn ohne `image:` überspringt `pull` den
   Service ohnehin, und `pull_policy: build` würde bei jedem `up` neu bauen.
   Das Dashboard liest die Datei über den Bind-Mount. Der Filter-Effekt
   verschwindet nach einem Stack-Neustart.
2. **Kein Default-Pfad als Sofortmaßnahme.** Ein
   `${CUSTOM_BOOKING_PATH:-./booking/custom}` würde eine fehlende oder falsch
   abgelegte `.env` verdecken: Statt eines lauten Build-Fehlers läuft dann
   still das Beispiel mit grünem `/health`.
3. **Ansage im Raum (und eine Zeile auf Folie 08c):** „Euer Service läuft
   *nur* als Container `booking-custom` im Stack. Nicht per `docker run`, nicht
   aus der IDE, sonst sehen ihn Dashboard und Gateway nicht. 8099 gehört dem
   Stack.“ Bei `port is already allocated`: `docker ps`, fremden Container
   stoppen.
4. **`docs/vorbereitung.md:47-55` sofort korrigieren:** Smoke-Test auf
   `-p 9090:8080` (wie der Dashboard-Fallback, `index.html:2663`), den Satz
   „auch später im Workshop unter 8099 erreichbar“ streichen. Die
   Vorbereitungsseite rendert die Markdown-Datei im Browser
   (`.github/workflows/deploy-pages.yml:54-58`), `pages/vorbereitung/` muss
   nicht geändert werden.
5. **Handgriffe für Teilnehmer**, die bereits eingeklinkt haben (sonst nichts
   tun):
   - wer custom.yml lokal geändert hat: `git checkout -- services/docker-compose.custom.yml`
   - `git pull`
   - Ctrl+C, `make docker-up-hub` (baut jetzt `services-booking-custom`)
   - nur wer im Hausaufgaben-Check weiterhin nichts sieht:
     `docker rmi larmic/workshop-microservices-booking:custom`. Das kann den
     nächsten Build verlängern, wenn die Layer nur an diesem Tag hängen.

   `make docker-down` ist nicht nötig. Falls nach dem Pull die alte custom.yml
   im Dashboard-Container hängt (Single-File-Mount, alter Inode), hilft es,
   denn ein Neustart des Containers lädt die Datei neu. Achtung: `docker-down`
   nutzt `down -v`.

Damit sind Hypothese 5 und der Hub-Fallstrick behoben, Hypothese 2 ist
entschärft. Die `.env` bleibt vorerst, wie sie ist.

### Umbau danach (vor dem nächsten Durchlauf)

Richtung: Variante T in entschärfter Form, also der Vorschlag des Trainers.

1. **Pfad direkt in `docker-compose.custom.yml`.** Die Zeile wird
   `context: ${CUSTOM_BOOKING_PATH:-./booking/custom}` mit einem deutschen
   Kommentar „EINKLINKEN: hier euren Pfad eintragen“. Teilnehmer ersetzen die
   ganze Zeile durch ihren Pfad. Die Variable bleibt nur als Brücke für alte
   `.env`-Dateien und das Gate. Der Trainer friert die Datei während einer
   Workshop-Reihe ein.
2. **`.env` nur noch für `COMPOSE_FILE`.** `.env.example` ohne
   `CUSTOM_BOOKING_PATH`; `cp .env.example .env` bleibt in der Vorbereitung,
   aber niemand editiert die Datei. Keine `compose.yaml` mit `include:`, wegen
   der Warnung „Found multiple config files“.
3. **Pfad-Konvention mit Rücksicht auf Windows.** Ein relativer Pfad wird
   empfohlen (Repos nebeneinander, fester Ordnername wie `team-service`),
   aber nicht vorgeschrieben. Wer unter Windows den Service unter `C:\` hat,
   trägt `/mnt/c/...` ein.
4. **Schritt 0 bekommt einen Container-Check im Dashboard.** Ist
   `booking-custom` laufend? Trägt sein Image das neue Label
   `workshop.example=true` (nur in `services/booking/custom/Dockerfile`)?
   Daraus ergibt sich ein Status in drei Stufen: nicht erreichbar,
   Beispiel-Lösung, euer Service. Der Log-Kasten bleibt als Zusatz.
5. **Hausaufgaben-Filter pro Tag statt pro Image**
   (`homework.go:514-528`), dazu ein Testfall für ein Image mit Workshop- und
   Teilnehmer-Tag. Compose-gebaute Images (`services-booking-custom`) im
   Picker nicht vorauswählen, weil das der Service eines Teamkollegen sein
   kann.
6. **CI-Job `build-booking-custom` und Hub-Tag `:custom` streichen.**
7. **Einklinken als Hausaufgabe:** nur, wenn der Trainer D5 aufheben will.
   Dann als optionale Generalprobe mit dem eigenen Service, geprüft durch den
   neuen Container-Check, und als Ersatz für den `docker run`-Smoke-Test.
8. **Port 8099 neu bewerten:** behalten als Direktzugang oder streichen und
   `curl` über das Gateway (`localhost/api/booking-custom/...`). Ohne Host-Port
   ist der Konflikt vollständig weg.

Variante R ist nicht empfohlen. `make einklinken` löst vor allem das
Schreiben einer Datei, die nach dem Umbau ohnehin niemand mehr schreibt. Die
Ampel hilft nur dem Teilnehmer selbst, nicht dem Trainer.

## Betroffene Stellen

Sofortmaßnahme:

- [ ] `services/docker-compose.custom.yml:6` (`image:`-Zeile löschen)
- [ ] `services/slides/chapters/08c-start.md` (Zeile „nur im Stack“, Hinweis `port is already allocated`)
- [ ] `docs/vorbereitung.md:47-55` (Smoke-Test 9090, Satz zu 8099 streichen)
- [ ] `docs/troubleshooting.md:75-112` (Hausaufgabe rot: alter `:custom`-Tag, `docker rmi`)
- [ ] `docs/troubleshooting.md:114-143` (Port belegt: Fall „Stack belegt 8099, eigener `docker run` scheitert“)
- [ ] `docs/custom-setup.md` (Satz „läuft ausschließlich als `booking-custom`“)

Umbau danach:

- [ ] `services/docker-compose.custom.yml:4-8` (Kontext mit Kommentar, Port-Entscheidung)
- [ ] `services/.env.example:1-5` (nur noch `COMPOSE_FILE`)
- [ ] `services/Makefile:55-67` (Hilfetexte „aus CUSTOM_BOOKING_PATH“)
- [ ] `scripts/gates.sh:95-96` (Variable nicht mehr nötig; optional prüfen, dass kein lokaler Pfad und keine `image:`-Zeile committet ist)
- [ ] `services/dashboard/handler/homework.go:285-290, 514-528, 564-584` und `homework_test.go`
- [ ] `services/dashboard/main.go:102-132` und `services/dashboard/static/index.html:1683-1721, 2604-2627, 2732-2760, 3720` (Schritt 0, Hausaufgaben-Karte, Untertitel)
- [ ] `services/booking/custom/Dockerfile:21` (Label `workshop.example=true`)
- [ ] `services/booking/custom/src/main/kotlin/workshop/booking/Application.kt:28-39` und `services/booking/custom/README.md` (Banner-Text verweist auf custom.yml statt `.env`)
- [ ] `services/slides/chapters/08c-start.md`, `08d-reiseplattform.md:41`, `09-story-01.md:44-46`
- [ ] `docs/vorbereitung.md:65-89` (`.env`-Schritt, Checkliste, LF-Hinweis)
- [ ] `docs/custom-setup.md` (komplett auf die Kontext-Zeile umschreiben)
- [ ] `docs/troubleshooting.md:22-31, 243-257` (`COMPOSE_FILE`, Windows/`.env`)
- [ ] `docs/stories/story-01-cloud-native-booking-service.md:9, 32`, `docs/stories/story-08-tracing.md:70-74`
- [ ] `docs/themen.md:91` (Moderation Schritt 0)
- [ ] `services/README.md:66, 75`
- [ ] `.claude/agents/docker-compose-validator.md:19, 29, 41, 122`
- [ ] `.github/workflows/docker.yml:214-232`, `.github/dockerhub/booking.md:33`
- [ ] `services/.gitattributes` (eol-Regel für `docker-compose.custom.yml` ergänzen, wenn Windows-Editoren sie anfassen)

## Offene Fragen und Entscheidungen für den Trainer

1. **Sofortmaßnahme heute:** `image:`-Zeile jetzt löschen und alle einmal
   `git pull` + Neustart machen lassen, oder nur die Ansage und die Ausnahme
   „wer nichts sieht: `docker rmi …:custom`“? Folgt noch ein zweiter Tag?
2. **Pfad in custom.yml statt `.env`:** Soll das der Weg ab dem nächsten
   Durchlauf sein, und wird custom.yml während einer Reihe eingefroren?
3. **Fester Ordnername** für den Team-Service (z. B. `../../team-service`),
   damit die Folie den Pfad wörtlich zeigen kann?
4. **D5 aufheben?** Einklinken als (optionale) Generalprobe in der Hausaufgabe
   oder weiterhin erst im Workshop?
5. **Port 8099:** als Direktzugang behalten oder streichen und nur über das
   Gateway testen?
6. **Erkennung der Beispiel-Lösung:** neues Label `workshop.example=true`
   einführen und den Container-Check in Schritt 0 bauen?
7. **Hub-Image `:custom`:** CI-Job streichen oder als Showcase behalten?
8. **Story 8:** `COMPOSE_FILE` weiter über die `.env`, oder die gedruckten
   Befehle auf `make`-Targets bzw. `-f`-Flags umstellen, damit die `.env`
   ganz entfallen kann?
