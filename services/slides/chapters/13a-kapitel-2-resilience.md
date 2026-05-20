<!-- .slide: class="chapter-slide" data-background-color="#0A0349" -->

<p class="chapter-num">Kapitel 2</p>

## Resilience

<p class="subtitle">Wenn Backends kippen.</p>

Note:
- Wir wechseln vom Fundament zum spannenden Teil. Stories 3 und 4 sind klassisches Resilience-Engineering &mdash; Circuit Breaker, Bulkhead, Timeout.
- Hook: &bdquo;Bis hierhin war's Hygiene. Jetzt wird's Architektur.&ldquo;
- Wichtigste Klammer f&uuml;rs Kapitel: Resilience-Patterns geh&ouml;ren in den <em>Aufrufer</em>, Schutz-Patterns in den <em>Aufgerufenen</em>. Wer das vermischt, sch&uuml;tzt nichts.
- Stories 5/6 sind streng genommen auch Resilience &mdash; wir trennen sie in Kapitel 3, weil dort die Komplexit&auml;t nicht mehr ein Aufruf, sondern eine Kette ist.
