FROM golang:1.25.4 AS builder
WORKDIR /app

# Install system dependencies for bun
RUN apt-get update && apt-get install -y unzip

# Install templ and bun for building assets
RUN go install github.com/a-h/templ/cmd/templ@latest
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

# Download FontAwesome CSS
RUN curl -fsSL https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css -o static/css/all.min.css

# Build Go binary
WORKDIR /app/cmd
RUN CGO_ENABLED=1 go build -o /app/app-binary


FROM golang:1.25.4 AS dev
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y unzip

# Install Go tools
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install bun for JavaScript bundling
COPY --from=oven/bun:1 /usr/local/bin/bun /usr/local/bin/bun
COPY --from=oven/bun:1 /usr/local/bin/bunx /usr/local/bin/bunx

CMD ["air"]


# Default target for production
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app-binary .
COPY --from=builder /app/static ./static
ENTRYPOINT ["/app/app-binary"]
