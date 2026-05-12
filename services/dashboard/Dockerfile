FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /go/src/app/
COPY . /go/src/app/

RUN go mod download
RUN go test -v ./dashboard/...
RUN CGO_ENABLED=0 go build -a -o main ./dashboard

FROM docker:28-cli

LABEL org.opencontainers.image.title="Workshop Microservices - Dashboard"
LABEL org.opencontainers.image.description="Dashboard for the team neusta microservices workshop – orchestrates demo services via docker compose"
LABEL org.opencontainers.image.vendor="Team Neusta"
LABEL org.opencontainers.image.authors="Lars Michaelis <l.michaelis@neusta.de>"
LABEL org.opencontainers.image.url="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.source="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /root/
COPY --from=builder /go/src/app/main .

EXPOSE 8080
ENTRYPOINT ["./main"]
