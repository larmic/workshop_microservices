## Die 12 Faktoren

<div class="factor-grid">

<div class="factor fragment">
<h3><span class="numeral">I</span> Codebase</h3>
<p>Ein Repo, viele Deployments. Kein geteilter Code zwischen Apps.</p>
<code>Git monorepo → dev / staging / prod</code>
<aside class="notes">Eine Codebase pro App in Versionskontrolle, daraus viele Deploys (Dev/Staging/Prod). Mehrere Apps teilen NIE den Source — gemeinsamer Code wird als Library extrahiert.<br><strong>Anti-Pattern:</strong> Ein Service liegt in mehreren Repos verstreut, oder zwei Apps teilen den selben Source-Ordner per Symlink / git submodule.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">II</span> Dependencies</h3>
<p>Abhängigkeiten explizit deklarieren und isolieren.</p>
<code>go.mod · package.json · requirements.txt</code>
<aside class="notes">Alle Abhängigkeiten explizit im Manifest deklarieren. Nicht auf System-Pakete verlassen. Isolation via Container oder venv. — pom.xml / package.json allein reichen NICHT: auch implizite Abhängigkeiten zählen.<br><strong>Anti-Pattern:</strong> App ruft <code>imagemagick</code>, <code>curl</code> oder <code>ffmpeg</code> als System-Binary — funktioniert auf dem Build-Server, fehlt im Prod-Image. Oder: „bei uns liegt das passende JAR halt im /opt/lib".</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">III</span> Config</h3>
<p>Konfiguration aus der Umgebung — niemals aus dem Code.</p>
<code>DB_URL = os.Getenv("DATABASE_URL")</code>
<aside class="notes">Alles was sich zwischen Deploys unterscheidet (DB-URL, Secrets, Hostnames) kommt aus Umgebungsvariablen. Niemals Config-Dateien im Repo.<br><strong>Anti-Pattern:</strong> <code>application-prod.properties</code> mit Klartext-Passwort im Git. Oder if-else-Switch auf den Hostname: <code>if (host == "prod") …</code>.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">IV</span> Backing Services</h3>
<p>DB, Queue, Cache sind austauschbare Ressourcen.</p>
<code>postgres:// → mysql:// (nur Config-Switch)</code>
<aside class="notes">DB, Cache, Queue werden per URL angesprochen. Lokales Postgres oder Cloud-RDS ist nur ein Config-Switch — kein Code-Change.<br><strong>Anti-Pattern:</strong> Hardcoded <code>jdbc:postgresql://localhost:5432/myapp</code> im Code. Oder: Service-spezifische Treiber-Konfiguration, sodass ein Wechsel von Postgres zu MySQL ein Refactoring auslöst.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">V</span> Build, Release, Run</h3>
<p>Drei strikt getrennte Stufen — kein Code-Edit in Prod.</p>
<code>build → release (+ config) → run</code>
<aside class="notes">Build erzeugt Artefakt aus Code. Release = Build + Config. Run führt aus. Strikt getrennt, kein „schnell auf Prod Code ändern".<br><strong>Anti-Pattern:</strong> SSH auf den Prod-Server, <code>git pull &amp;&amp; restart</code>. Oder: Hotfix direkt in der laufenden VM editiert — beim nächsten Deploy ist die Änderung weg.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">VI</span> Processes</h3>
<p>Stateless. State gehört in Backing Services.</p>
<code>Session → Redis, nicht in-memory</code>
<aside class="notes">App-Prozesse sind stateless. Was persistent sein muss, geht in Backing Services. Ermöglicht einfache horizontale Skalierung.<br><strong>Anti-Pattern:</strong> User-Session als HashMap im Heap. Datei-Upload ins lokale <code>/tmp</code>, im nächsten Request wieder lesen. Sticky-Sessions am Loadbalancer als Workaround.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">VII</span> Port Binding</h3>
<p>Service exportiert sich selbst über einen Port.</p>
<code>app.listen(8080) — kein App-Server davor</code>
<aside class="notes">Service stellt sich selbst über HTTP/TCP-Port bereit (embedded Server). Kein „in einen Tomcat/Apache deployen" — die App IST der Server.<br><strong>Anti-Pattern:</strong> WAR-File in einen extern verwalteten Tomcat deployen. Apache-vhost muss vom Ops-Team synchron gehalten werden, sonst läuft die App nicht.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">VIII</span> Concurrency</h3>
<p>Skalierung horizontal über Prozesse, nicht Threads.</p>
<code>kubectl scale --replicas=10</code>
<aside class="notes">Statt einen Riesen-Prozess zu skalieren, mehrere kleine Prozesse parallel starten. Unix-Process-Modell, horizontal statt vertikal.<br><strong>Anti-Pattern:</strong> Ein einziger JVM-Prozess mit 500-Thread-Pool und 32 GB Heap. Skalierung = größere Maschine (vertikal) statt mehr Instanzen.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">IX</span> Disposability</h3>
<p>Schneller Start, sauberer Shutdown bei SIGTERM.</p>
<code>signal.Notify(c, SIGTERM)</code>
<aside class="notes">Prozesse müssen jederzeit gestartet/gestoppt werden können. Schneller Boot, SIGTERM korrekt behandeln (laufende Requests sauber beenden).<br><strong>Anti-Pattern:</strong> 3-Minuten-Startzeit beim Booten (klassischer Spring-Monolith mit viel Hibernate-Init). SIGTERM ignoriert — Kubernetes killt mid-request, Daten gehen verloren.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">X</span> Dev / Prod Parity</h3>
<p>Dev und Prod so ähnlich wie möglich — gleiche Backing Services.</p>
<code>docker compose ≈ Kubernetes</code>
<aside class="notes">Gap zwischen Dev und Prod minimieren: gleiche Backing Services lokal wie produktiv (Docker hilft), kurze Zeit zwischen Commit und Deploy.<br><strong>Anti-Pattern:</strong> Dev nutzt H2 in-memory, Prod ist Oracle. Mail-Versand lokal als No-op gemockt, in Prod echter SMTP — Bug fällt erst auf Prod auf.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">XI</span> Logs</h3>
<p>Logs als Event-Stream nach stdout — keine Logfiles.</p>
<code>log.Println(...) → stdout → ELK</code>
<aside class="notes">App schreibt unstrukturiert auf stdout. Aggregation, Routing und Archivierung übernimmt die Plattform (ELK, Loki, CloudWatch).<br><strong>Anti-Pattern:</strong> Log4j-File-Appender schreibt nach <code>/var/log/myapp.log</code>, App rotiert selbst, beim Container-Neustart ist alles weg. Oder: eigener Log-Server, den die App direkt anspricht.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">XII</span> Admin Processes</h3>
<p>Admin-Tasks als One-off Prozesse, gleiche Umgebung.</p>
<code>kubectl exec … db:migrate</code>
<aside class="notes">One-off-Tasks (DB-Migrationen, Konsole, Cleanup) laufen in derselben Umgebung wie die App — gleicher Code, gleiche Config, gleiche Dependencies.<br><strong>Anti-Pattern:</strong> DB-Migration per <code>psql prod-db -c "ALTER TABLE ..."</code> von Hand. Oder: separates Admin-Tool mit eigenem Build, eigener Config, anderer Lib-Version — wird auf Prod inkompatibel.</aside>
</div>

</div>

<span class="show-all fragment" aria-hidden="true"></span>
