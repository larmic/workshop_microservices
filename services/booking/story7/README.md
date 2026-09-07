# Booking Service — Story 7: Choreography-Saga via Events

## Was diese Story zeigt

Story 6 implementiert die Kompensation **synchron**: Bei einem Fehlschlag im
Forward-Pfad ruft der Booking-Service nacheinander `DELETE /bookings/{id}`
gegen Flight, Hotel und Car auf und wartet auf jede Antwort, bevor er dem
Kunden `FAILED` zurückgibt.

Story 7 dreht die Verantwortung um: Booking publiziert pro Schritt ein
`CompensationRequested`-Event an die Backend-Services und ist **sofort
fertig**. Die Backends führen die Stornierung selbst in einer Goroutine aus.

## Im Vergleich

```
Story 6 (Orchestration, synchron)
─────────────────────────────────
  booking ─DELETE /bookings/F-...──► flight   (wartet auf 204)
  booking ─DELETE /bookings/H-...──► hotel    (wartet auf 204)
  booking ─DELETE /bookings/C-...──► car      (wartet auf 204)
  booking → Antwort an Kunden: FAILED


Story 7 (Choreography, asynchron — fire & forget)
─────────────────────────────────────────────────
  booking ─POST /events/compensation──► flight ──► 202 Accepted (sofort)
                                                    │
                                                    └─► storniert in Goroutine
  booking ─POST /events/compensation──► hotel  ──► 202 Accepted (sofort)
                                                    │
                                                    └─► storniert in Goroutine
  booking ─POST /events/compensation──► car    ──► 202 Accepted (sofort)
                                                    │
                                                    └─► storniert in Goroutine
  booking → Antwort an Kunden: FAILED  (ohne auf die Backends zu warten)
```

| Aspekt | Story 6 (Orchestration) | Story 7 (Choreography) |
|---|---|---|
| **Wer storniert?** | Booking-Service zentral | Jedes Backend selbst |
| **Wartet Booking auf Ergebnis?** | Ja, pro Schritt | Nur auf `202 Accepted` |
| **Antwortzeit für den Kunden bei Fehler** | Summe aller Kompensations-Latenzen | Konstant, unabhängig von Backend-Last |
| **Wo sitzt der Bug, wenn Stornierung schiefgeht?** | im Booking-Service | im jeweiligen Backend |
| **Robustheit bei Backend-Ausfall** | Booking sieht den Fehler und kann reagieren | Booking weiß nichts, Event ist verloren |

Der letzte Punkt ist absichtlich der unangenehmste — und der eigentliche
didaktische Kern: **Choreography löst nicht alles**. Was unsere
Schmalspur-Variante mit HTTP-Events alles **nicht** kann (Persistenz,
Redelivery, Dead-Letter, Fan-out), und warum echte Message-Broker (Kafka,
RabbitMQ, NATS) genau deshalb existieren, wird in
[`docs/questions/story7.md`](../../../docs/questions/story7.md)
diskutiert.

## Logging

Alle vier beteiligten Services schreiben strukturierte Log-Zeilen mit
denselben Korrelations-Feldern (`eventId`, `sagaId`). Damit lässt sich ein
einzelner Saga-Verlauf über alle Services hinweg per `grep` rekonstruieren:

```bash
docker compose logs booking-ref-story7 flight hotel car | grep S-abc123
```

Phasen im Log:

| Phase | Wo | Bedeutung |
|---|---|---|
| `publishing` | booking-ref-story7 | Booking hat das Event vorbereitet und versendet jetzt |
| `dispatched` | booking-ref-story7 | Backend hat `202 Accepted` zurückgegeben |
| `dispatch-failed` | booking-ref-story7 | Versand fehlgeschlagen (Netz / 5xx) — Step bleibt trotzdem `COMPENSATED` |
| `received` | flight/hotel/car | Event ist beim Backend angekommen, vor der Verarbeitung |
| `processing` | flight/hotel/car | Asynchrone Verarbeitung in der Goroutine startet |
| `done` | flight/hotel/car | Verarbeitung abgeschlossen |
