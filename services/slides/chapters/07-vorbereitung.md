## Vorbereitung

<p class="subtitle">Setup gemeinsam <span class="time-badge">&asymp; 45 min</span></p>

<div class="chip-row">
  <span class="chip">Git</span>
  <span class="chip">Docker + Compose</span>
  <span class="chip">IDE (optional)</span>
</div>

<div class="box">

### In wenigen Schritten startklar

```bash
git clone https://github.com/larmic/workshop_microservices.git
cp .env.example .env      # wird später angepasst
cd workshop_microservices/services
make docker-up-hub        # zieht fertige Images von Docker Hub
open http://localhost     # Dashboard = Workshop-Steuerzentrale

(https://github.com/larmic/workshop_microservices/blob/main/docs/vorbereitung.md)
```
</div>

<img class="dashboard-image" src="./assets/dashboard-ui.png" alt="Workshop-Dashboard mit Flight-, Hotel-, Car- und Booking-Services"/>


Note:
- Wir machen das Setup gemeinsam, Zeitbudget ~45 min — niemand sitzt nachher allein vor einem roten Terminal.
- `make docker-up-hub` zieht fertige Images, kein lokaler Build nötig (= schnell). Variante B `make docker-up` baut alles selbst.
- Dashboard auf `http://localhost` ist die Schaltzentrale für alle Stories: Start/Stop, Health, Chaos-Schalter, Links zu Swagger / Consul / Traefik (8080).
- Was hochfährt: Flight, Hotel, Car, Booking + Consul (Service Discovery) + Swagger UI (API-Docs) + Traefik (Gateway) + Dashboard (UI).
- Voraussetzungen: Git, Docker (inkl. Compose). IDE ist optional und nur nötig, wenn man die Stories selbst implementiert.
- Sprache/Framework für den Booking-Service ist **frei** (Go, Java, Quarkus, Node, …) — die Referenz hier ist nur ein Beispiel.
- Wer den vollen Setup-Pfad mit allen URLs und cURL-Health-Checks will, findet ihn in `docs/vorbereitung.md`.
