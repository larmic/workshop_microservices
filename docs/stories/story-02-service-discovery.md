# Story 2: Services dynamisch finden

**Zeitrahmen:** ca. 60 Minuten

## Thema

- **Service Discovery** → Dynamisches Auffinden von Services in Container-Umgebungen
- **Service Registry** → Zentrale Registrierung und Verwaltung von Service-Instanzen (Consul)
- **Client-Side Load Balancing** → Automatische Verteilung auf verfügbare Instanzen

---

## Ziel

Ein Booking-Service, der:
- Backend-Services (Flight, Hotel, Car) über logische Namen statt statischer URLs findet
- sich bei Consul registriert und seinen Health-Check bereitstellt
- bei Ausfall einer Service-Instanz automatisch eine andere verwendet
- sich beim Herunterfahren sauber deregistriert

---

## Aufgaben

### 1. Service-Registrierung bei Consul
- Booking-Service registriert sich bei Consul (`http://localhost:8500`)
- Health-Check wird bei Consul registriert (HTTP, TCP oder TTL)
- Service deregistriert sich beim Herunterfahren

### 2. Service Discovery statt statischer URLs
- Statische URLs aus [Story 1](story-01-cloud-native-booking-service.md) durch Service-Namen ersetzen
  - `flight-service`
  - `hotel-service`
  - `car-service`
- Backend-Services über Consul dynamisch auflösen

### 3. Ausfallsicherheit (optional)
- Client-Side Load Balancing bei mehreren Instanzen
- Consul Key-Value Store für zusätzliche Konfiguration
