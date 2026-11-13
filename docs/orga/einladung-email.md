# E-Mail-Einladung: Vorlage pro Termin

Vorlage für die Einladungs-Mail an angemeldete Teilnehmende. Platzhalter
ersetzen und versenden. Die Schritt-für-Schritt-Details leben in
[vorbereitung.md](../vorbereitung.md), die Mail benennt sie nur kompakt.

Platzhalter:

| Platzhalter    | Bedeutung                                                       | Beispiel                    |
|----------------|-----------------------------------------------------------------|-----------------------------|
| `{{DATUM}}`    | Workshop-Termin (beide Tage)                                    | 24. + 25.06.2026            |
| `{{UHRZEIT}}`  | Beginn und Ende                                                 | jeweils 9:00 bis 17:00 Uhr  |
| `{{ORT}}`      | Veranstaltungsort                                               | Schuppen Eins               |

---

**Betreff:** Workshop "Microservices richtig gemacht" · bitte Vorbereitung vor den Workshop-Tagen

Hallo zusammen,

schön, dass ihr dabei seid! Hier die wichtigsten Infos zum Workshop "Microservices richtig gemacht":

**Wann:** {{DATUM}}, {{UHRZEIT}}
**Wo:** {{ORT}} (kein Remote)
**Mitbringen:** Entwickler-Notebook (installiert: Git, Docker inkl. Docker Compose, IDE, Programmiersprache eurer Wahl)
**Was euch erwartet:** Service Discovery, Circuit Breaker, Bulkhead, Saga und Tracing. Ihr arbeitet in kleinen Teams an einem eigenen Service, dazu kommen eine Design-Session am Flipchart und drei Spiele.

**Vorbereitung (Pflicht, ca. 1 bis 2 Stunden, bis spätestens vor dem Workshop-Beginn):**

Damit wir am Workshop-Tag direkt mit den Inhalten starten, statt Setups zu debuggen, erledigt bitte vorab diese drei Punkte:

1. einen Mini-Service mit Health-Endpoint und Dockerfile aufsetzen,
2. die Workshop-Umgebung einmal starten,
3. beides einmal im Firmennetz testen.

Die Anleitung Schritt für Schritt: <https://larmic.github.io/workshop_microservices/vorbereitung/>

Falls etwas hakt, meldet euch einfach vorab bei mir. Bitte nicht erst am Workshop-Tag, da bleibt für Setup-Probleme wenig Zeit.

Ich freue mich auf zwei intensive Tage mit euch!

Viele Grüße
Lars
