# Load Balancing & Service Discovery

Dieses Dokument erklaert, wie Load Balancing und Service Discovery im Workshop-Setup zusammenspielen, insbesondere bei skalierten Services.

## Der Weg eines Requests

Wenn z.B. `http://localhost/api/car/cars` aufgerufen wird (via Swagger, Browser oder curl):

```
Browser / Swagger UI
        |
        v
  Traefik (Port 80)
        |   PathPrefix("/api/car") -> stripPrefix -> http://car:8080
        v
  Docker DNS ("car")
        |   Round-Robin auf alle Container mit Service-Name "car"
        v
  car-1, car-2 oder car-3  (je nach Skalierung)
```

## Schichten im Detail

### 1. Traefik (API Gateway / Reverse Proxy)

Traefik routet Requests anhand des URL-Pfads zum richtigen Service. Die Konfiguration in `traefik/dynamic.yml` definiert pro Service **nur einen Hostnamen**:

```yaml
car:
  loadBalancer:
    servers:
      - url: "http://car:8080"
```

Traefik kennt die einzelnen Container-Instanzen **nicht**. Es schickt alle Requests an den Docker-DNS-Namen `car`. Das stripPrefix-Middleware entfernt den Pfad-Prefix (`/api/car`), sodass der Service den Request als `/cars` erhaelt.

### 2. Docker DNS (das eigentliche Load Balancing)

Docker Compose erstellt fuer jeden Service-Namen einen DNS-Eintrag im internen Netzwerk. Bei `--scale car=3` loest der Hostname `car` auf **drei verschiedene Container-IPs** auf. Docker's eingebauter DNS-Server verteilt die Anfragen per **Round-Robin** auf alle laufenden Container.

Das bedeutet: Traefik fragt bei jedem Request den Hostnamen `car` auf, und Docker DNS liefert abwechselnd eine andere IP zurueck.

### 3. Consul (Service Registry)

Jeder Service registriert sich beim Start bei Consul mit einer **eindeutigen Service-ID** (z.B. `car-service-<container-hostname>`). Consul bietet:

- **Service Registry**: Ueberblick ueber alle laufenden Instanzen
- **Health Checks**: Regelmaessige Pruefung des `/health`-Endpoints (alle 10s)
- **UI**: Visuelle Darstellung unter `http://localhost:8500`

**Wichtig**: Traefik nutzt Consul in diesem Setup **nicht** fuer Service Discovery. Traefik verwendet den `file`-Provider (statische Konfiguration), nicht den `consul`-Provider. Consul dient hier primaer als Informationsquelle und wird von Booking Story 2 aktiv genutzt.

### 4. Booking Story 2 (Consul-basierte Discovery)

Im Gegensatz zu Booking Story 1 (hardcodierte URLs) nutzt Story 2 den `consul.Resolver`, um zur Laufzeit eine gesunde Service-Instanz zu finden:

```
Booking Story 2
    |   consul.ResolveServiceURL("car-service")
    v
Consul API (/v1/health/service/car-service?passing=true)
    |   Liefert Liste gesunder Instanzen
    v
Zufaellige Auswahl einer Instanz -> direkter HTTP-Call
```

## Uebersicht: Wer kennt die Instanzen?

| Schicht          | Aufgabe                                 | Kennt einzelne Instanzen? |
|------------------|-----------------------------------------|---------------------------|
| **Traefik**      | Routing (Pfad -> Service-Hostname)      | Nein, nur den Hostnamen   |
| **Docker DNS**   | Namensaufloesung + Round-Robin          | Ja, alle laufenden Container |
| **Consul**       | Service Registry + Health Checks        | Ja, alle registrierten Instanzen |
| **Booking Story 2** | Service Discovery via Consul        | Ja, waehlt zufaellige gesunde Instanz |

## Ausprobieren

```bash
# Car-Service auf 3 Instanzen skalieren (CLI)
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up -d --scale car=3

# Oder ueber das Dashboard: http://localhost/dashboard
```

Nach dem Skalieren:

- **Consul UI** (`http://localhost:8500`): Zeigt 3 registrierte `car-service`-Instanzen mit Health Status
- **Traefik Dashboard** (`http://localhost:8080`): Zeigt weiterhin nur `http://car:8080` als Backend
- **Swagger UI** (`http://localhost`): Requests an `/api/car/cars` landen abwechselnd auf verschiedenen Instanzen

Um zu sehen, welche Instanz antwortet, kann der `/info`-Endpoint genutzt werden:

```bash
# Mehrfach aufrufen -- die serviceAddress bleibt gleich, aber
# der Request landet auf unterschiedlichen Containern
curl http://localhost/api/car/info
curl http://localhost/api/car/info
curl http://localhost/api/car/info
```

## Diskussionspunkte fuer den Workshop

- **Docker DNS vs. dedizierter Load Balancer**: Docker DNS reicht fuer einfaches Round-Robin, bietet aber kein Health-Check-basiertes Routing, keine Gewichtung und keine Sticky Sessions.
- **Consul als Single Source of Truth**: Consul weiss ueber Health Checks, welche Instanzen wirklich gesund sind. Docker DNS kennt nur "Container laeuft", nicht "Container ist bereit".
- **File-Provider vs. Consul-Provider in Traefik**: Mit dem Consul-Provider koennte Traefik automatisch alle Instanzen erkennen und gezielt load-balancen, statt sich auf Docker DNS zu verlassen.
