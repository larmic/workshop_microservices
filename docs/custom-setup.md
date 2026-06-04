# Eigenen Booking-Service einklinken (Custom-Setup)

> **Workshop-Material:** Dieses Dokument nutzen wir gemeinsam im Workshop (ab
> Story 1). Für die Vorbereitung vorab gilt [vorbereitung.md](vorbereitung.md).

Im Dashboard kannst du pro Story zwischen *Reference* und *Custom* umschalten — die Buttons rufen dann den jeweils ausgewählten Service auf.

## Konzept

Stories bauen kumulativ aufeinander auf — ihr erweitert iterativ **einen** Service, statt für jede Story ein neues Projekt anzulegen. Per Default zeigt `.env` auf die eingecheckte Beispiel-Lösung (`services/booking/custom/`, in Kotlin/Ktor). Wer eine eigene Lösung in einer anderen Sprache baut, lässt diese in einem **separaten Verzeichnis** liegen (eigenes Repo, irgendwo auf der Platte) und passt nur den Pfad in `.env` an. Compose baut das Image dann aus diesem Pfad.

## Setup

1. **`.env` anlegen** im Verzeichnis `services/` (einmalig, in der Vorbereitung bereits passiert):
   ```bash
   cp .env.example .env
   ```
   Der Default-Eintrag zeigt auf die mitgelieferte Beispiel-Lösung:
   ```bash
   CUSTOM_BOOKING_PATH=./booking/custom
   ```
   Sobald ihr eure eigene Lösung einklinkt, tragt hier den (absoluten) Pfad ein — genau den Projektordner eures Mini-Service aus der Vorbereitung:
   ```bash
   CUSTOM_BOOKING_PATH=/absolute/path/to/your/booking-service
   ```
   Voraussetzung: Im angegebenen Verzeichnis liegt eine `Dockerfile`, der Service muss auf Port `8080` lauschen.

2. **Stack starten** wie in der Vorbereitung: `make docker-up-hub` (Referenz-Images von Docker Hub) oder `make docker-up` (alle Images lokal bauen). Der `booking-custom`-Container ist von Anfang an Teil des Stacks — kein separates Terminal nötig.

3. **Im Dashboard pro Story** zwischen *Reference* und *Custom* umschalten. Der Custom-Toggle ist grau, solange euer Service nicht erreichbar ist.

## Code-Iteration

Wenn ihr euren Service-Code geändert habt und nur den Custom-Container neu bauen wollt (ohne den Rest des Stacks anzufassen), läuft das in einem zweiten Terminal:

```bash
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml -f docker-compose.custom.yml up -d --build booking-custom
```

Alternativ einfach Ctrl+C im Stack-Terminal und erneut `make docker-up` — durch das Compose-Layer-Caching geht das meist schnell. Alle weiteren Make-Targets seht ihr mit `make help`.

## Interne URLs

Euer Service läuft im selben Docker-Netz wie die Referenz-Implementierung. Die internen Hostnamen sind identisch:

| Backend  | URL                          |
|----------|------------------------------|
| Flight   | `http://flight:8080`         |
| Hotel    | `http://hotel:8080`          |
| Car      | `http://car:8080`            |
| Consul   | `http://consul:8500`         |

Diese Werte werden eurem Container als Environment-Variablen `FLIGHT_SERVICE_URL`, `HOTEL_SERVICE_URL`, `CAR_SERVICE_URL`, `CONSUL_URL` übergeben.

## Tools & URLs

| Tool              | URL                              | Zweck                                |
|-------------------|----------------------------------|--------------------------------------|
| Dashboard         | <http://localhost>               | Workshop-Steuerzentrale              |
| Swagger UI        | <http://localhost/api>           | API-Dokumentation aller Services     |
| Consul UI         | <http://localhost/consul>        | Service Discovery & Health           |
| Traefik Dashboard | <http://localhost:8080>          | Routing & Reverse-Proxy-Monitoring   |
| Booking-Custom    | <http://localhost:8099>          | Direktzugriff auf euren Custom-Service (an Traefik vorbei) |

Schneller Health-Check vom Terminal:

```bash
curl http://localhost/api/flight/health
curl http://localhost/api/hotel/health
curl http://localhost/api/car/health
```

## Hinweise für Windows-Teilnehmer

- **Pfade** mit Forward-Slashes statt Backslashes: `CUSTOM_BOOKING_PATH=C:/Users/foo/projects/my-booking-service`
- **Line endings** der `.env` müssen LF sein (nicht CRLF) — sonst landet `\r` am Pfadende und der Build schlägt mit "no such file" fehl. Editor entsprechend einstellen oder `.env` über WSL/Git Bash erzeugen.
- **Build-Performance**: das Projekt sollte im WSL2-Filesystem liegen (z.B. `~/projects/...`), nicht unter `/mnt/c/...`.
- **Erreichbarkeit & Zertifikate**: Bei `localhost`-Problemen trotz laufender Container oder TLS-Fehlern beim `git clone` siehe [troubleshooting.md](troubleshooting.md).

## Stack stoppen

```bash
make docker-down
```

Stoppt alle Container — Referenz-Services, Infrastruktur und Custom-Service in einem Rutsch.

## Los geht's

Dein Arbeitsplatz ist eingerichtet. Starte mit [Story 1: Der erste cloud-native Booking-Service](stories/story-01-cloud-native-booking-service.md).
