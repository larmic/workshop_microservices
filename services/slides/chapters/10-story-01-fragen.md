## Story 1 &mdash; Recap

<p class="subtitle">Fragen nach der Umsetzung</p>

<div class="recap-grid">

<div class="factor fragment">
<h3><span class="numeral">1</span> Health-Check</h3>
<p>Unser <code>/health</code> gibt 200 zur&uuml;ck. Hei&szlig;t das, der Service ist <span class="hl">gesund</span>?</p>
<code>healthy &ne; useful</code>
<aside class="notes"><strong>Meine Antwort:</strong> Unser <code>/health</code> ist eher <strong>Liveness</strong> &mdash; &bdquo;Prozess lebt&ldquo;. Eine richtige <strong>Readiness</strong> m&uuml;sste pr&uuml;fen, ob alle Abh&auml;ngigkeiten erreichbar sind: Consul (Story 2), Backends (Story 3+), DB-Pool, Message-Broker. Kubernetes kennt zus&auml;tzlich <strong>Startup-Probes</strong> f&uuml;r langsame Booting-Phasen &mdash; bevor Readiness &uuml;berhaupt geprobed wird.<br><strong>Spicy:</strong> Wer nur Liveness baut, freut sich morgens &uuml;ber gr&uuml;nes Dashboard auf einem Service, der seit Stunden alle DB-Calls in den Timeout laufen l&auml;sst. <em>Healthy ist nicht dasselbe wie Useful.</em></aside>
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
<code>0.99&sup3; &asymp; 97 % &rarr; 21 Tage Downtime/Jahr</code>
<aside class="notes"><strong>Meine Antwort:</strong> In der ersten Iteration ist &bdquo;happy path&ldquo; das Richtige &mdash; wer alles antizipiert, liefert nichts. Aber: das ist <strong>nicht</strong> der Endzustand. Drei Backends sequenziell mit je 99 % Uptime: 0.99 &times; 0.99 &times; 0.99 &asymp; 97 % &mdash; &asymp; 21 Tage Downtime pro Jahr.<br><strong>Spicy:</strong> &bdquo;Kein Error-Handling&ldquo; ist eine <strong>didaktische Reduktion</strong>, keine Architektur-Empfehlung. Stories 3 / 4 / 5 schlie&szlig;en die L&uuml;cke (Circuit Breaker, Bulkhead, Saga).</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">4</span> Wenn ein Service umzieht?</h3>
<p>Backend-URLs stehen in ENV-Variablen. Was passiert, wenn Flight, Hotel oder Car woanders <span class="hl">hin wandern</span>?</p>
<code>FLIGHT_URL = ?</code>
<aside class="notes"><strong>Meine Antwort:</strong> Der Service ist dann nicht mehr erreichbar &mdash; neue URL, neuer Port, alte ENV-Variable zeigt ins Leere. Jede Adress&auml;nderung erzwingt ein Re-Deployment des Booking-Service. Skaliert nicht: bei 5 Backends &times; 4 Aufrufern sind das 20 Stellen zu pflegen. Bei jeder Topologie&auml;nderung.<br>Das bringt uns zu <strong>Story 2 &mdash; Service Discovery</strong>: Services registrieren sich selbst (z.&nbsp;B. bei Consul), der Booking-Service findet sie &uuml;ber den logischen Namen, nicht &uuml;ber eine hartkodierte URL.<br><strong>Spicy:</strong> Wer einmal Backend-URLs in YAML pflegt, wei&szlig; nach drei Jahren nicht mehr, ob die noch stimmen &mdash; und probiert es im Zweifel auf Prod aus.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">5</span> Config &mdash; ENV oder Datei?</h3>
<p>ENV-Vars, Property-Files, ConfigMaps &mdash; alle haben Liebhaber. Wo ist der <span class="hl">Unterschied</span>?</p>
<code>ENV &gt; File &gt; Code</code>
<aside class="notes"><strong>Meine Antwort:</strong> 12-Faktor sagt ENV &mdash; weil Container-Welt: immutable beim Start, kein Reload bei &Auml;nderung. <strong>Property-Files</strong> sind editierbar zur Laufzeit (Gefahr: kein Audit-Trail), brauchen Volume-Mounts. <strong>ConfigMaps</strong> (K8s) sind eine Mischung: per ENV oder Volume in den Pod, zentral verwaltet. <strong>Secret-Manager</strong> (Vault, AWS Secrets, GCP Secret Manager) f&uuml;r sensible Daten &mdash; Rotation, Audit, RBAC.<br><strong>Spicy:</strong> &bdquo;Wir lesen halt YAML rein&ldquo; funktioniert, bis Prod-YAML und Dev-YAML divergieren und niemand mehr wei&szlig;, warum. Klare Trennung: <em>nicht-sensible Config in ConfigMap, sensible in Secret-Manager, nichts im Git</em>.</aside>
</div>

<div class="factor fragment">
<h3><span class="numeral">6</span> OpenAPI &mdash; Pflicht oder Nice-to-have?</h3>
<p><code>GET /openapi</code> liefert die Spec. Ist das <span class="hl">Spielerei</span> oder essenziell?</p>
<code>API-First</code>
<aside class="notes"><strong>Meine Antwort:</strong> Pflicht, sobald mehr als ein Team / ein Mensch den Service aufruft. Vorteile: Tests gegen die Spec (Schema-Validation), Codegen f&uuml;r Clients (Java / TS / Go), Breaking Changes sind diff-bar, Versionierung sichtbar. Aber: die Spec ist nur dann wertvoll, wenn sie <strong>nicht</strong> nachtr&auml;glich aus dem Code generiert wird &mdash; das nennt sich <strong>API-First</strong>. Spec geh&ouml;rt <em>vor</em> den Code.<br><strong>Spicy:</strong> &bdquo;OpenAPI generieren wir uns mit Annotations&ldquo; &mdash; dann hat das Team die Spec nie gelesen. Was nicht im Vertrag steht, ist nicht das Verhalten. Wer die Spec sp&auml;ter generiert, dokumentiert, was er gebaut hat &mdash; nicht, was er bauen sollte.</aside>
</div>

<span class="show-all fragment" aria-hidden="true"></span>

</div>

<aside class="notes">
Vollst&auml;ndige Antworten und weitere Anekdoten: <code>docs/questions/story1.md</code>.
</aside>
