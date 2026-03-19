# Story 1: Der erste cloud-native Booking-Service

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

### 1. Health-Checks
- Health-Endpoint (`/health`) → HTTP 200 wenn gesund
- Logging auf stdout (nicht in Dateien)

### 2. Externe Konfiguration
- URLs der Backend-Services konfigurierbar (Flight, Hotel, Car)
- Umgebungsvariablen überschreiben Konfiguration
- Info-Endpoint zeigt aktuelle Config (ohne Secrets)

### 3. Erste Orchestrierung
- `GET /booking/offers` → ruft Flight-, Hotel- und CarService auf
- Kombinierte Ergebnisse als JSON zurückgeben
- Kein Error-Handling nötig (darf fehlschlagen)
