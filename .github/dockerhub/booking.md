# Workshop Microservices – Booking Service

Demo-Image aus dem Workshop **"Microservices richtig gemacht"** von [larmic/workshop_microservices](https://github.com/larmic/workshop_microservices).

## Inhalt

Booking-Service, der Buchungen über die Downstream-Services *flight*, *hotel* und *car* orchestriert. Jede Workshop-Story baut iterativ auf der vorherigen auf und führt ein neues Resilience- oder Architektur-Pattern ein.

## Tags

Die Story-Auswahl erfolgt **explizit über den Tag** (kein `latest`):

| Tag       | Inhalt                                                                |
|-----------|-----------------------------------------------------------------------|
| `story1`  | Cloud-Native Grundlagen (12-Factor, Health Checks)                    |
| `story3`  | Service Discovery via Consul                                          |
| `story4`  | Circuit Breaker                                                       |
| `story5`  | Bulkhead-Pattern                                                      |
| `story6`  | Saga (Orchestration)                                                  |
| `story7`  | Saga (Choreography)                                                   |
| `story8`  | Distributed Tracing                                                   |
| `custom`  | Beispiel-Custom-Lösung in Kotlin/Ktor (Story 1, alternative Sprache)  |

Pull-Beispiel:

```bash
docker pull larmic/workshop-microservices-booking:story4
docker pull larmic/workshop-microservices-booking:custom
```

## Hinweis

Reines Lehr- und Demo-Image. **Nicht** für den Produktivbetrieb gedacht.

## Quellcode & Doku

- Repository: https://github.com/larmic/workshop_microservices
- Stories: [`docs/stories/`](https://github.com/larmic/workshop_microservices/tree/main/docs/stories)
