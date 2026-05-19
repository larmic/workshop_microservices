package workshop.booking

import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.statement.bodyAsText
import io.ktor.http.HttpStatusCode
import io.ktor.http.isSuccess
import io.ktor.server.application.ApplicationCall
import io.ktor.server.request.httpMethod
import io.ktor.server.request.uri
import io.ktor.server.response.respond
import io.ktor.server.response.respondText
import kotlinx.coroutines.async
import kotlinx.coroutines.awaitAll
import kotlinx.coroutines.coroutineScope
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonElement
import org.slf4j.LoggerFactory

private val log = LoggerFactory.getLogger("BookingOffersHandler")

suspend fun ApplicationCall.handleBookingOffers(client: HttpClient, config: Config) {
    log.info("{} {} from {}", request.httpMethod.value, request.uri, request.local.remoteHost)

    val (flights, hotels, cars) = try {
        coroutineScope {
            val f = async { fetchJson(client, "${config.flightServiceUrl}/flights") }
            val h = async { fetchJson(client, "${config.hotelServiceUrl}/hotels") }
            val c = async { fetchJson(client, "${config.carServiceUrl}/cars") }
            awaitAll(f, h, c)
        }
    } catch (e: UpstreamException) {
        respondText(text = e.message ?: "upstream error", status = HttpStatusCode.InternalServerError)
        return
    }

    respond(BookingOffers(flights = flights, hotels = hotels, cars = cars))
}

private suspend fun fetchJson(client: HttpClient, url: String): JsonElement {
    val response = try {
        client.get(url)
    } catch (e: Exception) {
        throw UpstreamException("failed to fetch $url: ${e.message}")
    }
    val body = response.bodyAsText()
    if (!response.status.isSuccess()) {
        throw UpstreamException("upstream $url returned HTTP ${response.status.value}: ${body.trim()}")
    }
    return Json.parseToJsonElement(body)
}

private class UpstreamException(message: String) : RuntimeException(message)
