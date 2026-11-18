<div class="page">

<p class="kicker">Story 5 &middot; Rollenspiel</p>

<div class="page-head">

## Der geteilte Pool

<span class="badge design">&asymp; 25 min <small>Spiel</small></span>
</div>

<p class="subtitle">Booking schafft neun Requests gleichzeitig. Und Hotel wird <span class="hl">langsam</span>.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> l&auml;sst herein, h&ouml;chstens neun gleichzeitig, und gibt den Default-Chip.</li>
<li><strong>Flight, Hotel, Car</strong> halten ihr Schild und geben ihre Chips.</li>
<li><strong>Alle anderen</strong> sind Requests in der Schlange.</li>
</ol>
</div>
<div class="steps-col">
<h4>So l&auml;uft jeder Request</h4>
<p class="step-rule big">Chips abholen: Flight, Hotel, Car. Mit allen drei bist du fertig, hinten anstellen.</p>
<p class="step-note"><strong>Kein Hotel-Chip:</strong> bei Hotel stehen bleiben.</p>
<p class="step-note"><strong>Alle Pl&auml;tze voll:</strong> an der T&uuml;r warten.</p>
<p class="step-note"><strong>Mit Bulkhead:</strong> Hotel h&ouml;chstens drei, der Vierte bekommt sofort Default.</p>
</div>
</div>

</div>

</div>

Note:
- <strong>Hook:</strong> &bdquo;Story 4 hatte ein kaputtes Backend. Heute antwortet Hotel, nur langsam.&ldquo;
- <strong>Gruppe:</strong>
  - 14 Personen: Booking, Flight, Hotel, Car, zehn Requests. Neun parallel, mit Bulkhead drei pro Service. Der Trainer spielt nicht mit, er sagt die Stufen an und stoppt.
  - Es braucht einen Request mehr als Pl&auml;tze. 11&ndash;13: sechs parallel, zwei pro Service. 8&ndash;10: drei parallel, einer pro Service. Bei 13 geht auch neun, wenn der Trainer den zehnten Request spielt.
  - Unter 8 nicht spielbar. Ab 18 zwei Gruppen, die Chips teilen.
- <strong>Regeln ansagen:</strong>
  - Booking l&auml;sst den N&auml;chsten herein, sobald der Vorige fertig ist oder wartet. Alle Pl&auml;tze belegt: an der T&uuml;r warten, kein einziger Chip.
  - Immer Flight, dann Hotel, dann Car. Gibt Hotel nichts, bleibt man dort stehen.
  - Mit drei Chips fertig: Chips behalten, hinten anstellen.
- <strong>Aufbau:</strong> Booking vor der Schlange, Flight, Hotel, Car nebeneinander in dieser Reihenfolge. Auf jedem Schild ein Chip der Farbe, die Default-Chips h&auml;lt Booking. Trainer-Skript: Materialseite.
- <strong>Stufe 1, ohne Bulkhead, alles schnell:</strong> zwei Requests &uuml;ben den Ablauf.
- <strong>Stufe 2, Hotel langsam:</strong> Hotel dreht sein Schild. &bdquo;Hotel antwortet noch, etwa 2 Sekunden. In 2 Sekunden kommen Hunderte Requests. Im Spiel gibt Hotel bis zum STOPP nichts mehr.&ldquo; R1&ndash;R9 holen Flight und bleiben bei Hotel stehen, R10 steht an der T&uuml;r.
- <strong>STOPP 1:</strong> Standbild.
  - &bdquo;Ihr bei Hotel, was habt ihr in der Hand?&ldquo; Nur Flight.
  - &bdquo;R10, was hast du?&ldquo; Nichts, nicht einmal Flight, und Flight ist gesund.
  - &bdquo;Car, wie viele Chips hast du ausgegeben?&ldquo; Keinen, obwohl Car gesund ist.
  - &bdquo;Was ist kaputt?&ldquo; Nichts. Hotel ist nur langsam.
  - Chips zur&uuml;ck, neu anstellen. &bdquo;Dieselben neun Pl&auml;tze, aber jeder Service bekommt drei davon.&ldquo; Zusatzregel erkl&auml;ren, Hotel auf schnell.
- <strong>Stufe 3, mit Bulkhead, alles schnell:</strong> zwei Requests. &bdquo;Was hat sich ge&auml;ndert?&ldquo; Nichts, der Bulkhead kostet nichts, solange alles l&auml;uft.
- <strong>Stufe 4, Hotel wieder langsam:</strong> Vorhersage per Handzeichen. R1&ndash;R3 bleiben bei Hotel, ab R4 Flight, Default von Booking, Car, fertig, wieder anstellen.
- <strong>STOPP 2:</strong>
  - &bdquo;Wie viele stehen fest, wie viele sind frei?&ldquo; Stufe 2 alle fest, Stufe 4 drei fest, sieben frei.
  - &bdquo;Was hat die Grenze gekostet?&ldquo; Weniger Hotels. Verf&uuml;gbarkeit vor Vollst&auml;ndigkeit.
  - &bdquo;Hilft es, Hotel alle neun Pl&auml;tze zu geben?&ldquo; Nein, das ist Stufe 2.
  - &bdquo;Wem geh&ouml;rt die Grenze?&ldquo; Booking, Hotel wei&szlig; nichts davon. F&uuml;nf Booking-Instanzen: Hotel sieht f&uuml;nfzehn.
  - Zum Schluss: &bdquo;2 Sekunden sind um.&ldquo; Hotel gibt den drei Wartenden ihre Chips. Hotel war nie kaputt.
- Demo: Alle Services auf Normal, dann Hotel auf &bdquo;Langsam&ldquo;, Burst. &bdquo;hotel rej&ldquo; zeigt etwa zehn, die Ergebniszeile &bdquo;0 fail&ldquo;. Danach Hotel zur&uuml;ck. Details: <code>docs/stories/story-05-bulkhead.md</code>.
