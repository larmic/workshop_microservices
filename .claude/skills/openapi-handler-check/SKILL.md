---
name: openapi-handler-check
description: Prüft pro Service, ob alle in api/openapi.yaml deklarierten Paths in den Go-Handlern unter handler/ tatsächlich registriert sind und umgekehrt. Aufrufen nach Änderungen an einer openapi.yaml, einer Route in main.go/handler oder wenn die Konsistenz von Spec und Implementierung verifiziert werden soll.
---

# openapi-handler-check

Die Reference-Implementierung dient als Lehrmaterial. Drift zwischen
`api/openapi.yaml` und der tatsächlichen Handler-Registrierung in Go
ist hier doppelt schädlich: Teilnehmer:innen lernen falsche Pfade
oder finden Endpoints, die nicht dokumentiert sind.

## Geltungsbereich (Go-Services)

```
services/flight/        api/openapi.yaml   main.go + handler/*.go
services/hotel/         api/openapi.yaml   main.go + handler/*.go
services/car/           api/openapi.yaml   main.go + handler/*.go
services/booking/story1 api/openapi.yaml   main.go + handler/*.go
services/booking/story2 …                  …
services/booking/story7 …                  …
```

`services/booking/custom/` enthält eine Spec unter
`src/main/resources/openapi.yaml`. Diese gehört zu einer freien
Teilnehmer-Implementierung (anderes Tech-Stack möglich) und wird
von diesem Skill **übersprungen**, sofern der User es nicht explizit
verlangt.

## Ablauf pro Service

1. **Spec-Paths einlesen** — aus `services/<service>/api/openapi.yaml` alle Keys
   unter `paths:` extrahieren, inkl. Methoden (`get`, `post`, …).
   Pfad-Parameter (`{id}`) als Platzhalter behalten.
2. **Handler-Registrierungen einlesen** — in `services/<service>/main.go` und
   `services/<service>/handler/*.go` nach typischen Mux-Registrierungen suchen:
   - Standardbibliothek: `mux.HandleFunc("…")`, `http.HandleFunc("…")`,
     `mux.Handle("…", …)`
   - chi / gorilla / echo / gin: `r.Get("…", …)`, `r.Post("…", …)`,
     `r.HandleFunc("…")`, `e.GET("…")`, `g.POST("…")`, …
   - Method-spezifische Patterns von `net/http` 1.22+: `"GET /booking/offers"`
3. **Pfade normalisieren** — Pfad-Parameter aus Go (`{id}`, `:id`)
   und OpenAPI (`{id}`) auf eine Form bringen. Trailing Slashes
   konsistent abschneiden.
4. **Vergleich** — pro Pfad+Methode:
   - In Spec **und** Handler → ✓
   - Nur in Spec → ✗ Spec dokumentiert eine Route, die fehlt
   - Nur in Handler → ✗ Route existiert, ist aber nicht dokumentiert
5. **Bericht** — pro Service eine Sektion, am Ende ein Summary.

## Ausgabeformat

```
OpenAPI ⇄ Handler-Check
=======================

flight/
  ✓ GET    /flights
  ✓ GET    /flights/{id}
  ✗ POST   /flights         (nur in Spec, kein Handler in handler/flights.go)

booking/story3/
  ✗ GET    /booking/health  (nur als Handler in main.go, fehlt in openapi.yaml)

Summary: 11/13 Pfade konsistent, 2 Drift-Punkte.
```

## Bei Drift

- Schlage konkrete Edits vor: entweder Spec ergänzen (Vorlage aus einer
  benachbarten Story übernehmen, da die Stories konsistent sein sollen)
  oder Handler nachziehen.
- Führe Edits **nicht** ohne Rückfrage aus — der User entscheidet, ob
  die Wahrheit in Spec oder Code liegt.

## Story-übergreifende Konsistenz

Wenn alle sieben Booking-Stories geprüft werden: Markiere Routen, die
in `story1` existieren, aber in `story5` fehlen (und umgekehrt) als
**Info**, nicht als Fehler — Stories sollen sich didaktisch
unterscheiden, manche Endpoints kommen pro Story hinzu. Aber:
**Health-/Info-Endpoints** (`/health`, `/info`, `/openapi`) müssen in
jeder Story vorhanden sein, sonst Drift melden.

## Audience-Hinweis

Antworten auf Deutsch (Workshop-Zielgruppe: Architekt:innen / Tech-Leads,
siehe `services/CLAUDE.md`).
