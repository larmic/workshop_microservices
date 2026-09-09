<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Service Discovery</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Wer die Adresse fest verdrahtet, ruft irgendwann <span class="accent">ins Leere</span>.</p>
<p class="sub">Services finden sich &uuml;ber Namen, nicht &uuml;ber URLs.</p>
<p class="source"><code>&bdquo;Topology doesn&rsquo;t change.&ldquo;</code> &middot; die f&uuml;nfte Fallacy of Distributed Computing</p>
</div>
<div class="split-figure">
<img src="./assets/service-discovery.svg" alt="Ein Verzeichnis l&ouml;st einen Namen auf eine von mehreren Instanzen auf; die fr&uuml;here Instanz ist ausgefallen"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;In Story 1 standen die Backend-URLs in ENV-Variablen. Was passiert, wenn Flight umzieht, wenn ihr eine Instanz dazuskaliert, wenn der Container neu startet und eine andere IP zieht?&ldquo;
- Die Grafik: Ein Verzeichnis (links) kennt drei Eintr&auml;ge, der aktive ist Lime. Der Aufruf geht &uuml;ber den Namen an eine gesunde Instanz (gr&uuml;ner Punkt); die alte Instanz oben rechts ist weg, der gestrichelte Weg dorthin f&uuml;hrt ins Leere.
- Die Fallacies of Distributed Computing (Peter Deutsch, Sun, 1994): Nummer f&uuml;nf lautet &bdquo;Topology doesn&rsquo;t change&ldquo;. Genau diese Annahme steckt in jeder hartkodierten URL.
- Diskussions-Anker: Wer pflegt heute noch URLs per Hand, in YAML, ConfigMap oder Wiki? Wann hat das zuletzt Probleme gemacht?
- &Uuml;bergang zur Karten-Folie: &bdquo;F&uuml;nf Bausteine, die das Problem in den Griff bekommen.&ldquo;
