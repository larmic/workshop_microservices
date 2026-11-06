<div class="page">

<p class="kicker">Story 7 &middot; Rollenspiel</p>

<div class="page-head">

## Die Saga wird leise

<span class="badge design">&asymp; 25 min</span>
</div>

<p class="subtitle">Ihr seid die Kund:innen. Car sagt immer nein. Der <span class="hl">Becher</span> zeigt, was wirklich storniert wurde.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> bucht eine Reise nach der anderen. Sagt Car nein, storniert Booking r&uuml;ckw&auml;rts: Hotel, dann Flight.</li>
<li><strong>Hotel und Flight</strong> stornieren: bis f&uuml;nf z&auml;hlen, ein Gummib&auml;rchen in den Becher der Reise. Jedes Mal.</li>
<li><strong>Car</strong> sagt immer &bdquo;nein&ldquo;, in Runde 3 ist Car der Korb.</li>
<li><strong>Alle anderen</strong> haben je eine Reise und einen Becher. Nach jeder Runde <strong>Kassensturz</strong>: In jeden Becher geh&ouml;ren zwei.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li><strong>Runde 1, Anruf (wie Story 6):</strong> Booking wartet auf &bdquo;erledigt&ldquo;. Wer wartet? <span class="fragment">&rarr; die ganze Schlange, aber Booking wei&szlig; alles</span></li>
<li><strong>Runde 2, Zuruf (unsere Referenz):</strong> Das Backend ruft &bdquo;geh&ouml;rt!&ldquo;. Bleibt es still: bis drei z&auml;hlen, Strich, trotzdem &bdquo;storniert&ldquo;. Stimmt jeder Becher? <span class="fragment">&rarr; nein, und Booking hat es sogar gemerkt</span></li>
<li><strong>Runde 3, Korb:</strong> Booking wirft Zettel in den Korb. Der Korb liest vor, bis &bdquo;quittiert&ldquo; kommt, sonst noch einmal. Stimmt jeder Becher? <span class="fragment">&rarr; keiner fehlt, einer hat drei</span></li>
</ol>
</div>
</div>

</div>

</div>

Note:
- Hook: &bdquo;Story 6 habt ihr gebaut. Jetzt spielt ihr dieselbe Saga dreimal, jedes Mal ein bisschen leiser.&ldquo; Ab 12 Personen zwei Tische mit identischer Ansage. Jedes Fragment erst aufdecken, wenn beide Tische gez&auml;hlt haben.
- Material, Detailregeln (Reihenfolge, wann die n&auml;chste Reise dran ist) und Briefing: Story-Doc und Dashboard. Flight geheim: R2 bei Reise 2 das Storno &uuml;berh&ouml;ren, R3 bei Reise 2 stornieren, aber nicht quittieren. Die Zettelnummer regelt nur die Reihenfolge, Flight merkt sie sich nicht.
- R1 ist Story 6, R2 unsere Referenz: Zuruf = <code>POST /events/compensation</code>, &bdquo;geh&ouml;rt&ldquo; = 202, bis drei = 3-s-Timeout, Strich = <code>dispatch-failed</code>, trotzdem <code>COMPENSATED</code>.
- R3 ist Bonus: Korb = Broker mit Redelivery. &bdquo;storniert&ldquo; kam vor dem Gummib&auml;rchen, das ist Eventual Consistency. Die Referenz versucht genau einmal: keine Duplikate, keine Nachlieferung.
- Vereinfachungen: Car hat keinen Breaker, die Referenz bedient Anfragen parallel. Das Z&auml;hlen bis f&uuml;nf macht das Storno sp&uuml;rbar, unsere Backends loggen nur.
- &Uuml;berleitung: Becher stehen lassen, die Leiter sortiert sie. Die Demo steht in der Note dort.
