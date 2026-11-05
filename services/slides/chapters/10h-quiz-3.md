<div class="page">

<p class="kicker">Quiz &middot; 3 von 5</p>

## RESTful oder nicht?

<div class="page-body quiz">

<div class="quiz-code"><span class="m">GET</span> /booking/cancelBooking?id=4711</div>

<div class="quiz-answer fragment">
<div class="callout">Nicht RESTful</div>
<p>Verb im Pfad, Identit&auml;t als Query-Parameter, und ein GET, das etwas ver&auml;ndert: drei Regeln auf einmal. Wie es stattdessen aussieht, habt ihr gerade selbst entworfen.</p>
<code class="quiz-take">Aktion &ne; Pfad</code>
</div>

</div>

</div>

Note:
- <strong>Mechanik:</strong> Daumen wie gehabt. Erwartung: fast alle sagen nein. Wer Ja sagt, darf begr&uuml;nden, dann aufl&ouml;sen.
- <strong>Kern:</strong> Verb im Pfad, ID als Query-Parameter, GET mit Seiteneffekt. Am schwersten wiegt Letzteres.
- <strong>Anekdote:</strong> 2005 hat der Google Web Accelerator Links vorgeladen. Bei 37signals Backpack waren L&ouml;schen-Links einfache GETs, der Prefetcher hat Nutzern Daten gel&ouml;scht. &bdquo;GET ver&auml;ndert nichts&ldquo; ist Selbstschutz, keine Stilfrage.
- <strong>Bezug:</strong> Wer hat am Flipchart etwas &Auml;hnliches stehen?
- <strong>Weiter:</strong> &bdquo;Zwei Varianten, die es besser machen. Welche ist RESTful?&ldquo;
