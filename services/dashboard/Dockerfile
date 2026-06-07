FROM --platform=$BUILDPLATFORM golang:1.26.4-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /go/src/app/
COPY . /go/src/app/

ARG TARGETOS
ARG TARGETARCH

RUN go mod download

# Discord-Webhook fürs Retro-Easteregg: wird nur beim CI-Build (docker.yml)
# aus dem Repository-Secret gesetzt. Ohne Wert bleibt der Platzhalter in
# index.html stehen und die Benachrichtigung unterbleibt (z.B. lokale Builds).
ARG DISCORD_EASTEREGG_WEBHOOK_URL=""
RUN if [ -n "$DISCORD_EASTEREGG_WEBHOOK_URL" ]; then \
        sed -i "s|__DISCORD_EASTEREGG_WEBHOOK_URL__|${DISCORD_EASTEREGG_WEBHOOK_URL}|" dashboard/static/index.html; \
    fi

RUN go test -v ./dashboard/...
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main ./dashboard

FROM docker:29-cli

LABEL org.opencontainers.image.title="Workshop Microservices - Dashboard"
LABEL org.opencontainers.image.description="Dashboard for the team neusta microservices workshop – orchestrates demo services via docker compose"
LABEL org.opencontainers.image.vendor="Team Neusta"
LABEL org.opencontainers.image.authors="Lars Michaelis <l.michaelis@neusta.de>"
LABEL org.opencontainers.image.url="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.source="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.licenses="PolyForm-Noncommercial-1.0.0"

WORKDIR /root/
COPY --from=builder /go/src/app/main .

EXPOSE 8080
ENTRYPOINT ["./main"]
