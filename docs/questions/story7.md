# Workshop-Fragen: Choreography-Saga (Story 7)

Provokante Fragen rund um die asynchrone Kompensation via Events und die
Frage, was unsere Schmalspur-Variante (HTTP-POST als Bus-Ersatz) alles
**nicht** kann — und warum echte Message-Broker genau diese Lücken füllen.
Ziel: Den Sprung von Story 6 zu Story 7 nicht als „besseres Pattern"
verkaufen, sondern als bewusste Verschiebung der Komplexität.

---

## Frage 1 — Was passiert, wenn der Event-POST selbst fehlschlägt?

**Frage:** Booking schickt `CompensationRequested` per HTTP-POST an jedes
Backend, das schon gebucht hat, rückwärts (erst Hotel, dann Flight). Was,
wenn Flight seit der Buchung down ist und der POST in einen
Connection-Refused, einen 5xx oder den 3-s-Timeout läuft? In Story 6 hätte
Booking diesen Fehler dem Kunden zurückgegeben. Was tun wir in Story 7?
(Im Spiel: Beim Kassensturz liegt noch ein Chip bei Flight, obwohl Booking
„storniert“ gemeldet hat.)

**Antwort:** Aktuell genau **nichts**. Booking wartet beim Zuruf nur auf
die Annahme: Im Spiel steht der Wurf 1–4 der Gerufenen dafür. Bleibt sie
aus, merkt Booking das sogar und loggt `phase=dispatch-failed`. Der Step
wird trotzdem als `COMPENSATED` markiert, die Saga geht auf `FAILED` und
der Kunde bekommt seine Antwort. Es gibt genau einen Versuch, kein Retry.
Das ist *bewusst* die fragilste mögliche Variante, weil sie das
eigentliche Problem von Eventing-ohne-Broker unverstellt zeigt:

```
   Booking sagt: „storniert."
   Log:          dispatch-failed, eine Zeile unter vielen.
   Realität:     Event wurde nie empfangen.
                 Flight hat die Buchung noch im Bestand (in echtem System).
                 Niemand schaut hin, bis ein Kunde sich beschwert.
```

Die Lehrbuch-Antworten, die ein echtes System braucht:

| Mechanismus | Was er löst |
|---|---|
| **Retry mit Backoff** | transienten Backend-Hänger aussitzen |
| **Outbox-Pattern** | Event in derselben DB-Transaktion wie der Saga-State speichern, separater Worker publiziert es |
| **Dead-Letter-Queue** | nicht zustellbare Events sammeln, Operator kann manuell entscheiden |
| **At-least-once-Bus** | der Broker garantiert: Event wird zugestellt, auch wenn beim ersten Versuch alles brennt |

**Take-away:** Eventing **eliminiert** das „Backend kurz weg"-Problem
nicht — es **verschiebt** es. In Story 6 hat Booking den Schmerz gespürt
und konnte reagieren. In Story 7 steht er nur noch als Strich im Log. Das ist
das Gegenteil von „resilient", solange wir nichts gegen verlorene Events
haben. Wer Choreography ernst meint, fängt nicht beim Event-Versand an,
sondern bei der **Durability des Events**.

---

## Frage 2 — Warum ist Idempotenz bei Events ein Feature, kein Nice-to-have?

**Frage:** Unser Backend tut bei `POST /events/compensation` *nichts
außer loggen*. Welche Konsequenz hätte es, wenn dasselbe Event zweimal
zugestellt würde? Müssen wir uns überhaupt darum kümmern? (Im Spiel nur als
Gedankenexperiment auf der Leiter: Der Korb stellt noch einmal zu, der Service
zahlt zweimal.)

**Antwort:** In unserer Schmalspur-Variante mit einmaligem `POST`
**nie**: ein Versuch, kein Retry, also keine Duplikate, aber auch keine
Nachlieferung. Aber sobald ein echter Broker im Spiel ist, gilt
**at-least-once** als Standardgarantie — exactly-once ist die Ausnahme
und teuer. Doppelte Zustellung kann passieren bei:

```
1. Sender macht POST → Backend verarbeitet → 200 zurück, geht
   aber im Netz verloren → Sender retryt → Backend bekommt das
   Event ZWEIMAL.
2. Broker bestätigt Empfang, Consumer crasht VOR dem Commit
   des Offsets → Broker liefert das Event nochmal an die nächste
   Instanz.
3. Outbox-Worker hat Event versendet, crasht VOR dem Markieren als
   „verschickt" → Restart, schickt erneut.
```

Unser Backend ist heute **idempotent durch Zustandslosigkeit** — es
loggt nur, kein State ändert sich. Das ist eine glückliche
Workshop-Eigenschaft, in Produktion ist sie die Ausnahme:

| In echt | Was passiert ohne Idempotenz |
|---|---|
| Refund auf Kreditkarte | Kunde bekommt zweimal das Geld zurück |
| Lager-Reservierung freigeben | Bestand wird zweimal hochgezählt |
| Inkasso-Stop | Mehrfach „gestoppt" liefert Race-Conditions im nachgelagerten System |

Der Standard-Mechanismus ist ein **Dedup-Key** — das ist genau das Feld
`eventId` in unserer `CompensationEvent`-Struktur. Der Backend-Service
führt eine kleine Tabelle der zuletzt verarbeiteten `eventId`s und
ignoriert Wiederholungen. Was im Workshop nicht implementiert ist, aber
**in jedem Production-Event-Handler dazu gehört**:

```
def handle(event):
    if seen(event.eventId):
        return                              # bereits verarbeitet, skip
    with transaction:
        do_business_logic(event)
        mark_seen(event.eventId)            # in derselben Transaktion!
```

**Take-away:** Idempotenz ist nicht „nice to have" — sie ist das, was
einen Event-Handler von einem **Random-Number-Generator unter Last**
unterscheidet. Wer Events publiziert, muss davon ausgehen, dass sie
mehrfach ankommen. Wer das ignoriert, hat einen Bug, der erst unter Last
zuschlägt — also genau dann, wenn niemand ihn debuggen kann.

---

## Frage 3 — Warum brauchen wir echte Broker (Kafka, RabbitMQ, NATS)?

**Frage:** Wir machen Eventing mit HTTP-POST gegen die Backends via
Consul. Das ist schnell, einfach, kostet keine zusätzliche Infrastruktur.
Wofür braucht man dann überhaupt einen echten Message-Broker?

**Antwort:** Für genau die Dinge, die unsere Variante **strukturell
nicht leisten kann** — und das sind mehr, als man auf den ersten Blick
denkt:

| Aspekt | HTTP-Events (Workshop) | Echter Broker (Kafka / RabbitMQ / NATS / …) |
|---|---|---|
| **Persistenz** | Event existiert nur im RAM des Senders | Nachricht auf Disk, überlebt Sender- und Empfänger-Crash |
| **Redelivery** | nicht eingebaut — verloren ist verloren | At-least-once garantiert, optional Exactly-once |
| **Fan-out** | Sender muss alle Subscriber kennen und einzeln anrufen | Ein Topic, beliebig viele Subscriber, dynamisch dazukommen lassen |
| **Backpressure** | langsamer Empfänger blockiert / bricht den Sender | Queue puffert, Empfänger zieht im eigenen Tempo |
| **Dead-Letter** | gibt es nicht | Standardfeature für nicht zustellbare Nachrichten |
| **Consumer-Crash** | Event ist weg | Nachricht bleibt in Queue, neuer Consumer übernimmt am letzten Offset |
| **Ordering** | keine Garantie | per Partition / Topic-Sequence |
| **Replay** | unmöglich (Event existiert nicht mehr) | Konsument kann an alten Offset springen und neu konsumieren |
| **Entkopplung in der Zeit** | Sender und Empfänger müssen gleichzeitig erreichbar sein | Sender kann publizieren, wenn Empfänger schläft / down ist |
| **Routing-Logik** | im Sender hartcodiert | per Topic / Routing-Key zentral konfigurierbar |
| **Infrastruktur-Aufwand** | null — nur HTTP | zusätzlicher Service, Operations-Overhead, Cluster-Tuning |

Die letzte Zeile ist der einzige Punkt, in dem unsere Variante gewinnt —
und genau auf diesen Punkt fallen viele Teams herein:

```
   „Wir wollen Eventing, aber keinen Broker betreiben."
   → 6 Monate später: Outbox-Tabelle, Retry-Worker, Dedup-Logik,
     DLQ-Inbox, Replay-Skript, alles selbst gebaut.
   → 12 Monate später: das ist ein schlechter Broker.
```

**Take-away:** Choreography-Saga **ohne Broker** ist eine pädagogische
Übung. Die Architekturidee — Verantwortung verteilen, Backends
selbstständig kompensieren, Booking nicht zum Single Point of
Coordination machen — ist richtig und wichtig. Aber die Idee braucht
**durable messaging**, sonst wird sie zur Verbesserung der schlechten Art:
„Es fühlt sich lockerer gekoppelt an, ist aber stiller im Fehlerfall."
Wer Broker scheut, baut sich einen — und das wird meistens schlechter,
als das fertige Werkzeug zu nehmen.

---

## Frage 4 — Brauchen wir überhaupt einen Dirigenten?

**Frage:** Volle Choreographie hieße: Hotel reagiert auf Flight, Car auf
Hotel, niemand dirigiert. Wer merkt es, wenn Hotel nicht weiß, dass es
reagieren soll?

**Antwort:** Niemand. Für unsere Buchung ist **Orchestrierung** der
bessere Weg. Drei Schritte in fester Reihenfolge mit Kompensation sind ein
Prozess, und den will man an einer Stelle lesen, testen und im Log
verfolgen. Volle Choreographie verteilt die Reihenfolge auf die Köpfe der
Services: Jeder muss wissen, auf wen er hört. Fehlt ein Abonnement, bricht
die Kette, und niemand bekommt einen Fehler — es passiert einfach nichts.
Das ist der **verteilte Monolith**: lose gekoppelt zur Laufzeit, fest
gekoppelt im Kopf.

Was unsere Referenz macht (im Spiel der Zuruf), ist deshalb der
sinnvolle Mittelweg: Booking dirigiert den Hinweg und entscheidet auch,
wer ein Storno bekommt, storniert wird per Zuruf. Streng genommen ist
`CompensationRequested` damit ein Command und keine echte Choreografie.

**Faustregel nach Chris Richardson** (*Microservices Patterns*, 2018,
Kapitel 4):

| Wann | Was |
|---|---|
| einfache Saga, wenige Beteiligte, keine Reihenfolge (Benachrichtigung, Read-Model) | Choreographie |
| Ablauf mit mehreren Schritten oder Verzweigungen | Orchestrierung, sonst werden die Abhängigkeiten zyklisch und niemand sieht mehr den Prozess |

**Take-away:** Choreographie klingt eleganter, weil kein Service den
anderen kennt. In Wahrheit kennt jeder jeden, nur steht es nirgends.
Brücke zu Story 8: Wer merkt es eigentlich, wenn ein Event verloren geht?

---

## Sammelthemen für die Diskussion

- Wie würdet ihr **Exactly-once** in unserer aktuellen Lösung herstellen?
  (Hint: gar nicht. Idempotenz beim Empfänger + at-least-once beim Sender
  ist die praktische Annäherung.)
- Welcher Teil unserer Story-7-Implementierung würde komplett wegfallen,
  wenn wir einen echten Broker hätten? (Hint: der `compensateSingle`-Code
  würde zum `bus.Publish(topic, event)`-Einzeiler. Service-Discovery
  brauchen wir auch nicht mehr — der Broker macht das implizit.)
- Falls jemand nach **CQRS** fragt (Reserve in `docs/themen.md`, Block 7): Wenn wir das
  Read-Model aus Events füttern und ein Event verloren geht — wie
  repariert man das Read-Model? (Hint: Event-Replay; setzt aber
  **persistente** Events voraus, also einen Broker mit Log-Semantik
  wie Kafka.)
