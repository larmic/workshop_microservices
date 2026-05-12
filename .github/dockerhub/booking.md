# Workshop Microservices – Booking Service

Demo-Image aus dem **Microservices-Workshop** von [larmic/workshop_microservices](https://github.com/larmic/workshop_microservices).

## Inhalt

Booking-Service, der Buchungen über die Downstream-Services *flight*, *hotel* und *car* orchestriert. Jede Workshop-Story baut iterativ auf der vorherigen auf und führt ein neues Resilience- oder Architektur-Pattern ein.

## Tags

Die Story-Auswahl erfolgt **explizit über den Tag** (kein `latest`):

| Tag       | Inhalt                                                                |
|-----------|-----------------------------------------------------------------------|
| `story1`  | Cloud-Native Grundlagen (12-Factor, Health Checks)                    |
| `story2`  | Service Discovery via Consul                                          |
| `story3`  | Circuit Breaker                                                       |
| `story4`  | Bulkhead-Pattern                                                      |
| `story5`  | Saga (Orchestration)                                                  |
| `story6`  | Saga (Choreography)                                                   |
| `story7`  | Distributed Tracing                                                   |

Pull-Beispiel:

```bash
docker pull larmic/workshop-microservices-booking:story3
```

## Hinweis

Reines Lehr- und Demo-Image. **Nicht** für den Produktivbetrieb gedacht.

## Quellcode & Doku

- Repository: https://github.com/larmic/workshop_microservices
- Stories: [`docs/stories/`](https://github.com/larmic/workshop_microservices/tree/main/docs/stories)
