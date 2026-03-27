FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /go/src/app/
COPY . /go/src/app/

RUN go mod download
RUN go test -v ./dashboard/...
RUN CGO_ENABLED=0 go build -a -o main ./dashboard

FROM docker:28-cli

LABEL org.opencontainers.image.source="https://github.com/larmic/workshop-microservices"
LABEL org.opencontainers.image.description="Workshop Microservices - Dashboard"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /root/
COPY --from=builder /go/src/app/main .

EXPOSE 8080
ENTRYPOINT ["./main"]
