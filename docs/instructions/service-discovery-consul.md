# Service Discovery mit Consul — Workshop-Notizen

> Trainer-Notizen für Story 2. Reihenfolge folgt einer typischen Slide-Sequenz; jeder Abschnitt hat einen kleinen Hinweis, was an der Tafel/auf der Folie passieren sollte.

## 1. Worum geht es?

**Die Analogie:** Das Telefonbuch. Du willst „Pizzeria Roma" anrufen, kennst aber die Nummer nicht auswendig. Du schlägst im Telefonbuch unter dem **Namen** nach und bekommst eine **Nummer** zurück. Zieht die Pizzeria um, ändert sich der Eintrag im Telefonbuch — du als Anrufer musst nichts wissen. Eine Service Registry wie Consul ist genau das: ein **lebendes Telefonbuch** für Services.

**Das Problem:** Service A ruft Service B unter einer fest verdrahteten URL `http://flight:8081`. Was passiert, wenn …
- … B unter einem anderen Port hochfährt?
- … wir B auf 3 Instanzen skalieren — wer entscheidet, welche der drei der Aufrufer trifft?
- … eine der drei Instanzen gerade kaputt ist — soll der Aufrufer das selbst rausfinden?
- … in einer Container-Umgebung jeder Restart eine andere IP bekommt?

Hardcodierte URLs sind in der Konferenz-Folie schön. Im Betrieb sind sie der häufigste Grund für nächtliche Anrufe.

**Die Lösung:** Eine zentrale **Service Registry**. Jeder Service registriert sich beim Start mit seinem logischen Namen (z.B. `flight-service`) und seiner aktuellen Adresse. Aufrufer fragen die Registry: „Gib mir alle gesunden Instanzen von `flight-service`." Die Registry weiß auch, **wer gerade gesund ist** — über regelmäßige Health-Checks.

> 🎯 *Demo-Einstieg:* Im Workshop-Stack das Dashboard öffnen, `hotel`-Service auf 3 Replicas skalieren. Parallel die Consul-UI (`http://localhost:8500`) zeigen — dort ploppen die neuen Instanzen live auf, jede mit grünem Health-Check. Dann eine killen → wird rot → verschwindet aus `passing=true`. Der Booking-Service routet automatisch auf die übrigen.

---

## 2. Die drei Bausteine

```
   ┌──────────────┐    1. Register      ┌──────────────┐
   │   Service B  │────────────────────▶│              │
   │  (Provider)  │                     │              │
   │              │◀────────────────────│    Consul    │
   │              │   2. Health-Check   │   (Registry) │
   └──────────────┘     periodisch      │              │
                                        │              │
   ┌──────────────┐                     │              │
   │   Service A  │   3. Resolve        │              │
   │   (Caller)   │────────────────────▶│              │
   │              │◀────────────────────│              │
   │              │  Liste gesunder     └──────────────┘
   └──────┬───────┘     Instanzen
          │
          │ 4. HTTP-Call direkt an gewählte Instanz
          ▼
   ┌──────────────┐
   │  Service B   │
   │  (Instanz #2)│
   └──────────────┘
```

### Registry — das Telefonbuch
- Speichert pro Service: **Name** (logisch, z.B. `flight-service`), **ID** (eindeutig pro Instanz), **Address**, **Port**.
- Neue Instanz? Eintrag dazu. Instanz weg? Eintrag raus.
- Im Workshop: Consul-Dev-Container, läuft im Compose-Netz unter dem DNS-Namen `consul`.

### Health-Check — der Wahrheits-Beweis
- Consul fragt periodisch (Workshop: alle 10s) den `/health`-Endpoint jeder Instanz.
- Status `passing` / `warning` / `critical`.
- Ein Caller, der `?passing=true` filtert, sieht nur gesunde Instanzen.

### Resolver — die Telefonbuchabfrage im Caller
- Vor jedem Outbound-Call: „Welche `flight-service`-Instanzen sind gerade gesund?"
- Antwort: Liste mit URLs. Eine wird ausgewählt (Random / Round-Robin / gewichtet).
- Direktcall an die gewählte Instanz — **kein** weiterer Hop über Consul.

> ⚠️ **Stolperstein für Diskussion:** Consul ist nur *Telefonbuch*, nicht *Vermittlungsstelle*. Die HTTP-Requests fließen direkt zwischen den Services. Fällt Consul aus, können bestehende Caller mit gecachten Daten weiter funktionieren — neue Resolves schlagen aber fehl.

---

## 3. Was muss konfiguriert werden?

### Beim Registrieren

| Feld | Workshop-Wert | Bedeutung |
|---|---|---|
| `Name` | `flight-service`, `hotel-service`, `car-service` | Logischer Service-Name. Das ist, wonach Caller fragen. **Mehrere Instanzen teilen denselben Name.** |
| `ID` | `{Name}-{Hostname}` z.B. `flight-service-flight-1` | Eindeutige Instanz-ID. **Pro Instanz unterschiedlich**, sonst überschreibt Instanz #2 den Eintrag von #1. |
| `Address` | Container-Hostname (z.B. `flight-1`) | Wohin Caller später ihre Requests schicken. |
| `Port` | `8080` (intern im Container) | Port auf dem die Instanz lauscht. |
| `Check.HTTP` | `http://{Address}:{Port}/health` | Endpoint, den Consul für Health-Checks aufruft. |
| `Check.Interval` | `10s` | Wie oft Consul prüft. Zu kurz = Last + Lärm; zu lang = lange falsche „grün"-Anzeige nach Crash. |
| `Check.DeregisterCriticalServiceAfter` | (nicht gesetzt im Workshop) | Wenn Health-Check X Minuten lang `critical` ist, automatisch entfernen. Verhindert „Zombie-Einträge" nach Crashes. **Diskussionspunkt!** |

> 🎯 *Folie:* Diese Tabelle zeigen und betonen: Es gibt keine universellen Werte. Sie hängen ab von:
> - Wie schnell muss der Failover sein? (Intervall)
> - Wie viele Instanzen × Services hast du? (Last auf Consul)
> - Wie zuverlässig ist Graceful Shutdown? (DeregisterCriticalServiceAfter)

### Eindeutige Service-ID — wichtiger als man denkt

Wenn zwei Instanzen sich mit derselben ID registrieren, **überschreibt** die zweite die erste — Consul denkt, es gäbe nur eine. Workshop-Code nutzt `{Name}-{Hostname}`, weil Docker pro Container einen unique Hostnamen vergibt. In Kubernetes bietet sich der Pod-Name an, in VMs die Maschine + Port-Kombination.

---

## 4. Pseudocode

Drei Lifecycle-Schritte: **Registrieren** beim Start, **Deregistrieren** beim Shutdown, **Auflösen** bei jedem Outbound-Call.

### A) Registrieren beim Start

```kotlin
fun register(consulUrl, service): String {
    val payload = mapOf(
        "Name"    to service.name,                              // "flight-service"
        "ID"      to "${service.name}-${hostname()}",           // eindeutig pro Instanz
        "Address" to hostname(),
        "Port"    to service.port,
        "Check"   to mapOf(
            "HTTP"     to "http://${hostname()}:${service.port}/health",
            "Interval" to "10s",
        ),
    )

    var backoff = 1.seconds
    repeat(5) {
        if (httpPut("$consulUrl/v1/agent/service/register", payload)) {
            return payload["ID"]
        }
        sleep(backoff)
        backoff *= 2
    }
    throw RegistrationFailed
}
```

> 💡 **Warum Retry?** Bei `docker compose up` startet alles parallel — der Service kann schneller hochfahren als Consul. Ohne Retry wäre das Workshop-Setup flaky. Diskussionsfrage: „Was wäre die Alternative?" → `depends_on` mit Healthcheck, Init-Container, Sidecar.

### B) Deregistrieren beim Shutdown

```kotlin
fun shutdown() {
    deregister(consulUrl, serviceId)            // 1. erst aus der Registry
    server.shutdown(timeout = 5.seconds)        // 2. dann den HTTP-Server
}

fun deregister(consulUrl, serviceId) {
    httpPut("$consulUrl/v1/agent/service/deregister/$serviceId")
}
```

> ⚠️ **Reihenfolge ist entscheidend.** Erst aus Consul raus, *dann* den Server schließen. Sonst fragt ein Caller bei Consul nach, bekommt die sterbende Instanz zurück und läuft in Connection Refused.

> 💡 **Was bei `kill -9`?** Dann läuft `deregister` nicht. Der Eintrag bleibt bis zum nächsten Health-Check (im Workshop bis zu 10s). Lösung: `DeregisterCriticalServiceAfter` mit setzen — Consul räumt selbst auf.

### C) Auflösen mit Client-Side Load Balancing

```kotlin
fun resolve(consulUrl, name): String {
    val instances = httpGet("$consulUrl/v1/health/service/$name?passing=true")

    if (instances.isEmpty()) {
        throw NoHealthyInstance(name)
    }

    val pick = instances.random()
    return "http://${pick.address}:${pick.port}"
}
```

> 💡 **`?passing=true`** filtert bereits in Consul auf gesunde Instanzen. Ohne den Filter bekäme der Caller auch `warning`/`critical`-Instanzen mit zurück.

> 🎯 *Folie:* Pseudocode kurz zeigen, dann auf den echten Code in `services/shared/consul/register.go` und `services/shared/consul/resolver.go` verweisen — exakt diese Logik in ~80 Zeilen Go.

---

## 5. Server-Side vs. Client-Side Discovery

Das ist der **wichtigste pädagogische Übergang** der Story 2. Wer trifft die Entscheidung „welche Instanz wird angerufen"?

### Client-Side Discovery (Workshop-Variante)

Der Aufrufer (Booking) fragt Consul direkt und wählt die Instanz aus.

```
Booking ──▶ Consul: "Wer ist `flight-service`?"
Booking ◀── Consul: [flight-1, flight-2, flight-3]
Booking ──▶ flight-2  (zufällig gewählt)
```

**Pro:**
- Kein zusätzlicher Netzwerk-Hop.
- Lastverteilungs-Logik liegt im Caller (Random, Round-Robin, gewichtet, latency-aware …).

**Contra:**
- Jeder Caller muss Discovery-Code haben — pro Sprache/Stack.
- Caching-Strategie ist Aufgabe des Callers.

### Server-Side Discovery

Ein Reverse-Proxy / API-Gateway trifft die Entscheidung. Der Aufrufer ruft eine stabile virtuelle URL.

```
Booking ──▶ Gateway/Mesh: "ich will an `flight-service`"
                  │
                  └──▶ Consul: "Wer ist gesund?"
                  └──▶ flight-2
```

**Pro:**
- Caller weiß nichts von Discovery — einfache Codepfade.
- Zentrale Konfiguration (Retries, Circuit Breaker, mTLS) im Gateway.

**Contra:**
- Zusätzlicher Hop, weiterer Single-Point-of-Failure (zumindest logisch).
- Komplexere Infrastruktur.

> 💡 Im Workshop-Setup macht **Traefik** Server-Side Routing für externen Traffic (Ports 80/8080), aber **mit statischer File-Konfiguration** — nicht mit dem Consul-Provider. Der Booking-Service macht **Client-Side Discovery** für interne Calls. Beide Welten parallel — gutes Beispiel.
>
> Vertiefung dazu: siehe `services/load-balancing.md`.

---

## 6. Selbst implementieren oder Library?

| Aspekt | Selbst (Workshop) | Library |
|---|---|---|
| Lerneffekt | ⭐⭐⭐⭐⭐ Man versteht, was die Library tut | ⭐⭐ Bleibt Black Box |
| Korrektheit | Subtile Bugs lauern (Retry-Storm, Cache-Invalidation, Timeout-Pfusch) | Erprobt, viele Edge-Cases abgedeckt |
| Features | Was du baust | Watches, KV-Store, Sessions, ACLs, Streaming-Updates, Connection-Pooling |
| Zeilen Code | ~80 (Go) | 0 + ein paar Zeilen Konfiguration |

**Empfehlung für die Praxis:** Library nehmen.

**Empfehlung für den Workshop:** Einmal selbst bauen (oder den Reference-Code lesen), damit die Mechanik klar ist. Danach im echten Projekt auf Bewährtes setzen.

### Standardlibraries je Stack

| Stack | Library |
|---|---|
| Go | `hashicorp/consul/api` (offizieller Client) |
| Java (framework-frei) | `consul-api` (orbitz / ecwid) |
| Spring | **Spring Cloud Consul** (Auto-Registration via Annotation, Property-Refresh aus KV) |
| Quarkus | `quarkus-smallrye-stork` mit Consul-Service-Discovery |
| .NET | `Consul.NET` |
| Python | `python-consul2` |
| Node.js | `consul` (npm) |

> 💬 *Diskussionsfrage an die Teilnehmer:* „Welche von denen nutzt ihr in eurem Tagesgeschäft? Habt ihr schon mal Consul gegen Eureka / etcd / K8s-native getauscht? Was war der Treiber?"

---

## 7. Ausblick: DNS-Discovery, Kubernetes, Service Mesh

Consul ist eine Möglichkeit. Es gibt mehr.

### Consul-DNS-Interface

Consul beantwortet auch DNS-Queries: `flight-service.service.consul` → A-Records aller gesunden Instanzen. Für Sprachen ohne Consul-Client kann man so Discovery „kostenlos" haben — die Sprache braucht nur einen DNS-Resolver, den hat sie eh.

**Nachteil:** DNS-Caching (TTL) verzögert Failover; SRV-Records sind nicht überall gut unterstützt.

### Kubernetes-native Service Discovery

In K8s ist Service Discovery **eingebaut**: ein `Service`-Objekt bekommt einen DNS-Namen (`flight-service.default.svc.cluster.local`) und routet auf die Pods. Health-Checks via `readinessProbe`. Kein separater Consul nötig — der Cluster ist die Registry.

> 🔭 *Im Workshop:* Erwähnen, nicht vertiefen. „Wenn ihr in K8s seid, gilt vieles aus dieser Story implizit. Aber das Verständnis von Registry / Health / Resolver hilft, K8s zu durchschauen."

### Service Mesh (Istio, Linkerd, Consul Connect)

Discovery wird komplett in einen Sidecar-Proxy ausgelagert. Anwendungscode ruft einfach `localhost:port` — der Proxy weiß alles. Plus: mTLS, Circuit Breaker, Retries, Tracing zentral.

**Tradeoff:** Operativ teuer; nicht jedes Projekt braucht das.

---

## 8. Stolpersteine

> ⚠️ **„Service ist registriert, aber nicht bereit"**
> Health-Check liefert 200, aber die DB-Verbindung ist noch nicht aufgebaut. Caller bekommt 500. Lösung: `/health` echt implementieren — *Liveness* (Prozess lebt) und *Readiness* (Prozess kann Anfragen bedienen) trennen, nur Readiness an Consul melden.

> ⚠️ **„Zombie-Instanzen"**
> Container wird mit `kill -9` beendet, Deregister wird übersprungen. Eintrag bleibt minutenlang in Consul, wird als „critical" markiert, aber nicht entfernt. Lösung: `Check.DeregisterCriticalServiceAfter: "1m"` mitschicken — Consul räumt selbst auf.

> ⚠️ **„Resolver cached zu lange"**
> Naiver Cache („einmal aufgelöst, behalten") verzögert Failover dramatisch. Workshop-Code löst pro Request frisch auf — einfach und korrekt, aber bei hohem Traffic Last auf Consul. Pragmatischer Mittelweg: Cache mit kurzer TTL (1–5s) + Consul-Watches für Echtzeit-Updates.

> ⚠️ **„Hostname vs. IP — wer kann das auflösen?"**
> Im Compose-Netz ist `flight-1` ein gültiger DNS-Name, vom Host nicht. Wenn der Booking-Service außerhalb des Compose-Netzes läuft (etwa direkt aus der IDE), bekommt er Adressen, die er nicht erreichen kann. Lösung im Workshop: alles im selben Compose-Netz halten. Lösung in Produktion: registrierte Adresse muss aus Sicht der Caller routbar sein.

---

## 9. Diskussionsfragen für den Workshop

Zum Abschluss, wenn Zeit bleibt:

1. **Was, wenn Consul selbst ausfällt?** Single-Point-of-Failure! Antwort: Consul-Cluster (3 oder 5 Server), Caller-seitiges Caching für Übergangszeiträume.
2. **Wie testet ihr Discovery?** Service skalieren, Container killen, schauen ob das System weiterläuft. Chaos Engineering — vergleiche Story 3 / Dashboard.
3. **Client-Side vs. Server-Side: was passt zu eurem Stack?** Spring Cloud → Client-Side ist trivial. Polyglot mit Mesh-Vorhandensein → Server-Side. Kubernetes → eingebaut.
4. **Wie lange darf eine tote Instanz in der Registry stehen?** Tradeoff `Check.Interval` (Last) ↔ Failover-Geschwindigkeit ↔ `DeregisterCriticalServiceAfter` (Aufräumen).
5. **Service-Discovery-Daten persistieren?** Consul macht das per Default in Raft. Im Workshop läuft Consul aber im Dev-Mode (`agent -dev`) — alles im RAM. Was heißt das? Restart = leer. Diskussion: Wann ist das ein Problem?
6. **Selbst-Registrierung vs. Drittes-Tool-registriert** Im Workshop registriert sich der Service selbst (Self-Registration). Alternative: Sidecar / Orchestrator (z.B. K8s) macht das — Service weiß nichts von Consul (Third-Party-Registration). Vor- und Nachteile.

---

## 10. Wo der Code liegt (Reference-Implementierung)

| Was | Pfad |
|---|---|
| Register/Deregister + Retry/Backoff | `services/shared/consul/register.go` |
| Resolver (Random Load Balancing) | `services/shared/consul/resolver.go` |
| Lifecycle (Register beim Start, Deregister vor Server-Shutdown, Signal-Handling) | `services/flight/main.go` (analog `hotel/main.go`, `car/main.go`) |
| Resolver-Verwendung im Caller | `services/booking/story2/main.go`, `services/booking/story2/handler/booking.go` |
| Consul-Container im Compose-Setup | `services/docker-compose.infra.yml` (Service `consul`, Dev-Mode, UI auf Port 8500) |
| Begleitend: Zusammenspiel Traefik / Docker-DNS / Consul | `services/load-balancing.md` |
| Story-Beschreibung (was Teilnehmer bauen) | `docs/stories/story-02-service-discovery.md` |
