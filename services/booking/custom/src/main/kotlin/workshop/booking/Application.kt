package workshop.booking

import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO as ClientCIO
import io.ktor.client.plugins.HttpTimeout
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.Application
import io.ktor.server.application.install
import io.ktor.server.cio.CIO
import io.ktor.server.engine.embeddedServer
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.response.respond
import io.ktor.server.response.respondText
import io.ktor.server.routing.get
import io.ktor.server.routing.routing
import org.slf4j.LoggerFactory

private val log = LoggerFactory.getLogger("Application")

private const val TIMEOUT_MILLIS = 5000L

// Shown on every start: this example runs whenever CUSTOM_BOOKING_PATH still
// points to it, and then everything looks green in the dashboard.
private val EXAMPLE_BANNER = """
    |
    |****************************************************************
    |***                                                          ***
    |***   BEISPIEL-LÖSUNG (Kotlin/Ktor) läuft als booking-custom ***
    |***                                                          ***
    |***   Seht ihr das, ist euer eigener Service NICHT           ***
    |***   eingeklinkt. CUSTOM_BOOKING_PATH in services/.env      ***
    |***   auf euren Service setzen und den Stack neu starten.    ***
    |***                                                          ***
    |***   Anleitung: Dashboard http://localhost, Schritt 0 Start ***
    |***                                                          ***
    |****************************************************************
    |""".trimMargin()

fun main() {
    val config = loadConfig()

    val httpClient = HttpClient(ClientCIO) {
        install(HttpTimeout) {
            requestTimeoutMillis = TIMEOUT_MILLIS
            connectTimeoutMillis = TIMEOUT_MILLIS
            socketTimeoutMillis = TIMEOUT_MILLIS
        }
        expectSuccess = false
    }

    log.warn(EXAMPLE_BANNER)
    log.info("BookingService starting on port 8080...")
    embeddedServer(CIO, port = 8080, host = "0.0.0.0") {
        module(config, httpClient)
    }.start(wait = true)
}

fun Application.module(config: Config, httpClient: HttpClient) {
    install(ContentNegotiation) { json() }

    val openapiSpec: String = Application::class.java.classLoader
        .getResourceAsStream("openapi.yaml")
        ?.bufferedReader()
        ?.use { it.readText() }
        ?: error("openapi.yaml not found on classpath")

    routing {
        get("/health") { call.respond(HttpStatusCode.OK) }
        get("/info") { call.respond(config) }
        get("/openapi") {
            call.respondText(openapiSpec, contentType = ContentType("text", "yaml"))
        }
        get("/booking/offers") { call.handleBookingOffers(httpClient, config) }
    }
}
