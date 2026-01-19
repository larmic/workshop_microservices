# Microservices Workshop - Ideen

## Slide: Was ist ein Microservice? Definition

- Gibt es eine einheitliche Definition?
- Was ist ein µService? Was ist ein SelfContainedSystem?
- Unterschied zur SOA
  - (Was sind Microservices? Architektur im Überblick)
  - Microservice-Architektur ist eine Weiterentwicklung von serviceorientierter Architektur (SOA).
  - Die zwei Ansätze brechen beide große, komplexe Anwendungen in kleinere Komponenten auf, mit denen es sich leichter arbeiten lässt.
  - Wegen ihrer Ähnlichkeit werden SOA und Microservice-Architektur oft verwechselt.
  - Ein Hauptunterschied zwischen den beiden ist ihr Umfang: Eine SOA ist ein unternehmensweiter Ansatz der Architektur, während es sich bei Microservices um eine Implementierungsstrategie für die Teams in der Anwendungsentwicklung handelt.

## Slide: Wann sind Microservices sinnvoll?

- Frage an die Kursteilnehmer und Antworten sammeln
- Beste Antworten:
  - Technisch notwendig, weil ein Teilaspekt der Gesamtarchitektur mehr Bedarf hat
  - Beispiel: Check24: Einträge einstellen -> eher selten, Einträge suchen -> sehr häufig
- Mehrere Teams arbeiten an dem Produkt
  - Liefert unabhängige Team-Entscheidungen (z.B. unterschiedliche Programmier-Sprachen und -Stile)
- Skalierung auf weitere Menschen ist einfacher, da man nicht die Gesamtarchitektur kennen muss
- (Wann nicht?)

## Slide: Vorteile einer Microservice-Architektur

- Hohe Skalierbarkeit (siehe Check 24)
- Leichter Updates der einzelnen Komponenten (z.B. bei Sicherheitslücken)
- Häufigeres Deployment, weil in der Regel einfacher
- Schnellere Einarbeitung in einzelne Services ist möglich
- APIs dokumentieren einzelne Aspekte der Gesamtarchitektur (im Monolithen wird dies seltener getan)

## Slide: Nachteile einer Microservice-Architektur (im Gegensatz zu einem Monolithen)

- Richtigen Schnitt der Services finden -> DDD + Event Storming können hier helfen
  - Diese Frage stellt sich bei einem Monolithen nicht
- Performance HTTP und MQ ist generell langsamer als ein direkter DB Zugriff
  - Darf ein Service auf die DB eines anderen Services zugreifen? -> beantworten wir später
- Prozesse benötigen ggf. mehrere Services
  - Diese Frage stellt sich bei einem Monolithen nicht
- Single-Sign-On OAuth, SAML, ...
  - Diese Frage stellt sich bei einem Monolithen häufig nicht
- Asynchrone vs synchrone Kommunikation oder HTTP vs MQ
  - Nachteil synchron: der andere Dienst muss online sein
  - Nachteil asynchron: es gibt kein direktes Feedback, ob die Anfrage bearbeitet wurde
  - Diese Frage stellt sich bei einem Monolithen nicht
- Monitoring
  - zentrales Logging-System
  - Tracing der Prozesse (Jaeger, Zipkin, ...)
    - Muss von allen Microservices in der Gesamtarchitektur unterstützt werden -> benötigt globale Architektur-Abstimmung
  - Überwachung der Systeme (Health-Status)
- Debugging
  - Geht nur für den einzelnen Service bzw. wird komplex, wenn man mehrere Services zeitgleich debuggen möchte
- Infrastruktur
  - Container, Kubernetes, ....
  - Komplexer als einen Monolithen hinzustellen
  - dafür für die einzelnen Teams später einfacher, wenn die Infrastruktur geklärt und vorhanden ist

## Slide: Brauchen wir Microservices?

- Starte mit einem Monolithen
- Erstelle Microservices, wenn es wirklich sinnvoll ist (und nicht nur, weil gerade hip)
- Dieser Workshop wurde erstellt, um zu zeigen, wie man Microservices richtig machen kann (wenn man sie denn Braucht)

## Themen

Das sollen Aufgaben werden, die an dem eigenen Rechner umgesetzt werden können. Es fehlt noch eine gute Reihenfolge, damit ein roter Pfaden sichtbar wird.

Eine Idee ist, das ganze am Beispiel einer Reisebuchung durchzuführen (Hotel, Flug und Auto Buchen)

- Twelve-Factor-App
  - Kurzer Überblick
- Health-Checks
  - ein Must-Have
- 1 DB pro Service
- External Configuration
  - Wann ist dies notwendig? Eigentlich immer?
- Circuit Breaker
- API-Gateway
- Service Discovery / Service Registry
- Saga Pattern
- Bulkhead Pattern
- Backends for Frontends (BFF)
- Eventsourcing / Event-Driven Architecture
- CQRS
- Downtimeless Deployment
- API First Ansatz?
- Kulturwandel durch MS:
  - MS ohne DevOps? (Begriffsklärung!) - Antwort: Nein! (Könnte in den Infrastrukturteil)
  - "You build it you run it" (arbeiten wir jetzt 24/7? Wartungsteams?)
  - muss ich "cloud native" sein, um MS bauen und betreiben zu können?
