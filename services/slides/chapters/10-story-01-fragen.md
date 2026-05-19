## Story 1 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Health-Check</h3>
<p>Unser <code>/health</code> gibt 200 zur&uuml;ck. Hei&szlig;t das, der Service ist gesund?</p>
<code>/health &rarr; 200</code>
<aside class="notes"><strong>Meine Antwort:</strong> Unser <code>/health</code> ist eher <strong>Liveness</strong>. Eine <strong>Readiness</strong> m&uuml;sste z.&nbsp;B. pr&uuml;fen, ob Consul erreichbar ist (Story 2), ob die Backends antworten (Story 3+), usw.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">2</span> Logs auf stdout</h3>
<p>Sauber f&uuml;r einen Prozess &mdash; was passiert, wenn <span class="hl">50 Instanzen</span> gleichzeitig loggen?</p>
<code>stdout &rarr; Aggregator</code>
<aside class="notes"><strong>Meine Antwort:</strong> Die Verantwortung wandert nach oben &mdash; Container-Runtime, Log-Forwarder (Fluentd/Vector), Aggregator (Loki/ELK/Datadog). Stolperfallen: Container-Logs sind rotated (zu sp&auml;t = weg), Plain Text vs. JSON, ohne <strong>Trace-IDs</strong> ist der 50-Pod-Stream nur Rauschen (Br&uuml;cke zu Story 7), und Cloud-Logging kostet pro GB/Monat &mdash; ein lautes <code>log.Println</code> in der hei&szlig;en Schleife produziert vierstellige Rechnungen.<br><strong>Spicy:</strong> &bdquo;Logs auf stdout&ldquo; ist die richtige Default-Wahl &mdash; verschiebt die Komplexit&auml;t nur. Wer den Aggregations-Layer nicht durchdacht hat, fliegt nach drei Monaten auf die Nase.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">3</span> Kein Error-Handling</h3>
<p>In Story 1 darf alles fehlschlagen. Wie <span class="hl">ehrlich</span> ist das wirklich?</p>
<code>0.99 &times; 0.99 &times; 0.99 &asymp; 97 % &mdash; bis zu 21 Tage Downtime pro Jahr.</code>
<aside class="notes"><strong>Meine Antwort:</strong> In der ersten Iteration ist &bdquo;happy path&ldquo; das Richtige &mdash; wer alles antizipiert, liefert nichts. Aber: das ist <strong>nicht</strong> der Endzustand. Drei Backends sequenziell mit je 99 % Uptime: 0.99 &times; 0.99 &times; 0.99 &asymp; 97 % &mdash; &asymp; 21 Tage Downtime pro Jahr.<br><strong>Spicy:</strong> &bdquo;Kein Error-Handling&ldquo; ist eine <strong>didaktische Reduktion</strong>, keine Architektur-Empfehlung. Stories 3 / 4 / 5 schlie&szlig;en die L&uuml;cke (Circuit Breaker, Bulkhead, Saga).</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Wenn ein Service umzieht?</h3>
<p>Backend-URLs stehen in ENV-Variablen. Was passiert, wenn Flight, Hotel oder Car woanders hin wandern?</p>
<code>FLIGHT_URL = ?</code>
<aside class="notes"><strong>Meine Antwort:</strong> Der Service ist dann ggf. nicht mehr erreichbar &mdash; neue URL, neuer Port, alte ENV-Variable zeigt ins Leere. Jede Adress&auml;nderung erzwingt ein Re-Deployment des Booking-Service.<br>Das bringt uns zu <strong>Story 2 &mdash; Service Discovery</strong>: Services registrieren sich selbst (z.&nbsp;B. bei Consul), der Booking-Service findet sie &uuml;ber den logischen Namen, nicht &uuml;ber eine hartkodierte URL.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story1.md</code>.
</aside>
