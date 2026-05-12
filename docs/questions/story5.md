# Workshop-Fragen: Saga (Story 5)

Provokante Fragen rund um das Saga Pattern und die Frage, was passiert,
wenn das „alles oder nichts" *selbst* anfängt zu wackeln. Ziel: nicht
nur den Happy Path verstehen, sondern die Failure-Modi, die in echten
Systemen die schwierigen sind — und die Brücke zu Story 6
(Choreography-Saga via Events) bewusst sehen.

---

## Frage 1 — Was passiert, wenn die Kompensation selbst fehlschlägt?

**Frage:** Wir kompensieren mit `DELETE /bookings/{id}` gegen Hotel,
Flight, Car. Was, wenn dieser Aufruf seinerseits einen 5xx wirft oder
in einen Timeout läuft? Der Flug ist gebucht, das Hotel hat „nein"
gesagt, und jetzt kann ich den Flug nicht mehr stornieren — was nun?

**Antwort:** Saga macht eine **starke Annahme**: Forward-Steps dürfen
scheitern, Kompensationen **müssen letztlich gelingen**. Ohne diese
Annahme bricht das ganze Konstrukt zusammen. Die Lehrbuch-Antworten:

| Strategie | Was sie tut | Wann sie greift |
|---|---|---|
| **Idempotenz** | mehrfacher `DELETE` → immer dasselbe Ergebnis | bauliche Voraussetzung |
| **Retry mit Backoff** | transienten Fehler aussitzen | Service kurz weg, Netz-Wackler |
| **Persistenter Saga-Log** | Crash darf keine offene Kompensation verlieren | Booking-Service stürzt mitten in `COMPENSATING` ab |
| **Dead-Letter / Operator-Inbox** | Mensch greift ein | alle Retries erschöpft |
| **Pivot zur fachlichen Alternative** | Statt Storno → Gutschein | Flug ist schon abgehoben, „technische" Storno unmöglich |
| **Compensation-by-design** | Reservierung als Status, nicht als Löschung | macht „Kompensation kann nicht scheitern" zur Eigenschaft des Modells |

**Take-away:** Eine Saga ohne Plan für gescheiterte Kompensation ist
**keine Saga, sondern eine optimistische Hoffnung**. Die Frage „was bei
Misserfolg" ist nicht optional — sie ist das eigentliche Engineering an
dem Pattern.

---

## Frage 2 — Wir haben keinen Retry implementiert. Ist das schlimm?

**Frage:** In unserer Workshop-Implementierung ruft Booking die
Kompensation **genau einmal** auf. Schlägt sie fehl, ist die Saga
einfach `FAILED`. Verletzen wir damit das Pattern?

**Antwort:** Ja, formal verletzen wir das letzte Akzeptanzkriterium von
Story 5. Pragmatisch ist die Entscheidung trotzdem vertretbar — *wenn*
sie bewusst getroffen wird. Was wir damit aufgeben und wie wir es
absichern:

```
   ohne Retry
   ─────────────────────────────────────────────────────────
   booking ─DELETE→ hotel  ─►  500 Internal Server Error
                              │
                              └─►  Saga-Status = FAILED
                                   Flug bleibt fälschlich gebucht
                                   Niemand merkt es ohne Monitoring
```

Was *zwingend* dazugehört, wenn man Retry weglässt:

1. **Saga-Status persistieren** (`PENDING`/`COMPENSATING`/`FAILED`) —
   sonst weiß niemand, dass eine Saga überhaupt hängt.
2. **Alert auf Sagas im Status `FAILED` mit unfertiger Kompensation** —
   das ist der Operator-Eingriff.
3. **Idempotente Kompensations-Endpoints** — damit ein manueller
   Retry („Operator klickt nochmal") gefahrlos möglich ist.

**Take-away:** Retry weglassen ist erlaubt — aber dann muss das
**Monitoring der Retry sein**. Was du nicht im Code hast, musst du im
Dashboard haben. Was du in keinem von beiden hast, hast du nicht.

---

## Frage 3 — Wäre Eventing nicht die natürlichere Antwort als sync HTTP?

**Frage:** Booking macht einen synchronen `DELETE` gegen Hotel. Wenn
Hotel kurz weg ist, muss Booking das selbst auffangen (Retry-Schleife,
Timeout, Saga-State). Wäre es nicht viel sauberer, wenn Booking ein
Event `CancelBooking` rauslegt, Hotel das aus dem Broker zieht und
Booking damit „fertig" ist?

**Antwort:** Doch — und genau das ist der Sprung von **Orchestration
über sync HTTP** (Story 5) zu **Event-getriebener Saga** (Story 6).
Was du dabei gewinnst und was sich verschiebt:

```
   heute (sync Orchestration)
   ─────────────────────────────────────────────────────────
   booking ─DELETE──► hotel       Retry-Schleife in booking
   booking ◄──204─── hotel        Booking trägt die Verantwortung

   mit Eventing
   ─────────────────────────────────────────────────────────
   booking ─CancelBookingCommand─► [bus] ─► hotel
                                              │
   booking ◄─BookingCancelledEvent─[bus] ◄────┘
                                              │
                                   oder       └─CancellationFailedEvent
```

| Aspekt | sync HTTP | Eventing |
|---|---|---|
| Retry-Logik | im Booking-Code | im Broker (Redelivery) |
| Hotel down | Booking schlägt fehl | Event wartet in Queue |
| Booking-Crash mitten drin | offene Saga, Recovery aus DB | Events bleiben im Bus, Resume kostenfrei |
| Latenz | sofortige Antwort | eventually consistent |
| Topologie | Booking kennt Hotel-URL | Booking kennt nur Topic |
| Komplexität im Stack | gering | Broker, Outbox, DLQ |

**Aber** — der Punkt, der oft untergeht: Booking ist **nicht** fertig
nach „Event raus". Der Kunde will am Ende eine Antwort:
„Reise gebucht? storniert? steckt fest?". Also:

- Hotel publiziert `BookingCancelled` zurück.
- Booking konsumiert das Reply-Event und führt seinen Saga-Status nach.
- Booking braucht weiterhin **Timeout-Erkennung** (Reply nach X Minuten
  nicht da → eskalieren).

Das ist im Endeffekt **asynchrone Orchestrierung**, nicht „Booking ist
seine Verantwortung los". Die Verantwortung *für die Ausführung* wandert
zu Hotel; die Verantwortung *für den Gesamtstatus gegenüber dem Kunden*
bleibt bei Booking.

**Take-away:** Eventing ist die robustere Architektur — aber sie
verschiebt die Komplexität, sie eliminiert sie nicht. Story 5 zeigt
die Saga **isoliert** (sync, ein Konzept). Story 6 schaltet das
Eventing dazu und macht die Kompensation asynchron — das ist die
**Choreography-Saga**. Story 7 nimmt sich dann **CQRS** vor, und zwar
bewusst als **eigenes** Konzept: ein Read-Model, das die Event-
Infrastruktur aus Story 6 *nutzen kann*, aber nicht muss — denn CQRS
ist „Lese- und Schreibmodell trennen", nicht „Events haben". Drei
Schritte, drei Konzepte, sauber getrennt.

---

## Frage 4 — Warum liefert `DELETE /bookings/{id}` bei unbekannter ID 204, nicht 404?

**Frage:** Wir geben `204 No Content` auch zurück, wenn die ID nie
existiert hat. Verschleiert das nicht echte Fehler? Wäre `404` nicht
ehrlicher?

**Antwort:** **Nein** — und der Grund kommt direkt aus der Saga-Mechanik:

1. **Idempotenz ist ein Feature, nicht ein Schönheitsmakel.**
   Saga-Retries dürfen mehrfach denselben Storno absetzen. Jeder Aufruf
   nach dem ersten würde bei `404`-Antwort die Retry-Logik fälschlich
   als „Fehler" interpretieren — obwohl der Effekt längst eingetreten
   ist.
2. **Ohne State kann der Service ehrlich gar nicht zwischen
   „nie existiert", „bereits storniert" und „gerade erst gebucht und
   schon wieder weg" unterscheiden.** Die Workshop-Services persistieren
   nichts. `404` wäre also geraten, nicht gewusst.
3. **REST-konvention für DELETE ist umstritten, aber Saga-Praxis ist
   klar:** Kompensations-Endpoints sind „at-least-once safe". Status
   2xx auf jeden plausiblen Aufruf, 4xx nur bei strukturell falschen
   Anfragen (z.B. ID-Format invalid).

```
                  Aufrufer-Sicht
   ─────────────────────────────────────────────
   1. DELETE /bookings/F-7c1a9f   →  204   "ok"
   2. DELETE /bookings/F-7c1a9f   →  204   "auch ok, idempotent"
   3. DELETE /bookings/F-deadbe   →  204   "kannte ich nicht — egal"

                  Aufrufer-Logik
   ─────────────────────────────────────────────
   if status != 204:
       retry()    ← klare, einfache Regel
```

**Take-away:** Idempotenz schlägt Ehrlichkeit. Ein
Kompensations-Endpoint, der „korrekte" 4xx liefert, zwingt jeden
Aufrufer dazu, die 4xx wieder als „eigentlich ok" zu interpretieren —
das ist die schlechtere Stelle für die Sonderlogik.

---

## Frage 5 — Wie merkt man überhaupt, dass eine Saga gerade hängt?

**Frage:** Wenn Forward fehlschlägt und Kompensation auch, sieht der
Kunde eine Fehlermeldung und Booking schreibt eine Logzeile. Reicht das?

**Antwort:** Nein — Logzeilen sind die Untergrenze. Eine Saga lebt von
**Beobachtbarkeit ihres Zustands**, nicht von Einzel-Logs. Was im
Workshop nicht implementiert ist, was aber in jeder produktiven
Saga-Implementierung gehört:

| Signal | Was es zeigt | Alarm-Schwelle |
|---|---|---|
| **Saga-Status-Verteilung** (Counter) | wieviele in `PENDING`, `COMPLETED`, `FAILED`, `COMPENSATING` | `COMPENSATING > 0` für > N Min |
| **Saga-Dauer** (Histogramm) | wie lange läuft eine Saga im Schnitt | p99 > erwarteter Wert |
| **Compensation-Erfolgsrate** | wieviele Kompensationen klappen beim ersten Versuch | < 99 % → Backend krank |
| **Retry-Counter pro Schritt** | wie oft wird derselbe `DELETE` wiederholt | unerwartete Spitzen |
| **DLQ-Tiefe** (bei Eventing) | wieviele Events sind unprozessierbar | > 0 → manueller Eingriff |
| **Distributed Tracing** | welcher Schritt war Auslöser | jeder Span trägt Saga-ID |

**Spicy Take-away:** „Wir loggen das" ist die Antwort von Teams, die
noch keine hängende Saga im Produktionsbetrieb gesehen haben. Eine
hängende Saga ist nicht laut — sie ist **still**. Sie schreibt keine
Fehlermeldung mehr, weil der Step, der sie geschrieben hätte, nicht mehr
läuft. Das einzige, was sie sichtbar macht, ist ein **Counter, der zu
lange auf einem Wert stehen bleibt.**

---

## Frage 6 — Wer ist der Owner einer Saga: Booking oder die Backends?

**Frage:** In unserer Implementierung wissen Hotel/Flight/Car gar nichts
davon, dass sie Teil einer Saga sind. Sie kennen nur ihre eigenen
`POST` und `DELETE`. Ist das richtig so, oder müsste die Saga-Idee
weiter „nach unten" durchgereicht werden?

**Antwort:** Bei **Orchestration-Saga** ist es genau richtig so —
und das ist eines der Kernargumente für Orchestration:

```
   booking (Orchestrator) ──────►  weiß alles
       │                            ▪ Reihenfolge der Schritte
       │                            ▪ wer wurde schon aufgerufen
       │                            ▪ was muss kompensiert werden
       │                            ▪ aktuellen Saga-Status
       │
       ├─POST /bookings─►  hotel    weiß: nur „eine Buchung"
       ├─POST /bookings─►  flight   weiß: nur „eine Buchung"
       └─POST /bookings─►  car      weiß: nur „eine Buchung"
```

Hotel, Flight, Car bleiben **dumm und einfach**. Sie kennen ihre
eigenen lokalen Transaktionen — Buchen, Stornieren — und sonst nichts.
Das ist gut, weil:

- die Backends in *anderen* Sagen mitspielen können (z.B. „Premium-Reise
  + Versicherung") ohne Code-Änderung.
- ein Bug in der Saga-Logik nur an **einer** Stelle steckt (Booking),
  nicht in drei.
- ein neues Backend (z.B. „Mietboot") ohne Anpassung der Bestehenden
  hinzukommen kann.

Bei **Choreography-Saga** (Story 6) wandert ein Stück Saga-Wissen in die
Backends — sie reagieren auf Events anderer. Das ist eleganter aber
auch verwobener: ein Bug in der Saga-Logik kann jetzt in jedem der
beteiligten Services sitzen.

**Take-away:** Orchestration konzentriert das Wissen. Choreography
verteilt es. Beides ist legitim, aber **Wissen verteilen ohne Plan**
führt zu „verteilter Monolith" — der schlimmsten beider Welten.

---

## Frage 7 — Braucht eine Saga zwingend eine Datenbank?

**Frage:** Wir halten den Saga-Status in unserem Workshop in-memory. In
Produktion wäre das nicht akzeptabel. Heißt das, jede Saga braucht eine
„echte" Datenbank? Und wenn ja: welche Art?

**Antwort:** Hier zuerst die wichtigste Klarstellung: **Persistenz ist
kein Saga-spezifisches Thema.** Jeder Service, der einen mehrstufigen
Prozess steuert, kann mitten im Ablauf abstürzen und seinen
In-Memory-State verlieren. Saga macht das Problem nur **sichtbarer**,
weil die einzelnen Schritte externe Seiteneffekte (gebuchte Flüge,
gestartete Zahlungen) hinterlassen. Aber der **Mechanismus** dahinter
— „mein Prozess kann mitten drin sterben, ich brauche durable State,
um aufzuräumen" — gilt auch für Batch-Jobs, Workflows, Long-Running
HTTP-Handler oder Importer.

Was die Saga **wirklich** braucht, ist also nicht „eine Datenbank" im
engen Sinne, sondern eine **durable, transaktional konsistente Ablage
ihres Fortschritts** — derselbe Anspruch, den jedes mehrstufige System
hat. Einziger zwingender Grund: **Crash-Recovery des Orchestrators.**

```
   Ohne Persistenz                          Mit Persistenz
   ─────────────────                        ─────────────────
   Saga läuft im RAM                        Saga liegt in DB
        │                                        │
   Flug gebucht ✓                           Flug gebucht ✓
                                            → Saga-State persistiert
        │                                        │
   Orchestrator stürzt 💥                    Orchestrator stürzt 💥
        │                                        │
   Niemand weiß mehr, dass der              Neuer Orchestrator startet,
   Flug gebucht ist                         liest offene Sagas, kompensiert
        ↓                                        ↓
   Orphan booking, kein Aufräumen           Aufräumen funktioniert
```

Ohne Persistenz ist die Saga damit zwar **keine vollständige Saga**,
sondern ein optimistisches Skript — aber das wäre **jeder andere
mehrstufige Prozess** ohne State-Persistenz auch.

**Spektrum der „Datenbank" — von leichtgewichtig zu vollausgestattet:**

| Variante | Beispiele | Wann sinnvoll |
|---|---|---|
| Relationale DB | Postgres, MySQL — Saga-Tabelle + Steps-Tabelle | Standard, wenn die App ohnehin eine DB hat |
| Document Store | MongoDB, DynamoDB — Saga als JSON-Dokument | wenn der Saga-Ablauf variabel ist |
| Event Log | Kafka Compacted Topics, Event Sourcing | Saga-State ist abgeleitet aus Events |
| Embedded Store | SQLite, BoltDB | Single-Node-Deployments, Edge-Services |
| KV-Store | Redis (AOF), etcd, **Consul KV** | leichtgewichtig, aber schwächere Transaktionalität |
| Workflow Engine | Temporal, Cadence, Camunda, Axon Framework | „Saga as Code" — Engine löst Recovery für dich |

**Der elegante Ausweg in modernen Stacks:** Workflow Engines drehen
das Modell um. Statt selbst State, Retry und Recovery zu coden, schreibt
man die Geschäftslogik als Workflow, die Engine kümmert sich um den Rest:

```
   Manuell (was wir gerade gebaut haben)
   ─────────────────────────────────────
   Code:   forward; if err { state=COMPENSATING; compensate }
   Du:     bist verantwortlich für State, Persistenz, Retry, Recovery

   Mit Temporal / Camunda / ...
   ─────────────────────────────────────
   Code:   workflow.Execute(BookFlight)
           workflow.Execute(BookHotel)
           workflow.Execute(BookCar)
   Engine: persistiert jeden Schritt, retryt, startet nach Crash neu,
           garantiert exactly-once-Semantik
```

**Bezug zum Workshop:** Unser Booking-Service hat **gar keine
DB-Schicht** — auch Story 1–4 nicht. Wenn wir Persistenz workshop-
konform nachziehen wollten, wären die kleinsten Schritte:

1. **SQLite-Datei im Container** — eine Datei, kein neues Infra-Stück.
2. **Consul KV nutzen** — Consul ist seit Story 2 im Stack. Saga-State
   unter `kv/saga/{id}` ablegen. Pragmatisch, aber Atomicity bei
   Step-Updates ist schwächer.
3. **Postgres oder Redis** — saubere Lösung, aber neue Compose-
   Komponente und mehr Boilerplate.

Für einen 60-Minuten-Slot wäre das alles zu viel. Deshalb in-memory —
und der Punkt steht hier statt im Code.

**Take-away:** Persistenz ist **nicht das Saga-Pattern, sondern
generelle Orchestrator-Robustheit**. Jeder mehrstufige Prozess braucht
sie — Saga macht das Problem nur sichtbarer, weil ihre Schritte
außerhalb des eigenen Service Spuren hinterlassen. Welche Technologie
genau, ist sekundär — entscheidend ist die Zusicherung **„nach Crash
kann ich aufräumen"**. In Produktion ist die wichtigere Frage selten
„brauche ich eine DB?", sondern **„schreibe ich die Orchestrator-
Mechanik selbst, oder nutze ich eine Workflow Engine?"**. Letzteres
unterschätzen Teams regelmäßig — und schreiben dann monatelang das,
was Temporal seit Jahren in Produktion löst.

---

## Sammelthemen für die Diskussion

- Welche Schritte einer Reisebuchung sind eigentlich **gar nicht
  technisch kompensierbar**? (Hint: Zahlung. Refund ist eine fachliche
  Gegenbuchung, kein Delete.)
- Wenn Hotel und Flight nach demselben Crash beide „wegen Saga"
  storniert wurden — wie würde der Kunde davon erfahren? Welcher Service
  schickt die E-Mail?
- Was passiert, wenn der Orchestrator selbst während `COMPENSATING`
  abstürzt? Was muss persistiert sein, damit der Resume nach Neustart
  korrekt ist?
- Würdest du den **gleichen Saga-Code** wiederverwenden für eine
  Reise-Stornierung, die der Kunde aktiv anstößt (also Forward-Step
  „Storno", nicht Compensation)? Oder ist das ein anderer Use-Case?
  (Hint: oft ist es derselbe — die Saga ist nur eine Sequenz von
  lokalen Transaktionen. Die Richtung ist Konvention, nicht Pattern.)
- Brücke zu Story 6: Welcher Teil der heutigen Saga-Implementierung
  würde komplett wegfallen, wenn wir auf Eventing umstellen, und welcher
  Teil bleibt unverändert?
