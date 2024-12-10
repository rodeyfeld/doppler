FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /app/cmd
RUN go build -o /app/app-binary

FROM debian:bookworm-slim AS runner
WORKDIR /app
COPY --from=builder /app/app-binary .
ENTRYPOINT ["/app/app-binary"]
