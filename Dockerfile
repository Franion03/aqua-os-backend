# ── Build ──────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/aquaos-backend ./cmd/server

# ── Runtime ────────────────────────────────────────────────────────
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/aquaos-backend .
EXPOSE 8080
ENTRYPOINT ["./aquaos-backend"]