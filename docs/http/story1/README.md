# HTTP Client Tests — Story 1

HTTP-basierte Tests für die Services aus Story 1 (Flight, Hotel, Car, Booking).

## Voraussetzung

Services müssen laufen:

```bash
cd services && make docker-up
```

## Ausführung

**Headless via CLI (Docker):**

```bash
make test
```

**In IntelliJ / WebStorm:**

`requests.http` direkt öffnen und einzelne Requests ausführen.

## Referenz

- [jetbrains/intellij-http-client](https://hub.docker.com/r/jetbrains/intellij-http-client) — Docker-Image für headless-Ausführung
