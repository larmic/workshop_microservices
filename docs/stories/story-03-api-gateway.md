# Story 3: Ein Gateway für alle

**Thema:** API-Gateway
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Externe Clients (Web-App, Mobile-App) sollen nicht direkt mit den einzelnen Services kommunizieren. Ein API-Gateway dient als zentraler Einstiegspunkt, übernimmt Routing und kann übergreifende Aspekte wie Rate-Limiting implementieren.

## User Story

Als **Frontend-Entwickler**
möchte ich **alle Backend-Services über einen einzigen Einstiegspunkt ansprechen können**,
damit **ich mich nicht um die interne Service-Topologie kümmern muss und zentrale Policies (Rate-Limiting, CORS) einheitlich angewendet werden**.

## Akzeptanzkriterien

- [ ] Ein API-Gateway ist implementiert und läuft auf Port 8080
- [ ] Anfragen an `/api/flights/**` werden an den FlightService weitergeleitet
- [ ] Anfragen an `/api/hotels/**` werden an den HotelService weitergeleitet
- [ ] Anfragen an `/api/cars/**` werden an den CarService weitergeleitet
- [ ] Anfragen an `/api/bookings/**` werden an den BookingService weitergeleitet
- [ ] Das Gateway nutzt Service Discovery (keine hardcoded URLs)
- [ ] CORS ist für Frontend-Clients konfiguriert
- [ ] Basis Rate-Limiting ist implementiert (z.B. max. 100 Requests/Minute pro Client)

## Technische Hinweise

- **Empfohlene Technologien:**
  - Spring Cloud Gateway
  - Kong (Docker)
  - Traefik (Docker)
  - Eigene Implementierung mit Reverse-Proxy-Pattern
- **Routing-Konfiguration (Beispiel Spring Cloud Gateway):**
  ```yaml
  spring:
    cloud:
      gateway:
        routes:
          - id: flight-service
            uri: lb://flight-service
            predicates:
              - Path=/api/flights/**
  ```

## Bonus (optional)

- Implementiere Request-/Response-Logging im Gateway
- Füge einen `/health`-Aggregator hinzu, der den Status aller Services zusammenfasst
