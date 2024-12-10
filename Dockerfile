FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download


COPY . . 

WORKDIR /app/cmd
RUN go build -o /app/app-binary

# Goose is a SQLite Migrations handler
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN touch /app/doppler.db
RUN ls -lart
RUN goose -dir /app/internal/db/sql/ sqlite3 /app/doppler.db up


FROM debian:bookworm-slim AS runner
WORKDIR /app
COPY --from=builder /app/app-binary .
COPY --from=builder /app/doppler.db .
COPY --from=builder /app/static ./static
ENTRYPOINT ["/app/app-binary"]
