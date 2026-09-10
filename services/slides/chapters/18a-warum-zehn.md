<div class="page">

<p class="kicker">Story 5 &middot; &Uuml;bung</p>

<div class="page-head">

## Warum ausgerechnet zehn?

<span class="badge">&asymp; 20 min</span>
</div>

<div class="page-body littles">

<div class="formula">L = <span class="v">&lambda;</span> &times; <span class="v">W</span></div>

<p class="formula-text">L = gleichzeitig laufende Aufrufe, also ben&ouml;tigte Slots. &lambda; = Aufrufe pro Sekunde. W = Dauer eines Aufrufs. Der Booking-Service schickt <strong>20 Hotel-Aufrufe pro Sekunde</strong>.</p>

<div class="ltable">
<div class="ltable-head">Hotel antwortet in</div>
<div class="ltable-head">&lambda;</div>
<div class="ltable-head">W</div>
<div class="ltable-head">L = Slots</div>
<div>50 ms</div><div><code>20 / s</code></div><div><code>0,05 s</code></div><div class="swap"><code class="q fragment fade-out" data-fragment-index="1">?</code><code class="fragment" data-fragment-index="1">1</code></div>
<div>500 ms</div><div><code>20 / s</code></div><div><code>0,5 s</code></div><div class="swap"><code class="q fragment fade-out" data-fragment-index="1">?</code><code class="fragment" data-fragment-index="1">10</code></div>
<div>3 s Timeout</div><div><code>20 / s</code></div><div><code>3 s</code></div><div class="swap"><code class="q fragment fade-out" data-fragment-index="1">?</code><code class="fragment" data-fragment-index="1">60</code></div>
</div>

<div class="swap swap-foot">
<p class="ask fragment fade-out" data-fragment-index="1">Legt euch vorher fest: Braucht ein langsames Backend <strong>mehr</strong>, <strong>weniger</strong> oder <strong>gleich viele</strong> Slots?</p>
<div class="callout fragment" data-fragment-index="1">Zehnmal langsamer hei&szlig;t zehnmal mehr Slots, bei gleichem Verkehr.</div>
</div>

</div>

</div>

Note:
- Erst festlegen lassen (mehr, weniger, gleich viele), dann rechnen. Die meisten tippen auf &bdquo;weniger&ldquo;, weil ein langsames Backend geschont werden soll. Little&rsquo;s Law sagt das Gegenteil.
- Ein Klick zeigt die Ergebnisse: 1, 10, 60. Bei 3 s Timeout und 20 Aufrufen pro Sekunde h&auml;ngen 60 Aufrufe gleichzeitig. Ein Pool von zehn lehnt dann f&uuml;nf von sechs ab.
- Die zehn aus dem Code ist eine Workshop-Zahl, die f&uuml;rs Demo passt. In Produktion aus drei Faktoren ableiten: Connection-Pool des HTTP-Clients (mehr Slots als Connections sind sinnlos), Backend-Kapazit&auml;t (Replicas mal maxConcurrent kleiner gleich Kapazit&auml;t) und Little&rsquo;s Law aus gemessener Latenz und Ziel-Durchsatz.
- Kontraintuitiv und der Merksatz der &Uuml;bung: Langsamere Backends brauchen mehr Slots, nicht weniger. Wer den Pool klein macht, &bdquo;um das Backend zu schonen&ldquo;, lehnt ab, was das Backend h&auml;tte schaffen k&ouml;nnen.
- Spicy: Wer <code>maxConcurrent</code> aus dem Bauch setzt, hat den Bulkhead nicht implementiert, sondern dekoriert.
