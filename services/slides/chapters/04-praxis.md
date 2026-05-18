## Brauchen wir Microservices?

<p class="subtitle">TN Skills — gestartet mit Microservices, heute Modulith</p>

<div class="cols">
<div>

### Damals

- 2 Teams, je ein Service (fast 3 mit C#)
- PHP & Kotlin
- REST zwischen Services
- Eigene Pipelines & Deployments
- **~14 Leute** — kein Aushelfen über Teams

</div>
<div>

### Heute

- 1 Modulith, 4 Domänen
- 1 Pipeline, 1 Deployment
- **5 Leute**
- (UI war initial separat)

</div>
</div>

<div class="box">

- Microservices sind <span class="hl">nicht immer</span> die richtige Wahl
- Mit einem <span class="hl">Monolithen zu starten</span> ist meist die bessere Option
- Dieser Workshop zeigt das <span class="hl">wie</span> — nicht das <span class="hl">ob</span>

</div>

Note:
- TN Skills: 2 Teams, je ein Service; Reporting-Service in C# wurde angefangen, aber abgebrochen — Komplexität für die Team-Größe zu hoch.
- Heute: gleiches Produkt, 5 statt 14 Leute, ein Modulith. Deployment & Betrieb deutlich einfacher, Feature-Velocity gestiegen.
- Take-away für den Workshop: Microservices sind ein Werkzeug, kein Ziel. Conway's Law (Org-Struktur dominiert Architektur) ist die einzige zwingend gute Begründung.
