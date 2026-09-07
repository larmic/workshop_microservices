## REST vs. RESTful

<p class="subtitle">So nicht &middot; so</p>

| Prinzip | So nicht | So |
|---|---|---|
| Ressourcen statt Verben | `GET /getUser?id=1423` | `GET /users/1423` |
| Methode tr&auml;gt Semantik | `POST /users/1423/delete` | `DELETE /users/1423` |
| Idempotenz | `POST /bookings/1423/cancel` (zweimal: doppelt?) | `DELETE /bookings/1423` (zweimal: gleich) |
| GET ohne Seiteneffekte | `GET /bookings/1423/cancel` | `POST /bookings/1423/cancellation` |
| Status-Codes | `200 OK` + `{ "error": "not found" }` | `404 Not Found` + Details im Body |
| Sub-Ressourcen und Filter | `GET /getBookingsForCustomer?id=7` | `GET /customers/7/bookings?status=confirmed` |

Note:
- Die Tabelle ist die Zusammenfassung der sechs Karten. Kurz stehen lassen, dann die Aufgabe stellen.
- Zeile 4 rechts ist bewusst die Sub-Ressource, nicht das DELETE: das ist der Spoiler f&uuml;r Quiz 2, ohne ihn zu verraten.
- &Uuml;berleitung zur Story-Karte: &bdquo;Jetzt ihr. Storno und Umbuchung, am Flipchart, in Teams.&ldquo;
