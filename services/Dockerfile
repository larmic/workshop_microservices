FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /go/src/app/
COPY . /go/src/app/

ARG SERVICE_PATH

RUN go mod download
RUN go test -v ./${SERVICE_PATH}/...
RUN CGO_ENABLED=0 go build -a -o main ./${SERVICE_PATH}

FROM scratch

ARG SERVICE_DESC="Workshop Microservices"
ARG SERVICE_PORT=8080

LABEL org.opencontainers.image.title="${SERVICE_DESC}"
LABEL org.opencontainers.image.description="${SERVICE_DESC} – part of the team neusta microservices workshop"
LABEL org.opencontainers.image.vendor="Team Neusta"
LABEL org.opencontainers.image.authors="Lars Michaelis <l.michaelis@neusta.de>"
LABEL org.opencontainers.image.url="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.source="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.licenses="PolyForm-Noncommercial-1.0.0"

WORKDIR /root/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/main .

EXPOSE ${SERVICE_PORT}
ENTRYPOINT ["./main"]
