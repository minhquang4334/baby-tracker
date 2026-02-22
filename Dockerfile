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

# ── Stage 2: Build Go binary ─────────────────────────────────────────────────
FROM golang:1.22-bookworm AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy built static assets from previous stage
COPY --from=frontend-builder /app/static ./static

# Copy Go source
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o baby-care .

# ── Stage 3: Minimal runtime image ───────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=go-builder /app/baby-care .

# Data directory for SQLite (mount a Railway volume here)
RUN mkdir -p /data

ENV DB_PATH=/data/baby-care.db

EXPOSE 8080

CMD ["/app/baby-care", "--db", "/data/baby-care.db"]
