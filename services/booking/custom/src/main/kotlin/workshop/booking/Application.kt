package workshop.booking

import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO as ClientCIO
import io.ktor.client.plugins.HttpTimeout
import io.ktor.http.ContentType
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

fun main() {
    val config = loadConfig()

    val httpClient = HttpClient(ClientCIO) {
        install(HttpTimeout) {
            requestTimeoutMillis = config.timeout.toLong()
            connectTimeoutMillis = config.timeout.toLong()
            socketTimeoutMillis = config.timeout.toLong()
        }
        expectSuccess = false
    }

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
        get("/health") { call.respond(HealthResponse(status = "UP")) }
        get("/info") { call.respond(config) }
        get("/openapi") {
            call.respondText(openapiSpec, contentType = ContentType("text", "yaml"))
        }
        get("/booking/offers") { call.handleBookingOffers(httpClient, config) }
    }
}
