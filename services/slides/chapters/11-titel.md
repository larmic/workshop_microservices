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
- <strong>Hook:</strong> &bdquo;Erinnert ihr euch an Frage 2 aus dem Story-1-Recap? Vielleicht ist Flight gar nicht tot, sondern hat nur eine neue Adresse.&ldquo;
- <strong>Kernbotschaft:</strong> Bisher steht die URL fest in der Konfiguration. Container starten neu, skalieren, ziehen um: Die Topologie &auml;ndert sich st&auml;ndig. Wer sich die Adresse merkt, ruft ins Leere.
- <strong>&Uuml;berleitung:</strong> Die Grafik zeigt die Idee: Ein Verzeichnis l&ouml;st den Namen auf eine lebende Instanz auf. Was so eine Registry leistet, kommt jetzt.
