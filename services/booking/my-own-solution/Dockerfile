# --- Build stage ---
FROM eclipse-temurin:21-jdk-alpine AS builder

WORKDIR /app

COPY gradlew settings.gradle.kts build.gradle.kts gradle.properties ./
COPY gradle ./gradle
RUN ./gradlew --no-daemon dependencies > /dev/null 2>&1 || true

COPY src ./src
RUN ./gradlew --no-daemon shadowJar

# --- Runtime stage ---
FROM eclipse-temurin:21-jre-alpine

ARG SERVICE_DESC="Booking My-Own-Solution (Kotlin/Ktor)"

LABEL org.opencontainers.image.title="${SERVICE_DESC}"
LABEL org.opencontainers.image.description="${SERVICE_DESC} – part of the team neusta microservices workshop"
LABEL org.opencontainers.image.vendor="Team Neusta"
LABEL org.opencontainers.image.authors="Lars Michaelis <l.michaelis@neusta.de>"
LABEL org.opencontainers.image.licenses="PolyForm-Noncommercial-1.0.0"

WORKDIR /app
COPY --from=builder /app/build/libs/booking-my-own-solution-all.jar app.jar

EXPOSE 8080
ENTRYPOINT ["java", "-XX:+UseSerialGC", "-Xss512k", "-jar", "app.jar"]
