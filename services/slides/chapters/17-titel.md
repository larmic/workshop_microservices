<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Bulkhead</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Ein langsamer Service ist gef&auml;hrlicher als ein <span class="accent">kaputter</span>.</p>
<p class="sub">Ein eigener Pool pro Backend. L&auml;uft einer voll, bleiben die anderen trocken.</p>
<figure class="foot-quote">
<blockquote>Kaskadierende Ausf&auml;lle entstehen an Ressourcen-Pools, aus denen kein Aufruf zur&uuml;ckkehrt.</blockquote>
<figcaption>Michael Nygard, <em>Release It!</em>, 2007 (sinngem&auml;&szlig;)</figcaption>
</figure>
</div>
<div class="split-figure">
<img src="./assets/bulkhead.svg" alt="Schiffsrumpf mit drei Abteilen f&uuml;r Flight, Hotel und Car, das Hotel-Abteil ist geflutet"/>
</div>
</div>

</div>


Note:
- <strong>Tag 2:</strong> Gestern war Flight tot, der Circuit Breaker hat geholfen. Heute ist Hotel nur langsam, und der Breaker merkt nichts.
- <strong>Frage:</strong> &bdquo;Was ist schlimmer: ein Backend, das sofort 500 liefert, oder eins, das 20 s braucht?&ldquo; Kurz abstimmen lassen.
  - Das schnelle 500 gibt den Slot sofort zur&uuml;ck.
  - Das langsame h&auml;lt ihn fest, bis der Pool leer ist.
- <strong>Weiter:</strong> Wie ein Schiff das l&ouml;st: Schotten.
