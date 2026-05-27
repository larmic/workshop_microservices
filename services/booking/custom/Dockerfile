# --- Build stage ---
FROM eclipse-temurin:25.0.3_9-jdk-alpine AS builder

WORKDIR /app

COPY gradlew settings.gradle.kts build.gradle.kts gradle.properties ./
COPY gradle ./gradle
RUN sed -i 's/\r$//' gradlew && chmod +x gradlew
RUN ./gradlew --no-daemon dependencies > /dev/null 2>&1 || true

COPY src ./src
RUN ./gradlew --no-daemon shadowJar

# --- Runtime stage ---
FROM eclipse-temurin:25.0.3_9-jre-alpine

ARG SERVICE_DESC="Booking Custom (Kotlin/Ktor)"

LABEL org.opencontainers.image.title="${SERVICE_DESC}"
LABEL org.opencontainers.image.description="${SERVICE_DESC} – part of the team neusta microservices workshop"
LABEL org.opencontainers.image.vendor="Team Neusta"
LABEL org.opencontainers.image.authors="Lars Michaelis <l.michaelis@neusta.de>"
LABEL org.opencontainers.image.licenses="PolyForm-Noncommercial-1.0.0"

WORKDIR /app
COPY --from=builder /app/build/libs/booking-custom-all.jar app.jar

EXPOSE 8080
ENTRYPOINT ["java", "-XX:+UseSerialGC", "-Xss512k", "-jar", "app.jar"]
