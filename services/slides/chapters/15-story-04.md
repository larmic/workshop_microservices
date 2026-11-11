<div class="page">

<p class="kicker">Story 4 &middot; Rollenspiel</p>

<div class="page-head">

## Wenn der Flug ausf&auml;llt

<span class="badge design">&asymp; 60 min <small>davon 20 Spiel</small></span>
</div>

<p class="subtitle">Ihr seid die Requests. Flight ist krank, die Antwort steht auf der <span class="hl">Karte</span>.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> ruft Nummern, gibt Karten aus, f&uuml;hrt das Board.</li>
<li><strong>Alle anderen</strong> sind Requests in der Schlange.</li>
<li><strong>Dunkle Karte:</strong> STOPP, wir reden.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<p class="step-rule"><strong>Bei jeder Karte:</strong> Nummer rufen, warten, Antwort vorlesen, ablegen, hinten anstellen.</p>
<ol>
<li><strong>Ohne Breaker, 01&ndash;03</strong> <span class="fragment">&rarr; 500, auch kein Hotel</span></li>
<li><strong>STOPP 1, das Board:</strong> Rot aufs Board. Platz 5 = OPEN.</li>
<li><strong>OPEN:</strong> Wie lange wartet man? <span class="fragment">&rarr; gar nicht</span></li>
<li><strong>STOPP 2, 30 s um:</strong> Gelb probt. Erst Timeout, dann CLOSED.</li>
<li><strong>31.02.:</strong> Kommt das aufs Board? <span class="fragment">&rarr; ja, bei uns schon</span></li>
</ol>
</div>
</div>

</div>

</div>

Note:
- <strong>Hook:</strong> &bdquo;Story 3 hat uns Services finden lassen. Jetzt haben wir einen gefunden, und er antwortet nicht.&ldquo;
- <strong>Setup:</strong> Ab 4 Personen, derselbe Stapel f&uuml;r jede Gruppengr&ouml;&szlig;e. Booking bekommt Stapel 01&ndash;27 (01 oben) und die drei Zustandskarten, keine Legende. Das Board kommt erst mit STOPP 1. Ein Tisch, ein Haufen. Der Trainer spielt nicht mit, er sagt die Regeln an und leitet die STOPPs. Kartentabelle: Trainer-Skript auf der Materialseite.
- <strong>Regeln ansagen</strong>, bei jeder Karte in dieser Reihenfolge:
  - Booking ruft die Nummer und gibt die Karte dem N&auml;chsten in der Schlange. Die n&auml;chste erst, wenn die vorige liegt.
  - Warten, was auf der Karte steht. Bei &bdquo;5 s&ldquo; oder &bdquo;3 s&ldquo; laut z&auml;hlen: &bdquo;einundzwanzig, zweiundzwanzig, dreiundzwanzig&ldquo;. Bei &bdquo;sofort&ldquo; nicht.
  - Die Antwort laut vorlesen.
  - Karte ablegen, hinten anstellen. Wei&szlig;e Karten kommen auf den Haufen.
  - Dunkle Karte: Booking ruft &bdquo;STOPP&ldquo; und gibt sie dem Trainer. Wer noch z&auml;hlt, z&auml;hlt fertig.
- <strong>STOPP 1 (Karte 04):</strong> &bdquo;Was habt ihr bekommen? Wie lange? Hotel und Car waren gesund, was hattet ihr davon?&ldquo; Story 3 bricht beim ersten Fehler mit 500 ab, Hotel und Car werden gar nicht gefragt.
  - &bdquo;Booking bekommt einen Circuit Breaker und 3 s Timeout.&ldquo; Board hinlegen, CLOSED ins Zustandsfeld. Das Board ist das Dashboard, das sieht der Betrieb.
  - Dazu ansagen: Rot hei&szlig;t, Flight hat einen Fehler geliefert. Rote Karten kommen aufs Board, auf den n&auml;chsten freien Platz. Liegt Platz 5, ruft Booking &bdquo;OPEN!&ldquo; und legt OPEN ins Zustandsfeld.
- <strong>STOPP 2 (Karte 13, Recap 02):</strong> Am Board: &bdquo;Welche Karten liegen hier? Ab welcher Nummer gab es <code>flights: []</code>? Was ist ab 10 anders? Was war der Breaker, was der Fallback?&ldquo;
  - Fallback ab dem ersten Fehler, auch bei CLOSED. Der Breaker spart nur das Warten: OPEN fragt Flight nicht, darum <code>X-Circuit-Open</code>. 09 hat Flight noch gefragt.
  - Einmal sagen: &bdquo;Das war &uuml;brigens immer HTTP 200 mit <code>X-Fallback</code>.&ldquo; Dann: &bdquo;30 Sekunden sind um.&ldquo; HALF_OPEN hinlegen.
  - Dazu ansagen: Gelb ist die Probe. Karte zuerst aufs Probe-Feld, dann z&auml;hlen und vorlesen, danach auf den Haufen. Booking: &bdquo;Timeout!&ldquo; hei&szlig;t OPEN, die 30 s laufen von vorn. &bdquo;Flight antwortet!&ldquo; hei&szlig;t Board abr&auml;umen, CLOSED.
- <strong>STOPP 3 (Karte 17, Recap 03):</strong> &bdquo;Wer lag auf dem Probe-Feld? Was bekamen 15 und 16, wie schnell? Was, wenn alle zu Flight gegangen w&auml;ren?&ldquo; Genau ein Probe-Slot (CompareAndSwap), der Rest sofort Fallback. Nach der gescheiterten Probe laufen die 30 s von vorn.
  - Dann: &bdquo;30 Sekunden sind um, Flight ist zur&uuml;ck.&ldquo; HALF_OPEN.
- <strong>STOPP 4 (Karte 20):</strong> Vorhersage, nicht aufl&ouml;sen: &bdquo;Flight ist gesund. Die n&auml;chsten haben sich vertippt, Datum 31.02. Kommen die aufs Board?&ldquo; Handzeichen, weiterspielen. 21&ndash;25 sind rot, also Board, keine neue Regel.
- <strong>ENDE (Karte 27, Recap 01 + 04):</strong> Board bleibt liegen. &bdquo;War Flight bei 21 bis 25 kaputt? Wer wurde bei 26 bestraft?&ldquo; 4xx ist kein Ausfall, unser Code z&auml;hlt alles ab 400.
  - &bdquo;Eine zweite Booking-Instanz hat ein eigenes Board. Was liegt darauf? Was wei&szlig; Booking nach einem Neustart?&ldquo; Der Zustand lebt im RAM einer Instanz. Geteilter Zustand hilft, kostet aber Kopplung.
- <strong>Demo:</strong>
  - Flight auf &bdquo;H&auml;ngt&ldquo; (10 s, Timeout 3 s): f&uuml;nf Calls warten je 3 s, immer 200 mit <code>flights: []</code> und <code>X-Fallback: flight</code>, ab dem 6. sofort mit <code>X-Circuit-Open: flight</code>. Schneller: &bdquo;Fehler&ldquo;.
  - Flight zur&uuml;ck auf &bdquo;Normal&ldquo;: nach 30 s HALF_OPEN, das Panel zeigt es schon ohne Call. Die erfolgreiche Probe schlie&szlig;t.
  - Details: <code>docs/stories/story-04-circuit-breaker.md</code>.
