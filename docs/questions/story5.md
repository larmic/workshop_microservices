# Workshop-Fragen: Saga (Story 5)

Provokante Fragen rund um das Saga Pattern und die Frage, was passiert,
wenn das „alles oder nichts" *selbst* anfängt zu wackeln. Ziel: nicht
nur den Happy Path verstehen, sondern die Failure-Modi, die in echten
Systemen die schwierigen sind — und die Brücke zu Story 6 (Events &
CQRS) bewusst sehen.

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
die Saga **isoliert** (sync, ein Konzept). Story 6 schaltet den Broker
dazu (asynchron, ein zweites Konzept). Bewusst zwei Schritte, nicht
einer.

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
