<div class="page">

<p class="kicker">Microservices sind auch nur Monolithen</p>

## &hellip; mit <span class="hl">Netzwerkproblemen</span>

<div class="page-body">

<p class="lead">&rarr; und ein Methodenaufruf wird zum Netzwerkaufruf</p>

<div class="cards cards-3 cards-display">
<div class="card">
<h3>Langsam</h3>
<p>blockiert alles, was hinter ihm wartet</p>
</div>
<div class="card">
<h3>L&uuml;gt</h3>
<p>200 zur&uuml;ck, obwohl nichts mehr geht</p>
</div>
<div class="card">
<h3>Kommt nie an</h3>
<p>niemand wei&szlig;, ob etwas passiert ist</p>
</div>
</div>

</div>

</div>

Note:
- These provokant in den Raum stellen: stimmt das? Wo trifft sie zu, wo nicht?
- Der Kern: Ein Methodenaufruf im Monolithen ist schnell, ehrlich und kommt an. Ein Netzwerkaufruf kann langsam sein, l&uuml;gen (200 trotz Fehler) oder verschwinden. Alles, was danach kommt, ist der Umgang mit diesen drei F&auml;llen.
- Diskussions-Anker: Welche Probleme h&auml;tte man im Monolithen auch, und welche entstehen erst durch das Netzwerk?
- Weiter nach unten: die zwei Fragen, die der Workshop beantwortet.
