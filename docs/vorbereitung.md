# Workshop-Vorbereitung

Diese Vorbereitung ist **Pflicht** und wird **vor** dem Workshop-Termin erledigt. Aufwand: ca. 1 bis 2 Stunden. Sie stellt sicher, dass wir am Workshop-Tag direkt mit den Inhalten starten, statt Setup-Probleme zu lösen.

## Voraussetzungen

- **Git**
- **Docker** inkl. Docker Compose (harte Voraussetzung, ohne läuft die Workshop-Umgebung nicht)
- **IDE / Editor** nach Wahl
- Eigene Sprache/Framework deiner Wahl, in der du den Booking-Service umsetzen möchtest (Java/Spring Boot, Quarkus, Go, Node.js, ...)
- Du kannst in deiner Sprache ein neues Projekt aufsetzen, das einen HTTP-Endpoint bereitstellt und per Dockerfile ein Image baut. Genau das ist die Pflicht-Hausaufgabe in Aufgabe 1, im Workshop selbst üben wir es nicht.

## 1. Eigenen Mini-Service bauen (Pflicht-Hausaufgabe)

Setze in deiner Sprache/deinem Framework ein neues, leeres Projekt auf. Diesen Mini-Service erweiterst du im Workshop über alle 7 Stories hinweg. Wie er in die Workshop-Umgebung eingebunden wird, machen wir gemeinsam im Workshop.

Akzeptanzkriterien:

- [ ] Neues Projekt in deiner Sprache/deinem Framework
- [ ] `GET /health` antwortet mit HTTP 200
- [ ] Eine `Dockerfile` liegt im Projekt, `docker build` läuft fehlerfrei durch
- [ ] Der Container startet und lauscht auf Port `8080`

Mehr nicht: Body und Format der Health-Antwort sind egal, ein leeres HTTP 200 genügt. Alles Weitere kommt erst im Workshop.

Smoke-Test:

```bash
docker build -t my-booking-service .
docker run --rm -p 8080:8080 my-booking-service

# zweites Terminal:
curl -i http://localhost:8080/health   # erwartet: HTTP/1.1 200
```

Danach den Container stoppen (Ctrl+C): der Host-Port `8080` wird gleich von der Workshop-Umgebung belegt.

## 2. Workshop-Umgebung auschecken und prüfen

```bash
git clone https://github.com/larmic/workshop_microservices.git
cd workshop_microservices/services
cp .env.example .env
make docker-up-hub
```

Das zieht fertige Images von Docker Hub und startet alle Container. Die `.env` funktioniert so, wie sie ist. Falls `make` nicht verfügbar ist:

```bash
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml pull
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml -f docker-compose.custom.yml build booking-custom
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml -f docker-compose.custom.yml up
```

Erfolgreich vorbereitet bist du, wenn:

- [ ] Repo geklont, `.env` angelegt
- [ ] Der Stack läuft ohne Fehler
- [ ] Das **Dashboard** unter <http://localhost> lädt

![Workshop Dashboard](assets/dashboard.png)

Hinweise für Windows-Teilnehmer: Line endings der `.env` müssen LF sein (nicht CRLF), und das Repo sollte im WSL2-Filesystem liegen (z.B. `~/projects/...`, nicht `/mnt/c/...`).

> Wenn etwas nicht startet oder das Dashboard nicht lädt: siehe
> [troubleshooting.md](troubleshooting.md) (Port-Konflikte, Image-Pull,
> Windows/WSL2-Erreichbarkeit, Zertifikate).

## 3. Im Firmennetz testen

Führe `git clone` und `docker build` einmal **in dem Netz aus, in dem du auch am Workshop-Tag arbeitest** (Firmen-WLAN, VPN). Firmen-Proxies und TLS-Zertifikate fallen sonst erst am Workshop-Tag auf, besonders unter Windows/WSL2. Bei Problemen: [troubleshooting.md](troubleshooting.md).

## 4. Bei Problemen

Wenn etwas hakt und [troubleshooting.md](troubleshooting.md) nicht weiterhilft: melde dich **vor** dem Workshop kurz per Mail oder Chat beim Trainer. Am Workshop-Tag selbst bleibt für Setup-Probleme wenig Zeit.

---

Das war's. Alles Weitere machen wir gemeinsam im Workshop.
