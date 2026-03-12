FROM golang:1.24-bookworm AS builder
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev && rm -rf /var/lib/apt/lists/*

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN CGO_ENABLED=1 GOOS=linux go build -o /service-desk ./cmd/server

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates libc6 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /service-desk /app/service-desk
COPY --from=builder /app/.env.example /app/.env.example
RUN mkdir -p /app/data /app/uploads

ENV PORT=8080
ENV DB_PATH=/app/data/mirea_service.db
EXPOSE 8080

CMD ["/app/service-desk"]
