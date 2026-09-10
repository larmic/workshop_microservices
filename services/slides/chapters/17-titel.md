<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark">

<p class="kicker">Bulkhead</p>

<div class="page-body split">
<div class="split-text">
<p class="statement">Ein langsamer Service ist gef&auml;hrlicher als ein <span class="accent">kaputter</span>.</p>
<p class="sub">Ein eigener Pool pro Backend. L&auml;uft einer voll, bleiben die anderen trocken.</p>
<p class="source">Michael Nygard, <em>Release It!</em>: Kaskadierende Ausf&auml;lle entstehen an Ressourcen-Pools, aus denen kein Aufruf zur&uuml;ckkehrt.</p>
</div>
<div class="split-figure">
<img src="./assets/bulkhead.svg" alt="Schiffsrumpf mit f&uuml;nf Schotten, das zweite Abteil ist geflutet"/>
</div>
</div>

</div>

Note:
- Hook: &bdquo;Story 4 sch&uuml;tzt gegen <em>kaputte</em> Backends. Aber was, wenn das Backend gar nicht kaputt ist, nur langsam? Hotel antwortet in 2 s, fehlerfrei. Der Circuit Breaker bleibt CLOSED. Trotzdem h&auml;ngen unter Last alle Threads in Hotel-Calls fest, Flight und Car bekommen nichts mehr durch.&ldquo;
- Die Grafik: f&uuml;nf Schotten, das zweite Abteil l&auml;uft voll, die anderen vier bleiben trocken. Genau das soll der Pool pro Backend leisten.
- Nygard, Release It! (2007, 2. Auflage 2018): Kaskadierende Ausf&auml;lle entstehen fast immer an Ressourcen-Pools, deren Aufrufe nicht zur&uuml;ckkommen, Threads, Connections, Sockets. Bulkhead ist seine Antwort darauf.
- Provokation f&uuml;r sp&auml;ter: &bdquo;Async-Aufrufe l&ouml;sen das nicht. Eine Goroutine, die an HTTP h&auml;ngt, ist genauso teuer wie ein blockierter Thread.&ldquo; Kommt im Bonus nach dem Recap.
