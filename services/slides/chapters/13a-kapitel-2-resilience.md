<!-- .slide: data-background-color="#0A0349" -->

<div class="page dark chapter">

<p class="kicker">Kapitel 2 von 3</p>

<div class="page-body">
<h1 class="chapter-title">Resilience</h1>
<p class="chapter-sub">Wenn Backends kippen.</p>
</div>

</div>

Note:
- <strong>Br&uuml;cke:</strong> Kapitel 1 hat euren Service zum Laufen gebracht und auffindbar gemacht. Discovery sagt, wo ein Service ist. Ob er antwortet, sagt sie nicht.
- <strong>Ab jetzt:</strong> Was tut ihr, wenn die <em>anderen</em> nicht laufen? Kaputt &rarr; Circuit Breaker (Story 4), langsam &rarr; Bulkhead (Story 5).
- <strong>Format:</strong> Zwei Spiele, kein Code. Story 4 noch heute, Story 5 morgen fr&uuml;h.
- <strong>Optional, Kriegsgeschichte:</strong> CrowdStrike, Juli 2024: ein fehlerhaftes Update, weltweit stehen Check-in-Systeme, Bordkarten teils von Hand. Eine Abh&auml;ngigkeit kippt, alles kippt mit.
- <strong>Weiter:</strong> Zuerst der Fall &bdquo;Backend ist tot&ldquo;.
