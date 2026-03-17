# Story 2: Services dynamisch finden

**Thema:** Service Discovery / Service Registry
**Zeitrahmen:** ca. 60 Minuten

## Kontext

In einer dynamischen Container-Umgebung ändern sich IP-Adressen und Ports der Services ständig. Der Booking-Service soll die Backend-Services nicht mehr über statische URLs ansprechen, sondern über eine Service Registry dynamisch finden.

## User Story

Als **Booking-Service**
möchte ich **die Backend-Services über ihren logischen Namen finden können**,
damit **ich nicht von statischen IP-Adressen oder Ports abhängig bin und Services elastisch skaliert werden können**.

## Akzeptanzkriterien

- [ ] Der Booking-Service registriert sich bei Consul
- [ ] Der Booking-Service findet FlightService, HotelService und CarService über Consul
- [ ] Die statischen URLs aus **[Story 1](story-01-cloud-native-booking-service.md)** werden durch Service-Namen ersetzt
- [ ] Der Health-Check wird bei Consul registriert
- [ ] Bei Ausfall eines Service-Instanz wird automatisch eine andere verwendet (falls verfügbar)
- [ ] Der Service deregistriert sich beim Herunterfahren

## Technische Hinweise

- **Consul UI:** `http://localhost:8500` (wird vom Trainer bereitgestellt)
- **Service-Namen:**
  - `flight-service`
  - `hotel-service`
  - `car-service`
- **Empfohlene Libraries:**
  - Spring Boot: `spring-cloud-starter-consul-discovery`
  - Quarkus: `quarkus-consul-config`
  - Alternativ: Direkte Consul HTTP API
- **Consul Health-Check-Typen:** HTTP, TCP, TTL

## Bonus (optional)

- Implementiere Client-Side Load Balancing bei mehreren Instanzen
- Nutze Consul Key-Value Store für zusätzliche Konfiguration
