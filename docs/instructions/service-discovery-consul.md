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

## 4. Die drei Abläufe in Regeln

Drei Lifecycle-Schritte: **Registrieren** beim Start, **Deregistrieren** beim Shutdown, **Auflösen** bei jedem Outbound-Call.

### A) Registrieren beim Start

1. Registrierungs-Daten zusammenstellen:
   - **Name** — logischer Service-Name (z.B. `flight-service`).
   - **ID** — eindeutig pro Instanz, im Workshop `{Name}-{Hostname}`.
   - **Address** und **Port** — wie der Service erreichbar ist (im Compose-Netz: Container-Hostname + interner Port).
   - **Health-Check** — URL (`http://{Address}:{Port}/health`) und Intervall (`10s`).
2. Daten als JSON per `PUT` an `{consulURL}/v1/agent/service/register` schicken.
3. Bei Fehlschlag: erneut versuchen, mit wachsender Wartezeit:
   - 1. Versuch sofort, dann Wartezeit beginnt bei 1 Sekunde.
   - Nach jedem fehlgeschlagenen Versuch die Wartezeit verdoppeln (1s → 2s → 4s → 8s → 16s).
   - Nach maximal 5 Versuchen aufgeben und mit Fehler abbrechen.
4. Bei Erfolg: die vergebene Service-ID merken — sie wird beim Shutdown wieder gebraucht.

> 💡 **Warum Retry?** Bei `docker compose up` startet alles parallel. Es kann sein, dass der Service schneller hochfährt als Consul. Ohne Retry wäre die Race ein flaky Workshop-Setup. Echte Workshop-Frage: „Was wäre die Alternative?" — Antworten: `depends_on` mit Healthcheck, Init-Container, Sidecar.

### B) Deregistrieren beim Shutdown

Beim Empfang von `SIGINT` / `SIGTERM` (z.B. `docker stop`):

1. **Zuerst** aus der Registry austragen — `PUT` an `{consulURL}/v1/agent/service/deregister/{serviceID}`.
2. **Erst danach** den HTTP-Server graceful herunterfahren (im Workshop mit 5s Timeout für laufende Requests).

> ⚠️ **Reihenfolge ist entscheidend.** Erst aus Consul raus, *dann* den Server schließen. Sonst:
> - Caller A fragt Consul: „Wer ist gesund?" → bekommt die sterbende Instanz, weil sie noch registriert ist.
> - Server ist aber schon mitten im Shutdown.
> - Caller A bekommt Connection Refused.
>
> 🎯 *Folie:* Den Sequenz-Diagramm-Vergleich zeigen — „falsche Reihenfolge" vs. „richtige Reihenfolge".

> 💡 **Was bei `kill -9`?** Dann läuft Schritt 1 nicht. Der Eintrag bleibt in Consul, bis der nächste Health-Check fehlschlägt — also bis zu `Check.Interval` lang sieht jemand eine tote Instanz als „gesund". Lösung: `DeregisterCriticalServiceAfter` setzen, damit Consul selbst aufräumt.

### C) Auflösen mit Client-Side Load Balancing

Bei jedem Outbound-Call zu einem Backend-Service:

1. `GET` an `{consulURL}/v1/health/service/{name}?passing=true`.
2. Antwort enthält eine Liste gesunder Instanzen — jede mit `Address` und `Port`.
3. Liste leer? → Fehler werfen, der Aufrufer entscheidet (Fallback / Circuit Breaker — siehe Story 3).
4. Liste nicht leer? → eine Instanz auswählen (im Workshop: zufällig aus der Liste) und URL `http://{Address}:{Port}` zurückgeben.
5. HTTP-Call direkt an die gewählte URL — **nicht** über Consul.

> 💡 **Warum `?passing=true`?** Ohne den Filter bekommt der Caller auch `warning`/`critical`-Instanzen — und routet potentiell auf eine kaputte. Mit dem Filter macht Consul die Vorauswahl.

> 🎯 *Folie:* Diese Regelliste zeigen, dann auf den echten Code in `services/shared/consul/register.go` und `services/shared/consul/resolver.go` verweisen — exakt diese Logik in ~80 Zeilen Go.

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
