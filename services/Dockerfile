FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /go/src/app/
COPY go.mod /go/src/app/

ARG SERVICE_PATH
COPY ${SERVICE_PATH}/main.go /go/src/app/
COPY ${SERVICE_PATH}/api /go/src/app/api

RUN go mod download
RUN go test -v ./...
RUN CGO_ENABLED=0 go build -a -o main .

FROM scratch

ARG SERVICE_DESC="Workshop Microservices"
ARG SERVICE_PORT=8080

LABEL org.opencontainers.image.source="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.description="${SERVICE_DESC}"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /root/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/main .

EXPOSE ${SERVICE_PORT}
ENTRYPOINT ["./main"]
