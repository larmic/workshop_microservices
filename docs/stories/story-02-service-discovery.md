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
- die Consul HTTP API nutzt, um gesunde Service-Instanzen zu ermitteln
- bei Ausfall einer Service-Instanz automatisch eine andere verwendet

---

## Aufgaben

### 1. Consul Resolver implementieren
- Eigenes `consul`-Package mit einem `Resolver` erstellen
- Consul HTTP API (`/v1/health/service/{name}?passing=true`) abfragen
- Service-URL aus der Antwort (Address + Port) zusammenbauen

### 2. Service Discovery statt statischer URLs
- Statische URLs aus [Story 1](story-01-cloud-native-booking-service.md) durch Service-Namen ersetzen
  - `flight-service`
  - `hotel-service`
  - `car-service`
- Backend-Services über den Consul Resolver dynamisch auflösen
- Handler in eigenes `handler`-Package extrahieren

### 3. Ausfallsicherheit (optional)
- Client-Side Load Balancing bei mehreren Instanzen (zufällige Auswahl)
- Consul Key-Value Store für zusätzliche Konfiguration
