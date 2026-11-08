# Story 1: Der erste cloud-native Booking-Service

> Wir brauchen einen Booking-Service, der Flüge, Hotels und Mietwagen bündelt. Wo der am Ende läuft – auf einem Server im Keller, in der Cloud oder auf dem Raspberry Pi meiner Nichte – wissen wir noch nicht. Erstmal soll er auf dem Entwicklerrechner laufen, ohne dass wir dafür Spezialwerkzeug brauchen.
>
> Ach ja, und bitte mit einem Health-Check. Wir hatten schon zu oft die Situation, dass alle dachten, der Service läuft einwandfrei – bis ein Kunde anrief und fragte, warum seit drei Stunden nichts mehr geht. Wir möchten den Service später automatisch überwachen lassen, egal auf welcher Plattform er landet. Denn: Vertrauen ist gut, ein Health-Endpoint ist besser.

**Zeitrahmen:** ca. 90 Minuten

**Schon da:** Aus der Hausaufgabe ([vorbereitung.md](../vorbereitung.md)) bringt ihr den Mini-Service mit `GET /health` (HTTP 200), `Dockerfile` und Port `8080` mit. Im Dashboard, Schritt 0 „Start“, habt ihr ihn als `booking-custom` eingeklinkt und getestet. Story 1 baut darauf auf.

## Thema

- **Twelve-Factor-App** → Methodologie für moderne, cloud-native Anwendungen
- **Health-Checks** → Überwachbarkeit und automatischer Neustart bei Problemen
- **Externe Konfiguration** → Ein Build, viele Umgebungen

---

## Ziel

Ein Booking-Service, der:
- in verschiedenen Umgebungen ohne Code-Änderung läuft
- überwacht und automatisch neu gestartet werden kann
- Backend-Services (Flight, Hotel, Car) orchestriert

---

## Aufgaben

### 1. Betrieb
- `/health` habt ihr schon, der Statuscode 200 genügt, der Body ist egal
- Logging auf stdout (nicht in Dateien)

### 2. Externe Konfiguration
- Backend-URLs aus Umgebungsvariablen lesen: `FLIGHT_SERVICE_URL`, `HOTEL_SERVICE_URL`, `CAR_SERVICE_URL` (setzt `docker-compose.custom.yml`, ihr lest sie nur aus)
- `GET /info` zeigt die aktuelle Konfiguration: die drei URLs

### 3. API-Vertrag
- `GET /openapi` liefert die OpenAPI-Spec aus, Vorlage: [`services/booking/story1/api/openapi.yaml`](../../services/booking/story1/api/openapi.yaml), passend zum eigenen Service anpassen

### 4. Erste Orchestrierung
- `GET /booking/offers` → ruft Flight-, Hotel- und CarService auf
- Kombinierte Ergebnisse als JSON zurückgeben
- Kein Error-Handling nötig (darf fehlschlagen)
