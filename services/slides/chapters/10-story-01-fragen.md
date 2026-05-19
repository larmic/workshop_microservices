## Story 1 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Health-Check</h3>
<p>Unser <code>/health</code> gibt 200 zur&uuml;ck. Hei&szlig;t das, der Service ist gesund?</p>
<code>/health &rarr; 200</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Ein Build, viele Umgebungen</h3>
<p>Wie viele ENV-Variablen sind in der Realit&auml;t noch Twelve-Factor &mdash; und ab wann ist es nur die neue YAML-Wand?</p>
<code>ENV-Wildwuchs</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Logs auf stdout</h3>
<p>Sauber f&uuml;r einen Prozess &mdash; was passiert, wenn <span class="hl">50 Instanzen</span> gleichzeitig log&iuml;en?</p>
<code>stdout &rarr; Aggregator</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Der /info-Endpoint</h3>
<p>Wir zeigen die aktuelle Config. Was, wenn das aus Versehen im Internet erreichbar ist?</p>
<code>/info</code>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Kein Error-Handling</h3>
<p>In Story 1 darf alles fehlschlagen. Wie <span class="hl">ehrlich</span> ist das wirklich?</p>
<code>0.99&sup3; &asymp; 97 %</code>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

Note:
- Frage 1 &mdash; **Meine Antwort:** Praktisch nichts. Unser `/health` sagt nur &bdquo;der HTTP-Server l&auml;uft und kann eine Route bedienen&ldquo; &mdash; nicht, ob Backends erreichbar sind, ob der Service arbeitsf&auml;hig ist (DB), ob nur noch 5 % freie Threads &uuml;brig sind. In der Cloud-Welt unterscheidet man **Liveness** (Prozess lebt? &rarr; Container neu starten) und **Readiness** (kann Last sehen? &rarr; aus LB nehmen). Spicy: ein zu cleverer Readiness-Check kippt alle Instanzen gleichzeitig aus dem LB, sobald ein Backend wackelt &mdash; und die Plattform ist offline, obwohl die Services noch funktionieren w&uuml;rden.
- Frage 2 &mdash; **Meine Antwort:** Das Pattern ist nicht das Problem, der **Konfigurationsumfang** ist es. ENV-Wildwuchs entsteht graduell (3&ndash;5 neue Variablen pro Integration). Twelve-Factor sagt nichts &uuml;ber Secret-Handling. &bdquo;Ein Build, viele Umgebungen&ldquo; funktioniert nur, solange die Umgebungen wirklich gleich sind. Take-away: nicht &bdquo;m&ouml;glichst viel &uuml;ber ENV&ldquo;, sondern **&bdquo;keine umgebungsspezifischen Werte im Image&ldquo;**. Feature-Flags per ENV = eigentlich mehrere Services in einem.
- Frage 3 &mdash; **Meine Antwort:** Die Verantwortung wandert nach oben &mdash; Container-Runtime, Log-Forwarder (Fluentd/Vector), Aggregator (Loki/ELK/Datadog). Stolperfallen: Container-Logs sind rotated (zu sp&auml;t = weg), Plain Text vs. JSON, ohne **Trace-IDs** ist der 50-Pod-Stream nur Rauschen (Br&uuml;cke zu Story 7), und Cloud-Logging kostet pro GB/Monat &mdash; ein lautes `log.Println` in der hei&szlig;en Schleife produziert vierstellige Rechnungen. Spicy: &bdquo;Logs auf stdout&ldquo; ist die richtige Default-Wahl &mdash; verschiebt die Komplexit&auml;t nur. Wer den Aggregations-Layer nicht durchdacht hat, fliegt nach drei Monaten auf die Nase.
- Frage 4 &mdash; **Meine Antwort:** Klassisches **Information Disclosure**. `consulUrl: http://consul:8500` &rarr; &bdquo;Aha, ihr habt Consul, vermutlich ohne ACL, freue mich auf den Pivot.&ldquo; `timeout: 3000` &rarr; trivial DoS-bar. Spring-Actuator-Bingo: `/actuator/env` versehentlich offen, inkl. DB-Passw&ouml;rter. Take-away: Mgmt-Endpoints (`/health`, `/info`, `/admin/*`) brauchen eine Antwort auf &bdquo;wer darf das aufrufen?&ldquo; &mdash; separater Mgmt-Port, Auth oder per ENV deaktivierbar. Bei uns im Workshop ist alles offen &mdash; bewusst f&uuml;rs Debugging. In Produktion w&auml;re das fahrl&auml;ssig.
- Frage 5 &mdash; **Meine Antwort:** In der ersten Iteration ist &bdquo;happy path&ldquo; das Richtige &mdash; wer alles antizipiert, liefert nichts. Aber: das ist **nicht** der Endzustand. Drei Backends sequenziell mit je 99 % Uptime: 0.99 &times; 0.99 &times; 0.99 &asymp; 97 % &mdash; &asymp; 21 Tage Downtime pro Jahr. &bdquo;Kein Error-Handling&ldquo; ist eine **didaktische Reduktion**, keine Architektur-Empfehlung. Stories 3 / 4 / 5 schlie&szlig;en die L&uuml;cke (Circuit Breaker, Bulkhead, Saga).
- Diskussionsfutter f&uuml;r die Runde: Wie viele Konfigurationswerte hat euer schlimmster Service in der Praxis? Wer schreibt bei euch die Health-Checks &mdash; Entwickler oder Ops? Wenn ihr `/info` auf einem produktiven System aufruft: was seht ihr da, was ihr nicht sehen solltet?
- Vollst&auml;ndige Antworten und weitere Anekdoten: `docs/questions/story1.md`.
