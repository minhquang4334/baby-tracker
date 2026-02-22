# ── Stage 1: Build frontend ──────────────────────────────────────────────────
FROM node:20-slim AS frontend-builder

WORKDIR /app

COPY frontend/package*.json frontend/
RUN cd frontend && npm ci

COPY frontend/ frontend/
COPY static/.gitkeep static/

RUN cd frontend && npx esbuild src/main.ts \
      --bundle \
      --outfile=../static/app.js \
      --minify \
      --loader:.css=css \
      --external:*.ttf \
      --external:*.woff \
      --external:*.woff2 && \
    cp index.html ../static/index.html

# ── Stage 2: Build Go binary (pure Go — no CGO, no gcc) ──────────────────────
FROM golang:1.26-alpine AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --from=frontend-builder /app/static ./static
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o baby-care .

# ── Stage 3: Minimal runtime image ───────────────────────────────────────────
FROM alpine:3

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=go-builder /app/baby-care .

# Data directory for SQLite (mount a Railway volume here)
RUN mkdir -p /data

EXPOSE 8080

CMD ["/app/baby-care", "--db", "/data/baby-care.db"]
