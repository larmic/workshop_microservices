# Workshop-Vorbereitung

Diese Anleitung beschreibt, wie du das Repo auscheckst, die Umgebung startest und mit dem Workshop loslegst. **Abschnitt 0 ist Pflicht und wird vor dem Workshop-Termin erledigt.**

## Voraussetzungen

- **Git**
- **Docker** inkl. Docker Compose (harte Voraussetzung, ohne läuft die Workshop-Umgebung nicht)
- **IDE / Editor** nach Wahl
- Eigene Sprache/Framework deiner Wahl, in der du den Booking-Service umsetzen möchtest (Java/Spring Boot, Quarkus, Go, Node.js, ...)
- Du kannst in deiner Sprache ein neues Projekt aufsetzen, das einen HTTP-Endpoint bereitstellt und per Dockerfile ein Image baut. Genau das ist die Pflicht-Hausaufgabe in Abschnitt 0, im Workshop selbst üben wir es nicht.

## 0. Vor dem Workshop (Pflicht)

Diese Aufgaben erledigst du **vor** dem Workshop-Termin. Aufwand: ca. 1 bis 2 Stunden. Sie stellen sicher, dass wir am Workshop-Tag direkt mit den Inhalten starten, statt Setup-Probleme zu lösen.

### 0.1 Greenfield-Mini-Service (Pflicht-Hausaufgabe)

Setze in deiner Sprache/deinem Framework ein neues, leeres Projekt auf. Dieser Mini-Service ist genau die Basis, die du im Workshop über alle 7 Stories hinweg erweiterst. Merke dir den Projektordner: genau diesen Pfad trägst du im Workshop als `CUSTOM_BOOKING_PATH` ein (siehe Abschnitt 5), weggeworfen wird hier nichts.

Akzeptanzkriterien:

- [ ] Neues Projekt in deiner Sprache/deinem Framework
- [ ] `GET /health` antwortet mit HTTP 200
- [ ] Eine `Dockerfile` liegt im Projekt, `docker build` läuft fehlerfrei durch
- [ ] Der Container startet und lauscht auf Port `8080`

Mehr nicht: Body und Format der Health-Antwort sind egal, ein leeres HTTP 200 genügt. Logging, Info-Endpoint und Aufrufe der Backend-Services kommen erst im Workshop ab Story 1.

Smoke-Test:

```bash
docker build -t my-booking-service .
docker run --rm -p 8080:8080 my-booking-service

# zweites Terminal:
curl -i http://localhost:8080/health   # erwartet: HTTP/1.1 200
```

Danach den Container stoppen (Ctrl+C): der Host-Port `8080` wird gleich vom Workshop-Stack belegt (Traefik-Dashboard).

### 0.2 Workshop-Stack einmal starten

Klone das Repo und fahre den Stack einmal hoch. Die Schritt-für-Schritt-Anleitung steht in den Abschnitten 1 bis 3.

- [ ] Repo geklont, `.env` angelegt
- [ ] `make docker-up-hub` läuft durch
- [ ] Dashboard unter <http://localhost> erreichbar

Die per `cp .env.example .env` angelegte `.env` funktioniert mit ihrem Default-Pfad sofort. Dein eigener Booking-Pfad aus 0.1 wird erst im Workshop eingetragen (Abschnitt 5).

Windows/WSL2: beachte vorab die Windows-Hinweise in Abschnitt 5 (Line endings der `.env` müssen LF sein, Repo ins WSL2-Filesystem legen).

### 0.3 Im Firmennetz testen

Führe `git clone` und `docker build` einmal **in dem Netz aus, in dem du auch am Workshop-Tag arbeitest** (Firmen-WLAN, VPN). Firmen-Proxies und TLS-Zertifikate fallen sonst erst am Workshop-Tag auf, besonders unter Windows/WSL2. Bei Problemen: [troubleshooting.md](troubleshooting.md).

### 0.4 Bei Problemen

Wenn etwas hakt und [troubleshooting.md](troubleshooting.md) nicht weiterhilft: melde dich **vor** dem Workshop kurz per Mail oder Chat beim Trainer. Am Workshop-Tag selbst bleibt für Setup-Probleme wenig Zeit.

## 1. Repository klonen

```bash
git clone https://github.com/larmic/workshop_microservices.git
cd workshop_microservices/services
cp .env.example .env
```

Alle weiteren Befehle werden aus dem Verzeichnis `services/` ausgeführt.

## 2. Umgebung starten

Es gibt zwei Wege, die Workshop-Umgebung hochzufahren — beide bringen die gleichen Container ans Laufen (Flight, Hotel, Car, Booking-Referenzlösungen, **Booking-Custom**, Traefik, Consul, Swagger UI, Dashboard).

Der Booking-Custom-Service zeigt per Default auf die eingecheckte Beispiel-Lösung unter `services/booking/custom/`. Wer eine eigene Lösung baut, ändert dafür nur den Pfad in `.env` (siehe Abschnitt 5).

### Variante A — Images von Docker Hub ziehen (empfohlen für Teilnehmer)

Schnell und ohne lokalen Build der Referenz-Services. Der Booking-Custom-Service wird trotzdem lokal gebaut (kein Hub-Image dafür):

```bash
make docker-up-hub
```

Im Hintergrund: `docker compose ... pull` lädt fertige Images von Docker Hub (`larmic/workshop-microservices-*`), `docker compose ... build booking-custom` baut den Custom-Service, dann `docker compose ... up`.

### Variante B — Images lokal bauen

Wenn du Änderungen am Code oder an den Dockerfiles vornimmst, baust du die Images lokal:

```bash
make docker-up
```

Das entspricht `docker compose ... up --build` über alle vier Compose-Dateien.

### Ohne Makefile

Falls `make` nicht verfügbar ist, gehen beide Varianten auch direkt:

```bash
# Variante A (Docker Hub für Referenz, lokal bauen für Custom)
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml pull
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml -f docker-compose.custom.yml build booking-custom
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml -f docker-compose.custom.yml up

# Variante B (alles lokal bauen)
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml -f docker-compose.custom.yml up --build
```

Alle weiteren Make-Targets siehst du mit `make help`.

> Wenn etwas nicht startet oder das Dashboard nicht lädt: siehe
> [troubleshooting.md](troubleshooting.md) (Port-Konflikte, Image-Pull,
> Windows/WSL2-Erreichbarkeit, Zertifikate).

## 3. Dashboard öffnen

Das **Dashboard** ist das zentrale Tool für den Workshop. Hier startest, stoppst und skalierst du Services, siehst Health-Status und kommst in einem Klick zu den anderen Tools.

> **<http://localhost>**

![Workshop Dashboard](assets/dashboard.png)

## 4. Weitere URLs

| Tool              | URL                              | Zweck                                |
|-------------------|----------------------------------|--------------------------------------|
| Dashboard         | <http://localhost>               | Workshop-Steuerzentrale              |
| Swagger UI        | <http://localhost/api>           | API-Dokumentation aller Services     |
| Consul UI         | <http://localhost/consul>        | Service Discovery & Health           |
| Traefik Dashboard | <http://localhost:8080>          | Routing & Reverse-Proxy-Monitoring   |
| Booking-Custom    | <http://localhost:8099>          | Direktzugriff auf deinen Custom-Service (an Traefik vorbei) |

Schneller Health-Check vom Terminal:

```bash
curl http://localhost/api/flight/health
curl http://localhost/api/hotel/health
curl http://localhost/api/car/health
```

## 5. Eigenen Booking-Service einklinken (Custom-Setup, während des Workshops)

Im Dashboard kannst du pro Story zwischen *Reference* und *Custom* umschalten — die Buttons rufen dann den jeweils ausgewählten Service auf.

### Konzept

Stories bauen kumulativ aufeinander auf — ihr erweitert iterativ **einen** Service, statt für jede Story ein neues Projekt anzulegen. Per Default zeigt `.env` auf die eingecheckte Beispiel-Lösung (`services/booking/custom/`, in Kotlin/Ktor). Wer eine eigene Lösung in einer anderen Sprache baut, lässt diese in einem **separaten Verzeichnis** liegen (eigenes Repo, irgendwo auf der Platte) und passt nur den Pfad in `.env` an. Compose baut das Image dann aus diesem Pfad.

### Setup

1. **`.env` anlegen** im Verzeichnis `services/` (einmalig):
   ```bash
   cp .env.example .env
   ```
   Der Default-Eintrag zeigt auf die mitgelieferte Beispiel-Lösung:
   ```bash
   CUSTOM_BOOKING_PATH=./booking/custom
   ```
   Sobald ihr eine eigene Lösung baut, tragt hier den (absoluten) Pfad ein:
   ```bash
   CUSTOM_BOOKING_PATH=/absolute/path/to/your/booking-service
   ```
   Voraussetzung: Im angegebenen Verzeichnis liegt eine `Dockerfile`, der Service muss auf Port `8080` lauschen.

2. **Stack starten** wie üblich (siehe Abschnitt 2): `make docker-up-hub` oder `make docker-up`. Der `booking-custom`-Container ist von Anfang an Teil des Stacks — kein separates Terminal nötig.

3. **Im Dashboard pro Story** zwischen *Reference* und *Custom* umschalten. Der Custom-Toggle ist grau, solange euer Service nicht erreichbar ist.

### Code-Iteration

Wenn ihr euren Service-Code geändert habt und nur den Custom-Container neu bauen wollt (ohne den Rest des Stacks anzufassen), läuft das in einem zweiten Terminal:

```bash
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml -f docker-compose.custom.yml up -d --build booking-custom
```

Alternativ einfach Ctrl+C im Stack-Terminal und erneut `make docker-up` — durch das Compose-Layer-Caching geht das meist schnell.

### Interne URLs

Euer Service läuft im selben Docker-Netz wie die Referenz-Implementierung. Die internen Hostnamen sind identisch:

| Backend  | URL                          |
|----------|------------------------------|
| Flight   | `http://flight:8080`         |
| Hotel    | `http://hotel:8080`          |
| Car      | `http://car:8080`            |
| Consul   | `http://consul:8500`         |

Diese Werte werden eurem Container als Environment-Variablen `FLIGHT_SERVICE_URL`, `HOTEL_SERVICE_URL`, `CAR_SERVICE_URL`, `CONSUL_URL` übergeben.

### Hinweise für Windows-Teilnehmer

- **Pfade** mit Forward-Slashes statt Backslashes: `CUSTOM_BOOKING_PATH=C:/Users/foo/projects/my-booking-service`
- **Line endings** der `.env` müssen LF sein (nicht CRLF) — sonst landet `\r` am Pfadende und der Build schlägt mit "no such file" fehl. Editor entsprechend einstellen oder `.env` über WSL/Git Bash erzeugen.
- **Build-Performance**: das Projekt sollte im WSL2-Filesystem liegen (z.B. `~/projects/...`), nicht unter `/mnt/c/...`.
- **Erreichbarkeit & Zertifikate**: Bei `localhost`-Problemen trotz laufender Container oder TLS-Fehlern beim `git clone` siehe [troubleshooting.md](troubleshooting.md).

### Stack stoppen

```bash
make docker-down
```

Stoppt alle Container — Referenz-Services, Infrastruktur und Custom-Service in einem Rutsch.

## 6. Los geht's

Dein Arbeitsplatz ist eingerichtet. Starte mit [Story 1: Der erste cloud-native Booking-Service](stories/story-01-cloud-native-booking-service.md).
