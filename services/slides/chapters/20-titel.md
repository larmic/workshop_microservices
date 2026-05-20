<!-- .slide: data-background-image="./assets/saga.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.40" data-background-repeat="no-repeat" -->

## Saga

<p class="subtitle">Alles oder nichts &mdash; aber richtig</p>

Note:
- Hook: &bdquo;Stories 3 und 4 helfen <em>einem</em> Aufruf. Aber sobald ich mehrere Schritte habe (Flug + Hotel + Auto) und einer kippt mitten drin, brauche ich was anderes. Das Hotel sagt nein, der Flug ist schon gebucht &mdash; was nun?&ldquo;
- Demo-Vorschau: Hotel auf &bdquo;Fehler&ldquo;, dann <code>POST /booking/bookings</code> &mdash; im Dashboard wandert die Saga durch <code>PENDING</code> &rarr; <code>COMPENSATING</code> &rarr; <code>FAILED</code>, und der schon gebuchte Flug wird sichtbar storniert.
- &Uuml;bergang zur Karten-Slide: &bdquo;F&uuml;nf Bausteine, die uns Konsistenz &uuml;ber Service-Grenzen geben &mdash; ohne 2PC.&ldquo;
