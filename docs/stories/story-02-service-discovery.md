# Story 2: Services dynamisch finden

> Der Booking-Service aus Story 1 läuft – aber die Backend-URLs in der Konfiguration werden zum Problem. Jedes Mal, wenn der Hotel-Service neu deployed wird, ändert sich der Port. Und wenn wir zwei Instanzen vom Flight-Service hochfahren wollen, weiß unser Booking-Service davon erstmal nichts.
>
> Wir hätten gerne, dass sich die Backend-Services selbst registrieren und vom Booking-Service über einen logischen Namen gefunden werden. Fällt eine Instanz aus, soll automatisch eine andere übernehmen – ohne dass jemand nachts um drei eine Konfigurationsdatei anfasst. Statische URLs sind eine schöne Idee für Folien, aber nicht für den Betrieb.

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

### 3. Buchung durchführen
- `POST /booking/bookings` im Booking-Service
  - Request: `{ flightId, hotelId, carId, customerName }`
  - Ruft die drei Backend-Services über den Consul Resolver auf:
    - `POST /bookings` am Flight-Service
    - `POST /bookings` am Hotel-Service
    - `POST /bookings` am Car-Service
  - Antwortet mit aggregierter Buchung: `{ bookingId, customerName, flight, hotel, car }`
  - Kein Error-Handling (darf fehlschlagen)
- Backend-Services bestätigen die Buchung mit eigener Booking-ID und loggen sie auf stdout (keine Persistenz)

### 4. Ausfallsicherheit (optional)
- Client-Side Load Balancing bei mehreren Instanzen (zufällige Auswahl)
- Consul Key-Value Store für zusätzliche Konfiguration
