## Story 2

<p class="subtitle">Design-Session: REST vs. RESTful <span class="time-badge">&asymp; 55 min</span></p>

<div class="cols">
<div>

<div class="story-card">

#### Kontext

Der Booking-Service liefert Angebote. Als N&auml;chstes wird gebucht, und wo gebucht wird, wird storniert und umgebucht. Bevor jemand losprogrammiert, entsteht die API am Flipchart. Kein Code, kein Deployment, kein Port.

#### User Story

Als <em>Produktverantwortliche</em> m&ouml;chte ich <em>eine API f&uuml;r Storno und Umbuchung, die HTTP so nutzt, wie es gemeint ist</em>, damit <em>Retries sicher sind, Fehler sichtbar werden und wir sie nicht in drei Monaten umbauen</em>.

</div>

</div>
<div>

<div class="story-card">

#### Aufgabe

- Teams zu 3 bis 4 Personen, ein Flipchart pro Team
- Szenario: eine Kundin storniert die ganze Reise (Grund soll nachvollziehbar bleiben), ein Kunde bucht nur den Flug um
- Client wiederholt Anfragen bei Timeout

#### Ergebnisformat

Eine Tabelle, eine Zeile pro Endpoint:

- Methode und Pfad
- Request-Kern, Antwort und Status-Code
- idempotent ja/nein
- Begr&uuml;ndung (auch f&uuml;r bewusste Abweichungen)

</div>

</div>
</div>

Note:
- Hook: &bdquo;Beim letzten Mal hat ein Prefetcher &uuml;ber einen GET-Link Buchungen storniert. Das soll uns nicht noch einmal passieren.&ldquo;
- Wiedererkennung: dieselbe Karte im Dashboard unter Story 2 &rarr; &bdquo;Story lesen&ldquo;. Dort gibt es bewusst keinen Spickzettel, keine Buttons und keinen Service.
- Time-Box: 3 Minuten Aufgabe stellen, 20 Minuten Teamarbeit, dann je Team 2 Minuten Vorstellung.
- Spaltenk&ouml;pfe der Tabelle vorab auf die Flipcharts zeichnen, das spart f&uuml;nf Minuten.
- W&auml;hrend der Teamarbeit herumgehen und nur Fragen stellen: &bdquo;Was passiert beim zweiten Aufruf?&ldquo;, &bdquo;Welcher Code, wenn das Hotel ablehnt?&ldquo;
- Vollst&auml;ndige Aufgabenbeschreibung: <code>docs/stories/story-02-api-design-session.md</code>. Trainer-Hinweis mit L&ouml;sungsraum: <code>docs/instructions/rest-vs-restful.md</code>.
