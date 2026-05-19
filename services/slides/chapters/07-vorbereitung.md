## Vorbereitung

<p class="subtitle">Setup gemeinsam <span class="time-badge">&asymp; 45 min</span></p>

<div class="chip-row">
  <span class="chip">Git</span>
  <span class="chip">Docker + Compose</span>
  <span class="chip">IDE (optional)</span>
</div>

<div class="box">

### In 3 Schritten startklar

```bash
git clone https://github.com/larmic/workshop_microservices.git
cd workshop_microservices/services
make docker-up-hub        # zieht fertige Images von Docker Hub
open http://localhost     # Dashboard = Workshop-Steuerzentrale
```

</div>

<p class="subtitle">Was im Hintergrund hochf&auml;hrt</p>

<div class="chip-row">
  <span class="chip brand">Docker</span>
  <span class="chip brand">Services &middot; Flight &middot; Hotel &middot; Car &middot; Booking</span>
  <span class="chip brand">Consul</span>
  <span class="chip brand">Swagger UI</span>
  <span class="chip brand">OpenAPI</span>
</div>

&rarr; [docs/vorbereitung.md](https://github.com/larmic/workshop_microservices/blob/main/docs/vorbereitung.md)

Note:
- Wir machen das Setup gemeinsam, Zeitbudget ~45 min — niemand sitzt nachher allein vor einem roten Terminal.
- `make docker-up-hub` zieht fertige Images, kein lokaler Build nötig (= schnell). Variante B `make docker-up` baut alles selbst.
- Dashboard auf `http://localhost` ist die Schaltzentrale für alle Stories: Start/Stop, Health, Links zu Swagger / Consul / Traefik (8080).
- Voraussetzungen: Git, Docker (inkl. Compose). IDE ist optional und nur nötig, wenn man die Stories selbst implementiert.
- Sprache/Framework für den Booking-Service ist **frei** (Go, Java, Quarkus, Node, …) — die Referenz hier ist nur ein Beispiel.
- Wer den vollen Setup-Pfad mit allen URLs und cURL-Health-Checks will, findet ihn in `docs/vorbereitung.md`.
