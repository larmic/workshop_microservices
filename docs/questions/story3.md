# Workshop-Fragen: Circuit Breaker (Story 3)

Provokante Fragen rund um den Circuit Breaker. Ziel: das Pattern nicht
als magische Resilienz-Box akzeptieren, sondern verstehen, **welche
Annahmen** drin stecken — und wo sie kippen.

---

## Frage 1 — Wir öffnen den Circuit nach 5 aufeinanderfolgenden Fehlern. Wo kommt diese Zahl her?

**Frage:** Akzeptanzkriterium der Story sagt „nach 5 aufeinanderfolgenden
Fehlern". Warum 5? Wäre 3 nicht sicherer, oder 20 stabiler?

**Antwort:** Die ehrliche Antwort: **niemand weiß es genau, ohne die
Last und das Backend zu kennen.** 5 ist eine sinnvolle Default-Zahl,
aber sie steht im Spannungsfeld zwischen:

- **Zu klein (z.B. 1):** Ein einzelner zufälliger Fehler (Netz-Glitch,
  GC-Pause, transientes 5xx) öffnet den Circuit. Wir verlieren das
  Backend für 30 s, obwohl alles in Ordnung ist. → **False positives**.
- **Zu groß (z.B. 50):** Bei 100 req/s und einem komplett toten Backend
  brennen wir 50 Fehler-Calls in <1 s durch, jeder mit 3 s Timeout — der
  Booking-Service hängt parallel an 50 Threads, bevor der CB überhaupt
  reagiert. → **Bulkhead lässt grüßen** (Story 4).

```
Schwellenwert vs. Reaktionszeit / Robustheit:

  klein  ──┬───────────────────────────┬──  groß
           │                           │
       reaktiv,                    robust,
       aber flackernd              aber langsam
       bei Netz-Jitter             bei echten Ausfällen
```

**Andere Implementierungen** (Resilience4j) nutzen statt „N
aufeinanderfolgende Fehler" eine **Failure-Rate über ein Sliding
Window** (z.B. „>50 % der letzten 100 Calls"). Das ist robuster gegen
intermittierende Fehler — aber komplexer und braucht Last, um überhaupt
zu reagieren.

**Spicy Take-away:** Die 5 ist eine **Designentscheidung**, kein
Naturgesetz. In Produktion gehört der Wert ans aktuelle Last- und
Failure-Profil getuned, nicht ins erste Code-Beispiel kopiert.

---

## Frage 2 — Im HALF_OPEN-State lassen wir genau einen Probe-Call durch. Warum nicht alle?

**Frage:** Nach 30 s öffnet der CB für einen einzigen Test-Call. Warum
nicht einfach wieder normal aufmachen und schauen, was passiert?

**Antwort:** Weil ein **Probe-Storm** das Backend gleich wieder
umlegen würde. Szenario:

```
t=0    Backend kippt unter 1000 req/s
t=0    CB OPEN → 30 s Pause
t=30s  CB HALF_OPEN
       │
       │  ohne Probe-Lock:
       │   ─►  alle 1000 wartenden Requests laufen los
       │   ─►  Backend-Wiederanlauf wird sofort wieder erschlagen
       │
       │  mit Probe-Lock (atomare CompareAndSwap auf einen Slot):
       │   ─►  genau 1 Probe-Call erreicht das Backend
       │   ─►  alle anderen werden mit "in_flight"-Reject abgekürzt
       │   ─►  Backend hat Atemraum für die Antwort
       │   ─►  Erfolgreich? → CLOSED, Fehler? → OPEN für weitere 30 s
```

Im Code (`circuitbreaker.go`) macht das genau eine atomic-Bool:

```go
if cb.probeInFlight.CompareAndSwap(false, true) {
    // genau einer kommt rein
}
```

**Spicy Take-away:** Das HALF_OPEN-Detail ist der unterschätzte Teil
des Patterns. Eine naive Implementierung („ist Wartezeit abgelaufen?
Dann CLOSED") trägt zur **Cascading-Failure** bei, statt sie zu
verhindern. Genau hier trennt sich Library-Qualität von schnell
selbstgeschriebenem Code.

---

## Frage 3 — Was zählt überhaupt als „Fehler"? Eine 404 ist doch kein Backend-Problem.

**Frage:** Unser CB zählt jeden Fehler. Aber: ein `404 Not Found` ist
kein Defekt des Backends — der Aufrufer hat eine ID benutzt, die nicht
existiert. Sollte das den CB öffnen?

**Antwort:** **Nein**, und das ist eine häufige Falle. Eine grobe
Klassifizierung:

| HTTP-Status   | CB-relevant? | Begründung                                    |
|---------------|--------------|-----------------------------------------------|
| **5xx**       | ja           | Server-Fehler — Backend ist krank             |
| **Timeout**   | ja           | Backend hängt — symptomatisch für Überlast    |
| **Conn-Refused** | ja        | Backend nicht erreichbar                      |
| **4xx (außer 429)** | nein   | Client-Fehler — der Aufrufer ist schuld       |
| **429 Too Many** | jein      | Backend ist überlastet — diskutabel           |

In unserem Workshop-Code ist die Logik **bewusst grob**: alles, was
`error != nil` zurückgibt, zählt als Fehler. Im
`fetchJSON`/`postJSON`-Helper:

```go
if resp.StatusCode >= 400 {
    return nil, fmt.Errorf("backend returned %d: %s", ...)
}
```

Das heißt: aktuell **würde ein 404 den CB triggern**. Für die Demo
egal, in Produktion wäre das ein Bug.

**Spicy Take-away:** „Fehler" ist nicht binär. Wer den CB nur an
`error != nil` hängt, baut ihm ein zu sensibles Frühwarnsystem. Reale
Implementierungen (Resilience4j: `recordExceptions` /
`ignoreExceptions`) machen genau hier eine bewusste Klassifikation.

---

## Frage 4 — Bei OPEN liefern wir `flights: []` statt einer Fehlermeldung. Täuschen wir den User?

**Frage:** Wenn Flight ausfällt, sieht der User „keine Flüge verfügbar".
Tatsächlich ist nur unser Backend kaputt. Ist das ehrlich?

**Antwort:** Streng genommen: **nein, das ist ein UX-Antipattern.** Der
User kann nicht zwischen „es gibt keine Flüge auf dieser Strecke" und
„unser System hat ein Problem" unterscheiden — und macht im
Zweifelsfall eine schlechte Buchungs-Entscheidung („dann fahre ich
halt nur Hotel + Auto, ohne Flug").

Ehrlichere Varianten:

1. **Teilantwort mit Hinweis-Feld.** JSON enthält `flights: []` plus
   `unavailable: ["flight"]` oder `errors: [{ service: "flight",
   reason: "temporarily unavailable" }]`. UI kann eine
   Ersatznachricht anzeigen.
2. **Aussagekräftige Header.** Wir setzen schon `X-Circuit-Open:
   flight` — die UI könnte das auswerten und anzeigen.
3. **Stale-while-revalidate.** Letztes erfolgreiches Ergebnis cachen
   und im Fallback ausliefern, mit Hinweis „Daten von vor 2 Min".

Der Workshop-Code ist die einfachste Variante (leeres Array, keine
Erklärung) — bewusst, weil's um das Pattern geht, nicht um perfekte
UX. **In Produktion gehört das aufgehübscht.**

**Spicy Take-away:** Resilience-Patterns sind **nicht ehrlich von
selbst.** Sie liefern dem Aufrufer eine Antwort, die wie Erfolg
aussieht. Wer das nicht durch UX und Header transparent macht, baut
gut funktionierende Systeme, in denen der Mensch falsche Entscheidungen
trifft.

---

## Frage 5 — Wir warten 30 Sekunden, bevor wir's noch mal probieren. Wo kommt diese Zahl her — und ist sie sinnvoll?

**Frage:** Akzeptanzkriterium: 30 s in OPEN. Aber:
- wenn das Backend in 1 s wieder gesund ist, blockieren wir 29 s
  unnötig
- wenn es 5 min braucht, hämmern wir alle 30 s erneut drauf

Was ist die richtige Wartezeit?

**Antwort:** Es gibt keine. **30 s ist ein Kompromiss**, der für den
typischen Fall „Container restartet, Service ist in 10–20 s wieder
oben" funktioniert. Daneben gibt es zwei smarter:

- **Exponential Backoff:** Erste Wartezeit z.B. 5 s, dann 10 s, 20 s,
  40 s, 80 s, gedeckelt bei 5 min. Wenn das Backend tot bleibt,
  hämmern wir es nicht jede Minute. Sobald ein Probe erfolgreich ist,
  Reset auf 5 s.
- **Adaptive Wartezeit:** Wartezeit am letzten erfolgreichen
  Antwort-RTT orientieren — Backend, das normal in 50 ms antwortet,
  darf nach 1 s wieder probiert werden; eines mit 3 s RTT erst nach
  30 s.

```
   Konstante Wartezeit (Workshop):
       OPEN ── 30s ── HALF_OPEN ── fail ── OPEN ── 30s ── HALF_OPEN ──
       (gleich gefährlich, egal wie tot das Backend ist)

   Exponential Backoff:
       OPEN ── 5s ── HALF_OPEN ── fail ── OPEN ── 10s ── HALF_OPEN ── fail
       ── OPEN ── 20s ── … (das tote Backend bekommt zunehmend Ruhe)
```

**Spicy Take-away:** Eine fixe Wartezeit ist die **schlechteste der
guten Optionen**. Sie ist einfach zu implementieren und reicht für die
Demo. Echte Resilience-Libraries (Resilience4j, Polly) liefern
Backoff-Strategien out of the box — sich darauf nicht zu verlassen,
ist ein Symptom für „CB selbst gebaut".

---

## Frage 6 — Der CB-Zustand lebt im Speicher. Was passiert beim Restart des Booking-Service?

**Frage:** Wenn ich den Booking-Service neu starte, ist sein CB-Status
wieder CLOSED — auch wenn das Hotel-Backend gerade tot ist. Ist das
nicht gefährlich?

**Antwort:** Ja, in der Tendenz schon. Was passiert nach einem Restart:

```
t=0    Booking-Service startet (CB: CLOSED)
t=1    Erster Aufruf gegen totes Hotel → Fehler 1/5
t=2..5 vier weitere Aufrufe → Fehler 5/5 → CB OPEN
       (in der Zwischenzeit haben 5 Aufrufer 3 s gewartet)
t=…    Normaler Betrieb
```

In den ersten Sekunden nach Restart ist der Service also „blind" und
verbraucht 5 Aufrufe lang Threads, Connections und Wartezeit auf der
Aufruferseite, bevor er reagiert. Das skaliert übel:

- **Bei mehreren Replicas** macht jede für sich diese Lernphase.
- **Bei Rolling Deploy** sind temporär alle Replicas in der Lernphase
  gleichzeitig.

Lösungsansätze, von „simpel" bis „aufwendig":

1. **Längere Timeouts beim ersten Aufruf hinnehmen** und die paar
   Sekunden verschmerzen. Default unseres Workshops.
2. **Initial Probe beim Service-Start.** Vor dem ersten Real-Traffic
   einen Probe-Call gegen jedes Backend, damit der CB schon bei Start
   einen aktuellen Wert hat.
3. **CB-State in einer Shared Cache (Redis o.ä.).** Dann sehen alle
   Replicas die gleiche Sicht auf jedes Backend. Aufwand hoch,
   Konsistenz-Probleme inklusive.
4. **Service Mesh (z.B. Envoy).** Der Sidecar hält den CB-State, der
   App-Restart hat keinen Einfluss. → Brücke zu Story 7+.

**Spicy Take-away:** Lokale CB-Statistik ist **per Definition
ephemerer Zustand**. Wer das vergisst, hat nach jedem Deploy ein paar
Sekunden Blind-Phase, die im Monitoring als Latency-Spike auftaucht
und niemand kann erklären, woher.

---

## Frage 7 — Wir haben drei CBs (Flight, Hotel, Car). Warum nicht einen pro Endpoint? Oder einen für alles?

**Frage:** Granularität: pro Backend? Pro Endpoint? Pro Aufruf? Wo ist
der richtige Schnitt?

**Antwort:** Es gibt drei Stufen mit jeweils anderen Trade-offs:

| Granularität           | Beispiel                          | Vorteil                              | Nachteil                                |
|------------------------|-----------------------------------|--------------------------------------|-----------------------------------------|
| **Global (ein CB)**    | „backend"                         | trivial einfach                      | ein kranker Service kappt alle Backends |
| **Pro Service**        | Flight / Hotel / Car (unsere Wahl)| isoliert Ausfall pro Service         | innerhalb des Service kein Schutz       |
| **Pro Endpoint**       | `flight.search` / `flight.book`   | sehr feingranular                    | viele CBs, schwer überblickbar          |

```
   Global:                Pro Service:           Pro Endpoint:
   ┌──[CB]──┐             ┌─[CB]── Flight        ┌─[CB1]── Flight.search
   │        ├──Flight     │                      │
   │   App  ├──Hotel      App ─[CB]── Hotel      App ─[CB2]── Flight.book
   │        ├──Car        │                      │
   └────────┘             └─[CB]── Car           └─[CB3]── Hotel.search …
```

**Pro Service** ist die übliche Default-Wahl, weil ein „kranker"
Service oft alle seine Endpoints betrifft (Connection-Pool tot,
Container down). Pro Endpoint lohnt sich, wenn ein Service mehrere sehr
unterschiedliche Workloads hat (z.B. „lese" vs. „schreibe", einer
schnell, einer langsam).

**Spicy Take-away:** Granularität ist eine **Designentscheidung,
keine Pattern-Eigenschaft.** Der Default „pro Service" ist meist
richtig, aber wer einen langsamen `book`-Endpoint hat und einen
schnellen `search`, fasst die zusammen — und kappt search, weil book
unter Last steht.

---

## Sammelthemen für die Diskussion

- Welcher Schwellenwert ist bei euch konfiguriert (oder eben nicht)?
  Wer hat ihn gesetzt — und nach welcher Begründung?
- Habt ihr einen CB, der im OPEN-State hängen geblieben ist? Was war
  die Ursache — kaputtes Backend oder zu aggressives Tuning?
- Wer schaut sich bei euch CB-Metriken an? Sind sie im Standard-
  Dashboard, oder muss man sie suchen?
