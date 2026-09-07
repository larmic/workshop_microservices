<div class="page">

<p class="kicker">Der Kanon</p>

<div class="page-head">

## Die 12 Faktoren

<div class="callout fragment" data-fragment-index="1">Diese drei baut ihr gleich in Story 1.</div>
</div>

<div class="page-body">

<!-- Ein Klick: die neun nicht behandelten Faktoren blenden ab (.dim), die drei
     aus Story 1 heben sich hervor (.pick), der Hinweis oben erscheint. -->
<div class="cards cards-4">
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">I</span> Codebase</h3>
<p>Ein Repo, viele Deployments.</p>
<code>1 Repo &rarr; dev &middot; stage &middot; prod</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">II</span> Dependencies</h3>
<p>Explizit deklarieren und isolieren.</p>
<code>go.mod &middot; package.json</code>
</div>
<div class="card fragment custom pick" data-fragment-index="1">
<h3><span class="num">III</span> Config</h3>
<p>Aus der Umgebung, nicht aus dem Code.</p>
<code>os.Getenv("DATABASE_URL")</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">IV</span> Backing Services</h3>
<p>Angeh&auml;ngte Ressourcen, per URL konfiguriert.</p>
<code>lokale DB &rarr; RDS: neue URL</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">V</span> Build, Release, Run</h3>
<p>Drei strikt getrennte Stufen.</p>
<code>build &rarr; release &rarr; run</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">VI</span> Processes</h3>
<p>Zustandslos. State geh&ouml;rt ins Backing Service.</p>
<code>Session &rarr; Redis</code>
</div>
<div class="card fragment custom pick" data-fragment-index="1">
<h3><span class="num">VII</span> Port Binding</h3>
<p>Der Service bringt seinen Server selbst mit.</p>
<code>listen(8080)</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">VIII</span> Concurrency</h3>
<p>Nach au&szlig;en skalieren, nicht nach oben.</p>
<code>mehr Prozesse statt mehr RAM</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">IX</span> Disposability</h3>
<p>Schnell starten, sauber beenden.</p>
<code>SIGTERM &rarr; sauber beenden</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">X</span> Dev / Prod Parity</h3>
<p>Dev und Prod so &auml;hnlich wie m&ouml;glich.</p>
<code>docker compose &asymp; Kubernetes</code>
</div>
<div class="card fragment custom pick" data-fragment-index="1">
<h3><span class="num">XI</span> Logs</h3>
<p>Event-Stream nach stdout, keine Logfiles.</p>
<code>log &rarr; stdout</code>
</div>
<div class="card fragment custom dim" data-fragment-index="1">
<h3><span class="num">XII</span> Admin Processes</h3>
<p>One-off-Prozesse in derselben Umgebung.</p>
<code>kubectl exec &hellip; db:migrate</code>
</div>
</div>

</div>

</div>

Note:
- Alle zw&ouml;lf Faktoren stehen sofort da, kurz durchgehen. Ein Klick blendet die neun ab, die Story 1 nicht abdeckt, und zeigt den Hinweis.
- Story 1 setzt drei Faktoren konkret um: Config aus Umgebungsvariablen mit <code>/info</code>-Endpoint (III), eigener HTTP-Server auf Port 8080 im Dockerfile (VII), Logs auf stdout (XI). Der Health-Check ist kein Faktor, sondern Betriebsthema; Disposability (sauberer Shutdown bei SIGTERM) wird im Workshop nicht behandelt. Dependencies (II), Build/Release/Run (V) und Dev/Prod Parity (X) erledigt das Docker-Setup nebenbei, ohne dass wir sie thematisieren.
- <strong>I Codebase:</strong> Eine Codebase pro App in Versionskontrolle, daraus viele Deploys (Dev/Staging/Prod). Mehrere Apps teilen NIE den Source — gemeinsamer Code wird als Library extrahiert. <strong>Anti-Pattern:</strong> Ein Service liegt in mehreren Repos verstreut, oder zwei Apps teilen den selben Source-Ordner per Symlink / git submodule.
- <strong>II Dependencies:</strong> Alle Abhängigkeiten explizit im Manifest deklarieren. Nicht auf System-Pakete verlassen. Isolation via Container oder venv. — pom.xml / package.json allein reichen NICHT: auch implizite Abhängigkeiten zählen. <strong>Anti-Pattern:</strong> App ruft <code>imagemagick</code>, <code>curl</code> oder <code>ffmpeg</code> als System-Binary — funktioniert auf dem Build-Server, fehlt im Prod-Image. Oder: „bei uns liegt das passende JAR halt im /opt/lib".
- <strong>III Config:</strong> Alles was sich zwischen Deploys unterscheidet (DB-URL, Secrets, Hostnames) kommt aus Umgebungsvariablen. Niemals Config-Dateien im Repo. <strong>Anti-Pattern:</strong> <code>application-prod.properties</code> mit Klartext-Passwort im Git. Oder if-else-Switch auf den Hostname: <code>if (host == "prod") …</code>.
- <strong>IV Backing Services:</strong> DB, Cache, Queue werden per URL angesprochen. Lokales Postgres oder Cloud-RDS ist nur ein Config-Switch — kein Code-Change. <strong>Anti-Pattern:</strong> Hardcoded <code>jdbc:postgresql://localhost:5432/myapp</code> im Code. Oder: Service-spezifische Treiber-Konfiguration, sodass ein Wechsel von Postgres zu MySQL ein Refactoring auslöst.
- <strong>V Build, Release, Run:</strong> Build erzeugt Artefakt aus Code. Release = Build + Config. Run führt aus. Strikt getrennt, kein „schnell auf Prod Code ändern". <strong>Anti-Pattern:</strong> SSH auf den Prod-Server, <code>git pull &amp;&amp; restart</code>. Oder: Hotfix direkt in der laufenden VM editiert — beim nächsten Deploy ist die Änderung weg.
- <strong>VI Processes:</strong> App-Prozesse sind stateless. Was persistent sein muss, geht in Backing Services. Ermöglicht einfache horizontale Skalierung. <strong>Anti-Pattern:</strong> User-Session als HashMap im Heap. Datei-Upload ins lokale <code>/tmp</code>, im nächsten Request wieder lesen. Sticky-Sessions am Loadbalancer als Workaround.
- <strong>VII Port Binding:</strong> Service stellt sich selbst über HTTP/TCP-Port bereit (embedded Server). Kein „in einen Tomcat/Apache deployen" — die App IST der Server. <strong>Anti-Pattern:</strong> WAR-File in einen extern verwalteten Tomcat deployen. Apache-vhost muss vom Ops-Team synchron gehalten werden, sonst läuft die App nicht.
- <strong>VIII Concurrency:</strong> Statt einen Riesen-Prozess zu skalieren, mehrere kleine Prozesse parallel starten. Unix-Process-Modell, horizontal statt vertikal. <strong>Anti-Pattern:</strong> Ein einziger JVM-Prozess mit 500-Thread-Pool und 32 GB Heap. Skalierung = größere Maschine (vertikal) statt mehr Instanzen.
- <strong>IX Disposability:</strong> Prozesse müssen jederzeit gestartet/gestoppt werden können. Schneller Boot, SIGTERM korrekt behandeln (laufende Requests sauber beenden). <strong>Anti-Pattern:</strong> 3-Minuten-Startzeit beim Booten (klassischer Spring-Monolith mit viel Hibernate-Init). SIGTERM ignoriert — Kubernetes killt mid-request, Daten gehen verloren.
- <strong>X Dev / Prod Parity:</strong> Gap zwischen Dev und Prod minimieren: gleiche Backing Services lokal wie produktiv (Docker hilft), kurze Zeit zwischen Commit und Deploy. <strong>Anti-Pattern:</strong> Dev nutzt H2 in-memory, Prod ist Oracle. Mail-Versand lokal als No-op gemockt, in Prod echter SMTP — Bug fällt erst auf Prod auf.
- <strong>XI Logs:</strong> App schreibt unstrukturiert auf stdout. Aggregation, Routing und Archivierung übernimmt die Plattform (ELK, Loki, CloudWatch). <strong>Anti-Pattern:</strong> Log4j-File-Appender schreibt nach <code>/var/log/myapp.log</code>, App rotiert selbst, beim Container-Neustart ist alles weg. Oder: eigener Log-Server, den die App direkt anspricht.
- <strong>XII Admin Processes:</strong> One-off-Tasks (DB-Migrationen, Konsole, Cleanup) laufen in derselben Umgebung wie die App — gleicher Code, gleiche Config, gleiche Dependencies. <strong>Anti-Pattern:</strong> DB-Migration per <code>psql prod-db -c "ALTER TABLE ..."</code> von Hand. Oder: separates Admin-Tool mit eigenem Build, eigener Config, anderer Lib-Version — wird auf Prod inkompatibel.
