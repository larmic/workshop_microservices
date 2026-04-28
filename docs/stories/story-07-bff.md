# Story 7: Mobile First

**Thema:** Backends for Frontends (BFF)
**Zeitrahmen:** ca. 60 Minuten

## Kontext

Die mobile App benötigt andere Daten als die Web-Anwendung: kleinere Payloads, zusammengefasste Informationen, optimiert für langsame Verbindungen. Ein spezialisiertes Backend für das mobile Frontend löst dieses Problem.

## User Story

Als **Mobile-App-Entwickler**
möchte ich **ein Backend, das genau die Daten liefert, die meine App braucht**,
damit **ich keine überflüssigen Daten übertragen muss und die App auch bei schlechter Verbindung performant bleibt**.

## Akzeptanzkriterien

- [ ] Ein Mobile-BFF-Service ist implementiert
- [ ] Das BFF aggregiert Daten aus mehreren Backend-Services
- [ ] Die Payload ist kompakter als die der einzelnen Services
- [ ] Nur für Mobile relevante Felder werden zurückgegeben
- [ ] Das BFF cached häufig angefragte, selten ändernde Daten
- [ ] GraphQL oder spezialisierte REST-Endpoints sind verfügbar

## Technische Hinweise

- **BFF vs. API Gateway:**
  - API Gateway: Routing, übergreifende Policies
  - BFF: Frontend-spezifische Aggregation und Transformation
- **Payload-Optimierung:**
  - Felder auswählen (nur benötigte)
  - Daten zusammenfassen (z.B. "Reise" statt separate Flight/Hotel/Car)
  - Pagination für Listen
- **Empfohlene Technologien:**
  - GraphQL für flexible Abfragen
  - REST mit Sparse Fieldsets

## Bonus (optional)

- Implementiere ein separates BFF für die Web-Anwendung
- Füge Offline-Fähigkeit durch intelligentes Caching hinzu
