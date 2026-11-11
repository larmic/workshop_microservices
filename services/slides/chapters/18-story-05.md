<div class="page">

<p class="kicker">Story 5 &middot; Rollenspiel</p>

<div class="page-head">

## Der geteilte Pool

<span class="badge design">&asymp; 60 min <small>davon 25 Spiel</small></span>
</div>

<p class="subtitle">Zwei St&uuml;hle sind die Threads von Booking. Und Hotel wird <span class="hl">langsam</span>.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> ruft die Nummern und gibt die Karten aus.</li>
<li><strong>Hotel, Flight, Car</strong> halten ihr Schild und z&auml;hlen mit: Hotel, wer sich setzt. Flight und Car, wer abklatscht.</li>
<li><strong>Alle anderen</strong> sind Requests in der Schlange vor Booking.</li>
</ol>
<h4>Bei jeder Nummer</h4>
<ol>
<li><strong>Booking</strong> ruft die Nummer.</li>
<li><strong>Wer sitzt</strong> und die Nummer auf der Karte hat: aufstehen, Car abklatschen, Karte auf den Tisch, hinten anstellen.</li>
<li><strong>Wer hinter den St&uuml;hlen wartet</strong>, setzt sich auf den freien Stuhl und klatscht Flight ab.</li>
<li><strong>Booking</strong> gibt die Karte dem N&auml;chsten. Stuhl frei: Flight abklatschen, hinsetzen. Sonst hinter die St&uuml;hle.</li>
</ol>
</div>
<div class="steps-col">
<h4>Ablauf</h4>
<ol>
<li><strong>Probe 01&ndash;03, Hotel schnell:</strong> Man sitzt eine Nummer. Bei STOPP 1 dreht Hotel das Schild auf &bdquo;langsam&ldquo;, ab jetzt vier.</li>
<li><strong>Runde 1, zwei St&uuml;hle f&uuml;r alle:</strong> Was seht ihr bei STOPP 2? <span class="fragment">&rarr; vier stehen dahinter, Flight und Car haben fast nichts zu tun</span></li>
<li><strong>Die Wand kommt</strong> zwischen die St&uuml;hle. Ab Runde 2 gilt: Hotel-Stuhl besetzt? <strong>Nicht warten</strong>, sofort Car abklatschen, fertig ohne Hotels.</li>
<li><strong>Runde 2, ein Stuhl f&uuml;r Hotel:</strong> Wie viele der n&auml;chsten elf bekommen Hotels? <span class="fragment">&rarr; drei von elf bekommen Hotels</span></li>
<li><strong>Circuit Breaker:</strong> Was zeigt er nach beiden Runden? <span class="fragment">&rarr; CLOSED, langsam ist kein Fehler</span></li>
</ol>
</div>
</div>

</div>

</div>

Note:
- <strong>Hook:</strong> &bdquo;Story 4 hatte ein kaputtes Backend. Heute antwortet Hotel, nur langsam.&ldquo;
- <strong>Gruppe:</strong>
  - 10 Personen: Booking, Hotel, Flight, Car, sechs Requests. Der Trainer spielt nicht mit, er leitet die STOPPs und schiebt die Wand.
  - 7: Booking und sechs Requests, die Schilder lehnen, der Trainer z&auml;hlt. 8: Hotel als Person. 9: eine Person h&auml;lt Flight und Car.
  - 11&ndash;13: mehr Requests, sonst gleich. Ab 14 zwei Gruppen mit eigenem Satz. Unter 7 nicht spielbar.
- <strong>Aufbau:</strong> Zwei St&uuml;hle, Sitzfl&auml;che zu Booking. Ein Tisch, alle Karten auf einen Haufen. Die Wand liegt bereit (Pinnwand, Flipchart, Tisch hochkant). Karten, Schilder, Trainer-Skript: Materialseite.
- <strong>STOPP 1 (Karte 04):</strong> &bdquo;Wie viele St&uuml;hle waren gleichzeitig besetzt?&ldquo; Einer: eine Ankunft pro Nummer mal eine Nummer Dauer.
  - Hotel dreht auf &bdquo;langsam&ldquo;. Vorhersage per Handzeichen. Z&auml;hler auf 0.
- <strong>STOPP 2 (Karte 15):</strong> Standbild, zwei sitzen (09, 10), vier stehen dahinter (11&ndash;14).
  - &bdquo;Ihr hinten, worauf wartet ihr?&ldquo; Auf einen Stuhl. Dabei wollt ihr zuerst nur zu Flight, und Flight ist gesund.
  - &bdquo;Flight, Car, wie viele kamen an?&ldquo; 6 und 4.
  - &bdquo;Hat der Breaker etwas gemerkt?&ldquo; Nein: 2 s liegen unter 3 s Timeout, CLOSED.
  - &bdquo;Karte 14, wann stehst du auf?&ldquo; Bei 26, also 6 s nach dem Ankommen. Genau der Timeout.
- <strong>Die Wand:</strong> Z&auml;hler merken, dann auf 0. Alle Karten auf den Tisch, alle neu anstellen.
  - Wand zwischen die St&uuml;hle schieben: &bdquo;Dieselben zwei St&uuml;hle, nur mit einer Wand.&ldquo; Hotel auf die eine Seite, Flight und Car auf die andere.
  - Zusatzregel erkl&auml;ren, dann Vorhersage: Wie viele der elf bekommen Hotels?
- <strong>Runde 2 (16&ndash;26):</strong> Auf dem Hotel-Stuhl sitzen nacheinander 16, 20, 24. Alle anderen klatschen Flight und Car ab und sind sofort fertig. Niemand steht an.
- <strong>STOPP 3 (Karte 27):</strong> Z&auml;hler abfragen, weiter auf der n&auml;chsten Folie.
- Demo: Alle Services auf Normal, dann Hotel auf &bdquo;Langsam&ldquo;, Burst. &bdquo;hotel rej&ldquo; zeigt etwa zehn, &bdquo;failed&ldquo; bleibt 0. Danach Hotel zur&uuml;ck. Details: <code>docs/stories/story-05-bulkhead.md</code>.
