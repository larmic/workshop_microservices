# Workshop-Vorbereitung

Diese Anleitung beschreibt, wie du das Repo auscheckst, die Umgebung startest und mit dem Workshop loslegst.

## Voraussetzungen

- **Git**
- **Docker** inkl. Docker Compose
- **IDE / Editor** nach Wahl (optional, nur wenn du die Stories selbst implementieren willst)
- Eigene Sprache/Framework deiner Wahl, in der du den Booking-Service umsetzen möchtest (Java/Spring Boot, Quarkus, Go, Node.js, ...)

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

## 5. Eigenen Booking-Service einklinken (Custom-Setup)

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

### Stack stoppen

```bash
make docker-down
```

Stoppt alle Container — Referenz-Services, Infrastruktur und Custom-Service in einem Rutsch.

## 6. Los geht's

Dein Arbeitsplatz ist eingerichtet. Starte mit [Story 1: Der erste cloud-native Booking-Service](stories/story-01-cloud-native-booking-service.md).
