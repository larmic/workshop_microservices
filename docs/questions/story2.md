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

## Frage 3 — Wir machen Client-Side Load Balancing. Warum gibt es dann überhaupt Server-Side LB und Service Mesh?

**Frage:** Unser Resolver wählt zufällig eine gesunde Instanz. Das ist
einfach und funktioniert. Warum baut die halbe Branche stattdessen
Service Mesh mit Envoy / Istio / Linkerd?

**Antwort:** Weil Client-Side LB **mehrere unangenehme Eigenschaften**
hat, die in der Demo nicht auffallen:

1. **Logik wird in jedem Sprach-Stack einzeln gebaut.** Go, Java,
   Python, Node — jeder Client braucht den gleichen LB-Code. In Story 4
   bauen wir Bulkhead und CB ebenfalls im Service. Bei 5 Sprachen sind
   das 5 Implementierungen, die divergieren.
2. **Kein zentrales Tuning.** Wenn ihr das Retry-Verhalten ändern
   wollt, müsst ihr alle Services rebuilden und redeployen.
3. **Beobachtbarkeit ist verteilt.** Jeder Client hat seine eigenen
   Metriken. Ein einheitliches Bild über die Plattform fehlt.
4. **Sicherheit (mTLS).** Jeder Client muss Zertifikate handhaben.

Ein Service Mesh nimmt all das in den **Sidecar-Proxy** und macht es
sprachunabhängig:

```
Client-Side LB (unser Workshop):
   App ──► (eigener LB-Code) ──► Backend-Instanz

Service Mesh:
   App ──► Sidecar (Envoy) ──► Sidecar (Envoy) ──► App
            ↑ macht LB,         ↑ macht LB,
              CB, Retry,          CB, Retry,
              mTLS, Tracing       mTLS, Tracing
```

**Spicy Take-away:** Client-Side LB ist die richtige Lösung, **wenn
ihr eine homogene Sprachlandschaft und wenig Services habt**. Sobald
ihr in die Größenordnung „mehrere Teams, mehrere Sprachen, viele
Services" kommt, ist Service Mesh günstiger — die initialen Kosten
(Operations, Komplexität) sind aber hoch. Das ist eine **bewusste
Architektur-Entscheidung, kein No-Brainer in beide Richtungen**.

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
   davon nichts. → Brücke zu Story 8 (Downtimeless Deployment).

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
