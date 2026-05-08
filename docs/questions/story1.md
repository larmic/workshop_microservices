# Workshop-Fragen: Cloud-Native Booking-Service (Story 1)

Provokante Fragen für die Workshop-Diskussion. Ziel: Twelve-Factor und
Health-Checks **nicht als Checkliste abhaken**, sondern hinterfragen,
was die Patterns eigentlich versprechen — und wo sie in der Praxis
brüchig werden.

---

## Frage 1 — Unser Health-Check gibt 200 zurück. Heißt das, der Service ist gesund?

**Frage:** Wenn `/health` HTTP 200 zurückgibt, ist der Service dann
wirklich „gesund"? Was sagt das überhaupt aus?

**Antwort:** Praktisch nichts. Unser `/health` sagt nur: "der HTTP-Server
läuft und kann eine Route bedienen." Er sagt **nicht**:

- ob die Backend-Services Flight/Hotel/Car erreichbar sind,
- ob der Service überhaupt arbeitsfähig ist (z.B. DB-Connection),
- ob er nur noch 5 % freie Threads hat und gleich kippt.

In der Cloud-Welt unterscheidet man zwei Arten:

| Check         | Frage                          | Konsequenz bei Fail                |
|---------------|--------------------------------|------------------------------------|
| **Liveness**  | "Lebt der Prozess noch?"       | Container neu starten              |
| **Readiness** | "Kann der Prozess Last sehen?" | Aus dem Load-Balancer rausnehmen   |

Unser `/health` ist eher Liveness. Eine Readiness müsste z.B. prüfen,
ob Consul erreichbar ist (Story 2), ob die Backends antworten (Story 3+),
usw.

**Spicy Take-away:** Ein zu „cleverer" Readiness-Check kann dafür
sorgen, dass alle Instanzen gleichzeitig aus dem LB fliegen, sobald ein
Backend wackelt — und dann ist die ganze Plattform offline, obwohl die
Services selbst noch funktionieren würden. Health-Checks brauchen ein
durchdachtes Failure-Modell, sie sind kein Auto-Pilot.

---

## Frage 2 — Twelve-Factor sagt „ein Build, viele Umgebungen". Wie viele ENV-Variablen sind dafür realistisch?

**Frage:** Wir konfigurieren Backend-URLs per ENV. In der echten Welt
hat ein Service schnell 30+ ENV-Variablen. Ist das noch Twelve-Factor —
oder ist das die neue `application-prod.properties`, nur in YAML?

**Antwort:** Das Pattern selbst ist nicht das Problem, der
**Konfigurationsumfang** ist es. Drei Beobachtungen:

1. **ENV-Wildwuchs entsteht graduell.** Jede neue Integration bringt
   3–5 neue Variablen. Nach zwei Jahren ist die Helm-Chart eine
   200-Zeilen-Wand und niemand weiß mehr, was wirklich zur Laufzeit
   genutzt wird.
2. **Twelve-Factor sagt nichts über Sekret-Handling.** ENV ist für
   Secrets okay, aber nur, wenn die Plattform sie schützt
   (Vault/Secrets-Manager). Sonst landen sie im `docker inspect` und
   im Log.
3. **„Ein Build, viele Umgebungen" funktioniert nur, wenn die
   Umgebungen wirklich gleich sind.** Sobald Test eine andere
   Authentifizierung hat als Prod, ist das Versprechen schon gebrochen.

**Spicy Take-away:** Das Ziel ist nicht „möglichst viel über ENV
konfigurieren", sondern **„keine umgebungsspezifischen Werte im
Container-Image"**. Wenn dein Service dasselbe Verhalten in allen
Umgebungen hat und nur Endpoints/Credentials per ENV bekommt, bist du
im Geist von Twelve-Factor. Wenn dein Service über ENV
**Feature-Flags** umstellt, hast du eigentlich mehrere Services in
einem.

---

## Frage 3 — Logging auf stdout — was passiert, wenn 50 Service-Instanzen das gleichzeitig machen?

**Frage:** Twelve-Factor sagt: Logs als Stream auf stdout, der
Prozess kümmert sich nicht ums Routing. Klingt sauber. Was passiert
real, wenn 50 Pods das gleichzeitig produzieren?

**Antwort:** Die Verantwortung wandert nach oben — und damit auch das
Problem.

```
   Service-Instanz  ──stdout──►  Container-Runtime
                                 (Docker/Kubernetes)
                                       │
                                       ▼
                              Log-Forwarder (Fluentd/Vector/Filebeat)
                                       │
                                       ▼
                              Log-Aggregator (Loki, ELK, Datadog, …)
```

Auf jeder Stufe gibt es Stolperfallen:

- **Container-Logs sind round-rotated.** Wer zu spät schaut, sieht
  nichts mehr. Verteilte Bugs werden schwierig.
- **Strukturiertes JSON oder Plain Text?** Plain Text ist menschlich
  lesbar, aber im Aggregator schwer durchsuchbar. JSON ist umgekehrt.
- **Trace-IDs.** Ohne sie ist der Stream aus 50 Pods Rauschen — du
  siehst Hotel-Failures, aber nicht zu welchem Booking-Request sie
  gehörten. → Brücke zu OpenTelemetry / Distributed Tracing
  (späteres Workshop-Thema oder Bonus).
- **Kosten.** Cloud-Logging ist pro GB/Monat teuer. Ein zu lautes
  `log.Println` in einer heißen Schleife kann eine vierstellige
  Rechnung pro Monat produzieren.

**Spicy Take-away:** „Logs auf stdout" ist die richtige Default-Wahl —
aber das Pattern verschiebt nur die Komplexität. Wer das Aggregations-
Layer **nicht** durchdacht hat, fliegt damit nach drei Monaten auf die
Nase.

---

## Frage 4 — Der `/info`-Endpoint zeigt die aktuelle Config. Was könnte daran problematisch sein?

**Frage:** Wir haben einen `/info` gebaut, der Service-Konfiguration
zurückgibt. Was, wenn das aus Versehen im Internet erreichbar ist?

**Antwort:** Die ehrliche Antwort: **das ist ein klassisches
Information-Disclosure-Problem.** Schon harmlose Felder können einem
Angreifer das Leben leichter machen:

- `consulUrl: http://consul:8500` → "Aha, ihr habt Consul, vermutlich
  ohne ACL, gleicher Cluster, freue mich auf den Pivot."
- `timeout: 3000` → "Eure Synchron-Calls sind kurz, ich kann mit
  langsamen Antworten den Service trivial DoSen."
- Spring-Boot-Actuator-Bingo: viele Teams stellen `/actuator/env`
  versehentlich offen, inkl. DB-Passwörtern.

**Take-away:** Endpoints wie `/health`, `/info`, `/admin/...` brauchen
eine Antwort auf "wer darf das eigentlich aufrufen?". Mögliche
Strategien:

- Auf separatem Port / Mgmt-Interface, nicht im öffentlichen Routing.
- Hinter Authentifizierung.
- Im Image enthalten, aber per ENV deaktivierbar (`MGMT_ENDPOINTS=off`).

In unserem Workshop-Setup ist alles offen — bewusst, weil es das
Debugging einfach hält. **In Produktion wäre das fahrlässig.**

---

## Frage 5 — „Kein Error Handling, darf fehlschlagen" steht in der Story. Wie ehrlich ist das wirklich?

**Frage:** Story 1 sagt explizit: kein Error-Handling, der Aufruf
darf fehlschlagen. Heißt das, in der ersten Iteration eines Service
sollte man auf Resilienz verzichten?

**Antwort:** Ja — und genau das wird gerne falsch verstanden. Zwei
Aussagen, die nebeneinander stehen müssen:

1. **In der ersten Iteration ist „happy path" das Richtige.** Wenn man
   früh anfängt, jede Failure-Mode zu antizipieren, baut man ein Monster
   und liefert nichts. Story 1 zeigt **das Skelett**.
2. **Aber: das ist NICHT der Endzustand.** Wer in Produktion einen
   Service ohne CB / Bulkhead / Retry / Timeout deployed, der drei
   Backends sequenziell aufruft, hat ein **multiplikatives
   Verfügbarkeitsproblem**: bei 99 % Uptime pro Backend liegt die
   aggregierte Verfügbarkeit bei 0,99³ = 97 %. Das sind 21 Tage
   Downtime im Jahr.

```
Backend-Verfügbarkeiten kombinieren sich multiplikativ:

   Flight 99 %  ┐
   Hotel  99 %  ├─► Aggregat: 0.99 × 0.99 × 0.99 ≈ 97.0 %
   Car    99 %  ┘                                ~21 Tage Downtime/Jahr
```

**Spicy Take-away:** „Kein Error-Handling" ist eine **didaktische**
Reduktion, keine Architektur-Empfehlung. Der Workshop führt euch in
Story 3 / 4 / 5 schrittweise an die Resilienz-Patterns heran — aber
in Produktion fehlt euch ohne diese Patterns ein zentraler Teil eures
Service.

---

## Sammelthemen für die Diskussion

- Wie viele Konfigurationswerte hat euer schlimmster Service in der
  Praxis? Wie viele davon nutzt ihr noch?
- Wer schreibt bei euch die Health-Checks — Entwickler oder Ops? Was
  passiert dadurch (oder nicht)?
- Wenn ihr `/info` auf einem produktiven System aufruft: was sehe ich
  da, was ich nicht sehen sollte?
