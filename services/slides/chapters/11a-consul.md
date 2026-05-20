<!-- .slide: data-background-image="./assets/service_discovery.png" data-background-size="contain" data-background-position="center" data-background-opacity="0.18" data-background-repeat="no-repeat" -->

## Consul

<p class="subtitle">HTTP-Registry &mdash; so sieht's konkret aus</p>

<div class="cols">
<div>

<pre class="cheatsheet"><span class="cmd">REGISTER (beim Start):</span>
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

<span class="cmd">DEREGISTER (beim Shutdown):</span>
  PUT consul:8500/v1/agent/service/deregister/{serviceID}
</pre>

</div>
<div>

<pre class="cheatsheet"><span class="cmd">resolve(name):</span>
  res = GET consul:8500/v1/health/service/{name}?passing=true
  // res = [ {Service:{Address, Port}}, ... ]
  instance = randomPick(res)
  return "http://{instance.Address}:{instance.Port}"

<span class="cmd">GET /booking/offers:</span>
  for svc in [flight-service, hotel-service, car-service]:
    url = resolve(svc)
    results[svc] = GET url/...

<span class="cmd">POST /booking/bookings { flightId, hotelId, carId, customerName }:</span>
  for svc in [flight-service, hotel-service, car-service]:
    url = resolve(svc)
    POST url/bookings
  &rarr; { bookingId, customerName, flight, hotel, car }
</pre>

</div>
</div>

Note:
- Linke Spalte = Lifecycle der Instanz (anmelden / abmelden). Rechte Spalte = das, was der Caller bei jedem Request tut (resolvieren + downstream rufen).
- Identischer Pseudo-Code findet sich im Dashboard unter Story 2 &rarr; &bdquo;Spickzettel&ldquo;. Wiedererkennungseffekt gewollt.
- `?passing=true` ist der Knackpunkt: ohne den Filter bekommt ihr auch unhealthy Instanzen zur&uuml;ck. Health-Check ist nur dann wertvoll, wenn der Client ihn auch respektiert.
- `DeregisterCriticalServiceAfter: 1m` &mdash; Consul r&auml;umt tote Eintr&auml;ge selbst weg, falls der Shutdown-Hook nicht durchkommt (Kill -9, OOM, etc.). Wichtig f&uuml;r Selbstheilung.
- Im Reference-Image (`booking/story2/`) ist genau das implementiert: siehe `services/shared/consul/register.go` und `resolver.go`.
- Trade-off besprechen: vor jedem Call resolvieren = stets aktuell, aber Last auf Consul. Alternative: Cache mit Watch / TTL.
