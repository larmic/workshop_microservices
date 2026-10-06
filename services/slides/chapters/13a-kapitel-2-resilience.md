<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark chapter">

<p class="kicker">Kapitel 2 von 3</p>

<div class="page-body">
<h1 class="chapter-title">Resilience</h1>
<p class="chapter-sub">Wenn Backends kippen.</p>
</div>

</div>

Note:
- Wir wechseln vom Fundament zum spannenden Teil. Stories 4 und 5 sind klassisches Resilience-Engineering: Circuit Breaker, Bulkhead, Timeout.
- Hook: &bdquo;Bis hierhin war es Hygiene. Jetzt wird es Architektur.&ldquo;
- Wichtigste Klammer f&uuml;r das Kapitel: Resilience-Patterns geh&ouml;ren in den <em>Aufrufer</em>, Schutz-Patterns in den <em>Aufgerufenen</em>. Wer das vermischt, sch&uuml;tzt nichts.
- Stories 6 und 7 sind streng genommen auch Resilience. Wir trennen sie in Kapitel 3, weil dort die Komplexit&auml;t nicht mehr ein Aufruf ist, sondern eine Kette.
- Tagesgrenze liegt mitten im Kapitel: Story 4 heute, Story 5 morgen fr&uuml;h.
