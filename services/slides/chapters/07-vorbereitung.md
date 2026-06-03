## Vorbereitung

<p class="subtitle">Vorbereitung verifizieren <span class="time-badge">&asymp; 30 min</span></p>

<div class="chip-row">
  <span class="chip">Git</span>
  <span class="chip">Docker + Compose</span>
  <span class="chip">IDE</span>
  <span class="chip">Mini-Service (Hausaufgabe)</span>
</div>

<div class="box setup">

### In wenigen Schritten startklar

```bash
git clone https://github.com/larmic/workshop_microservices.git
cd workshop_microservices/services
cp .env.example .env      # wird später angepasst
make docker-up-hub        # zieht fertige Images von Docker Hub
open http://localhost     # Dashboard = Workshop-Steuerzentrale
```

<p class="qr-url">Voller Setup-Pfad: <a href="https://github.com/larmic/workshop_microservices/blob/main/docs/vorbereitung.md">docs/vorbereitung.md</a></p>
</div>

<img class="dashboard-image" src="./assets/dashboard-ui.png" alt="Workshop-Dashboard mit Flight-, Hotel-, Car- und Booking-Services"/>


Note:
- Das Setup ist Pflicht-Hausaufgabe vor dem Workshop (`docs/vorbereitung.md`, Abschnitt 0). Hier nur verifizieren, Zeitbudget ~30 min: Mini-Service läuft, Stack läuft. Restprobleme sofort einsammeln, niemand sitzt nachher allein vor einem roten Terminal.
- `make docker-up-hub` zieht fertige Images, kein lokaler Build nötig (= schnell). Variante B `make docker-up` baut alles selbst.
- Dashboard auf `http://localhost` ist die Schaltzentrale für alle Stories: Start/Stop, Health, Chaos-Schalter, Links zu Swagger / Consul / Traefik (8080).
- Was hochfährt: Flight, Hotel, Car, Booking + Consul (Service Discovery) + Swagger UI (API-Docs) + Traefik (Gateway) + Dashboard (UI).
- Voraussetzungen: Git, Docker (inkl. Compose), IDE und der vorab als Hausaufgabe aufgesetzte Greenfield-Mini-Service (Dockerfile, `/health`, Port 8080).
- Sprache/Framework für den Booking-Service ist **frei** (Go, Java, Quarkus, Node, …) — die Referenz hier ist nur ein Beispiel.
- Wer den vollen Setup-Pfad mit allen URLs und cURL-Health-Checks will, findet ihn in `docs/vorbereitung.md`.
