package workshop.booking

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Config(
    val service: String,
    @SerialName("flightServiceUrl") val flightServiceUrl: String,
    @SerialName("hotelServiceUrl") val hotelServiceUrl: String,
    @SerialName("carServiceUrl") val carServiceUrl: String,
    val timeout: Int,
)

fun loadConfig(): Config = Config(
    service = "booking-1",
    flightServiceUrl = env("FLIGHT_SERVICE_URL", "http://localhost:8081"),
    hotelServiceUrl = env("HOTEL_SERVICE_URL", "http://localhost:8082"),
    carServiceUrl = env("CAR_SERVICE_URL", "http://localhost:8083"),
    timeout = 5000,
)

private fun env(name: String, default: String): String =
    System.getenv(name)?.takeIf { it.isNotBlank() } ?: default
