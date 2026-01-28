FROM golang:1.25.6-alpine AS builder
LABEL stage=intermediate

RUN apk add --no-cache ca-certificates

WORKDIR /go/src/app/
COPY go.mod /go/src/app/
COPY car/main.go /go/src/app/
COPY car/api /go/src/app/api

RUN go mod download
RUN go test -v ./...

ARG TARGETPLATFORM

RUN if [ "$TARGETPLATFORM" = "linux/arm/v7" ] ; then \
        echo "I am building linux/arm/v7 with CGO_ENABLED=0 GOARCH=arm GOARM=7" ; \
        env CGO_ENABLED=0 GOARCH=arm GOARM=7 go build -a -o main . ; \
        echo "Build done" ; \
    fi

RUN if [ "$TARGETPLATFORM" = "linux/arm64" ] ; then \
        echo "I am building linux/arm64 with CGO_ENABLED=0 GOARCH=arm64" ; \
        env CGO_ENABLED=0 GOARCH=arm64 go build -a -o main . ; \
        echo "Build done" ; \
    fi

RUN if [ "$TARGETPLATFORM" = "linux/amd64" ] ; then \
        echo "I am building linux/amd64 with CGO_ENABLED=0 GOARCH=amd64" ; \
        env CGO_ENABLED=0 GOARCH=amd64 go build -a -o main . ; \
        echo "Build done" ; \
    fi

FROM scratch
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/main .

EXPOSE 8083
ENTRYPOINT ["./main"]
