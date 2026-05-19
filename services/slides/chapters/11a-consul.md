## Consul

<p class="subtitle">HTTP-Registry &mdash; so sieht's konkret aus</p>

<div class="cols">
<div>

### Konzepte

- **Agent** l&auml;uft auf jedem Host / als Sidecar
- **Service-Registry** + **Health-Check-Engine** in einem
- HTTP-API auf Port **8500** (DNS-Interface auf 8600)
- **Selbst-Heilung**: `DeregisterCriticalServiceAfter` r&auml;umt tote Instanzen
- **Datacenter-f&auml;hig**: ein Agent pro DC, Cross-DC-Replikation
- Im Workshop: <span class="hl">register / deregister / resolve</span>

</div>
<div>

<pre>REGISTER (beim Start):
  PUT consul:8500/v1/agent/service/register
  {
    Name:    "flight-service",
    ID:      "flight-service-{hostname}",
    Address: "{hostname}",
    Port:    8080,
    Check: {
      HTTP:     "http://{host}:8080/health",
      Interval: "10s",
      DeregisterCriticalServiceAfter: "1m"
    }
  }

DEREGISTER (beim Shutdown):
  PUT consul:8500/v1/agent/service/deregister/{serviceID}

RESOLVE (vor jedem Call):
  res = GET consul:8500/v1/health/service/{name}?passing=true
  // res = [ {Service:{Address, Port}}, ... ]
  instance = randomPick(res)
  return "http://{instance.Address}:{instance.Port}"</pre>

</div>
</div>

Note:
- Identischer Pseudo-Code findet sich im Dashboard unter Story 2 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- `?passing=true` ist der Knackpunkt: ohne den Filter bekommt ihr auch unhealthy Instanzen zur&uuml;ck. Health-Check ist nur dann wertvoll, wenn der Client ihn auch respektiert.
- `DeregisterCriticalServiceAfter: 1m` &mdash; Consul r&auml;umt tote Eintr&auml;ge selbst weg, falls der Shutdown-Hook nicht durchkommt (Kill -9, OOM, etc.). Wichtig f&uuml;r Selbstheilung.
- Im Reference-Image (`booking/story2/`) ist genau das implementiert: siehe `services/shared/consul/register.go` und `resolver.go`.
- Trade-off besprechen: vor jedem Call resolvieren = stets aktuell, aber Last auf Consul. Alternative: Cache mit Watch / TTL.
