<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Service Discovery</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Wer die Adresse fest verdrahtet, ruft irgendwann <span class="accent">ins Leere</span>.</p>
<p class="sub">Services finden sich &uuml;ber Namen, nicht &uuml;ber URLs.</p>
<figure class="foot-quote">
<blockquote>Topology doesn&rsquo;t change.</blockquote>
<figcaption>Irrtum Nr. 5 der <em>Fallacies of Distributed Computing</em>, L Peter Deutsch u.&nbsp;a., Sun Microsystems, 1994</figcaption>
</figure>
</div>
<div class="split-figure">
<img src="./assets/service-discovery.svg" alt="Ein Verzeichnis l&ouml;st einen Namen auf eine von mehreren Instanzen auf; die fr&uuml;here Instanz ist ausgefallen"/>
</div>
</div>

</div>


Note:
- <strong>Hook:</strong> &bdquo;Erinnert ihr euch an Frage 2 im Story-1-Recap? Vielleicht ist Flight nicht tot, sondern umgezogen.&ldquo;
- Container starten neu, skalieren, ziehen um: Die Adresse &auml;ndert sich st&auml;ndig.
- <strong>Herkunft:</strong> Die Fallacies sind acht Annahmen, die jeder beim ersten verteilten System trifft und die alle falsch sind. L Peter Deutsch und Kollegen bei Sun haben sie in den 90ern gesammelt.
  - Nr. 5 hei&szlig;t: Die Topologie &auml;ndert sich nicht. Bei Containern &auml;ndert sie sich st&auml;ndig.
- <strong>&Uuml;berleitung:</strong> Ein Verzeichnis l&ouml;st den Namen auf eine lebende Instanz auf. Was so eine Registry leistet, kommt jetzt.
