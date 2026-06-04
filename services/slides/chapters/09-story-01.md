## Story 1

<p class="subtitle">Erst lauff&auml;hig, dann smart <span class="time-badge">&asymp; 90 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

Der erste Booking-Service: l&auml;uft auf dem Entwicklerrechner, in der Cloud, auf einem Raspberry Pi &mdash; &uuml;berall gleich. Mit Health-Check (damit Probleme fr&uuml;h auffallen) und externer Konfiguration (damit ein Build in jeder Umgebung l&auml;uft).

#### User Story

Als <em>Betriebsteam</em> m&ouml;chte ich einen Service, der <em>in beliebigen Umgebungen ohne Code-&Auml;nderung l&auml;uft und sich automatisch &uuml;berwachen l&auml;sst</em>, damit <em>Probleme fr&uuml;h erkannt werden und Deployment-Pfade nicht spezifisch sein m&uuml;ssen</em>.

</div>

</div>
<div>

<div class="story-card">

#### Akzeptanzkriterien

- `GET /health` liefert 200 wenn der Service gesund ist
- Logs auf <code>stdout</code>
- Backend-URLs (Flight/Hotel/Car) konfigurierbar &uuml;ber Umgebungsvariablen
- `GET /info` zeigt die aktuelle Konfiguration
- `GET /openapi` liefert die OpenAPI-Spec aus &mdash; [Vorlage](https://github.com/larmic/workshop_microservices/blob/main/services/booking/story1/api/openapi.yaml)
- `Dockerfile` im Service-Verzeichnis baut ein lauff&auml;higes Image (Service lauscht auf Port `8080`)
- Pfad zum Service in `services/.env` als `CUSTOM_BOOKING_PATH` eingetragen, sodass `make docker-up-hub` den Service mitbaut
- `GET /booking/offers` aggregiert Flight, Hotel und Car und liefert kombinierte Ergebnisse

#### Setup-Hinweis (einmalig)

- `Dockerfile` im Service-Verzeichnis, Service lauscht auf Port `8080`
- `CUSTOM_BOOKING_PATH` in `services/.env` zeigt auf dieses Verzeichnis
- Backend-URLs (`FLIGHT_SERVICE_URL`, `HOTEL_SERVICE_URL`, `CAR_SERVICE_URL`, `CONSUL_URL`) kommen automatisch aus `docker-compose.custom.yml`, nicht selbst setzen
- Details: `docs/custom-setup.md`

</div>

</div>
</div>

Note:
- Hook zum Einsteigen vorlesen: &bdquo;Alle dachten, der Service l&auml;uft &mdash; bis ein Kunde anrief und fragte, warum seit drei Stunden nichts mehr geht. Vertrauen ist gut, ein Health-Endpoint ist besser.&ldquo;
- Wiedererkennung: dieselbe Karte (Kontext / User Story / Akzeptanzkriterien) findet ihr im Dashboard unter &bdquo;Story lesen&ldquo; &mdash; identischer Text, gleiches Layout.
- Sprache und Framework sind frei (Go, Java, Quarkus, Node, &hellip;). Die Referenz unter `services/booking/story1/` ist nur ein Go-Beispiel.
- Kein Error-Handling in Story 1 &mdash; das ist Absicht. Wir bauen das Skelett; Resilienz kommt in Stories 3&ndash;5.
- Setup-Hinweis steht absichtlich nur hier in Story 1. Ab Story 2 erweitern wir denselben Service iterativ, das Custom-Setup bleibt unver&auml;ndert.
- Time-Box 90 min inkl. Setup. Dashboard `http://localhost` zeigt den Story-1-Modus inkl. Spickzettel mit Pseudo-Code.
- Vollst&auml;ndige Aufgabenbeschreibung: `docs/stories/story-01-cloud-native-booking-service.md`.
