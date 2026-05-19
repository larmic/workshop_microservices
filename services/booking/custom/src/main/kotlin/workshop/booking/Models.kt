package workshop.booking

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement

@Serializable
data class HealthResponse(val status: String = "UP")

@Serializable
data class BookingOffers(
    val flights: JsonElement,
    val hotels: JsonElement,
    val cars: JsonElement,
)
