# Build stage
FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /campusworld-bedrock .

# Runtime
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /campusworld-bedrock /app/campusworld-bedrock
COPY config.toml /app/config.toml
EXPOSE 19132/udp
VOLUME ["/app/world", "/app/players"]
ENTRYPOINT ["/app/campusworld-bedrock"]
