## Story 2

<p class="subtitle">Services finden sich selbst <span class="time-badge">&asymp; 60 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

Statische URLs sind eine sch&ouml;ne Idee f&uuml;r Folien, aber nicht f&uuml;r den Betrieb. Backend-Services registrieren sich bei <em>Consul</em>; der Booking-Service findet sie &uuml;ber logische Namen statt fest verdrahteter Adressen. F&auml;llt eine Instanz aus, &uuml;bernimmt automatisch eine andere.

#### User Story

Als <em>Betriebsteam</em> m&ouml;chte ich <em>Backend-Services &uuml;ber logische Namen finden, statt URLs per Hand zu pflegen</em>, damit <em>Skalierung und Ausf&auml;lle ohne n&auml;chtliche Config-Edits funktionieren</em>.

</div>

</div>
<div>

<div class="story-card">

#### Akzeptanzkriterien

- Resolver fragt `GET /v1/health/service/{name}?passing=true` bei Consul
- Statische URLs ersetzt durch logische Namen: `flight-service`, `hotel-service`, `car-service`
- Bestehender `GET /booking/offers` aus Story 1 umgebaut, l&ouml;st Flight/Hotel/Car &uuml;ber den Consul-Resolver auf statt &uuml;ber statische Env-URLs
- `POST /booking/bookings` orchestriert die drei Backend-Services und liefert eine aggregierte Buchung
- Optional: Client-Side Load Balancing bei mehreren Instanzen (zuf&auml;llige Auswahl)

</div>

</div>
</div>

Note:
- Hook: &bdquo;Story 1 hatte URLs in ENV-Variablen. Was passiert, wenn Flight umzieht? Re-Deploy. Was passiert, wenn ihr Flight skaliert? Ein Backend bekommt allen Traffic.&ldquo;
- Wiedererkennung: identische Karte (Kontext / User Story / Akzeptanzkriterien) im Dashboard unter Story 2 &rarr; &bdquo;Story lesen&ldquo;.
- Sprache und Framework wieder frei. Referenz unter `services/booking/story2/` (Go).
- Self-Registration im Code &mdash; Flight/Hotel/Car melden sich beim Start aktiv an. Trade-off zu Sidecar / Plattform-basierter Registrierung diskutieren wir im Recap.
- Client-Side LB: Resolver bekommt eine Liste, w&auml;hlt eine Instanz zuf&auml;llig. Keine Load-Balancer-Magie n&ouml;tig.
- Time-Box 60 min. Dashboard `http://localhost` zeigt den Story-2-Modus inkl. Spickzettel.
- Vollst&auml;ndige Aufgabenbeschreibung: `docs/stories/story-02-service-discovery.md`.
