# Troubleshooting

Häufige Probleme beim Starten der Workshop-Umgebung und wie man sie eingrenzt.
Voraussetzung und Setup stehen in [vorbereitung.md](vorbereitung.md). Der gesamte
Zugriff läuft über Traefik auf **Port 80** (`http://localhost`), das
Traefik-Dashboard auf **Port 8080**.

Jeder Eintrag ist mit zwei Markern versehen:

- **Geltung:** `macOS/Linux/Windows` (überall) oder `nur Windows/WSL2`
- **Status:** `verifiziert` (am 2026-06-02 gegen das Repo und lokal auf macOS
  getestet) oder `berichtet` (aus dem Kickoff-Feedback vom 2026-05-30, unter
  Windows aufgetreten, hier nicht reproduziert, mit Links zur Eigendiagnose)

> Die Windows/WSL2-Einträge sind aus Teilnehmer-Rückmeldungen abgeleitet und noch
> nicht von uns unter Windows verifiziert. Wer sie unter Windows prüft, möge sie
> bestätigen oder korrigieren.

## Erste Diagnose (immer zuerst)

Alle Befehle aus dem Verzeichnis `services/`. Damit `docker compose` ohne die
vielen `-f`-Flags funktioniert, einmal pro Terminal die Compose-Dateien als Env
setzen:

```bash
export COMPOSE_FILE=docker-compose.yml:docker-compose.infra.yml:docker-compose.reference.yml:docker-compose.custom.yml
```

Dann der Reihe nach:

```bash
# 1. Welche Container laufen, welche sind unhealthy oder beendet?
docker compose ps

# 2. Logs eines auffälligen Containers (Beispiel Traefik)
docker compose logs traefik

# 3. Erreichbarkeit der Tools (HTTP-Status, ohne Body)
curl -sS -o /dev/null -w "%{http_code}\n" http://localhost          # Dashboard, erwartet 200
curl -sS -o /dev/null -w "%{http_code}\n" http://localhost/api/flight/health   # erwartet 200
```

Erwartetes Bild eines gesunden Stacks: `consul` ist `healthy`, `traefik` läuft mit
Port-Mapping `0.0.0.0:80->80` und `0.0.0.0:8080->8080`, `http://localhost`
antwortet mit `200`.

## Dashboard / `http://localhost` lädt nicht

**Geltung:** macOS/Linux/Windows · **Status:** verifiziert

Mögliche Ursachen und Diagnose:

1. **Container noch nicht oben.** Traefik startet erst, wenn das Dashboard
   gestartet ist (`depends_on`), und registriert seine Routen über den
   File-Provider. Direkt nach `up` kann es ein paar Sekunden dauern. Prüfen, ob
   die Routen registriert sind:
   ```bash
   curl -s http://localhost:8080/api/http/routers | grep -o '"name":"[^"]*"'
   ```
   Erwartet werden u. a. `root@file`, `dashboard@file`, `swagger-ui@file`,
   `consul@file`. Fehlen sie, kurz warten und neu laden.
2. **Dashboard-Container nicht gesund.** `docker compose ps` und
   `docker compose logs dashboard`. Das Dashboard braucht Zugriff auf den
   Docker-Socket (`/var/run/docker.sock`), siehe
   `services/docker-compose.infra.yml`.
3. **Browser-Cache / alte Weiterleitung.** Hartes Neuladen oder ein privates
   Fenster ausprobieren.

## Port belegt (Adresse bereits in Verwendung)

**Geltung:** macOS/Linux/Windows · **Status:** verifiziert

Symptom: `docker compose up` bricht ab mit `bind: address already in use` (oder
`Ports are not available`). Der Stack belegt diese Host-Ports:

| Port | Dienst |
|---|---|
| 80 | Traefik (alle Tools laufen darüber) |
| 8080 | Traefik-Dashboard |
| 8500 | Consul |
| 8084 | Swagger UI (direkt) |
| 8085 bis 8091 | Booking-Referenz Story 1 bis 7 (direkt) |
| 8099 | Booking-Custom (direkt, an Traefik vorbei) |

Belegung finden:

```bash
# macOS / Linux
lsof -nP -iTCP:80 -sTCP:LISTEN

# Windows (PowerShell)
netstat -ano | findstr :80
```

Häufigster Kandidat für Port 80 ist ein lokaler Webserver oder ein anderer
Reverse-Proxy. Den belegenden Prozess beenden oder die Umgebung freiräumen
(`docker compose down`), dann neu starten.

## Image-Pull schlägt fehl

**Geltung:** macOS/Linux/Windows · **Status:** verifiziert

`make docker-up-hub` zieht die Referenz-Images von Docker Hub
(`larmic/workshop-microservices-*`). Schlägt der Pull fehl:

- **Docker-Daemon läuft?** `docker version` muss eine Server-Version zeigen.
- **Rate-Limit / Login.** Bei `toomanyrequests` einmal `docker login`.
- **Netz / Proxy.** Hinter Firmen-Proxy den Docker-Proxy konfigurieren (Docker
  Desktop: Settings, Resources, Proxies).
- **Fallback:** lokal bauen statt ziehen mit `make docker-up`.

## Nur Windows/WSL2: `http://localhost` nicht erreichbar, obwohl alle Container laufen

**Geltung:** nur Windows/WSL2 · **Status:** berichtet (2026-05-30, nicht reproduziert)

Beobachtung im Kickoff: Alle Container liefen, Traefik zeigte die Route korrekt,
trotzdem kam unter Windows kein Zugriff auf `http://localhost` zustande. Erst nach
mehrmaligem Neustarten klappte es. Die Traefik-Konfiguration selbst ist unauffällig
(`services/traefik/dynamic.yml`); die Ursache liegt sehr wahrscheinlich im
WSL2/Docker-Desktop-Networking bzw. im Timing der Port-Weiterleitung.

Vorgehen:

1. **Kurz warten und neu laden.** Die Host-zu-Container-Weiterleitung über WSL2
   steht manchmal erst Sekunden nach dem Container-Start. Erst die Routen prüfen
   (siehe oben), dann neu laden.
2. **Stack einmal neu starten.** `docker compose down` und erneut hochfahren.
3. **`127.0.0.1` statt `localhost`** im Browser testen. In manchen WSL2-Setups
   reagiert nur eine der beiden Adressen
   ([docker/for-win#13182](https://github.com/docker/for-win/issues/13182),
   [microsoft/WSL#4983](https://github.com/microsoft/WSL/issues/4983)).
4. **Projekt im WSL2-Dateisystem.** Liegt das Repo unter `/mnt/c/...`, ist nicht
   nur der Build langsam, es gibt auch mehr Networking-Reibung. Empfohlen ist
   `~/projects/...` innerhalb der WSL-Distribution (siehe
   [vorbereitung.md](vorbereitung.md), Windows-Hinweise).
5. **Docker-Desktop-Integration.** Settings, Resources, WSL Integration: die
   genutzte Distribution muss aktiviert sein. Hintergrund:
   [Docker Desktop WSL2 Backend](https://docs.docker.com/desktop/wsl/).

## Nur Windows/WSL2: Zertifikatsfehler beim `git clone` oder Build

**Geltung:** nur Windows/WSL2 · **Status:** berichtet (2026-05-30, nicht reproduziert)

Beobachtung im Kickoff: In der WSL war kein Zertifikat hinterlegt, `git clone`
(bzw. Downloads im Build) scheiterten an der TLS-Prüfung. Workaround vor Ort war
der Wechsel ins Gast-WLAN. Typische Ursache ist ein Firmen-Proxy oder eine
Firewall, die TLS aufbricht und ein eigenes (selbstsigniertes) Zertifikat
vorlegt, das die WSL-Distribution nicht kennt.

Lösungen, sauberste zuerst:

1. **Anderes Netz testen.** Gast-WLAN oder Hotspot bestätigt schnell, ob ein
   Proxy die Ursache ist.
2. **Firmen-Root-CA in der WSL hinterlegen** (dauerhafte, korrekte Lösung). Das
   `.crt` der IT besorgen und in den Trust-Store der Distribution legen:
   ```bash
   sudo cp firmen-root-ca.crt /usr/local/share/ca-certificates/
   sudo update-ca-certificates
   ```
3. **Nur Git auf ein CA-Bundle zeigen lassen** (falls der System-Trust nicht
   geändert werden darf):
   ```bash
   git config --global http.sslCAInfo /pfad/zu/firmen-root-ca.pem
   ```
4. **Nicht empfohlen:** `git config --global http.sslVerify false` deaktiviert die
   Prüfung komplett und ist ein Sicherheitsrisiko. Höchstens kurzfristig in einem
   vertrauenswürdigen Netz.

Weiterführend:
[Git und SSL-Zertifikate](https://www.codestudy.net/blog/how-to-solve-ssl-certificate-self-signed-certificate-when-cloning-repo-from-github/).

## Weitere Windows-Stolpersteine

**Geltung:** nur Windows/WSL2 · **Status:** verifiziert (in [vorbereitung.md](vorbereitung.md) dokumentiert)

Pfade mit Forward-Slashes, `.env` mit LF statt CRLF, Projekt im WSL2-Dateisystem:
siehe die Windows-Hinweise in [vorbereitung.md](vorbereitung.md).

## Stack komplett zurücksetzen

**Geltung:** macOS/Linux/Windows · **Status:** verifiziert

Wenn nichts mehr zusammenpasst, hilft ein sauberer Neustart:

```bash
make docker-down     # stoppt alles, entfernt Container, Netzwerke, anonyme Volumes
make docker-up-hub   # zieht Images neu und startet
```
