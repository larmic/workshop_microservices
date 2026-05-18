## Die 12 Faktoren

<div class="factor-grid">

<div class="factor fragment">
<h3><span class="numeral">I</span> Codebase</h3>
<p>Ein Repo, viele Deployments. Kein geteilter Code zwischen Apps.</p>
<code>Git monorepo → dev / staging / prod</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">II</span> Dependencies</h3>
<p>Abhängigkeiten explizit deklarieren und isolieren.</p>
<code>go.mod · package.json · requirements.txt</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">III</span> Config</h3>
<p>Konfiguration aus der Umgebung — niemals aus dem Code.</p>
<code>DB_URL = os.Getenv("DATABASE_URL")</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">IV</span> Backing Services</h3>
<p>DB, Queue, Cache sind austauschbare Ressourcen.</p>
<code>postgres:// → mysql:// (nur Config-Switch)</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">V</span> Build, Release, Run</h3>
<p>Drei strikt getrennte Stufen — kein Code-Edit in Prod.</p>
<code>build → release (+ config) → run</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">VI</span> Processes</h3>
<p>Stateless. State gehört in Backing Services.</p>
<code>Session → Redis, nicht in-memory</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">VII</span> Port Binding</h3>
<p>Service exportiert sich selbst über einen Port.</p>
<code>app.listen(8080) — kein App-Server davor</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">VIII</span> Concurrency</h3>
<p>Skalierung horizontal über Prozesse, nicht Threads.</p>
<code>kubectl scale --replicas=10</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">IX</span> Disposability</h3>
<p>Schneller Start, sauberer Shutdown bei SIGTERM.</p>
<code>signal.Notify(c, SIGTERM)</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">X</span> Dev / Prod Parity</h3>
<p>Dev und Prod so ähnlich wie möglich — gleiche Backing Services.</p>
<code>docker compose ≈ Kubernetes</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">XI</span> Logs</h3>
<p>Logs als Event-Stream nach stdout — keine Logfiles.</p>
<code>log.Println(...) → stdout → ELK</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">XII</span> Admin Processes</h3>
<p>Admin-Tasks als One-off Prozesse, gleiche Umgebung.</p>
<code>kubectl exec … db:migrate</code>
</div>

</div>

Note:
- **I Codebase** — Eine Codebase pro App in Versionskontrolle, daraus viele Deploys (Dev/Staging/Prod). Mehrere Apps teilen NIE den Source — wenn gemeinsamer Code, dann als Library extrahieren.
- **II Dependencies** — Alle Abhängigkeiten explizit im Manifest deklarieren (go.mod, package.json, requirements.txt). Nicht auf System-Pakete verlassen. Isolation via Container oder venv.
- **III Config** — Alles was sich zwischen Deploys unterscheidet (DB-URL, Secrets, Hostnames) kommt aus Umgebungsvariablen. Niemals Config-Dateien im Repo.
- **IV Backing Services** — DB, Cache, Queue werden per URL angesprochen. Lokales Postgres oder Cloud-RDS ist nur ein Config-Switch — kein Code-Change.
- **V Build, Release, Run** — Build erzeugt Artefakt aus Code. Release = Build + Config. Run führt aus. Strikt getrennt, kein „schnell auf Prod Code ändern".
- **VI Processes** — App-Prozesse sind stateless. Was persistent sein muss, geht in Backing Services. Ermöglicht einfache horizontale Skalierung.
- **VII Port Binding** — Service stellt sich selbst über HTTP/TCP-Port bereit (embedded Server). Kein „in einen Tomcat/Apache deployen" — die App IST der Server.
- **VIII Concurrency** — Statt einen Riesen-Prozess zu skalieren, mehrere kleine Prozesse parallel starten. Unix-Process-Modell, horizontal statt vertikal.
- **IX Disposability** — Prozesse müssen jederzeit gestartet/gestoppt werden können. Schneller Boot, SIGTERM korrekt behandeln (laufende Requests sauber beenden).
- **X Dev / Prod Parity** — Gap zwischen Dev und Prod minimieren: gleiche Backing Services lokal wie produktiv (Docker hilft), kurze Zeit zwischen Commit und Deploy.
- **XI Logs** — App schreibt unstrukturiert auf stdout. Aggregation, Routing und Archivierung übernimmt die Plattform (ELK, Loki, CloudWatch).
- **XII Admin Processes** — One-off-Tasks (DB-Migrationen, Konsole, Cleanup) laufen in derselben Umgebung wie die App — gleicher Code, gleiche Config, gleiche Dependencies.
