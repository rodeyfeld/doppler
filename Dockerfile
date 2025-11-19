FROM golang:1.25.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . . 

WORKDIR /app/cmd
RUN go build -o /app/app-binary


FROM golang:1.25.4 AS dev
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y unzip

# Install Go tools
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install bun for JavaScript bundling
RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/root/.bun/bin:${PATH}"

CMD ["air"]


# Default target for production
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app-binary .
COPY --from=builder /app/static ./static
ENTRYPOINT ["/app/app-binary"]
