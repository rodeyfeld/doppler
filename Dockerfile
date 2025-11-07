FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . . 

WORKDIR /app/cmd
RUN go build -o /app/app-binary


FROM debian:bookworm-slim AS runner
WORKDIR /app
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/app-binary .
COPY --from=builder /app/static ./static
ENTRYPOINT ["/app/app-binary"]
