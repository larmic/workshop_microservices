# E-Mail-Einladung: Vorlage pro Termin

Vorlage für die Einladungs-Mail an angemeldete Teilnehmer:innen. Platzhalter
ersetzen und versenden. Die Schritt-für-Schritt-Details leben in
[vorbereitung.md](../vorbereitung.md), die Mail benennt sie nur kompakt.

Platzhalter:

| Platzhalter    | Bedeutung                                                       | Beispiel                    |
|----------------|-----------------------------------------------------------------|-----------------------------|
| `{{DATUM}}`    | Workshop-Termin (beide Tage)                                    | 24. + 25.06.2026            |
| `{{UHRZEIT}}`  | Beginn und Ende                                                 | jeweils 9:00 bis 17:00 Uhr  |
| `{{ORT}}`      | Veranstaltungsort                                               | Schuppen Eins               |
| `{{DEADLINE}}` | Stichtag für die Vorbereitung, Vorschlag: 3 Werktage vor Termin | 19.06.2026                  |
| `{{KONTAKT}}`  | Erreichbarkeit des Trainers                                     | Mail-Adresse / Teams        |

---

**Betreff:** Microservices-Workshop am {{DATUM}} · bitte Vorbereitung bis {{DEADLINE}}

Hallo zusammen,

schön, dass ihr dabei seid! Hier die wichtigsten Infos zum Microservices-Workshop:

**Wann:** {{DATUM}}, {{UHRZEIT}}
**Wo:** {{ORT}}
**Mitbringen:** Entwickler-Notebook (installiert: Git, Docker inkl. Docker Compose, IDE, Programmiersprache eurer Wahl)

**Vorbereitung (Pflicht, ca. 1 bis 2 Stunden, bis spätestens {{DEADLINE}}):**

Damit wir am Workshop-Tag direkt mit den Inhalten starten, statt Setups zu debuggen, erledigt bitte vorab diese drei Punkte:

1. **Mini-Service bauen:** In eurer Sprache ein neues Projekt aufsetzen mit `GET /health` (antwortet HTTP 200) und einer Dockerfile. Das Image baut, der Container läuft auf Port 8080. Diesen Service erweitert ihr im Workshop Schritt für Schritt.
2. **Workshop-Umgebung einmal starten:** Repo klonen, Stack hochfahren, Dashboard unter <http://localhost> öffnen.
3. **Im Firmennetz testen:** `git clone` und `docker build` einmal in dem Netz ausführen, in dem ihr auch am Workshop-Tag arbeitet. Proxy- und Zertifikatsprobleme fallen sonst erst vor Ort auf, besonders unter Windows/WSL2.

Die genaue Anleitung mit Akzeptanzkriterien (Abschnitt "0. Vor dem Workshop (Pflicht)"):
<https://larmic.github.io/workshop_microservices/vorbereitung/>

Das Repo zum Klonen: <https://github.com/larmic/workshop_microservices>

Falls etwas hakt: Erste Hilfe gibt die Troubleshooting-Seite (besonders für Windows/WSL2):
<https://larmic.github.io/workshop_microservices/vorbereitung/?doc=troubleshooting.md>
Wenn das nicht reicht, meldet euch einfach vorab bei mir: {{KONTAKT}}. Bitte nicht erst am Workshop-Tag, da bleibt für Setup-Probleme wenig Zeit.

Ich freue mich auf zwei intensive Tage mit euch!

Viele Grüße
Lars
