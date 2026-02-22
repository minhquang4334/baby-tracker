# Baby Care Tracker

A fast, minimal baby care tracking web app built for parents who want to log sleep, feeding, diapers, and growth with as few taps as possible. Runs as a single Go binary with SQLite — no cloud, no accounts, no subscriptions.

## Features

- **2-tap diaper logging** — FAB → Wet / Dirty / Mixed, saved instantly with current time
- **Sleep timer** — start with 2 taps, stop with 1 tap on the live banner
- **Feeding tracker** — left/right breast timers, or bottle with quantity (pre-filled from last use)
- **Growth chart** — Canvas-rendered weight & length chart against WHO standards
- **Vietnamese baby guide** — feeding schedule, WHO growth tables, milestones, vaccination calendar, diaper sizes, and warning signs all in the Guide tab
- **Smart auto-stop** — starting a breast feed automatically stops active sleep, and vice versa
- **Day history** — color-coded timeline with date navigator and filter pills (All / Sleep / Feeding / Diaper)
- **GMT+7 native** — all times stored and displayed in Ho Chi Minh City timezone
- **Offline-capable** — single binary + SQLite file, works without internet after first load
- **Mobile-first UI** — 480px max-width, 44px+ touch targets, Inter font, warm color palette

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.26, `net/http` (native method+path routing) |
| Database | SQLite via `go-sqlite3` (WAL mode, FK enforcement) |
| Frontend | Vanilla TypeScript, no framework |
| Bundler | esbuild (outputs a single `app.js` + inlined CSS) |
| Deployment | Single Go binary with `embed.FS` static assets |

## Screens

| Screen | Route | Description |
|--------|-------|-------------|
| Onboarding | `#/onboarding` | First-visit form: baby name, date of birth, gender |
| Dashboard | `#/dashboard` | Today's summary cards, active timer banners, recent activity |
| History | `#/history` | Date navigator, filter pills, color-coded vertical timeline |
| Growth | `#/growth` | Canvas line chart (weight/length toggle), data table |
| Guide | `#/guide` | Vietnamese baby development reference (WHO tables, milestones, vaccinations) |

## Project Structure

```
.
├── main.go                        # Entry point: flags, store init, server start
├── embed.go                       # //go:embed static/*
├── go.mod / go.sum
├── Makefile
├── Dockerfile                     # Multi-stage build (Node → Go/gcc → debian-slim)
├── railway.toml                   # Railway deploy config
├── .github/workflows/deploy.yml   # CI/CD pipeline
│
├── internal/
│   ├── server/server.go           # Mux routes + SPA fallback
│   ├── handler/                   # HTTP handlers (JSON in/out)
│   │   ├── handler.go             # Shared helpers (writeJSON, writeError)
│   │   ├── children.go
│   │   ├── sleep.go
│   │   ├── feeding.go
│   │   ├── diaper.go
│   │   ├── growth.go
│   │   └── summary.go
│   ├── middleware/middleware.go    # Logger, CORS
│   ├── model/                     # Go structs matching DB tables
│   └── store/                     # SQLite queries (one file per domain)
│       ├── store.go               # Open, migrations, GMT+7 timezone helpers
│       ├── children.go
│       ├── sleep.go               # Auto-stops active feeding on sleep start
│       ├── feeding.go             # Auto-stops active sleep on breast feed start
│       ├── diaper.go
│       ├── growth.go
│       └── summary.go
│
├── frontend/
│   ├── index.html                 # SPA shell
│   ├── package.json               # devDeps: esbuild, typescript
│   ├── tsconfig.json
│   ├── src/
│   │   ├── main.ts                # Bootstrap: API check → route to onboarding or dashboard
│   │   ├── router.ts              # Hash-based SPA router
│   │   ├── api.ts                 # Typed fetch wrapper
│   │   ├── state.ts               # Pub/sub signals (child, activeSleep, activeFeeding)
│   │   ├── types/models.ts        # TypeScript interfaces matching Go models
│   │   ├── utils/
│   │   │   ├── date.ts            # GMT+7 formatters, nowISO(), localInputToISO()
│   │   │   └── dom.ts             # h() element factory, $(), $$()
│   │   └── components/
│   │       ├── nav.ts             # Bottom tab bar (Home, History, Growth, Guide)
│   │       ├── onboarding.ts
│   │       ├── dashboard.ts
│   │       ├── quick-add.ts       # FAB + radial menu
│   │       ├── sleep-modal.ts
│   │       ├── feeding-modal.ts
│   │       ├── diaper-modal.ts
│   │       ├── growth-modal.ts
│   │       ├── growth-chart.ts
│   │       ├── history.ts
│   │       ├── guide.ts
│   │       └── toast.ts
│   └── css/
│       ├── main.css               # Design tokens (CSS custom properties), reset
│       ├── components.css
│       ├── dashboard.css
│       ├── modal.css
│       ├── history.css
│       ├── growth-chart.css
│       └── onboarding.css
│
└── static/                        # esbuild output — embedded into Go binary at compile time
    ├── index.html
    ├── app.js
    └── app.css
```

## Local Development

### Prerequisites

- Go 1.26+
- Node.js 20+
- gcc (required for CGO / sqlite3)
  - macOS: `xcode-select --install`
  - Ubuntu/Debian: `sudo apt-get install build-essential`

### Setup

```bash
# Install all dependencies (npm + go mod)
make deps
```

### Development (hot reload)

Open two terminals:

```bash
# Terminal 1 — watch TypeScript and rebuild static/app.js on save
make frontend-watch

# Terminal 2 — run the Go server (serves the embedded static files)
make run
```

Open http://localhost:8080 in your browser. After editing TypeScript files, refresh the browser. After editing Go files, restart `make run`.

> **Note:** Static assets are embedded at Go compile time. The `make run` target uses `go run .` which recompiles on each start, so changes to TypeScript are picked up after `esbuild` rebuilds and you restart the Go server.

### Production build (single binary)

```bash
make all          # runs: deps → frontend → backend
./baby-care       # serves on :8080, DB at ~/.baby-care/data.db
```

### CLI flags

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` (or `$PORT` env var) | HTTP listen port |
| `--db` | `~/.baby-care/data.db` | SQLite database file path |

```bash
./baby-care --port 3000 --db /var/data/baby.db
```

## API Reference

Base path: `/api/v1`
All requests and responses are JSON. Timestamps are RFC3339 with `+07:00` offset.

### Child

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/child` | Create child profile (first-visit onboarding) |
| `GET` | `/child` | Get the child profile |
| `PUT` | `/child` | Update child profile |

### Sleep

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/sleep` | Start a sleep (sets `end_time = null`). Also stops any active feeding. |
| `GET` | `/sleep` | List sleep logs (supports `?date=YYYY-MM-DD`) |
| `GET` | `/sleep/active` | Get in-progress sleep (no `end_time`) |
| `PUT` | `/sleep/{logId}` | Update sleep (stop: set `end_time`) |
| `DELETE` | `/sleep/{logId}` | Delete sleep log |

`POST /sleep` response includes a `stopped_feeding` field when a breast feed was auto-stopped.

### Feeding

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/feeding` | Start/log a feeding. Breast feeds auto-stop active sleep. |
| `GET` | `/feeding` | List feeding logs (supports `?date=YYYY-MM-DD`) |
| `GET` | `/feeding/active` | Get in-progress breast feed |
| `PUT` | `/feeding/{logId}` | Update feeding (stop breast feed, edit bottle) |
| `DELETE` | `/feeding/{logId}` | Delete feeding log |

`POST /feeding` response includes a `stopped_sleep` field when sleep was auto-stopped.

Feed types: `breast_left`, `breast_right`, `bottle`

### Diaper

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/diaper` | Log a diaper change |
| `GET` | `/diaper` | List diaper logs (supports `?date=YYYY-MM-DD`) |
| `PUT` | `/diaper/{logId}` | Update diaper log |
| `DELETE` | `/diaper/{logId}` | Delete diaper log |

Diaper types: `wet`, `dirty`, `mixed`

### Growth

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/growth` | Log a growth measurement |
| `GET` | `/growth` | List all growth logs |
| `PUT` | `/growth/{logId}` | Update growth log |
| `DELETE` | `/growth/{logId}` | Delete growth log |

### Summary

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/summary` | Aggregated day stats (`?date=YYYY-MM-DD`, defaults to today GMT+7) |

Summary response includes total sleep hours, feeding count + breakdown, diaper count, and latest growth measurement.

### Health check

```
GET /health  →  200 OK
```

## Database Schema

All IDs are UUID v4. All timestamps are RFC3339 strings with `+07:00` offset. Migrations run as `CREATE TABLE IF NOT EXISTS` on startup.

```sql
CREATE TABLE children (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  date_of_birth TEXT NOT NULL,
  gender TEXT NOT NULL CHECK(gender IN ('male','female','other')),
  photo_url TEXT,
  notes TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE sleep_logs (
  id TEXT PRIMARY KEY,
  child_id TEXT NOT NULL REFERENCES children(id),
  start_time TEXT NOT NULL,
  end_time TEXT,                   -- NULL = currently sleeping
  duration_minutes INTEGER,
  notes TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE feeding_logs (
  id TEXT PRIMARY KEY,
  child_id TEXT NOT NULL REFERENCES children(id),
  feed_type TEXT NOT NULL CHECK(feed_type IN ('breast_left','breast_right','bottle')),
  start_time TEXT NOT NULL,
  end_time TEXT,                   -- NULL = currently feeding (breast only)
  duration_minutes INTEGER,
  quantity_ml INTEGER,             -- bottle only
  notes TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE diaper_logs (
  id TEXT PRIMARY KEY,
  child_id TEXT NOT NULL REFERENCES children(id),
  diaper_type TEXT NOT NULL CHECK(diaper_type IN ('wet','dirty','mixed')),
  changed_at TEXT NOT NULL,
  notes TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE growth_logs (
  id TEXT PRIMARY KEY,
  child_id TEXT NOT NULL REFERENCES children(id),
  measured_on TEXT NOT NULL,
  weight_grams INTEGER,
  length_mm INTEGER,
  head_circumference_mm INTEGER,
  notes TEXT,
  created_at TEXT NOT NULL
);
```

## Deployment (Railway)

The app ships as a Docker container. Railway builds the image from the `Dockerfile` in the repo root.

### One-time Railway setup

1. Create a new Railway project → **Deploy from GitHub repo** → select your repo
2. Railway detects the `Dockerfile` and starts the first build automatically
3. **Add a persistent volume** so the SQLite database survives redeploys:
   - Service → **Volumes** → Add Volume → Mount Path: `/data`
4. The app reads `$PORT` automatically (set by Railway) — no extra config needed

### CI/CD (GitHub Actions)

The workflow in `.github/workflows/deploy.yml` runs on every push:

```
push to main
  └─ build job
       ├─ npm ci + esbuild (TypeScript → app.js)
       ├─ go vet ./...
       └─ CGO_ENABLED=1 go build
  └─ deploy job (main branch only, after build passes)
       └─ POST $RAILWAY_DEPLOY_HOOK
```

**To wire up the deploy hook:**

1. In Railway: Service → Settings → **Deploy Hook** → copy the URL
2. In GitHub repo: Settings → Secrets and variables → Actions → **New secret**
   - Name: `RAILWAY_DEPLOY_HOOK`
   - Value: the URL from step 1

After this, every push to `main` runs CI first. If the build and vet pass, GitHub Actions calls the Railway deploy hook to trigger a redeploy. Pull requests run the build check only (no deploy).

### Dockerfile overview

The image uses a 3-stage build to keep the final image small:

| Stage | Base | Purpose |
|-------|------|---------|
| `frontend-builder` | `node:20-slim` | `npm ci` + esbuild bundle |
| `go-builder` | `golang:1.26-bookworm` | `go build` with CGO (gcc included) |
| runtime | `debian:bookworm-slim` | Final image — just the binary + libc |

## Design System

| Token | Value | Used for |
|-------|-------|----------|
| `--color-primary` | `#E8507A` | Action buttons, active states |
| `--color-sleep` | `#8B5CF6` | Sleep cards, timeline entries |
| `--color-feeding` | `#EC4899` | Feeding cards, timeline entries |
| `--color-diaper` | `#10B981` | Diaper cards, timeline entries |
| `--color-growth` | `#F59E0B` | Growth cards, timeline entries |
| `--color-bg` | `#FEFCFB` | App background (warm off-white) |
| Font | Inter | All text |
| Touch targets | 44×44px min | All interactive elements |
| Container | 480px max-width | Centered on desktop (phone-like) |

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make all` | Full build: deps → frontend → backend binary |
| `make deps` | `npm install` + `go mod download` |
| `make frontend` | One-shot esbuild bundle to `static/` |
| `make frontend-watch` | esbuild in watch mode (dev) |
| `make backend` | `CGO_ENABLED=1 go build -o baby-care .` |
| `make run` | `go run . --port 8080` |
| `make clean` | Remove binary and built static assets |

## License

MIT
