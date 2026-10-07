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

