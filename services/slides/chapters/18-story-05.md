<div class="page">

<p class="kicker">Story 5 &middot; Rollenspiel</p>

<div class="page-head">

## Der geteilte Pool

<span class="badge design">&asymp; 25 min <small>Spiel</small></span>
</div>

<p class="subtitle">Zwei St&uuml;hle sind die Threads von Booking. Und Hotel wird <span class="hl">langsam</span>.</p>

<div class="page-body">

<div class="steps">
<div class="steps-col">
<h4>Rollen</h4>
<ol>
<li><strong>Booking</strong> ruft die Nummern und gibt die Karten aus.</li>
<li><strong>Hotel, Flight, Car</strong> halten ihr Schild und z&auml;hlen mit.</li>
<li><strong>Alle anderen</strong> sind Requests in der Schlange.</li>
</ol>
</div>
<div class="steps-col">
<h4>So l&auml;uft jede Nummer</h4>
<p class="step-rule big">Wer sie auf der Karte hat, steht auf. Wer wartet, r&uuml;ckt nach. Dann die n&auml;chste Karte.</p>
<p class="step-note"><strong>Dunkle Karte:</strong> STOPP, wir reden.</p>
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
- <strong>Regeln ansagen</strong>, bei jeder Nummer in dieser Reihenfolge:
  - Booking ruft die Nummer.
  - Wer sitzt und die Nummer auf der Karte hat: aufstehen, Car abklatschen, Karte auf den Tisch, hinten anstellen.
  - Wer hinter den St&uuml;hlen wartet, setzt sich und klatscht Flight ab. Wer zuerst kam, zuerst.
  - Booking gibt die Karte dem N&auml;chsten. Stuhl frei: Flight abklatschen, hinsetzen. Sonst hinter die St&uuml;hle.
  - Hotel z&auml;hlt, wer sich setzt. Flight und Car z&auml;hlen, wer abklatscht.
- <strong>Aufbau:</strong> Zwei St&uuml;hle, Sitzfl&auml;che zu Booking. Ein Tisch, alle Karten auf einen Haufen. Die Wand liegt bereit (Pinnwand, Flipchart, Tisch hochkant). Karten, Schilder, Trainer-Skript: Materialseite.
- <strong>STOPP 1 (Karte 04):</strong> &bdquo;Wie viele St&uuml;hle waren gleichzeitig besetzt?&ldquo; Einer: eine Ankunft pro Nummer mal eine Nummer Dauer.
  - Hotel dreht auf &bdquo;langsam&ldquo;. Vorhersage per Handzeichen. Z&auml;hler auf 0.
- <strong>STOPP 2 (Karte 15):</strong> Standbild, zwei sitzen (09, 10), vier stehen dahinter (11&ndash;14).
  - &bdquo;Ihr hinten, worauf wartet ihr?&ldquo; Auf einen Stuhl. Dabei wollt ihr zuerst nur zu Flight, und Flight ist gesund.
  - &bdquo;Flight, Car, wie viele kamen an?&ldquo; 6 und 4.
  - &bdquo;Hat der Breaker etwas gemerkt?&ldquo; Nein: 2 s liegen unter 3 s Timeout, CLOSED.
  - &bdquo;Karte 14, wann stehst du auf?&ldquo; Bei 26, also 6 s nach dem Ankommen. Das Doppelte des Timeouts, der Client w&auml;re l&auml;ngst weg.
- <strong>Die Wand:</strong> Z&auml;hler merken, dann auf 0. Alle Karten auf den Tisch, alle neu anstellen.
  - Wand zwischen die St&uuml;hle schieben: &bdquo;Dieselben zwei St&uuml;hle, nur mit einer Wand.&ldquo; Hotel auf die eine Seite, Flight und Car auf die andere.
  - Zusatzregel erkl&auml;ren, dann Vorhersage: Wie viele der elf bekommen Hotels?
- <strong>Runde 2 (16&ndash;26):</strong> Auf dem Hotel-Stuhl sitzen nacheinander 16, 20, 24. Alle anderen klatschen Flight und Car ab und sind sofort fertig. Niemand steht an.
- <strong>STOPP 3 (Karte 27):</strong> Z&auml;hler abfragen, erst Runde 1 (gemerkt), dann Runde 2.
  - Beantwortet (Car) 4 &rarr; 10, davon mit Hotels 4 &rarr; 2, hinter den St&uuml;hlen 4 &rarr; 0, Flight 6 &rarr; 11.
  - &bdquo;Welche Runde war besser?&ldquo; Runde 1 mehr Hotels, Runde 2 mehr beantwortet und niemand stand an.
  - &bdquo;Wie viele St&uuml;hle h&auml;tte Hotel gebraucht?&ldquo; Vier, Hotel hat einen. Deshalb bekommt jeder Vierte Hotels (16, 20, 24). Little&rsquo;s Law zum Anfassen: ein Stuhl = 10 Slots, langsam br&auml;uchte es 40.
  - &bdquo;Breaker?&ldquo; CLOSED, ein besetzter Stuhl ist kein Fehler.
  - &bdquo;Wem geh&ouml;ren die St&uuml;hle?&ldquo; Booking, nicht Hotel. F&uuml;nf Booking-Instanzen h&auml;tten f&uuml;nf Hotel-St&uuml;hle.
- Demo: Alle Services auf Normal, dann Hotel auf &bdquo;Langsam&ldquo;, Burst. &bdquo;hotel rej&ldquo; zeigt etwa zehn, die Ergebniszeile &bdquo;0 fail&ldquo;. Danach Hotel zur&uuml;ck. Details: <code>docs/stories/story-05-bulkhead.md</code>.
