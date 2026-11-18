# Konzept: Einklinken in „Start (0)“ vereinfachen

Stand: 18.11.2026, während eines laufenden Workshops.

**Worum es geht:** Der Schritt „Start (0) – eigenen Service einklinken“ läuft
im Workshop holprig. Dieses Dokument wertet die Beobachtungen des Trainers aus,
vergleicht drei Lösungswege und empfiehlt eine Sofortmaßnahme für den laufenden
Workshop sowie einen Umbau danach.

**Ergänzt im Lauf des Tages:** Nachträge aus Fragen im Workshop (Consul,
Env-Variablen, Story 3) und die konsolidierte Feedback-Liste des zweiten
Durchlaufs (Abschnitt „Feedback des zweiten Durchlaufs, konsolidiert“), die
auch das Story-4-Review vom 17.11. einbezieht.

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

## Nachtrag: Consul-Adresse in Story 3 (Browser vs. Container)

**Beobachtung im Workshop:** Im Browser liefert
`http://localhost:8500/v1/health/service/flight-service?passing=true` die
gesunden Flight-Instanzen. Im eigenen Service funktioniert dieselbe URL nicht,
dort muss es `http://consul:8500/...` heißen. Das steht nirgends so erklärt,
dass man es vor dem ersten Fehlschlag findet.

**Stimmt das?** Ja.

- `services/docker-compose.infra.yml:17-18` veröffentlicht Consul mit
  `8500:8500` auf dem Host. Deshalb klappt `localhost:8500` im Browser und per
  `curl` vom Laptop.
- Im Container zeigt `localhost` auf den Container selbst. Consul ist dort nur
  über den Compose-Servicenamen `consul` erreichbar.
  `docker-compose.custom.yml` setzt dafür `CONSUL_URL=http://consul:8500`.
- Dasselbe gilt für die Antwort: Die Referenz-Services registrieren als
  `Address` ihren Container-Hostnamen (`shared/consul/register.go:37-47`, eine
  Container-ID). Die URL aus `Address` + `Port` funktioniert also nur im
  Docker-Netz, im Browser nicht. Wer im Browser die JSON-Antwort ansieht und
  die Adresse ausprobiert, läuft in denselben Fehler, nur umgekehrt.
- Wer seinen Service aus der IDE statt im Stack startet (siehe Hypothese 2),
  bräuchte `localhost:8500`. Die Adressen aus Consul kann er trotzdem nicht
  auflösen. Das ist ein weiterer Grund für „nur als `booking-custom` im Stack“.

**Was heute dazu im Material steht:**

- `docs/stories/story-03-service-discovery.md:29`: „im Stack
  `http://consul:8500`“, ohne den Unterschied zu `localhost`
- `docs/custom-setup.md:47-53` (Tabelle interne URLs) und `:61` (Consul UI
  über `localhost/consul`): beide Sichten, aber getrennt und ohne Hinweis
- Folie 11a `11a-consul.md:13-21`: `consul:8500` ohne Erklärung
- Dashboard Story 3 (`index.html:1871`): „im Stack http://consul:8500“
- `docs/instructions/service-discovery-consul.md:19`: Consul-UI unter
  `http://localhost:8500`

**Vorschlag:** Den Unterschied an einer Stelle als Merksatz erklären und von
den anderen dorthin verlinken, zum Beispiel:

> Vom Laptop (Browser, `curl`): `localhost:8500`. Aus eurem Container: `consul:8500`,
> also `CONSUL_URL`. Im Container ist `localhost` euer eigener Container. Die
> Adressen in der Consul-Antwort sind Container-Namen und gelten nur im
> Docker-Netz.

Passende Orte: die Tabelle „Interne URLs“ in `docs/custom-setup.md` (zweite
Spalte „vom Laptop“) und Aufgabe 1 in Story 3 inklusive Dashboard-Text. Auf
Folie 11a reicht eine Zeile in den Notes. Dasselbe Muster gilt für Flight,
Hotel und Car (`flight:8080` im Container, `localhost/api/flight` vom
Laptop). Die Tabelle würde also alle Backends abdecken.

**Aufwand:** klein.

## Nachtrag: Env-Variablen in `docker-compose.custom.yml` anpassen

**Feedback aus dem Workshop:** Man kann ja auch die Env-Variablen in
`docker-compose.custom.yml` anpassen. Ist das eine gute Idee?

**Geht das?** Ja. Der Block `environment:` (`docker-compose.custom.yml:9-13`)
gilt nur für `booking-custom`. Nach einer Änderung reicht
`make docker-rebuild-custom`, denn `up -d` erzeugt den Container neu, sobald
sich seine Konfiguration geändert hat. Der Rest des Stacks läuft weiter. Das
Dashboard liest custom.yml nur wegen `image:` (Hausaufgaben-Filter) und beim
Scaling, nicht wegen `environment:`.

**Bewertung: für eigene Variablen gut, die vier vorgegebenen nicht anfassen.**

- **Dafür:**
  - Das passt zu Faktor III der 12-Faktor-App (Konfiguration in der Umgebung)
    und zu Story 1.
  - Teams können eigene Werte von außen setzen, ohne neu zu bauen, etwa
    Timeout, Breaker-Schwelle, Bulkhead-Größe oder Log-Level. Damit lässt sich
    im Spiel Gesehenes live nachstellen, zum Beispiel ein kürzerer Timeout in
    Story 4.
- **Dagegen bzw. Grenzen:**
  - `FLIGHT_SERVICE_URL`, `HOTEL_SERVICE_URL`, `CAR_SERVICE_URL` und
    `CONSUL_URL` sind der Vertrag mit dem Stack. Werden sie auf `localhost`
    umgebogen, bricht der Service (siehe Nachtrag Consul). Ab Story 3 bleiben
    die drei Backend-URLs einfach gesetzt, der Service nutzt sie nur nicht
    mehr. Löschen ist nicht nötig.
  - Die Datei ist versioniert. Lokale Änderungen machen die Arbeitskopie dirty
    und können beim `git pull` kollidieren, genau wie beim Pfad (siehe
    Hypothese 3). Wer nach jeder Story `git pull` macht, stolpert hier
    zuerst.
  - Im Team hat nur der Laptop, auf dem getippt wird, die Änderung. Beim
    Wechsel pro Story fehlt sie auf dem nächsten Rechner, wenn sie nicht im
    Team-Repo steht.
  - Secrets oder Webhooks gehören nicht in eine versionierte Datei.
- **Hinweis zur Referenz:** Die Referenz-Stories fallen ohne `CONSUL_URL`
  auf `http://localhost:8500` zurück (`services/booking/story4/main.go:23`,
  ebenso Story 8). Im Container ist dieser Default falsch, wer die Variable
  löscht, merkt es also erst zur Laufzeit.

**Vorschlag:**

1. **Ansage im Workshop:** „Eigene Variablen dürft ihr in
   `docker-compose.custom.yml` ergänzen. Die vier vorgegebenen bleiben, wie sie
   sind.“ Gern als Zeile in `docs/custom-setup.md` unter „Interne URLs“.
2. **Wenn der Pfad nach dem Umbau ohnehin in custom.yml steht,** gehört die
   Datei dem Team. Env-Anpassungen passen dann ohne Zusatzaufwand dazu, die
   Git-Kosten sind dieselben.
3. **Alternative ohne dirty Arbeitskopie:** In custom.yml
   `env_file: { path: booking-custom.env, required: false }` eintragen und
   `services/booking-custom.env` in `.gitignore` aufnehmen. Teams schreiben
   eigene Werte dorthin, die versionierte Datei bleibt unberührt. Kosten: eine
   weitere Datei im mentalen Modell. Die Syntax braucht außerdem Compose
   2.24 oder neuer, das wäre in `docs/vorbereitung.md` und mit
   `scripts/gates.sh compose` zu prüfen. Eignet sich eher für den Umbau, nicht
   für heute.

**Aufwand:** klein (Ansage und Doku-Zeile), mittel mit `env_file`-Variante.

## Nachtrag: Muss der POST in Story 3 prüfen, ob Flug, Hotel und Auto existieren?

**Frage aus dem Workshop:** Ist es für die folgenden Stories wichtig, dass
`POST /booking/bookings` vor dem Buchen prüft, ob Flug, Hotel und Auto
vorhanden sind?

**Antwort: Nein, und eine solche Prüfung ist sogar kontraproduktiv.**

- **Die Backends prüfen selbst kaum:** `POST /bookings` bei Flight nimmt jede
  `flightId` an und antwortet mit `201` (`services/flight/handler/flight.go:45-70`),
  Hotel genauso. Einzige fachliche Ablehnung ist `carId` `C9` mit `409`
  (`services/car/handler`, `unavailableCarID`), gedacht für Story 6.
- **Story 3** verlangt ausdrücklich „Kein Error-Handling (darf fehlschlagen)“.
- **Story 4:** Der Pre-Check im POST prüft nur, ob ein Breaker OPEN ist
  (dann sofort `503`), nicht die Existenz.
- **Story 5** arbeitet mit `GET /booking/offers`, nicht mit dem POST.
- **Story 6** ersetzt den POST aus Story 3 ganz. Fehler kommen per Chaos-Modus
  oder `C9`, Booking reagiert auf den Status des Backends und kompensiert.
- **Story 7 und 8:** keine Rolle.

**Warum eine Vorab-Prüfung den Lernzielen schadet:**

- Ob gebucht werden kann, entscheidet der Service, dem die Daten gehören.
  Prüft Booking vorher, wird diese Logik verdoppelt.
- Prüfen und Buchen sind zwei Calls, dazwischen kann das Auto weg sein. Dieses
  Check-then-act-Problem löst nur „buchen und auf `409` reagieren“, am Ende
  mit Kompensation (Story 6).
- Jeder zusätzliche Call kann ausfallen, mehr Fläche für Story 4 und 5.

**Sinnvoll ist dagegen:** den Request-Body validieren (fehlende Pflichtfelder
`400`) und Backend-Fehler als Status-Code durchreichen, nicht als `200` mit
Fehler-Body (Story 2). Das zahlt sich ab Story 6 aus.

**Umgesetzt (18.11.):** Hinweis in `docs/stories/story-03-service-discovery.md`
(Aufgabe 3) und im Dashboard unter Story 3, „Buchung orchestrieren“.
Folie 12 „Story 3“ ist unverändert, dort genügt die Ansage.

## Nachtrag: Kann Consul direkt eine zufällige Instanz liefern?

**Frage aus dem Workshop:** Gibt es eine Consul-URL, die statt einer Liste
gleich eine zufällige Instanz zurückgibt? Heute muss jeder selbst einen
Randomizer einbauen.

**Befund** (geprüft mit `hashicorp/consul:2.1` im Dev-Mode und drei
registrierten `flight-service`-Instanzen in einem Wegwerf-Container):

- **Health-API:** `/v1/health/service/{name}?passing=true` liefert immer die
  Liste und kennt kein `limit`.
- **Prepared Query:** `GET /v1/query/flight-service/execute?limit=1` liefert
  pro Aufruf zufällig eine gesunde Instanz (sechs Aufrufe: f2, f2, f3, f1,
  f2, f1). Die Antwort ist trotzdem ein Array (`Nodes`) mit einem Element.
  Ohne Vorbereitung antwortet Consul `404 Query not found`. Die Query muss
  einmal angelegt werden, mit einem Template für alle Namen:
  `POST /v1/query {"Name":"","Template":{"Type":"name_prefix_match"},"Service":{"Service":"${name.full}","OnlyPassing":true}}`.
- **DNS:** `flight-service.service.consul` auf Port 8600 liefert nur gesunde
  Instanzen in zufälliger Reihenfolge. Im Docker-Netz antwortet aber Dockers
  eingebauter DNS, die Teilnehmer-Container lösen `.consul` ohne zusätzliche
  DNS-Konfiguration nicht auf. Clients mit DNS-Cache (z. B. JVM) bleiben
  außerdem an einer Instanz hängen.

**Ist die Prepared Query besser als der heutige Weg? Empfehlung: nein, als
Bonus ja.**

- **Der Randomizer ist nicht die Hürde.** Er ist eine Zeile
  (`list.random()`, `rand.Intn(len(list))`). Die eigentliche Arbeit bleibt in
  beiden Wegen gleich: Consul abfragen, Antwort parsen, URL aus `Address` +
  `Port` bauen, leere Antwort behandeln. Bei der Prepared Query ist das JSON
  sogar tiefer verschachtelt (`Nodes[0].Service.Address`).
- **Das Lernziel ist Client-Side Load Balancing.** Folie 11 „Service
  Discovery“ (Notes: „Der Booking-Service fragt Consul selbst und wählt eine
  Instanz“), Folie 12 „Story 3“, Story-Doc Aufgabe 1 und
  `docs/questions/story3.md` Frage 3 bauen darauf auf. Mit der Prepared Query
  wählt die Registry, nicht der Aufrufer. Das verschiebt die Topologie, über
  die Frage 3 diskutiert.
- **Mehr Infrastruktur im Stack:** Die Query muss beim Start angelegt werden
  (Init-Schritt in Compose). Consul läuft mit `tmpfs`
  (`docker-compose.infra.yml:19-20`), nach jedem Neustart ist sie weg. Startet
  ein Service vor dem Init, bekommt er `404`. Die Referenz-Implementierung
  (`shared/consul/resolver.go`) und die Folie 11a „Consul“ mit
  `GET consul:8500/v1/health/service/{name}?passing=true` müssten ebenfalls
  umgestellt werden.
- **Die Liste ist später nützlich:** Mit der Liste kann der Client selbst
  entscheiden, etwa Round-Robin, bei Fehler eine andere Instanz versuchen
  oder die Liste kurz cachen, statt bei jedem Request Consul zu fragen. Das
  sind gute Anknüpfungen an Story 4 und 5.
- **Dafür spricht:** weniger Code bei den Teilnehmern, Health-Filter und
  Zufall an einer Stelle, Consul-Feature sichtbar gemacht.

**Vorschlag:** Den heutigen Weg behalten. Die Prepared Query (und DNS) als
Bonus bzw. Diskussionsimpuls in Story 3 aufnehmen: „Consul kann auch selbst
wählen. Was verliert ihr, wenn die Registry entscheidet?“ Passt zu Frage 3 in
`docs/questions/story3.md`.

> **Antwort für den Trainer:** Ihr verliert die Entscheidung im Moment des
> Aufrufs, und Consul wird Teil jedes Requests.
>
> 1. **Consul im heißen Pfad:** Eine Liste kann Booking kurz cachen und
>    weiterarbeiten, auch wenn Consul hängt. Mit `limit=1` braucht jeder
>    Request Consul: ein Hop mehr und ein neuer Single Point of Failure (Story
>    4). Wer das `limit=1`-Ergebnis cacht, schickt für die Cache-Dauer alles an
>    eine Instanz, dann ist die Verteilung weg.
> 2. **Kein Retry auf eine andere Instanz:** Mit der Liste nimmt man bei einem
>    Fehler einfach die nächste. Mit `limit=1` muss man Consul erneut fragen
>    und bekommt vielleicht dieselbe.
> 3. **Keine eigene Strategie:** Consul würfelt nur. Round-Robin, Gewichtung
>    nach Last oder dieselbe Instanz pro Kunde gehen nicht.
> 4. **Eigenes Wissen bleibt ungenutzt:** Consul sieht nur den Health-Check
>    alle paar Sekunden. Booking merkt sofort, dass eine Instanz bei ihm langsam
>    ist oder Fehler liefert, und könnte sie aussortieren (Outlier Detection,
>    die Idee des Circuit Breakers pro Instanz).
> 5. **Kopplung an ein Consul-Feature:** Prepared Queries gibt es nur bei
>    Consul. „Liste gesunder Instanzen“ liefert jede Registry, auch
>    Kubernetes.
>
> **Gewinn:** weniger Code pro Service und eine zentrale Regel, die man ändern
> kann, ohne Clients neu auszurollen. Bei vielen Sprachen und Teams ein echtes
> Argument.
>
> **Überleitung zu Frage 3:** Ein Service Mesh nimmt euch genau diese Punkte
> ab, ohne sie zu verlieren. Der Sidecar läuft neben dem Client, hält die
> Liste, macht Retry und sortiert schlechte Instanzen aus. Die Entscheidung
> bleibt beim Aufrufer, steht aber nicht mehr in eurem Code.

## Feedback des zweiten Durchlaufs, konsolidiert

Der Workshop lief am 18.11. zum zweiten Mal, das Feedback war durchweg
positiv. Die Verbesserungspunkte (vom Trainer nach dem Workshop
zusammengefasst) sind hier mit diesem Dokument und dem Story-4-Review vom
17.11. (`2026-11-17-story4-spiel.md`) abgeglichen. Nichts davon ist
umgesetzt.

Kennzeichnung: **neu** = bisher nirgends notiert, **zusammengeführt** =
deckt sich mit einem bestehenden Punkt, **⚠** = Widerspruch, den der Trainer
entscheiden muss.

### A. Circuit-Breaker-Rollenspiel (Folie 15 „Wenn der Flug ausfällt“)

**A1 Visueller Einstieg** (neu)

- **Feedback:** Die Story steht nur als Text an der Wand, die Regeln sind
  nicht trivial. Gewünscht ist ein Bild oder eine kurze Bildfolge (etwa drei
  Schritte). Booking steht in der Mitte, drumherum die Requests mit Karten,
  dazu Flight, Hotel und Car. Pfeile zeigen, wer wann vorliest und wie
  Booking reagiert.
- **Bestand:** Eine Aufstellungsskizze gibt es schon, aber nur auf der
  Druckseite (`pages/material/index.html:380`, „Aufbau im Raum“), nicht auf
  der Folie. Sie kann die Basis sein.
- **⚠** Heute stehen die Requests in einer Schlange vor dem Booking-Tisch. Die
  Backends sind keine Personen, ihre Antwort steht auf der Karte. Sollen
  Flight, Hotel und Car nur im Bild erscheinen oder als Personen mitspielen?
  Letzteres ändert den Ablauf.
- **Dazu aus dieser Session:** Wer die rote Karte aufs Board legt, steht
  nirgends ausdrücklich (Story-Doc und Druckseite: „Rote Karten kommen aufs
  Board“, Passiv). Gemeint ist der Request, Folie 15 sagt aber „Booking …
  führt das Board“. Die Pfeile im Bild sollten das klären.

**A2 Teil 1 (ohne Breaker) spürbar machen** (neu)

- **Feedback:** Die drei Karten vor STOPP 1 zeigen nicht, warum das Verhalten
  ein Problem ist. Idee: Booking hat begrenzte Kapazität (etwa drei „Hände“),
  jeder hängende Request blockiert eine, weitere müssen warten. Der Stau wird
  sichtbar, der Breaker wirkt danach wie eine Erlösung.
- **⚠** Begrenzte Kapazität ist der Kern des Story-5-Spiels (Bulkhead, Stühle
  als Slots). In Story 4 nimmt sie Story 5 die Pointe vorweg und vermischt
  Timeout mit Erschöpfung des Thread-Pools. Heute trägt Teil 1: Hotel und Car
  sind gesund und trotzdem weg, und jeder wartet die volle Zeit.
- **Alternative ohne Überschneidung:** Die Wartezeit laut und sichtbar
  aufaddieren („15 Sekunden Kundenzeit für nichts“).
- **Zusammenhang:** Story-4-Review Punkt 2 (3 s Timeout in beiden Phasen,
  damit bei STOPP 1 nur der Breaker neu ist).

**A3 Probe / HALF_OPEN vereinfachen** (zusammengeführt)

- **Feedback:** Der Probe-Durchlauf geht aus den Karten nicht hervor.
  Gewünscht ist eine farblich abgesetzte Probe-Karte: Eine Person geht vor,
  nur wenn sie durchkommt, dürfen alle anderen wieder.
- **Bestand:** Das gibt es schon, die gelben Karten 14 und 18 und das
  Probe-Feld auf dem Board. Das Feedback zeigt, dass sie nicht als solche
  erkannt werden.
- **Vorschlag:** Die Probe-Karte deutlicher machen statt neu erfinden, etwa
  mit dem Wort PROBE groß und dem Ablauf als Satz auf der Karte. Zusammen mit
  Story-4-Review Punkt 1 („Flight: Timeout!“ bzw. „Flight antwortet!“).

**A4 Kartenlayout** (zusammengeführt mit Story-4-Review Punkte 1 und 3)

- **Feedback:** Warten und Vorlesen wirken wie eine Aktion, und der
  betroffene Downstream-Service ist nicht erkennbar. Vorschlag: zwei Felder
  „1. Warten“ (Uhr, Sekunden groß) und „2. Vorlesen“ (Sprechblase), dazu
  festes Icon und feste Farbe pro Service, Flight etwa immer rot. Ziel: Wer
  die Karte zum ersten Mal sieht, weiß sofort, was zu tun ist.
- **Neu:** die Teilung in „Warten“ und „Vorlesen“.
- **Schon notiert (17.11.):** Zeile „was ist mit Flight“ zuerst, Leiste mit
  den Status-Codes von Flight, Hotel und Car.
- **⚠** „Flight immer rot“ kollidiert doppelt. In der Palette ist Flight blau
  (`#1456b3`, `services/CLAUDE.md`, Gate `palette`), und Rot heißt auf den
  Karten schon „Fehler, aufs Board“. Vorschlag: Flight-Icon in Blau, den
  Ausfall per Symbol zeigen statt per Farbe.
- **Offen (17.11.):** Passt das alles auf A7?
- **Folge:** Alle vier Karton-Bogen neu drucken.

### B. Hausaufgabe und Einklinken

**B1 Einklinken in die Hausaufgabe** (zusammengeführt mit Hypothese 4 und
offener Frage 4)

- **Feedback:** Das Grundgerüst als Hausaufgabe hat gut funktioniert. Die
  Integration ins Dashboard hat im Workshop fast 45 Minuten gekostet. Sie
  gehört in die Hausaufgabe, mit klarem Erfolgskriterium („eigener Service
  im Dashboard grün“).
- **Neu:** die Zahl, rund 45 Minuten.
- **Voraussetzung:** Die Entscheidung D5 („Einklinken erst im Workshop“) fällt.
- **⚠** Heute einigen sich die Teams erst im Workshop auf einen Service (Notes
  zu Folie 08c). Wer zu Hause seinen eigenen einklinkt, stellt im Workshop
  trotzdem einmal auf den Team-Service um. Lösung: Zu Hause wird der Ablauf
  geübt, im Workshop nur der Pfad getauscht.
- **Erfolgskriterium:** Dafür passt der vorgeschlagene Container-Check in
  Schritt 0 mit drei Stufen (nicht erreichbar, Beispiel-Lösung, euer
  Service), siehe „Umbau danach“.

**B2 Doku korrigieren und einmal komplett durchspielen** (teils
zusammengeführt)

- **Feedback:** Die Integrationsdoku ist nicht gut genug und enthält Fehler.
  Überarbeiten und einmal mit frischem Checkout exakt nach Doku durchspielen.
- **Schon notiert:** `docs/vorbereitung.md:53` verspricht 8099, der
  Smoke-Test nutzt 8099 (siehe Hypothese 2 und Sofortmaßnahme).
- **Neu:** der vollständige Durchlauf mit frischem Checkout. Unter Linux mit
  Docker kann ihn Claude Code im Container machen. Windows und macOS braucht
  einen eigenen Durchlauf.

**B3 Portkonflikt mit dem Beispiel-Service** (zusammengeführt mit Hypothese 2
und offener Frage 5)

- **Feedback:** Die fest integrierte Beispiel-Lösung belegt den Port, den die
  Gruppen nutzen sollen. Ideen: eigener reservierter Port für das Beispiel,
  ein Compose-Profil, das standardmäßig aus ist, und/oder feste Ports pro
  Gruppe.
- **Fehlt im Feedback:** Die Hauptursache ist die `image:`-Zeile (Hypothese
  5), nicht der Port.
- **Eigener Port fürs Beispiel:** deckt sich mit „Port 8099 behalten oder
  streichen“.
- **Compose-Profil, standardmäßig aus:** neu und gut. Ohne eingeklinkten
  Service ist „Mein Service“ grau, statt dass still das Beispiel läuft. Der
  BEISPIEL-LÖSUNG-Kasten als Prüfmerkmal würde überflüssig. **⚠** Das
  widerspricht der Rückfallebene „Beispiel-Lösung läuft immer mit“ in diesem
  Konzept. Es braucht eine Entscheidung.
- **⚠ Feste Ports pro Gruppe:** lösen das Problem nicht. Jede Gruppe hat
  ihren Stack auf dem eigenen Laptop, Ports kollidieren nur innerhalb eines
  Laptops. Das eigentliche Problem ist der Service außerhalb des Stacks
  (`docker run`, IDE), den das Dashboard nie erreicht.

**B4 Bestehende Punkte, die im Feedback fehlen** (alle oben in diesem
Dokument)

- `image:`-Zeile löschen (Sofortmaßnahme), Hub-Image `:custom` samt CI-Job
  streichen
- Pfad in `custom.yml` statt `.env`, `COMPOSE_FILE` für Story 8 (offene
  Fragen 2 und 8)
- Consul: `localhost:8500` vom Laptop, `consul:8500` aus dem Container
  (offene Frage 9)
- Eigene Env-Variablen in `custom.yml` (offene Frage 10)
- Story 3: keine Existenzprüfung im POST (erledigt, `073fcb6`), Prepared
  Query als Bonus (offene Frage 11)

### C. Easter-Egg-Webhook (neu, geprüft)

- **Feedback:** Der Discord-Webhook liegt als Secret in GitHub und wird beim
  CI-Build eingebaut. Lokal fehlt er. Prüfen: Wofür ist er da, was passiert
  lokal ohne ihn, funktioniert dort etwas falsch? Er soll den Teilnehmern
  nicht offengelegt werden.
- **Wofür:** Wer im Dashboard den Konami-Code eingibt, schaltet das Theme
  „Inselwelt“ frei. Einmal pro Browser geht eine anonyme Nachricht an
  Discord (`services/dashboard/static/index.html:3814-3845`). Ein zweiter
  Webhook (`DISCORD_WEBHOOK_URL`) gehört zum Feedback-Formular auf GitHub
  Pages (`pages/feedback/index.html:286`, eingesetzt in
  `.github/workflows/deploy-pages.yml:67-85`).
- **Lokal ohne Secret:** Der Platzhalter `__DISCORD_EASTEREGG_WEBHOOK_URL__`
  bleibt stehen, die Funktion bricht vorher ab
  (`index.html:3832`). Kein Fehler, keine Konsolenmeldung, sonst nichts
  betroffen. Das gilt für `make docker-up` (Compose übergibt kein Build-Arg)
  und `make docker-build-dashboard` ohne Variable. Es funktioniert also
  nichts falsch.
- **⚠ Offenlegung lässt sich so nicht verhindern.** `docker.yml:211` schreibt
  die URL ins Docker-Hub-Image des Dashboards, das alle Teilnehmer per
  `docker-up-hub` ziehen. Dort steht sie im Klartext im ausgelieferten
  JavaScript, sichtbar im Seitenquelltext von `localhost`. Gleiches gilt für
  den Feedback-Webhook auf der öffentlichen Pages-Seite. Wer die URL hat,
  kann beliebige Nachrichten in den Kanal schicken oder den Webhook
  entfernen. Eine URL in Code, der im Browser der Teilnehmer läuft, ist
  grundsätzlich nicht geheim zu halten.
- **Optionen:**
  1. Risiko akzeptieren und die Webhooks nach jedem Workshop neu erzeugen
     (Secret in GitHub tauschen, Image neu bauen).
  2. Den Aufruf über das Go-Backend des Dashboards leiten, mit fester
     Nachricht. Das ist schwerer auszulesen, steckt aber weiterhin im Image.
  3. Ein kleines Relay außerhalb (etwa ein Worker mit Rate-Limit), das nur die
     feste Nachricht weitergibt. Nur hier liegt die echte URL nicht beim
     Teilnehmer.

### Entscheidungen vor der Umsetzung

Siehe offene Fragen 12 bis 17 unten.

**Aufwand:** klein als Bonus-Hinweis, mittel als Umstellung (Init-Schritt,
Referenz, Folie 11a, Story-Doc, Dashboard).

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
9. **Consul-Adresse (Nachtrag Story 3):** Merksatz „Laptop `localhost`,
   Container `consul`“ zentral in `docs/custom-setup.md` und von Story 3,
   Dashboard und Folie 11a verlinken, oder an jeder Stelle einzeln erklären?
10. **Env-Variablen (Nachtrag):** Eigene Variablen direkt in custom.yml
    erlauben (einfach, aber die Arbeitskopie wird dirty), oder eine
    ungetrackte `booking-custom.env` per `env_file` einführen?
11. **Consul wählt selbst (Nachtrag):** Prepared Query mit `limit=1` nur als
    Bonus in Story 3 erwähnen, oder den Standardweg darauf umstellen?
12. **A1 Bildfolge Story 4:** Flight, Hotel und Car nur im Bild, oder als
    Personen im Spiel?
13. **A2 Teil 1:** Kapazitätsidee („drei Hände“) trotz Überschneidung mit
    dem Bulkhead-Spiel aus Story 5, oder Wartezeit sichtbar aufaddieren?
14. **A4 Farbe für Flight:** Blau laut Palette mit Ausfall-Symbol, oder Rot
    (dann Palette und Kartenfarben anpassen)?
15. **B1 Einklinken als Hausaufgabe:** D5 aufheben?
16. **B3 Beispiel-Lösung:** per Compose-Profil standardmäßig aus, oder als
    Rückfallebene immer an?
17. **C Webhook:** Rotation nach jedem Workshop, Backend-Proxy oder externes
    Relay?
