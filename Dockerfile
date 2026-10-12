FROM golang:1.25.4 AS builder
WORKDIR /app

# Install system dependencies for bun
RUN apt-get update && apt-get install -y unzip

# Install templ and bun for building assets
RUN go install github.com/a-h/templ/cmd/templ@v0.3.960
COPY --from=oven/bun:1 /usr/local/bin/bun ./bun
COPY --from=oven/bun:1 /usr/local/bin/bunx ./bunx
ENV PATH="/app:${PATH}"

# Copy dependency files
COPY go.mod go.sum package.json bun.lockb ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ files
RUN templ generate

# Build JavaScript and CSS
RUN ./bun install --frozen-lockfile
RUN ./bun run build:js
RUN ./bun run build:css

# Build Go binary
WORKDIR /app/cmd
RUN CGO_ENABLED=1 go build -o /app/app-binary


FROM golang:1.25.4 AS dev
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y unzip

# Install Go tools
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/a-h/templ/cmd/templ@v0.3.960

# Install bun for JavaScript bundling
COPY --from=oven/bun:1 /usr/local/bin/bun /usr/local/bin/bun
COPY --from=oven/bun:1 /usr/local/bin/bunx /usr/local/bin/bunx

CMD ["air"]


FROM debian:bookworm-slim AS runtime
WORKDIR /app
COPY --from=builder /app/app-binary .
COPY --from=builder /app/static ./static
ENTRYPOINT ["/app/app-binary"]


# A sporocarp: doppler on throwaway test data, for a mycelia spore to try a
# change end to end (built from spore/* branches only; see sporocarp/boot.sh).
# Garage runs in the same container, and the seed fills it and SQLite with
# sample posts. Nothing here touches a real secret or the real blog.
FROM builder AS sporocarp-tools
WORKDIR /app/cmd
RUN CGO_ENABLED=1 go build -o /app/seed-binary ./seed \
  && CGO_ENABLED=0 go build -o /app/garage-init ./garageinit

FROM runtime AS sporocarp
COPY --from=dxflrs/garage:v2.1.0 /garage /garage
COPY --from=sporocarp-tools /app/seed-binary /app/garage-init /app/
COPY sporocarp/garage.toml /etc/garage.toml
COPY sporocarp/boot.sh /app/sporocarp-boot.sh
RUN useradd --uid 1000 --no-create-home sporocarp \
  && mkdir -p /app/data /shared /var/lib/garage/meta /var/lib/garage/data \
  && chown -R 1000:1000 /app/data /shared /var/lib/garage
USER 1000
ENV DB_PATH=data/doppler.db \
    S3_ENDPOINT=127.0.0.1:3900 \
    S3_REGION=garage \
    S3_BUCKET=doppler-photos \
    S3_USE_SSL=false
ENTRYPOINT ["/app/sporocarp-boot.sh"]


# Default target: production
FROM runtime
