# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder
WORKDIR /build
COPY cmd/garageinit/main.go .
RUN CGO_ENABLED=0 go build -o garage-init main.go

FROM dxflrs/garage:v2.1.0
COPY --from=builder /build/garage-init /garage-init
ENTRYPOINT ["/garage-init"]
