# Workshop-Fragen: Service Discovery (Story 2)

Provokante Fragen rund um Consul, Self-Registration und Client-Side
Load Balancing. Die Story löst ein konkretes Problem (statische URLs
in der Config), führt aber **drei neue Probleme** ein — und genau
darüber sollte der Workshop sprechen.

---

## Frage 1 — Wir haben jetzt eine Service Registry. Ist das nicht ein neuer Single Point of Failure?

**Frage:** In Story 1 war die Backend-URL hartkodiert in der ENV. Doof,
aber zumindest unabhängig. Jetzt fragen wir bei jedem Aufruf Consul.
Was passiert, wenn Consul kippt — fällt dann unser ganzes System?

**Antwort:** Ja, **wenn man es naiv implementiert**. Genau so haben wir
es im Workshop gebaut: jeder Aufruf geht zuerst gegen Consul, und wenn
Consul nicht antwortet, fällt der ganze Booking-Service.

```
Workshop-Stand (naiv):
   /booking/offers ──► Consul (ResolveServiceURL)
                       │
                       ├─► OK ──► Backend
                       └─► Fehler ──► 500 zurück, kein Fallback
```

In Produktion baut man das **mehrstufig**:

1. **Lokaler Cache der Service-Locations.** Letzten erfolgreich
   aufgelösten Endpoint behalten und benutzen, wenn Consul kurz weg ist.
2. **Consul-Agent pro Node.** Statt mit dem zentralen Server zu reden,
   spricht jeder Service mit einem lokalen Agent. Der Agent puffert.
3. **TTL für Cache-Einträge.** Damit ausgefallene Instanzen nicht
   ewig im Cache stehen.

**Spicy Take-away:** Service Discovery löst das Problem statischer
URLs — und schafft sich dadurch eine neue zentrale Komponente, die
hochverfügbar sein muss. Wer Consul / Eureka / etcd „mal eben"
einführt, ohne über deren Resilienz nachzudenken, hat das Problem
nur verschoben, nicht gelöst.

---

## Frage 2 — Self-Registration setzt voraus, dass der Service seine eigene Adresse kennt. Ist das nicht ein Bruch von Twelve-Factor?

**Frage:** In Story 1 war's noch sauber: der Service kümmert sich um
sich selbst, externe Konfiguration kommt von außen. Jetzt soll der
Service plötzlich seine eigene IP/Port kennen und sich aktiv bei
Consul anmelden. Ist das nicht eine Regression?

**Antwort:** Halb. Es gibt zwei Self-Registration-Modelle:

| Modell                 | Wer registriert?            | Vorteil                        | Nachteil                       |
|------------------------|-----------------------------|--------------------------------|--------------------------------|
| **Im Service-Code**    | Der Service selbst (z.B. `register-on-startup`) | Funktioniert ohne Plattform | Service kennt seine Umgebung  |
| **Sidecar / Plattform**| Container-Runtime (z.B. K8s Service / Consul-Connect) | Service ist umgebungsneutral | Setzt Plattform voraus        |

Im Workshop nutzen wir Variante 1: Flight/Hotel/Car melden sich beim
Start aktiv an. Das ist okay für Demos, aber in einer Container-
Plattform wie Kubernetes ist Variante 2 (Sidecar) sauberer:

```
Ohne Sidecar (Workshop-Stil):
   ┌─────────────┐     "Hi Consul, ich bin
   │ Flight-App  │      flight-service auf 10.0.0.5:8080,
   └─────┬───────┘      hier mein Health-Check"
         │
         ▼
      Consul

Mit Sidecar (Plattform-Stil):
   ┌──────────────┐ ┌────────────┐
   │ Flight-App   │ │  Sidecar   │ ◄── überwacht App,
   │ kennt nur    │◄┤  (registry │     registriert sich,
   │ localhost    │ │   client)  │     deregistriert beim Stop
   └──────────────┘ └─────┬──────┘
                          ▼
                       Consul
```

**Spicy Take-away:** Self-Registration im Code ist die schnellere
Lösung, aber sie verteilt Plattform-Wissen in jeden Service. Das ist
in Ordnung für 5 Services, wird zur Hölle bei 50. Der Branchen-Trend
geht klar zu „Plattform übernimmt das" (Service Mesh, K8s Services).

---

## Frage 3 — Wir machen Client-Side Load Balancing. Warum gibt es dann überhaupt Service Mesh — oder brauchen wir am Ende beides?

**Frage:** Unser Resolver wählt zufällig eine gesunde Instanz. Das ist
einfach und funktioniert. Warum baut die halbe Branche stattdessen
Service Mesh mit Envoy / Istio / Linkerd? Und ist das überhaupt ein
Gegensatz?

**Antwort:** „Client-Side LB vs. Service Mesh" ist eine **falsche
Dichotomie**. Ein Service Mesh **macht** Client-Side LB — nur eben
im Sidecar statt im App-Prozess. Man muss zwei Achsen sauber trennen:

| Achse | Optionen |
|-------|----------|
| **Wer entscheidet?** (Topologie der Auswahl) | Client-Side (Aufrufer wählt) ↔ Server-Side (zentraler LB/Proxy davor wählt) |
| **Wo läuft der Code?** (Deployment) | Library im App-Prozess ↔ Sidecar-Prozess ↔ Plattform-Service |

Client-Side LB gibt es als Library (Spring Cloud LoadBalancer, Netflix
Ribbon, unser Workshop-Resolver), als reinen Sidecar (Envoy ohne
Mesh-Drumherum) und implizit in Plattformen (`kube-proxy` lokal pro
Knoten). „Client-Side" ist also eine **Entscheidungs-Topologie**, keine
**Deployment-Topologie**.

Damit entkräftet sich auch der naive Reflex „Library schlecht, Sidecar
gut, also Mesh". Sobald die LB-Logik in einem **reinen LB-Sidecar**
sitzt, ist das Sprach-Stack-Argument schon erledigt — dafür braucht
man noch kein Mesh.

### Was unterscheidet ein Mesh wirklich?

Service Mesh ist nicht „Sidecar statt Library", sondern **Sidecar
plus deutlich mehr als nur LB**:

1. **L7-Resilience eingebaut** — Retry, Timeout, Circuit Breaker,
   Outlier Detection, Rate Limit. Stories 3 / 4 / 5 bauen wir im
   Workshop in der App; ein Mesh nimmt das ab.
2. **mTLS by default** zwischen allen Services — Zero-Trust, ohne
   dass die App Zertifikate sieht.
3. **Zentrale Control Plane** — eine YAML/CRD-Änderung, alle Sidecars
   ziehen nach. Retry-Verhalten ändern = kein Redeploy.
4. **Traffic-Splitting** für Canary / Blue-Green / A-B — die Antwort
   auf Frage 5.
5. **Einheitliche Observability** — jeder Hop emittiert dieselben
   Metriken / Traces aus demselben Layer.

### Brauchen wir beides?

Das ist keine Oder-Frage: **Mesh enthält Service Discovery.** Die
Control Plane muss wissen, welche Endpoints es gibt — in K8s übernimmt
das der API-Server, mit Consul Connect übernimmt das Consul selbst,
in standalone Envoy ein xDS-Server.

Die ehrliche Schichtung:

| Stufe | Was steckt drin | Wer betreibt |
|-------|-----------------|--------------|
| **1. Nur Discovery** (Workshop) | Registry + Resolver-Library | App-Team |
| **2. Discovery + LB-Sidecar** | Wie 1, aber Auflösung im Sidecar | App-Team + Plattform |
| **3. Service Mesh** | Discovery + LB + Retry/CB/Timeout + mTLS + Routing + Observability | Plattform-Team |

Stufe 2 ist ein oft übersehener Mittelweg — sinnvoll bei
polyglotter Landschaft ohne den vollen Mesh-Betrieb.

```
Stufe 1 (Workshop):
   App ──► (eigene Resolver-Library) ──► Backend

Stufe 2 (LB-Sidecar):
   App ──► Sidecar (Envoy als LB) ──► Backend
            ↑ nur LB + Discovery

Stufe 3 (Mesh):
   App ──► Sidecar (Envoy) ──► Sidecar (Envoy) ──► App
            ↑ LB, CB, Retry,    ↑ LB, CB, Retry,
              mTLS, Tracing,      mTLS, Tracing,
              Authz, Routing      Authz, Routing
   ▲ alle gesteuert von einer zentralen Control Plane
```

**Spicy Take-away:** Die echte Frage ist nicht „Mesh ja/nein", sondern
„**brauchen wir mTLS, Traffic-Splitting und L7-Resilience wirklich
überall, oder reicht Stufe 1 oder 2?**". Wer die Frage nicht stellt
und direkt zu Istio greift, hat sechs Monate später einen Sidecar-
Wildwuchs und kein Team, das den Mesh sauber operiert. Linkerd ist
schlanker, Consul Connect die Variante für Nicht-K8s-Welten — und ein
LB-Sidecar ohne Mesh ist ein legitimer Mittelweg.

---

## Frage 4 — Consul prüft mit einem Health-Check, ob der Service gesund ist. Was, wenn der Service lügt?

**Frage:** Unser `/health` aus Story 1 gibt einfach 200 zurück. Consul
nimmt das als Beweis, dass der Service gesund ist. Was passiert, wenn
der Service kaputt ist, aber `/health` trotzdem 200 sagt?

**Antwort:** Dann routet Consul munter Traffic auf eine kaputte
Instanz. Consul ist **keine Wahrheits-Instanz**, sondern nur ein
Endpoint-Indexer. Drei typische Failure-Modes:

1. **Zombie-Service.** HTTP-Server lebt, Worker-Pool ist tot, jeder
   Request hängt 30 s. Health-Check sagt OK, weil er nur einen
   leichtgewichtigen Endpoint trifft.
2. **Backend-Abhängigkeit weg.** Service kann technisch antworten,
   aber er kommt nicht mehr an seine Datenbank. Health-Check würde
   das nur erkennen, wenn er die DB mitprüft.
3. **Langsame Degradation.** P99-Latenz steigt von 50 ms auf 5 s,
   Service ist „technisch" gesund, faktisch unbrauchbar.

In jedem Fall hilft euch Consul nicht — es muss der **richtige
Check** angeschlossen sein. Optionen, in steigender Schärfe:

- TCP-Ping (lebt der Port?) → schwächster Check
- HTTP-Get auf `/health` → unser Workshop-Default
- HTTP-Get mit echter Probe-Logik (Schema-Read auf DB, etc.)
- Externe Synthetic-Checks (echter Booking-Roundtrip alle 30 s)

**Spicy Take-away:** Service Discovery ist **so gut wie der
schlechteste Health-Check, der ihr darin pflegt**. Ein Flight-Service,
der lügt, ist schlimmer als gar keine Service Registry — weil ihr
euch in Sicherheit wiegt.

---

## Frage 5 — Mit logischen Namen wie `flight-service` — wie deploye ich eine neue Version, ohne dass alle Anfragen sofort drauf gehen?

**Frage:** Im Booking-Service steht „löse `flight-service` auf". Wenn
ich jetzt eine neue Version v2 deployen will und erst mal nur 5 % des
Traffics darauf schicken möchte (Canary) — wie?

**Antwort:** Mit unserem aktuellen Stand: **gar nicht.** Wir haben
einen logischen Namen und eine zufällige Auswahl unter allen gesunden
Instanzen. Sobald ich `flight-service-v2` zusätzlich registriere,
bekommt es **sofort den gleichen Traffic-Anteil** wie v1.

Lösungsoptionen, jeweils mit Komplexität:

1. **Tags / Metadaten in Consul.** v1 trägt `version=1`, v2 trägt
   `version=2`. Resolver liest Tag und gewichtet. Funktioniert, aber
   muss in jedem Client gebaut werden.
2. **Separate logische Namen** (`flight-service-v1`, `flight-service-v2`).
   Sauberer, aber jetzt muss der Aufrufer wissen, welche Version er
   will — Versionierung leakt nach oben.
3. **Service Mesh / API-Gateway.** Traffic-Splitting wird zentral
   konfiguriert (z.B. `90 % → v1, 10 % → v2`), die Aufrufer wissen
   davon nichts. → Brücke zum Thema Downtimeless Deployment.

**Spicy Take-away:** Service Discovery löst „wo läuft das?", nicht
„welche Version will ich?". Sobald ihr Canary / Blue-Green / A/B-Tests
braucht, kommt eine zweite Schicht ins Spiel. Wer das nicht im Modell
hat, debuggt später Tage daran, warum 5 % der User „komische Fehler"
sehen.

---

## Frage 6 — Die Backend-Services registrieren sich beim Start in Consul. Was passiert beim STOP?

**Frage:** Beim Start melden sich Flight/Hotel/Car bei Consul an. Was
passiert beim Beenden? Verschwindet der Eintrag?

**Antwort:** Im Workshop-Code: **nein, nicht aktiv.** Wenn ein Container
einfach gestoppt wird (`docker stop`, `kubectl delete`), bleibt der
Eintrag in Consul stehen, bis dessen TTL abläuft oder der
Health-Check nach mehreren Misses ausschlägt.

Das hat ein konkretes Symptom:

```
t=0     Hotel-Instanz wird gestoppt
t=0     Eintrag in Consul ist noch DA, Health-Check läuft, aber
        Hotel antwortet nicht mehr
t=0..n  Booking-Service löst hotel-service auf, wählt zufällig die
        TOTE Instanz, Aufruf läuft in den Connection-Refused
t=n     Consul markiert Instanz nach n fehlgeschlagenen Checks als
        unhealthy → wird aus Resolver-Ergebnis entfernt
```

Bis dahin (typisch 10–30 s) läuft Traffic ins Leere. Lösungen:

1. **Graceful Shutdown** im Service: vor dem Stop aktiv bei Consul
   `deregister` rufen. Dann ist der Eintrag sofort weg.
2. **Schneller Health-Check-Intervall** in Consul (z.B. 1 s statt 10 s).
   Senkt Erkennungszeit, kostet Last.
3. **Out-of-Service-Mode**: Service nimmt für ein paar Sekunden keine
   neuen Requests mehr an, beendet laufende, dann erst Stop. Damit
   verlieren wir während des Deploys keine Anfragen.

**Spicy Take-away:** Service Discovery ist immer ein **eventually-
consistent** System. Es gibt **immer** ein Zeitfenster, in dem Aufrufer
auf tote Endpoints stoßen. Das ist genau einer der Gründe, warum wir
in Story 3 / 4 Resilience-Patterns brauchen — sie überbrücken dieses
Fenster, ohne dass der User es spürt.

---

## Sammelthemen für die Diskussion

- Welche Service Registry / Discovery nutzt ihr? Wie hochverfügbar ist
  die wirklich (Mehrheit reicht? RAFT? Multi-Region)?
- Wer ist bei euch verantwortlich, wenn ein Service in der Registry
  „fehlt" — App-Team oder Plattform-Team?
- Habt ihr schon mal erlebt, dass ein toter Service über Stunden in der
  Registry blieb? Was war die Ursache?
